package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"

	"go_pco_arrivals/internal/config"
	"go_pco_arrivals/internal/database"
	"go_pco_arrivals/internal/handlers"
	"go_pco_arrivals/internal/middleware"
	"go_pco_arrivals/internal/services"
	"go_pco_arrivals/internal/utils"
)

func main() {
	// Validate environment variables
	if err := validateEnvironment(); err != nil {
		log.Fatal("Environment validation failed:", err)
	}

	// Initialize logger
	appLogger := utils.NewLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", "error", err)
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		appLogger.Fatal("Failed to run database migrations", "error", err)
	}

	// Get the appropriate database connection for services
	var gormDB *gorm.DB
	if db.GetType() == database.SQLiteDB {
		sqliteDB := db.(*database.SQLiteDatabase)
		gormDB = sqliteDB.GetGormDB()
		// Configure database connection pool for SQLite
		if err := database.ConfigureConnectionPool(gormDB, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, time.Duration(cfg.Database.ConnMaxLifetime)*time.Second); err != nil {
			appLogger.Fatal("Failed to configure database pool", "error", err)
		}
	}

	// Initialize logger
	logger := utils.NewLogger()

	// Initialize services
	pcoService := services.NewPCOService(cfg, gormDB, logger)
	authService := services.NewAuthService(cfg, gormDB, logger, pcoService)
	notificationService := services.NewNotificationService(gormDB, pcoService)
	billboardService := services.NewBillboardService(cfg, gormDB, logger, pcoService, nil) // TODO: Add WebSocket service

	// Initialize WebSocket hub
	wsHub := services.NewWebSocketHub()
	go wsHub.Run()

	// Initialize cleanup service
	cleanupService := services.NewCleanupService(db, notificationService)
	go cleanupService.Start()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "PCO Arrivals Billboard",
		ServerHeader: "PCO-Arrivals-Billboard/1.0",
		ErrorHandler: middleware.ErrorHandler,
	})

	// Add middleware
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.PerformanceMonitoring())
	app.Use(middleware.Compression())
	app.Use(middleware.ValidateInput())
	app.Use(middleware.SanitizeInput())
	app.Use(middleware.RateLimit())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.Server.CORSOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))
	// TODO: Add proper logging middleware
	app.Use(recover.New())

	// Add auth service to app context for middleware access
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("auth_service", authService)
		return c.Next()
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, gormDB, logger, authService, pcoService)
	apiHandler := handlers.NewAPIHandler(gormDB, pcoService, notificationService, billboardService, wsHub, logger)
	staticHandler := handlers.NewStaticHandler()
	websocketHandler := handlers.NewWebSocketHandler(wsHub, authService)
	healthHandler := handlers.NewHealthHandler(gormDB)
	billboardHandler := handlers.NewBillboardHandler(cfg, gormDB, logger, billboardService, pcoService)

	// Setup routes
	setupRoutes(app, authHandler, apiHandler, staticHandler, websocketHandler, healthHandler, billboardHandler)

	// Start server
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Stop WebSocket hub
	wsHub.Stop()

	// Stop cleanup service
	cleanupService.Stop()

	// Close database connection
	if err := db.Close(); err != nil {
		appLogger.Error("Failed to close database connection", "error", err)
	}

	appLogger.Info("Server stopped")
}

func validateEnvironment() error {
	required := []string{"PCO_CLIENT_ID", "PCO_CLIENT_SECRET", "PCO_REDIRECT_URI"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("required environment variable %s is not set", env)
		}
	}
	return nil
}

func setupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, apiHandler *handlers.APIHandler, staticHandler *handlers.StaticHandler, websocketHandler *handlers.WebSocketHandler, healthHandler *handlers.HealthHandler, billboardHandler *handlers.BillboardHandler) {
	// Health check
	app.Get("/health", healthHandler.Health)
	app.Get("/health/detailed", healthHandler.DetailedHealth)

	// Authentication routes
	auth := app.Group("/auth")
	auth.Get("/login", authHandler.InitiateOAuth)
	auth.Get("/callback", authHandler.OAuthCallback)
	auth.Get("/status", authHandler.GetAuthStatus)
	auth.Post("/logout", authHandler.Logout)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Get("/profile", authHandler.GetUserProfile)
	auth.Put("/profile", authHandler.UpdateUserProfile)

	// API routes
	api := app.Group("/api", middleware.RequireAuth())
	api.Get("/events", apiHandler.GetEvents)
	api.Get("/events/:id", apiHandler.GetEvent)
	api.Post("/events", middleware.RequireAdmin(), apiHandler.CreateEvent)
	api.Put("/events/:id", middleware.RequireAdmin(), apiHandler.UpdateEvent)
	api.Delete("/events/:id", middleware.RequireAdmin(), apiHandler.DeleteEvent)

	api.Get("/notifications", apiHandler.GetNotifications)
	api.Post("/notifications", apiHandler.CreateNotification)
	api.Delete("/notifications/:id", apiHandler.DeleteNotification)

	// Admin-specific routes
	api.Get("/notifications/active", apiHandler.GetNotifications)
	api.Get("/security-codes", apiHandler.GetSecurityCodes)
	api.Post("/security-codes", apiHandler.AddSecurityCode)
	api.Delete("/security-codes/:code", apiHandler.RemoveSecurityCode)
	api.Get("/billboard/control", apiHandler.GetBillboardControl)
	api.Post("/billboard/launch", apiHandler.LaunchBillboard)
	api.Post("/billboard/clear", apiHandler.ClearBillboard)
	api.Get("/billboard/stats/:locationId", apiHandler.GetCheckInStats)

	api.Get("/check-ins", apiHandler.GetCheckIns)
	api.Get("/check-ins/location/:locationId", apiHandler.GetCheckInsByLocation)
	api.Get("/check-ins/event/:eventId", apiHandler.GetCheckInsByEvent)

	api.Get("/locations", apiHandler.GetLocations)
	api.Get("/locations/:id", apiHandler.GetLocation)
	api.Get("/locations/:locationId/status", apiHandler.GetLocationStatus)
	api.Get("/locations/:locationId/analytics", apiHandler.GetLocationAnalytics)
	api.Get("/locations/overview", apiHandler.GetLocationsOverview)

	// Test endpoint for WebSocket broadcasts (development only)
	app.Get("/test/websocket", apiHandler.TestWebSocketBroadcast)

	// Billboard routes
	billboard := app.Group("/billboard")
	billboard.Get("/state/:locationID", billboardHandler.GetBillboardState)
	billboard.Get("/check-ins/:locationID", billboardHandler.GetRecentCheckIns)
	billboard.Get("/stats/:locationID", billboardHandler.GetCheckInStats)
	billboard.Post("/sync/:locationID", billboardHandler.SyncPCOCheckIns)
	billboard.Get("/locations", billboardHandler.GetLocations)
	billboard.Post("/locations", billboardHandler.AddLocation)
	billboard.Get("/location/:locationID", billboardHandler.GetLocationBillboard)
	billboard.Post("/cleanup", billboardHandler.CleanupOldData)
	billboard.Get("/status", billboardHandler.GetSystemStatus)

	// Location-specific billboard
	billboard.Get("/location/:locationId", billboardHandler.GetLocationBillboard)

	// WebSocket routes
	app.Get("/ws", websocket.New(websocketHandler.HandleWebSocket))
	app.Get("/ws/billboard/:locationId", websocket.New(websocketHandler.HandleBillboardWebSocket))

	// Static files
	app.Get("/", staticHandler.ServeIndex)
	app.Get("/admin", middleware.RequireAuth(), staticHandler.ServeAdmin)
	app.Get("/billboard", staticHandler.ServeBillboard)
	app.Get("/location/:locationId", staticHandler.ServeLocationBillboard)
	app.Get("/login", staticHandler.ServeLogin)
	app.Get("/offline", staticHandler.ServeOffline)

	// Static assets
	app.Static("/static", "./web/static")
	app.Get("/manifest.json", staticHandler.ServeManifest)
	app.Get("/sw.js", staticHandler.ServeServiceWorker)

	// 404 handler
	app.Use(staticHandler.Handle404)
}
