# Building a Personal Website with Go/Gin and Vue.js

I'll guide you through building a personal website with the technologies you've specified. Let's break this down into manageable steps.

## 1. Project Requirements and Structure

Let's start by setting up the project structure:

```bash
mkdir -p /Users/walter.fan/workspace/walter/idaas
cd /Users/walter.fan/workspace/walter/idaas
```

Here's the recommended project structure:

```
idaas/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes.go
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   └── service/
│   ├── pkg/
│   │   ├── auth/
│   │   ├── logger/
│   │   └── validator/
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── composables/
│   │   ├── layouts/
│   │   ├── router/
│   │   ├── store/
│   │   ├── views/
│   │   ├── App.vue
│   │   └── main.js
│   ├── .env
│   ├── package.json
│   └── tailwind.config.js
└── README.md
```

## 2. Backend Development

### Initialize Go Module

```bash
mkdir -p /Users/walter.fan/workspace/walter/idaas/backend
cd /Users/walter.fan/workspace/walter/idaas/backend
go mod init github.com/walterfan/idaas/backend
```

### Install Dependencies

```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/gin-contrib/cors
go get -u github.com/go-playground/validator/v10
```

### Create Database Models

Let's create the models for our application:

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/models/user.go
package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"-"`
	Name     string `gorm:"size:255" json:"name"`
	Bio      string `gorm:"type:text" json:"bio"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Role     string `gorm:"size:20;default:'user'" json:"role"` // admin or user
}

// BeforeSave hashes the password before saving to the database
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword compares the provided password with the stored hash
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
```

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/models/project.go
package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	ImageURL    string `gorm:"size:255" json:"image_url"`
	ProjectURL  string `gorm:"size:255" json:"project_url"`
	GithubURL   string `gorm:"size:255" json:"github_url"`
	Featured    bool   `gorm:"default:false" json:"featured"`
	Tags        string `gorm:"size:255" json:"tags"` // Comma-separated tags
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"-"`
}
```

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/models/post.go
package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	gorm.Model
	Title     string    `gorm:"size:255;not null" json:"title"`
	Slug      string    `gorm:"size:255;not null;unique" json:"slug"`
	Content   string    `gorm:"type:text" json:"content"`
	Excerpt   string    `gorm:"size:500" json:"excerpt"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	Published bool      `gorm:"default:false" json:"published"`
	PublishedAt *time.Time `json:"published_at"`
	Tags      string    `gorm:"size:255" json:"tags"` // Comma-separated tags
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}
```

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/models/contact.go
package models

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	Name    string `gorm:"size:255;not null" json:"name"`
	Email   string `gorm:"size:255;not null" json:"email"`
	Subject string `gorm:"size:255" json:"subject"`
	Message string `gorm:"type:text;not null" json:"message"`
	Read    bool   `gorm:"default:false" json:"read"`
}
```

### Database Configuration

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/config/database.go
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/walterfan/idaas/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Create data directory if it doesn't exist
	dataDir := "./data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	dbPath := filepath.Join(dataDir, "personal_website.db")
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = database.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Post{},
		&models.Contact{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = database
	log.Println("Database connected and migrated successfully")
}
```

### JWT Authentication

```go:/Users/walter.fan/workspace/walter/idaas/backend/pkg/auth/jwt.go
package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(getEnv("JWT_SECRET", "your-secret-key"))

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint, role string) (string, error) {
	// Set expiration time
	expirationHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Create claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "personal-website",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Helper function to get environment variables with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
```

### Middleware

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/idaas/backend/pkg/auth"
)

// AuthMiddleware verifies the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID and role in the context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminMiddleware checks if the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
```

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/middleware/cors.go
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// CorsMiddleware configures CORS for the application
func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
```

### Handlers

Let's create some basic handlers for our API:

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/handlers/auth.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/idaas/backend/internal/config"
	"github.com/walterfan/idaas/backend/internal/models"
	"github.com/walterfan/idaas/backend/pkg/auth"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// Login handles user authentication and returns a JWT token
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
			"role":     user.Role,
		},
	})
}

// Register creates a new user account
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Create new user
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Name:     input.Name,
		Role:     "user", // Default role
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
		},
	})
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
			"bio":      user.Bio,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	})
}

// UpdateProfile updates the current user's profile
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if input.Name != "" {
		user.Name = input.Name
	}
	user.Bio = input.Bio
	user.Avatar = input.Avatar

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
			"bio":      user.Bio,
			"avatar":   user.Avatar,
		},
	})
}
```

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/handlers/project.go
package handlers

import (
	"net/http"
	"strconv"

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
```

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/handlers/post.go
package handlers

import (
	"net/http"
	"strconv"
	"strings"
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
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	
	var post models.Post
	if err

----
Let me continue with the backend development by completing the post handler and adding the contact handler, routes, and main server file.

### Completing the Post Handler

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/handlers/post.go
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
```

### Contact Handler

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/handlers/contact.go
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
```

### API Routes

```go:/Users/walter.fan/workspace/walter/idaas/backend/internal/api/routes.go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/walterfan/idaas/backend/internal/api/handlers"
	"github.com/walterfan/idaas/backend/internal/api/middleware"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Apply global middleware
	router.Use(middleware.CorsMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Public routes
	public := router.Group("/api")
	{
		// Auth routes
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/register", handlers.Register)
		
		// Project routes
		public.GET("/projects", handlers.GetProjects)
		public.GET("/projects/:id", handlers.GetProject)
		
		// Blog post routes
		public.GET("/posts", handlers.GetPosts)
		public.GET("/posts/:slug", handlers.GetPost)
		
		// Contact route
		public.POST("/contact", handlers.SubmitContact)
	}
	
	// Protected routes (require authentication)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User profile routes
		protected.GET("/profile", handlers.GetProfile)
		protected.PUT("/profile", handlers.UpdateProfile)
		
		// Project management routes
		protected.POST("/projects", handlers.CreateProject)
		protected.PUT("/projects/:id", handlers.UpdateProject)
		protected.DELETE("/projects/:id", handlers.DeleteProject)
		
		// Blog post management routes
		protected.POST("/posts", handlers.CreatePost)
		protected.PUT("/posts/:id", handlers.UpdatePost)
		protected.DELETE("/posts/:id", handlers.DeletePost)
	}
	
	// Admin routes (require admin role)
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Contact management routes
		admin.GET("/contacts", handlers.GetContacts)
		admin.GET("/contacts/:id", handlers.GetContact)
		admin.PUT("/contacts/:id/read", handlers.MarkContactAsRead)
		admin.DELETE("/contacts/:id", handlers.DeleteContact)
	}
}
```

### Main Server File

```go:/Users/walter.fan/workspace/walter/idaas/backend/cmd/server/main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/walterfan/idaas/backend/internal/api"
	"github.com/walterfan/idaas/backend/internal/config"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	
	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)
	
	// Connect to database
	config.ConnectDatabase()
	
	// Create a new Gin router
	router := gin.Default()
	
	// Setup routes
	api.SetupRoutes(router)
	
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Start the server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

