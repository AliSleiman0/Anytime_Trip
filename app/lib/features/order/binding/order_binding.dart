import 'package:get/get.dart';
import '../controller/order_controller.dart';
import '../service/order_api.dart';

class OrderBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => OrderApi());
    Get.lazyPut(() => OrderController());
  }
}
