import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../model/account_ui_model.dart';

class AccountApi {
  final ApiClient _apiClient = ApiClient();

  Future<ProfileModel> getProfile() async {
    try {
      final response = await _apiClient.get(Endpoints.profile);
      return ProfileModel.fromJson(response.data);
    } catch (e) {
      rethrow;
    }
  }
}
