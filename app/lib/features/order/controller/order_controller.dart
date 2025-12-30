import 'package:get/get.dart';
import '../model/order_ui_model.dart';
import '../service/order_api.dart';

class OrderController extends GetxController {
  final OrderApi _orderApi = Get.find();

  final orders = <OrderModel>[].obs;
  final isLoading = false.obs;

  @override
  void onInit() {
    super.onInit();
    fetchOrders();
  }

  Future<void> fetchOrders() async {
    try {
      isLoading.value = true;
      final result = await _orderApi.getOrders();
      orders.value = result;
    } catch (e) {
      // Handle error
    } finally {
      isLoading.value = false;
    }
  }
}
