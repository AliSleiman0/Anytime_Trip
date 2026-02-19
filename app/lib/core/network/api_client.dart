import 'package:dio/dio.dart';
import 'endpoints.dart';
import '../storage/storage_service.dart';

class ApiClient {
  late Dio _dio;
  final StorageService _storage = StorageService();
  
  static final ApiClient _instance = ApiClient._internal();
  factory ApiClient() => _instance;

  ApiClient._internal() {
    _dio = Dio(
      BaseOptions(
        baseUrl: Endpoints.baseUrl,
        connectTimeout: const Duration(seconds: 30),
        receiveTimeout: const Duration(seconds: 30),
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );

    // Add interceptors
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          // Add auth token if available
          final token = await _storage.getToken();
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          print('[ApiClient] Request: ${options.method} ${options.baseUrl}${options.path}');
          print('[ApiClient] Has token: ${token != null}');
          return handler.next(options);
        },
        onResponse: (response, handler) {
          print('[ApiClient] Response: ${response.statusCode} - ${response.requestOptions.path}');
          return handler.next(response);
        },
        onError: (error, handler) {
          print('[ApiClient] Error interceptor triggered');
          print('[ApiClient] Error type: ${error.type}');
          print('[ApiClient] Error message: ${error.message}');
          print('[ApiClient] Response status: ${error.response?.statusCode}');
          print('[ApiClient] Response data: ${error.response?.data}');
          
          // Extract error message from response
          if (error.response?.data != null) {
            final data = error.response!.data;
            if (data is Map<String, dynamic> && data['error'] != null) {
              print('[ApiClient] Throwing error from response: ${data['error']}');
              throw Exception(data['error']);
            }
          }
          return handler.next(error);
        },
      ),
    );
  }

  // GET request
  Future<Response> get(String path, {Map<String, dynamic>? queryParameters}) async {
    try {
      return await _dio.get(path, queryParameters: queryParameters);
    } catch (e) {
      rethrow;
    }
  }

  // POST request
  Future<Response> post(String path, {dynamic data}) async {
    try {
      return await _dio.post(path, data: data);
    } catch (e) {
      rethrow;
    }
  }

  // PUT request
  Future<Response> put(String path, {dynamic data}) async {
    try {
      return await _dio.put(path, data: data);
    } catch (e) {
      rethrow;
    }
  }

  // DELETE request
  Future<Response> delete(String path) async {
    try {
      return await _dio.delete(path);
    } catch (e) {
      rethrow;
    }
  }
}
