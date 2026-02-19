import '../network/api_client.dart';
import '../network/endpoints.dart';

class CarService {
  final ApiClient _apiClient = ApiClient();

  // Search for available cars
  Future<List<dynamic>> searchCars({
    required String pickupLocation,
    String? dropoffLocation,
    String? pickupTime,
    String? dropoffTime,
    String? carType,
    int? passengers,
  }) async {
    try {
      // Build query parameters
      Map<String, dynamic> queryParams = {
        'pickup_location': pickupLocation,
      };

      if (dropoffLocation != null && dropoffLocation.isNotEmpty) {
        queryParams['dropoff_location'] = dropoffLocation;
      }
      if (pickupTime != null && pickupTime.isNotEmpty) {
        queryParams['pickup_time'] = pickupTime;
      }
      if (dropoffTime != null && dropoffTime.isNotEmpty) {
        queryParams['dropoff_time'] = dropoffTime;
      }
      if (carType != null && carType.isNotEmpty) {
        queryParams['car_type'] = carType;
      }
      if (passengers != null) {
        queryParams['passengers'] = passengers.toString();
      }

      final response = await _apiClient.get(
        Endpoints.searchCars,
        queryParameters: queryParams,
      );

      if (response.data['success'] == true) {
        return response.data['data'] ?? [];
      }
      throw Exception('Failed to search cars');
    } catch (e) {
      print('[CAR_SERVICE] Error searching cars: $e');
      throw Exception('Error searching cars: $e');
    }
  }
}
