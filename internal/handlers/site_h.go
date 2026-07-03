package handlers

import (
	"gocheck/internal/database"
	"gocheck/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSites(c *gin.Context) {
	userID, _ := c.Get("userID")

	rows, err := database.DB.Query("SELECT id, name, url, interval FROM sites WHERE user_id=$1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania z bazy"})
		return
	}
	defer rows.Close()

	var sites []models.Site
	for rows.Next() {
		var s models.Site
		rows.Scan(&s.ID, &s.Name, &s.URL, &s.Interval)
		sites = append(sites, s)
	}

	if sites == nil {
		sites = []models.Site{}
	}

	c.JSON(http.StatusOK, sites)
}

func CreateSite(c *gin.Context) {
	userID, _ := c.Get("userID")

	var newSite models.Site
	if err := c.BindJSON(&newSite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zły format JSON"})
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO sites (user_id, name, url, interval) VALUES ($1, $2, $3, $4) RETURNING id",
		userID, newSite.Name, newSite.URL, newSite.Interval,
	).Scan(&newSite.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd zapisu do bazy"})
		return
	}

	c.JSON(http.StatusCreated, newSite)
}

func DeleteSite(c *gin.Context) {
	id := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM sites WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd przy usuwaniu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usunięto stronę"})
}

func UpdateSite(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var updatedData models.Site
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zły format JSON"})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE sites SET name=$1, url=$2, interval=$3 WHERE id=$4 AND user_id=$5",
		updatedData.Name, updatedData.URL, updatedData.Interval, id, userID,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas aktualizacji w bazie"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nie znaleziono strony o tym ID lub nie masz do niej uprawnień"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strona zaktualizowana pomyślnie"})
}