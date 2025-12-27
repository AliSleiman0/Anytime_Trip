import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../core/utils/helpers.dart';
import '../../../core/widgets/custom_button.dart';
import '../controller/product_controller.dart';
import '../widget/product_price_widget.dart';

class ProductDetailsPage extends GetView<ProductController> {
  const ProductDetailsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final productId = Get.parameters['id'] ?? '';
    
    if (productId.isNotEmpty) {
      controller.fetchProductDetails(productId);
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Product Details'),
      ),
      body: Obx(() {
        if (controller.isLoading.value) {
          return const Center(child: CircularProgressIndicator());
        }

        final product = controller.selectedProduct.value;
        if (product == null) {
          return const Center(child: Text('Product not found'));
        }

        return SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                height: 300,
                width: double.infinity,
                color: Colors.grey[200],
                child: product.image.isNotEmpty
                    ? Image.network(product.image, fit: BoxFit.cover)
                    : const Icon(Icons.image, size: 100),
              ),
              Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      product.name,
                      style: const TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    ProductPriceWidget(price: product.price),
                    const SizedBox(height: 16),
                    Chip(label: Text(product.category)),
                    const SizedBox(height: 16),
                    const Text(
                      'Description',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(product.description),
                    const SizedBox(height: 24),
                    CustomButton(
                      text: 'Add to Cart',
                      onPressed: () {
                        Helpers.showSnackbar(
                          'Success',
                          'Product added to cart',
                        );
                      },
                    ),
                  ],
                ),
              ),
            ],
          ),
        );
      }),
    );
  }
}
