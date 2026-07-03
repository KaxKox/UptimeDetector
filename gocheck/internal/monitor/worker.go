package monitor

import (
	"log"
	"net/http"
	"time"

	"gocheck/internal/database"
)

func Start() {
	for {
		runChecks()

		time.Sleep(15 * time.Second)
	}
}

func runChecks() {
	rows, err := database.DB.Query("SELECT id, url FROM sites")
	if err != nil {
		log.Println("Błąd pobierania stron do sprawdzenia:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		rows.Scan(&id, &url)

		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start).Milliseconds()

		statusCode := 0
		if err != nil {
			log.Printf("Strona %s nie odpowiada!\n", url)
		} else {
			statusCode = resp.StatusCode
			log.Printf("Strona %s odpowiedziała kodem %d w %d ms\n", url, statusCode, duration)
			resp.Body.Close()
		}

		_, dbErr := database.DB.Exec(
			"INSERT INTO checks (site_id, status_code, response_time_ms) VALUES ($1, $2, $3)",
			id, statusCode, duration,
		)
		if dbErr != nil {
			log.Println("Błąd zapisu logu do bazy:", dbErr)
		}
	}
}