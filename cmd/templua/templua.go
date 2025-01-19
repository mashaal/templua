package main

import (
	"fmt"
	"log"
	"net/http"
	"templua/templates"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	lt := templates.NewLuaTemplate()
	defer lt.Close()

	// Homepage handler
	e.GET("/", func(c echo.Context) error {
		// read the homepage template
		templatePath := filepath.Join("templates", "home.lua")
		templateBytes, err := os.ReadFile(templatePath)
		if err != nil {
			log.Printf("Failed to read template: %v", err)
			return fmt.Errorf("failed to read template: %v", err)
		}
		template := string(templateBytes)
		log.Printf("Template content:\n%s", template)

		vars := map[string]interface{}{
			"heading": "Welcome to Dynamic Templua!",
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
