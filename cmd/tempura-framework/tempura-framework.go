package main

import (
	"fmt"
	"log"
	"net/http"
	"templua/framework"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/mashaal/tempura"
)

func main() {
	e := echo.New()
	lt := tempura.NewLuaTemplate()
	defer lt.Close()

	// Register custom components
	lt.RegisterCustomComponent("Card", "framework/components/card.lua")

	// Initialize live reload
	lr, err := framework.NewLiveReload()
	if err != nil {
		log.Fatalf("Failed to initialize live reload: %v", err)
	}
	defer lr.Close()

	// Get absolute path to framework directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	frameworkDir := filepath.Join(wd, "framework")
	log.Printf("Watching framework directory: %s", frameworkDir)

	// Watch framework directory
	if err := lr.WatchDir(frameworkDir); err != nil {
		log.Fatalf("Failed to watch framework directory: %v", err)
	}

	// WebSocket endpoint for live reload
	e.GET("/ws", lr.HandleWebSocket)

	// Homepage handler
	e.GET("/", func(c echo.Context) error {
		// read the homepage template
		templatePath := filepath.Join("framework", "home.lua")
		templateBytes, err := os.ReadFile(templatePath)
		if err != nil {
			log.Printf("Failed to read template: %v", err)
			return fmt.Errorf("failed to read template: %v", err)
		}
		template := string(templateBytes)
		log.Printf("Template content:\n%s", template)

		vars := map[string]interface{}{
			"heading": "Welcome to Templua",
		}

		html, err := lt.RenderHTMLWithVars(template, vars)
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			return fmt.Errorf("failed to render template: %v", err)
		}

		log.Printf("Final rendered HTML:\n%s", html)
		return c.HTML(http.StatusOK, html)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
