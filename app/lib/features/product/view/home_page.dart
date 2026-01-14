import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../app/routes/app_routes.dart';
import '../controller/product_controller.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage>
    with TickerProviderStateMixin, AutomaticKeepAliveClientMixin {
  late TabController _tabController;
  int _selectedTripType = 0; // 0: Roundtrip, 1: One way, 2: Multi-City
  
  // Hotel input controllers
  late TextEditingController _hotelDestinationController;
  late TextEditingController _checkinDateController;
  late TextEditingController _checkoutDateController;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _tabController.addListener(() {
      setState(() {});
    });
    _hotelDestinationController = TextEditingController();
    _checkinDateController = TextEditingController();
    _checkoutDateController = TextEditingController();
  }

  @override
  void dispose() {
    _tabController.dispose();
    _hotelDestinationController.dispose();
    _checkinDateController.dispose();
    _checkoutDateController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final size = MediaQuery.of(context).size;

    return DefaultTabController(
      length: 3,
      child: Scaffold(
        extendBodyBehindAppBar: true,
        backgroundColor: Colors.white,
        appBar: AppBar(
          backgroundColor: Colors.transparent,
          elevation: 0,
          leadingWidth: 0,
          leading: const SizedBox.shrink(),
          title: Padding(
            padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 16),
            child: Container(
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(25),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.15),
                    blurRadius: 12,
                    offset: const Offset(0, 4),
                    spreadRadius: 1,
                  ),
                ],
              ),
              child: TextField(
                decoration: InputDecoration(
                  hintText: 'Search flights...',
                  border: InputBorder.none,
                  prefixIcon: Padding(
                    padding: const EdgeInsets.all(12),
                    child: Image.asset(
                      'assets/images/logo-search.png',
                      width: 20,
                      height: 20,
                      fit: BoxFit.contain,
                    ),
                  ),
                  contentPadding: const EdgeInsets.symmetric(
                    vertical: 12,
                    horizontal: 16,
                  ),
                  suffixIcon: const Icon(
                    Icons.close,
                    color: Colors.grey,
                    size: 20,
                  ),
                ),
              ),
            ),
          ),
          titleSpacing: 0,
          actions: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8),
              child: Center(
                child: GestureDetector(
                  onTap: () {},
                  child: const Icon(
                    Icons.notifications,
                    color: Color(0xFF1e5a8e),
                    size: 24,
                  ),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.only(right: 16),
              child: Center(
                child: GestureDetector(
                  onTap: () => Get.toNamed(AppRoutes.ACCOUNT),
                  child: CircleAvatar(
                    radius: 18,
                    backgroundColor: Colors.grey[300],
                    child: const Icon(
                      Icons.person,
                      color: Color(0xFF1e5a8e),
                      size: 20,
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
        body: Column(
          children: [
            SizedBox(height: size.height * 0.12),

            // Tab Bar (custom pills)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: Row(
                children: [
                  Expanded(
                    child: Container(
                      decoration: BoxDecoration(
                        color: _tabController.index == 0
                            ? const Color(0xFF1e5a8e)
                            : Colors.white,
                        borderRadius: BorderRadius.circular(25),
                        border: _tabController.index != 0
                            ? Border.all(
                                color: Colors.grey[400]!,
                                width: 1,
                              )
                            : null,
                        boxShadow: [
                          BoxShadow(
                            color: Colors.black.withOpacity(0.1),
                            blurRadius: 8,
                            offset: const Offset(0, 2),
                          ),
                        ],
                      ),
                      child: Material(
                        color: Colors.transparent,
                        child: InkWell(
                          onTap: () {
                            _tabController.animateTo(0);
                            setState(() {});
                          },
                          borderRadius: BorderRadius.circular(25),
                          child: Padding(
                            padding:
                                const EdgeInsets.symmetric(vertical: 12),
                            child: Center(
                              child: Text(
                                'flights'.tr,
                                style: TextStyle(
                                  fontSize: 13,
                                  fontWeight: FontWeight.w600,
                                  color: _tabController.index == 0
                                      ? Colors.white
                                      : Colors.grey[700],
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Container(
                      decoration: BoxDecoration(
                        color: _tabController.index == 1
                            ? const Color(0xFF1e5a8e)
                            : Colors.white,
                        borderRadius: BorderRadius.circular(25),
                        border: _tabController.index != 1
                            ? Border.all(
                                color: Colors.grey[400]!,
                                width: 1,
                              )
                            : null,
                        boxShadow: [
                          BoxShadow(
                            color: Colors.black.withOpacity(0.1),
                            blurRadius: 8,
                            offset: const Offset(0, 2),
                          ),
                        ],
                      ),
                      child: Material(
                        color: Colors.transparent,
                        child: InkWell(
                          onTap: () {
                            _tabController.animateTo(1);
                            setState(() {});
                          },
                          borderRadius: BorderRadius.circular(25),
                          child: Padding(
                            padding:
                                const EdgeInsets.symmetric(vertical: 12),
                            child: Center(
                              child: Text(
                                'cars'.tr,
                                style: TextStyle(
                                  fontSize: 13,
                                  fontWeight: FontWeight.w600,
                                  color: _tabController.index == 1
                                      ? Colors.white
                                      : Colors.grey[700],
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Container(
                      decoration: BoxDecoration(
                        color: _tabController.index == 2
                            ? const Color(0xFF1e5a8e)
                            : Colors.white,
                        borderRadius: BorderRadius.circular(25),
                        border: _tabController.index != 2
                            ? Border.all(
                                color: Colors.grey[400]!,
                                width: 1,
                              )
                            : null,
                        boxShadow: [
                          BoxShadow(
                            color: Colors.black.withOpacity(0.1),
                            blurRadius: 8,
                            offset: const Offset(0, 2),
                          ),
                        ],
                      ),
                      child: Material(
                        color: Colors.transparent,
                        child: InkWell(
                          onTap: () {
                            _tabController.animateTo(2);
                            setState(() {});
                          },
                          borderRadius: BorderRadius.circular(25),
                          child: Padding(
                            padding:
                                const EdgeInsets.symmetric(vertical: 12),
                            child: Center(
                              child: Text(
                                'hotels'.tr,
                                style: TextStyle(
                                  fontSize: 13,
                                  fontWeight: FontWeight.w600,
                                  color: _tabController.index == 2
                                      ? Colors.white
                                      : Colors.grey[700],
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 20),

            // Tab Views
            Expanded(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24),
                child: TabBarView(
                  controller: _tabController,
                  children: [
                    // Flights Tab (scrollable)
                    SingleChildScrollView(
                      key: const PageStorageKey('flights_tab'),
                      physics: const AlwaysScrollableScrollPhysics(),
                      child: Padding(
                        padding: const EdgeInsets.only(bottom: 24),
                        child: _buildFlightsTab(context, size),
                      ),
                    ),

                    // Cars Tab (centered card, no stretching)
                    _buildCarsTab(),

                    // Hotels Tab (centered card, no stretching)
                    _buildHotelsTab(),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  // ----------------- FLIGHTS TAB -----------------

  Widget _buildFlightsTab(BuildContext context, Size size) {
    return Column(
      children: [
        // Trip Type Selection
        Container(
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(15),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 8,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              Expanded(
                child: GestureDetector(
                  onTap: () => setState(() => _selectedTripType = 0),
                  child: _buildTripTypeButton(
                    'Roundtrip',
                    _selectedTripType == 0,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: GestureDetector(
                  onTap: () => setState(() => _selectedTripType = 1),
                  child: _buildTripTypeButton(
                    'One way',
                    _selectedTripType == 1,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: GestureDetector(
                  onTap: () => setState(() => _selectedTripType = 2),
                  child: _buildTripTypeButton(
                    'Multi-City',
                    _selectedTripType == 2,
                  ),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),

        // From - To Row (side by side)
        Row(
          children: [
            // From Field
            Expanded(
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(15),
                  border: Border.all(color: const Color(0xFF1e5a8e)),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 8,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'From',
                      style: TextStyle(
                        fontSize: 11,
                        color: const Color(0xFFD32F2F),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        const Icon(Icons.location_on, color: Color(0xFFD32F2F), size: 20),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Text(
                            'Bey-Lebanon',
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(width: 12),
            // To Field
            Expanded(
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(15),
                  border: Border.all(color: const Color(0xFF1e5a8e)),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 8,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'To',
                      style: TextStyle(
                        fontSize: 11,
                        color: const Color(0xFFD32F2F),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        const Icon(Icons.location_on, color: Color(0xFFD32F2F), size: 20),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Text(
                            'Bey-Lebanon',
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),

        // Depart - Return
        Container(
          width: double.infinity,
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(15),
            border: Border.all(color: const Color(0xFF1e5a8e)),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 8,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Depart - Return',
                    style: TextStyle(
                      fontSize: 11,
                      color: const Color(0xFFD32F2F),
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      const Icon(Icons.calendar_today, color: Color(0xFFD32F2F), size: 16),
                      const SizedBox(width: 8),
                      Text(
                        '07 Nov 22 - 13 Nov 22',
                        style: const TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w600,
                          color: Colors.black,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              Icon(Icons.edit, color: Colors.grey[400], size: 20),
            ],
          ),
        ),
        const SizedBox(height: 16),

        // Number of Travelers
        Container(
          width: double.infinity,
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(15),
            border: Border.all(color: const Color(0xFF1e5a8e)),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 8,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Number of Travelers',
                style: TextStyle(
                  fontSize: 11,
                  color: const Color(0xFFD32F2F),
                  fontWeight: FontWeight.w500,
                ),
              ),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    children: const [
                      Icon(Icons.person, color: Color(0xFFD32F2F), size: 20),
                      SizedBox(width: 8),
                      Text(
                        '1 Adult, 2 children, 1 infant',
                        style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w600,
                          color: Colors.black,
                        ),
                      ),
                    ],
                  ),
                  Text(
                    'Maximum: 9 passengers',
                    style: TextStyle(
                      fontSize: 11,
                      color: Colors.grey[500],
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),

        // Class Selection
        Container(
          width: double.infinity,
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(15),
            border: Border.all(color: const Color(0xFF1e5a8e)),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 8,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Class',
                style: TextStyle(
                  fontSize: 11,
                  color: const Color(0xFFD32F2F),
                  fontWeight: FontWeight.w500,
                ),
              ),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    'Business',
                    style: TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: Colors.black,
                    ),
                  ),
                  Icon(Icons.expand_more, color: const Color(0xFF1e5a8e), size: 24),
                ],
              ),
            ],
          ),
        ),
        const SizedBox(height: 24),

        // Search Button
        SizedBox(
          width: double.infinity,
          height: 50,
          child: ElevatedButton(
            onPressed: () {},
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF1e5a8e),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(25),
              ),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(Icons.search, color: Colors.white, size: 20),
                const SizedBox(width: 8),
                Text(
                  'search'.tr,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildTripTypeButton(String label, bool isSelected) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.transparent,
        border: isSelected
            ? const Border(
                bottom: BorderSide(
                  color: Color(0xFF1e5a8e),
                  width: 3,
                ),
              )
            : null,
      ),
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Center(
        child: Text(
          label,
          style: TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            color: isSelected
                ? const Color(0xFF1e5a8e)
                : Colors.black87,
          ),
        ),
      ),
    );
  }

  Widget _buildFlightCard({
    required String title,
    required String fromLocation,
    required String toLocation,
    required String departDate,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(15),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Title Row
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                title,
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
              const Text(
                'Lebaq',
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: Color(0xFF1e5a8e),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  // ----------------- CARS & HOTELS TABS -----------------

  Widget _buildCarsTab() {
    return Center(
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(15),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 8,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        padding: const EdgeInsets.all(32),
        child: Text(
          'Cars feature coming soon',
          style: TextStyle(
            fontSize: 16,
            color: Colors.grey[600],
            fontWeight: FontWeight.w500,
          ),
        ),
      ),
    );
  }

  Widget _buildHotelsTab() {
    return SingleChildScrollView(
      physics: const AlwaysScrollableScrollPhysics(),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Hotel Destination
            Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(15),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 8,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
              child: Row(
                children: [
                  Icon(Icons.location_on, color: const Color(0xFFD32F2F), size: 24),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextField(
                      controller: _hotelDestinationController,
                      decoration: InputDecoration(
                        hintText: 'Going to',
                        hintStyle: TextStyle(color: Colors.grey[500], fontSize: 16),
                        border: InputBorder.none,
                        isDense: true,
                      ),
                      style: const TextStyle(fontSize: 16, color: Colors.black),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),

            // Check-in Date
            Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(15),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 8,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
              child: Row(
                children: [
                  Icon(Icons.calendar_today, color: const Color(0xFFD32F2F), size: 20),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextField(
                      controller: _checkinDateController,
                      decoration: InputDecoration(
                        hintText: 'Check-in Date',
                        hintStyle: TextStyle(color: Colors.grey[500], fontSize: 14),
                        border: InputBorder.none,
                        isDense: true,
                      ),
                      style: const TextStyle(fontSize: 14, color: Colors.black),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),

            // Check-out Date
            Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(15),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 8,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
              child: Row(
                children: [
                  Icon(Icons.calendar_today, color: const Color(0xFFD32F2F), size: 20),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextField(
                      controller: _checkoutDateController,
                      decoration: InputDecoration(
                        hintText: 'Check-out Date',
                        hintStyle: TextStyle(color: Colors.grey[500], fontSize: 14),
                        border: InputBorder.none,
                        isDense: true,
                      ),
                      style: const TextStyle(fontSize: 14, color: Colors.black),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),

            // Number of Rooms and Guests
            Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(15),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 8,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
              child: Row(
                children: [
                  Icon(Icons.people, color: const Color(0xFFD32F2F), size: 24),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '2 rooms, 4 guests',
                          style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w500, color: Colors.black),
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          '1 room, 2 guests, 1 child',
                          style: TextStyle(fontSize: 11, color: Colors.grey[600]),
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),

            // Hotel Class
            Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(15),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 8,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
              child: Row(
                children: [
                  Icon(Icons.star, color: const Color(0xFFD32F2F), size: 24),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Text(
                      '5 Star',
                      style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Colors.black),
                    ),
                  ),
                  Icon(Icons.check_circle, color: const Color(0xFF1e5a8e), size: 24),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // Search Button
            SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton(
                onPressed: () {},
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1e5a8e),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.search, color: Colors.white),
                    const SizedBox(width: 8),
                    Text(
                      'search'.tr,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 32),

            // Get the best travel experience
            Text(
              'Get the best travel experience',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: Colors.grey[800],
              ),
            ),
            const SizedBox(height: 16),

            // Popular Hotel Cards
            _buildHotelDestinationCard('Paris, France', 'assets/images/paris.jpg'),
            const SizedBox(height: 12),
            _buildHotelDestinationCard('Dubai, UAE', 'assets/images/dubai.jpg'),
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildHotelDestinationCard(String location, String imagePath) {
    return Container(
      width: double.infinity,
      height: 200,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(15),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(15),
        child: Stack(
          children: [
            Container(
              color: Colors.grey[300],
              child: const Icon(Icons.image, size: 60, color: Colors.grey),
            ),
            Positioned(
              bottom: 0,
              left: 0,
              right: 0,
              child: Container(
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.bottomCenter,
                    end: Alignment.topCenter,
                    colors: [
                      Colors.black.withOpacity(0.7),
                      Colors.transparent,
                    ],
                  ),
                ),
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
                child: Row(
                  children: [
                    Expanded(
                      child: Text(
                        location,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 15,
                          fontWeight: FontWeight.bold,
                        ),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Container(
                      decoration: BoxDecoration(
                        color: const Color(0xFF1e5a8e),
                        borderRadius: BorderRadius.circular(6),
                      ),
                      padding: const EdgeInsets.all(6),
                      child: const Icon(Icons.arrow_forward, color: Colors.white, size: 16),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