### Create a .env file for the backend

```plaintext:/Users/walter.fan/workspace/walter/idaas/backend/.env
# Server configuration
PORT=8080
GIN_MODE=debug

# JWT configuration
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION_HOURS=24
```

## 3. Frontend Development (Vue.js/Tailwind CSS)

Now let's set up the frontend with Vue.js 3 and Tailwind CSS.

### Initialize Vue.js Project

```bash
mkdir -p /Users/walter.fan/workspace/walter/idaas/frontend
cd /Users/walter.fan/workspace/walter/idaas/frontend
npm init vue@latest .
```

When prompted, select:
- Vue version: 3.x
- TypeScript: No
- JSX: No
- Vue Router: Yes
- Pinia: Yes
- Vitest: No
- End-to-End Testing: No
- ESLint: Yes
- Prettier: Yes

### Install Dependencies

```bash
cd /Users/walter.fan/workspace/walter/idaas/frontend
npm install
npm install axios tailwindcss postcss autoprefixer @headlessui/vue @heroicons/vue
npm install -D @tailwindcss/forms @tailwindcss/typography
```

### Configure Tailwind CSS

```bash
npx tailwindcss init -p
```

Update the Tailwind configuration:

```js:/Users/walter.fan/workspace/walter/idaas/frontend/tailwind.config.js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
  darkMode: 'class',
}
```

Add Tailwind directives to your CSS:

```css:/Users/walter.fan/workspace/walter/idaas/frontend/src/assets/main.css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer components {
  .btn {
    @apply px-4 py-2 rounded font-medium focus:outline-none focus:ring-2 focus:ring-offset-2;
  }
  
  .btn-primary {
    @apply btn bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
  }
  
  .btn-secondary {
    @apply btn bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-500;
  }
  
  .btn-danger {
    @apply btn bg-red-600 text-white hover:bg-red-700 focus:ring-red-500;
  }
  
  .form-input {
    @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500;
  }
  
  .form-label {
    @apply block text-sm font-medium text-gray-700 mb-1;
  }
}
```

### Create API Service

```js:/Users/walter.fan/workspace/walter/idaas/frontend/src/services/api.js
import axios from 'axios';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to include the auth token in requests
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add a response interceptor to handle common errors
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      // If unauthorized, clear token and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default {
  // Auth
  login(credentials) {
    return apiClient.post('/auth/login', credentials);
  },
  register(userData) {
    return apiClient.post('/auth/register', userData);
  },
  getProfile() {
    return apiClient.get('/profile');
  },
  updateProfile(profileData) {
    return apiClient.put('/profile', profileData);
  },
  
  // Projects
  getProjects(params = {}) {
    return apiClient.get('/projects', { params });
  },
  getProject(id) {
    return apiClient.get(`/projects/${id}`);
  },
  createProject(projectData) {
    return apiClient.post('/projects', projectData);
  },
  updateProject(id, projectData) {
    return apiClient.put(`/projects/${id}`, projectData);
  },
  deleteProject(id) {
    return apiClient.delete(`/projects/${id}`);
  },
  
  // Blog Posts
  getPosts(params = {}) {
    return apiClient.get('/posts', { params });
  },
  getPost(slug) {
    return apiClient.get(`/posts/${slug}`);
  },
  createPost(postData) {
    return apiClient.post('/posts', postData);
  },
  updatePost(id, postData) {
    return apiClient.put(`/posts/${id}`, postData);
  },
  deletePost(id) {
    return apiClient.delete(`/posts/${id}`);
  },
  
  // Contact
  submitContact(contactData) {
    return apiClient.post('/contact', contactData);
  },
  getContacts(params = {}) {
    return apiClient.get('/admin/contacts', { params });
  },
  getContact(id) {
    return apiClient.get(`/admin/contacts/${id}`);
  },
  markContactAsRead(id) {
    return apiClient.put(`/admin/contacts/${id}/read`);
  },
  deleteContact(id) {
    return apiClient.delete(`/admin/contacts/${id}`);
  },
};
```

