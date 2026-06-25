package admin

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	adminmodel "travel/backend/internal/models/admin"

	"github.com/gofiber/fiber/v2"
)

// GetHomeTravelFragment serves the Travel Experience fragment
func (h *AdminHandler) GetHomeTravelFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-travel-frag.html")))
}

// UploadTravelImage accepts image uploads via multipart/form-data (HTMX file upload)
func (h *AdminHandler) UploadTravelImage(c *fiber.Ctx) error {
	m, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).SendString("invalid multipart form")
	}

	// Determine suffix for this upload: prefer HX header, fallback to any file field present
	suffix := c.Get("X-Image-Slot")
	if suffix == "" {
		for fieldName := range m.File {
			if strings.HasPrefix(fieldName, "travel_image_") {
				suffix = strings.TrimPrefix(fieldName, "travel_image_")
				break
			}
		}
	}
	if suffix == "" {
		suffix = "1"
	}

	// Fetch only the file for the target slot to avoid reusing previous selections
	fieldName := "travel_image_" + suffix
	fh, err := c.FormFile(fieldName)
	if err != nil || fh == nil {
		// Fallback to any available file but keep the computed suffix
		for _, fhArr := range m.File {
			if len(fhArr) > 0 {
				fh = fhArr[0]
				break
			}
		}
		if fh == nil {
			return c.Status(400).SendString("no file uploaded")
		}
	}

	// Only accept image files by extension
	ext := strings.ToLower(path.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg":
		// allowed
	default:
		return c.Status(400).SendString("only image files allowed")
	}

	// Save file with stable name: travel-expirience-image-<suffix><ext>
	dstDir := filepath.Clean("static/admin/images")
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return c.Status(500).SendString("unable to create image folder")
	}
	filename := fmt.Sprintf("travel-expirience-image-%s%s", suffix, ext)
	dstPath := filepath.Join(dstDir, filename)
	if err := c.SaveFile(fh, dstPath); err != nil {
		return c.Status(500).SendString("failed to save file")
	}

	// Prepare image url exposed by static server
	imageURL := "/static/admin/images/" + filename

	// Reply with background-image div with checkmark overlay to indicate upload success
	html := fmt.Sprintf(`<div id="current-image-%s" class="w-full h-48 rounded-xl overflow-hidden border border-[#CAD9EA] bg-[#e7eef6] relative" style="background-image:url('%s'); background-size:cover; background-position:center;"><div class="absolute top-2 right-2 bg-green-500 rounded-full p-1.5"><svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg></div></div><input type="hidden" name="travel_image_path_%s" value="%s" />`, suffix, imageURL, suffix, imageURL)
	return c.Status(200).SendString(html)
}

// GetHomeBannerFragment serves the Banner fragment
func (h *AdminHandler) GetHomeBannerFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-banner-frag.html")))
}

// UploadBannerImage accepts image uploads via multipart/form-data (HTMX file upload)
func (h *AdminHandler) UploadBannerImage(c *fiber.Ctx) error {
	m, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).SendString("invalid multipart form")
	}

	// Determine suffix for this upload: prefer HX header, fallback to any file field present
	suffix := c.Get("X-Banner-Slot")
	if suffix == "" {
		for fieldName := range m.File {
			if strings.HasPrefix(fieldName, "banner_image_") {
				suffix = strings.TrimPrefix(fieldName, "banner_image_")
				break
			}
		}
	}
	if suffix == "" {
		suffix = "1"
	}

	// Fetch only the file for the target slot to avoid reusing previous selections
	fieldName := "banner_image_" + suffix
	fh, err := c.FormFile(fieldName)
	if err != nil || fh == nil {
		// Fallback to any available file but keep the computed suffix
		for _, fhArr := range m.File {
			if len(fhArr) > 0 {
				fh = fhArr[0]
				break
			}
		}
		if fh == nil {
			return c.Status(400).SendString("no file uploaded")
		}
	}

	// Only accept image files by extension
	ext := strings.ToLower(path.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg":
		// allowed
	default:
		return c.Status(400).SendString("only image files allowed")
	}

	// Save file with stable name: banner-image-<suffix><ext>
	dstDir := filepath.Clean("static/admin/images")
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return c.Status(500).SendString("unable to create image folder")
	}
	filename := fmt.Sprintf("banner-image-%s%s", suffix, ext)
	dstPath := filepath.Join(dstDir, filename)
	if err := c.SaveFile(fh, dstPath); err != nil {
		return c.Status(500).SendString("failed to save file")
	}

	// Prepare image url exposed by static server
	imageURL := "/static/admin/images/" + filename

	// Reply with background-image div with checkmark overlay to indicate upload success
	html := fmt.Sprintf(`<div id="current-banner-%s" class="w-full h-48 rounded-xl overflow-hidden border border-[#CAD9EA] bg-[#e7eef6] relative" style="background-image:url('%s'); background-size:cover; background-position:center;"><div class="absolute top-2 right-2 bg-green-500 rounded-full p-1.5"><svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg></div></div><input type="hidden" name="banner_image_path_%s" value="%s" />`, suffix, imageURL, suffix, imageURL)
	return c.Status(200).SendString(html)
}

