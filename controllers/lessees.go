package controllers

import (
	"flatman-api/models"
	"flatman-api/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GetLessees(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lessees []models.Lessee
	var flatIDs []int64

	models.DB.Model(&models.Flat{}).Where("user_id = ?", userId).Pluck("id", &flatIDs)

	if err := models.DB.Find(&lessees, "flat_id IN (?)", flatIDs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lessees)
}

func GetLessee(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var lessee models.Lessee

	if err := models.DB.Find(&lessee, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of the lessee
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, lessee.FlatID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only view lessees of your own flats."})
		return
	}

	c.JSON(http.StatusOK, lessee)
}

type SaveLesseeInput struct {
	FlatID     uint    `json:"flat_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	PostalCode string  `json:"postal_code" binding:"required"`
	From       string  `json:"from" binding:"required"`
	Until      string  `json:"until" binding:"required"`
	Tin        string  `json:"tin" binding:"required"`
	Rent       float64 `json:"rent" binding:"required"`
}

func SaveLessee(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input SaveLesseeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of flat
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, input.FlatID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add lessees to your own flats."})
		return
	}

	newLessee := models.Lessee{}
	copier.Copy(&newLessee, &input)

	if err := models.DB.Create(&newLessee).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newLessee)
}

type UpdateLesseeInput struct {
	FlatID     uint    `json:"flat_id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	PostalCode string  `json:"postal_code"`
	From       string  `json:"from"`
	Until      string  `json:"until"`
	Tin        string  `json:"tin"`
	Rent       float64 `json:"rent"`
}

func UpdateLessee(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input UpdateLesseeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var lessee models.Lessee

	if err := models.DB.Find(&lessee, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of the lessee
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, lessee.FlatID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update lessees of your own flats."})
		return
	}

	// Check if currently authenticated user is the owner of flat if it tries to update the flat_id
	if input.FlatID != 0 {
		if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, input.FlatID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only assign lessees to your own flats."})
			return
		}
	}

	if err := models.DB.Model(&lessee).Updates(input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lessee)
}

func DeleteLessee(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var lessee models.Lessee

	if err := models.DB.Find(&lessee, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of the lessee
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, lessee.FlatID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete lessees from your own flats."})
		return
	}

	if err := models.DB.Delete(&lessee).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
