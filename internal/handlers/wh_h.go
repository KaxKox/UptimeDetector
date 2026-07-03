package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gocheck/internal/database"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type LiveLog struct {
	SiteName   string `json:"site_name"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	TimeMs     int    `json:"response_time_ms"`
}

func WsHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Błąd WebSocket:", err)
		return
	}
	defer ws.Close()

	log.Println("Ktoś podłączył się do Dashboardu na żywo!")

	for {
		query := `
			SELECT s.name, s.url, c.status_code, c.response_time_ms
			FROM checks c
			JOIN sites s ON c.site_id = s.id
			ORDER BY c.checked_at DESC
			LIMIT 5
		`
		
		rows, err := database.DB.Query(query)
		if err == nil {
			var logs []LiveLog
			for rows.Next() {
				var l LiveLog
				rows.Scan(&l.SiteName, &l.URL, &l.StatusCode, &l.TimeMs)
				logs = append(logs, l)
			}
			rows.Close()

			err = ws.WriteJSON(logs)
			if err != nil {
				log.Println("Klient się rozłączył")
				break
			}
		}

		time.Sleep(5 * time.Second)
	}
}