import 'package:flutter/material.dart';
import '../../../core/widgets/unified_ui_components.dart';
import '../../../core/widgets/policy_dialogs.dart';
import '../../account/view/bookings_page.dart';
import 'home_page.dart';

class FlightCheckoutScreen extends StatefulWidget {
  final String from;
  final String to;
  final String departureTime;
  final String arrivalTime;
  final String airline;
  final String stops;
  final String tripType;
  final int price;
  final String departureDate;
  final String? returnDate;
  final String? returnDepartureTime;
  final String? returnArrivalTime;
  final String? returnAirline;
  final String? returnStops;

  const FlightCheckoutScreen({
    Key? key,
    required this.from,
    required this.to,
    required this.departureTime,
    required this.arrivalTime,
    required this.airline,
    required this.stops,
    required this.tripType,
    required this.price,
    required this.departureDate,
    this.returnDate,
    this.returnDepartureTime,
    this.returnArrivalTime,
    this.returnAirline,
    this.returnStops,
  }) : super(key: key);

  @override
  State<FlightCheckoutScreen> createState() => _FlightCheckoutScreenState();
}

class _FlightCheckoutScreenState extends State<FlightCheckoutScreen> {
  int _selectedPayment = 0;
  bool _showDetails = true;
  final TextEditingController _firstNameController = TextEditingController();
  final TextEditingController _lastNameController = TextEditingController();
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F5F5),
      appBar: UnifiedAppBar(
        title: 'Checkout',
        onBackPressed: () => Navigator.pop(context),
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            // Header with Destination and Dates
            Container(
              color: Colors.white,
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Center(
                    child: Text(
                      widget.to,
                      style: const TextStyle(fontSize: 28, fontWeight: FontWeight.w700, color: Colors.black87),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Center(
                    child: Text(
                      'Mon, 26 Sep${widget.tripType == 'roundtrip' ? ' - Wed, 28 Sep' : ''}',
                      style: const TextStyle(fontSize: 13, color: Colors.grey),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Divider(color: Colors.grey.shade300, thickness: 1),
                  const SizedBox(height: 16),
                  
                  // Trip Type and Passenger Info
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        widget.tripType == 'roundtrip' ? 'Roundtrip Flight' : 'One Way Flight',
                        style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700, color: Color(0xFF1e5a8e)),
                      ),
                      Text(
                        '1 ticket: 1 adult',
                        style: const TextStyle(fontSize: 13, color: Colors.grey),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // Departure Flight
                  if (_showDetails) ...[
                    _buildFlightDetailRow('Departure', widget.from, widget.to, widget.departureTime, widget.arrivalTime, widget.airline, 'Mon, 26 Sep'),
                    const SizedBox(height: 16),
                  ],
                  
                  // Return Flight (if roundtrip)
                  if (_showDetails && widget.tripType == 'roundtrip')
                    Column(
                      children: [
                        _buildFlightDetailRow('Return', widget.to, widget.from, widget.returnDepartureTime ?? '', widget.returnArrivalTime ?? '', widget.returnAirline ?? '', 'Wed, 28 Sep'),
                        const SizedBox(height: 16),
                      ],
                    ),
                  
                  // Price Summary Section
                  if (_showDetails) ...[
                    _buildPriceSummarySection(),
                    const SizedBox(height: 16),
                  ],
                  
                  // Show Less Details Button
                  Center(
                    child: ElevatedButton(
                      onPressed: () => setState(() => _showDetails = !_showDetails),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF2c5f8d),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                        padding: const EdgeInsets.symmetric(horizontal: 40, vertical: 12),
                      ),
                      child: Text(
                        _showDetails ? 'Show Less Details' : 'Show More Details',
                        style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w600, color: Colors.white),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 8),

            // Passenger Information
            Container(
              color: Colors.white,
              padding: const EdgeInsets.all(16),
              margin: const EdgeInsets.symmetric(vertical: 8),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 12),
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.grey.shade300),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text('First Traveler Information (Adult 1)', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e))),
                        const SizedBox(height: 4),
                        const Text('Traveler names must match government-issued photo ID exactly.', style: TextStyle(fontSize: 12, color: Colors.grey, fontStyle: FontStyle.italic)),
                        const SizedBox(height: 16),
                        _buildTravelerForm(0),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),
                  if (widget.tripType == 'roundtrip')
                    Container(
                      padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 12),
                      decoration: BoxDecoration(
                        border: Border.all(color: Colors.grey.shade300),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text('Second Traveler Information (Adult 1)', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e))),
                          const SizedBox(height: 4),
                          const Text('Traveler names must match government-issued photo ID exactly.', style: TextStyle(fontSize: 12, color: Colors.grey, fontStyle: FontStyle.italic)),
                          const SizedBox(height: 16),
                          _buildTravelerForm(1),
                        ],
                      ),
                    ),
                ],
              ),
            ),
            const SizedBox(height: 8),

            // Price Summary
            Container(
              color: Colors.white,
              padding: const EdgeInsets.all(16),
              margin: const EdgeInsets.symmetric(vertical: 8),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Price Summary', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e))),
                  const SizedBox(height: 12),
                  _buildPriceRow('Base Fare', '\$${widget.price}'),
                  _buildPriceRow('Taxes', '\$0.00'),
                  _buildPriceRow('Fees', '\$0.00'),
                  const Divider(height: 16),
                  _buildPriceRow('Total', '\$${widget.price}', isBold: true),
                ],
              ),
            ),

            // Pay With Section
            Container(
              color: Colors.white,
              padding: const EdgeInsets.all(16),
              margin: const EdgeInsets.symmetric(vertical: 8),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  const Center(
                    child: Text(
                      'Pay With',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w700,
                        color: Colors.black87,
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
                  
                  // Payment Options Row 1
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
                  
                  // Payment Options Row 2
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
                  const SizedBox(height: 20),
                  
                  // Credit Card Form
                  _buildPaymentForm(),
                ],
              ),
            ),

            // Trip Revision
            Container(
              color: Colors.white,
              padding: const EdgeInsets.all(16),
              margin: const EdgeInsets.symmetric(vertical: 8),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Trip Revision', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e))),
                  const SizedBox(height: 12),
                  
                  // Point 1
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('1. ', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Colors.black87)),
                      Expanded(
                        child: Text(
                          'Review your trip details to make sure the dates and times are correct.',
                          style: const TextStyle(fontSize: 13, color: Colors.black87),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  
                  // Point 2
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('2. ', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Colors.black87)),
                      Expanded(
                        child: Text(
                          'Check your spelling. Flight passenger names must match government-issued photo ID exactly.\nTraveler 1: John Doe',
                          style: const TextStyle(fontSize: 13, color: Colors.black87),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  
                  // Point 3
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('3. ', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Colors.black87)),
                      Expanded(
                        child: Text(
                          'Review the terms of your booking:',
                          style: const TextStyle(fontSize: 13, color: Colors.black87),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  
                  // Terms sub-items
                  Padding(
                    padding: const EdgeInsets.only(left: 16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          '• In case of a no-show or cancellation, you may be entitled to a refund of airport taxes and fees included in the price of the flight purchased. In this instance, you can request such a refund from us, and we will submit your request to the airline on your behalf.',
                          style: TextStyle(fontSize: 12, color: Colors.black87),
                        ),
                        const SizedBox(height: 8),
                        const Text(
                          '• View Price Drop Protection details & disclosures',
                          style: TextStyle(fontSize: 12, color: Colors.black87),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),
                  
                  // Acknowledgement text
                  RichText(
                    text: TextSpan(
                      style: const TextStyle(fontSize: 12, color: Colors.black87),
                      children: [
                        const TextSpan(text: 'By clicking on the button below, I acknowledge that I have reviewed the '),
                        WidgetSpan(
                          child: GestureDetector(
                            onTap: () {
                              PolicyDialogs.showPrivacyPolicy(context);
                            },
                            child: const Text(
                              'Privacy Statement',
                              style: TextStyle(
                                color: Color(0xFF1e5a8e),
                                decoration: TextDecoration.underline,
                                fontSize: 12,
                              ),
                            ),
                          ),
                        ),
                        const TextSpan(text: ' and Government Travel Advice and have reviewed and accept the above Rules & Restrictions and '),
                        WidgetSpan(
                          child: GestureDetector(
                            onTap: () {
                              PolicyDialogs.showTermsAndConditions(context);
                            },
                            child: const Text(
                              'Terms of Use',
                              style: TextStyle(
                                color: Color(0xFF1e5a8e),
                                decoration: TextDecoration.underline,
                                fontSize: 12,
                              ),
                            ),
                          ),
                        ),
                        const TextSpan(text: '.'),
                      ],
                    ),
                  ),
                ],
              ),
            ),

            // Complete Booking Button
            Padding(
              padding: const EdgeInsets.all(16),
              child: SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    _showBookingSuccessDialog();
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF1e5a8e),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(25),
                    ),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                  child: const Text('Complete Booking', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600, color: Colors.white)),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFlightDetailRow(String label, String from, String to, String departureTime, String arrivalTime, String airline, String date) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('$from (${from.substring(0, 3).toUpperCase()}) to $to (${to.substring(0, 3).toUpperCase()})', style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Colors.black87)),
        const SizedBox(height: 4),
        Text('$departureTime - $arrivalTime', style: const TextStyle(fontSize: 13, color: Colors.grey)),
        const SizedBox(height: 4),
        Text(airline, style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e))),
        const SizedBox(height: 2),
        Text(widget.stops, style: const TextStyle(fontSize: 12, color: Colors.grey)),
        const SizedBox(height: 8),
        Align(
          alignment: Alignment.bottomRight,
          child: Text(date, style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Color(0xFFD32F2F))),
        ),
      ],
    );
  }

  Widget _buildPriceSummarySection() {
    return Container(
      color: const Color(0xFFFAFAFA),
      padding: const EdgeInsets.symmetric(vertical: 12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 0),
            child: Text('Price Summary', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e))),
          ),
          const SizedBox(height: 12),
          _buildPriceRow('Traveler 1: Adult', '\$${widget.price}'),
          _buildPriceRow('Fare Price', '\$${(widget.price * 0.6).toStringAsFixed(0)}'),
          _buildPriceRow('Taxes', '\$287'),
          _buildPriceRow('Service Fee', '\$500'),
          const SizedBox(height: 8),
          _buildPriceRow('Total', '\$${(widget.price + 787)}', isBold: true, isTotal: true),
        ],
      ),
    );
  }

  Widget _buildFlightInfo(String label, String from, String to, String departureTime, String arrivalTime, String airline, String stops) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Colors.grey)),
        const SizedBox(height: 8),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(from, style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w600)),
                Text(departureTime, style: const TextStyle(fontSize: 12, color: Colors.grey)),
              ],
            ),
            Column(
              children: [
                const Icon(Icons.flight, size: 20, color: Color(0xFF1e5a8e)),
                Text(stops, style: const TextStyle(fontSize: 10, color: Colors.grey)),
              ],
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(to, style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w600)),
                Text(arrivalTime, style: const TextStyle(fontSize: 12, color: Colors.grey)),
              ],
            ),
          ],
        ),
        const SizedBox(height: 4),
        Text(airline, style: const TextStyle(fontSize: 12, color: Colors.grey)),
      ],
    );
  }

  Widget _buildPriceRow(String label, String amount, {bool isBold = false, bool isTotal = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6, horizontal: 0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(
              fontSize: isTotal ? 14 : 13,
              fontWeight: isBold ? FontWeight.w600 : FontWeight.w400,
              color: Colors.black87,
            ),
          ),
          Text(
            amount,
            style: TextStyle(
              fontSize: isTotal ? 14 : 13,
              fontWeight: isBold ? FontWeight.w600 : FontWeight.w400,
              color: isTotal ? const Color(0xFFD32F2F) : Colors.black87,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTravelerForm(int travelerIndex) {
    return Column(
      children: [
        // Title and First Name Row
        Row(
          children: [
            // Title Dropdown
            SizedBox(
              width: 80,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Title', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Container(
                    decoration: const BoxDecoration(
                      border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                    ),
                    child: DropdownButtonFormField<String>(
                      decoration: const InputDecoration(
                        border: InputBorder.none,
                        enabledBorder: InputBorder.none,
                        focusedBorder: InputBorder.none,
                        isDense: true,
                        contentPadding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      ),
                      value: 'Mr',
                      items: ['Mr', 'Mrs', 'Ms', 'Dr'].map((String value) {
                        return DropdownMenuItem<String>(
                          value: value,
                          child: Text(value, style: const TextStyle(fontSize: 14)),
                        );
                      }).toList(),
                      onChanged: (String? newValue) {
                        // Will be populated from DB later
                      },
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('First Name', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Container(
                    decoration: const BoxDecoration(
                      border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                    ),
                    child: TextField(
                      decoration: InputDecoration(
                        hintText: 'John',
                        border: InputBorder.none,
                        enabledBorder: InputBorder.none,
                        focusedBorder: InputBorder.none,
                        isDense: true,
                        contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        // Last Name Row
        Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Last Name', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Container(
                    decoration: const BoxDecoration(
                      border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                    ),
                    child: TextField(
                      decoration: InputDecoration(
                        border: InputBorder.none,
                        enabledBorder: InputBorder.none,
                        focusedBorder: InputBorder.none,
                        isDense: true,
                        contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        
        // Passport Number and Passport Expiry Date Row
        Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Passport number', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Container(
                    decoration: const BoxDecoration(
                      border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                    ),
                    child: TextField(
                      decoration: InputDecoration(
                        border: InputBorder.none,
                        enabledBorder: InputBorder.none,
                        focusedBorder: InputBorder.none,
                        isDense: true,
                        contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Passport Expiry Date', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Expanded(
                        child: Container(
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: TextField(
                            decoration: InputDecoration(
                              hintText: 'Day',
                              border: InputBorder.none,
                              enabledBorder: InputBorder.none,
                              focusedBorder: InputBorder.none,
                              isDense: true,
                              contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 4),
                      const Text('-', style: TextStyle(fontSize: 12, color: Colors.grey)),
                      const SizedBox(width: 4),
                      Expanded(
                        child: Container(
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: TextField(
                            decoration: InputDecoration(
                              hintText: 'Month',
                              border: InputBorder.none,
                              enabledBorder: InputBorder.none,
                              focusedBorder: InputBorder.none,
                              isDense: true,
                              contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 4),
                      const Text('-', style: TextStyle(fontSize: 12, color: Colors.grey)),
                      const SizedBox(width: 4),
                      Expanded(
                        child: Container(
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: TextField(
                            decoration: InputDecoration(
                              hintText: 'Year',
                              border: InputBorder.none,
                              enabledBorder: InputBorder.none,
                              focusedBorder: InputBorder.none,
                              isDense: true,
                              contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        
        // Country/Territory Code and Phone Number Row
        Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Country/ Territory Code', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Container(
                    decoration: const BoxDecoration(
                      border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                    ),
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text('Select', style: TextStyle(fontSize: 12, color: Colors.grey)),
                        Icon(Icons.check_circle, color: Colors.blue.shade600, size: 18),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Phone Number', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Container(
                    decoration: const BoxDecoration(
                      border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                    ),
                    child: TextField(
                      decoration: InputDecoration(
                        border: InputBorder.none,
                        enabledBorder: InputBorder.none,
                        focusedBorder: InputBorder.none,
                        isDense: true,
                        contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        
        // Date of Birth and Gender Row
        Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Date of birth', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Expanded(
                        child: Container(
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: TextField(
                            decoration: InputDecoration(
                              hintText: 'Day',
                              border: InputBorder.none,
                              enabledBorder: InputBorder.none,
                              focusedBorder: InputBorder.none,
                              isDense: true,
                              contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 4),
                      const Text('-', style: TextStyle(fontSize: 12, color: Colors.grey)),
                      const SizedBox(width: 4),
                      Expanded(
                        child: Container(
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: TextField(
                            decoration: InputDecoration(
                              hintText: 'Month',
                              border: InputBorder.none,
                              enabledBorder: InputBorder.none,
                              focusedBorder: InputBorder.none,
                              isDense: true,
                              contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 4),
                      const Text('-', style: TextStyle(fontSize: 12, color: Colors.grey)),
                      const SizedBox(width: 4),
                      Expanded(
                        child: Container(
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: TextField(
                            decoration: InputDecoration(
                              hintText: 'Year',
                              border: InputBorder.none,
                              enabledBorder: InputBorder.none,
                              focusedBorder: InputBorder.none,
                              isDense: true,
                              contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Gender', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Flexible(
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Radio<String>(
                              value: 'Male',
                              groupValue: 'Male',
                              onChanged: (value) {},
                              activeColor: const Color(0xFF1e5a8e),
                            ),
                            const Flexible(
                              child: Text('Male', style: TextStyle(fontSize: 12), overflow: TextOverflow.ellipsis),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 4),
                      Flexible(
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Radio<String>(
                              value: 'Female',
                              groupValue: 'Male',
                              onChanged: (value) {},
                              activeColor: const Color(0xFF1e5a8e),
                            ),
                            const Flexible(
                              child: Text('Female', style: TextStyle(fontSize: 12), overflow: TextOverflow.ellipsis),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildPaymentOption(String name, String iconPath) {
    final isSelected = _selectedPayment == 0 && name == 'Whish Money' ||
                       _selectedPayment == 1 && name == 'OMT Pay' ||
                       _selectedPayment == 2 && name == 'Crypto Pay' ||
                       _selectedPayment == 3 && name == 'Visa';
    
    return GestureDetector(
      onTap: () {
        setState(() {
          if (name == 'Whish Money') _selectedPayment = 0;
          else if (name == 'OMT Pay') _selectedPayment = 1;
          else if (name == 'Crypto Pay') _selectedPayment = 2;
          else if (name == 'Visa') _selectedPayment = 3;
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
          mainAxisSize: MainAxisSize.min,
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

  Widget _buildPaymentForm() {
    return Column(
      children: [
        // Name on Card and Debit/Credit card number Row
        Row(
          children: [
            Expanded(
              child: Container(
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                child: TextField(
                  decoration: InputDecoration(
                    hintText: 'Name on Card',
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                child: const Text('Select Card', style: TextStyle(fontSize: 12, color: Colors.grey)),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        
        // Expiration Date and Security Code Row
        Row(
          children: [
            Expanded(
              child: Row(
                children: [
                  Expanded(
                    child: Container(
                      decoration: const BoxDecoration(
                        border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                      ),
                      child: TextField(
                        decoration: InputDecoration(
                          hintText: 'MM',
                          border: InputBorder.none,
                          enabledBorder: InputBorder.none,
                          focusedBorder: InputBorder.none,
                          isDense: true,
                          contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 4),
                  const Text('-', style: TextStyle(fontSize: 12, color: Colors.grey)),
                  const SizedBox(width: 4),
                  Expanded(
                    child: Container(
                      decoration: const BoxDecoration(
                        border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                      ),
                      child: TextField(
                        decoration: InputDecoration(
                          hintText: 'YY',
                          border: InputBorder.none,
                          enabledBorder: InputBorder.none,
                          focusedBorder: InputBorder.none,
                          isDense: true,
                          contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextField(
                obscureText: true,
                decoration: InputDecoration(
                  hintText: 'Security Code',
                  border: InputBorder.none,
                  enabledBorder: InputBorder.none,
                  focusedBorder: InputBorder.none,
                  isDense: true,
                  contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        
        // Country/Territory Code and Billing Address Row
        Row(
          children: [
            Expanded(
              child: Container(
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Select', style: TextStyle(fontSize: 12, color: Colors.grey)),
                    Icon(Icons.check_circle, color: Colors.blue.shade600, size: 18),
                  ],
                ),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Container(
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                child: TextField(
                  decoration: InputDecoration(
                    hintText: 'Billing Address 1',
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  ),
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        
        // Billing Address 2 and City Row
        Row(
          children: [
            Expanded(
              child: Container(
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                child: TextField(
                  decoration: InputDecoration(
                    hintText: 'Billing Address 2',
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Container(
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                child: TextField(
                  decoration: InputDecoration(
                    hintText: 'City',
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  ),
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        
        // State and ZIP Code Row
        Row(
          children: [
            Expanded(
              child: Container(
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Select', style: TextStyle(fontSize: 12, color: Colors.grey)),
                    Icon(Icons.check_circle, color: Colors.blue.shade600, size: 18),
                  ],
                ),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Container(
                decoration: const BoxDecoration(
                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                ),
                child: TextField(
                  decoration: InputDecoration(
                    hintText: 'ZIP code',
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  ),
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  void _showBookingSuccessDialog() {
    final bookingReference = 'B${DateTime.now().millisecondsSinceEpoch.toString().substring(0, 10)}';
    
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Container(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Close button
                Align(
                  alignment: Alignment.topRight,
                  child: GestureDetector(
                    onTap: () {
                      Navigator.of(context).pushAndRemoveUntil(
                        MaterialPageRoute(builder: (context) => HomePage()),
                        (route) => false,
                      );
                      Navigator.push(
                        context,
                        MaterialPageRoute(builder: (context) => BookingsPage()),
                      );
                    },
                    child: Container(
                      width: 32,
                      height: 32,
                      decoration: const BoxDecoration(
                        color: Color(0xFFD32F2F),
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(Icons.close, color: Colors.white, size: 18),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                
                // Plane Logo
                Image.asset(
                  'assets/images/planez.png',
                  width: 80,
                  height: 80,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      width: 80,
                      height: 80,
                      decoration: BoxDecoration(
                        color: Colors.grey.shade200,
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(Icons.flight, size: 40, color: Color(0xFF1e5a8e)),
                    );
                  },
                ),
                const SizedBox(height: 24),
                
                // Success Message
                const Text(
                  'Your flight has been booked successfully.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 12),
                
                // Booking Reference Label
                const Text(
                  'Your booking reference is',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.grey,
                  ),
                ),
                const SizedBox(height: 8),
                
                // Booking Reference Number
                Text(
                  bookingReference,
                  textAlign: TextAlign.center,
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w700,
                    color: Color(0xFF1e5a8e),
                  ),
                ),
                const SizedBox(height: 16),
                
                // Email Confirmation Text
                const Text(
                  'You will receive a confirmation email including your itinerary details.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                ),
                const SizedBox(height: 8),
                
                // Booking Details Text
                const Text(
                  'You can check for the full details in your bookings page.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                ),
                const SizedBox(height: 24),
                
                // Done Button
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.of(context).pushAndRemoveUntil(
                        MaterialPageRoute(builder: (context) => HomePage()),
                        (route) => false,
                      );
                      Navigator.push(
                        context,
                        MaterialPageRoute(builder: (context) => BookingsPage()),
                      );
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF1e5a8e),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(25),
                      ),
                      padding: const EdgeInsets.symmetric(vertical: 12),
                    ),
                    child: const Text(
                      'Done',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: Colors.white,
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
}

