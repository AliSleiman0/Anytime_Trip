import 'package:firebase_auth/firebase_auth.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:get/get.dart';
import '../core/network/api_client.dart';
import '../core/network/endpoints.dart';
import '../core/storage/storage_service.dart';

class AuthService {
  final FirebaseAuth _auth = FirebaseAuth.instance;
  final ApiClient _apiClient = ApiClient();
  final StorageService _storage = StorageService();

  /// Signs in (or signs up) the user using Google on Android.
  /// Returns a map with user data and isNewUser flag.
  Future<Map<String, dynamic>?> signInWithGoogle() async {
    // Opens Google account picker
    final GoogleSignInAccount? googleUser = await GoogleSignIn().signIn();

    if (googleUser == null) {
      // user cancelled
      return null;
    }

    final GoogleSignInAuthentication googleAuth =
        await googleUser.authentication;

    final credential = GoogleAuthProvider.credential(
      idToken: googleAuth.idToken,
      accessToken: googleAuth.accessToken,
    );

    final UserCredential userCred =
        await _auth.signInWithCredential(credential);

    final User? firebaseUser = userCred.user;
    
    if (firebaseUser == null) return null;

    // Send Google user data to backend
    final response = await _apiClient.post(
      Endpoints.googleSignIn,
      data: {
        'uid': firebaseUser.uid,
        'email': firebaseUser.email,
        'name': firebaseUser.displayName,
        'photo_url': firebaseUser.photoURL,
      },
    );

    // Backend returns: { message, token, user: {...}, isNewUser: true/false }
    final userData = response.data['user'] as Map<String, dynamic>;
    final token = response.data['token'] as String;
    final isNewUser = response.data['isNewUser'] as bool? ?? false;
    
    // Save token and user data
    await _storage.saveToken(token);
    await _storage.saveUser(userData);

    return {
      'user': firebaseUser,
      'isNewUser': isNewUser,
      'userData': userData,
    };
  }

  Future<void> signOut() async {
    await GoogleSignIn().signOut();
    await _auth.signOut();
    await _storage.clearAll();
  }
}
