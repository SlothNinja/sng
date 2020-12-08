package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/atf"
	"github.com/SlothNinja/confucius"
	"github.com/SlothNinja/game"
	"github.com/SlothNinja/got"
	"github.com/SlothNinja/indonesia"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/rating/v2"
	"github.com/SlothNinja/restful"
	"github.com/SlothNinja/sn"
	"github.com/SlothNinja/tammany"
	gtype "github.com/SlothNinja/type"
	ucon "github.com/SlothNinja/user-controller/v2"
	"github.com/SlothNinja/welcome/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/patrickmn/go-cache"
)

const (
	NODE_ENV       = "NODE_ENV"
	production     = "production"
	userPrefix     = "user"
	gamesPrefix    = "games"
	ratingPrefix   = "rating"
	mailPrefix     = "mail"
	rootPath       = "/"
	hashKeyLength  = 64
	blockKeyLength = 32
	sessionName    = "sng-oauth"
)

func main() {
	if sn.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	c := context.Background()
	client := newClient(c)
	// userClient := user.NewClient(client.DS)

	mcache := cache.New(30*time.Minute, 10*time.Minute)

	s, err := getSecrets()
	if err != nil {
		panic(err.Error())
	}

	store := createCookieStore(s.HashKey, s.BlockKey)
	r := gin.Default()
	renderer := restful.ParseTemplates("templates/", ".tmpl")
	r.HTMLRender = renderer

	r.Use(
		sessions.Sessions(sessionName, store),
		restful.AddTemplates(renderer.Templates),
		// user2.GetCUserHandler(userClient),
		// user.GetCUserHandler(userClient.DS),
		cors.Default(),
	)

	// Welcome Page (index.html) route
	welcome.AddRoutes(r)

	// Games Routes
	r = game.NewClient(client.DS).AddRoutes(gamesPrefix, r)

	// User Routes
	r = ucon.NewClient(client.DS).AddRoutes(userPrefix, r)

	// Rating Routes
	r = rating.NewClient(client.DS).AddRoutes(ratingPrefix, r)

	// After The Flood
	r = atf.NewClient(client.DS, mcache).Register(gtype.ATF, r)

	// Guild of Thieves
	r = got.NewClient(client.DS, mcache).Register(gtype.GOT, r)

	// Tammany Hall
	r = tammany.NewClient(client.DS, mcache).Register(gtype.Tammany, r)

	// Indonesia
	r = indonesia.NewClient(client.DS, mcache).Register(gtype.Indonesia, r)

	// Confucius
	r = confucius.NewClient(client.DS, mcache).Register(gtype.Confucius, r)

	// warmup
	r.GET("_ah/warmup", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// login
	r.GET("login", login)

	r = staticRoutes(r)

	r.Run()
}

type Client struct {
	DS *datastore.Client
}

func newClient(c context.Context) Client {
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		panic(fmt.Sprintf("unable to connect to database: %v", err.Error()))
	}
	return Client{dsClient}
}

type secrets struct {
	HashKey   []byte
	BlockKey  []byte
	UpdatedAt time.Time
	Key       *datastore.Key `datastore:"__key__"`
}

func createCookieStore(hashKey, blockKey []byte) cookie.Store {
	if !sn.IsProduction() {
		log.Debugf("hashKey: %s\nblockKey: %s",
			base64.StdEncoding.EncodeToString(hashKey),
			base64.StdEncoding.EncodeToString(blockKey),
		)
	}
	store := cookie.NewStore(hashKey, blockKey)
	store.Options(sessions.Options{
		Domain: "slothninja.com",
		Path:   "/",
	})
	return store
}

func getSecrets() (*secrets, error) {
	resp, err := http.Get("http://luser.slothninja.com:8087/cookie")

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

// func (client Client) getSecrets(c context.Context) (secrets, error) {
// 	s := secrets{
// 		Key: secretsKey(),
// 	}
//
// 	err := client.DS.Get(c, s.Key, &s)
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
// 	_, err = client.DS.Put(c, s.Key, &s)
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
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	referer := c.Request.Referer()
	encodedReferer := base64.StdEncoding.EncodeToString([]byte(referer))

	c.Redirect(http.StatusSeeOther, "http://luser.slothninja.com:8087/login?redirect="+encodedReferer)
}
