import 'package:get/get.dart';
import '../controller/auth_controller.dart';
import '../service/auth_api.dart';

class AuthBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => AuthApi());
    Get.lazyPut(() => AuthController());
  }
}
