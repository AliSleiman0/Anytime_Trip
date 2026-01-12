package admin

import (
	"html/template"
	"net/url"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// Error renders the admin error page
func (h *AdminHandler) Error(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "error.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	// Get error message from query parameter, default to generic message
	encodedMsg := c.Query("error")
	errorMsg, err := url.QueryUnescape(encodedMsg)
	if err != nil || errorMsg == "" {
		errorMsg = "We're sorry, but an error occurred while processing your request. Please try again later or contact support if the problem persists."
	}

	data := fiber.Map{
		"Title":               "Error",
		"ErrorMessage":        errorMsg,
		"EncodedErrorMessage": encodedMsg,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetErrorFragment serves the error fragment for HTMX partial loads
func (h *AdminHandler) GetErrorFragment(c *fiber.Ctx) error {
	// Get error message from query parameter, default to generic message
	encodedMsg := c.Query("error")
	errorMsg, err := url.QueryUnescape(encodedMsg)
	if err != nil || errorMsg == "" {
		errorMsg = "We're sorry, but an error occurred while processing your request. Please try again later or contact support if the problem persists."
	}

	data := fiber.Map{
		"ErrorMessage": errorMsg,
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "error-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
