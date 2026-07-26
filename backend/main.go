package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"skillgap-ai/handlers"
	"skillgap-ai/services"
)

func main() {
	_ = godotenv.Load()

	if key := os.Getenv("MISTRAL_API_KEY"); key == "" || key == "your_mistral_api_key_here" {
		log.Fatal("❌ MISTRAL_API_KEY is missing or still set to the placeholder value. " +
			"Get a real key at https://console.mistral.ai/ and put it in backend/.env, then restart.")
	}

	dbPath := os.Getenv("BADGER_DB_PATH")
	if dbPath == "" {
		dbPath = "./data/badger"
	}

	store, err := services.OpenStore(dbPath)
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}
	defer store.Close()
	log.Printf("✅ BadgerDB ready — data dir: %s", dbPath)

	app := fiber.New(fiber.Config{
		AppName:      "SkillGap AI API",
		BodyLimit:    10 * 1024 * 1024, // 10MB, enough for a resume PDF
		ErrorHandler: fiberErrorHandler,
	})

	app.Use(logger.New())

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET,POST,DELETE,OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "SkillGap AI running", "version": "2.0.0"})
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	analyzeHandler := handlers.NewAnalyzeHandler(store)
	historyHandler := handlers.NewHistoryHandler(store)

	api := app.Group("/api")
	api.Post("/analyze", analyzeHandler.Analyze)
	api.Get("/history", historyHandler.GetHistory)
	api.Get("/history/:id", historyHandler.GetAnalysis)
	api.Delete("/history/:id", historyHandler.DeleteAnalysis)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	// Graceful shutdown: on Ctrl+C / SIGTERM, stop accepting new requests and
	// close BadgerDB cleanly (the deferred store.Close() above only fires if
	// main() returns normally, which this makes true instead of the process
	// being killed mid-write).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down…")
	_ = app.ShutdownWithTimeout(5 * time.Second)
}

func fiberErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{"detail": err.Error()})
}