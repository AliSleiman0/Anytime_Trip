class Endpoints {
  // Base URL - Use 10.0.2.2 for Android Emulator, localhost for iOS Simulator
  // For physical device, use your computer's IP address (e.g., 192.168.1.x)
  static const String baseUrl = 'http://10.0.2.2:8080/api';

  // Auth endpoints
  static const String login = '/app/login';
  static const String register = '/app/signup';
  static const String googleSignIn = '/app/google-signin';
  static const String appleSignIn = '/app/apple-signin';
  static const String logout = '/app/logout';
  static const String sendOTP = '/app/send-otp';
  static const String verifyOTP = '/app/verify-otp';
  static const String activateAccount = '/app/activate-account';
  static const String resetPassword = '/app/reset-password';
  static const String forgotPassword = '/app/forgot-password';
  static const String verifyResetCode = '/app/verify-reset-code';

  // Product endpoints
  static const String products = '/app/products';
  static String productDetails(String id) => '/app/products/$id';

  // Account endpoints
  static const String profile = '/app/profile';
  static const String updateProfile = '/app/profile/update';
  static const String uploadProfileImage = '/app/profile/upload-image';
  static const String getNotificationPreferences = '/app/notification-preferences';
  static const String saveNotificationPreferences = '/app/notification-preferences';
  static const String getSecurityPreferences = '/app/security-preferences';
  static const String saveSecurityPreferences = '/app/security-preferences';
  static const String paymentMethods = '/app/payment-methods';
  static String paymentMethodDetails(String id) => '/app/payment-methods/$id';
  static String setDefaultPaymentMethod(String id) => '/app/payment-methods/$id/set-default';

  // Order endpoints
  static const String orders = '/app/orders';
  static String orderDetails(String id) => '/app/orders/$id';
  static const String createOrder = '/app/orders/create';

  // Banner endpoints
  static const String banners = '/app/banners';
  static const String homepageBanners = '/app/homepage-banners';
  static const String searchBanners = '/app/search-banners';
  
  // Popular locations endpoint
  static const String popularLocations = '/app/popular-locations';

  // Car search endpoint
  static const String searchCars = '/app/search-cars';

  // Booking endpoints
  static const String myBookings = '/app/my-bookings';
  static const String myCarBookings = '/app/my-bookings/car';
  static const String myFlightBookings = '/app/my-bookings/flight';
  static const String myHotelBookings = '/app/my-bookings/hotel';
  static const String myTransferBookings = '/app/my-bookings/transfer';
  static String bookingDetails(String type, String id) => '/app/bookings/$type/$id';
  static String cancelBooking(String type, String id) => '/app/bookings/$type/$id/cancel';

  // Support endpoints
  static const String createSupportTicket = '/app/support/tickets';
  static const String userSupportTickets = '/app/support/tickets';
  static String supportTicketDetails(String ticketId) => '/app/support/tickets/$ticketId';
  static String replyToTicket(String ticketId) => '/app/support/tickets/$ticketId/reply';
  static String chatHistory(String ticketId) => '/app/support/chat/history?ticket_id=$ticketId';

  // Chatbot rating endpoints
  static const String submitChatbotRating = '/app/chatbot/rating';
  static const String getChatbotRatings = '/app/chatbot/ratings';
}