### Create Store with Pinia

```js:/Users/walter.fan/workspace/walter/idaas/frontend/src/stores/auth.js
import { defineStore } from 'pinia';
import api from '@/services/api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
    token: localStorage.getItem('token') || null,
    loading: false,
    error: null,
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user && state.user.role === 'admin',
  },
  
  actions: {
    async login(credentials) {
      this.loading = true;
      this.error = null;
      try {
        const response = await api.login(credentials);
        this.setUserData(response.data);
        return response;
      } catch (error) {
        this.error = error.response?.data?.error || 'Login failed';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async register(userData) {
      this.loading = true;
      this.error = null;
      try {
        const response = await api.register(userData);
        return response;
      } catch (error) {
        this.error = error.response?.data?.error || 'Registration failed';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchProfile() {
      if (!this.token) return;
      
      this.loading = true;
      try {
        const response = await api.getProfile();
        this.user = response.data.user;
        localStorage.setItem('user', JSON.stringify(this.user));
        return response;
      } catch (error) {
        console.error('Failed to fetch profile:', error);
      } finally {
        this.loading = false;
      }
    },
    
    setUserData(data) {
      this.token = data.token;
      this.user = data.user;
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
    },
    
    logout() {
      this.user = null;
      this.token = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    },
  },
});
```

### Create Router Configuration

```js:/Users/walter.fan/workspace/walter/idaas/frontend/src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

// Public pages
import HomePage from '@/views/HomePage.vue';
import ProjectsPage from '@/views/ProjectsPage.vue';
import ProjectDetailPage from '@/views/ProjectDetailPage.vue';
import BlogPage from '@/views/BlogPage.vue';
import PostDetailPage from '@/views/PostDetailPage.vue';
import ContactPage from '@/views/ContactPage.vue';
import LoginPage from '@/views/LoginPage.vue';
import RegisterPage from '@/views/RegisterPage.vue';

// Admin pages
import AdminDashboard from '@/views/admin/DashboardPage.vue';
import AdminProjects from '@/views/admin/ProjectsPage.vue';
import AdminProjectEdit from '@/views/admin/ProjectEditPage.vue';
import AdminPosts from '@/views/admin/PostsPage.vue';
import AdminPostEdit from '@/views/admin/PostEditPage.vue';
import AdminContacts from '@/views/admin/ContactsPage.vue';
import AdminProfile from '@/views/admin/ProfilePage.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Public routes
    {
      path: '/',
      name: 'home',
      component: HomePage,
      meta: { title: 'Home' }
    },
    {
      path: '/projects',
      name: 'projects',
      component: ProjectsPage,
      meta: { title: 'Projects' }
    },
    {
      path: '/projects/:id',
      name: 'project-detail',
      component: ProjectDetailPage,
      meta: { title: 'Project Detail' }
    },
    {
      path: '/blog',
      name: 'blog',
      component: BlogPage,
      meta: { title: 'Blog' }
    },
    {
      path: '/blog/:slug',
      name: 'post-detail',
      component: PostDetailPage,
      meta: { title: 'Blog Post' }
    },
    {
      path: '/contact',
      name: 'contact',
      component: ContactPage,
      meta: { title: 'Contact' }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
      meta: { title: 'Login', guest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterPage,
      meta: { title: 'Register', guest: true }
    },
    
    // Admin routes
    {
      path: '/admin',
      name: 'admin',
      component: AdminDashboard,
      meta: { requiresAuth: true, admin: true, title: 'Admin Dashboard' }
    },
    {
      path: '/admin/projects',
      name: 'admin-projects',
      component: AdminProjects,
      meta: { requiresAuth: true, admin: true, title: 'Manage Projects' }
    },
    {
      path: '/admin/projects/new',
      name: 'admin-project-new',
      component: AdminProjectEdit,
      meta: { requiresAuth: true, admin: true, title: 'New Project' }
    },
    {
      path: '/admin/projects/:id',
      name: 'admin-project-edit',
      component: AdminProjectEdit,
      meta: { requiresAuth: true, admin: true, title: 'Edit Project' }
    },
    {
      path: '/admin/posts',
      name: 'admin-posts',
      component: AdminPosts,
      meta: { requiresAuth: true, admin: true, title: 'Manage Posts' }
    },
    {
      path: '/admin/posts/new',
      name: 'admin-post-new',
      component: AdminPostEdit,
      meta: { requiresAuth: true, admin: true, title: 'New Post' }
    },
    {
      path: '/admin/posts/:id',
      name: 'admin-post-edit',
      component: AdminPostEdit,
      meta: { requiresAuth: true, admin: true, title: 'Edit Post' }
    },
    {
      path: '/admin/contacts',
      name: 'admin-contacts',
      component: AdminContacts,
      meta: { requiresAuth: true, admin: true, title: 'Manage Contacts' }
    },
    {
      path: '/admin/profile',
      name: 'admin-profile',
      component: AdminProfile,
      meta: { requiresAuth: true, title: 'Profile Settings' }
    },
    
    // 404 route
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
});

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const requiresAdmin = to.matched.some(record => record.meta.admin);
  const isGuestOnly = to.matched.some(record => record.meta.guest);
  
  // Set page title
  document.title = to.meta.title ? `${to.meta.title} | Personal Website` : 'Personal Website';
  
  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else if (requiresAdmin && !authStore.isAdmin) {
    next('/');
  } else if (isGuestOnly && authStore.isAuthenticated) {
    next('/');
  } else {
    next();
  }
});

export default router;
```

