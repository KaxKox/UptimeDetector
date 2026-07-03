package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gocheck/internal/database"
)

func ExportCSV(c *gin.Context) {
	userID, _ := c.Get("userID")

	query := `
		SELECT s.name, s.url, c.status_code, c.response_time_ms, c.checked_at
		FROM checks c
		JOIN sites s ON c.site_id = s.id
		WHERE s.user_id = $1
		ORDER BY c.checked_at DESC
	`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania danych"})
		return
	}
	defer rows.Close()

	c.Header("Content-Disposition", `attachment; filename="raport_stron.csv"`)
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{"Nazwa", "URL", "Kod Odpowiedzi", "Czas (ms)", "Data Sprawdzenia"})

	for rows.Next() {
		var name, url, checkedAt string
		var statusCode, timeMs int
		rows.Scan(&name, &url, &statusCode, &timeMs, &checkedAt)

		row := []string{
			name,
			url,
			strconv.Itoa(statusCode),
			strconv.Itoa(timeMs),
			checkedAt,
		}
		writer.Write(row)
	}
}

func ImportCSV(c *gin.Context) {
	userID, _ := c.Get("userID")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nie znaleziono pliku w żądaniu"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd otwarcia pliku"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zły format pliku CSV"})
		return
	}

	addedCount := 0

	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		
		name := record[0]
		url := record[1]
		interval, _ := strconv.Atoi(record[2])

		_, err := database.DB.Exec(
			"INSERT INTO sites (user_id, name, url, interval) VALUES ($1, $2, $3, $4)",
			userID, name, url, interval,
		)
		if err == nil {
			addedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Zaimportowano pomyślnie", "added": addedCount})
}