package main

import (
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/game"
	"github.com/SlothNinja/sn"
	gtype "github.com/SlothNinja/type"
	"github.com/gin-gonic/gin"
)

func (cl *Client) myGamesHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	obj := struct {
		Options struct {
			ItemsPerPage int `json:"itemsPerPage"`
		} `json:"options"`
		Status  string `json:"status"`
		Type    string `json:"type"`
		Forward string `json:"forward"`
		UserID  int64  `json:"userId"`
	}{}

	err := c.ShouldBind(&obj)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cu, err := cl.User.Current(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	t := gtype.ToType[obj.Type]

	items := obj.Options.ItemsPerPage
	if t == gtype.All {
		items = -1
	}

	var (
		gotHeaders []*game.IndexEntry
		gotCnt     int
	)

	if t == gtype.All || t == gtype.GOT {
		gotHeaders, gotCnt, _, err = cl.Got.GamesIndex(c, game.Options{
			ItemsPerPage: -1,
			Kind:         "Header",
			Status:       game.ToStatus[obj.Status],
			Type:         gtype.ToType[obj.Type],
			UserID:       obj.UserID,
		})
		if err != nil {
			sn.JErr(c, err)
			return
		}
	}

	forward, err := datastore.DecodeCursor(obj.Forward)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	gHeaders, gCnt, forward, err := cl.Game.GamesIndex(c, game.Options{
		ItemsPerPage: items,
		Kind:         "Game",
		Status:       game.ToStatus[obj.Status],
		Type:         gtype.ToType[obj.Type],
		Forward:      forward,
		UserID:       obj.UserID,
	})
	if err != nil {
		sn.JErr(c, err)
		return
	}

	if t == gtype.All || t == gtype.GOT {
		forward = datastore.Cursor{}
	}

	c.JSON(http.StatusOK, gin.H{
		"gheaders":   append(gHeaders, gotHeaders...),
		"totalItems": gotCnt + gCnt,
		"forward":    forward.String(),
		"cu":         cu,
	})
}