### Create Main Layout Components

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/components/layout/Navbar.vue
<template>
  <nav class="bg-white dark:bg-gray-800 shadow">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex">
          <div class="flex-shrink-0 flex items-center">
            <router-link to="/" class="text-xl font-bold text-gray-800 dark:text-white">
              My Portfolio
            </router-link>
          </div>
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <router-link 
              to="/" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path === '/' ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Home
            </router-link>
            <router-link 
              to="/projects" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path.startsWith('/projects') ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Projects
            </router-link>
            <router-link 
              to="/blog" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path.startsWith('/blog') ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Blog
            </router-link>
            <router-link 
              to="/contact" 
              class="inline-flex items-center px-1 pt-1 border-b-2" 
              :class="[$route.path === '/contact' ? 'border-blue-500 text-gray-900 dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200']"
            >
              Contact
            </router-link>
          </div>
        </div>
        <div class="hidden sm:ml-6 sm:flex sm:items-center">
          <!-- Dark mode toggle -->
          <button 
            @click="toggleDarkMode" 
            class="p-1 rounded-full text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200 focus:outline-none"
          >
            <span v-if="isDarkMode" class="sr-only">Switch to light mode</span>
            <span v-else class="sr-only">Switch to dark mode</span>
            <svg v-if="isDarkMode" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646
----

Let me continue with the frontend development by completing the Navbar component and creating other essential components for the personal website.

### Completing the Navbar Component

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/components/layout/Navbar.vue
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>

          <!-- Profile dropdown -->
          <div v-if="authStore.isAuthenticated" class="ml-3 relative">
            <div>
              <button 
                @click="isProfileOpen = !isProfileOpen" 
                class="flex text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <span class="sr-only">Open user menu</span>
                <div v-if="authStore.user?.avatar" class="h-8 w-8 rounded-full overflow-hidden">
                  <img :src="authStore.user.avatar" alt="User avatar" class="h-full w-full object-cover" />
                </div>
                <div v-else class="h-8 w-8 rounded-full bg-blue-500 flex items-center justify-center text-white">
                  {{ authStore.user?.name?.charAt(0) || authStore.user?.username?.charAt(0) || 'U' }}
                </div>
              </button>
            </div>
            
            <div 
              v-if="isProfileOpen" 
              class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white dark:bg-gray-700 ring-1 ring-black ring-opacity-5 focus:outline-none z-10"
            >
              <router-link 
                v-if="authStore.isAdmin" 
                to="/admin" 
                class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600"
                @click="isProfileOpen = false"
              >
                Dashboard
              </router-link>
              <router-link 
                to="/admin/profile" 
                class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600"
                @click="isProfileOpen = false"
              >
                Profile
              </router-link>
              <button 
                @click="logout" 
                class="block w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600"
              >
                Sign out
              </button>
            </div>
          </div>
          
          <div v-else class="flex items-center space-x-4">
            <router-link 
              to="/login" 
              class="text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200"
            >
              Login
            </router-link>
            <router-link 
              to="/register" 
              class="btn-primary"
            >
              Register
            </router-link>
          </div>
        </div>
        
        <!-- Mobile menu button -->
        <div class="flex items-center sm:hidden">
          <button 
            @click="toggleDarkMode" 
            class="p-1 mr-2 rounded-full text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200 focus:outline-none"
          >
            <svg v-if="isDarkMode" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>
          <button 
            @click="isMobileMenuOpen = !isMobileMenuOpen" 
            class="inline-flex items-center justify-center p-2 rounded-md text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
          >
            <span class="sr-only">Open main menu</span>
            <svg 
              v-if="!isMobileMenuOpen" 
              class="block h-6 w-6" 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor" 
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg 
              v-else 
              class="block h-6 w-6" 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor" 
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Mobile menu -->
    <div v-if="isMobileMenuOpen" class="sm:hidden">
      <div class="pt-2 pb-3 space-y-1">
        <router-link 
          to="/" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path === '/' ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Home
        </router-link>
        <router-link 
          to="/projects" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path.startsWith('/projects') ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Projects
        </router-link>
        <router-link 
          to="/blog" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path.startsWith('/blog') ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Blog
        </router-link>
        <router-link 
          to="/contact" 
          class="block pl-3 pr-4 py-2 border-l-4"
          :class="[$route.path === '/contact' ? 'border-blue-500 text-blue-700 dark:text-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
          @click="isMobileMenuOpen = false"
        >
          Contact
        </router-link>
      </div>
      
      <div v-if="authStore.isAuthenticated" class="pt-4 pb-3 border-t border-gray-200 dark:border-gray-600">
        <div class="flex items-center px-4">
          <div class="flex-shrink-0">
            <div v-if="authStore.user?.avatar" class="h-10 w-10 rounded-full overflow-hidden">
              <img :src="authStore.user.avatar" alt="User avatar" class="h-full w-full object-cover" />
            </div>
            <div v-else class="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center text-white">
              {{ authStore.user?.name?.charAt(0) || authStore.user?.username?.charAt(0) || 'U' }}
            </div>
          </div>
          <div class="ml-3">
            <div class="text-base font-medium text-gray-800 dark:text-white">{{ authStore.user?.name }}</div>
            <div class="text-sm font-medium text-gray-500 dark:text-gray-300">{{ authStore.user?.email }}</div>
          </div>
        </div>
        <div class="mt-3 space-y-1">
          <router-link 
            v-if="authStore.isAdmin" 
            to="/admin" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Dashboard
          </router-link>
          <router-link 
            to="/admin/profile" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Profile
          </router-link>
          <button 
            @click="logout" 
            class="block w-full text-left px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            Sign out
          </button>
        </div>
      </div>
      
      <div v-else class="pt-4 pb-3 border-t border-gray-200 dark:border-gray-600">
        <div class="space-y-1">
          <router-link 
            to="/login" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Login
          </router-link>
          <router-link 
            to="/register" 
            class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="isMobileMenuOpen = false"
          >
            Register
          </router-link>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const isProfileOpen = ref(false);
