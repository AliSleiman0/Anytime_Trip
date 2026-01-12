package admin

import (
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// CMSTravelExperience renders the homepage travel experience management page
func (h *AdminHandler) CMSTravelExperience(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-home-travel.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Homepage - Travel Experience",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSTravelExperienceFragment serves the travel experience fragment for HTMX partial loads
func (h *AdminHandler) GetCMSTravelExperienceFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-travel-frag.html")))
}

// CMSBanner renders the homepage banner management page
func (h *AdminHandler) CMSBanner(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-home-banner.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Homepage - Banner",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSBannerFragment serves the banner fragment for HTMX partial loads
func (h *AdminHandler) GetCMSBannerFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-banner-frag.html")))
}

// CMSPopularLocations renders the homepage popular locations management page
func (h *AdminHandler) CMSPopularLocations(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-home-popular.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Homepage - Popular Locations",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSPopularLocationsFragment serves the popular locations fragment for HTMX partial loads
func (h *AdminHandler) GetCMSPopularLocationsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-popular-frag.html")))
}
