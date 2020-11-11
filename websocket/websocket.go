package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andynl/youtube-stats/youtube"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Upgrade(w http.ResponseWriter, r *http.Request) (ws *websocket.Conn, err error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return ws, err
	}

	return ws, nil
}

func Writer(conn *websocket.Conn) {
	for {
		ticker := time.NewTimer(5 * time.Second)

		for t := range ticker.C {
			fmt.Println("Updating Stats: %+v\n", t)

			items, err := youtube.GetSubsribers()
			if err != nil {
				fmt.Println(err)
			}

			jsonString, err := json.Marshal(items)
			if err != nil {
				fmt.Println(err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
