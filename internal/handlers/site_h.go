package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gocheck/internal/models"
)

var sites = []models.Site{
	{ID: "1", Name: "Google", URL: "https://google.com", Interval: 60},
	{ID: "2", Name: "Facebook", URL: "https://facebook.com", Interval: 120},
}

func GetSites(c *gin.Context) {
	c.JSON(http.StatusOK, sites)
}

func CreateSite(c *gin.Context) {
	var newSite models.Site

	if err := c.BindJSON(&newSite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zły format danych JSON"})
		return
	}

	sites = append(sites, newSite)
	c.JSON(http.StatusCreated, newSite)
}

func UpdateSite(c *gin.Context) {
	id := c.Param("id")
	var updatedData models.Site

	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zły format danych JSON"})
		return
	}

	for i, site := range sites {
		if site.ID == id {
			sites[i].Name = updatedData.Name
			sites[i].URL = updatedData.URL
			sites[i].Interval = updatedData.Interval

			c.JSON(http.StatusOK, sites[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Nie znaleziono strony o podanym ID"})
}

func DeleteSite(c *gin.Context) {
	id := c.Param("id")

	for i, site := range sites {
		if site.ID == id {
			sites = append(sites[:i], sites[i+1:]...)
			
			c.JSON(http.StatusOK, gin.H{"message": "Strona usunięta pomyślnie"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Nie znaleziono strony o podanym ID"})
}