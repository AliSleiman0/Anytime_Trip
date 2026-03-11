import 'dart:async';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../app/routes/app_routes.dart';
import '../../../core/widgets/global_chatbot_overlay.dart';
import '../../../core/network/endpoints.dart';
import '../controller/product_controller.dart';
import '../service/banner_service.dart';
import '../service/popular_location_service.dart';
import '../model/banner_model.dart';
import '../model/popular_location_model.dart';
import '../../account/view/notifications_page.dart';
import 'cars_tab.dart';
import 'hotels_tab.dart';
import 'transfers_tab.dart';
import 'flight_search_results.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage>
    with TickerProviderStateMixin, AutomaticKeepAliveClientMixin {
  late TabController _tabController;
  int _selectedTripType = 0; // 0: Roundtrip, 1: One way, 2: Multi-City
  
  // Carousel for destination images
  late PageController _carouselController;
  int _currentCarouselIndex = 0;
  Timer? _carouselTimer;
  final BannerService _bannerService = BannerService();
  final PopularLocationService _popularLocationService = PopularLocationService();
  List<BannerModel> _banners = [];
  bool _bannersLoading = true;
  
  // Homepage banners (ads)
  List<BannerModel> _homepageBanners = [];
  bool _homepageBannersLoading = true;
  
  // Popular locations
  List<PopularLocationModel> _popularLocations = [];
  bool _popularLocationsLoading = true;
  
  // Fallback images if no banners from backend
  final List<String> _fallbackCarouselImages = [
    'assets/images/barcelona.jpg',
    'assets/images/paris.jpg',
    'assets/images/london.jpg',
  ];
  final List<String> _fallbackCarouselTitles = [
    'Barcelona, Spain',
    'Paris, France',
    'London, UK',
  ];
  
  // Hotel input controllers
  late TextEditingController _hotelDestinationController;
  late TextEditingController _checkinDateController;
  late TextEditingController _checkoutDateController;
  
  // Flight location controllers
  String _selectedFromLocation = 'Beirut, Lebanon';
  String _selectedToLocation = 'Atlanta, USA';
  
  // Flight date controllers
  DateTime? _departDate;
  DateTime? _returnDate;
  
  // Passenger counts
  int _adultsCount = 1;
  int _childrenCount = 0;
  int _infantsCount = 0;
  
  // Flight class
  String _selectedClass = 'Economy';
  final List<String> _flightClasses = ['All', 'Economy', 'Business', 'First-Class'];
  
  // Multi-city flights
  List<Map<String, dynamic>> _multiCityFlights = [
    {
      'from': 'Beirut, Lebanon',
      'to': 'Beirut, Lebanon',
      'date': null,
    },
    {
      'from': 'Beirut, Lebanon',
      'to': 'Beirut, Lebanon',
      'date': null,
    },
  ];
  
  String get _passengerSummary {
    List<String> parts = [];
    if (_adultsCount > 0) parts.add('$_adultsCount Adult${_adultsCount > 1 ? 's' : ''}');
    if (_childrenCount > 0) parts.add('$_childrenCount Child${_childrenCount > 1 ? 'ren' : ''}');
    if (_infantsCount > 0) parts.add('$_infantsCount Infant${_infantsCount > 1 ? 's' : ''}');
    return parts.isEmpty ? 'Select passengers' : parts.join(', ');
  }
  
  int get _totalPassengers => _adultsCount + _childrenCount + _infantsCount;
  
  String get _formattedDates {
    if (_departDate == null) {
      return _selectedTripType == 1 ? '07 Nov 22' : '07 Nov 22 - 13 Nov 22';
    }
    final departStr = '${_departDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_departDate!.month)} ${_departDate!.year.toString().substring(2)}';
    
    if (_selectedTripType == 1) {
      // One way - only show depart date
      return departStr;
    }
    
    // Roundtrip - show both dates
    final returnStr = _returnDate != null
        ? '${_returnDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_returnDate!.month)} ${_returnDate!.year.toString().substring(2)}'
        : 'Select';
    return '$departStr - $returnStr';
  }
  
  String _getMonthAbbr(int month) {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return months[month - 1];
  }
  
  final List<Map<String, String>> _destinations = [
    {'name': 'Beirut, Lebanon', 'airport': 'Rafic Hariri Intl.'},
    {'name': 'Atlanta, USA', 'airport': 'Hartsfield–Jackson Atlanta Intl.'},
    {'name': 'Madrid, Spain', 'airport': 'Madrid-Barajas Airport'},
    {'name': 'New York, USA', 'airport': 'JFK International'},
    {'name': 'Dubai, UAE', 'airport': 'Dubai International'},
    {'name': 'London, UK', 'airport': 'Heathrow Airport'},
    {'name': 'Paris, France', 'airport': 'Charles de Gaulle Airport'},
    {'name': 'Tokyo, Japan', 'airport': 'Narita International'},
  ];

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    _tabController.addListener(() {
      setState(() {});
    });
    _hotelDestinationController = TextEditingController();
    _checkinDateController = TextEditingController();
    _checkoutDateController = TextEditingController();
    
    // Initialize carousel
    _carouselController = PageController();
    _loadBanners();
    _loadHomepageBanners();
    _loadPopularLocations();
    _startCarouselTimer();
  }
  
  void _startCarouselTimer() {
    _carouselTimer?.cancel();
    _carouselTimer = Timer.periodic(const Duration(seconds: 3), (timer) {
      final imageCount = _banners.isNotEmpty ? _banners.length : _fallbackCarouselImages.length;
      if (_currentCarouselIndex < imageCount - 1) {
        _currentCarouselIndex++;
      } else {
        _currentCarouselIndex = 0;
      }
      if (_carouselController.hasClients) {
        _carouselController.animateToPage(
          _currentCarouselIndex,
          duration: const Duration(milliseconds: 350),
          curve: Curves.easeInOut,
        );
      }
    });
  }

  Future<void> _loadBanners() async {
    try {
      final banners = await _bannerService.getBanners();
      if (mounted) {
        setState(() {
          _banners = banners;
          _bannersLoading = false;
        });
      }
    } catch (e) {
      print('[HOME_PAGE] Failed to load banners: $e');
      if (mounted) {
        setState(() {
          _bannersLoading = false;
        });
      }
    }
  }

  Future<void> _loadHomepageBanners() async {
    try {
      final banners = await _bannerService.getHomepageBanners();
      if (mounted) {
        setState(() {
          _homepageBanners = banners;
          _homepageBannersLoading = false;
        });
      }
    } catch (e) {
      print('[HOME_PAGE] Failed to load homepage banners: $e');
      if (mounted) {
        setState(() {
          _homepageBannersLoading = false;
        });
      }
    }
  }

  Future<void> _loadPopularLocations() async {
    try {
      final locations = await _popularLocationService.getPopularLocations();
      if (mounted) {
        setState(() {
          _popularLocations = locations;
          _popularLocationsLoading = false;
        });
      }
    } catch (e) {
      print('[HOME_PAGE] Failed to load popular locations: $e');
      if (mounted) {
        setState(() {
          _popularLocationsLoading = false;
        });
      }
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    _hotelDestinationController.dispose();
    _checkinDateController.dispose();
    _checkoutDateController.dispose();
    _carouselController.dispose();
    _carouselTimer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final size = MediaQuery.of(context).size;

    return DefaultTabController(
      length: 3,
      child: Scaffold(
        extendBodyBehindAppBar: false,
        backgroundColor: Colors.white,
        appBar: AppBar(
          backgroundColor: const Color(0xFF1e5a8e),
          elevation: 0,
          leadingWidth: 60,
          leading: Padding(
            padding: const EdgeInsets.only(left: 16),
            child: Image.asset(
              'assets/images/main-logo.png',
              width: 40,
              height: 40,
              fit: BoxFit.contain,
            ),
          ),
          titleSpacing: 0,
          actions: [
            Expanded(
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    GestureDetector(
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(builder: (context) => const NotificationsPage()),
                        );
                      },
                      child: const Icon(
                        Icons.notifications,
                        color: Colors.white,
                        size: 24,
                      ),
                    ),
                    const SizedBox(width: 16),
                    GestureDetector(
                      onTap: () => Get.toNamed(AppRoutes.ACCOUNT),
                      child: const CircleAvatar(
                        radius: 18,
                        backgroundColor: Colors.white24,
                        child: Icon(
                          Icons.person,
                          color: Colors.white,
                          size: 20,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
        body: Column(
          children: [
            SizedBox(height: size.height * 0.02),

            // Tab Bar (custom pills) - horizontally scrollable
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: [
                    _buildTabButton(0, 'flights'.tr),
                    const SizedBox(width: 12),
                    _buildTabButton(1, 'cars'.tr),
                    const SizedBox(width: 12),
                    _buildTabButton(2, 'hotels'.tr),
                    const SizedBox(width: 12),
                    _buildTabButton(3, 'Transfers'),
                  ],
                ),
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

                    // Cars Tab
                    CarsTab(
                      popularLocations: _popularLocations,
                      popularLocationsLoading: _popularLocationsLoading,
                    ),

                    // Hotels Tab
                    HotelsTab(
                      popularLocations: _popularLocations,
                      popularLocationsLoading: _popularLocationsLoading,
                    ),

                    // Transfers Tab
                    TransfersTab(
                      popularLocations: _popularLocations,
                      popularLocationsLoading: _popularLocationsLoading,
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

  // ----------------- FLIGHTS TAB -----------------

  Widget _buildFlightsTab(BuildContext context, Size size) {
    // Show multi-city UI if trip type is 2
    if (_selectedTripType == 2) {
      return _buildMultiCityFlights(context, size);
    }
    
    // Regular roundtrip/one-way UI
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
              child: GestureDetector(
                onTap: () => _showLocationBottomSheet(context, isFrom: true),
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
                              _selectedFromLocation,
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Colors.black,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            // To Field
            Expanded(
              child: GestureDetector(
                onTap: () => _showLocationBottomSheet(context, isFrom: false),
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
                              _selectedToLocation,
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Colors.black,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),

        // Depart - Return
        GestureDetector(
          onTap: () => _showDateBottomSheet(context),
          child: Container(
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
                      _selectedTripType == 1 ? 'Depart' : 'Depart - Return',
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
                          _formattedDates,
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
                const Icon(Icons.calendar_month, color: Color(0xFFD32F2F), size: 24),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),

        // Number of Travelers
        GestureDetector(
          onTap: () => _showPassengersBottomSheet(context),
          child: Container(
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
                    Expanded(
                      child: Row(
                        children: [
                          const Icon(Icons.person, color: Color(0xFFD32F2F), size: 20),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              _passengerSummary,
                              style: const TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: Colors.black,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    ),
                    Text(
                      'Max: 9',
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
        ),
        const SizedBox(height: 16),

        // Class Selection
        GestureDetector(
          onTap: () => _showClassBottomSheet(context),
          child: Container(
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
                    Text(
                      _selectedClass,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: Colors.black,
                      ),
                    ),
                    const Icon(Icons.expand_more, color: Color(0xFF1e5a8e), size: 24),
                  ],
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 24),

        // Search Button
        SizedBox(
          width: double.infinity,
          height: 50,
          child: ElevatedButton(
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => FlightSearchResults(
                    from: _selectedFromLocation,
                    to: _selectedToLocation,
                    departureDate: _departDate != null
                        ? '${_departDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_departDate!.month)} ${_departDate!.year.toString().substring(2)}'
                        : '07 Nov 22',
                    returnDate: _selectedTripType == 1
                        ? null
                        : (_returnDate != null
                            ? '${_returnDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_returnDate!.month)} ${_returnDate!.year.toString().substring(2)}'
                            : null),
                    travelers: _totalPassengers,
                    travelClass: _selectedClass,
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
        const SizedBox(height: 32),

        // Carousel Section
        _buildDestinationCarousel(),
        const SizedBox(height: 24),

        // Advertisement
        _buildAdvertisement(),
        const SizedBox(height: 32),

        // Popular Locations
        _buildPopularLocations(),
      ],
    );
  }

  Widget _buildMultiCityFlights(BuildContext context, Size size) {
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

        // Number of Travelers
        GestureDetector(
          onTap: () => _showPassengersBottomSheet(context),
          child: Container(
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
                    Expanded(
                      child: Row(
                        children: [
                          const Icon(Icons.person, color: Color(0xFFD32F2F), size: 20),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              _passengerSummary,
                              style: const TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: Colors.black,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    ),
                    Text(
                      'Max: 9',
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
        ),
        const SizedBox(height: 16),

        // Dynamic flight sections
        ...List.generate(_multiCityFlights.length, (index) {
          return Column(
            children: [
              _buildFlightSection(context, index),
              const SizedBox(height: 16),
            ],
          );
        }),

        // Add Flight Button
        GestureDetector(
          onTap: () {
            if (_multiCityFlights.length < 5) {
              setState(() {
                _multiCityFlights.add({
                  'from': 'Beirut, Lebanon',
                  'to': 'Beirut, Lebanon',
                  'date': null,
                });
              });
            }
          },
          child: Text(
            '+Add Flight',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: const Color(0xFF1e5a8e),
            ),
          ),
        ),
        const SizedBox(height: 24),

        // Search Button
        SizedBox(
          width: double.infinity,
          height: 50,
          child: ElevatedButton(
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => FlightSearchResults(
                    from: _multiCityFlights[0]['from'] ?? 'Beirut, Lebanon',
                    to: _multiCityFlights[_multiCityFlights.length - 1]['to'] ?? 'Atlanta, USA',
                    departureDate: '07 Nov 22 - 13 Nov 22',
                    returnDate: null,
                    travelers: _totalPassengers,
                    travelClass: _selectedClass,
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
        const SizedBox(height: 32),

        // Carousel Section
        _buildDestinationCarousel(),
        const SizedBox(height: 24),

        // Advertisement
        _buildAdvertisement(),
        const SizedBox(height: 32),

        // Popular Locations
        _buildPopularLocations(),
      ],
    );
  }

  Widget _buildFlightSection(BuildContext context, int index) {
    final flight = _multiCityFlights[index];
    final isLastFlight = index == _multiCityFlights.length - 1;
    
    String getOrdinal(int num) {
      if (num == 0) return 'First';
      if (num == 1) return 'Second';
      if (num == 2) return 'Third';
      if (num == 3) return 'Fourth';
      return '${num + 1}th';
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Flight header
        Text(
          '${getOrdinal(index)} Flight',
          style: const TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w700,
            color: Colors.black,
          ),
        ),
        const SizedBox(height: 12),

        // From - To Row
        Row(
          children: [
            // From Field
            Expanded(
              child: GestureDetector(
                onTap: () => _showMultiCityLocationBottomSheet(context, index, isFrom: true),
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
                              flight['from'],
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Colors.black,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            // To Field
            Expanded(
              child: GestureDetector(
                onTap: () => _showMultiCityLocationBottomSheet(context, index, isFrom: false),
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
                              flight['to'],
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Colors.black,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),

        // Depart Date
        GestureDetector(
          onTap: () => _showMultiCityDateBottomSheet(context, index),
          child: Container(
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
                      'Depart',
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
                          flight['date'] != null
                              ? '${flight['date'].day.toString().padLeft(2, '0')} ${_getMonthAbbr(flight['date'].month)} ${flight['date'].year.toString().substring(2)}'
                              : '07 Nov 22 - 13 Nov 22',
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
                const Icon(Icons.calendar_month, color: Color(0xFFD32F2F), size: 24),
              ],
            ),
          ),
        ),

        // Class Selection (only for last flight)
        if (isLastFlight) ...[
          const SizedBox(height: 16),
          GestureDetector(
            onTap: () => _showClassBottomSheet(context),
            child: Container(
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
                      Text(
                        _selectedClass,
                        style: const TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                          color: Colors.black,
                        ),
                      ),
                      const Icon(Icons.expand_more, color: Color(0xFF1e5a8e), size: 24),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
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

  // ----------------- CAROUSEL, AD & POPULAR LOCATIONS -----------------

  void _showLocationBottomSheet(BuildContext context, {required bool isFrom}) {
    final TextEditingController searchController = TextEditingController();
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            final filteredDestinations = _destinations.where((destination) {
              final name = destination['name']!.toLowerCase();
              final airport = destination['airport']!.toLowerCase();
              final query = searchController.text.toLowerCase();
              return name.contains(query) || airport.contains(query);
            }).toList();

            return DraggableScrollableSheet(
              initialChildSize: 0.6,
              minChildSize: 0.4,
              maxChildSize: 0.95,
              builder: (BuildContext context, ScrollController scrollController) {
                return Container(
                  decoration: const BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(20),
                      topRight: Radius.circular(20),
                    ),
                  ),
                  child: Column(
                    children: [
                      // Handle bar
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        child: Container(
                          width: 40,
                          height: 4,
                          decoration: BoxDecoration(
                            color: Colors.grey[400],
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                      ),
                      // Header
                      Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 20,
                          vertical: 12,
                        ),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              isFrom ? 'Select Departure' : 'Select Destination',
                              style: const TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.w700,
                                color: Colors.black,
                              ),
                            ),
                            GestureDetector(
                              onTap: () {
                                isModalOpenNotifier.value = false;
                                Navigator.pop(context);
                              },
                              child: Container(
                                decoration: BoxDecoration(
                                  color: const Color(0xFFD32F2F),
                                  borderRadius: BorderRadius.circular(50),
                                ),
                                padding: const EdgeInsets.all(4),
                                child: const Icon(
                                  Icons.close,
                                  color: Colors.white,
                                  size: 20,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      // Search field
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20),
                        child: TextField(
                          controller: searchController,
                          autofocus: false,
                          onChanged: (value) {
                            setModalState(() {});
                          },
                          decoration: InputDecoration(
                            hintText: 'Search destinations...',
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(color: Colors.grey[300]!),
                            ),
                            enabledBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(color: Colors.grey[300]!),
                            ),
                            focusedBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(
                                color: Color(0xFF1e5a8e),
                                width: 2,
                              ),
                            ),
                            prefixIcon: const Icon(
                              Icons.search,
                              color: Colors.grey,
                            ),
                            suffixIcon: searchController.text.isNotEmpty
                                ? GestureDetector(
                                    onTap: () {
                                      searchController.clear();
                                      setModalState(() {});
                                    },
                                    child: const Icon(
                                      Icons.clear,
                                      color: Colors.grey,
                                    ),
                                  )
                                : const Icon(
                                    Icons.location_on,
                                    color: Color(0xFFD32F2F),
                                  ),
                            contentPadding: const EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 12,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),
                      // Destinations list
                      Expanded(
                        child: filteredDestinations.isEmpty
                            ? Center(
                                child: Padding(
                                  padding: const EdgeInsets.all(24),
                                  child: Column(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Icon(
                                        Icons.search_off,
                                        size: 64,
                                        color: Colors.grey[400],
                                      ),
                                      const SizedBox(height: 16),
                                      Text(
                                        'No destinations found',
                                        style: TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.w600,
                                          color: Colors.grey[600],
                                        ),
                                      ),
                                      const SizedBox(height: 8),
                                      Text(
                                        'Try searching with different keywords',
                                        style: TextStyle(
                                          fontSize: 13,
                                          color: Colors.grey[500],
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              )
                            : ListView.builder(
                                controller: scrollController,
                                padding: const EdgeInsets.symmetric(horizontal: 16),
                                itemCount: filteredDestinations.length,
                                itemBuilder: (context, index) {
                                  final destination = filteredDestinations[index];
                                  final isSelected = isFrom
                                      ? _selectedFromLocation == destination['name']
                                      : _selectedToLocation == destination['name'];

                                  return GestureDetector(
                                    onTap: () {
                                      setState(() {
                                        if (isFrom) {
                                          _selectedFromLocation = destination['name']!;
                                        } else {
                                          _selectedToLocation = destination['name']!;
                                        }
                                      });
                                      isModalOpenNotifier.value = false;
                                      Navigator.pop(context);
                                    },
                                    child: Container(
                                      margin: const EdgeInsets.only(bottom: 12),
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 16,
                                        vertical: 12,
                                      ),
                                      decoration: BoxDecoration(
                                        color: isSelected
                                            ? const Color(0xFF1e5a8e).withOpacity(0.1)
                                            : Colors.white,
                                        borderRadius: BorderRadius.circular(12),
                                        border: Border.all(
                                          color: isSelected
                                              ? const Color(0xFF1e5a8e)
                                              : Colors.grey[200]!,
                                          width: isSelected ? 2 : 1,
                                        ),
                                      ),
                                      child: Row(
                                        children: [
                                          Container(
                                            width: 40,
                                            height: 40,
                                            decoration: BoxDecoration(
                                              color: const Color(0xFFD32F2F)
                                                  .withOpacity(0.1),
                                              borderRadius: BorderRadius.circular(8),
                                            ),
                                            child: const Icon(
                                              Icons.location_on,
                                              color: Color(0xFFD32F2F),
                                              size: 20,
                                            ),
                                          ),
                                          const SizedBox(width: 12),
                                          Expanded(
                                            child: Column(
                                              crossAxisAlignment:
                                                  CrossAxisAlignment.start,
                                              children: [
                                                Text(
                                                  destination['name']!,
                                                  style: TextStyle(
                                                    fontSize: 14,
                                                    fontWeight: FontWeight.w600,
                                                    color: isSelected
                                                        ? const Color(0xFF1e5a8e)
                                                        : Colors.black,
                                                  ),
                                                ),
                                                const SizedBox(height: 4),
                                                Text(
                                                  destination['airport']!,
                                                  style: TextStyle(
                                                    fontSize: 12,
                                                    color: Colors.grey[600],
                                                  ),
                                                ),
                                              ],
                                            ),
                                          ),
                                          if (isSelected)
                                            const Icon(
                                              Icons.check_circle,
                                              color: Color(0xFF1e5a8e),
                                              size: 24,
                                            ),
                                        ],
                                      ),
                                    ),
                                  );
                                },
                              ),
                      ),
                      const SizedBox(height: 24),
                    ],
                  ),
                );
              },
            );
          },
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  void _showDateBottomSheet(BuildContext context) {
    DateTime? tempDepartDate = _departDate;
    DateTime? tempReturnDate = _returnDate;
    final bool isOneWay = _selectedTripType == 1;
    DateTime currentMonth = tempDepartDate ?? DateTime.now();
    bool isSelectingReturn = false;
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return DraggableScrollableSheet(
              initialChildSize: 0.7,
              minChildSize: 0.5,
              maxChildSize: 0.9,
              builder: (BuildContext context, ScrollController scrollController) {
                return Container(
                  decoration: const BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(20),
                      topRight: Radius.circular(20),
                    ),
                  ),
                  child: Column(
                    children: [
                      // Handle bar
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        child: Container(
                          width: 40,
                          height: 4,
                          decoration: BoxDecoration(
                            color: Colors.grey[400],
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                      ),
                      // Header with close button
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              isOneWay ? 'Select Departure Date' : 'Select Travel Dates',
                              style: const TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.w700,
                                color: Colors.black,
                              ),
                            ),
                            GestureDetector(
                              onTap: () {
                                isModalOpenNotifier.value = false;
                                Navigator.pop(context);
                              },
                              child: Container(
                                decoration: BoxDecoration(
                                  color: const Color(0xFFD32F2F),
                                  borderRadius: BorderRadius.circular(50),
                                ),
                                padding: const EdgeInsets.all(4),
                                child: const Icon(
                                  Icons.close,
                                  color: Colors.white,
                                  size: 20,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      // Date selection info
                      if (!isOneWay)
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                          child: Container(
                            padding: const EdgeInsets.all(12),
                            decoration: BoxDecoration(
                              color: const Color(0xFF1e5a8e).withOpacity(0.1),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceAround,
                              children: [
                                Expanded(
                                  child: Column(
                                    children: [
                                      Text(
                                        'Departure',
                                        style: TextStyle(
                                          fontSize: 12,
                                          color: Colors.grey[600],
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                      const SizedBox(height: 4),
                                      Text(
                                        tempDepartDate != null
                                            ? '${tempDepartDate!.day} ${_getMonthAbbr(tempDepartDate!.month)}'
                                            : 'Select',
                                        style: const TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.w700,
                                          color: Color(0xFFD32F2F),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                const Icon(Icons.arrow_forward, color: Color(0xFF1e5a8e)),
                                Expanded(
                                  child: Column(
                                    children: [
                                      Text(
                                        'Return',
                                        style: TextStyle(
                                          fontSize: 12,
                                          color: Colors.grey[600],
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                      const SizedBox(height: 4),
                                      Text(
                                        tempReturnDate != null
                                            ? '${tempReturnDate!.day} ${_getMonthAbbr(tempReturnDate!.month)}'
                                            : 'Select',
                                        style: const TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.w700,
                                          color: Color(0xFFD32F2F),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      // Date picker content
                      Expanded(
                        child: SingleChildScrollView(
                          controller: scrollController,
                          child: Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const SizedBox(height: 12),
                                if (isOneWay)
                                  _buildCalendar(
                                    context,
                                    tempDepartDate ?? DateTime.now(),
                                    (date) {
                                      setModalState(() {
                                        tempDepartDate = date;
                                      });
                                    },
                                  )
                                else
                                  _buildDateRangeCalendar(
                                    context,
                                    currentMonth,
                                    tempDepartDate,
                                    tempReturnDate,
                                    (newMonth) {
                                      setModalState(() {
                                        currentMonth = newMonth;
                                      });
                                    },
                                    (date) {
                                      setModalState(() {
                                        if (tempDepartDate == null || (tempReturnDate != null)) {
                                          // First selection or reset
                                          tempDepartDate = date;
                                          tempReturnDate = null;
                                          isSelectingReturn = true;
                                        } else if (date.isBefore(tempDepartDate!) || date.isAtSameMomentAs(tempDepartDate!)) {
                                          // Selected before or same as departure, reset
                                          tempDepartDate = date;
                                          tempReturnDate = null;
                                          isSelectingReturn = true;
                                        } else {
                                          // Second selection, set as return
                                          tempReturnDate = date;
                                          isSelectingReturn = false;
                                        }
                                      });
                                    },
                                    isSelectingReturn: isSelectingReturn,
                                  ),
                                const SizedBox(height: 24),
                              ],
                            ),
                          ),
                        ),
                      ),
                      // Done button
                      Padding(
                        padding: const EdgeInsets.all(20),
                        child: SizedBox(
                          width: double.infinity,
                          height: 50,
                          child: ElevatedButton(
                            onPressed: (tempDepartDate == null || (!isOneWay && tempReturnDate == null))
                                ? null
                                : () {
                                    setState(() {
                                      _departDate = tempDepartDate;
                                      if (!isOneWay) {
                                        _returnDate = tempReturnDate;
                                      } else {
                                        _returnDate = null;
                                      }
                                    });
                                    isModalOpenNotifier.value = false;
                                    Navigator.pop(context);
                                  },
                            style: ElevatedButton.styleFrom(
                              backgroundColor: const Color(0xFF1e5a8e),
                              disabledBackgroundColor: Colors.grey[300],
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
                      ),
                    ],
                  ),
                );
              },
            );
          },
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  Widget _buildCalendar(
    BuildContext context,
    DateTime selectedDate,
    Function(DateTime) onDateSelected, {
    DateTime? minDate,
  }) {
    final now = DateTime.now();
    final currentMonth = DateTime(selectedDate.year, selectedDate.month);
    final daysInMonth = DateTime(selectedDate.year, selectedDate.month + 1, 0).day;
    final firstDayOfWeek = DateTime(selectedDate.year, selectedDate.month, 1).weekday;
    final currentMonthNormalized = DateTime(now.year, now.month);
    final isPreviousMonthDisabled = currentMonth.year == currentMonthNormalized.year && 
        currentMonth.month == currentMonthNormalized.month;
    
    return Column(
      children: [
        // Month/Year header with navigation
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            IconButton(
              icon: Icon(
                Icons.chevron_left,
                color: isPreviousMonthDisabled ? Colors.grey[300] : Colors.black,
              ),
              onPressed: isPreviousMonthDisabled ? null : () {
                final newDate = DateTime(selectedDate.year, selectedDate.month - 1);
                onDateSelected(DateTime(newDate.year, newDate.month, selectedDate.day.clamp(1, DateTime(newDate.year, newDate.month + 1, 0).day)));
              },
            ),
            Text(
              '${_getMonthName(selectedDate.month)} ${selectedDate.year}',
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.black,
              ),
            ),
            IconButton(
              icon: const Icon(Icons.chevron_right),
              onPressed: () {
                final newDate = DateTime(selectedDate.year, selectedDate.month + 1);
                onDateSelected(DateTime(newDate.year, newDate.month, selectedDate.day.clamp(1, DateTime(newDate.year, newDate.month + 1, 0).day)));
              },
            ),
          ],
        ),
        const SizedBox(height: 8),
        // Weekday headers
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'].map((day) {
            return SizedBox(
              width: 40,
              child: Center(
                child: Text(
                  day,
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: Colors.grey[600],
                  ),
                ),
              ),
            );
          }).toList(),
        ),
        const SizedBox(height: 8),
        // Calendar grid
        ...List.generate((daysInMonth + firstDayOfWeek % 7 + 6) ~/ 7, (weekIndex) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: List.generate(7, (dayIndex) {
                final dayNumber = weekIndex * 7 + dayIndex - (firstDayOfWeek % 7) + 1;
                
                if (dayNumber < 1 || dayNumber > daysInMonth) {
                  return const SizedBox(width: 40, height: 40);
                }
                
                final date = DateTime(selectedDate.year, selectedDate.month, dayNumber);
                final isSelected = date.year == selectedDate.year &&
                    date.month == selectedDate.month &&
                    date.day == selectedDate.day;
                final isPast = minDate != null && date.isBefore(minDate);
                final isDisabled = isPast || date.isBefore(DateTime(now.year, now.month, now.day));
                
                return GestureDetector(
                  onTap: (isPast || isDisabled) ? null : () => onDateSelected(date),
                  child: Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: isSelected
                          ? const Color(0xFFD32F2F)
                          : Colors.transparent,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Center(
                      child: Text(
                        dayNumber.toString(),
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: isSelected ? FontWeight.w700 : FontWeight.w500,
                          color: (isPast || isDisabled)
                              ? Colors.grey[400]
                              : isSelected
                                  ? Colors.white
                                  : Colors.black,
                        ),
                      ),
                    ),
                  ),
                );
              }),
            ),
          );
        }),
      ],
    );
  }

  // New unified date range calendar
  Widget _buildDateRangeCalendar(
    BuildContext context,
    DateTime currentMonth,
    DateTime? departDate,
    DateTime? returnDate,
    Function(DateTime) onMonthChanged,
    Function(DateTime) onDateTap, {
    bool isSelectingReturn = false,
  }) {
    final now = DateTime.now();
    final daysInMonth = DateTime(currentMonth.year, currentMonth.month + 1, 0).day;
    final firstDayOfWeek = DateTime(currentMonth.year, currentMonth.month, 1).weekday;
    final currentMonthNormalized = DateTime(now.year, now.month);
    final isPreviousMonthDisabled = currentMonth.year == currentMonthNormalized.year && 
        currentMonth.month == currentMonthNormalized.month;
    
    return Column(
      children: [
        // Month/Year header with navigation
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            IconButton(
              icon: Icon(
                Icons.chevron_left,
                color: isPreviousMonthDisabled ? Colors.grey[300] : Colors.black,
              ),
              onPressed: isPreviousMonthDisabled ? null : () {
                final newDate = DateTime(currentMonth.year, currentMonth.month - 1);
                onMonthChanged(newDate);
              },
            ),
            Text(
              '${_getMonthName(currentMonth.month)} ${currentMonth.year}',
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.black,
              ),
            ),
            IconButton(
              icon: const Icon(Icons.chevron_right),
              onPressed: () {
                final newDate = DateTime(currentMonth.year, currentMonth.month + 1);
                onMonthChanged(newDate);
              },
            ),
          ],
        ),
        const SizedBox(height: 8),
        // Weekday headers
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'].map((day) {
            return SizedBox(
              width: 40,
              child: Center(
                child: Text(
                  day,
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: Colors.grey[600],
                  ),
                ),
              ),
            );
          }).toList(),
        ),
        const SizedBox(height: 8),
        // Calendar grid
        ...List.generate((daysInMonth + firstDayOfWeek % 7 + 6) ~/ 7, (weekIndex) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: List.generate(7, (dayIndex) {
                final dayNumber = weekIndex * 7 + dayIndex - (firstDayOfWeek % 7) + 1;
                
                if (dayNumber < 1 || dayNumber > daysInMonth) {
                  return const SizedBox(width: 40, height: 40);
                }
                
                final date = DateTime(currentMonth.year, currentMonth.month, dayNumber);
                final normalizedDate = DateTime(date.year, date.month, date.day);
                final normalizedDepart = departDate != null ? DateTime(departDate.year, departDate.month, departDate.day) : null;
                final normalizedReturn = returnDate != null ? DateTime(returnDate.year, returnDate.month, returnDate.day) : null;
                
                final isDepartDate = normalizedDepart != null && normalizedDate.isAtSameMomentAs(normalizedDepart);
                final isReturnDate = normalizedReturn != null && normalizedDate.isAtSameMomentAs(normalizedReturn);
                final isInRange = normalizedDepart != null && normalizedReturn != null &&
                    normalizedDate.isAfter(normalizedDepart) && normalizedDate.isBefore(normalizedReturn);
                final isDisabled = date.isBefore(DateTime(now.year, now.month, now.day));
                
                Color? bgColor;
                Color? textColor;
                BorderRadius? borderRadius;
                
                if (isDepartDate || isReturnDate) {
                  bgColor = const Color(0xFFD32F2F);
                  textColor = Colors.white;
                  if (isDepartDate && isReturnDate) {
                    borderRadius = BorderRadius.circular(8);
                  } else if (isDepartDate) {
                    borderRadius = const BorderRadius.only(
                      topLeft: Radius.circular(8),
                      bottomLeft: Radius.circular(8),
                    );
                  } else {
                    borderRadius = const BorderRadius.only(
                      topRight: Radius.circular(8),
                      bottomRight: Radius.circular(8),
                    );
                  }
                } else if (isInRange) {
                  bgColor = const Color(0xFFD32F2F).withOpacity(0.15);
                  textColor = Colors.black;
                  borderRadius = BorderRadius.zero;
                } else {
                  bgColor = Colors.transparent;
                  textColor = isDisabled ? Colors.grey[400] : Colors.black;
                  borderRadius = BorderRadius.circular(8);
                }
                
                return GestureDetector(
                  onTap: isDisabled ? null : () => onDateTap(date),
                  child: Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: bgColor,
                      borderRadius: borderRadius,
                    ),
                    child: Center(
                      child: Text(
                        dayNumber.toString(),
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: (isDepartDate || isReturnDate) ? FontWeight.w700 : FontWeight.w500,
                          color: textColor,
                        ),
                      ),
                    ),
                  ),
                );
              }),
            ),
          );
        }),
      ],
    );
  }

  String _getMonthName(int month) {
    const months = [
      'January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'
    ];
    return months[month - 1];
  }

  void _showPassengersBottomSheet(BuildContext context) {
    int tempAdults = _adultsCount;
    int tempChildren = _childrenCount;
    int tempInfants = _infantsCount;
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            int totalPassengers = tempAdults + tempChildren + tempInfants;
            
            return Container(
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.only(
                  topLeft: Radius.circular(20),
                  topRight: Radius.circular(20),
                ),
              ),
              padding: const EdgeInsets.all(20),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  // Handle bar
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: Colors.grey[400],
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 12),
                  // Close button
                  Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      GestureDetector(
                        onTap: () => Navigator.pop(context),
                        child: Container(
                          decoration: BoxDecoration(
                            color: const Color(0xFFD32F2F),
                            borderRadius: BorderRadius.circular(50),
                          ),
                          padding: const EdgeInsets.all(4),
                          child: const Icon(
                            Icons.close,
                            color: Colors.white,
                            size: 20,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  // Adults
                  _buildPassengerRow(
                    title: 'Adults',
                    subtitle: '12+',
                    count: tempAdults,
                    onDecrement: tempAdults > 1
                        ? () => setModalState(() => tempAdults--)
                        : null,
                    onIncrement: totalPassengers < 9
                        ? () => setModalState(() => tempAdults++)
                        : null,
                  ),
                  const SizedBox(height: 16),
                  // Children
                  _buildPassengerRow(
                    title: 'Children',
                    subtitle: '2 till 12',
                    count: tempChildren,
                    onDecrement: tempChildren > 0
                        ? () => setModalState(() => tempChildren--)
                        : null,
                    onIncrement: totalPassengers < 9
                        ? () => setModalState(() => tempChildren++)
                        : null,
                  ),
                  const SizedBox(height: 16),
                  // Infants
                  _buildPassengerRow(
                    title: 'Infants',
                    subtitle: '0 till 2',
                    count: tempInfants,
                    onDecrement: tempInfants > 0
                        ? () => setModalState(() => tempInfants--)
                        : null,
                    onIncrement: totalPassengers < 9
                        ? () => setModalState(() => tempInfants++)
                        : null,
                  ),
                  const SizedBox(height: 32),
                  // Done button
                  SizedBox(
                    width: double.infinity,
                    height: 50,
                    child: ElevatedButton(
                      onPressed: () {
                        setState(() {
                          _adultsCount = tempAdults;
                          _childrenCount = tempChildren;
                          _infantsCount = tempInfants;
                        });
                        isModalOpenNotifier.value = false;
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
                  const SizedBox(height: 8),
                ],
              ),
            );
          },
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  Widget _buildPassengerRow({
    required String title,
    required String subtitle,
    required int count,
    required VoidCallback? onDecrement,
    required VoidCallback? onIncrement,
  }) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Color(0xFF1e5a8e),
              ),
            ),
            const SizedBox(height: 4),
            Text(
              subtitle,
              style: TextStyle(
                fontSize: 12,
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
        Row(
          children: [
            // Decrement button
            GestureDetector(
              onTap: onDecrement,
              child: Container(
                width: 36,
                height: 36,
                decoration: BoxDecoration(
                  color: onDecrement != null
                      ? const Color(0xFF1e5a8e)
                      : Colors.grey[300],
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.remove,
                  color: onDecrement != null ? Colors.white : Colors.grey[500],
                  size: 20,
                ),
              ),
            ),
            const SizedBox(width: 20),
            // Count
            SizedBox(
              width: 30,
              child: Text(
                count.toString(),
                textAlign: TextAlign.center,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
            ),
            const SizedBox(width: 20),
            // Increment button
            GestureDetector(
              onTap: onIncrement,
              child: Container(
                width: 36,
                height: 36,
                decoration: BoxDecoration(
                  color: onIncrement != null
                      ? const Color(0xFF1e5a8e)
                      : Colors.grey[300],
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.add,
                  color: onIncrement != null ? Colors.white : Colors.grey[500],
                  size: 20,
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  void _showClassBottomSheet(BuildContext context) {
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.only(
              topLeft: Radius.circular(20),
              topRight: Radius.circular(20),
            ),
          ),
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              // Handle bar
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey[400],
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 12),
              // Close button
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  GestureDetector(
                    onTap: () {
                      isModalOpenNotifier.value = false;
                      Navigator.pop(context);
                    },
                    child: Container(
                      decoration: BoxDecoration(
                        color: const Color(0xFFD32F2F),
                        borderRadius: BorderRadius.circular(50),
                      ),
                      padding: const EdgeInsets.all(4),
                      child: const Icon(
                        Icons.close,
                        color: Colors.white,
                        size: 20,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              // Class options
              ..._flightClasses.map((className) {
                return GestureDetector(
                  onTap: () {
                    setState(() {
                      _selectedClass = className;
                    });
                    isModalOpenNotifier.value = false;
                    Navigator.pop(context);
                  },
                  child: Container(
                    width: double.infinity,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    child: Center(
                      child: Text(
                        className,
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: _selectedClass == className
                              ? FontWeight.w700
                              : FontWeight.w500,
                          color: const Color(0xFF1e5a8e),
                        ),
                      ),
                    ),
                  ),
                );
              }).toList(),
              const SizedBox(height: 16),
            ],
          ),
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  void _showMultiCityLocationBottomSheet(BuildContext context, int flightIndex, {required bool isFrom}) {
    final TextEditingController searchController = TextEditingController();
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            final filteredDestinations = _destinations.where((destination) {
              final name = destination['name']!.toLowerCase();
              final airport = destination['airport']!.toLowerCase();
              final query = searchController.text.toLowerCase();
              return name.contains(query) || airport.contains(query);
            }).toList();

            return DraggableScrollableSheet(
              initialChildSize: 0.6,
              minChildSize: 0.4,
              maxChildSize: 0.95,
              builder: (BuildContext context, ScrollController scrollController) {
                return Container(
                  decoration: const BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(20),
                      topRight: Radius.circular(20),
                    ),
                  ),
                  child: Column(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        child: Container(
                          width: 40,
                          height: 4,
                          decoration: BoxDecoration(
                            color: Colors.grey[400],
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 20,
                          vertical: 12,
                        ),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              isFrom ? 'Select Departure' : 'Select Destination',
                              style: const TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.w700,
                                color: Colors.black,
                              ),
                            ),
                            GestureDetector(
                              onTap: () => Navigator.pop(context),
                              child: Container(
                                decoration: BoxDecoration(
                                  color: const Color(0xFFD32F2F),
                                  borderRadius: BorderRadius.circular(50),
                                ),
                                padding: const EdgeInsets.all(4),
                                child: const Icon(
                                  Icons.close,
                                  color: Colors.white,
                                  size: 20,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20),
                        child: TextField(
                          controller: searchController,
                          autofocus: false,
                          onChanged: (value) {
                            setModalState(() {});
                          },
                          decoration: InputDecoration(
                            hintText: 'Search destinations...',
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(color: Colors.grey[300]!),
                            ),
                            enabledBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(color: Colors.grey[300]!),
                            ),
                            focusedBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(
                                color: Color(0xFF1e5a8e),
                                width: 2,
                              ),
                            ),
                            prefixIcon: const Icon(
                              Icons.search,
                              color: Colors.grey,
                            ),
                            suffixIcon: searchController.text.isNotEmpty
                                ? GestureDetector(
                                    onTap: () {
                                      searchController.clear();
                                      setModalState(() {});
                                    },
                                    child: const Icon(
                                      Icons.clear,
                                      color: Colors.grey,
                                    ),
                                  )
                                : const Icon(
                                    Icons.location_on,
                                    color: Color(0xFFD32F2F),
                                  ),
                            contentPadding: const EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 12,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),
                      Expanded(
                        child: filteredDestinations.isEmpty
                            ? Center(
                                child: Padding(
                                  padding: const EdgeInsets.all(24),
                                  child: Column(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Icon(
                                        Icons.search_off,
                                        size: 64,
                                        color: Colors.grey[400],
                                      ),
                                      const SizedBox(height: 16),
                                      Text(
                                        'No destinations found',
                                        style: TextStyle(
                                          fontSize: 16,
                                          fontWeight: FontWeight.w600,
                                          color: Colors.grey[600],
                                        ),
                                      ),
                                      const SizedBox(height: 8),
                                      Text(
                                        'Try searching with different keywords',
                                        style: TextStyle(
                                          fontSize: 13,
                                          color: Colors.grey[500],
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              )
                            : ListView.builder(
                                controller: scrollController,
                                padding: const EdgeInsets.symmetric(horizontal: 16),
                                itemCount: filteredDestinations.length,
                                itemBuilder: (context, index) {
                                  final destination = filteredDestinations[index];
                                  final currentLocation = isFrom 
                                      ? _multiCityFlights[flightIndex]['from']
                                      : _multiCityFlights[flightIndex]['to'];
                                  final isSelected = currentLocation == destination['name'];

                                  return GestureDetector(
                                    onTap: () {
                                      setState(() {
                                        if (isFrom) {
                                          _multiCityFlights[flightIndex]['from'] = destination['name']!;
                                        } else {
                                          _multiCityFlights[flightIndex]['to'] = destination['name']!;
                                        }
                                      });
                                      Navigator.pop(context);
                                    },
                                    child: Container(
                                      margin: const EdgeInsets.only(bottom: 12),
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 16,
                                        vertical: 12,
                                      ),
                                      decoration: BoxDecoration(
                                        color: isSelected
                                            ? const Color(0xFF1e5a8e).withOpacity(0.1)
                                            : Colors.white,
                                        borderRadius: BorderRadius.circular(12),
                                        border: Border.all(
                                          color: isSelected
                                              ? const Color(0xFF1e5a8e)
                                              : Colors.grey[200]!,
                                          width: isSelected ? 2 : 1,
                                        ),
                                      ),
                                      child: Row(
                                        children: [
                                          Container(
                                            width: 40,
                                            height: 40,
                                            decoration: BoxDecoration(
                                              color: const Color(0xFFD32F2F)
                                                  .withOpacity(0.1),
                                              borderRadius: BorderRadius.circular(8),
                                            ),
                                            child: const Icon(
                                              Icons.location_on,
                                              color: Color(0xFFD32F2F),
                                              size: 20,
                                            ),
                                          ),
                                          const SizedBox(width: 12),
                                          Expanded(
                                            child: Column(
                                              crossAxisAlignment:
                                                  CrossAxisAlignment.start,
                                              children: [
                                                Text(
                                                  destination['name']!,
                                                  style: TextStyle(
                                                    fontSize: 14,
                                                    fontWeight: FontWeight.w600,
                                                    color: isSelected
                                                        ? const Color(0xFF1e5a8e)
                                                        : Colors.black,
                                                  ),
                                                ),
                                                const SizedBox(height: 4),
                                                Text(
                                                  destination['airport']!,
                                                  style: TextStyle(
                                                    fontSize: 12,
                                                    color: Colors.grey[600],
                                                  ),
                                                ),
                                              ],
                                            ),
                                          ),
                                          if (isSelected)
                                            const Icon(
                                              Icons.check_circle,
                                              color: Color(0xFF1e5a8e),
                                              size: 24,
                                            ),
                                        ],
                                      ),
                                    ),
                                  );
                                },
                              ),
                      ),
                      const SizedBox(height: 24),
                    ],
                  ),
                );
              },
            );
          },
        );
      },
    );
  }

  void _showMultiCityDateBottomSheet(BuildContext context, int flightIndex) {
    DateTime tempDate = _multiCityFlights[flightIndex]['date'] ?? DateTime.now();
    
    // Calculate minimum date based on previous flight (next day after previous flight)
    DateTime? minDate;
    if (flightIndex > 0 && _multiCityFlights[flightIndex - 1]['date'] != null) {
      final previousFlightDate = _multiCityFlights[flightIndex - 1]['date'] as DateTime;
      // Set minimum date to the day AFTER the previous flight
      minDate = previousFlightDate.add(const Duration(days: 1));
      
      // Ensure tempDate is at least minDate
      if (tempDate.isBefore(minDate)) {
        tempDate = minDate;
      }
    }
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return DraggableScrollableSheet(
              initialChildSize: 0.65,
              minChildSize: 0.5,
              maxChildSize: 0.9,
              builder: (BuildContext context, ScrollController scrollController) {
                return Container(
                  decoration: const BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(20),
                      topRight: Radius.circular(20),
                    ),
                  ),
                  child: Column(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        child: Container(
                          width: 40,
                          height: 4,
                          decoration: BoxDecoration(
                            color: Colors.grey[400],
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              'Select Flight ${flightIndex + 1} Date',
                              style: const TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.w700,
                                color: Colors.black,
                              ),
                            ),
                            GestureDetector(
                              onTap: () {
                                isModalOpenNotifier.value = false;
                                Navigator.pop(context);
                              },
                              child: Container(
                                decoration: BoxDecoration(
                                  color: const Color(0xFFD32F2F),
                                  borderRadius: BorderRadius.circular(50),
                                ),
                                padding: const EdgeInsets.all(4),
                                child: const Icon(
                                  Icons.close,
                                  color: Colors.white,
                                  size: 20,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      Expanded(
                        child: SingleChildScrollView(
                          controller: scrollController,
                          child: Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                if (minDate != null)
                                  Padding(
                                    padding: const EdgeInsets.only(bottom: 8),
                                    child: Container(
                                      padding: const EdgeInsets.all(10),
                                      decoration: BoxDecoration(
                                        color: const Color(0xFF1e5a8e).withOpacity(0.1),
                                        borderRadius: BorderRadius.circular(8),
                                      ),
                                      child: Row(
                                        children: [
                                          const Icon(
                                            Icons.info_outline,
                                            color: Color(0xFF1e5a8e),
                                            size: 18,
                                          ),
                                          const SizedBox(width: 8),
                                          Expanded(
                                            child: Text(
                                              'Date must be after Flight ${flightIndex}',
                                              style: const TextStyle(
                                                fontSize: 12,
                                                color: Color(0xFF1e5a8e),
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ),
                                const Text(
                                  'Departure Date',
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: Color(0xFF1e5a8e),
                                  ),
                                ),
                                const SizedBox(height: 12),
                                _buildCalendar(
                                  context,
                                  tempDate,
                                  (date) {
                                    setModalState(() {
                                      tempDate = date;
                                    });
                                  },
                                  minDate: minDate,
                                ),
                                const SizedBox(height: 24),
                              ],
                            ),
                          ),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.all(20),
                        child: SizedBox(
                          width: double.infinity,
                          height: 50,
                          child: ElevatedButton(
                            onPressed: () {
                              setState(() {
                                _multiCityFlights[flightIndex]['date'] = tempDate;
                              });
                              isModalOpenNotifier.value = false;
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
                      ),
                    ],
                  ),
                );
              },
            );
          },
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  Widget _buildDestinationCarousel() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Section Title
        const Padding(
          padding: EdgeInsets.only(bottom: 12),
          child: Text(
            'Get the best travel experience',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w700,
              color: Colors.black,
            ),
          ),
        ),
        // Carousel
        SizedBox(
          height: 200,
          child: PageView.builder(
            controller: _carouselController,
            onPageChanged: (index) {
              setState(() {
                _currentCarouselIndex = index;
              });
            },
            itemCount: _banners.isNotEmpty ? _banners.length : _fallbackCarouselImages.length,
            itemBuilder: (context, index) {
              // Use backend banners if available, otherwise use fallback images
              final imageUrl = _banners.isNotEmpty
                  ? '${Endpoints.baseUrl.replaceAll('/api', '')}${_banners[index].imagePath}'
                  : null;
              final assetImage = _banners.isEmpty ? _fallbackCarouselImages[index] : null;
              final title = _banners.isNotEmpty ? _banners[index].title : _fallbackCarouselTitles[index];

              return Padding(
                padding: const EdgeInsets.only(right: 12),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(15),
                  child: Stack(
                    fit: StackFit.expand,
                    children: [
                      // Image - either from network (backend) or assets (fallback)
                      if (imageUrl != null)
                        Image.network(
                          imageUrl,
                          fit: BoxFit.cover,
                          loadingBuilder: (context, child, loadingProgress) {
                            if (loadingProgress == null) return child;
                            return Container(
                              color: Colors.grey[300],
                              child: Center(
                                child: CircularProgressIndicator(
                                  value: loadingProgress.expectedTotalBytes != null
                                      ? loadingProgress.cumulativeBytesLoaded /
                                          loadingProgress.expectedTotalBytes!
                                      : null,
                                ),
                              ),
                            );
                          },
                          errorBuilder: (context, error, stackTrace) {
                            // Fallback to asset if network image fails
                            return Image.asset(
                              _fallbackCarouselImages[index % _fallbackCarouselImages.length],
                              fit: BoxFit.cover,
                              errorBuilder: (context, error, stackTrace) {
                                return Container(
                                  color: Colors.grey[300],
                                  child: Center(
                                    child: Icon(
                                      Icons.image,
                                      size: 60,
                                      color: Colors.grey[500],
                                    ),
                                  ),
                                );
                              },
                            );
                          },
                        )
                      else if (assetImage != null)
                        Image.asset(
                          assetImage,
                          fit: BoxFit.cover,
                          errorBuilder: (context, error, stackTrace) {
                            return Container(
                              color: Colors.grey[300],
                              child: Center(
                                child: Icon(
                                  Icons.image,
                                  size: 60,
                                  color: Colors.grey[500],
                                ),
                              ),
                            );
                          },
                        ),
                      // Gradient overlay
                      Container(
                        decoration: BoxDecoration(
                          gradient: LinearGradient(
                            begin: Alignment.bottomCenter,
                            end: Alignment.topCenter,
                            colors: [
                              Colors.black.withOpacity(0.6),
                              Colors.transparent,
                            ],
                          ),
                        ),
                      ),
                      // Title
                      Positioned(
                        bottom: 16,
                        left: 16,
                        child: Text(
                          title,
                          style: const TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.w700,
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
        const SizedBox(height: 12),
        // Carousel indicators
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: List.generate(
            _banners.isNotEmpty ? _banners.length : _fallbackCarouselImages.length,
            (index) => Container(
              width: _currentCarouselIndex == index ? 24 : 8,
              height: 8,
              margin: const EdgeInsets.symmetric(horizontal: 4),
              decoration: BoxDecoration(
                color: _currentCarouselIndex == index
                    ? const Color(0xFF1e5a8e)
                    : Colors.grey[300],
                borderRadius: BorderRadius.circular(4),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildAdvertisement() {
    // Display first homepage banner if available, otherwise show fallback ad
    if (_homepageBanners.isNotEmpty) {
      final banner = _homepageBanners[0];
      final imageUrl = '${Endpoints.baseUrl.replaceAll('/api', '')}${banner.imagePath}';
      
      return ClipRRect(
        borderRadius: BorderRadius.circular(15),
        child: Image.network(
          imageUrl,
          width: double.infinity,
          fit: BoxFit.cover,
          loadingBuilder: (context, child, loadingProgress) {
            if (loadingProgress == null) return child;
            return Container(
              width: double.infinity,
              height: 100,
              decoration: BoxDecoration(
                color: Colors.grey[300],
                borderRadius: BorderRadius.circular(15),
              ),
              child: Center(
                child: CircularProgressIndicator(
                  value: loadingProgress.expectedTotalBytes != null
                      ? loadingProgress.cumulativeBytesLoaded /
                          loadingProgress.expectedTotalBytes!
                      : null,
                ),
              ),
            );
          },
          errorBuilder: (context, error, stackTrace) {
            // Fallback to static ad if network image fails
            return Image.asset(
              'assets/images/Ad.png',
              width: double.infinity,
              fit: BoxFit.cover,
              errorBuilder: (context, error, stackTrace) {
                return Container(
                  width: double.infinity,
                  height: 100,
                  decoration: BoxDecoration(
                    color: Colors.grey[300],
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Center(
                    child: Icon(
                      Icons.image,
                      size: 40,
                      color: Colors.grey[500],
                    ),
                  ),
                );
              },
            );
          },
        ),
      );
    }
    
    // Fallback to static ad
    return ClipRRect(
      borderRadius: BorderRadius.circular(15),
      child: Image.asset(
        'assets/images/Ad.png',
        width: double.infinity,
        fit: BoxFit.cover,
        errorBuilder: (context, error, stackTrace) {
          return Container(
            width: double.infinity,
            height: 100,
            decoration: BoxDecoration(
              color: Colors.grey[300],
              borderRadius: BorderRadius.circular(15),
            ),
            child: Center(
              child: Icon(
                Icons.image,
                size: 40,
                color: Colors.grey[500],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildPopularLocations() {
    // Show loading indicator while fetching data
    if (_popularLocationsLoading) {
      return const SizedBox(
        height: 220,
        child: Center(
          child: CircularProgressIndicator(
            valueColor: AlwaysStoppedAnimation<Color>(Color(0xFF1e5a8e)),
          ),
        ),
      );
    }

    // Show empty state if no locations available
    if (_popularLocations.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Section Title
        const Padding(
          padding: EdgeInsets.only(bottom: 16),
          child: Text(
            'Popular Locations',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w700,
              color: Colors.black,
            ),
          ),
        ),
        // Horizontal scrollable cards
        SizedBox(
          height: 220,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: _popularLocations.length,
            itemBuilder: (context, index) {
              final location = _popularLocations[index];
              return Container(
                width: 240,
                margin: EdgeInsets.only(
                  right: index < _popularLocations.length - 1 ? 16 : 0,
                ),
                child: Stack(
                  children: [
                    // Card with image
                    ClipRRect(
                      borderRadius: BorderRadius.circular(15),
                      child: Stack(
                        fit: StackFit.expand,
                        children: [
                          // Background image from backend
                          Image.network(
                            '${Endpoints.baseUrl.replaceAll('/api', '')}${location.imagePath}',
                            fit: BoxFit.cover,
                            errorBuilder: (context, error, stackTrace) {
                              return Container(
                                color: Colors.grey[300],
                                child: Icon(
                                  Icons.image,
                                  size: 60,
                                  color: Colors.grey[500],
                                ),
                              );
                            },
                            loadingBuilder: (context, child, loadingProgress) {
                              if (loadingProgress == null) return child;
                              return Container(
                                color: Colors.grey[200],
                                child: Center(
                                  child: CircularProgressIndicator(
                                    value: loadingProgress.expectedTotalBytes != null
                                        ? loadingProgress.cumulativeBytesLoaded /
                                            loadingProgress.expectedTotalBytes!
                                        : null,
                                    valueColor: const AlwaysStoppedAnimation<Color>(Color(0xFF1e5a8e)),
                                  ),
                                ),
                              );
                            },
                          ),
                          // Gradient overlay
                          Container(
                            decoration: BoxDecoration(
                              gradient: LinearGradient(
                                begin: Alignment.topCenter,
                                end: Alignment.bottomCenter,
                                colors: [
                                  Colors.black.withOpacity(0.3),
                                  Colors.black.withOpacity(0.7),
                                ],
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    // Title badge (if available)
                    if (location.title.isNotEmpty)
                      Positioned(
                        top: 12,
                        left: 12,
                        child: Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 10,
                            vertical: 6,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.9),
                            borderRadius: BorderRadius.circular(20),
                          ),
                          child: Text(
                            location.title,
                            style: const TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                              color: Color(0xFF1e5a8e),
                            ),
                          ),
                        ),
                      ),
                    // Location and Book button
                    Positioned(
                      bottom: 12,
                      left: 12,
                      right: 12,
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          // Location
                          Expanded(
                            child: Row(
                              children: [
                                const Icon(
                                  Icons.location_on,
                                  color: Colors.white,
                                  size: 18,
                                ),
                                const SizedBox(width: 4),
                                Expanded(
                                  child: Text(
                                    location.location,
                                    style: const TextStyle(
                                      fontSize: 14,
                                      fontWeight: FontWeight.w600,
                                      color: Colors.white,
                                    ),
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(width: 8),
                          // Book button
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 16,
                              vertical: 8,
                            ),
                            decoration: BoxDecoration(
                              color: const Color(0xFF1e5a8e),
                              borderRadius: BorderRadius.circular(20),
                            ),
                            child: const Text(
                              'Book',
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: Colors.white,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
