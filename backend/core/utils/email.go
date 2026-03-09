package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

// EmailService handles sending emails
type EmailService struct {
	smtpHost       string
	smtpPort       string
	smtpUsername   string
	senderEmail    string
	senderPassword string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	es := &EmailService{
		smtpHost:       getEnvOrDefault("SMTP_HOST", "mail-eu.smtp2go.com"),
		smtpPort:       getEnvOrDefault("SMTP_PORT", "2525"),
		smtpUsername:   getEnvOrDefault("SMTP_USERNAME", "anytimetravel.app"),
		senderEmail:    getEnvOrDefault("SENDER_EMAIL", "info@anytimetravel.app"),
		senderPassword: getEnvOrDefault("SENDER_PASSWORD", ""),
	}
	fmt.Printf("[EMAIL] SMTP Config - Host: %s, Port: %s, Username: %s, From: %s, HasPassword: %v\n",
		es.smtpHost, es.smtpPort, es.smtpUsername, es.senderEmail, es.senderPassword != "")
	return es
}

// SendPasswordResetEmail sends a password reset email
func (es *EmailService) SendPasswordResetEmail(recipientEmail, resetLink string) error {
	// If no SMTP credentials configured, log and skip (useful for development)
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Password reset link for %s: %s\n", recipientEmail, resetLink)
		return nil
	}

	subject := "Password Reset Request - Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; }
    .header { background-color: #336891; color: white; padding: 20px; text-align: center; }
    .content { padding: 20px; border: 1px solid #ddd; }
    .button { background-color: #336891; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; }
    .footer { text-align: center; color: #666; margin-top: 20px; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Password Reset Request</h1>
    </div>
    <div class="content">
      <p>Hello,</p>
      <p>We received a request to reset your password. Click the button below to create a new password:</p>
      <p style="text-align: center; margin: 30px 0;">
        <a href="%s" class="button">Reset Password</a>
      </p>
      <p>This link will expire in 1 hour.</p>
      <p>If you didn't request a password reset, you can safely ignore this email.</p>
      <p>Best regards,<br>Anytime Travel Team</p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
`, resetLink)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendPasswordResetCodeEmail sends a verification code for password reset
func (es *EmailService) SendPasswordResetCodeEmail(recipientEmail, code string) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Password reset code for %s: %s\n", recipientEmail, code)
		return nil
	}

	subject := "Your password reset code - Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 20px auto; background-color: white; }
    .header { background-color: #336891; color: white; padding: 30px 20px; text-align: center; }
    .header h1 { margin: 0; font-size: 28px; }
    .content { padding: 30px 20px; }
    .otp-box { background-color: #f9f9f9; border: 2px dashed #336891; padding: 20px; margin: 20px 0; text-align: center; }
    .otp-code { font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #336891; margin: 10px 0; }
    .footer { background-color: #f4f4f4; text-align: center; padding: 20px; color: #666; font-size: 12px; }
    .warning { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🔒 Password Reset</h1>
    </div>
    <div class="content">
      <p>Hello,</p>
      <p>We received a request to reset your password. Use the verification code below to complete the process:</p>
      
      <div class="otp-box">
        <p style="margin: 0; color: #666;">Your verification code:</p>
        <div class="otp-code">%s</div>
        <p style="margin: 0; color: #999; font-size: 12px;">This code expires in 10 minutes</p>
      </div>

      <div class="warning">
        <strong>⚠️ Security Notice:</strong><br>
        If you didn't request this password reset, please ignore this email or contact our support team immediately.
      </div>

      <p>Best regards,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, code)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendOTPEmail sends an OTP code via email for general verification purposes
func (es *EmailService) SendOTPEmail(recipientEmail, code, purpose string) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] OTP for %s (%s): %s\n", recipientEmail, purpose, code)
		return nil
	}

	subject := "Your verification code - Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 20px auto; background-color: white; }
    .header { background-color: #336891; color: white; padding: 30px 20px; text-align: center; }
    .header h1 { margin: 0; font-size: 28px; }
    .content { padding: 30px 20px; }
    .otp-box { background-color: #f9f9f9; border: 2px dashed #336891; padding: 20px; margin: 20px 0; text-align: center; }
    .otp-code { font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #336891; margin: 10px 0; }
    .footer { background-color: #f4f4f4; text-align: center; padding: 20px; color: #666; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Verification Code</h1>
    </div>
    <div class="content">
      <p>Hello,</p>
      <p>Use the verification code below to %s:</p>
      
      <div class="otp-box">
        <p style="margin: 0; color: #666;">Your verification code:</p>
        <div class="otp-code">%s</div>
        <p style="margin: 0; color: #999; font-size: 12px;">This code expires in 10 minutes</p>
      </div>

      <p>If you didn't request this code, please ignore this email.</p>

      <p>Best regards,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, purpose, code)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendNewBookingNotification sends a new booking notification email to a customer
func (es *EmailService) SendNewBookingNotification(recipientEmail, customerName, bookingID, bookingType string, amount float64, bookingDate time.Time) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] New booking notification for %s: %s (%s)\n", recipientEmail, bookingID, bookingType)
		return nil
	}

	subject := fmt.Sprintf("Booking Confirmation - %s", bookingID)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 20px auto; background-color: white; }
    .header { background-color: #336891; color: white; padding: 30px 20px; text-align: center; }
    .header h1 { margin: 0; font-size: 28px; }
    .content { padding: 30px 20px; }
    .booking-details { background-color: #f9f9f9; border-left: 4px solid #336891; padding: 20px; margin: 20px 0; }
    .booking-details h2 { margin-top: 0; color: #336891; }
    .detail-row { margin: 10px 0; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .amount { font-size: 24px; font-weight: bold; color: #336891; margin: 20px 0; }
    .footer { background-color: #f4f4f4; text-align: center; padding: 20px; color: #666; font-size: 12px; }
    .button { background-color: #336891; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>✓ Booking Confirmed</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>
      <p>Thank you for choosing Anytime Travel! We're pleased to confirm your booking.</p>
      
      <div class="booking-details">
        <h2>Booking Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Booking Type:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Booking Date:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="amount">Total Amount: $%.2f</div>
      </div>

      <p>Your booking has been successfully processed. You will receive additional details and confirmation documents shortly.</p>
      
      <p>If you have any questions or need to make changes to your booking, please don't hesitate to contact our support team.</p>
      
      <p style="text-align: center;">
        <a href="https://anytimetravel.app/bookings" class="button">View My Bookings</a>
      </p>

      <p>Best regards,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, bookingType, bookingDate.Format("Monday, January 2, 2006"), amount)

	return es.sendEmail(recipientEmail, subject, body)
}

// FlightSegment represents a flight segment in the journey
type FlightSegment struct {
	Airline          string
	FlightNumber     string
	Class            string
	DepartureCity    string
	DepartureCode    string
	DepartureAirport string
	DepartureDate    string
	DepartureTime    string
	ArrivalCity      string
	ArrivalCode      string
	ArrivalAirport   string
	ArrivalDate      string
	ArrivalTime      string
	Duration         string
}

// SendFlightBookingEmail sends a flight booking confirmation email
func (es *EmailService) SendFlightBookingEmail(
	recipientEmail, customerName, confirmationNumber string,
	bookingDate time.Time,
	travelers []string,
	outboundSegments []FlightSegment,
	returnSegments []FlightSegment,
	basePrice, taxes, serviceFee, totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Flight booking confirmation for %s: %s\n", recipientEmail, confirmationNumber)
		return nil
	}

	subject := fmt.Sprintf("Flight Booking Confirmed - %s", confirmationNumber)

	// Build travelers HTML
	travelersHTML := ""
	for _, traveler := range travelers {
		travelersHTML += fmt.Sprintf(`
          <p style="margin: 0; padding: 12px 0 0 14px; font-family: 'Roboto', Arial, sans-serif; font-size: 14px; line-height: 21px; color: #1D1D1D;">%s</p>`, traveler)
	}

	// Build outbound flights HTML
	outboundHTML := ""
	for _, segment := range outboundSegments {
		outboundHTML += fmt.Sprintf(`
        <div style="margin-bottom: 16px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
            <tr>
              <td style="vertical-align: top;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">%s</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              </td>
              <td width="30" style="width: 30px;"></td>
              <td style="vertical-align: top; text-align: right; white-space: nowrap;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Class</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #336891;">%s</p>
              </td>
            </tr>
          </table>
          <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
            <div style="margin-bottom: 20px;">
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
            <div style="text-align: center; margin: 12px 0; color: rgba(0, 0, 0, 0.7);">
              <span style="font-size: 12px;">✈️ %s</span>
            </div>
            <div>
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
          </div>
        </div>`,
			segment.Airline,
			segment.FlightNumber,
			segment.Class,
			segment.DepartureCode,
			segment.DepartureCity,
			segment.DepartureAirport,
			segment.DepartureDate,
			segment.DepartureTime,
			segment.Duration,
			segment.ArrivalCode,
			segment.ArrivalCity,
			segment.ArrivalAirport,
			segment.ArrivalDate,
			segment.ArrivalTime,
		)
	}

	// Build return flights HTML
	returnHTML := ""
	if len(returnSegments) > 0 {
		for _, segment := range returnSegments {
			returnHTML += fmt.Sprintf(`
        <div style="margin-bottom: 16px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
            <tr>
              <td style="vertical-align: top;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">%s</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              </td>
              <td width="30" style="width: 30px;"></td>
              <td style="vertical-align: top; text-align: right; white-space: nowrap;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Class</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #336891;">%s</p>
              </td>
            </tr>
          </table>
          <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
            <div style="margin-bottom: 20px;">
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
            <div style="text-align: center; margin: 12px 0; color: rgba(0, 0, 0, 0.7);">
              <span style="font-size: 12px;">✈️ %s</span>
            </div>
            <div>
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
          </div>
        </div>`,
				segment.Airline,
				segment.FlightNumber,
				segment.Class,
				segment.DepartureCode,
				segment.DepartureCity,
				segment.DepartureAirport,
				segment.DepartureDate,
				segment.DepartureTime,
				segment.Duration,
				segment.ArrivalCode,
				segment.ArrivalCity,
				segment.ArrivalAirport,
				segment.ArrivalDate,
				segment.ArrivalTime,
			)
		}

		returnHTML = fmt.Sprintf(`
      <div style="padding: 0 20px 14px; border-bottom: 1px solid rgba(51, 104, 145, 0.31);">
        <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Return Journey</h2>
        %s
      </div>`, returnHTML)
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; font-family: 'Roboto', Arial, sans-serif; background-color: #f4f4f4; }
    .email-container { max-width: 430px; margin: 0 auto; background: #FFFFFF; }
    @media only screen and (max-width: 430px) {
      .email-container { width: 100%% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div style="background: #336891; padding: 24px 24px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 16px;">
        <tr>
          <td style="vertical-align: middle;">
            <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.9);">ANYTIME TRAVEL & TOURISM</p>
            <h1 style="margin: 4px 0 0 0; font-size: 24px; font-weight: 700; color: #FFFFFF;">Booking Confirmed</h1>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td width="48" style="width: 48px; vertical-align: middle;">
            <div style="width: 48px; height: 48px; border: 4px solid #FFFFFF; border-radius: 50%%; text-align: center; line-height: 48px;">
              <span style="color: #FFFFFF; font-size: 28px; font-weight: bold;">✓</span>
            </div>
          </td>
        </tr>
      </table>
      
      <div style="border-top: 1px solid rgba(255, 255, 255, 0.3); padding-top: 13px; padding-bottom: 24px;">
        <div style="margin-bottom: 12px;">
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Confirmation Number</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
        <div>
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Booking Date</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
      </div>
    </div>

    <!-- Success Banner -->
    <div style="background: #E8F5E9; border-left: 2px solid #4CAF50; padding: 12px 19px;">
      <p style="margin: 0; font-size: 14px; font-weight: 700; color: #2E7D32;">✓ Booking Confirmed</p>
      <p style="margin: 4px 0 0 0; font-size: 12px; line-height: 18px; color: #1D1D1D;">Your flight has been confirmed. A copy has been sent to your email address.</p>
    </div>

    <!-- Travelers -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Travelers</h2>
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 0;">
        %s
      </div>
    </div>

    <!-- Outbound Journey -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Outbound Journey</h2>
      %s
    </div>

    <!-- Return Journey -->
    %s

    <!-- Price Summary -->
    <div style="padding: 0 20px; margin: 13px 0;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Price Summary</h2>
      <div style="margin-bottom: 13px;">
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Flight Ticket</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Taxes</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Service Fee</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="text-align: right; margin-top: 13px;">
          <span style="font-size: 14px; color: #D24124; font-weight: 700;">Total: $%.2f</span>
        </div>
      </div>
    </div>

    <!-- Important Information -->
    <div style="background: #E3F2FD; border-left: 2px solid #336891; padding: 12px 14px; margin-bottom: 17px;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Important Information</p>
      <div style="font-size: 12px; line-height: 18px; color: #1D1D1D;">
        <p style="margin: 0 0 6px 0;">• Arrive 3 hours before departure</p>
        <p style="margin: 0 0 6px 0;">• Check-in opens 24 hours before flight</p>
        <p style="margin: 0 0 6px 0;">• Valid passport required</p>
        <p style="margin: 0;">• Check baggage allowance for your class</p>
      </div>
    </div>

    <!-- Footer -->
    <div style="border-top: 1px solid rgba(51, 104, 145, 0.31); padding: 17px 0; text-align: center;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Thank you for booking with us!</p>
      <p style="margin: 0 0 8px 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">For support: info@anytimetravel.app</p>
      <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone: +961 1 234 567</p>
    </div>
  </div>
</body>
</html>`,
		confirmationNumber,
		bookingDate.Format("January 2, 2006"),
		travelersHTML,
		outboundHTML,
		returnHTML,
		basePrice,
		taxes,
		serviceFee,
		totalPrice,
	)

	return es.sendEmail(recipientEmail, subject, body)
}

// sendEmail sends an email using SMTP
func (es *EmailService) sendEmail(to, subject, body string) error {
	fmt.Printf("[EMAIL] Attempting to send email to %s using SMTP %s:%s with username %s\n",
		to, es.smtpHost, es.smtpPort, es.smtpUsername)

	auth := smtp.PlainAuth("", es.smtpUsername, es.senderPassword, es.smtpHost)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		es.senderEmail, to, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", es.smtpHost, es.smtpPort)
	err := smtp.SendMail(addr, auth, es.senderEmail, []string{to}, []byte(msg))

	if err != nil {
		fmt.Printf("[EMAIL] Failed to send email: %v\n", err)
		return err
	}

	fmt.Printf("[EMAIL] Successfully sent email to %s\n", to)
	return nil
}

// SendOneWayFlightBookingEmail sends a one-way flight booking confirmation email
func (es *EmailService) SendOneWayFlightBookingEmail(
	recipientEmail, customerName, confirmationNumber string,
	bookingDate time.Time,
	travelers []string,
	outboundSegments []FlightSegment,
	basePrice, taxes, serviceFee, totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] One-way flight booking confirmation for %s: %s\n", recipientEmail, confirmationNumber)
		return nil
	}

	subject := fmt.Sprintf("Flight Booking Confirmed - %s", confirmationNumber)

	// Build travelers HTML
	travelersHTML := ""
	for _, traveler := range travelers {
		travelersHTML += fmt.Sprintf(`
          <p style="margin: 0; padding: 12px 0 0 14px; font-family: 'Roboto', Arial, sans-serif; font-size: 14px; line-height: 21px; color: #1D1D1D;">%s</p>`, traveler)
	}

	// Build outbound flights HTML
	outboundHTML := ""
	for _, segment := range outboundSegments {
		outboundHTML += fmt.Sprintf(`
        <div style="margin-bottom: 16px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
            <tr>
              <td style="vertical-align: top;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">%s</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              </td>
              <td width="30" style="width: 30px;"></td>
              <td style="vertical-align: top; text-align: right; white-space: nowrap;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Class</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #336891;">%s</p>
              </td>
            </tr>
          </table>
          <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
            <div style="margin-bottom: 20px;">
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
            <div style="text-align: center; margin: 12px 0; color: rgba(0, 0, 0, 0.7);">
              <span style="font-size: 12px;">✈️ %s</span>
            </div>
            <div>
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
          </div>
        </div>`,
			segment.Airline,
			segment.FlightNumber,
			segment.Class,
			segment.DepartureCode,
			segment.DepartureCity,
			segment.DepartureAirport,
			segment.DepartureDate,
			segment.DepartureTime,
			segment.Duration,
			segment.ArrivalCode,
			segment.ArrivalCity,
			segment.ArrivalAirport,
			segment.ArrivalDate,
			segment.ArrivalTime,
		)
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; font-family: 'Roboto', Arial, sans-serif; background-color: #f4f4f4; }
    .email-container { max-width: 430px; margin: 0 auto; background: #FFFFFF; }
    @media only screen and (max-width: 430px) {
      .email-container { width: 100%% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div style="background: #336891; padding: 24px 24px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 16px;">
        <tr>
          <td style="vertical-align: middle;">
            <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.9);">ANYTIME TRAVEL & TOURISM</p>
            <h1 style="margin: 4px 0 0 0; font-size: 24px; font-weight: 700; color: #FFFFFF;">Booking Confirmed</h1>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td width="48" style="width: 48px; vertical-align: middle;">
            <div style="width: 48px; height: 48px; border: 4px solid #FFFFFF; border-radius: 50%%; text-align: center; line-height: 48px;">
              <span style="color: #FFFFFF; font-size: 28px; font-weight: bold;">✓</span>
            </div>
          </td>
        </tr>
      </table>
      
      <div style="border-top: 1px solid rgba(255, 255, 255, 0.3); padding-top: 13px; padding-bottom: 24px;">
        <div style="margin-bottom: 12px;">
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Confirmation Number</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
        <div>
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Booking Date</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
      </div>
    </div>

    <!-- Success Banner -->
    <div style="background: #E8F5E9; border-left: 2px solid #4CAF50; padding: 12px 19px;">
      <p style="margin: 0; font-size: 14px; font-weight: 700; color: #2E7D32;">✓ Booking Confirmed</p>
      <p style="margin: 4px 0 0 0; font-size: 12px; line-height: 18px; color: #1D1D1D;">Your one-way flight has been confirmed. A copy has been sent to your email address.</p>
    </div>

    <!-- Travelers -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Travelers</h2>
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 0;">
        %s
      </div>
    </div>

    <!-- Flight Details -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Flight Details</h2>
      %s
    </div>

    <!-- Price Summary -->
    <div style="padding: 0 20px; margin: 13px 0;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Price Summary</h2>
      <div style="margin-bottom: 13px;">
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Flight Ticket</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Taxes</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Service Fee</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="text-align: right; margin-top: 13px;">
          <span style="font-size: 14px; color: #D24124; font-weight: 700;">Total: $%.2f</span>
        </div>
      </div>
    </div>

    <!-- Important Information -->
    <div style="background: #E3F2FD; border-left: 2px solid #336891; padding: 12px 14px; margin-bottom: 17px;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Important Information</p>
      <div style="font-size: 12px; line-height: 18px; color: #1D1D1D;">
        <p style="margin: 0 0 6px 0;">• Arrive 3 hours before departure</p>
        <p style="margin: 0 0 6px 0;">• Check-in opens 24 hours before flight</p>
        <p style="margin: 0 0 6px 0;">• Valid passport required</p>
        <p style="margin: 0;">• Check baggage allowance for your class</p>
      </div>
    </div>

    <!-- Footer -->
    <div style="border-top: 1px solid rgba(51, 104, 145, 0.31); padding: 17px 0; text-align: center;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Thank you for booking with us!</p>
      <p style="margin: 0 0 8px 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">For support: info@anytimetravel.app</p>
      <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone: +961 1 234 567</p>
    </div>
  </div>
</body>
</html>`,
		confirmationNumber,
		bookingDate.Format("January 2, 2006"),
		travelersHTML,
		outboundHTML,
		basePrice,
		taxes,
		serviceFee,
		totalPrice,
	)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendRoundTripFlightBookingEmail sends a round-trip flight booking confirmation email
func (es *EmailService) SendRoundTripFlightBookingEmail(
	recipientEmail, customerName, confirmationNumber string,
	bookingDate time.Time,
	travelers []string,
	outboundSegments []FlightSegment,
	returnSegments []FlightSegment,
	basePrice, taxes, serviceFee, totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Round-trip flight booking confirmation for %s: %s\n", recipientEmail, confirmationNumber)
		return nil
	}

	subject := fmt.Sprintf("Flight Booking Confirmed - %s", confirmationNumber)

	// Build travelers HTML
	travelersHTML := ""
	for _, traveler := range travelers {
		travelersHTML += fmt.Sprintf(`
          <p style="margin: 0; padding: 12px 0 0 14px; font-family: 'Roboto', Arial, sans-serif; font-size: 14px; line-height: 21px; color: #1D1D1D;">%s</p>`, traveler)
	}

	// Build outbound flights HTML
	outboundHTML := ""
	for _, segment := range outboundSegments {
		outboundHTML += fmt.Sprintf(`
        <div style="margin-bottom: 16px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
            <tr>
              <td style="vertical-align: top;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">%s</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              </td>
              <td width="30" style="width: 30px;"></td>
              <td style="vertical-align: top; text-align: right; white-space: nowrap;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Class</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #336891;">%s</p>
              </td>
            </tr>
          </table>
          <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
            <div style="margin-bottom: 20px;">
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
            <div style="text-align: center; margin: 12px 0; color: rgba(0, 0, 0, 0.7);">
              <span style="font-size: 12px;">✈️ %s</span>
            </div>
            <div>
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
          </div>
        </div>`,
			segment.Airline,
			segment.FlightNumber,
			segment.Class,
			segment.DepartureCode,
			segment.DepartureCity,
			segment.DepartureAirport,
			segment.DepartureDate,
			segment.DepartureTime,
			segment.Duration,
			segment.ArrivalCode,
			segment.ArrivalCity,
			segment.ArrivalAirport,
			segment.ArrivalDate,
			segment.ArrivalTime,
		)
	}

	// Build return flights HTML
	returnHTML := ""
	for _, segment := range returnSegments {
		returnHTML += fmt.Sprintf(`
        <div style="margin-bottom: 16px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
            <tr>
              <td style="vertical-align: top;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">%s</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              </td>
              <td width="30" style="width: 30px;"></td>
              <td style="vertical-align: top; text-align: right; white-space: nowrap;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Class</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #336891;">%s</p>
              </td>
            </tr>
          </table>
          <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
            <div style="margin-bottom: 20px;">
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
            <div style="text-align: center; margin: 12px 0; color: rgba(0, 0, 0, 0.7);">
              <span style="font-size: 12px;">✈️ %s</span>
            </div>
            <div>
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
          </div>
        </div>`,
			segment.Airline,
			segment.FlightNumber,
			segment.Class,
			segment.DepartureCode,
			segment.DepartureCity,
			segment.DepartureAirport,
			segment.DepartureDate,
			segment.DepartureTime,
			segment.Duration,
			segment.ArrivalCode,
			segment.ArrivalCity,
			segment.ArrivalAirport,
			segment.ArrivalDate,
			segment.ArrivalTime,
		)
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; font-family: 'Roboto', Arial, sans-serif; background-color: #f4f4f4; }
    .email-container { max-width: 430px; margin: 0 auto; background: #FFFFFF; }
    @media only screen and (max-width: 430px) {
      .email-container { width: 100%% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div style="background: #336891; padding: 24px 24px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 16px;">
        <tr>
          <td style="vertical-align: middle;">
            <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.9);">ANYTIME TRAVEL & TOURISM</p>
            <h1 style="margin: 4px 0 0 0; font-size: 24px; font-weight: 700; color: #FFFFFF;">Booking Confirmed</h1>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td width="48" style="width: 48px; vertical-align: middle;">
            <div style="width: 48px; height: 48px; border: 4px solid #FFFFFF; border-radius: 50%%; text-align: center; line-height: 48px;">
              <span style="color: #FFFFFF; font-size: 28px; font-weight: bold;">✓</span>
            </div>
          </td>
        </tr>
      </table>
      
      <div style="border-top: 1px solid rgba(255, 255, 255, 0.3); padding-top: 13px; padding-bottom: 24px;">
        <div style="margin-bottom: 12px;">
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Confirmation Number</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
        <div>
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Booking Date</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
      </div>
    </div>

    <!-- Success Banner -->
    <div style="background: #E8F5E9; border-left: 2px solid #4CAF50; padding: 12px 19px;">
      <p style="margin: 0; font-size: 14px; font-weight: 700; color: #2E7D32;">✓ Booking Confirmed</p>
      <p style="margin: 4px 0 0 0; font-size: 12px; line-height: 18px; color: #1D1D1D;">Your round-trip flight has been confirmed. A copy has been sent to your email address.</p>
    </div>

    <!-- Travelers -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Travelers</h2>
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 0;">
        %s
      </div>
    </div>

    <!-- Outbound Journey -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Outbound Journey</h2>
      %s
    </div>

    <!-- Return Journey -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Return Journey</h2>
      %s
    </div>

    <!-- Price Summary -->
    <div style="padding: 0 20px; margin: 13px 0;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Price Summary</h2>
      <div style="margin-bottom: 13px;">
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Flight Ticket</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Taxes</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Service Fee</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="text-align: right; margin-top: 13px;">
          <span style="font-size: 14px; color: #D24124; font-weight: 700;">Total: $%.2f</span>
        </div>
      </div>
    </div>

    <!-- Important Information -->
    <div style="background: #E3F2FD; border-left: 2px solid #336891; padding: 12px 14px; margin-bottom: 17px;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Important Information</p>
      <div style="font-size: 12px; line-height: 18px; color: #1D1D1D;">
        <p style="margin: 0 0 6px 0;">• Arrive 3 hours before departure</p>
        <p style="margin: 0 0 6px 0;">• Check-in opens 24 hours before flight</p>
        <p style="margin: 0 0 6px 0;">• Valid passport required</p>
        <p style="margin: 0;">• Check baggage allowance for your class</p>
      </div>
    </div>

    <!-- Footer -->
    <div style="border-top: 1px solid rgba(51, 104, 145, 0.31); padding: 17px 0; text-align: center;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Thank you for booking with us!</p>
      <p style="margin: 0 0 8px 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">For support: info@anytimetravel.app</p>
      <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone: +961 1 234 567</p>
    </div>
  </div>
</body>
</html>`,
		confirmationNumber,
		bookingDate.Format("January 2, 2006"),
		travelersHTML,
		outboundHTML,
		returnHTML,
		basePrice,
		taxes,
		serviceFee,
		totalPrice,
	)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendMultiCityFlightBookingEmail sends a multi-city flight booking confirmation email
func (es *EmailService) SendMultiCityFlightBookingEmail(
	recipientEmail, customerName, confirmationNumber string,
	bookingDate time.Time,
	travelers []string,
	allSegments []FlightSegment,
	basePrice, taxes, serviceFee, totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Multi-city flight booking confirmation for %s: %s\n", recipientEmail, confirmationNumber)
		return nil
	}

	subject := fmt.Sprintf("Flight Booking Confirmed - %s", confirmationNumber)

	// Build travelers HTML
	travelersHTML := ""
	for _, traveler := range travelers {
		travelersHTML += fmt.Sprintf(`
          <p style="margin: 0; padding: 12px 0 0 14px; font-family: 'Roboto', Arial, sans-serif; font-size: 14px; line-height: 21px; color: #1D1D1D;">%s</p>`, traveler)
	}

	// Build all segments HTML
	segmentsHTML := ""
	for i, segment := range allSegments {
		segmentsHTML += fmt.Sprintf(`
        <div style="margin-bottom: 16px;">
          <p style="margin: 0 0 8px 0; font-size: 16px; font-weight: 700; color: #336891;">Flight %d</p>
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
            <tr>
              <td style="vertical-align: top;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">%s</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              </td>
              <td width="30" style="width: 30px;"></td>
              <td style="vertical-align: top; text-align: right; white-space: nowrap;">
                <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Class</p>
                <p style="margin: 0; font-size: 14px; font-weight: 700; color: #336891;">%s</p>
              </td>
            </tr>
          </table>
          <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
            <div style="margin-bottom: 20px;">
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
            <div style="text-align: center; margin: 12px 0; color: rgba(0, 0, 0, 0.7);">
              <span style="font-size: 12px;">✈️ %s</span>
            </div>
            <div>
              <div style="display: flex; align-items: center; margin-bottom: 4px;">
                <div style="width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; margin-right: 6px;"></div>
                <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
              </div>
              <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
              <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
              <div style="margin-top: 6px; font-size: 12px; color: #1D1D1D;">
                <span>📅 %s</span>
                <span style="margin-left: 12px; font-weight: 700;">🕐 %s</span>
              </div>
            </div>
          </div>
        </div>`,
			i+1,
			segment.Airline,
			segment.FlightNumber,
			segment.Class,
			segment.DepartureCode,
			segment.DepartureCity,
			segment.DepartureAirport,
			segment.DepartureDate,
			segment.DepartureTime,
			segment.Duration,
			segment.ArrivalCode,
			segment.ArrivalCity,
			segment.ArrivalAirport,
			segment.ArrivalDate,
			segment.ArrivalTime,
		)
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; font-family: 'Roboto', Arial, sans-serif; background-color: #f4f4f4; }
    .email-container { max-width: 430px; margin: 0 auto; background: #FFFFFF; }
    @media only screen and (max-width: 430px) {
      .email-container { width: 100%% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div style="background: #336891; padding: 24px 24px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 16px;">
        <tr>
          <td style="vertical-align: middle;">
            <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.9);">ANYTIME TRAVEL & TOURISM</p>
            <h1 style="margin: 4px 0 0 0; font-size: 24px; font-weight: 700; color: #FFFFFF;">Booking Confirmed</h1>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td width="48" style="width: 48px; vertical-align: middle;">
            <div style="width: 48px; height: 48px; border: 4px solid #FFFFFF; border-radius: 50%%; text-align: center; line-height: 48px;">
              <span style="color: #FFFFFF; font-size: 28px; font-weight: bold;">✓</span>
            </div>
          </td>
        </tr>
      </table>
      
      <div style="border-top: 1px solid rgba(255, 255, 255, 0.3); padding-top: 13px; padding-bottom: 24px;">
        <div style="margin-bottom: 12px;">
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Confirmation Number</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
        <div>
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Booking Date</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
      </div>
    </div>

    <!-- Success Banner -->
    <div style="background: #E8F5E9; border-left: 2px solid #4CAF50; padding: 12px 19px;">
      <p style="margin: 0; font-size: 14px; font-weight: 700; color: #2E7D32;">✓ Booking Confirmed</p>
      <p style="margin: 4px 0 0 0; font-size: 12px; line-height: 18px; color: #1D1D1D;">Your multi-city flight has been confirmed. A copy has been sent to your email address.</p>
    </div>

    <!-- Travelers -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Travelers</h2>
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 0;">
        %s
      </div>
    </div>

    <!-- Flight Segments -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Flight Itinerary</h2>
      %s
    </div>

    <!-- Price Summary -->
    <div style="padding: 0 20px; margin: 13px 0;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Price Summary</h2>
      <div style="margin-bottom: 13px;">
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Flight Ticket</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Taxes</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
          <span style="font-size: 14px; color: #1D1D1D;">Service Fee</span>
          <span style="font-size: 14px; color: #336891; font-weight: 400;">$%.2f</span>
        </div>
        <div style="text-align: right; margin-top: 13px;">
          <span style="font-size: 14px; color: #D24124; font-weight: 700;">Total: $%.2f</span>
        </div>
      </div>
    </div>

    <!-- Important Information -->
    <div style="background: #E3F2FD; border-left: 2px solid #336891; padding: 12px 14px; margin-bottom: 17px;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Important Information</p>
      <div style="font-size: 12px; line-height: 18px; color: #1D1D1D;">
        <p style="margin: 0 0 6px 0;">• Arrive 3 hours before each departure</p>
        <p style="margin: 0 0 6px 0;">• Check-in opens 24 hours before each flight</p>
        <p style="margin: 0 0 6px 0;">• Valid passport required</p>
        <p style="margin: 0;">• Check baggage allowance for your class</p>
      </div>
    </div>

    <!-- Footer -->
    <div style="border-top: 1px solid rgba(51, 104, 145, 0.31); padding: 17px 0; text-align: center;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Thank you for booking with us!</p>
      <p style="margin: 0 0 8px 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">For support: info@anytimetravel.app</p>
      <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone: +961 1 234 567</p>
    </div>
  </div>
</body>
</html>`,
		confirmationNumber,
		bookingDate.Format("January 2, 2006"),
		travelersHTML,
		segmentsHTML,
		basePrice,
		taxes,
		serviceFee,
		totalPrice,
	)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendHotelBookingEmail sends a hotel booking confirmation email
func (es *EmailService) SendHotelBookingEmail(
	recipientEmail, customerName, confirmationNumber string,
	bookingDate time.Time,
	guestNames []string,
	hotelName, hotelImage string,
	checkInDate, checkOutDate time.Time,
	nights, rooms int,
	nightPrice, taxes, destinationFee, serviceFee, totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Hotel booking confirmation for %s: %s\n", recipientEmail, confirmationNumber)
		return nil
	}

	subject := fmt.Sprintf("Hotel Booking Confirmed - %s", confirmationNumber)

	// Build guest list HTML
	guestsHTML := ""
	if len(guestNames) > 0 {
		guestsHTML = `
        <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
          <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Travelers</h2>
          <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 0;">`

		for _, guest := range guestNames {
			guestsHTML += fmt.Sprintf(`
            <p style="margin: 0; padding: 12px 0 0 14px; font-family: 'Roboto', Arial, sans-serif; font-size: 14px; line-height: 21px; color: #1D1D1D;">%s</p>`, guest)
		}

		guestsHTML += `
          </div>
        </div>`
	}

	// Guest count display
	guestCount := len(guestNames)
	guestText := "1 Guest"
	if guestCount > 1 {
		guestText = fmt.Sprintf("%d Guests", guestCount)
	}

	// Calculate nights text
	nightsText := "1 night"
	if nights > 1 {
		nightsText = fmt.Sprintf("%d nights", nights)
	}

	// Rooms text
	roomsText := "1 room"
	if rooms > 1 {
		roomsText = fmt.Sprintf("%d rooms", rooms)
	}

	// Format check-in and check-out dates
	checkInFormatted := checkInDate.Format("Monday, January 2, 2006")
	checkOutFormatted := checkOutDate.Format("Monday, January 2, 2006")

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; font-family: 'Roboto', Arial, sans-serif; background-color: #f4f4f4; }
    .email-container { max-width: 430px; margin: 0 auto; background: #FFFFFF; }
    @media only screen and (max-width: 430px) {
      .email-container { width: 100%% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div style="background: #336891; padding: 24px 24px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 16px;">
        <tr>
          <td width="60" style="width: 60px; vertical-align: middle; padding-right: 12px;">
            <img src="https://anytimetravel.app/static/admin/images/main-logo.png" alt="Anytime Travel" style="width: 50px; height: 50px; display: block;" />
          </td>
          <td style="vertical-align: middle;">
            <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.9);">ANYTIME TRAVEL & TOURISM</p>
            <h1 style="margin: 4px 0 0 0; font-size: 24px; font-weight: 700; color: #FFFFFF;">Booking Confirmed</h1>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td width="48" style="width: 48px; vertical-align: middle;">
            <div style="width: 48px; height: 48px; border: 4px solid #FFFFFF; border-radius: 50%%; text-align: center; line-height: 40px;">
              <span style="color: #FFFFFF; font-size: 28px; font-weight: bold;">✓</span>
            </div>
          </td>
        </tr>
      </table>
      
      <div style="border-top: 1px solid rgba(255, 255, 255, 0.3); padding-top: 13px; padding-bottom: 24px;">
        <div style="margin-bottom: 12px;">
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Confirmation Number</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
        <div>
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Booking Date</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
      </div>
    </div>

    <!-- Success Banner -->
    <div style="background: #E8F5E9; border-left: 2px solid #4CAF50; padding: 12px 19px;">
      <p style="margin: 0; font-size: 14px; font-weight: 700; color: #2E7D32;">✓ Booking Confirmed</p>
      <p style="margin: 4px 0 0 0; font-size: 12px; line-height: 18px; color: #1D1D1D;">Your hotel stay has been confirmed. A copy has been sent to your email.</p>
    </div>

    <!-- Travelers -->
    %s

    <!-- Hotel Stay -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin: 13px 0;">
        <tr>
          <td style="vertical-align: middle;">
            <h2 style="margin: 0; font-size: 20px; font-weight: 700; color: #336891;">Hotel Stay</h2>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td style="vertical-align: middle; text-align: right;">
            <p style="margin: 0; font-size: 14px; color: #1D1D1D;">%s</p>
          </td>
        </tr>
      </table>
      
      <!-- Hotel Image and Info -->
      <div style="background: #FFFFFF; border-radius: 9px; overflow: hidden; margin-bottom: 13px;">
        <img src="%s" alt="%s" style="width: 100%%; height: 180px; object-fit: cover; display: block;" />
        <div style="padding: 14px 16px;">
          <h3 style="margin: 0 0 12px 0; font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</h3>
          
          <!-- Check-in -->
          <div style="background: #D4DFE8; border-radius: 9px; padding: 12px 14px; margin-bottom: 8px;">
            <p style="margin: 0; font-size: 12px; color: #336891;">Check-in</p>
            <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
          </div>
          
          <!-- Check-out -->
          <div style="background: #D4DFE8; border-radius: 9px; padding: 12px 14px; margin-bottom: 8px;">
            <p style="margin: 0; font-size: 12px; color: #336891;">Check-out</p>
            <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
          </div>
          
          <p style="margin: 8px 0 0 0; font-size: 14px; color: #1D1D1D; text-align: right;">%s, %s</p>
        </div>
      </div>
    </div>

    <!-- Price Summary -->
    <div style="padding: 0 20px; margin: 13px 0;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Price Summary</h2>
      <div style="margin-bottom: 13px;">
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">%s</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">Taxes</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">Destination Fee</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">Service fees</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
        </table>
        <div style="text-align: right; margin-top: 13px;">
          <span style="font-size: 14px; color: #D24124; font-weight: 700;">Total: $%.2f</span>
        </div>
      </div>
    </div>

    <!-- Important Information -->
    <div style="background: #E3F2FD; border-left: 2px solid #336891; padding: 12px 14px; margin-bottom: 17px;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Important Information</p>
      <div style="font-size: 12px; line-height: 18px; color: #1D1D1D;">
        <p style="margin: 0 0 6px 0;">• Check-in time: 3:00 PM</p>
        <p style="margin: 0 0 6px 0;">• Check-out time: 11:00 AM</p>
        <p style="margin: 0 0 6px 0;">• Valid ID required at check-in</p>
        <p style="margin: 0 0 6px 0;">• Credit card may be required for incidentals</p>
        <p style="margin: 0;">• Free cancellation up to 48 hours before check-in</p>
      </div>
    </div>

    <!-- Footer -->
    <div style="border-top: 1px solid rgba(51, 104, 145, 0.31); padding: 17px 0; text-align: center;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Thank you for booking with us!</p>
      <p style="margin: 0 0 8px 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">For support: support@anytimetravel.com</p>
      <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone: +961 1 234 567</p>
    </div>
  </div>
</body>
</html>`,
		confirmationNumber,
		bookingDate.Format("January 2, 2006"),
		guestsHTML,
		guestText,
		hotelImage,
		hotelName,
		hotelName,
		checkInFormatted,
		checkOutFormatted,
		nightsText,
		roomsText,
		nightsText,
		nightPrice,
		taxes,
		destinationFee,
		serviceFee,
		totalPrice,
	)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendCarBookingEmail sends a car rental booking confirmation email
func (es *EmailService) SendCarBookingEmail(
	recipientEmail, customerName, confirmationNumber string,
	bookingDate time.Time,
	carType string,
	passengers int,
	pickupLocation, pickupAddress, pickupDate, pickupTime string,
	dropoffLocation, dropoffAddress, dropoffDate, dropoffTime string,
	driverName, driverPhone, driverLicense string,
	rentalPrice, taxes, totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Car booking confirmation for %s: %s\n", recipientEmail, confirmationNumber)
		return nil
	}

	subject := fmt.Sprintf("Car Rental Confirmed - %s", confirmationNumber)

	// Passenger text
	passengerText := "1 Passenger"
	if passengers > 1 {
		passengerText = fmt.Sprintf("%d Passengers", passengers)
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; font-family: 'Roboto', Arial, sans-serif; background-color: #f4f4f4; }
    .email-container { max-width: 430px; margin: 0 auto; background: #FFFFFF; }
    @media only screen and (max-width: 430px) {
      .email-container { width: 100%% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div style="background: #336891; padding: 24px 24px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 16px;">
        <tr>
          <td width="60" style="width: 60px; vertical-align: middle; padding-right: 12px;">
            <img src="https://anytimetravel.app/static/admin/images/main-logo.png" alt="Anytime Travel" style="width: 50px; height: 50px; display: block;" />
          </td>
          <td style="vertical-align: middle;">
            <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.9);">ANYTIME TRAVEL & TOURISM</p>
            <h1 style="margin: 4px 0 0 0; font-size: 24px; font-weight: 700; color: #FFFFFF;">Booking Confirmed</h1>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td width="48" style="width: 48px; vertical-align: middle;">
            <div style="width: 48px; height: 48px; border: 4px solid #FFFFFF; border-radius: 50%%; text-align: center; line-height: 40px;">
              <span style="color: #FFFFFF; font-size: 28px; font-weight: bold;">✓</span>
            </div>
          </td>
        </tr>
      </table>
      
      <div style="border-top: 1px solid rgba(255, 255, 255, 0.3); padding-top: 13px; padding-bottom: 24px;">
        <div style="margin-bottom: 12px;">
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Confirmation Number</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
        <div>
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Booking Date</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
      </div>
    </div>

    <!-- Success Banner -->
    <div style="background: #E8F5E9; border-left: 2px solid #4CAF50; padding: 12px 19px;">
      <p style="margin: 0; font-size: 14px; font-weight: 700; color: #2E7D32;">✓ Booking Confirmed</p>
      <p style="margin: 4px 0 0 0; font-size: 12px; line-height: 18px; color: #1D1D1D;">Your car rental has been confirmed. A copy has been sent to your email.</p>
    </div>

    <!-- Car Details -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin: 13px 0;">
        <tr>
          <td style="vertical-align: middle;">
            <h2 style="margin: 0; font-size: 20px; font-weight: 700; color: #336891;">%s</h2>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td style="vertical-align: middle; text-align: right;">
            <p style="margin: 0; font-size: 14px; color: #1D1D1D;">%s</p>
          </td>
        </tr>
      </table>
      
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 12px 0 14px;">
        <p style="margin: 0; padding-bottom: 2px; font-size: 14px; color: #1D1D1D;">%s</p>
        <p style="margin: 0; font-size: 14px; color: #1D1D1D;">%s - %s</p>
      </div>
    </div>

    <!-- Pickup Details -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Pickup Details</h2>
      <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
        <div style="margin-bottom: 4px;">
          <div style="display: inline-block; width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; vertical-align: middle; margin-right: 6px;"></div>
          <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
        </div>
        <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
        <div style="margin-top: 6px;">
          <span style="font-size: 12px; color: #1D1D1D;">📅 %s</span>
        </div>
        <div style="margin-top: 6px;">
          <span style="font-size: 12px; color: #1D1D1D;">🕐 </span>
          <span style="font-size: 12px; font-weight: 700; color: #1D1D1D;">%s</span>
        </div>
      </div>
    </div>

    <!-- Dropoff Details -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Dropoff Details</h2>
      <div style="background: #D4DFE8; border-radius: 9px; padding: 14px 16px;">
        <div style="margin-bottom: 4px;">
          <div style="display: inline-block; width: 14px; height: 14px; border: 1.17px solid #D24124; border-radius: 50%%; vertical-align: middle; margin-right: 6px;"></div>
          <span style="font-size: 18px; font-weight: 700; color: #1D1D1D;">%s</span>
        </div>
        <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
        <div style="margin-top: 6px;">
          <span style="font-size: 12px; color: #1D1D1D;">📅 %s</span>
        </div>
        <div style="margin-top: 6px;">
          <span style="font-size: 12px; color: #1D1D1D;">🕐 </span>
          <span style="font-size: 12px; font-weight: 700; color: #1D1D1D;">%s</span>
        </div>
      </div>
    </div>

    <!-- Who's Driving -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Who's Driving?</h2>
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 12px 0 14px;">
        <div style="margin-bottom: 8px;">
          <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Driver Name</p>
          <p style="margin: 0; font-size: 14px; color: #1D1D1D;">%s</p>
        </div>
        <div style="margin-bottom: 8px;">
          <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone Number</p>
          <p style="margin: 0; font-size: 14px; color: #1D1D1D;">%s</p>
        </div>
        <div style="padding-bottom: 8px;">
          <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">License Number</p>
          <p style="margin: 0; font-size: 14px; color: #1D1D1D;">%s</p>
        </div>
      </div>
    </div>

    <!-- Price Summary -->
    <div style="padding: 0 20px; margin: 13px 0;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Price Summary</h2>
      <div style="margin-bottom: 13px;">
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">Rental Price</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">Taxes</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
        </table>
        <div style="text-align: right; margin-top: 13px;">
          <span style="font-size: 14px; color: #D24124; font-weight: 700;">Total: $%.2f</span>
        </div>
      </div>
    </div>

    <!-- Important Information -->
    <div style="background: #E3F2FD; border-left: 2px solid #336891; padding: 12px 14px; margin-bottom: 17px;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Important Information</p>
      <div style="font-size: 12px; line-height: 18px; color: #1D1D1D;">
        <p style="margin: 0 0 6px 0;">• Valid driver's license required</p>
        <p style="margin: 0 0 6px 0;">• Fuel policy: Pick up full, return full</p>
        <p style="margin: 0 0 6px 0;">• Credit card required for deposit</p>
        <p style="margin: 0 0 6px 0;">• Free cancellation up to 24 hours before pickup</p>
        <p style="margin: 0;">• Additional drivers can be added at pickup</p>
      </div>
    </div>

    <!-- Footer -->
    <div style="border-top: 1px solid rgba(51, 104, 145, 0.31); padding: 17px 0; text-align: center;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Thank you for booking with us!</p>
      <p style="margin: 0 0 8px 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">For support: support@anytimetravel.com</p>
      <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone: +961 1 234 567</p>
    </div>
  </div>
</body>
</html>`,
		confirmationNumber,
		bookingDate.Format("January 2, 2006"),
		carType,
		passengerText,
		pickupLocation,
		pickupTime,
		dropoffTime,
		pickupLocation,
		pickupAddress,
		pickupDate,
		pickupTime,
		dropoffLocation,
		dropoffAddress,
		dropoffDate,
		dropoffTime,
		driverName,
		driverPhone,
		driverLicense,
		rentalPrice,
		taxes,
		totalPrice,
	)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendTransferBookingEmail sends a transfer booking confirmation email
func (es *EmailService) SendTransferBookingEmail(
	recipientEmail, customerName, customerPhone, confirmationNumber string,
	bookingDate time.Time,
	vehicleType, vehicleModel, vehicleDuration, vehicleImage string,
	vehiclePassengers int,
	vehicleMeetGreet bool,
	vehiclePrice float64,
	pickupName, pickupAddress, pickupDate, pickupTime string,
	dropoffName, dropoffAddress, dropoffDate, dropoffTime string,
	basePrice, taxes, taxOnFees, totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Transfer booking confirmation for %s: %s\n", recipientEmail, confirmationNumber)
		return nil
	}

	subject := fmt.Sprintf("Transfer Booking Confirmed - %s", confirmationNumber)

	// Passenger text
	passengerText := fmt.Sprintf("%d Passengers", vehiclePassengers)
	if vehiclePassengers == 1 {
		passengerText = "1 Passenger"
	}

	// Meet & Greet text
	meetGreetText := ""
	if vehicleMeetGreet {
		meetGreetText = `<p style="margin: 0; font-size: 14px; line-height: 16px; color: #1D1D1D;">Meet & Greet availability</p>`
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; font-family: 'Roboto', Arial, sans-serif; background-color: #f4f4f4; }
    .email-container { max-width: 430px; margin: 0 auto; background: #FFFFFF; }
    @media only screen and (max-width: 430px) {
      .email-container { width: 100%% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div style="background: #336891; padding: 24px 24px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 16px;">
        <tr>
          <td width="60" style="width: 60px; vertical-align: middle; padding-right: 12px;">
            <img src="https://anytimetravel.app/static/admin/images/main-logo.png" alt="Anytime Travel" style="width: 50px; height: 50px; display: block;" />
          </td>
          <td style="vertical-align: middle;">
            <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.9);">ANYTIME TRAVEL & TOURISM</p>
            <h1 style="margin: 4px 0 0 0; font-size: 24px; font-weight: 700; color: #FFFFFF;">Booking Confirmed</h1>
          </td>
          <td width="20" style="width: 20px;"></td>
          <td width="48" style="width: 48px; vertical-align: middle;">
            <div style="width: 48px; height: 48px; border: 4px solid #FFFFFF; border-radius: 50%%; text-align: center; line-height: 40px;">
              <span style="color: #FFFFFF; font-size: 28px; font-weight: bold;">✓</span>
            </div>
          </td>
        </tr>
      </table>
      
      <div style="border-top: 1px solid rgba(255, 255, 255, 0.3); padding-top: 13px; padding-bottom: 24px;">
        <div style="margin-bottom: 12px;">
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Confirmation Number</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
        <div>
          <p style="margin: 0; font-size: 12px; color: rgba(255, 255, 255, 0.75);">Booking Date</p>
          <p style="margin: 4px 0 0 0; font-size: 18px; font-weight: 700; color: #FFFFFF;">%s</p>
        </div>
      </div>
    </div>

    <!-- Success Banner -->
    <div style="background: #E8F5E9; border-left: 2px solid #4CAF50; padding: 12px 0 0 19px; margin: 20px 20px 0;">
      <p style="margin: 0; font-size: 14px; font-weight: 700; color: #2E7D32;">✓ Booking Confirmed</p>
      <p style="margin: 4px 0 0 0; font-size: 12px; line-height: 18px; color: #1D1D1D; padding-bottom: 12px;">Your transfer has been confirmed. A copy has been sent to your email.</p>
    </div>

    <!-- Who's Booking -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px; margin-top: 20px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Who's booking?</h2>
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 7px 12px 0 14px;">
        <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Primary Contact</p>
        <p style="margin: 0; padding: 8px 0; font-size: 14px; color: #1D1D1D;">%s</p>
        <p style="margin: 0; padding-bottom: 8px; font-size: 14px; color: #1D1D1D;">%s</p>
      </div>
    </div>

    <!-- Transfer Details -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin: 13px 0;">Transfer Details</h2>
      <div style="background: #D4DFE8; border-radius: 9px; padding: 5px 14px 14px; position: relative;">
        <img src="%s" alt="Vehicle" style="width: 100%%; height: auto; border-radius: 8px; margin-bottom: 12px;" />
        <div style="position: relative;">
          <p style="margin: 0; font-size: 14px; font-weight: 700; line-height: 17px; color: #D24124;">%s</p>
          <p style="margin: 0; padding-top: 4px; font-size: 14px; font-weight: 300; line-height: 16px; color: #1D1D1D;">%s</p>
          <p style="margin: 0; padding-top: 4px; font-size: 14px; line-height: 16px; color: #1D1D1D;">%s</p>
          <p style="margin: 0; padding-top: 4px; font-size: 14px; line-height: 16px; color: #1D1D1D;">Estimated time: %s</p>
          %s
          <p style="margin: 0; padding: 5px 0 10px 0; font-size: 20px; font-weight: 700; text-align: right; color: #1D1D1D;">$%.2f</p>
        </div>
      </div>
    </div>

    <!-- Pickup -->
    <div style="padding: 0 20px;">
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 10px 12px 0 14px; margin-bottom: 13px;">
        <p style="margin: 0; font-size: 14px; font-weight: 700; color: #D24124;">Pickup</p>
        <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
        <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
        <p style="margin: 4px 0 0 0; font-size: 12px; color: #1D1D1D;">%s</p>
        <p style="margin: 4px 0 0 0; padding-bottom: 10px; font-size: 12px; font-weight: 700; color: #1D1D1D;">%s</p>
      </div>
    </div>

    <!-- Dropoff -->
    <div style="padding: 0 20px; border-bottom: 1px solid rgba(51, 104, 145, 0.31); padding-bottom: 13px;">
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 10px 12px 0 14px;">
        <p style="margin: 0; font-size: 14px; font-weight: 700; color: #D24124;">Dropoff</p>
        <p style="margin: 4px 0 0 0; font-size: 14px; font-weight: 700; color: #1D1D1D;">%s</p>
        <p style="margin: 4px 0 0 0; font-size: 12px; font-weight: 300; color: #1D1D1D;">%s</p>
        <p style="margin: 4px 0 0 0; font-size: 12px; color: #1D1D1D;">%s</p>
        <p style="margin: 4px 0 0 0; padding-bottom: 10px; font-size: 12px; font-weight: 700; color: #1D1D1D;">%s</p>
      </div>
    </div>

    <!-- Price Summary -->
    <div style="padding: 0 20px; margin: 13px 0;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 13px;">Price Summary</h2>
      <div style="margin-bottom: 13px;">
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom: 8px;">
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">%s</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">Taxes</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
          <tr>
            <td style="font-size: 14px; color: #1D1D1D; padding-bottom: 8px;">Tax on fees</td>
            <td style="font-size: 14px; color: #336891; font-weight: 400; text-align: right; padding-bottom: 8px;">$%.2f</td>
          </tr>
        </table>
        <div style="text-align: right; margin-top: 13px;">
          <span style="font-size: 14px; color: #D24124; font-weight: 700;">Total: $%.2f</span>
        </div>
      </div>
    </div>

    <!-- Important Information -->
    <div style="padding: 0 20px;">
      <h2 style="font-size: 20px; font-weight: 700; color: #336891; margin-bottom: 8px;">Important Information</h2>
      <div style="background: #FFFFFF; border-left: 2px solid #D24124; padding: 12px 14px 0;">
        <div style="font-size: 12px; line-height: 18px; color: #1D1D1D;">
          <p style="margin: 0 0 6px 0;">• Driver will meet you at arrivals with a name sign</p>
          <p style="margin: 0 0 6px 0;">• Free waiting time: 60 minutes for airport pickups</p>
          <p style="margin: 0 0 6px 0;">• Driver contact details will be sent 24 hours before pickup</p>
          <p style="margin: 0 0 6px 0;">• All-inclusive price - no hidden fees</p>
          <p style="margin: 0; padding-bottom: 12px;">• Free cancellation up to 24 hours before pickup</p>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div style="border-top: 1px solid rgba(51, 104, 145, 0.31); padding: 17px 0; text-align: center; margin-top: 20px;">
      <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: 700; color: #336891;">Thank you for booking with us!</p>
      <p style="margin: 0 0 8px 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">For support: support@anytimetravel.com</p>
      <p style="margin: 0; font-size: 12px; color: rgba(0, 0, 0, 0.7);">Phone: +961 1 234 567</p>
    </div>
  </div>
</body>
</html>`,
		confirmationNumber,
		bookingDate.Format("January 2, 2006"),
		customerName,
		customerPhone,
		vehicleImage,
		vehicleType,
		vehicleModel,
		passengerText,
		vehicleDuration,
		meetGreetText,
		vehiclePrice,
		pickupName,
		pickupAddress,
		pickupDate,
		pickupTime,
		dropoffName,
		dropoffAddress,
		dropoffDate,
		dropoffTime,
		passengerText,
		basePrice,
		taxes,
		taxOnFees,
		totalPrice,
	)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendCarBookingReminder sends a car booking reminder email 48 hours before pickup
func (es *EmailService) SendCarBookingReminder(
	recipientEmail, customerName, bookingID, carType string,
	pickupLocation, pickupAddress, pickupDate, pickupTime string,
	totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Car booking reminder for %s: %s\n", recipientEmail, bookingID)
		return nil
	}

	subject := fmt.Sprintf("Reminder: Your Car Booking %s", bookingID)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 20px auto; background-color: white; }
    .header { background-color: #336891; color: white; padding: 30px 20px; text-align: center; }
    .header h1 { margin: 0; font-size: 28px; }
    .content { padding: 30px 20px; }
    .reminder-box { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 20px; margin: 20px 0; }
    .booking-details { background-color: #f9f9f9; border-left: 4px solid #336891; padding: 20px; margin: 20px 0; }
    .detail-row { margin: 10px 0; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { background-color: #f4f4f4; text-align: center; padding: 20px; color: #666; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header" style="display: flex; align-items: center; justify-content: center; gap: 15px;">
      <img src="https://anytimetravel.app/static/admin/images/main-logo.png" alt="Anytime Travel" style="width: 50px; height: 50px; vertical-align: middle;" />
      <h1>🚗 Upcoming Car Rental</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>
      
      <div class="reminder-box">
        <strong>⏰ Reminder:</strong> Your car rental is scheduled in 48 hours!
      </div>

      <p>This is a friendly reminder about your upcoming car rental reservation.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #336891;">Booking Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Car Type:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Location:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Address:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Date & Time:</span>
          <span class="detail-value">%s at %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Total Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p><strong>Important:</strong> Please ensure you have all required documents (driver's license, ID, payment method) ready for pickup.</p>
      
      <p>If you need to make any changes to your booking or have questions, please contact us as soon as possible.</p>

      <p>Safe travels,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, carType, pickupLocation, pickupAddress, pickupDate, pickupTime, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendFlightBookingReminder sends a flight booking reminder email 48 hours before departure
func (es *EmailService) SendFlightBookingReminder(
	recipientEmail, customerName, bookingID string,
	airline, flightNumber, departureCity, arrivalCity string,
	departureDate, departureTime string,
	totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Flight booking reminder for %s: %s\n", recipientEmail, bookingID)
		return nil
	}

	subject := fmt.Sprintf("Reminder: Your Flight %s", bookingID)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 20px auto; background-color: white; }
    .header { background-color: #336891; color: white; padding: 30px 20px; text-align: center; }
    .header h1 { margin: 0; font-size: 28px; }
    .content { padding: 30px 20px; }
    .reminder-box { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 20px; margin: 20px 0; }
    .booking-details { background-color: #f9f9f9; border-left: 4px solid #336891; padding: 20px; margin: 20px 0; }
    .detail-row { margin: 10px 0; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { background-color: #f4f4f4; text-align: center; padding: 20px; color: #666; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header" style="display: flex; align-items: center; justify-content: center; gap: 15px;">
      <img src="https://anytimetravel.app/static/admin/images/main-logo.png" alt="Anytime Travel" style="width: 50px; height: 50px; vertical-align: middle;" />
      <h1>✈️ Upcoming Flight</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>
      
      <div class="reminder-box">
        <strong>⏰ Reminder:</strong> Your flight departs in 48 hours!
      </div>

      <p>This is a friendly reminder about your upcoming flight reservation.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #336891;">Flight Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Flight:</span>
          <span class="detail-value">%s %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Route:</span>
          <span class="detail-value">%s → %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Departure Date & Time:</span>
          <span class="detail-value">%s at %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Total Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p><strong>Check-in Reminder:</strong> Please check in online 24 hours before departure and arrive at the airport at least 2-3 hours early for international flights.</p>
      
      <p>If you need assistance or have questions, please contact us.</p>

      <p>Have a wonderful flight,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, airline, flightNumber, departureCity, arrivalCity, departureDate, departureTime, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendHotelBookingReminder sends a hotel booking reminder email 48 hours before check-in
func (es *EmailService) SendHotelBookingReminder(
	recipientEmail, customerName, bookingID, hotelName string,
	checkInDate string,
	nights int,
	totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Hotel booking reminder for %s: %s\n", recipientEmail, bookingID)
		return nil
	}

	subject := fmt.Sprintf("Reminder: Your Hotel Booking %s", bookingID)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 20px auto; background-color: white; }
    .header { background-color: #336891; color: white; padding: 30px 20px; text-align: center; }
    .header h1 { margin: 0; font-size: 28px; }
    .content { padding: 30px 20px; }
    .reminder-box { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 20px; margin: 20px 0; }
    .booking-details { background-color: #f9f9f9; border-left: 4px solid #336891; padding: 20px; margin: 20px 0; }
    .detail-row { margin: 10px 0; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { background-color: #f4f4f4; text-align: center; padding: 20px; color: #666; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header" style="display: flex; align-items: center; justify-content: center; gap: 15px;">
      <img src="https://anytimetravel.app/static/admin/images/main-logo.png" alt="Anytime Travel" style="width: 50px; height: 50px; vertical-align: middle;" />
      <h1>🏨 Upcoming Hotel Stay</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>
      
      <div class="reminder-box">
        <strong>⏰ Reminder:</strong> Your hotel check-in is in 48 hours!
      </div>

      <p>This is a friendly reminder about your upcoming hotel reservation.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #336891;">Booking Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Hotel:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Check-in Date:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Number of Nights:</span>
          <span class="detail-value">%d</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Total Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p><strong>Important:</strong> Standard check-in time is usually 3:00 PM. Please have your ID and booking confirmation ready.</p>
      
      <p>If you need to arrange early check-in or have special requests, please contact the hotel directly.</p>

      <p>Enjoy your stay,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, hotelName, checkInDate, nights, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendTransferBookingReminder sends a transfer booking reminder email 48 hours before pickup
func (es *EmailService) SendTransferBookingReminder(
	recipientEmail, customerName, bookingID, vehicleType string,
	pickupLocation, pickupAddress, pickupDate, pickupTime string,
	totalPrice float64,
) error {
	if es.senderPassword == "" {
		fmt.Printf("[EMAIL] Transfer booking reminder for %s: %s\n", recipientEmail, bookingID)
		return nil
	}

	subject := fmt.Sprintf("Reminder: Your Transfer Booking %s", bookingID)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 20px auto; background-color: white; }
    .header { background-color: #336891; color: white; padding: 30px 20px; text-align: center; }
    .header h1 { margin: 0; font-size: 28px; }
    .content { padding: 30px 20px; }
    .reminder-box { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 20px; margin: 20px 0; }
    .booking-details { background-color: #f9f9f9; border-left: 4px solid #336891; padding: 20px; margin: 20px 0; }
    .detail-row { margin: 10px 0; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { background-color: #f4f4f4; text-align: center; padding: 20px; color: #666; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header" style="display: flex; align-items: center; justify-content: center; gap: 15px;">
      <img src="https://anytimetravel.app/static/admin/images/main-logo.png" alt="Anytime Travel" style="width: 50px; height: 50px; vertical-align: middle;" />
      <h1>🚖 Upcoming Transfer</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>
      
      <div class="reminder-box">
        <strong>⏰ Reminder:</strong> Your transfer is scheduled in 48 hours!
      </div>

      <p>This is a friendly reminder about your upcoming transfer reservation.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #336891;">Booking Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Vehicle Type:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Location:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Address:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Date & Time:</span>
          <span class="detail-value">%s at %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Total Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p><strong>Important:</strong> Your driver will be waiting at the specified pickup location. Please be ready 5-10 minutes before the scheduled time.</p>
      
      <p>If you need to make any changes or have questions, please contact us immediately.</p>

      <p>Safe travels,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, vehicleType, pickupLocation, pickupAddress, pickupDate, pickupTime, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendCarBookingCancellation sends a car booking cancellation confirmation email
func (es *EmailService) SendCarBookingCancellation(recipientEmail, customerName, bookingID, carType, pickupLocation, pickupDate, pickupTime string, totalPrice float64) error {
	subject := "Car Booking Cancelled - Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f9f9f9; }
    .header { background-color: #d32f2f; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
    .content { background-color: white; padding: 30px; border-radius: 0 0 8px 8px; }
    .booking-details { background-color: #f5f5f5; padding: 20px; margin: 20px 0; border-radius: 8px; }
    .detail-row { margin: 10px 0; display: flex; justify-content: space-between; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { text-align: center; margin-top: 30px; padding: 20px; color: #777; font-size: 12px; }
    .cancelled-badge { background-color: #d32f2f; color: white; padding: 10px 20px; border-radius: 20px; display: inline-block; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1 style="margin: 0;">Booking Cancelled</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>

      <div class="cancelled-badge">✓ Booking Successfully Cancelled</div>

      <p>Your car booking has been cancelled as requested.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #d32f2f;">Cancelled Booking Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Car Type:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Location:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Date & Time:</span>
          <span class="detail-value">%s at %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p>If this cancellation was made in error or you have any questions, please contact our support team immediately.</p>

      <p>We hope to serve you again soon!</p>

      <p>Best regards,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, carType, pickupLocation, pickupDate, pickupTime, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendFlightBookingCancellation sends a flight booking cancellation confirmation email
func (es *EmailService) SendFlightBookingCancellation(recipientEmail, customerName, bookingID, airline, flightNumber, departureCity, arrivalCity, departureDate, departureTime string, totalPrice float64) error {
	subject := "Flight Booking Cancelled - Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f9f9f9; }
    .header { background-color: #d32f2f; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
    .content { background-color: white; padding: 30px; border-radius: 0 0 8px 8px; }
    .booking-details { background-color: #f5f5f5; padding: 20px; margin: 20px 0; border-radius: 8px; }
    .detail-row { margin: 10px 0; display: flex; justify-content: space-between; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { text-align: center; margin-top: 30px; padding: 20px; color: #777; font-size: 12px; }
    .cancelled-badge { background-color: #d32f2f; color: white; padding: 10px 20px; border-radius: 20px; display: inline-block; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1 style="margin: 0;">Flight Booking Cancelled</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>

      <div class="cancelled-badge">✓ Booking Successfully Cancelled</div>

      <p>Your flight booking has been cancelled as requested.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #d32f2f;">Cancelled Flight Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Airline:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Flight Number:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Route:</span>
          <span class="detail-value">%s → %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Departure:</span>
          <span class="detail-value">%s at %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p><strong>Note:</strong> Please check with the airline regarding their cancellation policy and potential refunds or credits.</p>

      <p>If this cancellation was made in error or you have any questions, please contact our support team immediately.</p>

      <p>Best regards,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, airline, flightNumber, departureCity, arrivalCity, departureDate, departureTime, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendHotelBookingCancellation sends a hotel booking cancellation confirmation email
func (es *EmailService) SendHotelBookingCancellation(recipientEmail, customerName, bookingID, hotelName, checkInDate string, nights int, totalPrice float64) error {
	subject := "Hotel Booking Cancelled - Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f9f9f9; }
    .header { background-color: #d32f2f; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
    .content { background-color: white; padding: 30px; border-radius: 0 0 8px 8px; }
    .booking-details { background-color: #f5f5f5; padding: 20px; margin: 20px 0; border-radius: 8px; }
    .detail-row { margin: 10px 0; display: flex; justify-content: space-between; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { text-align: center; margin-top: 30px; padding: 20px; color: #777; font-size: 12px; }
    .cancelled-badge { background-color: #d32f2f; color: white; padding: 10px 20px; border-radius: 20px; display: inline-block; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1 style="margin: 0;">Hotel Booking Cancelled</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>

      <div class="cancelled-badge">✓ Booking Successfully Cancelled</div>

      <p>Your hotel reservation has been cancelled as requested.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #d32f2f;">Cancelled Reservation Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Hotel:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Check-in Date:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Duration:</span>
          <span class="detail-value">%d night(s)</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p><strong>Note:</strong> Please check the hotel's cancellation policy regarding potential refunds or credits.</p>

      <p>If this cancellation was made in error or you have any questions, please contact our support team immediately.</p>

      <p>Best regards,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, hotelName, checkInDate, nights, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

// SendTransferBookingCancellation sends a transfer booking cancellation confirmation email
func (es *EmailService) SendTransferBookingCancellation(recipientEmail, customerName, bookingID, vehicleType, pickupLocation, pickupAddress, pickupDate, pickupTime string, totalPrice float64) error {
	subject := "Transfer Booking Cancelled - Anytime Travel"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f9f9f9; }
    .header { background-color: #d32f2f; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
    .content { background-color: white; padding: 30px; border-radius: 0 0 8px 8px; }
    .booking-details { background-color: #f5f5f5; padding: 20px; margin: 20px 0; border-radius: 8px; }
    .detail-row { margin: 10px 0; display: flex; justify-content: space-between; }
    .detail-label { font-weight: bold; color: #555; }
    .detail-value { color: #333; }
    .footer { text-align: center; margin-top: 30px; padding: 20px; color: #777; font-size: 12px; }
    .cancelled-badge { background-color: #d32f2f; color: white; padding: 10px 20px; border-radius: 20px; display: inline-block; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1 style="margin: 0;">Transfer Booking Cancelled</h1>
    </div>
    <div class="content">
      <p>Dear %s,</p>

      <div class="cancelled-badge">✓ Booking Successfully Cancelled</div>

      <p>Your transfer booking has been cancelled as requested.</p>
      
      <div class="booking-details">
        <h2 style="margin-top: 0; color: #d32f2f;">Cancelled Transfer Details</h2>
        <div class="detail-row">
          <span class="detail-label">Booking ID:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Vehicle Type:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Location:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Address:</span>
          <span class="detail-value">%s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Pickup Date & Time:</span>
          <span class="detail-value">%s at %s</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Amount:</span>
          <span class="detail-value">$%.2f</span>
        </div>
      </div>

      <p>If this cancellation was made in error or you have any questions, please contact our support team immediately.</p>

      <p>Best regards,<br><strong>Anytime Travel Team</strong></p>
    </div>
    <div class="footer">
      <p>&copy; 2026 Anytime Travel. All rights reserved.</p>
      <p>info@anytimetravel.app | www.anytimetravel.app</p>
    </div>
  </div>
</body>
</html>
`, customerName, bookingID, vehicleType, pickupLocation, pickupAddress, pickupDate, pickupTime, totalPrice)

	return es.sendEmail(recipientEmail, subject, body)
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
