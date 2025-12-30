import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../model/order_ui_model.dart';

class OrderApi {
  final ApiClient _apiClient = ApiClient();

  Future<List<OrderModel>> getOrders() async {
    try {
      final response = await _apiClient.get(Endpoints.orders);
      final List<dynamic> data = response.data['orders'] ?? [];
      return data.map((json) => OrderModel.fromJson(json)).toList();
    } catch (e) {
      rethrow;
    }
  }
}
