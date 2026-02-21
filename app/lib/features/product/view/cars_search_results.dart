import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'car_details.dart';
import '../../../core/widgets/unified_ui_components.dart';
import '../../../core/network/endpoints.dart';
import '../service/banner_service.dart';
import '../model/banner_model.dart';
import '../../../core/services/car_service.dart';

class CarsSearchResults extends StatefulWidget {
  final String pickupLocation;
  final String dropoffLocation;
  final DateTime? pickupDate;
  final DateTime? dropoffDate;
  final String pickupTime;
  final String dropoffTime;

  const CarsSearchResults({
    super.key,
    required this.pickupLocation,
    required this.dropoffLocation,
    this.pickupDate,
    this.dropoffDate,
    required this.pickupTime,
    required this.dropoffTime,
  });

  @override
  State<CarsSearchResults> createState() => _CarsSearchResultsState();
}

class _CarsSearchResultsState extends State<CarsSearchResults> {
  // Filter state
  String _sortBy = 'Cheapest';
  String _selectedCarType = '';
  String _selectedRentalCompany = '';
  Map<String, bool> _features = {
    'Automatic Transmission': false,
    'Manual Transmission': false,
    'Air Conditioning': false,
    'Unlimited Mileage': false,
    'Free Cancellation': false,
    'Child Seat Available': false,
    'GPS Included': false,
    'Additional Driver Included': false,
  };

  // Search banners
  final BannerService _bannerService = BannerService();
  final CarService _carService = CarService();
  List<BannerModel> _searchBanners = [];
  bool _searchBannersLoading = true;
  
