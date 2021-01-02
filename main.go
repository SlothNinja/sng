package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/atf"
	"github.com/SlothNinja/confucius"
	"github.com/SlothNinja/game"
	"github.com/SlothNinja/got"
	"github.com/SlothNinja/indonesia"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/rating"
	"github.com/SlothNinja/restful"
	"github.com/SlothNinja/sn"
	"github.com/SlothNinja/tammany"
	gtype "github.com/SlothNinja/type"
	"github.com/SlothNinja/user"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	// ucon "github.com/SlothNinja/user-controller"
	"github.com/SlothNinja/welcome"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/patrickmn/go-cache"
)

const (
	// Environment Variables
	NODE_ENV        = "NODE_ENV"
	DS_PROJECT_ID   = "DS_PROJECT_ID"
	USER_PROJECT_ID = "USER_PROJECT_ID"
	COOKIE_URL      = "COOKIE_URL"
	LOGIN_HOST      = "LOGIN_HOST"
	DS_HOST         = "DS_HOST"
	DS_USER_HOST    = "DS_USER_HOST"

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

func main() {
	setGinMode()

	dsClient := getDSClient()
	// db, err := datastore.NewClient(context.Background(), "")
	// if err != nil {
	// 	panic(fmt.Sprintf("unable to connect to database: %v", err.Error()))
	// }

	mcache := cache.New(30*time.Minute, 10*time.Minute)
	userClient := getUserDSClient(mcache)

	s, err := getSecrets()
	if err != nil {
		panic(err.Error())
	}

	store := createCookieStore(s)
	r := gin.Default()
	renderer := restful.ParseTemplates("templates/", ".tmpl")
	r.HTMLRender = renderer

	r.Use(
		sessions.Sessions(sessionName, store),
		restful.AddTemplates(renderer.Templates),
		// user.GetCUserHandler(db),
	)

	// Welcome Page (index.html) route
	r = welcome.NewClient(userClient).AddRoutes(r)

	// Games Routes
	r = game.NewClient(userClient, dsClient).AddRoutes(gamesPrefix, r)

	// User Routes
	// r = ucon.NewClient(db).AddRoutes(userPrefix, r)

	// Rating Routes
	r = rating.NewClient(userClient, dsClient).AddRoutes(ratingPrefix, r)

	// After The Flood
	r = atf.NewClient(dsClient, userClient, mcache).Register(gtype.ATF, r)

	// Guild of Thieves
	r = got.NewClient(dsClient, userClient, mcache).Register(gtype.GOT, r)

	// Tammany Hall
	r = tammany.NewClient(dsClient, userClient, mcache).Register(gtype.Tammany, r)

	// Indonesia
	r = indonesia.NewClient(dsClient, userClient, mcache).Register(gtype.Indonesia, r)

	// Confucius
	r = confucius.NewClient(dsClient, userClient, mcache).Register(gtype.Confucius, r)

	// warmup
	r.GET("_ah/warmup", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// login
	r.GET("login", login)

	// logout
	r.GET("logout", logout)

	r = staticRoutes(r)

	r.Run()
}

type secrets struct {
	HashKey   []byte
	BlockKey  []byte
	UpdatedAt time.Time
	Key       *datastore.Key `datastore:"__key__"`
}

func setGinMode() {
	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		return
	}
	gin.SetMode(gin.DebugMode)
	return
}

func getUserDSClient(mcache *cache.Cache) user.Client {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	if sn.IsProduction() {
		log.Debugf("production")
		dsClient, err := datastore.NewClient(
			context.Background(),
			os.Getenv(USER_PROJECT_ID),
		)
		if err != nil {
			panic(fmt.Sprintf("unable to connect to user database: %v", err.Error()))
		}
		return user.NewClient(dsClient, mcache)

	}
	log.Debugf("development")
	dsClient, err := datastore.NewClient(
		context.Background(),
		os.Getenv(USER_PROJECT_ID),
		option.WithEndpoint(os.Getenv(DS_USER_HOST)),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithGRPCConnectionPool(50),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to connect to user database: %v", err.Error()))
	}
	return user.NewClient(dsClient, mcache)
}

func getDSClient() *datastore.Client {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	if sn.IsProduction() {
		log.Debugf("production")
		dsClient, err := datastore.NewClient(
			context.Background(),
			os.Getenv(DS_PROJECT_ID),
		)
		if err != nil {
			panic(fmt.Sprintf("unable to connect to database: %v", err.Error()))
		}
		return dsClient

	}
	log.Debugf("development")
	dsClient, err := datastore.NewClient(
		context.Background(),
		os.Getenv(DS_PROJECT_ID),
		option.WithEndpoint(os.Getenv(DS_HOST)),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithGRPCConnectionPool(50),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to connect to database: %v", err.Error()))
	}
	return dsClient
}

func getCookieURL() string {
	return os.Getenv(COOKIE_URL)
}

func getLoginHost() string {
	return os.Getenv(LOGIN_HOST)
}

func getSecrets() (secrets, error) {
	s := secrets{
		Key: secretsKey(),
	}

	resp, err := retryablehttp.Get(getCookieURL())

	if err != nil {
		return s, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return s, err
	}

	err = json.Unmarshal(body, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

// func getSecrets() (secrets, error) {
// 	s := secrets{
// 		Key: secretsKey(),
// 	}
//
// 	c := context.Background()
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		return s, err
// 	}
//
// 	err = dsClient.Get(c, s.Key, &s)
// 	if err == nil {
// 		return s, nil
// 	}
//
// 	if err != datastore.ErrNoSuchEntity {
// 		return s, err
// 	}
//
// 	log.Warningf("generated new secrets")
// 	s, err = genSecrets()
// 	if err != nil {
// 		return s, err
// 	}
//
// 	_, err = dsClient.Put(c, s.Key, &s)
// 	return s, err
// }

func secretsKey() *datastore.Key {
	return datastore.NameKey("Secrets", "root", nil)
}

func genSecrets() (secrets, error) {
	s := secrets{
		HashKey:  securecookie.GenerateRandomKey(hashKeyLength),
		BlockKey: securecookie.GenerateRandomKey(blockKeyLength),
		Key:      secretsKey(),
	}

	if s.HashKey == nil {
		return s, fmt.Errorf("generated hashKey was nil")
	}

	if s.BlockKey == nil {
		return s, fmt.Errorf("generated blockKey was nil")
	}

	return s, nil
}

func (s *secrets) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(s, ps)
}

func (s *secrets) Save() ([]datastore.Property, error) {
	s.UpdatedAt = time.Now()
	return datastore.SaveStruct(s)
}

func (s *secrets) LoadKey(k *datastore.Key) error {
	s.Key = k
	return nil
}

// staticHandler for local development since app.yaml is ignored
// static files are handled via app.yaml routes when deployed
func staticRoutes(r *gin.Engine) *gin.Engine {
	if sn.IsProduction() {
		return r
	}
	r.StaticFile("/favicon.ico", "public/favicon.ico")
	r.Static("/images", "public/images")
	r.Static("/javascripts", "public/javascripts")
	r.Static("/js", "public/js")
	r.Static("/stylesheets", "public/stylesheets")
	r.Static("/rules", "public/rules")
	return r
}

func login(c *gin.Context) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	referer := c.Request.Referer()
	log.Debugf("referer: %v", referer)
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

func createCookieStore(s secrets) cookie.Store {
	if !sn.IsProduction() {
		log.Debugf("hashKey: %s\nblockKey: %s",
			base64.StdEncoding.EncodeToString(s.HashKey),
			base64.StdEncoding.EncodeToString(s.BlockKey),
		)
	}
	store := cookie.NewStore(s.HashKey, s.BlockKey)
	opts := sessions.Options{
		Domain: "slothninja.com",
		Path:   "/",
	}
	if sn.IsProduction() {
		opts.Secure = true
	}
	store.Options(opts)
	return store
}