// SaveBannerImage saves banner images via multipart form with HTMX
func (h *AdminHandler) SaveBannerImage(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	form, err := c.MultipartForm()
	if err != nil {
		// fallback to plain form for single slot
		title := c.FormValue("banner_title_1")
		imgPath := c.FormValue("banner_image_path_1")
		if title == "" && imgPath == "" {
			return c.Status(400).JSON(fiber.Map{"error": "no data provided"})
		}
		if h.bannerRepo != nil {
			item := &adminmodel.Banner{ID: "banner_1", Title: title, ImagePath: imgPath}
			if err := h.bannerRepo.UpsertByID(ctx, item); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "db upsert failed"})
			}
		}
		return c.JSON(fiber.Map{"status": "saved"})
	}

	if form == nil {
		return c.Status(400).JSON(fiber.Map{"error": "no form data"})
	}

	// collect suffixes from incoming form keys
	suffixSet := map[int]struct{}{}
	for key := range form.Value {
		if strings.HasPrefix(key, "banner_title_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "banner_title_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
		if strings.HasPrefix(key, "banner_image_path_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "banner_image_path_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
		if strings.HasPrefix(key, "banner_image_delete_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "banner_image_delete_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
	}

	if len(suffixSet) == 0 {
		suffixSet[1] = struct{}{}
	}

	// sort suffixes
	suffixes := make([]int, 0, len(suffixSet))
	for idx := range suffixSet {
		suffixes = append(suffixes, idx)
	}
	sort.Ints(suffixes)

	getVal := func(key string) string {
		if vals, ok := form.Value[key]; ok && len(vals) > 0 {
			return vals[0]
		}
		return c.FormValue(key)
	}

	for _, idx := range suffixes {
		suffix := fmt.Sprintf("%d", idx)
		titleKey := "banner_title_" + suffix
		pathKey := "banner_image_path_" + suffix
		deleteKey := "banner_image_delete_" + suffix

		title := getVal(titleKey)
		img := getVal(pathKey)
		deleteFlag := strings.ToLower(strings.TrimSpace(getVal(deleteKey)))
		shouldDelete := deleteFlag == "1" || deleteFlag == "true" || deleteFlag == "on"

		id := "banner_" + suffix

		var existingItem *adminmodel.Banner
		if h.bannerRepo != nil {
			existingItem, _ = h.bannerRepo.FindByID(ctx, id)
		}

		if shouldDelete {
			if h.bannerRepo != nil {
				if err := h.bannerRepo.DeleteByID(ctx, id); err != nil {
					return c.Status(500).JSON(fiber.Map{"error": "db delete failed"})
				}
			}
			continue
		}

		if existingItem != nil {
			if title == "" && existingItem.Title != "" {
				title = existingItem.Title
			}
			if img == "" && existingItem.ImagePath != "" {
				img = existingItem.ImagePath
			}
		}

		if title != "" || img != "" {
			if h.bannerRepo != nil {
				item := &adminmodel.Banner{ID: id, Title: title, ImagePath: img}
				if err := h.bannerRepo.UpsertByID(ctx, item); err != nil {
					return c.Status(500).JSON(fiber.Map{"error": "db upsert failed"})
				}
			}
		}
	}

	return c.JSON(fiber.Map{"status": "saved"})
}

// GetHomePopularFragment serves the Popular Locations fragment
func (h *AdminHandler) GetHomePopularFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-popular-frag.html")))
}

// UploadPopularImage accepts image uploads for popular locations via multipart/form-data (HTMX file upload)
func (h *AdminHandler) UploadPopularImage(c *fiber.Ctx) error {
	m, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).SendString("invalid multipart form")
	}

	// Determine suffix for this upload: prefer HX header, then form value, fallback to any file field present
	suffix := c.Get("X-Popular-Slot")
	if suffix == "" {
		if vals, ok := m.Value["popular-slot"]; ok && len(vals) > 0 {
			suffix = vals[0]
		}
	}
	if suffix == "" {
		for fieldName := range m.File {
			if strings.HasPrefix(fieldName, "popular_image_") {
				suffix = strings.TrimPrefix(fieldName, "popular_image_")
				break
			}
		}
	}
	if suffix == "" {
		suffix = "1"
	}

	fieldName := "popular_image_" + suffix
	fh, err := c.FormFile(fieldName)
	if err != nil || fh == nil {
		for _, fhArr := range m.File {
			if len(fhArr) > 0 {
				fh = fhArr[0]
				break
			}
		}
		if fh == nil {
			return c.Status(400).SendString("no file uploaded")
		}
	}

	ext := strings.ToLower(path.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg":
	default:
		return c.Status(400).SendString("only image files allowed")
	}

	dstDir := filepath.Clean("static/admin/images")
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return c.Status(500).SendString("unable to create image folder")
	}
	// stable name per request
	filename := fmt.Sprintf("popular-location-image-%s%s", suffix, ext)
	dstPath := filepath.Join(dstDir, filename)
	if err := c.SaveFile(fh, dstPath); err != nil {
		return c.Status(500).SendString("failed to save file")
	}

	imageURL := "/static/admin/images/" + filename

	html := fmt.Sprintf(`<div id="current-popular-%s" class="w-full h-48 rounded-xl overflow-hidden border border-[#CAD9EA] bg-[#e7eef6] relative" style="background-image:url('%s'); background-size:cover; background-position:center;"><div class="absolute top-2 right-2 bg-green-500 rounded-full p-1.5"><svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg></div></div><input type="hidden" name="popular_image_path_%s" value="%s" />`, suffix, imageURL, suffix, imageURL)
	return c.Status(200).SendString(html)
}

