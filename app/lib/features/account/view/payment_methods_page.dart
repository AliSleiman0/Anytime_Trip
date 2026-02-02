import 'package:flutter/material.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/network/endpoints.dart';
import '../service/account_api.dart';
import '../model/account_ui_model.dart';

class PaymentMethodsPage extends StatefulWidget {
  const PaymentMethodsPage({super.key});

  @override
  State<PaymentMethodsPage> createState() => _PaymentMethodsPageState();
}

class _PaymentMethodsPageState extends State<PaymentMethodsPage> {
  late TextEditingController _cardHolderController;
  late TextEditingController _cardNumberController;
  late TextEditingController _expiryMonthController;
  late TextEditingController _expiryYearController;
  late TextEditingController _cvvController;
  late TextEditingController _countryController;
  late TextEditingController _cityController;
  late TextEditingController _stateController;
  late TextEditingController _addressController;
  late TextEditingController _postalCodeController;
  late TextEditingController _phoneController;
  
  final _api = AccountApi();
  List<PaymentMethod> _paymentMethods = [];
  bool _isLoading = true;
  String? _selectedCardBrand;
  bool _isDefault = false;

  @override
  void initState() {
    super.initState();
    _initControllers();
    _loadPaymentMethods();
  }

  void _initControllers() {
    _cardHolderController = TextEditingController();
    _cardNumberController = TextEditingController();
    _expiryMonthController = TextEditingController();
    _expiryYearController = TextEditingController();
    _cvvController = TextEditingController();
    _countryController = TextEditingController();
    _cityController = TextEditingController();
    _stateController = TextEditingController();
    _addressController = TextEditingController();
    _postalCodeController = TextEditingController();
    _phoneController = TextEditingController();
  }

  Future<void> _loadPaymentMethods() async {
    try {
      setState(() => _isLoading = true);
      final methods = await _api.getPaymentMethods();
      setState(() {
        _paymentMethods = methods;
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load payment methods: $e')),
        );
      }
    }
  }

