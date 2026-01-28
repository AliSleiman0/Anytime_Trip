import 'package:get/get.dart';
import '../../features/auth/binding/auth_binding.dart';
import '../../features/auth/view/login_page.dart';
import '../../features/auth/view/register_page.dart';
import '../../features/auth/view/forgot_password_page.dart';
import '../../features/auth/view/otp_confirmation_page.dart';
import '../../features/auth/view/set_new_password_page.dart';
import '../../features/auth/view/splash_page.dart';
import '../../features/product/binding/product_binding.dart';
import '../../features/product/view/home_page.dart';
import '../../features/product/view/product_page.dart';
import '../../features/product/view/product_details_page.dart';
import '../../features/account/binding/account_binding.dart';
import '../../features/account/view/account_page.dart';
import '../../features/order/binding/order_binding.dart';
import '../../features/order/view/order_page.dart';
import '../../features/order/view/order_details_page.dart';
import 'app_routes.dart';

class AppPages {
  static final pages = [
    GetPage(
      name: AppRoutes.SPLASH_SCREEN,
      page: () => const SplashPage(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: AppRoutes.SPLASH,
      page: () => const SplashPage(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: AppRoutes.LOGIN,
      page: () => const LoginPage(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: AppRoutes.REGISTER,
      page: () => const RegisterPage(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: AppRoutes.FORGOT_PASSWORD,
      page: () => const ForgotPasswordPage(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: AppRoutes.OTP_CONFIRMATION,
      page: () => const OtpConfirmationPage(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: AppRoutes.SET_NEW_PASSWORD,
      page: () => const SetNewPasswordPage(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: AppRoutes.HOME,
      page: () => const HomePage(),
      binding: ProductBinding(),
    ),
    GetPage(
      name: AppRoutes.PRODUCT,
      page: () => const ProductPage(),
      binding: ProductBinding(),
    ),
    GetPage(
      name: AppRoutes.PRODUCT_DETAILS,
      page: () => const ProductDetailsPage(),
      binding: ProductBinding(),
    ),
    GetPage(
      name: AppRoutes.ACCOUNT,
      page: () => const AccountPage(),
      binding: AccountBinding(),
    ),
    GetPage(
      name: AppRoutes.ORDER,
      page: () => const OrderPage(),
      binding: OrderBinding(),
    ),
    GetPage(
      name: AppRoutes.ORDER_DETAILS,
      page: () => const OrderDetailsPage(),
      binding: OrderBinding(),
    ),
  ];
}
