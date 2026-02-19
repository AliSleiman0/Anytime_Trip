import '../network/api_client.dart';
import '../network/endpoints.dart';

class BookingService {
  final ApiClient _apiClient = ApiClient();

  // Get all bookings
  Future<Map<String, dynamic>> getAllBookings() async {
    try {
      final response = await _apiClient.get(Endpoints.myBookings);
      if (response.data['success'] == true) {
        return response.data['data'];
      }
      throw Exception('Failed to fetch bookings');
    } catch (e) {
      throw Exception('Error fetching bookings: $e');
    }
  }

  // Get car bookings
  Future<List<dynamic>> getCarBookings() async {
    try {
      final response = await _apiClient.get(Endpoints.myCarBookings);
      if (response.data['success'] == true) {
        return response.data['data'] ?? [];
      }
      throw Exception('Failed to fetch car bookings');
    } catch (e) {
      throw Exception('Error fetching car bookings: $e');
    }
  }

  // Get flight bookings
  Future<List<dynamic>> getFlightBookings() async {
    try {
      final response = await _apiClient.get(Endpoints.myFlightBookings);
      if (response.data['success'] == true) {
        return response.data['data'] ?? [];
      }
      throw Exception('Failed to fetch flight bookings');
    } catch (e) {
      throw Exception('Error fetching flight bookings: $e');
    }
  }

  // Get hotel bookings
  Future<List<dynamic>> getHotelBookings() async {
    try {
      final response = await _apiClient.get(Endpoints.myHotelBookings);
      if (response.data['success'] == true) {
        return response.data['data'] ?? [];
      }
      throw Exception('Failed to fetch hotel bookings');
    } catch (e) {
      throw Exception('Error fetching hotel bookings: $e');
    }
  }

  // Get transfer bookings
  Future<List<dynamic>> getTransferBookings() async {
    try {
      final response = await _apiClient.get(Endpoints.myTransferBookings);
      if (response.data['success'] == true) {
        return response.data['data'] ?? [];
      }
      throw Exception('Failed to fetch transfer bookings');
    } catch (e) {
      throw Exception('Error fetching transfer bookings: $e');
    }
  }

  // Get specific booking by ID and type
  Future<Map<String, dynamic>> getBookingById(String type, String id) async {
    try {
      final response = await _apiClient.get(Endpoints.bookingDetails(type, id));
      if (response.data['success'] == true) {
        return response.data['data'];
      }
      throw Exception('Failed to fetch booking');
    } catch (e) {
      throw Exception('Error fetching booking: $e');
    }
  }

  // Cancel booking
  Future<Map<String, dynamic>> cancelBooking(String type, String id) async {
    try {
      final response = await _apiClient.post(Endpoints.cancelBooking(type, id));
      return {
        'success': response.data['success'] == true,
        'message': response.data['message'] ?? 'Booking cancelled successfully',
        'data': response.data['data'],
      };
    } catch (e) {
      return {
        'success': false,
        'message': 'Error cancelling booking: $e',
      };
    }
  }
}
