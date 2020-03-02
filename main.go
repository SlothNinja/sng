package main

import (

	// _ "net/http/pprof"

	// "bitbucket.org/SlothNinja/atf"
	// "bitbucket.org/SlothNinja/confucius"
	// "bitbucket.org/SlothNinja/got"
	// "bitbucket.org/SlothNinja/indonesia"
	// "bitbucket.org/SlothNinja/slothninja-games/sn/game"
	// "bitbucket.org/SlothNinja/slothninja-games/sn/rating"
	// "bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	// "bitbucket.org/SlothNinja/slothninja-games/sn/send"
	// "bitbucket.org/SlothNinja/slothninja-games/sn/type"
	// "bitbucket.org/SlothNinja/slothninja-games/sn/user"
	// "bitbucket.org/SlothNinja/slothninja-games/sn/user_controller"
	// "bitbucket.org/SlothNinja/tammany"

	"github.com/SlothNinja/game"
	"github.com/SlothNinja/got"
	"github.com/SlothNinja/restful"
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
	userPrefix     = "user"
	gamesPrefix    = "games"
	ratingPrefix   = "rating"
	mailPrefix     = "mail"
	rootPath       = "/"
	hashKeyLength  = 64
	blockKeyLength = 32
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
		sessions.Sessions("sngsession", store),
		// restful.CTXHandler(),
		restful.TemplateHandler(r),
		// user.GetGUserHandler,
		user.GetCUserHandler,
	)

	// Welcome Page (index.html) route
	welcome.AddRoutes(r)

	// // Mail route
	// send.AddRoutes(mailPrefix, r)

	// // Games Routes
	game.AddRoutes(gamesPrefix, r)

	// // User Routes
	user_controller.AddRoutes(userPrefix, r)

	// // Rating Routes
	// rating.AddRoutes(ratingPrefix, r)

	// // After The Flood
	// atf.Register(gtype.ATF, r)

	// // Guild of Thieves
	got.Register(gtype.GOT, r)

	// // Tammany Hall
	// tammany.Register(gtype.Tammany, r)

	// // Indonesia
	// indonesia.Register(gtype.Indonesia, r)

	// // Confucius
	// confucius.Register(gtype.Confucius, r)

	// http.Handle(rootPath, r)

	r.Run()
}
