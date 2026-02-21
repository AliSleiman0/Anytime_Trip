import '../network/api_client.dart';
import '../network/endpoints.dart';

class CarService {
  final ApiClient _apiClient = ApiClient();

  // Fetch available pickup and dropoff locations from the database
  Future<Map<String, List<String>>> getAvailableLocations() async {
    try {
      final response = await _apiClient.get(
        Endpoints.availableCarLocations,
      );

      if (response.data['success'] == true) {
        final data = response.data['data'] as Map<String, dynamic>;
        return {
          'pickup_locations': List<String>.from(data['pickup_locations'] ?? []),
          'dropoff_locations': List<String>.from(data['dropoff_locations'] ?? []),
        };
      }
      throw Exception('Failed to fetch available locations');
    } catch (e) {
      print('[CAR_SERVICE] Error fetching available locations: $e');
      throw Exception('Error fetching available locations: $e');
    }
  }

  // Search for available cars with date-based availability checking
  Future<List<dynamic>> searchCars({
    required String pickupLocation,
    String? dropoffLocation,
    DateTime? pickupDate,
    DateTime? dropoffDate,
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
      
      // Add date parameters
      if (pickupDate != null) {
        queryParams['pickup_date'] = pickupDate.toIso8601String();
      }
      if (dropoffDate != null) {
        queryParams['dropoff_date'] = dropoffDate.toIso8601String();
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

  // Create a car booking
  Future<Map<String, dynamic>> createCarBooking({
    required String carId,
    required String pickupLocation,
    required String dropoffLocation,
    required DateTime pickupDate,
    required DateTime dropoffDate,
    required String pickupTime,
    required String dropoffTime,
    required String driverName,
    required String phone,
    required String countryCode,
    required String email,
    required String paymentMethod,
    required double totalPrice,
  }) async {
    try {
      // Convert dates to RFC3339 format with timezone
      final pickupDateStr = pickupDate.toUtc().toIso8601String().replaceAll('Z', '+00:00');
      final dropoffDateStr = dropoffDate.toUtc().toIso8601String().replaceAll('Z', '+00:00');

      final bookingData = {
        'car_id': carId,
        'car_type': 'Rental',
        'passengers': 5,
        'status': 'confirmed',
        'customer': {
          'name': driverName,
          'email': email,
        },
        'driver': {
          'name': driverName,
          'phone_number': '$countryCode$phone',
          'license_number': '',
        },
        'pickup': {
          'location': pickupLocation,
          'address': pickupLocation,
          'date': pickupDateStr,
          'time': pickupTime,
        },
        'dropoff': {
          'location': dropoffLocation,
          'address': dropoffLocation,
          'date': dropoffDateStr,
          'time': dropoffTime,
        },
        'pricing': {
          'rental_price': totalPrice,
          'taxes': 0.0,
          'total': totalPrice,
        },
        'amount': totalPrice,
        'currency': '\$',
        'payment_status': 'paid',
        'details': '$paymentMethod payment',
      };

      print('[CAR_SERVICE] Booking payload: $bookingData');

      final response = await _apiClient.post(
        Endpoints.createCarBooking,
        data: bookingData,
      );

      print('[CAR_SERVICE] Booking response: ${response.data}');

      // Check if booking was created successfully
      // Backend returns 201 with booking data, not a success field
      if (response.statusCode == 201 && response.data['booking'] != null) {
        return response.data['booking'] ?? {};
      }
      
      if (response.data['success'] == true) {
        return response.data['data'] ?? {};
      }
      
      throw Exception('Failed to create car booking: ${response.data}');
    } catch (e) {
      print('[CAR_SERVICE] Error creating car booking: $e');
      throw Exception('Error creating car booking: $e');
    }
  }
}
