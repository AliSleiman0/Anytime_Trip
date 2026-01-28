import 'package:flutter/material.dart';
import 'flight_checkout_screen.dart';
import '../../../core/widgets/unified_ui_components.dart';

class FlightDetailsScreen extends StatelessWidget {
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

  const FlightDetailsScreen({
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
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: UnifiedAppBar(
        title: 'Flight Details',
        onBackPressed: () => Navigator.pop(context),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Center(
              child: Text(
                '$from TO $to',
                style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e)),
              ),
            ),
            const SizedBox(height: 4),
            Center(
              child: Text(
                tripType == 'roundtrip' ? '1 traveler' : '1 traveler',
                style: const TextStyle(fontSize: 12, color: Colors.grey),
              ),
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Flexible(child: _buildFlightCard(from, to, departureTime, arrivalTime, airline, stops, departureDate, 'Departure')),
                if (tripType.toLowerCase().trim() == 'roundtrip') ...[
                  Flexible(
                    child: Column(
                      children: [
                        Image.asset(
                          'assets/images/fl1.png',
                          width: 40,
                          height: 40,
                          errorBuilder: (context, error, stackTrace) {
                            return const Icon(Icons.arrow_forward, size: 24, color: Color(0xFF1e5a8e));
                          },
                        ),
                        const SizedBox(height: 8),
                        Image.asset(
                          'assets/images/fl2.png',
                          width: 40,
                          height: 40,
                          errorBuilder: (context, error, stackTrace) {
                            return const Icon(Icons.arrow_back, size: 24, color: Color(0xFF1e5a8e));
                          },
                        ),
                      ],
                    ),
                  ),
                  Flexible(child: _buildFlightCard(to, from, returnDepartureTime ?? '', returnArrivalTime ?? '', returnAirline ?? '', returnStops ?? '', returnDate ?? '', 'Return')),
                ],
              ],
            ),
            const SizedBox(height: 16),
            _buildFareSection(),
            const SizedBox(height: 8),
            const Text('Seats', style: TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e), fontSize: 18)),
            _buildSeatsSection(),
            const SizedBox(height: 8),
            const Text('Bags', style: TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e), fontSize: 18)),
            _buildBagsSection(),
            const Spacer(),
            Center(
              child: Column(
                children: [
                  Text('Total Price: ,${price}', style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Color(0xFFD32F2F))),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => FlightCheckoutScreen(
                              from: from,
                              to: to,
                              departureTime: departureTime,
                              arrivalTime: arrivalTime,
                              airline: airline,
                              stops: stops,
                              tripType: tripType,
                              price: price,
                              departureDate: departureDate,
                              returnDate: returnDate,
                              returnDepartureTime: returnDepartureTime,
                              returnArrivalTime: returnArrivalTime,
                              returnAirline: returnAirline,
                              returnStops: returnStops,
                            ),
                          ),
                        );
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1e5a8e),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                      child: const Text('Checkout', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Align(
                    alignment: Alignment.bottomLeft,
                    child: Opacity(
                      opacity: 0.3,
                      child: Image.asset('assets/images/plane.png', width: 60),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFlightCard(String from, String to, String dep, String arr, String airline, String stops, String date, String label) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 4),
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: Colors.white,
        border: Border.all(color: Color(0xFF1e5a8e), width: 1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('$from to $to', style: const TextStyle(fontWeight: FontWeight.w600)),
          Text(date, style: const TextStyle(fontSize: 12, color: Color(0xFF1e5a8e), fontWeight: FontWeight.w600)),
          Text('$dep - $arr', style: const TextStyle(fontSize: 12)),
          Text(airline, style: const TextStyle(fontSize: 12)),
          Text(stops, style: const TextStyle(fontSize: 12)),
          Text('Change Flight', style: const TextStyle(fontSize: 12, color: Color(0xFF1e5a8e), decoration: TextDecoration.underline)),
        ],
      ),
    );
  }

  Widget _buildFareSection() {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Your Fare: Basic', style: TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF1e5a8e))),
          Row(
            children: const [
              Icon(Icons.check_circle, color: Colors.green, size: 16),
              SizedBox(width: 4),
              Text('Seat choice allowed'),
              SizedBox(width: 16),
              Icon(Icons.check_circle, color: Colors.green, size: 16),
              SizedBox(width: 4),
              Text('Carry-on bag included (8kg)'),
            ],
          ),
          Row(
            children: const [
              Icon(Icons.cancel, color: Colors.red, size: 16),
              SizedBox(width: 4),
              Text('Non-refundable'),
              SizedBox(width: 16),
              Icon(Icons.cancel, color: Colors.red, size: 16),
              SizedBox(width: 4),
              Text('Changes not allowed'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSeatsSection() {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: const [
          Icon(Icons.check_circle, color: Colors.green, size: 16),
          SizedBox(width: 4),
          Text('Seat choice allowed'),
          Spacer(),
          Text('Choose Seat', style: TextStyle(color: Color(0xFF1e5a8e), decoration: TextDecoration.underline)),
        ],
      ),
    );
  }

  Widget _buildBagsSection() {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: const [
          Row(
            children: [
              Icon(Icons.check_circle, color: Colors.green, size: 16),
              SizedBox(width: 4),
              Text('Personal item included'),
            ],
          ),
          Row(
            children: [
              Icon(Icons.check_circle, color: Colors.green, size: 16),
              SizedBox(width: 4),
              Text('Carry-on bag included (8kg)'),
            ],
          ),
          Row(
            children: [
              Icon(Icons.check_circle, color: Colors.green, size: 16),
              SizedBox(width: 4),
              Text('1st checked bag included (20kg)'),
            ],
          ),
        ],
      ),
    );
  }
}
