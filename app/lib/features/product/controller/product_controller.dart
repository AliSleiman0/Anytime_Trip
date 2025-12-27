import 'package:get/get.dart';
import '../../../core/utils/helpers.dart';
import '../model/product_ui_model.dart';
import '../service/product_api.dart';

class ProductController extends GetxController {
  final ProductApi _productApi = Get.find();

  final products = <ProductModel>[].obs;
  final isLoading = false.obs;
  final selectedProduct = Rx<ProductModel?>(null);

  @override
  void onInit() {
    super.onInit();
    fetchProducts();
  }

  Future<void> fetchProducts() async {
    try {
      isLoading.value = true;
      final result = await _productApi.getProducts();
      products.value = result;
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> fetchProductDetails(String id) async {
    try {
      isLoading.value = true;
      final result = await _productApi.getProductDetails(id);
      selectedProduct.value = result;
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    } finally {
      isLoading.value = false;
    }
  }
}
