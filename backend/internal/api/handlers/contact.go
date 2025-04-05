package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/idaas/backend/internal/config"
	"github.com/walterfan/idaas/backend/internal/models"
)

// SubmitContact handles contact form submissions
func SubmitContact(c *gin.Context) {
	var input struct {
		Name    string `json:"name" binding:"required"`
		Email   string `json:"email" binding:"required,email"`
		Subject string `json:"subject"`
		Message string `json:"message" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	contact := models.Contact{
		Name:    input.Name,
		Email:   input.Email,
		Subject: input.Subject,
		Message: input.Message,
		Read:    false,
	}
	
	if err := config.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit contact message"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Contact message submitted successfully"})
}

// GetContacts returns all contact messages (admin only)
func GetContacts(c *gin.Context) {
	var contacts []models.Contact
	
	query := config.DB.Order("created_at desc")
	
	// Filter by read status if provided
	if readStatus := c.Query("read"); readStatus != "" {
		if readStatus == "true" {
			query = query.Where("read = ?", true)
		} else if readStatus == "false" {
			query = query.Where("read = ?", false)
		}
	}
	
	if err := query.Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact messages"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

// GetContact returns a specific contact message by ID (admin only)
func GetContact(c *gin.Context) {
	id := c.Param("id")
	
	var contact models.Contact
	if err := config.DB.First(&contact, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact message not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"contact": contact})
}

// MarkContactAsRead marks a contact message as read (admin only)
func MarkContactAsRead(c *gin.Context) {
	id := c.Param("id")
	
	var contact models.Contact
	if err := config.DB.First(&contact, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact message not found"})
		return
	}
	
	contact.Read = true
	
	if err := config.DB.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact message"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Contact message marked as read"})
}

// DeleteContact deletes a contact message (admin only)
func DeleteContact(c *gin.Context) {
	id := c.Param("id")
	
	var contact models.Contact
	if err := config.DB.First(&contact, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact message not found"})
		return
	}
	
	if err := config.DB.Delete(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact message"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Contact message deleted successfully"})
}