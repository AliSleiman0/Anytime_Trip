class Endpoints {
  // Base URL - Use 10.0.2.2 for Android Emulator, localhost for iOS Simulator
  // For physical device, use your computer's IP address (e.g., 192.168.1.x)
  static const String baseUrl = 'http://10.0.2.2:8080/api';

  // Auth endpoints
  static const String login = '/app/login';
  static const String register = '/app/signup';
  static const String googleSignIn = '/app/google-signin';
  static const String logout = '/app/logout';
  static const String sendOTP = '/app/send-otp';
  static const String verifyOTP = '/app/verify-otp';
  static const String activateAccount = '/app/activate-account';
  static const String resetPassword = '/app/reset-password';

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
}