  Future<void> _addPaymentMethod() async {
    if (_validateForm()) {
      try {
        final expiryMonth = int.tryParse(_expiryMonthController.text) ?? 0;
        final expiryYear = int.tryParse(_expiryYearController.text) ?? 0;
        
        final request = PaymentMethodRequest(
          cardHolderName: _cardHolderController.text,
          cardNumber: _cardNumberController.text,
          cardBrand: _selectedCardBrand ?? 'Visa',
          expiryMonth: expiryMonth,
          expiryYear: expiryYear,
          cvv: _cvvController.text,
          billingAddress: _addressController.text,
          city: _cityController.text,
          state: _stateController.text,
          postalCode: _postalCodeController.text,
          country: _countryController.text,
          phoneNumber: _phoneController.text,
          isDefault: _isDefault,
        );

        await _api.addPaymentMethod(request);
        
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Payment method saved successfully')),
          );
          _clearForm();
          _loadPaymentMethods();
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to save payment method: $e')),
          );
        }
      }
    }
  }

  bool _validateForm() {
    if (_cardHolderController.text.isEmpty ||
        _cardNumberController.text.isEmpty ||
        _expiryMonthController.text.isEmpty ||
        _expiryYearController.text.isEmpty ||
        _cvvController.text.isEmpty ||
        _addressController.text.isEmpty ||
        _cityController.text.isEmpty ||
        _stateController.text.isEmpty ||
        _postalCodeController.text.isEmpty ||
        _countryController.text.isEmpty ||
        _phoneController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please fill all required fields')),
      );
      return false;
    }
    return true;
  }

  void _clearForm() {
    _cardHolderController.clear();
    _cardNumberController.clear();
    _expiryMonthController.clear();
    _expiryYearController.clear();
    _cvvController.clear();
    _countryController.clear();
    _cityController.clear();
    _stateController.clear();
    _addressController.clear();
    _postalCodeController.clear();
    _phoneController.clear();
    _selectedCardBrand = 'Visa';
    _isDefault = false;
  }

  Future<void> _deletePaymentMethod(String id) async {
    try {
      await _api.deletePaymentMethod(id);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Payment method deleted')),
        );
        _loadPaymentMethods();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to delete: $e')),
        );
      }
    }
  }

  Future<void> _setDefault(String id) async {
    try {
      await _api.setDefaultPaymentMethod(id);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Default payment method updated')),
        );
        _loadPaymentMethods();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to update: $e')),
        );
      }
    }
  }

  @override
  void dispose() {
    _cardHolderController.dispose();
    _cardNumberController.dispose();
    _expiryMonthController.dispose();
    _expiryYearController.dispose();
    _cvvController.dispose();
    _countryController.dispose();
    _cityController.dispose();
    _stateController.dispose();
    _addressController.dispose();
    _postalCodeController.dispose();
    _phoneController.dispose();
    super.dispose();
  }

  String _getUserName() {
    final user = StorageService().getUser();
    return user?['name'] ?? 'User';
  }

  String? _getProfileImage() {
    final user = StorageService().getUser();
    return user?['profile_image'];
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
            Text(
              'Hello, ${_getUserName()}',
              style: const TextStyle(
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
                child: _getProfileImage() != null && _getProfileImage()!.isNotEmpty
                    ? Image.network(
                        '${Endpoints.baseUrl.replaceAll('/api', '')}${_getProfileImage()}',
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) {
                          return Image.asset(
                            'assets/images/defaultp.png',
                            fit: BoxFit.cover,
                            errorBuilder: (context, error, stackTrace) {
                              return const Icon(Icons.person, size: 20, color: Color(0xFF1e5a8e));
                            },
                          );
                        },
                      )
                    : Image.asset(
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
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const SizedBox(height: 16),
                    // Saved Payment Methods Section
                    const Text(
                      'Saved Payment Methods',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: Colors.black87,
                      ),
                    ),
                    const SizedBox(height: 12),
                    if (_paymentMethods.isEmpty)
                      Container(
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          border: Border.all(color: Colors.grey[300]!),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: const Text(
                          'No payment methods saved yet',
                          style: TextStyle(color: Colors.grey),
                        ),
                      )
                    else
                      ListView.builder(
                        shrinkWrap: true,
                        physics: const NeverScrollableScrollPhysics(),
                        itemCount: _paymentMethods.length,
                        itemBuilder: (context, index) {
                          final pm = _paymentMethods[index];
                          return _buildPaymentMethodCard(pm);
                        },
                      ),
                    const SizedBox(height: 24),
                    // Add New Payment Method Section
                    const Text(
                      'Add New Payment Method',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: Colors.black87,
                      ),
                    ),
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
                    _buildCardBrandDropdown(),
                    const SizedBox(height: 16),
                    // Expiry Date and CVV Row
                    Row(
                      children: [
                        Expanded(
                          flex: 1,
                          child: _buildLabeledTextField(
                            label: 'Expiry Month (MM)',
                            controller: _expiryMonthController,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          flex: 1,
                          child: _buildLabeledTextField(
                            label: 'Expiry Year (YYYY)',
                            controller: _expiryYearController,
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
                    _buildLabeledTextField(
                      label: 'Street Address',
                      controller: _addressController,
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
                    // State and Postal Code Row
                    Row(
                      children: [
                        Expanded(
                          child: _buildLabeledTextField(
                            label: 'State',
                            controller: _stateController,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _buildLabeledTextField(
                            label: 'Postal Code',
                            controller: _postalCodeController,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    _buildLabeledTextField(
                      label: 'Phone Number',
                      controller: _phoneController,
                    ),
                    const SizedBox(height: 16),
                    // Set as Default Checkbox
                    CheckboxListTile(
                      value: _isDefault,
                      onChanged: (value) {
                        setState(() => _isDefault = value ?? false);
                      },
                      title: const Text('Set as default payment method'),
                      contentPadding: EdgeInsets.zero,
                    ),
                    const SizedBox(height: 32),
                    // Save Changes Button
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: _addPaymentMethod,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF336891),
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(25),
                          ),
                          elevation: 2,
                        ),
                        child: const Text(
                          'Save Payment Method',
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

  Widget _buildPaymentMethodCard(PaymentMethod pm) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      pm.cardHolderName,
                      style: const TextStyle(
                        fontWeight: FontWeight.w600,
                        fontSize: 14,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '${pm.cardBrand} •••• ${pm.cardNumberLast4.isNotEmpty ? pm.cardNumberLast4 : 'XXXX'}',
                      style: const TextStyle(
                        color: Colors.grey,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
                if (pm.isDefault)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.green[100],
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: const Text(
                      'Default',
                      style: TextStyle(
                        color: Colors.green,
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              'Expires ${pm.expiryMonth.toString().padLeft(2, '0')}/${pm.expiryYear}',
              style: const TextStyle(fontSize: 12, color: Colors.grey),
            ),
            const SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                if (!pm.isDefault)
                  TextButton(
                    onPressed: () => _setDefault(pm.id),
                    child: const Text('Set Default'),
                  ),
                TextButton(
                  onPressed: () => _deletePaymentMethod(pm.id),
                  style: TextButton.styleFrom(foregroundColor: Colors.red),
                  child: const Text('Delete'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCardBrandDropdown() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Card Brand',
          style: TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            color: Color(0xFFD32F2F),
          ),
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            border: Border.all(color: const Color(0xFF336891), width: 2),
            borderRadius: BorderRadius.circular(8),
          ),
          child: DropdownButton<String>(
            value: _selectedCardBrand ?? 'Visa',
            isExpanded: true,
            underline: const SizedBox(),
            onChanged: (String? newValue) {
              setState(() => _selectedCardBrand = newValue);
            },
            items: ['Visa', 'Mastercard', 'Amex', 'Discover'].map((String value) {
              return DropdownMenuItem<String>(
                value: value,
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 14),
                  child: Text(value),
                ),
              );
            }).toList(),
          ),
        ),
      ],
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