  // Car results
  List<dynamic> _carResults = [];
  List<dynamic> _filteredCarResults = [];
  bool _carResultsLoading = true;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    _loadSearchBanners();
    _loadCarResults();
  }

  Future<void> _loadSearchBanners() async {
    try {
      final banners = await _bannerService.getHomepageBanners();
      if (mounted) {
        setState(() {
          // Use second banner (index 1) for search results
          _searchBanners = banners.length > 1 ? [banners[1]] : [];
          _searchBannersLoading = false;
        });
      }
    } catch (e) {
      print('[CARS_SEARCH] Failed to load search banners: $e');
      if (mounted) {
        setState(() {
          _searchBannersLoading = false;
        });
      }
    }
  }

  Future<void> _loadCarResults() async {
    try {
      setState(() {
        _carResultsLoading = true;
        _errorMessage = null;
      });

      // Search with date and time parameters for availability checking and cost calculation
      final cars = await _carService.searchCars(
        pickupLocation: widget.pickupLocation,
        dropoffLocation: widget.dropoffLocation,
        pickupDate: widget.pickupDate,
        dropoffDate: widget.dropoffDate,
        pickupTime: widget.pickupTime,
        dropoffTime: widget.dropoffTime,
      );

      if (mounted) {
        setState(() {
          _carResults = cars;
          _filteredCarResults = cars;
          _carResultsLoading = false;
        });
      }
    } catch (e) {
      print('[CARS_SEARCH] Failed to load car results: $e');
      if (mounted) {
        setState(() {
          _errorMessage = 'Failed to load cars. Please try again.';
          _carResultsLoading = false;
        });
      }
    }
  }

  // Mock car results
  final List<Map<String, dynamic>> _mockCarResults = [
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
  ];

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
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '${widget.pickupLocation} (1 traveler)',
              style: const TextStyle(
                color: Color(0xFF1e5a8e),
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
            Text(
              '${widget.pickupTime} - ${widget.dropoffTime}',
              style: const TextStyle(
                color: Colors.grey,
                fontSize: 12,
                fontWeight: FontWeight.w400,
              ),
            ),
          ],
        ),
        actions: [
          UnifiedFilterButton(
            onPressed: _showFilterBottomSheet,
          ),
        ],
      ),
      body: _carResultsLoading
          ? const Center(
              child: CircularProgressIndicator(
                valueColor: AlwaysStoppedAnimation<Color>(Color(0xFF1e5a8e)),
              ),
            )
          : _errorMessage != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(
                        Icons.error_outline,
                        size: 60,
                        color: Colors.red,
                      ),
                      const SizedBox(height: 16),
                      Text(
                        _errorMessage!,
                        textAlign: TextAlign.center,
                        style: const TextStyle(fontSize: 16),
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _loadCarResults,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF1e5a8e),
                        ),
                        child: const Text('Retry'),
                      ),
                    ],
                  ),
                )
              : _filteredCarResults.isEmpty
                  ? const Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.search_off,
                            size: 60,
                            color: Colors.grey,
                          ),
                          SizedBox(height: 16),
                          Text(
                            'No cars found',
                            style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600),
                          ),
                          SizedBox(height: 8),
                          Text(
                            'Try adjusting your search criteria',
                            style: TextStyle(color: Colors.grey),
                          ),
                        ],
                      ),
                    )
                  : ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _filteredCarResults.length + 1, // +1 for the "Results" header
        itemBuilder: (context, index) {
          if (index == 0) {
            // Results header
            return const Padding(
              padding: EdgeInsets.only(bottom: 16),
              child: Text(
                'Results',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w700,
                  color: Colors.black,
                ),
              ),
            );
          }

          final carIndex = index - 1;

          // Add ad banner after 1st car (if there's at least 1 result)
          if (carIndex == 0) {
            return Column(
              children: [
                _buildCarCard(_filteredCarResults[carIndex]),
                const SizedBox(height: 16),
                // Ad Banner from database
                _buildSearchBanner(),
                const SizedBox(height: 16),
              ],
            );
          }

          return Column(
            children: [
              _buildCarCard(_filteredCarResults[carIndex]),
              const SizedBox(height: 16),
            ],
          );
        },
      ),
    );
  }

  Widget _buildSearchBanner() {
    if (_searchBanners.isNotEmpty) {
      final banner = _searchBanners[0];
      final imageUrl = '${Endpoints.baseUrl.replaceAll('/api', '')}${banner.imagePath}';
      
      return ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: Image.network(
          imageUrl,
          width: double.infinity,
          height: 100,
          fit: BoxFit.cover,
          loadingBuilder: (context, child, loadingProgress) {
            if (loadingProgress == null) return child;
            return Container(
              height: 100,
              decoration: BoxDecoration(
                color: Colors.grey.shade200,
                borderRadius: BorderRadius.circular(12),
              ),
              child: const Center(
                child: CircularProgressIndicator(),
              ),
            );
          },
          errorBuilder: (context, error, stackTrace) {
            return Image.asset(
              'assets/images/Ad.png',
              width: double.infinity,
              height: 100,
              fit: BoxFit.cover,
              errorBuilder: (context, error, stackTrace) {
                return Container(
                  height: 100,
                  decoration: BoxDecoration(
                    color: Colors.grey.shade200,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Center(
                    child: Text('Advertisement'),
                  ),
                );
              },
            );
          },
        ),
      );
    }
    
    return ClipRRect(
      borderRadius: BorderRadius.circular(12),
      child: Image.asset(
        'assets/images/Ad.png',
        width: double.infinity,
        height: 100,
        fit: BoxFit.cover,
        errorBuilder: (context, error, stackTrace) {
          return Container(
            height: 100,
            decoration: BoxDecoration(
              color: Colors.grey.shade200,
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Center(
              child: Text('Advertisement'),
            ),
          );
        },
      ),
    );
  }

  Widget _buildCarCard(dynamic car) {
    // Handle both Map (from API) and dynamic types - convert to Map<String, dynamic>
    final Map<String, dynamic> carMap = car is Map<String, dynamic> 
        ? car 
        : (car is Map ? Map<String, dynamic>.from(car as Map) : {});
    
    final carType = carMap['car_type'] ?? carMap['category'] ?? 'Car';
    final carName = carMap['car_name'] ?? carMap['name'] ?? 'Vehicle';
    final passengers = carMap['passengers']?.toString() ?? '4';
    final transmission = carMap['transmission'] ?? 'Automatic';
    final mileage = carMap['mileage'] ?? 'Unlimited mileage';
    
    // Use calculated_cost if available, otherwise fallback to cost or price
    final calculatedCost = carMap['calculated_cost'];
    final costPerDay = carMap['cost_per_day'] ?? carMap['cost'] ?? 0;
    final displayCost = calculatedCost != null ? calculatedCost.toString() : costPerDay.toString();
    final currency = carMap['currency'] ?? '\$';
    
    final shuttle = carMap['shuttle_to_counter'] == true 
        ? 'Shuttle to counter and car' 
        : (carMap['shuttle'] ?? 'Walk to counter');
    
    // Determine if showing total or per-day price
    final isPriceTotal = calculatedCost != null;
    final priceLabel = isPriceTotal ? 'Total' : 'Per day';
    
    return GestureDetector(
      onTap: () {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => CarDetails(
              car: carMap,
              pickupLocation: widget.pickupLocation,
              dropoffLocation: widget.dropoffLocation,
              pickupDate: widget.pickupDate,
              dropoffDate: widget.dropoffDate,
              pickupTime: widget.pickupTime,
              dropoffTime: widget.dropoffTime,
            ),
          ),
        );
      },
      child: Container(
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
                  carType,
                  style: const TextStyle(
                    color: Color(0xFFD32F2F),
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  carName,
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '$passengers Passengers - $transmission',
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  mileage,
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  shuttle,
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
              Column(
                children: [
                  Text(
                    '$currency${double.parse(displayCost).toStringAsFixed(2)}',
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w700,
                      color: Colors.black,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    priceLabel,
                    style: TextStyle(
                      fontSize: 11,
                      color: Colors.grey[600],
                    ),
                  ),
                ],
              ),
            ],
          ),
        ],
      ),
    ),
    );
  }

  void _clearFilters() {
    setState(() {
      _sortBy = 'Cheapest';
      _selectedCarType = '';
      _selectedRentalCompany = '';
      _features = {
        'Automatic Transmission': false,
        'Manual Transmission': false,
        'Air Conditioning': false,
        'Unlimited Mileage': false,
        'Free Cancellation': false,
        'Child Seat Available': false,
        'GPS Included': false,
        'Additional Driver Included': false,
      };
      _filteredCarResults = _carResults;
    });
  }

  void _applyFilters() {
    setState(() {
      _filteredCarResults = _carResults.where((car) {
        final Map<String, dynamic> carMap = car is Map<String, dynamic> 
            ? car 
            : (car is Map ? Map<String, dynamic>.from(car as Map) : {});
        
        // Filter by car type
        if (_selectedCarType.isNotEmpty) {
          final carType = (carMap['car_type'] ?? carMap['category'] ?? '').toString().toLowerCase();
          if (!carType.contains(_selectedCarType.toLowerCase())) {
            return false;
          }
        }
        
        // Filter by transmission (from features)
        if (_features['Automatic Transmission'] == true) {
          final transmission = (carMap['transmission'] ?? '').toString().toLowerCase();
          if (!transmission.contains('automatic')) {
            return false;
          }
        }
        
        if (_features['Manual Transmission'] == true) {
          final transmission = (carMap['transmission'] ?? '').toString().toLowerCase();
          if (!transmission.contains('manual')) {
            return false;
          }
        }
        
        return true;
      }).toList();
      
      // Apply sorting
      if (_sortBy == 'Cheapest') {
        _filteredCarResults.sort((a, b) {
          final aMap = a is Map<String, dynamic> ? a : Map<String, dynamic>.from(a as Map);
          final bMap = b is Map<String, dynamic> ? b : Map<String, dynamic>.from(b as Map);
          
          final aPrice = (aMap['calculated_cost'] ?? aMap['cost_per_day'] ?? aMap['cost'] ?? 0).toDouble();
          final bPrice = (bMap['calculated_cost'] ?? bMap['cost_per_day'] ?? bMap['cost'] ?? 0).toDouble();
          
          return aPrice.compareTo(bPrice);
        });
      } else if (_sortBy == 'Most Expensive') {
        _filteredCarResults.sort((a, b) {
          final aMap = a is Map<String, dynamic> ? a : Map<String, dynamic>.from(a as Map);
          final bMap = b is Map<String, dynamic> ? b : Map<String, dynamic>.from(b as Map);
          
          final aPrice = (aMap['calculated_cost'] ?? aMap['cost_per_day'] ?? aMap['cost'] ?? 0).toDouble();
          final bPrice = (bMap['calculated_cost'] ?? bMap['cost_per_day'] ?? bMap['cost'] ?? 0).toDouble();
          
          return bPrice.compareTo(aPrice);
        });
      }
    });
  }

  void _showFilterBottomSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => StatefulBuilder(
        builder: (context, setModalState) {
          return Container(
            height: MediaQuery.of(context).size.height * 0.9,
            decoration: const BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
            ),
            child: Column(
              children: [
                // Header
                Padding(
                  padding: const EdgeInsets.all(16),
                  child: Row(
                    children: [
                      IconButton(
                        icon: const Icon(Icons.arrow_back_ios, size: 20),
                        onPressed: () => Navigator.pop(context),
                        padding: EdgeInsets.zero,
                        constraints: const BoxConstraints(),
                      ),
                      const SizedBox(width: 8),
                      const Text(
                        'Filter By:',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ],
                  ),
                ),

                const Divider(height: 1),

                // Filters content
                Expanded(
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      // Sort By
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text(
                            'Sort By',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                              color: Color(0xFFD32F2F),
                            ),
                          ),
                          GestureDetector(
                            onTap: () {
                              showMenu(
                                context: context,
                                position: RelativeRect.fromLTRB(
                                  MediaQuery.of(context).size.width - 200,
                                  150,
                                  20,
                                  0,
                                ),
                                items: [
                                  PopupMenuItem(
                                    value: 'Cheapest',
                                    child: Text('Cheapest'),
                                  ),
                                  PopupMenuItem(
                                    value: 'Most Expensive',
                                    child: Text('Most Expensive'),
                                  ),
                                ],
                              ).then((value) {
                                if (value != null) {
                                  setModalState(() {
                                    _sortBy = value;
                                  });
                                }
                              });
                            },
                            child: Container(
                              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                              decoration: BoxDecoration(
                                border: Border.all(color: Colors.grey.shade300),
                                borderRadius: BorderRadius.circular(8),
                              ),
                              child: Row(
                                children: [
                                  Text(
                                    _sortBy,
                                    style: const TextStyle(fontSize: 14),
                                  ),
                                  const SizedBox(width: 8),
                                  const Icon(Icons.arrow_drop_down, size: 20),
                                ],
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Car Type
                      const Text(
                        'Car Type',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Economy', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('SUV', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Van', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('Luxury', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Rental Companies
                      const Text(
                        'Rental Companies',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Hertz', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('Sixt', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Avis', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('Budget', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Features & Inclusions
                      const Text(
                        'Features & Inclusions',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Automatic Transmission', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Manual Transmission', setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Air Conditioning', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Unlimited Mileage', setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Free Cancellation', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Child Seat Available', setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('GPS Included', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Additional Driver Included', setModalState),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),
                    ],
                  ),
                ),

                // Buttons (Clear Filters and Done)
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
                  child: Row(
                    children: [
                      Expanded(
                        child: SizedBox(
                          height: 50,
                          child: OutlinedButton(
                            onPressed: () {
                              setModalState(() {
                                _sortBy = 'Cheapest';
                                _selectedCarType = '';
                                _selectedRentalCompany = '';
                                _features = {
                                  'Automatic Transmission': false,
                                  'Manual Transmission': false,
                                  'Air Conditioning': false,
                                  'Unlimited Mileage': false,
                                  'Free Cancellation': false,
                                  'Child Seat Available': false,
                                  'GPS Included': false,
                                  'Additional Driver Included': false,
                                };
                              });
                            },
                            style: OutlinedButton.styleFrom(
                              side: const BorderSide(color: Color(0xFF1e5a8e), width: 2),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(25),
                              ),
                            ),
                            child: const Text(
                              'Clear Filters',
                              style: TextStyle(
                                color: Color(0xFF1e5a8e),
                                fontSize: 16,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: SizedBox(
                          height: 50,
                          child: ElevatedButton(
                            onPressed: () {
                              Navigator.pop(context);
                              _applyFilters();
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
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildRadioOption(String label, String groupValue, Function(String?) onChanged) {
    return Row(
      children: [
        Radio<String>(
          value: label,
          groupValue: groupValue,
          onChanged: onChanged,
          activeColor: const Color(0xFF1e5a8e),
        ),
        Expanded(
          child: Text(
            label,
            style: const TextStyle(fontSize: 14),
          ),
        ),
      ],
    );
  }

  Widget _buildCheckboxOption(String label, StateSetter setModalState) {
    return Row(
      children: [
        Checkbox(
          value: _features[label] ?? false,
          onChanged: (value) {
            setModalState(() {
              _features[label] = value ?? false;
            });
          },
          activeColor: const Color(0xFF1e5a8e),
        ),
        Expanded(
          child: Text(
            label,
            style: const TextStyle(fontSize: 12),
          ),
        ),
      ],
    );
  }
}

