import 'package:flutter/material.dart';
import '../../../core/utils/helpers.dart';

class ProductPriceWidget extends StatelessWidget {
  final double price;
  final double? originalPrice;

  const ProductPriceWidget({
    super.key,
    required this.price,
    this.originalPrice,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(
          Helpers.formatPrice(price),
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: Theme.of(context).primaryColor,
          ),
        ),
        if (originalPrice != null) ...[
          const SizedBox(width: 8),
          Text(
            Helpers.formatPrice(originalPrice!),
            style: const TextStyle(
              fontSize: 14,
              decoration: TextDecoration.lineThrough,
              color: Colors.grey,
            ),
          ),
        ],
      ],
    );
  }
}
