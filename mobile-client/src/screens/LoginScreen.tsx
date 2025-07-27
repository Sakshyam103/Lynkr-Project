import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, StyleSheet, Alert, ActivityIndicator } from 'react-native';
import { useDispatch } from 'react-redux';
import { login } from '../store/slices/authSlice';
import { apiService } from '../services/api';

export default function LoginScreen({ route, navigation }: any) {
  const userType = route.params?.userType || 'user';
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const dispatch = useDispatch();

  const getRoleInfo = () => {
    switch (userType) {
      case 'brand':
        return {
          title: 'Brand Login',
          subtitle: 'Access your brand dashboard',
          icon: '🏢',
          color: '#28a745'
        };
      case 'organization':
        return {
          title: 'Organization Login',
          subtitle: 'Manage your events',
          icon: '🎪',
          color: '#ffc107'
        };
      default:
        return {
          title: 'User Login',
          subtitle: 'Join events and earn rewards',
          icon: '👤',
          color: '#007AFF'
        };
    }
  };

  const handleLogin = async () => {
    if (!email || !password) {
      Alert.alert('Error', 'Please fill in all fields');
      return;
    }

    setLoading(true);
    try {
      let response;
      
      // Call appropriate login endpoint based on user type
      switch (userType) {
        case 'brand':
          response = await apiService.brandLogin({ email, password });
          break;
        case 'organization':
          // Assuming organizations use the same endpoint as users for now
          response = await apiService.login({ email, password });
          break;
        default:
          response = await apiService.login({ email, password });
      }
      
      // Set token and user role in API service
      apiService.setToken(response.token);
      apiService.setUserRole(userType);
      
      // Dispatch login action with role
      dispatch(login({ 
        user: { ...response.user, role: userType }, 
        token: response.token 
      }));
      
    } catch (error) {
      Alert.alert('Error', 'Login failed. Please check your credentials.');
    } finally {
      setLoading(false);
    }
  };

  const roleInfo = getRoleInfo();

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity 
          style={styles.backButton}
          onPress={() => navigation.goBack()}
        >
          <Text style={styles.backText}>← Back</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.logoContainer}>
        <Text style={[styles.roleIcon, { color: roleInfo.color }]}>{roleInfo.icon}</Text>
        <Text style={styles.logo}>Lynkr</Text>
        <Text style={styles.tagline}>Make Sponsorships Smarter</Text>
      </View>

      <View style={styles.formContainer}>
        <Text style={styles.title}>{roleInfo.title}</Text>
        <Text style={styles.subtitle}>{roleInfo.subtitle}</Text>

        <TextInput
          style={styles.input}
          placeholder="Email"
          value={email}
          onChangeText={setEmail}
          keyboardType="email-address"
          autoCapitalize="none"
          placeholderTextColor="#999"
        />

        <TextInput
          style={styles.input}
          placeholder="Password"
          value={password}
          onChangeText={setPassword}
          secureTextEntry
          placeholderTextColor="#999"
        />

        <TouchableOpacity 
          style={[styles.button, { backgroundColor: roleInfo.color }, loading && styles.buttonDisabled]} 
          onPress={handleLogin}
          disabled={loading}
        >
          {loading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.buttonText}>Sign In</Text>
          )}
        </TouchableOpacity>

        <TouchableOpacity onPress={() => navigation.navigate('Register', { userType })}>
          <Text style={[styles.linkText, { color: roleInfo.color }]}>
            Don't have an account? Sign up as {userType}
          </Text>
        </TouchableOpacity>

        <View style={styles.roleIndicator}>
          <Text style={styles.roleIndicatorText}>
            Signing in as: <Text style={[styles.roleType, { color: roleInfo.color }]}>{userType.toUpperCase()}</Text>
          </Text>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    paddingTop: 50,
    paddingHorizontal: 20,
  },
  backButton: {
    alignSelf: 'flex-start',
  },
  backText: {
    fontSize: 16,
    color: '#007AFF',
  },
  logoContainer: {
    alignItems: 'center',
    paddingTop: 40,
    paddingBottom: 40,
  },
  roleIcon: {
    fontSize: 48,
    marginBottom: 16,
  },
  logo: {
    fontSize: 48,
    fontWeight: 'bold',
    color: '#007AFF',
    marginBottom: 8,
  },
  tagline: {
    fontSize: 16,
    color: '#666',
    fontStyle: 'italic',
  },
  formContainer: {
    flex: 1,
    backgroundColor: '#fff',
    borderTopLeftRadius: 30,
    borderTopRightRadius: 30,
    padding: 30,
    paddingTop: 40,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginBottom: 30,
  },
  input: {
    borderWidth: 1,
    borderColor: '#e0e0e0',
    padding: 16,
    borderRadius: 12,
    marginBottom: 16,
    fontSize: 16,
    backgroundColor: '#f8f9fa',
  },
  button: {
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
    marginBottom: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 4,
  },
  buttonDisabled: {
    opacity: 0.7,
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  linkText: {
    textAlign: 'center',
    fontSize: 16,
    marginBottom: 20,
  },
  roleIndicator: {
    backgroundColor: '#f8f9fa',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  roleIndicatorText: {
    fontSize: 14,
    color: '#666',
  },
  roleType: {
    fontWeight: 'bold',
  },
});