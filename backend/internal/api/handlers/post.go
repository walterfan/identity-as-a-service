package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/walterfan/idaas/backend/internal/config"
	"github.com/walterfan/idaas/backend/internal/models"
)

// GetPosts returns all blog posts
func GetPosts(c *gin.Context) {
	var posts []models.Post
	
	query := config.DB.Order("created_at desc")
	
	// For non-admin users, only show published posts
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		query = query.Where("published = ?", true)
	}
	
	// Filter by tag if provided
	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	var total int64
	query.Model(&models.Post{}).Count(&total)
	
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetPost returns a specific blog post by slug
func GetPost(c *gin.Context) {
	slug := c.Param("slug")
	
	var post models.Post
	if err := config.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// For non-admin users, only show published posts
	role, exists := c.Get("role")
	if (!exists || role != "admin") && !post.Published {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"post": post})
}

// CreatePost creates a new blog post
func CreatePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var input struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Excerpt   string `json:"excerpt"`
		ImageURL  string `json:"image_url"`
		Published bool   `json:"published"`
		Tags      string `json:"tags"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Generate slug from title
	postSlug := slug.Make(input.Title)
	
	// Check if slug already exists
	var count int64
	config.DB.Model(&models.Post{}).Where("slug = ?", postSlug).Count(&count)
	if count > 0 {
		postSlug = postSlug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}
	
	// Create excerpt if not provided
	excerpt := input.Excerpt
	if excerpt == "" && len(input.Content) > 150 {
		excerpt = input.Content[:150] + "..."
	} else if excerpt == "" {
		excerpt = input.Content
	}
	
	var publishedAt *time.Time
	if input.Published {
		now := time.Now()
		publishedAt = &now
	}
	
	post := models.Post{
		Title:       input.Title,
		Slug:        postSlug,
		Content:     input.Content,
		Excerpt:     excerpt,
		ImageURL:    input.ImageURL,
		Published:   input.Published,
		PublishedAt: publishedAt,
		Tags:        input.Tags,
		UserID:      userID.(uint),
	}
	
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"post": post})
}

// UpdatePost updates an existing blog post
// UpdatePost updates an existing blog post
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Check if the user is the owner or an admin
	role, _ := c.Get("role")
	if post.UserID != userID.(uint) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this post"})
		return
	}
	
	var input struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		Excerpt   string `json:"excerpt"`
		ImageURL  string `json:"image_url"`
		Published *bool  `json:"published"`
		Tags      string `json:"tags"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update fields if provided
	if input.Title != "" && input.Title != post.Title {
		post.Title = input.Title
		// Update slug if title changes
		post.Slug = slug.Make(input.Title)
		
		// Check if slug already exists
		var count int64
		config.DB.Model(&models.Post{}).Where("slug = ? AND id != ?", post.Slug, post.ID).Count(&count)
		if count > 0 {
			post.Slug = post.Slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
		}
	}
	
	if input.Content != "" {
		post.Content = input.Content
	}
	
	if input.Excerpt != "" {
		post.Excerpt = input.Excerpt
	} else if input.Content != "" && input.Excerpt == "" {
		// Update excerpt if content changed but excerpt not provided
		if len(input.Content) > 150 {
			post.Excerpt = input.Content[:150] + "..."
		} else {
			post.Excerpt = input.Content
		}
	}
	
	post.ImageURL = input.ImageURL
	
	if input.Published != nil && *input.Published != post.Published {
		post.Published = *input.Published
		if *input.Published && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}
	
	if input.Tags != "" {
		post.Tags = input.Tags
	}
	
	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"post": post})
}

// DeletePost deletes a blog post
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Check if the user is the owner or an admin
	role, _ := c.Get("role")
	if post.UserID != userID.(uint) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this post"})
		return
	}
	
	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}