class Endpoints {
  // Base URL
  static const String baseUrl = 'http://localhost:8080/api';

  // Auth endpoints
  static const String login = '/app/login';
  static const String register = '/app/register';
  static const String logout = '/app/logout';

  // Product endpoints
  static const String products = '/app/products';
  static String productDetails(String id) => '/app/products/$id';

  // Account endpoints
  static const String profile = '/app/profile';
  static const String updateProfile = '/app/profile/update';

  // Order endpoints
  static const String orders = '/app/orders';
  static String orderDetails(String id) => '/app/orders/$id';
  static const String createOrder = '/app/orders/create';
}
