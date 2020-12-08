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
	NODE_ENV         = "NODE_ENV"
	UserProjectID    = "UserProjectID"
	production       = "production"
	userPrefix       = "user"
	gamesPrefix      = "games"
	ratingPrefix     = "rating"
	mailPrefix       = "mail"
	rootPath         = "/"
	hashKeyLength    = 64
	blockKeyLength   = 32
	sessionName      = "sng-oauth"
	EmulatorUserHost = "EmulatorUserHost"
	msgEnter         = "Entering"
	msgExit          = "Exiting"
)

func main() {
	setGinMode()
	userClient := getUserClient()

	db, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		panic(fmt.Sprintf("unable to connect to database: %v", err.Error()))
	}

	mcache := cache.New(30*time.Minute, 10*time.Minute)

	s, err := getSecrets()
	if err != nil {
		panic(err.Error())
	}

	store := cookie.NewStore(s.HashKey, s.BlockKey)
	// store := sessions.NewCookieStore([]byte("secret123"))

	r := gin.Default()
	renderer := restful.ParseTemplates("templates/", ".tmpl")
	r.HTMLRender = renderer

	r.Use(
		sessions.Sessions(sessionName, store),
		restful.AddTemplates(renderer.Templates),
		// user.GetCUserHandler(db),
	)

	// Welcome Page (index.html) route
	welcome.AddRoutes(r)

	// Games Routes
	r = game.NewClient(db).AddRoutes(gamesPrefix, r)

	// User Routes
	// r = ucon.NewClient(db).AddRoutes(userPrefix, r)

	// Rating Routes
	r = rating.NewClient(userClient, db).AddRoutes(ratingPrefix, r)

	// After The Flood
	r = atf.NewClient(db, userClient, mcache).Register(gtype.ATF, r)

	// Guild of Thieves
	r = got.NewClient(db, userClient, mcache).Register(gtype.GOT, r)

	// Tammany Hall
	r = tammany.NewClient(db, userClient, mcache).Register(gtype.Tammany, r)

	// Indonesia
	r = indonesia.NewClient(db, userClient, mcache).Register(gtype.Indonesia, r)

	// Confucius
	r = confucius.NewClient(db, userClient, mcache).Register(gtype.Confucius, r)

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

func getUserClient() *datastore.Client {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	if sn.IsProduction() {
		log.Debugf("production")
		userClient, err := datastore.NewClient(
			context.Background(),
			os.Getenv(UserProjectID),
		)
		if err != nil {
			panic(fmt.Sprintf("unable to connect to user database: %v", err.Error()))
		}
		return userClient

	}
	log.Debugf("development")
	userClient, err := datastore.NewClient(
		context.Background(),
		os.Getenv(UserProjectID),
		option.WithEndpoint(os.Getenv(EmulatorUserHost)),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithGRPCConnectionPool(50),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to connect to user database: %v", err.Error()))
	}
	return userClient
}

func getSecrets() (*secrets, error) {
	resp, err := retryablehttp.Get("http://luser.slothninja.com:8087/cookie")

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	s := new(secrets)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
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

	c.Redirect(http.StatusSeeOther, "http://luser.slothninja.com:8087/login?redirect="+encodedReferer)
}

func logout(c *gin.Context) {
	log.Debugf(msgEnter)
	defer log.Debugf(msgExit)

	referer := c.Request.Referer()
	encodedReferer := base64.StdEncoding.EncodeToString([]byte(referer))

	c.Redirect(http.StatusSeeOther, "http://luser.slothninja.com:8087/logout?redirect="+encodedReferer)
}
