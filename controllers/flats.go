package controllers

import (
	"flatman-api/models"
	"flatman-api/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GetFlats(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var flats []models.Flat

	if err := models.DB.Find(&flats, "user_id = ?", userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flats)
}

func GetFlat(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var flat models.Flat

	if err := models.DB.Where("user_id = ?", userId).Find(&flat, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flat)
}

type SaveFlatInput struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Floor   uint   `json:"floor" binding:"required"`
}

func SaveFlat(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input SaveFlatInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newFlat := models.Flat{}

	copier.Copy(&newFlat, &input)
	newFlat.UserID = userId

	if err := models.DB.Create(&newFlat).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newFlat)
}

type UpdateFlatInput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Floor   uint   `json:"floor"`
}

func UpdateFlat(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input UpdateFlatInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var flat models.Flat

	if err := models.DB.Where("user_id = ?", userId).Find(&flat, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Model(&flat).Updates(input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flat)
}

func DeleteFlat(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var flat models.Flat

	if err := models.DB.First(&flat, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if flat.UserID != userId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you can only delete your own flats"})
		return
	}

	if err := models.DB.Delete(&flat).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
