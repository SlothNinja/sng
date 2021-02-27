package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/atf"
	"github.com/SlothNinja/confucius"
	"github.com/SlothNinja/cookie"
	"github.com/SlothNinja/game"
	"github.com/SlothNinja/got"
	"github.com/SlothNinja/indonesia"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/mlog"
	"github.com/SlothNinja/rating"
	"github.com/SlothNinja/restful"
	"github.com/SlothNinja/sn"
	"github.com/SlothNinja/tammany"
	gtype "github.com/SlothNinja/type"
	"github.com/SlothNinja/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const (
	// Environment Variables
	NODE_ENV           = "NODE_ENV"
	DS_PROJECT_ID      = "DS_PROJECT_ID"
	USER_PROJECT_ID    = "USER_PROJECT_ID"
	COOKIE_URL         = "COOKIE_URL"
	LOGIN_HOST         = "LOGIN_HOST"
	DS_HOST            = "DS_HOST"
	DS_USER_HOST       = "DS_USER_HOST"
	googleCloudProject = "GOOGLE_CLOUD_PROJECT"

	production     = "production"
	userPrefix     = "user"
	gamesPrefix    = "games"
	ratingPrefix   = "rating"
	mailPrefix     = "mail"
	rootPath       = "/"
	hashKeyLength  = 64
	blockKeyLength = 32
	sessionName    = "sng-oauth"
	msgEnter       = "Entering"
	msgExit        = "Exiting"
)

type Client struct {
	*sn.Client
	User   *user.Client
	Game   *game.Client
	MLog   *mlog.Client
	Rating *rating.Client
}

func NewClient(dClient *datastore.Client, uClient *user.Client, logger *log.Logger, cache *cache.Cache, router *gin.Engine) *Client {
	client := &Client{
		Client: sn.NewClient(dClient, logger, cache, router),
		User:   uClient,
		MLog:   mlog.NewClient(dClient, uClient, logger, cache),
		Rating: rating.NewClient(dClient, uClient, logger, cache, router, "rating"),
	}
	return client.staticRoutes()
}

func main() {
	setGinMode()

	logClient := newLogClient()
	defer logClient.Close()

	logger := logClient.Logger("sng")
	cache := cache.New(30*time.Minute, 10*time.Minute)
	store, err := cookie.NewClient(logger, cache).NewStore()

	if err != nil {
		logger.Panicf("unable create cookie store: %v", err)
	}

	router := gin.Default()
	renderer := restful.ParseTemplates("templates/", ".tmpl")
	router.HTMLRender = renderer

	router.Use(
		sessions.Sessions(sessionName, store),
		restful.AddTemplates(renderer.Templates),
	)

	userClient := user.NewClient(logger, cache)
	dsClient := newDSClient(logger)
	client := NewClient(dsClient, userClient, logger, cache, router)

	// Welcome Page (index.html) route
	// welcome.NewClient(dsClient, userClient, logger, cache, router)

	const afterLoad = true
	// Game routes
	client.Game = game.NewClient(dsClient, userClient, logger, cache, router, "games", afterLoad)

	// After The Flood
	atf.NewClient(dsClient, userClient, client.Game, client.MLog, client.Rating,
		logger, cache, router, gtype.ATF)

	// Guild of Thieves
	got.NewClient(dsClient, userClient, client.MLog, client.Rating, logger, cache, router, gtype.GOT)

	// Tammany Hall
	tammany.NewClient(dsClient, userClient, client.Game, client.MLog, client.Rating,
		logger, cache, router, gtype.Tammany)

	// Indonesia
	indonesia.NewClient(dsClient, userClient, client.Game, client.MLog, client.Rating,
		logger, cache, router, gtype.Indonesia)

	// Confucius
	confucius.NewClient(dsClient, userClient, client.Game, client.MLog, client.Rating,
		logger, cache, router, gtype.Confucius)

	// warmup
	router.GET("_ah/warmup", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// login
	router.GET("login", login)

	// logout
	router.GET("logout", logout)

	// home
	router.GET("home", client.homeHandler)

	router.Run()
}

func setGinMode() {
	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		return
	}
	gin.SetMode(gin.DebugMode)
	return
}

func newDSClient(log *log.Logger) *datastore.Client {
	client, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		log.Panicf("unable to create datastore client: %v", err)
	}
	return client
}

func getLoginHost() string {
	return os.Getenv(LOGIN_HOST)
}

// staticHandler for local development since app.yaml is ignored
// static files are handled via app.yaml routes when deployed
func (client *Client) staticRoutes() *Client {
	if sn.IsProduction() {
		return client
	}
	// client.Router.StaticFile("/favicon.ico", "public/favicon.ico")
	client.Router.Static("/images", "public/images")
	client.Router.Static("/javascripts", "public/javascripts")
	client.Router.Static("/jsold", "public/js")
	client.Router.Static("/stylesheets", "public/stylesheets")
	client.Router.Static("/rules", "public/rules")
	client.Router.StaticFile("/", "dist/index.html")
	client.Router.StaticFile("/app.js", "dist/app.js")
	client.Router.StaticFile("/favicon.ico", "dist/favicon.ico")
	client.Router.Static("/img", "dist/img")
	client.Router.Static("/js", "dist/js")
	client.Router.Static("/css", "dist/css")
	return client
}

func login(c *gin.Context) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	referer := c.Request.Referer()
	encodedReferer := base64.StdEncoding.EncodeToString([]byte(referer))

	c.Redirect(http.StatusSeeOther, getLoginHost()+"/login?redirect="+encodedReferer)
}

func logout(c *gin.Context) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	referer := c.Request.Referer()
	encodedReferer := base64.StdEncoding.EncodeToString([]byte(referer))

	c.Redirect(http.StatusSeeOther, getLoginHost()+"/logout?redirect="+encodedReferer)
}

func getProjectID() string {
	return os.Getenv(googleCloudProject)
}

func newLogClient() *log.Client {
	client, err := log.NewClient(getProjectID())
	if err != nil {
		log.Panicf("unable to create logging client: %v", err)
	}
	return client
}

func (cl *Client) homeHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"cu": cu})
}
