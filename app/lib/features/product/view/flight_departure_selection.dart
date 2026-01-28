import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'flight_return_selection_screen.dart';
import 'special_services_screen.dart';

class FlightDepartureSelection extends StatefulWidget {
  final String from;
  final String to;
  final String departureTime;
  final String arrivalTime;
  final String duration;
  final String stops;
  final String tripType; // 'oneway', 'roundtrip', 'multicity'
  final int price;
  final String airline;
  final String departureDate;
  final String? returnDate;

  const FlightDepartureSelection({
    Key? key,
    required this.from,
    required this.to,
    required this.departureTime,
    required this.arrivalTime,
    required this.duration,
    required this.stops,
    required this.tripType,
    required this.price,
    this.airline = 'Middle East Airlines',
    required this.departureDate,
    this.returnDate,
  }) : super(key: key);

  @override
  State<FlightDepartureSelection> createState() => _FlightDepartureSelectionState();
}

class _FlightDepartureSelectionState extends State<FlightDepartureSelection> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  
  final List<Map<String, dynamic>> _cabinClasses = [
    {
      'name': 'Economy',
      'displayName': 'Cabin: Economy',
      'flexibility': 'Flexible',
      'seatChoice': 'Seat choice available',
      'cabinBaggage': '1 Piece of cabin baggage',
      'underseatBag': '1 Piece of Underseat Bag Only',
      'checkedBaggage': '20 Kg Check-in Baggage',
      'cancellationFee': '\$105',
      'changeFee': '\$105',
    },
    {
      'name': 'Business',
      'displayName': 'Cabin: Business',
      'flexibility': 'Flexible',
      'seatChoice': 'Premium seat selection',
      'cabinBaggage': '2 Pieces of cabin baggage',
      'underseatBag': '1 Piece of Underseat Bag',
      'checkedBaggage': '32 Kg Check-in Baggage',
      'cancellationFee': '\$75',
      'changeFee': '\$75',
    },
    {
      'name': 'First Class',
      'displayName': 'Cabin: First Class',
      'flexibility': 'Fully Flexible',
      'seatChoice': 'Suite selection included',
      'cabinBaggage': '2 Pieces of cabin baggage',
      'underseatBag': '1 Piece of Underseat Bag',
      'checkedBaggage': '40 Kg Check-in Baggage',
      'cancellationFee': 'Free',
      'changeFee': 'Free',
    },
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: Column(
          children: [
            // Header with back button, title and flight info
            _buildHeader(),
            
            // Class Tabs
            TabBar(
              controller: _tabController,
              indicatorColor: const Color(0xFF1e5a8e),
              indicatorWeight: 2,
              dividerColor: Colors.transparent,
              labelColor: const Color(0xFF1e5a8e),
              unselectedLabelColor: Colors.grey,
              labelStyle: const TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
              unselectedLabelStyle: const TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w400,
              ),
              tabs: const [
                Tab(text: 'Economy'),
                Tab(text: 'Business'),
                Tab(text: 'First Class'),
              ],
            ),
            
            // Content
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  _buildClassContent(0),
                  _buildClassContent(1),
                  _buildClassContent(2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 12),
      child: Column(
        children: [
          // Title row with back button
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.arrow_back_ios, color: Color(0xFF1e5a8e), size: 20),
                onPressed: () => Navigator.pop(context),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
              Expanded(
                child: Text(
                  'Select Departure Flight',
                  textAlign: TextAlign.center,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                ),
              ),
              const SizedBox(width: 28), // Balance for back button
            ],
          ),
          const SizedBox(height: 8),
          // Flight info row directly below title
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: Row(
              children: [
                // Time and flight path
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            widget.departureTime,
                            style: const TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w500,
                              color: Colors.black87,
                            ),
                          ),
                          const SizedBox(width: 6),
                          Container(
                            width: 20,
                            height: 1,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(width: 4),
                          Icon(Icons.flight, size: 14, color: Colors.grey.shade600),
                          const SizedBox(width: 4),
                          Container(
                            width: 20,
                            height: 1,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(width: 6),
                          Text(
                            widget.arrivalTime,
                            style: const TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w500,
                              color: Colors.black87,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 2),
                      Row(
                        children: [
                          Icon(Icons.airline_seat_recline_normal, size: 12, color: Colors.grey.shade500),
                          const SizedBox(width: 4),
                          Text(
                            widget.airline,
                            style: TextStyle(
                              fontSize: 11,
                              color: Colors.grey.shade600,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                // Duration and stops
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      widget.duration,
                      style: const TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                        color: Colors.black87,
                      ),
                    ),
                    Text(
                      widget.stops,
                      style: TextStyle(
                        fontSize: 11,
                        color: Colors.grey.shade600,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildClassContent(int classIndex) {
    final cabinClass = _cabinClasses[classIndex];
    
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Flight Card with ccl.png background
          _buildFlightCard(cabinClass),
          
          const SizedBox(height: 24),
          
          // Terms and conditions
          _buildTermsAndConditions(),
        ],
      ),
    );
  }

  Widget _buildFlightCard(Map<String, dynamic> cabinClass) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: Stack(
          children: [
            // Background Image (ccl.png)
            Positioned.fill(
              child: Image.asset(
                'assets/images/ccl.png',
                fit: BoxFit.cover,
              ),
            ),
            
            // Content
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Top row: Time and Flexibility
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Left: Time and duration
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              '${widget.departureTime}- ${widget.arrivalTime}',
                              style: const TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w700,
                                color: Colors.black87,
                              ),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              '${widget.duration}, ${widget.stops}',
                              style: const TextStyle(
                                fontSize: 12,
                                color: Colors.black54,
                              ),
                            ),
                          ],
                        ),
                      ),
                      
                      // Right: Flexibility and Cabin type
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.end,
                        children: [
                          Text(
                            cabinClass['flexibility'],
                            style: const TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                              color: Colors.black54,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(
                            cabinClass['displayName'],
                            style: const TextStyle(
                              fontSize: 12,
                              color: Colors.black54,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  
                  const SizedBox(height: 16),
                  
                  // Seat section
                  const Text(
                    'Seat',
                    style: TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 6),
                  Row(
                    children: [
                      const Icon(
                        Icons.check_circle,
                        color: Color(0xFF4CAF50),
                        size: 16,
                      ),
                      const SizedBox(width: 6),
                      Text(
                        cabinClass['seatChoice'],
                        style: const TextStyle(
                          fontSize: 12,
                          color: Colors.black87,
                        ),
                      ),
                    ],
                  ),
                  
                  const SizedBox(height: 16),
                  
                  // Bags section
                  const Text(
                    'Bags',
                    style: TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 6),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Left column
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              cabinClass['cabinBaggage'],
                              style: const TextStyle(
                                fontSize: 12,
                                color: Colors.black87,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              cabinClass['underseatBag'],
                              style: const TextStyle(
                                fontSize: 12,
                                color: Colors.black87,
                              ),
                            ),
                          ],
                        ),
                      ),
                      // Right column
                      Text(
                        cabinClass['checkedBaggage'],
                        style: const TextStyle(
                          fontSize: 12,
                          color: Colors.black87,
                        ),
                      ),
                    ],
                  ),
                  
                  const SizedBox(height: 16),
                  
                  // Flexibility section
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      // Left: Fees
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Flexibility',
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: Color(0xFF1e5a8e),
                              ),
                            ),
                            const SizedBox(height: 6),
                            Row(
                              children: [
                                const Text(
                                  'Cancellation Fee',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Colors.black54,
                                  ),
                                ),
                                const SizedBox(width: 8),
                                Text(
                                  cabinClass['cancellationFee'],
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: cabinClass['cancellationFee'] == 'Free' 
                                        ? const Color(0xFF4CAF50) 
                                        : Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 4),
                            Row(
                              children: [
                                const Text(
                                  'Change Fee',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Colors.black54,
                                  ),
                                ),
                                const SizedBox(width: 8),
                                Text(
                                  cabinClass['changeFee'],
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: cabinClass['changeFee'] == 'Free' 
                                        ? const Color(0xFF4CAF50) 
                                        : Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                      
                      // Right: Select button
                      ElevatedButton(
                        onPressed: () {
                          _onSelectFlight(cabinClass);
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF1e5a8e),
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(horizontal: 28, vertical: 10),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(6),
                          ),
                          elevation: 0,
                        ),
                        child: const Text(
                          'Select',
                          style: TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
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
      ),
    );
  }

  Widget _buildTermsAndConditions() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Terms and conditions',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: Color(0xFF1e5a8e),
          ),
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: Text(
                'Terms and conditions apply',
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.grey.shade600,
                ),
              ),
            ),
            Expanded(
              child: Text(
                'Terms and conditions apply',
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.grey.shade600,
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  void _onSelectFlight(Map<String, dynamic> cabinClass) {
    if (widget.tripType == 'roundtrip') {
      // Roundtrip: push return flight selection, then special services after return selected
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => FlightReturnSelectionScreen(
            from: widget.to,
            to: widget.from,
            departureTime: '14:30', // Example, should be dynamic
            arrivalTime: '18:45',   // Example, should be dynamic
            duration: widget.duration,
            stops: widget.stops,
            price: widget.price,
            airline: widget.airline,
            tripType: widget.tripType,
            departureDate: widget.departureDate,
            returnDate: widget.returnDate,
          ),
        ),
      ).then((_) {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => SpecialServicesScreen(
              tripType: widget.tripType,
              departureDate: widget.departureDate,
              returnDate: widget.returnDate,
            ),
          ),
        );
      });
    } else {
      // One-way and multicity: go directly to special services
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => SpecialServicesScreen(
            tripType: widget.tripType,
            departureDate: widget.departureDate,
            returnDate: widget.returnDate,
          ),
        ),
      );
    }
  }

  void _showFlightSelectedConfirmation(Map<String, dynamic> cabinClass) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(
                Icons.check_circle,
                color: Color(0xFF4CAF50),
                size: 60,
              ),
              const SizedBox(height: 16),
              const Text(
                'Flight Selected!',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w700,
                  color: Colors.black87,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                '${widget.from} → ${widget.to}',
                style: const TextStyle(
                  fontSize: 14,
                  color: Colors.black54,
                ),
              ),
              Text(
                '${cabinClass['name']}',
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: Color(0xFF1e5a8e),
                ),
              ),
              const SizedBox(height: 20),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    Navigator.pop(context);
                    _navigateToCheckout(cabinClass);
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF1e5a8e),
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text(
                    'Continue to Checkout',
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
        );
      },
    );
  }

  void _navigateToCheckout(Map<String, dynamic> cabinClass) {
    // Navigate to flight checkout (similar to transfer_checkout.dart)
    Navigator.pop(context, {
      'cabinClass': cabinClass,
      'price': widget.price,
    });
  }
}

