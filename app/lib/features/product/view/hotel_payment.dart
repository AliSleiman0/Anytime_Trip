import 'package:flutter/material.dart';

class HotelPayment extends StatefulWidget {
  final Map<String, dynamic> hotel;
  final Map<String, dynamic> room;
  final String checkinDate;
  final String checkoutDate;
  final int adults;
  final int children;
  final int rooms;

  const HotelPayment({
    super.key,
    required this.hotel,
    required this.room,
    required this.checkinDate,
    required this.checkoutDate,
    required this.adults,
    required this.children,
    required this.rooms,
  });

  @override
  State<HotelPayment> createState() => _HotelPaymentState();
}

class _HotelPaymentState extends State<HotelPayment> {
  String _selectedPaymentMethod = 'Whish Money';
  String _selectedTraveler = '';
  String _selectedCountryCode = '';
  String _selectedExpiryMonth = 'Month';
  String _selectedExpiryYear = 'Year';
  String _selectedBillingCountry = '';
  String _selectedState = '';
  
  final TextEditingController _firstNameController = TextEditingController();
  final TextEditingController _lastNameController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  final TextEditingController _nameOnCardController = TextEditingController();
  final TextEditingController _cardNumberController = TextEditingController();
  final TextEditingController _securityCodeController = TextEditingController();
  final TextEditingController _billingAddress1Controller = TextEditingController();
  final TextEditingController _billingAddress2Controller = TextEditingController();
  final TextEditingController _cityController = TextEditingController();
  final TextEditingController _zipCodeController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    final roomPrice = widget.room['price'].toDouble();
    final taxes = roomPrice * 0.17;
    final reservationFee = roomPrice * 0.0434;
    final taxOnFees = reservationFee * 0.164;
    final total = roomPrice + taxes + reservationFee + taxOnFees;

    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios, color: Colors.black, size: 20),
          onPressed: () => Navigator.pop(context),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Text(
              widget.hotel['name'],
              style: const TextStyle(
                color: Colors.black,
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
            ),
            Text(
              '${widget.checkinDate} - ${widget.checkoutDate}',
              style: const TextStyle(
                color: Colors.grey,
                fontSize: 11,
              ),
            ),
          ],
        ),
        centerTitle: true,
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Hotel/Room section
                  const Text(
                    'Room',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Colors.black87,
                    ),
                  ),
                  const SizedBox(height: 12),
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.grey.shade300),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              widget.checkinDate,
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Colors.black87,
                              ),
                            ),
                            const SizedBox(height: 4),
                            const Text(
                              '1 night, 1 room',
                              style: TextStyle(
                                fontSize: 12,
                                color: Colors.black54,
                              ),
                            ),
                          ],
                        ),
                        Text(
                          widget.checkoutDate,
                          style: const TextStyle(
                            fontSize: 12,
                            color: Colors.black54,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // Price Summary
                  const Text(
                    'Price Summary',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildPriceLine('1 night', '\$${roomPrice.toStringAsFixed(2)}'),
                  const Divider(height: 20),
                  _buildPriceLine('Taxes', '\$${taxes.toStringAsFixed(2)}'),
                  const Divider(height: 20),
                  _buildPriceLine('Reservation Fee', '\$${reservationFee.toStringAsFixed(2)}'),
                  const Divider(height: 20),
                  _buildPriceLine('Tax on fees', '\$${taxOnFees.toStringAsFixed(2)}'),
                  const Divider(height: 20),
                  _buildPriceLine('Total', '\$${total.toStringAsFixed(2)}', isTotal: true),
                  const SizedBox(height: 24),

                  // Who's checking in?
                  const Text(
                    'Who\'s checking in?',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildDropdown('Select Traveler', _selectedTraveler, (val) {
                    setState(() {
                      _selectedTraveler = val ?? '';
                    });
                  }),
                  const SizedBox(height: 12),
                  _buildTextField('First Name', _firstNameController),
                  const SizedBox(height: 12),
                  _buildTextField('Last Name', _lastNameController),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildDropdown('Country/ Territory Code', _selectedCountryCode, (val) {
                          setState(() {
                            _selectedCountryCode = val ?? '';
                          });
                        }),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildTextField('Phone Number', _phoneController),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),

                  // Pay With
                  const Text(
                    'Pay With',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildPaymentOption('Whish Money', 'whish.png'),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildPaymentOption('OMT Pay', 'omt.png'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildPaymentOption('Crypto Pay', 'crypto.png'),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildPaymentOption('Visa', 'visa.png'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),

                  // Card Details
                  const Text(
                    'CREDIT CARD INFO',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: Colors.black54,
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildTextField('Name on Card', _nameOnCardController),
                  const SizedBox(height: 12),
                  _buildTextField('0000 0000 0000 0000', _cardNumberController, keyboardType: TextInputType.number),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        flex: 2,
                        child: Row(
                          children: [
                            Expanded(
                              child: _buildDropdown('Month', _selectedExpiryMonth, (val) {
                                setState(() {
                                  _selectedExpiryMonth = val ?? 'Month';
                                });
                              }, items: ['Month', '01', '02', '03', '04', '05', '06', '07', '08', '09', '10', '11', '12']),
                            ),
                            const Padding(
                              padding: EdgeInsets.symmetric(horizontal: 8),
                              child: Text('-', style: TextStyle(fontSize: 18)),
                            ),
                            Expanded(
                              child: _buildDropdown('Year', _selectedExpiryYear, (val) {
                                setState(() {
                                  _selectedExpiryYear = val ?? 'Year';
                                });
                              }, items: ['Year', '2024', '2025', '2026', '2027', '2028', '2029', '2030']),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildTextField('Security Code', _securityCodeController, keyboardType: TextInputType.number),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  _buildDropdown('Country/ Territory Code', _selectedBillingCountry, (val) {
                    setState(() {
                      _selectedBillingCountry = val ?? '';
                    });
                  }),
                  const SizedBox(height: 12),
                  _buildTextField('Billing Address 1', _billingAddress1Controller),
                  const SizedBox(height: 12),
                  _buildTextField('Billing Address 2', _billingAddress2Controller),
                  const SizedBox(height: 12),
                  _buildTextField('City', _cityController),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildDropdown('State', _selectedState, (val) {
                          setState(() {
                            _selectedState = val ?? '';
                          });
                        }),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildTextField('ZIP code', _zipCodeController, keyboardType: TextInputType.number),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),
                ],
              ),
            ),
          ),

          // Complete Booking button
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
                  _showBookingConfirmationDialog();
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1e5a8e),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(25),
                  ),
                ),
                child: const Text(
                  'Complete Booking',
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

  Widget _buildPriceLine(String label, String amount, {bool isTotal = false}) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 14,
            color: isTotal ? const Color(0xFFD32F2F) : Colors.black87,
            fontWeight: isTotal ? FontWeight.w600 : FontWeight.w400,
          ),
        ),
        Text(
          amount,
          style: TextStyle(
            fontSize: 14,
            color: isTotal ? const Color(0xFFD32F2F) : Colors.black87,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }

  Widget _buildTextField(String hint, TextEditingController controller, {TextInputType? keyboardType}) {
    return SizedBox(
      height: 45,
      child: TextField(
        controller: controller,
        keyboardType: keyboardType,
        decoration: InputDecoration(
          hintText: hint,
          hintStyle: const TextStyle(fontSize: 13, color: Colors.grey),
          contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(6),
            borderSide: BorderSide(color: Colors.grey.shade300),
          ),
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(6),
            borderSide: BorderSide(color: Colors.grey.shade300),
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(6),
            borderSide: const BorderSide(color: Color(0xFF1e5a8e)),
          ),
        ),
      ),
    );
  }

  Widget _buildDropdown(String hint, String value, Function(String?) onChanged, {List<String>? items}) {
    return Container(
      height: 45,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.shade300),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Row(
        children: [
          Expanded(
            child: Text(
              value.isEmpty ? hint : value,
              style: TextStyle(
                fontSize: 13,
                color: value.isEmpty ? Colors.grey : Colors.black87,
              ),
            ),
          ),
          const Icon(Icons.arrow_drop_down, color: Color(0xFF1e5a8e), size: 20),
        ],
      ),
    );
  }

  Widget _buildPaymentOption(String name, String iconPath) {
    final isSelected = _selectedPaymentMethod == name;
    return GestureDetector(
      onTap: () {
        setState(() {
          _selectedPaymentMethod = name;
        });
      },
      child: Container(
        height: 50,
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 10),
        decoration: BoxDecoration(
          border: Border.all(
            color: isSelected ? const Color(0xFF1e5a8e) : Colors.grey.shade300,
            width: isSelected ? 2 : 1,
          ),
          borderRadius: BorderRadius.circular(6),
          color: Colors.white,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.asset(
              'assets/images/$iconPath',
              width: 20,
              height: 20,
              errorBuilder: (context, error, stackTrace) {
                return const Icon(Icons.payment, size: 20, color: Colors.grey);
              },
            ),
            const SizedBox(width: 6),
            Flexible(
              child: Text(
                name,
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                  color: isSelected ? const Color(0xFF1e5a8e) : Colors.black87,
                ),
                overflow: TextOverflow.ellipsis,
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showBookingConfirmationDialog() {
    final bookingReference = '1872306517801';
    
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          child: Container(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    GestureDetector(
                      onTap: () {
                        Navigator.pop(context);
                        Navigator.pop(context);
                      },
                      child: Container(
                        decoration: const BoxDecoration(
                          color: Color(0xFFD32F2F),
                          shape: BoxShape.circle,
                        ),
                        padding: const EdgeInsets.all(6),
                        child: const Icon(
                          Icons.close,
                          color: Colors.white,
                          size: 18,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Container(
                  width: 100,
                  height: 100,
                  decoration: BoxDecoration(
                    color: const Color(0xFF1e5a8e).withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Icon(
                    Icons.apartment,
                    size: 60,
                    color: Color(0xFF1e5a8e),
                  ),
                ),
                const SizedBox(height: 24),
                const Text(
                  'Your booking has been completed successfully.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.black87,
                    height: 1.5,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Your booking reference is:',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 12),
                Text(
                  bookingReference,
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w700,
                    color: Colors.black,
                    letterSpacing: 1.2,
                  ),
                ),
                const SizedBox(height: 16),
                const Text(
                  'You will receive confirmation email including your itinerary details',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.black87,
                    height: 1.5,
                  ),
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 50,
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.pop(context);
                      Navigator.pop(context);
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF1e5a8e),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(25),
                      ),
                    ),
                    child: const Text(
                      'Done',
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
          ),
        );
      },
    );
  }

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _phoneController.dispose();
    _nameOnCardController.dispose();
    _cardNumberController.dispose();
    _securityCodeController.dispose();
    _billingAddress1Controller.dispose();
    _billingAddress2Controller.dispose();
    _cityController.dispose();
    _zipCodeController.dispose();
    super.dispose();
  }
}
