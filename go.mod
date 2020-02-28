module github.com/SlothNinja/slothninja-games

require (
	github.com/SlothNinja/slothninja-games/sn/restful v0.0.0-20200227212226-2cd0a8e24f94
	github.com/SlothNinja/slothninja-games/welcome v0.0.0-20200227212226-2cd0a8e24f94
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.5.0
	github.com/gorilla/securecookie v1.1.1
	google.golang.org/appengine v1.6.5
)

replace github.com/SlothNinja/slothninja-games/sn/restful => ./restful
