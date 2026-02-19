import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../model/popular_location_model.dart';

class PopularLocationService {
  final ApiClient _apiClient = ApiClient();

  Future<List<PopularLocationModel>> getPopularLocations() async {
    try {
      final response = await _apiClient.get(Endpoints.popularLocations);
      final List<dynamic> data = response.data['data'] ?? [];
      return data.map((location) => PopularLocationModel.fromJson(location as Map<String, dynamic>)).toList();
    } catch (e) {
      print('[POPULAR_LOCATION_SERVICE_ERROR] Failed to get popular locations: $e');
      // Return empty list instead of throwing to prevent app crashes
      return [];
    }
  }
}
