package main

import (
	"net/http"

	"github.com/SlothNinja/atf"
	"github.com/SlothNinja/confucius"
	"github.com/SlothNinja/game"
	"github.com/SlothNinja/got"
	"github.com/SlothNinja/indonesia"
	"github.com/SlothNinja/rating"
	"github.com/SlothNinja/restful"
	"github.com/SlothNinja/send"
	"github.com/SlothNinja/tammany"
	gtype "github.com/SlothNinja/type"
	"github.com/SlothNinja/user"
	user_controller "github.com/SlothNinja/user-controller"
	"github.com/SlothNinja/welcome"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"google.golang.org/appengine"
)

const (
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
	// if appengine.IsDevAppServer() {
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	hashKey := securecookie.GenerateRandomKey(hashKeyLength)
	if hashKey == nil {
		panic("generated hashKey was nil")
	}

	blockKey := securecookie.GenerateRandomKey(blockKeyLength)
	if blockKey == nil {
		panic("generated blockKey was nil")
	}

	store := cookie.NewStore(hashKey, blockKey)
	// store := sessions.NewCookieStore([]byte("secret123"))

	r := gin.Default()
	r.Use(
		sessions.Sessions(sessionName, store),
		// restful.CTXHandler(),
		restful.TemplateHandler(r),
		// user.GetGUserHandler,
		user.GetCUserHandler,
	)

	// Welcome Page (index.html) route
	welcome.AddRoutes(r)

	// Mail route
	send.AddRoutes(mailPrefix, r)

	// Games Routes
	game.AddRoutes(gamesPrefix, r)

	// User Routes
	user_controller.AddRoutes(userPrefix, r)

	// Rating Routes
	rating.AddRoutes(ratingPrefix, r)

	// After The Flood
	atf.Register(gtype.ATF, r)

	// Guild of Thieves
	got.Register(gtype.GOT, r)

	// Tammany Hall
	tammany.Register(gtype.Tammany, r)

	// Indonesia
	indonesia.Register(gtype.Indonesia, r)

	// Confucius
	confucius.Register(gtype.Confucius, r)

	http.Handle(rootPath, r)
	appengine.Main()
}