const isMobileMenuOpen = ref(false);
const isDarkMode = ref(false);

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (isProfileOpen.value && !event.target.closest('.relative')) {
    isProfileOpen.value = false;
  }
};

onMounted(() => {
  // Check if dark mode is enabled in localStorage
  const darkMode = localStorage.getItem('darkMode') === 'true';
  isDarkMode.value = darkMode;
  
  if (darkMode) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
  
  document.addEventListener('click', handleClickOutside);
});

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value;
  localStorage.setItem('darkMode', isDarkMode.value);
  
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
};

const logout = () => {
  authStore.logout();
  isProfileOpen.value = false;
  isMobileMenuOpen.value = false;
  router.push('/');
};
</script>
```

### Create Footer Component

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/components/layout/Footer.vue
<template>
  <footer class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
    <div class="max-w-7xl mx-auto py-12 px-4 overflow-hidden sm:px-6 lg:px-8">
      <nav class="flex flex-wrap justify-center -mx-5 -my-2">
        <div class="px-5 py-2">
          <router-link to="/" class="text-base text-gray-500 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
            Home
          </router-link>
        </div>
        <div class="px-5 py-2">
          <router-link to="/projects" class="text-base text-gray-500 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
            Projects
          </router-link>
        </div>
        <div class="px-5 py-2">
          <router-link to="/blog" class="text-base text-gray-500 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
            Blog
          </router-link>
        </div>
        <div class="px-5 py-2">
          <router-link to="/contact" class="text-base text-gray-500 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
            Contact
          </router-link>
        </div>
      </nav>
      <div class="mt-8 flex justify-center space-x-6">
        <!-- Social media links - customize as needed -->
        <a href="#" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
          <span class="sr-only">GitHub</span>
          <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
          </svg>
        </a>
        <a href="#" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
          <span class="sr-only">Twitter</span>
          <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84" />
          </svg>
        </a>
        <a href="#" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
          <span class="sr-only">LinkedIn</span>
          <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 2.761 2.239 5 5 5h14c2.762 0 5-2.239 5-5v-14c0-2.761-2.238-5-5-5zm-11 19h-3v-11h3v11zm-1.5-12.268c-.966 0-1.75-.79-1.75-1.764s.784-1.764 1.75-1.764 1.75.79 1.75 1.764-.783 1.764-1.75 1.764zm13.5 12.268h-3v-5.604c0-3.368-4-3.113-4 0v5.604h-3v-11h3v1.765c1.396-2.586 7-2.777 7 2.476v6.759z" />
          </svg>
        </a>
      </div>
      <p class="mt-8 text-center text-base text-gray-400">
        &copy; {{ new Date().getFullYear() }} My Portfolio. All rights reserved.
      </p>
    </div>
  </footer>
</template>
```

### Create Main Layout

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/layouts/MainLayout.vue
<template>
  <div class="min-h-screen flex flex-col bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <Navbar />
    <main class="flex-grow">
      <slot></slot>
    </main>
    <Footer />
  </div>
</template>

