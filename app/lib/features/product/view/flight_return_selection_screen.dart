import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'special_services_screen.dart';

class FlightReturnSelectionScreen extends StatefulWidget {
  final String from;
  final String to;
  final String departureTime;
  final String arrivalTime;
  final String duration;
  final String stops;
  final int price;
  final String airline;
  final String tripType;
  final String departureDate;
  final String? returnDate;

  const FlightReturnSelectionScreen({
    Key? key,
    required this.from,
    required this.to,
    required this.departureTime,
    required this.arrivalTime,
    required this.duration,
    required this.stops,
    required this.price,
    required this.airline,
    this.tripType = 'roundtrip',
    required this.departureDate,
    this.returnDate,
  }) : super(key: key);

  @override
  State<FlightReturnSelectionScreen> createState() => _FlightReturnSelectionScreenState();
}

class _FlightReturnSelectionScreenState extends State<FlightReturnSelectionScreen> with SingleTickerProviderStateMixin {
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
      'cancellationFee': '105',
      'changeFee': '105',
    },
    {
      'name': 'Business',
      'displayName': 'Cabin: Business',
      'flexibility': 'Flexible',
      'seatChoice': 'Premium seat selection',
      'cabinBaggage': '2 Pieces of cabin baggage',
      'underseatBag': '1 Piece of Underseat Bag',
      'checkedBaggage': '32 Kg Check-in Baggage',
      'cancellationFee': '75',
      'changeFee': '75',
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
            _buildHeader(),
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
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: Row(
              children: [
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
                          // After selecting return flight, go to special services
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
}
