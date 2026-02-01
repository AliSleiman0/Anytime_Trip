import 'package:flutter/material.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/network/endpoints.dart';

class BookingsPage extends StatefulWidget {
  const BookingsPage({super.key});

  @override
  State<BookingsPage> createState() => _BookingsPageState();
}

class _BookingsPageState extends State<BookingsPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  final List<Map<String, dynamic>> _bookings = [
    {
      'type': 'flights',
      'status': 'Upcoming',
      'from': 'Bey',
      'to': 'DXB',
      'departureTime': '8:15 am',
      'arrivalTime': '10:25 am',
      'duration': '1 hr',
      'stops': 'One way',
      'airline': 'Emirates',
      'cabin': 'Cabin: Economy',
      'bags': '1 carry-on bag included (7 kg)\n2 checked bags (23 kg each)',
      'cancellationFee': '\$250',
      'ticketPrice': '\$1,300',
      'changeFee': '\$250',
      'upgradeFee': '\$250',
      'totalPrice': '\$1,800',
    },
    {
      'type': 'flights',
      'status': 'Upcoming',
      'from': 'Bey',
      'to': 'DXB',
      'departureTime': '7:21 am',
      'arrivalTime': '9:35 am',
      'duration': '1 hr',
      'stops': 'One way',
      'airline': 'Middle East Airlines',
      'cabin': 'Cabin: Economy',
      'bags': '1 carry-on bag included (7 kg)\n1 checked bag (23 kg each)',
      'cancellationFee': '\$195',
      'ticketPrice': '\$1,195',
      'changeFee': '\$195',
      'upgradeFee': '\$195',
      'totalPrice': '\$1,595',
    },
    {
      'type': 'flights',
      'status': 'Cancelled',
      'from': 'Bey',
      'to': 'DOH',
      'departureTime': '6:00 am',
      'arrivalTime': '8:00 am',
      'duration': '2 hr',
      'stops': 'One way',
      'airline': 'Qatar Airways',
      'cabin': 'Cabin: Economy',
      'bags': '1 carry-on bag included (10 kg)\n2 checked bags (23 kg each)',
      'cancellationFee': '\$200',
      'ticketPrice': '\$1,300',
      'changeFee': '\$200',
      'upgradeFee': '\$200',
      'totalPrice': '\$1,400',
    },
    {
      'type': 'flights',
      'status': 'Postponed',
      'from': 'Bey',
      'to': 'DOH',
      'departureTime': '9:00 am',
      'arrivalTime': '11:15 am',
      'duration': '2 hr',
      'stops': 'One way',
      'airline': 'Qatar Airways',
      'cabin': 'Cabin: Economy',
      'bags': '1 carry-on bag included (10 kg)\n2 checked bags (23 kg each)',
      'cancellationFee': '\$200',
      'ticketPrice': '\$1,300',
      'changeFee': '\$200',
      'upgradeFee': '\$200',
      'totalPrice': '\$1,400',
    },
    {
      'type': 'cars',
      'status': 'Upcoming',
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'type': 'cars',
      'status': 'Voided',
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'type': 'cars',
      'status': 'Cancelled',
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'type': 'cars',
      'status': 'Cancelled',
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'type': 'hotels',
      'status': 'Upcoming',
      'hotelName': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'amenities': 'Breakfast + Pool Included',
      'bookedDate': 'Booked for 11/11/2025',
      'price': '\$603 total',
    },
    {
      'type': 'hotels',
      'status': 'Voided',
      'hotelName': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'amenities': 'Breakfast + Pool Included',
      'bookedDate': 'Booked for 11/11/2025',
      'price': '\$603 total',
    },
    {
      'type': 'hotels',
      'status': 'Cancelled',
      'hotelName': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'amenities': 'Breakfast + Pool Included',
      'bookedDate': 'Booked for 11/11/2025',
      'price': '\$603 total',
    },
    {
      'type': 'hotels',
      'status': 'Cancelled',
      'hotelName': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'amenities': 'Breakfast + Pool Included',
      'bookedDate': 'Booked for 11/11/2025',
      'price': '\$603 total',
    },
    {
      'type': 'transfers',
      'status': 'Upcoming',
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'estimatedTime': 'Estimated time: 40 min',
      'meetGreet': 'Meet & Greet availability',
      'price': '\$50',
    },
    {
      'type': 'transfers',
      'status': 'Upcoming',
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'estimatedTime': 'Estimated time: 40 min',
      'meetGreet': 'Meet & Greet availability',
      'price': '\$50',
    },
    {
      'type': 'transfers',
      'status': 'Cancelled',
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'estimatedTime': 'Estimated time: 40 min',
      'meetGreet': 'Meet & Greet availability',
      'price': '\$50',
    },
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    _tabController.addListener(() {
      setState(() {});
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
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
      body: Column(
        children: [
          // Title
          Padding(
            padding: const EdgeInsets.all(16),
            child: const Text(
              'Bookings History',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w600,
                color: Colors.black87,
              ),
            ),
          ),

          // Tabs - Custom style like home page
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                children: [
                  _buildTabButton(0, 'Flights'),
                  const SizedBox(width: 12),
                  _buildTabButton(1, 'Cars'),
                  const SizedBox(width: 12),
                  _buildTabButton(2, 'Hotels'),
                  const SizedBox(width: 12),
                  _buildTabButton(3, 'Transfers'),
                ],
              ),
            ),
          ),

          const SizedBox(height: 20),

          // Tab Content
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                _buildFlightsTab(),
                _buildCarsTab(),
                _buildHotelsTab(),
                _buildTransfersTab(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTabButton(int index, String label) {
    final isSelected = _tabController.index == index;
    return GestureDetector(
      onTap: () {
        _tabController.animateTo(index);
        setState(() {});
      },
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
        decoration: BoxDecoration(
          color: isSelected ? const Color(0xFF1e5a8e) : Colors.white,
          borderRadius: BorderRadius.circular(25),
          border: !isSelected
              ? Border.all(color: Colors.grey[400]!, width: 1)
              : null,
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.1),
              blurRadius: 8,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        child: Text(
          label,
          style: TextStyle(
            fontSize: 13,
            fontWeight: FontWeight.w600,
            color: isSelected ? Colors.white : Colors.grey[700],
          ),
        ),
      ),
    );
  }

  Widget _buildFlightsTab() {
    final flights = _bookings.where((b) => b['type'] == 'flights').toList();
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: flights.map((flight) => _buildBookingCard(flight)).toList(),
        ),
      ),
    );
  }

  Widget _buildCarsTab() {
    final cars = _bookings.where((b) => b['type'] == 'cars').toList();
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: cars.map((car) => _buildCarBookingCard(car)).toList(),
        ),
      ),
    );
  }

  Widget _buildHotelsTab() {
    final hotels = _bookings.where((b) => b['type'] == 'hotels').toList();
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: hotels.map((hotel) => _buildHotelBookingCard(hotel)).toList(),
        ),
      ),
    );
  }

  Widget _buildTransfersTab() {
    final transfers = _bookings.where((b) => b['type'] == 'transfers').toList();
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: transfers.map((transfer) => _buildTransferBookingCard(transfer)).toList(),
        ),
      ),
    );
  }

  Widget _buildBookingCard(Map<String, dynamic> booking) {
    Color statusColor = _getStatusColor(booking['status']);
    return Padding(
      padding: const EdgeInsets.only(bottom: 16, top: 8),
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          Container(
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
                  // Background Image
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
                        // Top row: Routes and details
                        Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    '${booking['from']}-${booking['to']}',
                                    style: const TextStyle(
                                      fontSize: 15,
                                      fontWeight: FontWeight.w700,
                                      color: Colors.black87,
                                    ),
                                  ),
                                  const SizedBox(height: 2),
                                  Text(
                                    '${booking['departureTime']} - ${booking['arrivalTime']} (${booking['duration']})',
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
                                  booking['airline'],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black87,
                                  ),
                                ),
                                const SizedBox(height: 2),
                                Text(
                                  booking['cabin'],
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
                        Text(
                          booking['bags'],
                          style: const TextStyle(
                            fontSize: 12,
                            color: Colors.black87,
                          ),
                        ),

                        const SizedBox(height: 16),

                        // Fees section
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Cancellation Fee',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFD32F2F),
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                Text(
                                  booking['cancellationFee'],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Change Fee',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFD32F2F),
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                Text(
                                  booking['changeFee'],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),

                        const SizedBox(height: 12),

                        // Prices section
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Ticket Price',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFD32F2F),
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                Text(
                                  booking['ticketPrice'],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Upgrade Fee',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Color(0xFFD32F2F),
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                Text(
                                  booking['upgradeFee'],
                                  style: const TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),

                        const SizedBox(height: 12),

                        // Total Price
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            const Text(
                              'Total Price',
                              style: TextStyle(
                                fontSize: 12,
                                fontWeight: FontWeight.w600,
                                color: Colors.black87,
                              ),
                            ),
                            Text(
                              booking['totalPrice'],
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w700,
                                color: Colors.black87,
                              ),
                            ),
                          ],
                        ),

                        const SizedBox(height: 16),

                        // Buttons - Only for Upcoming status
                        if (booking['status'] == 'Upcoming')
                          Row(
                            children: [
                              Expanded(
                                child: ElevatedButton(
                                  onPressed: () => _showRefundDialog(context),
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: const Color(0xFFD32F2F),
                                    padding: const EdgeInsets.symmetric(vertical: 10),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(20),
                                    ),
                                  ),
                                  child: const Text(
                                    'Request refund',
                                    style: TextStyle(
                                      color: Colors.white,
                                      fontSize: 13,
                                      fontWeight: FontWeight.w600,
                                    ),
                                  ),
                                ),
                              ),
                              const SizedBox(width: 12),
                              Expanded(
                                child: OutlinedButton(
                                  onPressed: () => _showCancelConfirmation(context),
                                  style: OutlinedButton.styleFrom(
                                    side: const BorderSide(color: Color(0xFFD32F2F), width: 2),
                                    padding: const EdgeInsets.symmetric(vertical: 10),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(20),
                                    ),
                                  ),
                                  child: const Text(
                                    'Cancel',
                                    style: TextStyle(
                                      color: Color(0xFFD32F2F),
                                      fontSize: 13,
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
                ],
              ),
            ),
          ),

          // Status Badge - positioned on the border
          Positioned(
            top: -8,
            right: 16,
            child: Container(
              decoration: BoxDecoration(
                color: statusColor,
                borderRadius: BorderRadius.circular(6),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.2),
                    blurRadius: 4,
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              child: Text(
                booking['status'],
                style: const TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: Colors.white,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'Upcoming':
        return const Color(0xFF336891);
      case 'Cancelled':
        return const Color(0xFF9E9E9E);
      case 'Postponed':
        return const Color(0xFF9E9E9E);
      case 'Voided':
        return const Color(0xFF9E9E9E);
      default:
        return const Color(0xFF336891);
    }
  }

  Widget _buildHotelBookingCard(Map<String, dynamic> hotel) {
    Color statusColor = _getStatusColor(hotel['status']);
    return Padding(
      padding: const EdgeInsets.only(bottom: 16, top: 8),
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: const Color(0xFFE8F0F7),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        hotel['hotelName'],
                        style: const TextStyle(
                          color: Color(0xFFD32F2F),
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        hotel['location'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        hotel['amenities'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        hotel['bookedDate'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 12),
                      Text(
                        hotel['price'],
                        style: const TextStyle(
                          color: Color(0xFF2196F3),
                          fontSize: 16,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                      const SizedBox(height: 12),
                      // Buttons - Only for Upcoming status
                      if (hotel['status'] == 'Upcoming')
                        Row(
                          children: [
                            Expanded(
                              flex: 1,
                              child: ElevatedButton(
                                onPressed: () => _showRefundDialog(context),
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: const Color(0xFFD32F2F),
                                  padding: const EdgeInsets.symmetric(vertical: 8),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                ),
                                child: const Text(
                                  'Request refund',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                              ),
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              flex: 1,
                              child: OutlinedButton(
                                onPressed: () => _showCancelConfirmation(context),
                                style: OutlinedButton.styleFrom(
                                  side: const BorderSide(color: Color(0xFFD32F2F), width: 2),
                                  padding: const EdgeInsets.symmetric(vertical: 8),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                ),
                                child: const Text(
                                  'Cancel',
                                  style: TextStyle(
                                    color: Color(0xFFD32F2F),
                                    fontSize: 12,
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
                const SizedBox(width: 12),
                ClipRRect(
                  borderRadius: BorderRadius.circular(8),
                  child: Image.asset(
                    'assets/images/hotel2.jpg',
                    width: 100,
                    height: 80,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) {
                      return Container(
                        width: 100,
                        height: 80,
                        color: Colors.grey.shade300,
                        child: const Icon(Icons.hotel, size: 30),
                      );
                    },
                  ),
                ),
              ],
            ),
          ),
          // Status Badge - positioned on the border
          Positioned(
            top: -8,
            right: 12,
            child: Container(
              decoration: BoxDecoration(
                color: statusColor,
                borderRadius: BorderRadius.circular(6),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.2),
                    blurRadius: 4,
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              child: Text(
                hotel['status'],
                style: const TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: Colors.white,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCarBookingCard(Map<String, dynamic> car) {
    Color statusColor = _getStatusColor(car['status']);
    return Padding(
      padding: const EdgeInsets.only(bottom: 16, top: 8),
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: const Color(0xFFE8F0F7),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        car['category'],
                        style: const TextStyle(
                          color: Color(0xFFD32F2F),
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        car['name'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '${car['passengers']} - ${car['transmission']}',
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        car['mileage'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        car['shuttle'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 12),
                      // Buttons - Only for Upcoming status
                      if (car['status'] == 'Upcoming')
                        Row(
                          children: [
                            Expanded(
                              flex: 1,
                              child: ElevatedButton(
                                onPressed: () => _showRefundDialog(context),
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: const Color(0xFFD32F2F),
                                  padding: const EdgeInsets.symmetric(vertical: 8),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                ),
                                child: const Text(
                                  'Request refund',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                              ),
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              flex: 1,
                              child: OutlinedButton(
                                onPressed: () => _showCancelConfirmation(context),
                                style: OutlinedButton.styleFrom(
                                  side: const BorderSide(color: Color(0xFFD32F2F), width: 2),
                                  padding: const EdgeInsets.symmetric(vertical: 8),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                ),
                                child: const Text(
                                  'Cancel',
                                  style: TextStyle(
                                    color: Color(0xFFD32F2F),
                                    fontSize: 12,
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
                      car['price'],
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
          ),
          // Status Badge - positioned on the border
          Positioned(
            top: -8,
            right: 12,
            child: Container(
              decoration: BoxDecoration(
                color: statusColor,
                borderRadius: BorderRadius.circular(6),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.2),
                    blurRadius: 4,
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              child: Text(
                car['status'],
                style: const TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: Colors.white,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTransferBookingCard(Map<String, dynamic> transfer) {
    Color statusColor = _getStatusColor(transfer['status']);
    return Padding(
      padding: const EdgeInsets.only(bottom: 16, top: 8),
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: const Color(0xFFE8F0F7),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        transfer['category'],
                        style: const TextStyle(
                          color: Color(0xFFD32F2F),
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        transfer['name'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        transfer['passengers'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        transfer['estimatedTime'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        transfer['meetGreet'],
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 12),
                      // Buttons - Only for Upcoming status
                      if (transfer['status'] == 'Upcoming')
                        Row(
                          children: [
                            Expanded(
                              flex: 1,
                              child: ElevatedButton(
                                onPressed: () => _showRefundDialog(context),
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: const Color(0xFFD32F2F),
                                  padding: const EdgeInsets.symmetric(vertical: 8),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                ),
                                child: const Text(
                                  'Request refund',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                              ),
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              flex: 1,
                              child: OutlinedButton(
                                onPressed: () => _showCancelConfirmation(context),
                                style: OutlinedButton.styleFrom(
                                  side: const BorderSide(color: Color(0xFFD32F2F), width: 2),
                                  padding: const EdgeInsets.symmetric(vertical: 8),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                ),
                                child: const Text(
                                  'Cancel',
                                  style: TextStyle(
                                    color: Color(0xFFD32F2F),
                                    fontSize: 12,
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
                      transfer['price'],
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
          ),
          // Status Badge - positioned on the border
          Positioned(
            top: -8,
            right: 12,
            child: Container(
              decoration: BoxDecoration(
                color: statusColor,
                borderRadius: BorderRadius.circular(6),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.2),
                    blurRadius: 4,
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              child: Text(
                transfer['status'],
                style: const TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: Colors.white,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }


  void _showCancelConfirmation(BuildContext context) {
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
                  'Are you sure you want to cancel your ticket?',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 12),
                const Text(
                  'Canceling may cause additional fees',
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
                          'Back',
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
                          _showCancelSuccess(context);
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF336891),
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(20),
                          ),
                        ),
                        child: const Text(
                          'Cancel',
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

  void _showCancelSuccess(BuildContext context) {
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
                  'Ticket canceled successfully',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
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
                      'Back',
                      style: TextStyle(
                        color: Color(0xFF336891),
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

  void _showRefundDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          child: Stack(
            clipBehavior: Clip.none,
            children: [
              Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const SizedBox(height: 40),
                    
                    // Refund Icon
                    Image.asset(
                      'assets/images/refund.png',
                      height: 60,
                      width: 60,
                      errorBuilder: (context, error, stackTrace) {
                        return Container(
                          height: 60,
                          width: 60,
                          decoration: BoxDecoration(
                            color: const Color(0xFF336891),
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(
                            Icons.currency_exchange,
                            color: Colors.white,
                            size: 36,
                          ),
                        );
                      },
                    ),
                    
                    const SizedBox(height: 16),
                    
                    const Text(
                      'Request A Refund?',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: Colors.black87,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    
                    const SizedBox(height: 8),
                    
                    const Text(
                      'Your refund amount will be determined by our admin team based on the booking details and policy.',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w400,
                        color: Colors.black54,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    
                    const SizedBox(height: 24),
                    
                    // Request Button
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: () {
                          Navigator.pop(context);
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF336891),
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                        child: const Text(
                          'Request',
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
              
              // Close Button
              Positioned(
                top: -12,
                right: -12,
                child: GestureDetector(
                  onTap: () => Navigator.pop(context),
                  child: Container(
                    decoration: const BoxDecoration(
                      color: Color(0xFFD32F2F),
                      shape: BoxShape.circle,
                    ),
                    padding: const EdgeInsets.all(8),
                    child: const Icon(
                      Icons.close,
                      color: Colors.white,
                      size: 20,
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