<script setup>
import Navbar from '@/components/layout/Navbar.vue';
import Footer from '@/components/layout/Footer.vue';
</script>
```

### Create Admin Layout

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/layouts/AdminLayout.vue
<template>
  <div class="min-h-screen flex flex-col bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <Navbar />
    <div class="flex-grow flex">
      <!-- Sidebar -->
      <div class="hidden md:flex md:flex-shrink-0">
        <div class="flex flex-col w-64 border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
          <div class="h-0 flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
            <div class="flex items-center flex-shrink-0 px-4">
              <h2 class="text-lg font-semibold">Admin Dashboard</h2>
            </div>
            <nav class="mt-5 flex-1 px-2 bg-white dark:bg-gray-800 space-y-1">
              <router-link 
                to="/admin" 
                class="group flex items-center px-2 py-2 text-sm font-medium rounded-md"
                :class="[$route.path === '/admin' ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
              >
                <svg class="mr-3 h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                </svg>
                Dashboard
              </router-link>
              <router-link 
                to="/admin/projects" 
                class="group flex items-center px-2 py-2 text-sm font-medium rounded-md"
                :class="[$route.path.startsWith('/admin/projects') ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
              >
                <svg class="mr-3 h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
                Projects
              </router-link>
              <router-link 
                to="/admin/posts" 
                class="group flex items-center px-2 py-2 text-sm font-medium rounded-md"
                :class="[$route.path.startsWith('/admin/posts') ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
              >
                <svg class="mr-3 h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                </svg>
                Blog Posts
              </router-link>
              <router-link 
                to="/admin/contacts" 
                class="group flex items-center px-2 py-2 text-sm font-medium rounded-md"
                :class="[$route.path.startsWith('/admin/contacts') ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
              >
                <svg class="mr-3 h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                Contact Messages
              </router-link>
              <router-link 
                to="/admin/profile" 
                class="group flex items-center px-2 py-2 text-sm font-medium rounded-md"
                :class="[$route.path.startsWith('/admin/profile') ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']"
              >
                <svg class="mr-3 h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                Profile
              </router-link>
              <router-link 
                to="/" 
                class="group flex items-center px-2 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-md"
              >
                <svg class="mr-3 h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                </svg>
                Back to Site
              </router-link>
            </nav>
          </div>
        </div>
      </div>
      
      <!-- Mobile sidebar -->
      <div class="md:hidden fixed bottom-0 left-0 right-0 z-10 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
        <nav class="flex justify-around">
          <router-link 
            to="/admin" 
            class="flex flex-col items-center py-2 px-3"
            :class="[$route.path === '/admin' ? 'text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-400']"
          >
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
----

Let me continue with the mobile sidebar navigation and then create the main pages for the frontend.

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/layouts/AdminLayout.vue
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            <span class="text-xs mt-1">Dashboard</span>
          </router-link>
          <router-link 
            to="/admin/projects" 
            class="flex flex-col items-center py-2 px-3"
            :class="[$route.path.startsWith('/admin/projects') ? 'text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-400']"
          >
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <span class="text-xs mt-1">Projects</span>
          </router-link>
          <router-link 
            to="/admin/posts" 
            class="flex flex-col items-center py-2 px-3"
            :class="[$route.path.startsWith('/admin/posts') ? 'text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-400']"
          >
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
            </svg>
            <span class="text-xs mt-1">Posts</span>
          </router-link>
          <router-link 
            to="/admin/contacts" 
            class="flex flex-col items-center py-2 px-3"
            :class="[$route.path.startsWith('/admin/contacts') ? 'text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-400']"
          >
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            <span class="text-xs mt-1">Messages</span>
          </router-link>
          <router-link 
            to="/admin/profile" 
            class="flex flex-col items-center py-2 px-3"
            :class="[$route.path.startsWith('/admin/profile') ? 'text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-400']"
          >
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            <span class="text-xs mt-1">Profile</span>
          </router-link>
        </nav>
      </div>
      
      <!-- Main content -->
      <div class="flex-1 overflow-auto">
        <div class="py-6 px-4 sm:px-6 lg:px-8">
          <slot></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import Navbar from '@/components/layout/Navbar.vue';
</script>
```

### Create App.vue

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/App.vue
<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();

