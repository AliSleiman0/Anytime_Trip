import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../controller/order_controller.dart';
import '../../chatbot/view/chatbot_button.dart';

class OrderPage extends GetView<OrderController> {
  const OrderPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Orders'),
      ),
      body: Stack(
        children: [
          Obx(() {
            if (controller.isLoading.value) {
              return const Center(child: CircularProgressIndicator());
            }

            if (controller.orders.isEmpty) {
              return const Center(child: Text('No orders yet'));
            }

            return ListView.builder(
              itemCount: controller.orders.length,
              itemBuilder: (context, index) {
                final order = controller.orders[index];
                return ListTile(
                  title: Text('Order #${order.id}'),
                  subtitle: Text(order.status),
                  trailing: Text('\$${order.total.toStringAsFixed(2)}'),
                );
              },
            );
          }),
          const ChatbotButton(),
        ],
      ),
    );
  }
}