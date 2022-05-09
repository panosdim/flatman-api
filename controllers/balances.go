package controllers

import (
	"flatman-api/models"
	"flatman-api/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GetBalances(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var balances []models.Balance
	var flatIDs []int64

	models.DB.Model(&models.Flat{}).Where("user_id = ?", userId).Pluck("id", &flatIDs)

	if err := models.DB.Find(&balances, "flat_id IN (?)", flatIDs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balances)
}

func GetBalance(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var balance models.Balance

	if err := models.DB.Find(&balance, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of the balance
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, balance.FlatID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only view balance of your own flats."})
		return
	}

	c.JSON(http.StatusOK, balance)
}

type SaveBalanceInput struct {
	FlatID  uint    `json:"flat_id" binding:"required"`
	Date    string  `json:"date" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
	Comment string  `json:"comment" binding:"required"`
}

func SaveBalance(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input SaveBalanceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of balance
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, input.FlatID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add balance to your own flats."})
		return
	}

	newBalance := models.Balance{}
	copier.Copy(&newBalance, &input)

	if err := models.DB.Create(&newBalance).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBalance)
}

type UpdateBalanceInput struct {
	FlatID  uint    `json:"flat_id"`
	Date    string  `json:"date"`
	Amount  float64 `json:"amount"`
	Comment string  `json:"comment"`
}

func UpdateBalance(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input UpdateBalanceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var balance models.Balance

	if err := models.DB.Find(&balance, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of the balance
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, balance.FlatID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update balances of your own flats."})
		return
	}

	// Check if currently authenticated user is the owner of flat if it tries to update the flat_id
	if input.FlatID != 0 {
		if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, input.FlatID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only assign balances to your own flats."})
			return
		}
	}

	if err := models.DB.Model(&balance).Updates(input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

func DeleteBalance(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	var balance models.Balance

	if err := models.DB.Find(&balance, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if currently authenticated user is the owner of the balance
	if err := models.DB.Where("user_id = ?", userId).Find(&models.Flat{}, balance.FlatID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete balances from your own flats."})
		return
	}

	if err := models.DB.Delete(&balance).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