// Return Flight Selection for Round Trip
class FlightReturnSelection extends StatefulWidget {
  final String from;
  final String to;
  final String departureTime;
  final String arrivalTime;
  final String duration;
  final String stops;
  final String tripType;
  final Map<String, dynamic> outboundCabinClass;
  final int outboundPrice;

  const FlightReturnSelection({
    Key? key,
    required this.from,
    required this.to,
    required this.departureTime,
    required this.arrivalTime,
    required this.duration,
    required this.stops,
    required this.tripType,
    required this.outboundCabinClass,
    required this.outboundPrice,
  }) : super(key: key);

  @override
  State<FlightReturnSelection> createState() => _FlightReturnSelectionState();
}

class _FlightReturnSelectionState extends State<FlightReturnSelection> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  
  final List<Map<String, dynamic>> _cabinClasses = [
    {
      'name': 'Economy',
      'displayName': 'Cabin: Economy',
      'flexibility': 'Flexible',
      'seatChoice': 'Seat choice available',
      'cabinBaggage': '1 Piece of cabin baggage',
      'underseatBag': '1 Piece of Underseat Bag Only',
      'checkedBaggage': '20 Kg Check-in Baggage',
      'cancellationFee': '\$105',
      'changeFee': '\$105',
    },
    {
      'name': 'Business',
      'displayName': 'Cabin: Business',
      'flexibility': 'Flexible',
      'seatChoice': 'Premium seat selection',
      'cabinBaggage': '2 Pieces of cabin baggage',
      'underseatBag': '1 Piece of Underseat Bag',
      'checkedBaggage': '32 Kg Check-in Baggage',
      'cancellationFee': '\$75',
      'changeFee': '\$75',
    },
    {
      'name': 'First Class',
      'displayName': 'Cabin: First Class',
      'flexibility': 'Fully Flexible',
      'seatChoice': 'Suite selection included',
      'cabinBaggage': '2 Pieces of cabin baggage',
      'underseatBag': '1 Piece of Underseat Bag',
      'checkedBaggage': '40 Kg Check-in Baggage',
      'cancellationFee': 'Free',
      'changeFee': 'Free',
    },
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: Column(
          children: [
            // Header with back button, title and flight info
            _buildHeader(),
            
            // Outbound Flight Summary
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
              color: const Color(0xFFE8F5E9),
              child: Row(
                children: [
                  const Icon(Icons.check_circle, color: Color(0xFF4CAF50), size: 18),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      'Outbound: ${widget.to} → ${widget.from} • ${widget.outboundCabinClass['name']}',
                      style: const TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                        color: Color(0xFF2E7D32),
                      ),
                    ),
                  ),
                  Text(
                    '\$${widget.outboundPrice}',
                    style: const TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF2E7D32),
                    ),
                  ),
                ],
              ),
            ),
            
            // Class Tabs
            TabBar(
              controller: _tabController,
              indicatorColor: const Color(0xFF1e5a8e),
              indicatorWeight: 2,
              dividerColor: Colors.transparent,
              labelColor: const Color(0xFF1e5a8e),
              unselectedLabelColor: Colors.grey,
              labelStyle: const TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
              unselectedLabelStyle: const TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w400,
              ),
              tabs: const [
                Tab(text: 'Economy'),
                Tab(text: 'Business'),
                Tab(text: 'First Class'),
              ],
            ),
            
            // Content
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  _buildClassContent(0),
                  _buildClassContent(1),
                  _buildClassContent(2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 12),
      child: Column(
        children: [
          // Title row with back button
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.arrow_back_ios, color: Color(0xFF1e5a8e), size: 20),
                onPressed: () => Navigator.pop(context),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
              Expanded(
                child: Text(
                  'Select Return Flight',
                  textAlign: TextAlign.center,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                ),
              ),
              const SizedBox(width: 28),
            ],
          ),
          const SizedBox(height: 8),
          // Flight info row directly below title
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: Row(
              children: [
                // Time and flight path
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            widget.departureTime,
                            style: const TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w500,
                              color: Colors.black87,
                            ),
                          ),
                          const SizedBox(width: 6),
                          Container(
                            width: 20,
                            height: 1,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(width: 4),
                          Icon(Icons.flight, size: 14, color: Colors.grey.shade600),
                          const SizedBox(width: 4),
                          Container(
                            width: 20,
                            height: 1,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(width: 6),
                          Text(
                            widget.arrivalTime,
                            style: const TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w500,
                              color: Colors.black87,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 2),
                      Row(
                        children: [
                          Icon(Icons.airline_seat_recline_normal, size: 12, color: Colors.grey.shade500),
                          const SizedBox(width: 4),
                          Text(
                            'Return Flight',
                            style: TextStyle(
                              fontSize: 11,
                              color: Colors.grey.shade600,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                // Duration and stops
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      widget.duration,
                      style: const TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                        color: Colors.black87,
                      ),
                    ),
                    Text(
                      widget.stops,
                      style: TextStyle(
                        fontSize: 11,
                        color: Colors.grey.shade600,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildClassContent(int classIndex) {
    final cabinClass = _cabinClasses[classIndex];
    
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildFlightCard(cabinClass),
          const SizedBox(height: 24),
          _buildTermsAndConditions(),
        ],
      ),
    );
  }

  Widget _buildFlightCard(Map<String, dynamic> cabinClass) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: Stack(
          children: [
            Positioned.fill(
              child: Image.asset(
                'assets/images/ccl.png',
                fit: BoxFit.cover,
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              '${widget.departureTime}- ${widget.arrivalTime}',
                              style: const TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w700,
                                color: Colors.black87,
                              ),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              '${widget.duration}, ${widget.stops}',
                              style: const TextStyle(
                                fontSize: 12,
                                color: Colors.black54,
                              ),
                            ),
                          ],
                        ),
                      ),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.end,
                        children: [
                          Text(
                            cabinClass['flexibility'],
                            style: const TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                              color: Colors.black54,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(
                            cabinClass['displayName'],
                            style: const TextStyle(
                              fontSize: 12,
                              color: Colors.black54,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Seat',
                    style: TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 6),
                  Row(
                    children: [
                      const Icon(Icons.check_circle, color: Color(0xFF4CAF50), size: 16),
                      const SizedBox(width: 6),
                      Text(
                        cabinClass['seatChoice'],
                        style: const TextStyle(fontSize: 12, color: Colors.black87),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Bags',
                    style: TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 6),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              cabinClass['cabinBaggage'],
                              style: const TextStyle(fontSize: 12, color: Colors.black87),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              cabinClass['underseatBag'],
                              style: const TextStyle(fontSize: 12, color: Colors.black87),
                            ),
                          ],
                        ),
                      ),
                      Text(
                        cabinClass['checkedBaggage'],
                        style: const TextStyle(fontSize: 12, color: Colors.black87),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Flexibility',
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: Color(0xFF1e5a8e),
                              ),
                            ),
                            const SizedBox(height: 6),
                            Row(
                              children: [
                                const Text(
                                  'Cancellation Fee',
                                  style: TextStyle(fontSize: 12, color: Colors.black54),
                                ),
                                const SizedBox(width: 8),
                                Text(
                                  cabinClass['cancellationFee'],
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: cabinClass['cancellationFee'] == 'Free' 
                                        ? const Color(0xFF4CAF50) 
                                        : Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 4),
                            Row(
                              children: [
                                const Text(
                                  'Change Fee',
                                  style: TextStyle(fontSize: 12, color: Colors.black54),
                                ),
                                const SizedBox(width: 8),
                                Text(
                                  cabinClass['changeFee'],
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: cabinClass['changeFee'] == 'Free' 
                                        ? const Color(0xFF4CAF50) 
                                        : Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                      ElevatedButton(
                        onPressed: () {
                          _onSelectFlight(cabinClass);
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF1e5a8e),
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(horizontal: 28, vertical: 10),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(6),
                          ),
                          elevation: 0,
                        ),
                        child: const Text(
                          'Select',
                          style: TextStyle(fontSize: 14, fontWeight: FontWeight.w600),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTermsAndConditions() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Terms and conditions',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: Color(0xFF1e5a8e),
          ),
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: Text(
                'Terms and conditions apply',
                style: TextStyle(fontSize: 12, color: Colors.grey.shade600),
              ),
            ),
            Expanded(
              child: Text(
                'Terms and conditions apply',
                style: TextStyle(fontSize: 12, color: Colors.grey.shade600),
              ),
            ),
          ],
        ),
      ],
    );
  }

  void _onSelectFlight(Map<String, dynamic> cabinClass) {
    final int totalPrice = widget.outboundPrice + widget.outboundPrice; // Using same price for return
    
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(
                Icons.check_circle,
                color: Color(0xFF4CAF50),
                size: 60,
              ),
              const SizedBox(height: 16),
              const Text(
                'Flights Selected!',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w700,
                  color: Colors.black87,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Round Trip: ${widget.to} ↔ ${widget.from}',
                style: const TextStyle(
                  fontSize: 14,
                  color: Colors.black54,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                'Total: \$${totalPrice.toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (Match m) => '${m[1]},')}',
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w700,
                  color: Color(0xFF1e5a8e),
                ),
              ),
              const SizedBox(height: 20),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    Navigator.pop(context);
                    Navigator.pop(context, {
                      'outboundCabinClass': widget.outboundCabinClass,
                      'returnCabinClass': cabinClass,
                      'totalPrice': totalPrice,
                    });
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF1e5a8e),
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text(
                    'Continue to Checkout',
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
        );
      },
    );
  }
}