onMounted(() => {
  // Check if user is logged in and fetch profile
  if (authStore.isAuthenticated) {
    authStore.fetchProfile();
  }
});
</script>
```

### Create Main Pages

Let's create the main pages for the website:

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/views/HomePage.vue
<template>
  <MainLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Hero section -->
      <div class="py-12 md:py-20">
        <div class="text-center">
          <h1 class="text-4xl tracking-tight font-extrabold text-gray-900 dark:text-white sm:text-5xl md:text-6xl">
            <span class="block">Hi, I'm</span>
            <span class="block text-blue-600 dark:text-blue-400">Your Name</span>
          </h1>
          <p class="mt-3 max-w-md mx-auto text-base text-gray-500 dark:text-gray-300 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl">
            I'm a full-stack developer passionate about building web applications with modern technologies.
          </p>
          <div class="mt-5 max-w-md mx-auto sm:flex sm:justify-center md:mt-8">
            <div class="rounded-md shadow">
              <router-link to="/projects" class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 md:py-4 md:text-lg md:px-10">
                View My Work
              </router-link>
            </div>
            <div class="mt-3 rounded-md shadow sm:mt-0 sm:ml-3">
              <router-link to="/contact" class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-blue-600 bg-white hover:bg-gray-50 dark:text-blue-400 dark:bg-gray-800 dark:hover:bg-gray-700 md:py-4 md:text-lg md:px-10">
                Contact Me
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Featured projects section -->
      <div class="py-12">
        <div class="text-center">
          <h2 class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl">
            Featured Projects
          </h2>
          <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 dark:text-gray-300 sm:mt-4">
            Check out some of my recent work
          </p>
        </div>

        <div class="mt-10" v-if="projects.length > 0">
          <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            <div v-for="project in projects" :key="project.id" class="flex flex-col rounded-lg shadow-lg overflow-hidden bg-white dark:bg-gray-800">
              <div class="flex-shrink-0">
                <img v-if="project.imageURL" :src="project.imageURL" alt="" class="h-48 w-full object-cover">
                <div v-else class="h-48 w-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                  <svg class="h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </div>
              </div>
              <div class="flex-1 p-6 flex flex-col justify-between">
                <div class="flex-1">
                  <p class="text-sm font-medium text-blue-600 dark:text-blue-400">
                    {{ project.tags ? project.tags.split(',')[0] : 'Project' }}
                  </p>
                  <router-link :to="`/projects/${project.id}`" class="block mt-2">
                    <p class="text-xl font-semibold text-gray-900 dark:text-white">{{ project.title }}</p>
                    <p class="mt-3 text-base text-gray-500 dark:text-gray-300">{{ project.description }}</p>
                  </router-link>
                </div>
                <div class="mt-6 flex items-center">
                  <div class="flex space-x-2">
                    <a v-if="project.projectURL" :href="project.projectURL" target="_blank" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white">
                      <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                      </svg>
                    </a>
                    <a v-if="project.githubURL" :href="project.githubURL" target="_blank" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white">
                      <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                        <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
                      </svg>
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-10 text-center">
            <router-link to="/projects" class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700">
              View All Projects
            </router-link>
          </div>
        </div>
        <div v-else class="mt-10 text-center text-gray-500 dark:text-gray-400">
          <p>No projects to display yet.</p>
        </div>
      </div>

      <!-- Latest blog posts section -->
      <div class="py-12">
        <div class="text-center">
          <h2 class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl">
            Latest Blog Posts
          </h2>
          <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 dark:text-gray-300 sm:mt-4">
            Thoughts, ideas, and tutorials
          </p>
        </div>

        <div class="mt-10" v-if="posts.length > 0">
          <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            <div v-for="post in posts" :key="post.id" class="flex flex-col rounded-lg shadow-lg overflow-hidden bg-white dark:bg-gray-800">
              <div class="flex-shrink-0">
                <img v-if="post.imageURL" :src="post.imageURL" alt="" class="h-48 w-full object-cover">
                <div v-else class="h-48 w-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                  <svg class="h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                  </svg>
                </div>
              </div>
              <div class="flex-1 p-6 flex flex-col justify-between">
                <div class="flex-1">
                  <p class="text-sm font-medium text-blue-600 dark:text-blue-400">
                    {{ post.tags ? post.tags.split(',')[0] : 'Blog' }}
                  </p>
                  <router-link :to="`/blog/${post.slug}`" class="block mt-2">
                    <p class="text-xl font-semibold text-gray-900 dark:text-white">{{ post.title }}</p>
                    <p class="mt-3 text-base text-gray-500 dark:text-gray-300">{{ post.excerpt }}</p>
                  </router-link>
                </div>
                <div class="mt-6 flex items-center">
                  <div class="flex-shrink-0">
                    <span class="sr-only">{{ post.user?.name || 'Author' }}</span>
                    <div class="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center text-white">
                      {{ post.user?.name?.charAt(0) || 'A' }}
                    </div>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900 dark:text-white">
                      {{ post.user?.name || 'Author' }}
                    </p>
                    <div class="flex space-x-1 text-sm text-gray-500 dark:text-gray-400">
                      <time :datetime="post.publishedAt">
                        {{ formatDate(post.publishedAt) }}
                      </time>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-10 text-center">
            <router-link to="/blog" class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700">
              View All Posts
            </router-link>
          </div>
        </div>
        <div v-else class="mt-10 text-center text-gray-500 dark:text-gray-400">
          <p>No blog posts to display yet.</p>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import MainLayout from '@/layouts/MainLayout.vue';
import api from '@/services/api';

const projects = ref([]);
const posts = ref([]);

onMounted(async () => {
  try {
    // Fetch featured projects
    const projectsResponse = await api.getProjects({ featured: true });
    projects.value = projectsResponse.data.projects.slice(0, 3);
    
    // Fetch latest blog posts
    const postsResponse = await api.getPosts({ limit: 3 });
    posts.value = postsResponse.data.posts;
  } catch (error) {
    console.error('Error fetching data:', error);
  }
});

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
};
</script>
```

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/views/LoginPage.vue
<template>
  <MainLayout>
    <div class="min-h-full flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
            Sign in to your account
          </h2>
          <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
            Or
            <router-link to="/register" class="font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400 dark:hover:text-blue-300">
              create a new account
            </router-link>
          </p>
        </div>
        <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
          <div v-if="error" class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800 dark:text-red-200">{{ error }}</h3>
              </div>
            </div>
          </div>
          
          <div class="rounded-md shadow-sm -space-y-px">
            <div>
              <label for="username" class="sr-only">Username</label>
              <input id="username" name="username" type="text" autocomplete="username" required v-model="username" class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm" placeholder="Username">
            </div>
            <div>
              <label for="password" class="sr-only">Password</label>
              <input id="password" name="password" type="password" autocomplete="current-password" required v-model="password" class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm" placeholder="Password">
            </div>
          </div>

          <div>
            <button type="submit" :disabled="loading" class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed">
              <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
                <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </span>
              <span v-else class="absolute left-0 inset-y-0 flex items-center pl-3">
                <svg class="h-5 w-5 text-blue-500 group-hover:text-blue-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
              </span>
              Sign in
            </button>
          </div>
        </form>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import MainLayout from '@/layouts/MainLayout.vue';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');

