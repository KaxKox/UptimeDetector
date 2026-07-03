package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gocheck/internal/models"
	"gocheck/internal/database"
)

func GetSites(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, url, interval FROM sites")
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
	var newSite models.Site
	if err := c.BindJSON(&newSite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zły format JSON"})
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO sites (name, url, interval) VALUES ($1, $2, $3) RETURNING id",
		newSite.Name, newSite.URL, newSite.Interval,
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