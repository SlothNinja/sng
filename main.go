package main

import (
	"context"
	"fmt"
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
	user_controller "github.com/SlothNinja/user-controller"
	"github.com/SlothNinja/welcome"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
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

	db, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		panic(fmt.Sprintf("unable to connect to database: %v", err.Error()))
	}

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
		user.GetCUserHandler(db),
	)

	// Welcome Page (index.html) route
	welcome.AddRoutes(r)

	// Games Routes
	game.AddRoutes(gamesPrefix, r)

	// User Routes
	user_controller.NewClient(db).AddRoutes(userPrefix, r)

	// Rating Routes
	r = rating.NewClient(db).AddRoutes(ratingPrefix, r)

	// After The Flood
	r = atf.NewClient(db).Register(gtype.ATF, r)

	// Guild of Thieves
	r = got.NewClient(db).Register(gtype.GOT, r)

	// Tammany Hall
	r = tammany.NewClient(db).Register(gtype.Tammany, r)

	// Indonesia
	r = indonesia.NewClient(db).Register(gtype.Indonesia, r)

	// Confucius
	r = confucius.NewClient(db).Register(gtype.Confucius, r)

	r.Run()
}

type secrets struct {
	HashKey   []byte
	BlockKey  []byte
	UpdatedAt time.Time
	Key       *datastore.Key `datastore:"__key__"`
}

func getSecrets() (secrets, error) {
	s := secrets{
		Key: secretsKey(),
	}

	c := context.Background()
	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return s, err
	}

	err = dsClient.Get(c, s.Key, &s)
	if err == nil {
		return s, nil
	}

	if err != datastore.ErrNoSuchEntity {
		return s, err
	}

	log.Warningf("generated new secrets")
	s, err = genSecrets()
	if err != nil {
		return s, err
	}

	_, err = dsClient.Put(c, s.Key, &s)
	return s, err
}

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