const handleLogin = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    await authStore.login({
      username: username.value,
      password: password.value
    });
    
    // Redirect to admin dashboard if admin, otherwise to home
    if (authStore.isAdmin) {
      router.push('/admin');
    } else {
      router.push('/');
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to login. Please check your credentials.';
  } finally {
    loading.value = false;
  }
};
</script>
```

Let's create the contact page as well:

```vue:/Users/walter.fan/workspace/walter/idaas/frontend/src/views/ContactPage.vue
<template>
  <MainLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="max-w-3xl mx-auto">
        <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white sm:text-4xl">
          Get in Touch
        </h1>
        <p class="mt-4 text-lg text-gray-500 dark:text-gray-300">
          Have a question or want to work together? Fill out the form below and I'll get back to you as soon as possible.
        </p>
        
        <div class="mt-12">
          <form @submit.prevent="submitForm" class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-8">
            <div class="sm:col-span-2">
              <div v-if="formStatus.success" class="rounded-md bg-green-50 dark:bg-green-900/20 p-4 mb-6">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg class="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-green-800 dark:text-green-200">
                      Your message has been sent successfully! I'll get back to you soon.
                    </p>
                  </div>
                </div>
              </div>
              
              <div v-if="formStatus.error" class="rounded-md bg-red-50 dark:bg-red-900/20 p-4 mb-6">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414


-----
I'll continue adding tasks to the fabfile.py to help with the idaas project. Let me add tasks for setting up the backend and frontend, as well as deployment tasks.

```python:/Users/walter.fan/workspace/walter/idaas/fabfile.py
go mod init github.com/walterfan/idaas/backend
mkdir -p internal/api/handlers internal/api/middleware internal/config internal/models
mkdir -p cmd/server
touch .env
    """
    for cmd in cmds.strip().split('\n'):
        run_cmd(c, cmd)

@task(hosts=default_hosts)
def init_frontend(c):
    cmds = """
mkdir -p /Users/walter.fan/workspace/walter/idaas/frontend
cd /Users/walter.fan/workspace/walter/idaas/frontend
npm init vue@latest .
npm install
npm install axios tailwindcss postcss autoprefixer @headlessui/vue @heroicons/vue
npm install -D @tailwindcss/forms @tailwindcss/typography
npx tailwindcss init -p
    """
    for cmd in cmds.strip().split('\n'):
        run_cmd(c, cmd)

@task(hosts=default_hosts)
def run_backend(c):
    cmd = "cd /Users/walter.fan/workspace/walter/idaas/backend && go run cmd/server/main.go"
    c.local(cmd)

@task(hosts=default_hosts)
def run_frontend(c):
    cmd = "cd /Users/walter.fan/workspace/walter/idaas/frontend && npm run dev"
    c.local(cmd)

@task(hosts=default_hosts)
def build_frontend(c):
    cmd = "cd /Users/walter.fan/workspace/walter/idaas/frontend && npm run build"
    c.local(cmd)

@task(hosts=default_hosts)
def build_backend(c):
    cmd = "cd /Users/walter.fan/workspace/walter/idaas/backend && go build -o bin/server cmd/server/main.go"
    c.local(cmd)

@task(hosts=default_hosts)
def test_backend(c):
    cmd = "cd /Users/walter.fan/workspace/walter/idaas/backend && go test ./..."
    c.local(cmd)

@task(hosts=default_hosts)
def deploy(c):
    # Build frontend
    build_frontend(c)
    
    # Build backend
    build_backend(c)
    
    # Deploy to server
    c.local("rsync -avz --delete /Users/walter.fan/workspace/walter/idaas/frontend/dist/ walter@www.fanyamin.com:/var/www/html/identity-service/")
    c.local("rsync -avz /Users/walter.fan/workspace/walter/idaas/backend/bin/server walter@www.fanyamin.com:/opt/identity-service/")
    
    # Restart service on remote server
    conn.run("sudo systemctl restart identity-service")

@task(hosts=default_hosts)
def setup_db(c):
    # Create database schema
    cmd = "cd /Users/walter.fan/workspace/walter/idaas/backend && go run cmd/migrations/main.go"
    c.local(cmd)

@task(hosts=default_hosts)
def generate_api_docs(c):
    cmd = "cd /Users/walter.fan/workspace/walter/idaas/backend && swag init -g cmd/server/main.go -o api/docs"
    c.local(cmd)

@task(hosts=default_hosts)
def clean(c):
    # Clean build artifacts
    c.local("rm -rf /Users/walter.fan/workspace/walter/idaas/frontend/dist")
    c.local("rm -rf /Users/walter.fan/workspace/walter/idaas/backend/bin")
```

These tasks will help you manage the development workflow for your idaas project:

1. `init_backend` - Sets up the Go backend directory structure
2. `init_frontend` - Initializes the Vue.js frontend project
3. `run_backend` - Runs the backend server locally
4. `run_frontend` - Runs the frontend development server
5. `build_frontend` - Builds the frontend for production
6. `build_backend` - Compiles the backend Go code
7. `test_backend` - Runs backend tests
8. `deploy` - Builds and deploys both frontend and backend to your server
9. `setup_db` - Runs database migrations
10. `generate_api_docs` - Generates API documentation using Swagger
11. `clean` - Removes build artifacts

You can run these tasks using the `fab` command, for example:
```
fab init-backend
fab run-frontend
```