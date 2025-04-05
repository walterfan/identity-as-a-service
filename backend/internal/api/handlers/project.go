package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/idaas/backend/internal/config"
	"github.com/walterfan/idaas/backend/internal/models"
)

// GetProjects returns all projects
func GetProjects(c *gin.Context) {
	var projects []models.Project
	
	query := config.DB.Order("created_at desc")
	
	// Filter by featured if requested
	if featured := c.Query("featured"); featured == "true" {
		query = query.Where("featured = ?", true)
	}
	
	// Filter by tag if provided
	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	
	if err := query.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// GetProject returns a specific project by ID
func GetProject(c *gin.Context) {
	id := c.Param("id")
	
	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// CreateProject creates a new project
func CreateProject(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		ProjectURL  string `json:"project_url"`
		GithubURL   string `json:"github_url"`
		Featured    bool   `json:"featured"`
		Tags        string `json:"tags"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	project := models.Project{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		ProjectURL:  input.ProjectURL,
		GithubURL:   input.GithubURL,
		Featured:    input.Featured,
		Tags:        input.Tags,
		UserID:      userID.(uint),
	}
	
	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// UpdateProject updates an existing project
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	
	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	
	// Check if the user is the owner or an admin
	role, _ := c.Get("role")
	if project.UserID != userID.(uint) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this project"})
		return
	}
	
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		ProjectURL  string `json:"project_url"`
		GithubURL   string `json:"github_url"`
		Featured    *bool  `json:"featured"`
		Tags        string `json:"tags"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update fields if provided
	if input.Title != "" {
		project.Title = input.Title
	}
	if input.Description != "" {
		project.Description = input.Description
	}
	project.ImageURL = input.ImageURL
	project.ProjectURL = input.ProjectURL
	project.GithubURL = input.GithubURL
	if input.Featured != nil {
		project.Featured = *input.Featured
	}
	if input.Tags != "" {
		project.Tags = input.Tags
	}
	
	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// DeleteProject deletes a project
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	
	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	
	// Check if the user is the owner or an admin
	role, _ := c.Get("role")
	if project.UserID != userID.(uint) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this project"})
		return
	}
	
	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}