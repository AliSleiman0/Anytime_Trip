import 'package:get/get.dart';
import '../controller/account_controller.dart';
import '../service/account_api.dart';

class AccountBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => AccountApi());
    Get.lazyPut(() => AccountController());
  }
}
