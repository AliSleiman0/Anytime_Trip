package services

import (
	"context"
	"fmt"

	"Anytime_Travel/backend/core/utils"
	appmodels "Anytime_Travel/backend/internal/models/app"
	apprepo "Anytime_Travel/backend/internal/repository/app"
)

// BookingService handles booking operations with notification support
type BookingService struct {
	hotelBookingRepo    *apprepo.HotelBookingRepository
	flightBookingRepo   *apprepo.FlightBookingRepository
	carBookingRepo      *apprepo.CarBookingRepository
	transferBookingRepo *apprepo.TransferBookingRepository
	notificationHelper  *utils.NotificationHelper
}

// NewBookingService creates a new booking service
func NewBookingService(
	hotelBookingRepo *apprepo.HotelBookingRepository,
	flightBookingRepo *apprepo.FlightBookingRepository,
	carBookingRepo *apprepo.CarBookingRepository,
	transferBookingRepo *apprepo.TransferBookingRepository,
	notificationHelper *utils.NotificationHelper,
) *BookingService {
	return &BookingService{
		hotelBookingRepo:    hotelBookingRepo,
		flightBookingRepo:   flightBookingRepo,
		carBookingRepo:      carBookingRepo,
		transferBookingRepo: transferBookingRepo,
		notificationHelper:  notificationHelper,
	}
}

// CreateHotelBooking creates a hotel booking and sends notification email
func (s *BookingService) CreateHotelBooking(ctx context.Context, booking *appmodels.HotelBooking) error {
	// Create the booking in database
	if err := s.hotelBookingRepo.Create(ctx, booking); err != nil {
		return fmt.Errorf("failed to create hotel booking: %w", err)
	}

	// Send notification email (non-blocking - log errors but don't fail the booking)
	if err := s.notificationHelper.SendHotelBookingEmail(ctx, booking); err != nil {
		fmt.Printf("[BOOKING] Warning: Failed to send booking notification email: %v\n", err)
		// Don't return error - booking was successful, email is secondary
	}

	return nil
}

// CreateFlightBooking creates a flight booking and sends notification email
func (s *BookingService) CreateFlightBooking(ctx context.Context, booking *appmodels.FlightBooking) error {
	// Create the booking in database
	if err := s.flightBookingRepo.Create(ctx, booking); err != nil {
		return fmt.Errorf("failed to create flight booking: %w", err)
	}

	// Send notification email (non-blocking - log errors but don't fail the booking)
	if err := s.notificationHelper.SendFlightBookingEmail(ctx, booking); err != nil {
		fmt.Printf("[BOOKING] Warning: Failed to send booking notification email: %v\n", err)
		// Don't return error - booking was successful, email is secondary
	}

	return nil
}

// CreateCarBooking creates a car booking and sends notification email
func (s *BookingService) CreateCarBooking(ctx context.Context, booking *appmodels.CarBooking) error {
	// Create the booking in database
	if err := s.carBookingRepo.Create(ctx, booking); err != nil {
		return fmt.Errorf("failed to create car booking: %w", err)
	}

	// Send notification email (non-blocking - log errors but don't fail the booking)
	if err := s.notificationHelper.SendCarBookingEmail(ctx, booking); err != nil {
		fmt.Printf("[BOOKING] Warning: Failed to send booking notification email: %v\n", err)
		// Don't return error - booking was successful, email is secondary
	}

	return nil
}

// CreateTransferBooking creates a transfer booking and sends notification email
func (s *BookingService) CreateTransferBooking(ctx context.Context, booking *appmodels.TransferBooking) error {
	// Create the booking in database
	if err := s.transferBookingRepo.Create(ctx, booking); err != nil {
		return fmt.Errorf("failed to create transfer booking: %w", err)
	}

	// Send notification email (non-blocking - log errors but don't fail the booking)
	if err := s.notificationHelper.SendTransferBookingEmail(ctx, booking); err != nil {
		fmt.Printf("[BOOKING] Warning: Failed to send booking notification email: %v\n", err)
		// Don't return error - booking was successful, email is secondary
	}

	return nil
}