// SavePopular saves popular locations via multipart form with HTMX
func (h *AdminHandler) SavePopular(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	form, err := c.MultipartForm()
	if err != nil {
		// fallback to plain form for single slot
		title := c.FormValue("popular_title_1")
		location := c.FormValue("popular_location_1")
		imgPath := c.FormValue("popular_image_path_1")
		if title == "" && location == "" && imgPath == "" {
			return c.Status(400).JSON(fiber.Map{"error": "no data provided"})
		}
		if h.popularRepo != nil {
			item := &adminmodel.PopularLocation{ID: "popular_1", Title: title, Location: location, ImagePath: imgPath}
			if err := h.popularRepo.UpsertByID(ctx, item); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "db upsert failed"})
			}
		}
		return c.JSON(fiber.Map{"status": "saved"})
	}

	if form == nil {
		return c.Status(400).JSON(fiber.Map{"error": "no form data"})
	}

	suffixSet := map[int]struct{}{}
	for key := range form.Value {
		if strings.HasPrefix(key, "popular_title_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "popular_title_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
		if strings.HasPrefix(key, "popular_location_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "popular_location_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
		if strings.HasPrefix(key, "popular_image_path_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "popular_image_path_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
		if strings.HasPrefix(key, "popular_image_delete_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "popular_image_delete_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
	}

	if len(suffixSet) == 0 {
		suffixSet[1] = struct{}{}
	}

	suffixes := make([]int, 0, len(suffixSet))
	for idx := range suffixSet {
		suffixes = append(suffixes, idx)
	}
	sort.Ints(suffixes)

	getVal := func(key string) string {
		if vals, ok := form.Value[key]; ok && len(vals) > 0 {
			return vals[0]
		}
		return c.FormValue(key)
	}

	for _, idx := range suffixes {
		suffix := fmt.Sprintf("%d", idx)
		titleKey := "popular_title_" + suffix
		locationKey := "popular_location_" + suffix
		pathKey := "popular_image_path_" + suffix
		deleteKey := "popular_image_delete_" + suffix

		title := getVal(titleKey)
		location := getVal(locationKey)
		img := getVal(pathKey)
		deleteFlag := strings.ToLower(strings.TrimSpace(getVal(deleteKey)))
		shouldDelete := deleteFlag == "1" || deleteFlag == "true" || deleteFlag == "on"

		id := "popular_" + suffix

		var existingItem *adminmodel.PopularLocation
		if h.popularRepo != nil {
			existingItem, _ = h.popularRepo.FindByID(ctx, id)
		}

		if shouldDelete {
			if h.popularRepo != nil {
				if err := h.popularRepo.DeleteByID(ctx, id); err != nil {
					return c.Status(500).JSON(fiber.Map{"error": "db delete failed"})
				}
			}
			continue
		}

		if existingItem != nil {
			if title == "" && existingItem.Title != "" {
				title = existingItem.Title
			}
			if location == "" && existingItem.Location != "" {
				location = existingItem.Location
			}
			if img == "" && existingItem.ImagePath != "" {
				img = existingItem.ImagePath
			}
		}

		if title != "" || location != "" || img != "" {
			if h.popularRepo != nil {
				item := &adminmodel.PopularLocation{ID: id, Title: title, Location: location, ImagePath: img}
				if err := h.popularRepo.UpsertByID(ctx, item); err != nil {
					return c.Status(500).JSON(fiber.Map{"error": "db upsert failed"})
				}
			}
		}
	}

	return c.JSON(fiber.Map{"status": "saved"})
}

// GetHomePage returns a placeholder JSON of the homepage model (will be wired to repo later)
func (h *AdminHandler) GetHomePage(c *fiber.Ctx) error {
	// TODO: wire to HomePageRepository to fetch real data
	return c.JSON(fiber.Map{"message": "homepage get placeholder"})
}

// SaveHomePage accepts homepage payload and returns success placeholder
func (h *AdminHandler) SaveHomePage(c *fiber.Ctx) error {
	// Expect form values e.g. travel_image_title_1 and travel_image_path_1
	// Support multiple by iterating form values and upserting per slot
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// simple pattern: look for keys travel_image_title_<n>
	form, err := c.MultipartForm()
	if err != nil {
		// fallback to normal form values
		// iterate known single slot
		title := c.FormValue("travel_image_title_1")
		imgPath := c.FormValue("travel_image_path_1")
		if title == "" && imgPath == "" {
			return c.Status(400).JSON(fiber.Map{"error": "no data provided"})
		}

		te := map[string]string{"id": "image_1", "title": title, "image_path": imgPath}
		// upsert via repo
		if h.travelRepo != nil {
			item := &adminmodel.TravelExperience{ID: "image_1", Title: title, ImagePath: imgPath}
			if err := h.travelRepo.UpsertByID(ctx, item); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "db upsert failed"})
			}
		}
		return c.JSON(fiber.Map{"status": "saved", "data": te})
	}

	// handle multipart form: process any detected slot and preserve existing values when not changed
	if form == nil {
		return c.Status(400).JSON(fiber.Map{"error": "no form data"})
	}

	// collect suffixes from incoming form keys
	suffixSet := map[int]struct{}{}
	for key := range form.Value {
		if strings.HasPrefix(key, "travel_image_title_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "travel_image_title_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
		if strings.HasPrefix(key, "travel_image_path_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "travel_image_path_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
		if strings.HasPrefix(key, "travel_image_delete_") {
			if idx, err := strconv.Atoi(strings.TrimPrefix(key, "travel_image_delete_")); err == nil {
				suffixSet[idx] = struct{}{}
			}
		}
	}

	// if no suffixes found, default to slot 1
	if len(suffixSet) == 0 {
		suffixSet[1] = struct{}{}
	}

	// sort suffixes for deterministic processing
	suffixes := make([]int, 0, len(suffixSet))
	for idx := range suffixSet {
		suffixes = append(suffixes, idx)
	}
	sort.Ints(suffixes)

	// helper to fetch value from multipart form or plain form
	getVal := func(key string) string {
		if vals, ok := form.Value[key]; ok && len(vals) > 0 {
			return vals[0]
		}
		return c.FormValue(key)
	}

	for _, idx := range suffixes {
		suffix := fmt.Sprintf("%d", idx)
		titleKey := "travel_image_title_" + suffix
		pathKey := "travel_image_path_" + suffix
		deleteKey := "travel_image_delete_" + suffix

		title := getVal(titleKey)
		img := getVal(pathKey)
		deleteFlag := strings.ToLower(strings.TrimSpace(getVal(deleteKey)))
		shouldDelete := deleteFlag == "1" || deleteFlag == "true" || deleteFlag == "on"

		id := "image_" + suffix

		// Fetch existing record to preserve unchanged values
		var existingItem *adminmodel.TravelExperience
		if h.travelRepo != nil {
			existingItem, _ = h.travelRepo.FindByID(ctx, id)
		}

		if shouldDelete {
			if h.travelRepo != nil {
				if err := h.travelRepo.DeleteByID(ctx, id); err != nil {
					return c.Status(500).JSON(fiber.Map{"error": "db delete failed"})
				}
			}
			// also skip any further processing for this slot
			continue
		}

		// Preserve existing values if new values are empty
		if existingItem != nil {
			if title == "" && existingItem.Title != "" {
				title = existingItem.Title
			}
			if img == "" && existingItem.ImagePath != "" {
				img = existingItem.ImagePath
			}
		}

		// Upsert only if we have either a title or an image path for this slot
		if title != "" || img != "" {
			if h.travelRepo != nil {
				item := &adminmodel.TravelExperience{ID: id, Title: title, ImagePath: img}
				if err := h.travelRepo.UpsertByID(ctx, item); err != nil {
					return c.Status(500).JSON(fiber.Map{"error": "db upsert failed"})
				}
			}
		}
	}

	return c.JSON(fiber.Map{"status": "saved"})
}

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
	// Render fragment template with current DB values for all images (dynamic count)
	tmplPath := filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-travel-frag.html"))
	tmpl, err := template.New("cms-home-travel-frag.html").ParseFiles(tmplPath)
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	type travelSection struct {
		Index int
		Item  *adminmodel.TravelExperience
	}

	sections := map[int]*adminmodel.TravelExperience{}
	maxIndex := 0

	// Load existing records
	if h.travelRepo != nil {
		if items, err := h.travelRepo.ListAll(ctx); err == nil {
			for _, it := range items {
				idx := parseImageIndex(it.ID)
				if idx == 0 {
					continue
				}
				sections[idx] = it
				if idx > maxIndex {
					maxIndex = idx
				}
			}
		}
	}

	// Ensure at least three sections for initial UI
	if maxIndex < 3 {
		maxIndex = 3
	}

	for i := 1; i <= maxIndex; i++ {
		if _, ok := sections[i]; !ok {
			sections[i] = &adminmodel.TravelExperience{}
		}
	}

	ordered := make([]travelSection, 0, len(sections))
	for i := 1; i <= maxIndex; i++ {
		ordered = append(ordered, travelSection{Index: i, Item: sections[i]})
	}

	data := struct {
		Sections  []travelSection
		NextIndex int
	}{
		Sections:  ordered,
		NextIndex: maxIndex + 1,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// parseImageIndex extracts the numeric suffix from IDs like "image_4"
func parseImageIndex(id string) int {
	parts := strings.Split(id, "_")
	if len(parts) != 2 {
		return 0
	}
	idx, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return idx
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
	// Render fragment template with current DB values for banners
	tmplPath := filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-banner-frag.html"))
	tmpl, err := template.New("cms-home-banner-frag.html").ParseFiles(tmplPath)
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	banner1 := &adminmodel.Banner{}
	banner2 := &adminmodel.Banner{}

	if h.bannerRepo != nil {
		if items, err := h.bannerRepo.ListAll(ctx); err == nil {
			for _, it := range items {
				if it.ID == "banner_1" {
					banner1 = it
				} else if it.ID == "banner_2" {
					banner2 = it
				}
			}
		}
	}

	data := struct {
		Banner1 *adminmodel.Banner
		Banner2 *adminmodel.Banner
	}{
		Banner1: banner1,
		Banner2: banner2,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
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
	// Render fragment template with current DB values for popular locations
	tmplPath := filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-popular-frag.html"))
	tmpl, err := template.New("cms-home-popular-frag.html").ParseFiles(tmplPath)
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	type popularSection struct {
		Index int
		Item  *adminmodel.PopularLocation
	}

	sections := map[int]*adminmodel.PopularLocation{}
	maxIndex := 0

	if h.popularRepo != nil {
		if items, err := h.popularRepo.ListAll(ctx); err == nil {
			for _, it := range items {
				idx := parseImageIndex(it.ID)
				if idx == 0 {
					continue
				}
				sections[idx] = it
				if idx > maxIndex {
					maxIndex = idx
				}
			}
		}
	}

	if maxIndex < 3 {
		maxIndex = 3
	}

	for i := 1; i <= maxIndex; i++ {
		if _, ok := sections[i]; !ok {
			sections[i] = &adminmodel.PopularLocation{}
		}
	}

	ordered := make([]popularSection, 0, len(sections))
	for i := 1; i <= maxIndex; i++ {
		ordered = append(ordered, popularSection{Index: i, Item: sections[i]})
	}

	data := struct {
		Sections  []popularSection
		NextIndex int
	}{
		Sections:  ordered,
		NextIndex: maxIndex + 1,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
