import 'package:get/get.dart';
import '../controller/product_controller.dart';
import '../service/product_api.dart';

class ProductBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => ProductApi());
    Get.lazyPut(() => ProductController());
  }
}
