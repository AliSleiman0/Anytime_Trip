import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'transfer_checkout.dart';
import '../../../core/widgets/unified_ui_components.dart';

class TransferDetails extends StatelessWidget {
  final Map<String, dynamic> transfer;
  final String pickupLocation;
  final String dropoffLocation;
  final String pickupTime;
  final String transferType;

  const TransferDetails({
    super.key,
    required this.transfer,
    required this.pickupLocation,
    required this.dropoffLocation,
    required this.pickupTime,
    required this.transferType,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: UnifiedAppBar(
        title: transfer['category'] ?? 'Transfer',
        onBackPressed: () => Navigator.pop(context),
      ),
      body: Column(
        children: [
          const Divider(height: 1, color: Colors.grey),
          
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Transfer Card
                  _buildTransferCard(),
                  const SizedBox(height: 24),

                  // Transfer Details Section
                  const Text(
                    'Transfer Details',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF1e5a8e),
                      decoration: TextDecoration.underline,
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildDetailRow('Provider Name:', 'John Morrison'),
                  _buildDetailRow('Price:', transfer['price'] ?? '\$50'),
                  _buildDetailRow('Transfer Type:', 'Private'),
                  _buildDetailRow('', 'Shuttle to the car located in the airport'),
                  const SizedBox(height: 24),

                  // Pickup & Drop-off Details Section
                  const Text(
                    'Pickup & Drop-off Details',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF1e5a8e),
                      decoration: TextDecoration.underline,
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildDetailRow('Pickup Location:', pickupLocation.isNotEmpty ? pickupLocation : 'BEY Airport'),
                  _buildDetailRow('Drop-off Location:', dropoffLocation.isNotEmpty ? dropoffLocation : 'BEY Airport'),
                  _buildDetailRow('Pick up time:', pickupTime.isNotEmpty ? pickupTime : '9:00 AM'),
                  _buildDetailRow('Return time:', '9:00 PM'),
                ],
              ),
            ),
          ),

          // Book Button
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.white,
              boxShadow: [
                BoxShadow(
                  color: Colors.grey.shade200,
                  blurRadius: 4,
                  offset: const Offset(0, -2),
                ),
              ],
            ),
            child: SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => TransferCheckout(
                        transfer: transfer,
                        pickupLocation: pickupLocation,
                        dropoffLocation: dropoffLocation,
                        pickupTime: pickupTime,
                        transferType: transferType,
                      ),
                    ),
                  );
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1e5a8e),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(25),
                  ),
                ),
                child: const Text(
                  'Book',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTransferCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFFE8F0F7),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: const Color(0xFF1e5a8e).withOpacity(0.3)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  transfer['category'] ?? 'Midsize SUV',
                  style: const TextStyle(
                    color: Color(0xFFD32F2F),
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  transfer['name'] ?? 'Toyota RAV 4 or similar',
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  transfer['passengers'] ?? '5 Passengers',
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  transfer['estimatedTime'] ?? 'Estimated time: 40 min',
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  transfer['meetGreet'] ?? 'Meet & Greet availability',
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 12),
          Column(
            children: [
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: Image.asset(
                  'assets/images/urus.jpeg',
                  width: 80,
                  height: 60,
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      width: 80,
                      height: 60,
                      color: Colors.grey.shade300,
                      child: const Icon(Icons.directions_car, size: 30),
                    );
                  },
                ),
              ),
              const SizedBox(height: 8),
              Text(
                transfer['price'] ?? '\$50',
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w700,
                  color: Colors.black,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    if (label.isEmpty) {
      return Padding(
        padding: const EdgeInsets.only(bottom: 8),
        child: Text(
          value,
          style: const TextStyle(
            fontSize: 13,
            color: Colors.black87,
          ),
        ),
      );
    }
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: const TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w500,
              color: Colors.black87,
            ),
          ),
          const SizedBox(width: 4),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontSize: 13,
                color: Colors.black87,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
