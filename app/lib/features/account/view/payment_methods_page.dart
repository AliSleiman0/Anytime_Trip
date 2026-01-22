import 'package:flutter/material.dart';

class PaymentMethodsPage extends StatefulWidget {
  const PaymentMethodsPage({super.key});

  @override
  State<PaymentMethodsPage> createState() => _PaymentMethodsPageState();
}

class _PaymentMethodsPageState extends State<PaymentMethodsPage> {
  late TextEditingController _cardHolderController;
  late TextEditingController _cardNumberController;
  late TextEditingController _expiryDateController;
  late TextEditingController _cvvController;
  late TextEditingController _countryController;
  late TextEditingController _cityController;
  late TextEditingController _streetAddress1Controller;
  late TextEditingController _streetAddress2Controller;
  late TextEditingController _buildingController;
  late TextEditingController _floorController;

  @override
  void initState() {
    super.initState();
    _cardHolderController = TextEditingController(text: 'Michael John Doe');
    _cardNumberController = TextEditingController(text: '**** **** **** 0281');
    _expiryDateController = TextEditingController(text: '12/27');
    _cvvController = TextEditingController(text: '547');
    _countryController = TextEditingController(text: 'Lebanon');
    _cityController = TextEditingController(text: 'Beirut');
    _streetAddress1Controller = TextEditingController(text: 'Saint Therese Main R...');
    _streetAddress2Controller = TextEditingController(text: '');
    _buildingController = TextEditingController(text: 'John Doe Building');
    _floorController = TextEditingController(text: '3');
  }

  @override
  void dispose() {
    _cardHolderController.dispose();
    _cardNumberController.dispose();
    _expiryDateController.dispose();
    _cvvController.dispose();
    _countryController.dispose();
    _cityController.dispose();
    _streetAddress1Controller.dispose();
    _streetAddress2Controller.dispose();
    _buildingController.dispose();
    _floorController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF1e5a8e), size: 24),
          onPressed: () => Navigator.pop(context),
        ),
        title: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text(
              'Hello, User123',
              style: TextStyle(
                color: Colors.black87,
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(width: 12),
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: const Color(0xFFE0E0E0),
              ),
              child: ClipOval(
                child: Image.asset(
                  'assets/images/defaultp.png',
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return const Icon(Icons.person, size: 20, color: Color(0xFF1e5a8e));
                  },
                ),
              ),
            ),
          ],
        ),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 16),

              // Card Information Section
              _buildLabeledTextField(
                label: 'Card Holder Name',
                controller: _cardHolderController,
              ),
              const SizedBox(height: 16),

              _buildLabeledTextField(
                label: 'Card Number',
                controller: _cardNumberController,
              ),
              const SizedBox(height: 16),

              // Expiry Date and CVV Row
              Row(
                children: [
                  Expanded(
                    flex: 1,
                    child: _buildLabeledTextField(
                      label: 'Expiry Date',
                      controller: _expiryDateController,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    flex: 1,
                    child: _buildLabeledTextField(
                      label: 'CVV',
                      controller: _cvvController,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),

              // Billing Address Section
              const Text(
                'Billing Address',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: Colors.black87,
                ),
              ),
              const SizedBox(height: 16),

              // Country and City Row
              Row(
                children: [
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'Country',
                      controller: _countryController,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'City',
                      controller: _cityController,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // Street Address Row
              Row(
                children: [
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'Street Address 1',
                      controller: _streetAddress1Controller,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'Street Address 2',
                      controller: _streetAddress2Controller,
                      isOptional: true,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // Building and Floor Row
              Row(
                children: [
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'Building',
                      controller: _buildingController,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'Floor',
                      controller: _floorController,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 32),

              // Save Changes Button
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _showSaveConfirmationDialog,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF336891),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(25),
                    ),
                    elevation: 2,
                  ),
                  child: const Text(
                    'Save Changes',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }

  void _showSaveConfirmationDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text(
                  'Save Payment Information?',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 16),
                const Text(
                  'Save this card to your account for one-click payments. Your information is securely stored and protected by encryption.',
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w400,
                    color: Colors.black54,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 24),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () => Navigator.pop(context),
                        style: OutlinedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          side: const BorderSide(color: Color(0xFF336891), width: 2),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(20),
                          ),
                        ),
                        child: const Text(
                          'Cancel',
                          style: TextStyle(
                            color: Color(0xFF336891),
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: () {
                          Navigator.pop(context);
                          // Add save logic here
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF336891),
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(20),
                          ),
                        ),
                        child: const Text(
                          'Save',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildLabeledTextField({
    required String label,
    required TextEditingController controller,
    bool isOptional = false,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Text(
              label,
              style: const TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: Color(0xFFD32F2F),
              ),
            ),
            if (isOptional)
              const Text(
                ' (Optional)',
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  color: Colors.grey,
                ),
              ),
          ],
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            border: Border.all(color: const Color(0xFF336891), width: 2),
            borderRadius: BorderRadius.circular(8),
            color: Colors.white,
          ),
          child: TextField(
            controller: controller,
            decoration: const InputDecoration(
              border: InputBorder.none,
              contentPadding: EdgeInsets.symmetric(horizontal: 14, vertical: 14),
            ),
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: Colors.black87,
            ),
          ),
        ),
      ],
    );
  }
}
