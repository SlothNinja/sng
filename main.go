package main

import (
	"context"
	"encoding/base64"
	"math/rand"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"github.com/SlothNinja/atf"
	"github.com/SlothNinja/confucius"
	"github.com/SlothNinja/cookie"
	"github.com/SlothNinja/game"
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
	COOKIE_URL       = "COOKIE_URL"
	LOGIN_HOST       = "LOGIN_HOST"
	DS_HOST          = "DS_HOST"
	DS_USER_HOST     = "DS_USER_HOST"
	NodeEnv          = "NODE_ENV"
	SNGProjectIDEnv  = "SNG_PROJECT_ID"
	SNGDSURLEnv      = "SNG_DS_URL"
	SNGHostURLEnv    = "SNG_HOST_URL"
	UserProjectIDEnv = "USER_PROJECT_ID"
	UserDSURLEnv     = "USER_DS_URL"
	UserHostURLEnv   = "USER_HOST_URL"
	GotProjectIDEnv  = "GOT_PROJECT_ID"
	GotDSURLEnv      = "GOT_DS_URL"
	GotHostURLEnv    = "GOT_HOST_URL"

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
	Got    *game.Client
	User   *user.Client
	Game   *game.Client
	MLog   *mlog.Client
	Rating *rating.Client
}

func NewClient(ctx context.Context) *Client {
	logClient := newLogClient()
	snClient := sn.NewClient(ctx, sn.Options{
		ProjectID: getSNGProjectID(),
		DSURL:     getSNGDSURL(),
		Logger:    logClient.Logger("sng"),
		Cache:     cache.New(30*time.Minute, 10*time.Minute),
		Router:    gin.New(),
	})

	gotClient := sn.NewClient(ctx, sn.Options{
		ProjectID: getGotProjectID(),
		DSURL:     getGotDSURL(),
		Logger:    snClient.Log,
		Cache:     snClient.Cache,
		Router:    snClient.Router,
	})

	uClient := user.NewClient(sn.NewClient(ctx, sn.Options{
		ProjectID: getUserProjectID(),
		DSURL:     getUserDSURL(),
		Logger:    snClient.Log,
		Cache:     snClient.Cache,
		Router:    snClient.Router,
	}))

	store, err := cookie.NewClient(uClient.Client).NewStore(ctx)
	if err != nil {
		snClient.Log.Panicf("unable create cookie store: %v", err)
	}

	renderer := restful.ParseTemplates("templates/", ".tmpl")
	snClient.Router.HTMLRender = renderer

	snClient.Router.Use(
		sessions.Sessions(sessionName, store),
		restful.AddTemplates(renderer.Templates),
		gin.LoggerWithWriter(snClient.Log.StandardLogger(logging.Debug).Writer()),
		gin.RecoveryWithWriter(snClient.Log.StandardLogger(logging.Critical).Writer()),
	)

	const afterLoad = true
	nClient := &Client{
		Client: snClient,
		User:   uClient,
		MLog:   mlog.NewClient(snClient, uClient),
		Rating: rating.NewClient(snClient, uClient, "rating"),
		Game:   game.NewClient(snClient, uClient, "games", afterLoad),
		Got:    game.NewClient(gotClient, uClient, "gotgames", !afterLoad),
	}

	// After The Flood
	atf.NewClient(snClient, uClient, nClient.Game, nClient.Rating, gtype.ATF)

	// Tammany Hall
	tammany.NewClient(snClient, uClient, nClient.Game, nClient.Rating, gtype.Tammany)

	// Indonesia
	indonesia.NewClient(snClient, uClient, nClient.Game, nClient.Rating, gtype.Indonesia)

	// Confucius
	confucius.NewClient(snClient, uClient, nClient.Game, nClient.Rating, gtype.Confucius)
	return nClient.addRoutes()
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		cl := NewClient(ctx)
		defer cl.Close()
		cl.Router.Run()
	} else {
		gin.SetMode(gin.DebugMode)
		cl := NewClient(ctx)
		defer cl.Close()
		cl.Router.RunTLS(getPort(), "cert.pem", "key.pem")
	}
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

// staticHandler for local development since app.yaml is ignored
// static files are handled via app.yaml routes when deployed
func (client *Client) addRoutes() *Client {
	if !sn.IsProduction() {
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
	}

	// warmup
	client.Router.GET("_ah/warmup", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// login
	client.Router.GET("login", client.login)

	// logout
	client.Router.GET("logout", client.logout)

	// home
	client.Router.GET("home", client.homeHandler)

	client.Router.POST("games", client.myGamesHandler)

	return client
}

func (cl *Client) login(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	referer := c.Request.Referer()
	encodedReferer := base64.StdEncoding.EncodeToString([]byte(referer))

	path := getUserHostURL() + "/login?redirect=" + encodedReferer
	c.Redirect(http.StatusSeeOther, path)
}

func (cl *Client) logout(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	user.Logout(c)
	c.Redirect(http.StatusSeeOther, "/")
}

func newLogClient() *log.Client {
	client, err := log.NewClient(getSNGProjectID())
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

func getSNGProjectID() string {
	return os.Getenv(SNGProjectIDEnv)
}

func getSNGHostURL() string {
	return os.Getenv(SNGHostURLEnv)
}

func getSNGDSURL() string {
	return os.Getenv(SNGDSURLEnv)
}

func getUserProjectID() string {
	return os.Getenv(UserProjectIDEnv)
}

func getUserDSURL() string {
	return os.Getenv(UserDSURLEnv)
}

func getUserHostURL() string {
	return os.Getenv(UserHostURLEnv)
}

func getGotProjectID() string {
	return os.Getenv(GotProjectIDEnv)
}

func getGotDSURL() string {
	return os.Getenv(GotDSURLEnv)
}

func getGotHostURL() string {
	return os.Getenv(GotHostURLEnv)
}

func getPort() string {
	return ":" + os.Getenv("PORT")
}
