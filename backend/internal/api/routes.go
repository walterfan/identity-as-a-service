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