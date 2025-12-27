import 'package:get/get.dart';
import '../model/account_ui_model.dart';
import '../service/account_api.dart';

class AccountController extends GetxController {
  final AccountApi _accountApi = Get.find();

  final profile = Rx<ProfileModel?>(null);
  final isLoading = false.obs;

  @override
  void onInit() {
    super.onInit();
    fetchProfile();
  }

  Future<void> fetchProfile() async {
    try {
      isLoading.value = true;
      final result = await _accountApi.getProfile();
      profile.value = result;
    } catch (e) {
      // Handle error
    } finally {
      isLoading.value = false;
    }
  }
}
