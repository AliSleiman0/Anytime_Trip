import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../model/banner_model.dart';

class BannerService {
  final ApiClient _apiClient = ApiClient();

  Future<List<BannerModel>> getBanners() async {
    try {
      final response = await _apiClient.get(Endpoints.banners);
      final List<dynamic> data = response.data['data'] ?? [];
      return data.map((banner) => BannerModel.fromJson(banner as Map<String, dynamic>)).toList();
    } catch (e) {
      print('[BANNER_SERVICE_ERROR] Failed to get banners: $e');
      // Return empty list instead of throwing to prevent app crashes
      return [];
    }
  }

  Future<List<BannerModel>> getHomepageBanners() async {
    try {
      final response = await _apiClient.get(Endpoints.homepageBanners);
      final List<dynamic> data = response.data['data'] ?? [];
      return data.map((banner) => BannerModel.fromJson(banner as Map<String, dynamic>)).toList();
    } catch (e) {
      print('[BANNER_SERVICE_ERROR] Failed to get homepage banners: $e');
      // Return empty list instead of throwing to prevent app crashes
      return [];
    }
  }

  Future<List<BannerModel>> getSearchBanners() async {
    try {
      final response = await _apiClient.get(Endpoints.searchBanners);
      final List<dynamic> data = response.data['data'] ?? [];
      return data.map((banner) => BannerModel.fromJson(banner as Map<String, dynamic>)).toList();
    } catch (e) {
      print('[BANNER_SERVICE_ERROR] Failed to get search banners: $e');
      // Return empty list instead of throwing to prevent app crashes
      return [];
    }
  }
}
