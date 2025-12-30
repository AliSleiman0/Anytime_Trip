import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../model/product_ui_model.dart';

class ProductApi {
  final ApiClient _apiClient = ApiClient();

  Future<List<ProductModel>> getProducts() async {
    try {
      final response = await _apiClient.get(Endpoints.products);
      final List<dynamic> data = response.data['products'] ?? [];
      return data.map((json) => ProductModel.fromJson(json)).toList();
    } catch (e) {
      rethrow;
    }
  }

  Future<ProductModel> getProductDetails(String id) async {
    try {
      final response = await _apiClient.get(Endpoints.productDetails(id));
      return ProductModel.fromJson(response.data);
    } catch (e) {
      rethrow;
    }
  }
}
