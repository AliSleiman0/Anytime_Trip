package admin

import (
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"sort"
	"time"

	app "Anytime_Travel/backend/internal/models/app"

	"github.com/gofiber/fiber/v2"
)

// GetPayments renders the payments page (full page)
func (h *AdminHandler) GetPayments(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "payments.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{"Title": "Payments and Transactions"}
	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetPaymentsFragment serves the payments fragment with dynamic data
func (h *AdminHandler) GetPaymentsFragment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payments, err := h.paymentRepo.GetAllPayments(ctx)
	if err != nil {
		return c.Status(500).SendString("Error fetching payments: " + err.Error())
	}

	// keep a copy of the full, unfiltered payments for building filter lists
	allPayments := make([]app.Payment, len(payments))
	copy(allPayments, payments)

	// read filters from query params
	dateFilter := c.Query("date")
	methodFilter := c.Query("method")
	typeFilter := c.Query("type")
	statusFilter := c.Query("status")

	// apply simple server-side filtering
	filtered := make([]app.Payment, 0, len(payments))
	for _, p := range payments {
		// date match (if provided)
		if dateFilter != "" {
			if pd := p.BookingDate; !pd.IsZero() {
				// try parse incoming date (YYYY-MM-DD)
				if t, err := time.Parse("2006-01-02", dateFilter); err == nil {
					if pd.Year() != t.Year() || pd.Month() != t.Month() || pd.Day() != t.Day() {
						continue
					}
				}
			}
		}
		if methodFilter != "" && methodFilter != "all" {
			if p.PaymentMethod != methodFilter {
				continue
			}
		}
		if typeFilter != "" && typeFilter != "all" {
			if p.Type != typeFilter {
				continue
			}
		}
		if statusFilter != "" && statusFilter != "all" {
			if p.PaymentStatus != statusFilter {
				continue
			}
		}
		filtered = append(filtered, p)
	}
	// use filtered list for rendering/stats
	payments = filtered

	// compute simple stats
	var overall float64
	var refunds int
	var cancellations int
	for _, p := range payments {
		overall += p.Amount
		switch p.PaymentStatus {
		case "Refunded", "Refund":
			refunds++
		case "Cancelled", "Canceled":
			cancellations++
		}
	}

	// compute percentage changes by comparing current period to previous 30 days
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	sixtyDaysAgo := now.AddDate(0, 0, -60)

	var prevOverall float64
	var prevPaymentsCount int
	var prevRefunds int
	var prevCancellations int

	for _, p := range allPayments {
		if p.BookingDate.After(sixtyDaysAgo) && p.BookingDate.Before(thirtyDaysAgo) {
			prevOverall += p.Amount
			prevPaymentsCount++
			switch p.PaymentStatus {
			case "Refunded", "Refund":
				prevRefunds++
			case "Cancelled", "Canceled":
				prevCancellations++
			}
		}
	}

	// calculate percentage changes
	calcPercentage := func(current, previous float64) (float64, bool) {
		if previous == 0 {
			if current > 0 {
				return 100.0, true
			}
			return 0, true
		}
		change := ((current - previous) / previous) * 100
		return change, change >= 0
	}

	revenuePercent, revenuePositive := calcPercentage(overall, prevOverall)
	paymentsPercent, paymentsPositive := calcPercentage(float64(len(payments)), float64(prevPaymentsCount))
	refundsPercent, refundsPositive := calcPercentage(float64(refunds), float64(prevRefunds))
	cancellationsPercent, cancellationsPositive := calcPercentage(float64(cancellations), float64(prevCancellations))

	// build dynamic filter lists from the full (unfiltered) payments
	methodSet := make(map[string]struct{})
	typeSet := make(map[string]struct{})
	statusSet := make(map[string]struct{})
	for _, p := range allPayments {
		if p.PaymentMethod != "" {
			methodSet[p.PaymentMethod] = struct{}{}
		}
		if p.Type != "" {
			typeSet[p.Type] = struct{}{}
		}
		if p.PaymentStatus != "" {
			statusSet[p.PaymentStatus] = struct{}{}
		}
	}

	paymentMethods := make([]string, 0, len(methodSet))
	bookingTypes := make([]string, 0, len(typeSet))
	paymentStatuses := make([]string, 0, len(statusSet))
	for k := range methodSet {
		paymentMethods = append(paymentMethods, k)
	}
	for k := range typeSet {
		bookingTypes = append(bookingTypes, k)
	}
	for k := range statusSet {
		paymentStatuses = append(paymentStatuses, k)
	}
	sort.Strings(paymentMethods)
	sort.Strings(bookingTypes)
	sort.Strings(paymentStatuses)

	data := fiber.Map{
		"OverallRevenue":         formatCurrency("$", overall),
		"OverallRevenuePercent":  fmt.Sprintf("%.2f", revenuePercent),
		"OverallRevenuePositive": revenuePositive,
		"PaymentsCount":          len(payments),
		"PaymentsPercent":        fmt.Sprintf("%.2f", paymentsPercent),
		"PaymentsPositive":       paymentsPositive,
		"RefundsCount":           refunds,
		"RefundsPercent":         fmt.Sprintf("%.2f", refundsPercent),
		"RefundsPositive":        refundsPositive,
		"CancellationsCount":     cancellations,
		"CancellationsPercent":   fmt.Sprintf("%.2f", cancellationsPercent),
		"CancellationsPositive":  cancellationsPositive,
		"Payments":               payments,
		"PaymentMethods":         paymentMethods,
		"BookingTypes":           bookingTypes,
		"PaymentStatuses":        paymentStatuses,
		"Filters":                map[string]string{"Date": dateFilter, "Method": methodFilter, "Type": typeFilter, "Status": statusFilter},
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "payments-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template: " + err.Error())
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetAllPayments returns payments as JSON
func (h *AdminHandler) GetAllPayments(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payments, err := h.paymentRepo.GetAllPayments(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch payments: %v", err)})
	}
	return c.JSON(payments)
}
