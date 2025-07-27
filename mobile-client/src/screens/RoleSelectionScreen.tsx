import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';

export default function RoleSelectionScreen({ navigation }: any) {
  const selectRole = (role: string) => {
    navigation.navigate('Login', { userType: role });
  };

  return (
    <View style={styles.container}>
      <View style={styles.logoContainer}>
        <Text style={styles.logo}>Lynkr</Text>
        <Text style={styles.tagline}>Make Sponsorships Smarter</Text>
      </View>

      <View style={styles.content}>
        <Text style={styles.title}>Welcome to Lynkr</Text>
        <Text style={styles.subtitle}>How would you like to continue?</Text>

        <TouchableOpacity 
          style={styles.roleCard}
          onPress={() => selectRole('user')}
        >
          <Text style={styles.roleIcon}>👤</Text>
          <Text style={styles.roleTitle}>I'm a User</Text>
          <Text style={styles.roleDescription}>
            Attend events, create content, and earn rewards
          </Text>
        </TouchableOpacity>

        <TouchableOpacity 
          style={styles.roleCard}
          onPress={() => selectRole('brand')}
        >
          <Text style={styles.roleIcon}>🏢</Text>
          <Text style={styles.roleTitle}>I'm a Brand</Text>
          <Text style={styles.roleDescription}>
            Sponsor events, view analytics, and manage campaigns
          </Text>
        </TouchableOpacity>

        {/*<TouchableOpacity */}
        {/*  style={styles.roleCard}*/}
        {/*  onPress={() => selectRole('organization')}*/}
        {/*>*/}
        {/*  <Text style={styles.roleIcon}>🎪</Text>*/}
        {/*  <Text style={styles.roleTitle}>I'm an Organization</Text>*/}
        {/*  <Text style={styles.roleDescription}>*/}
        {/*    Create and manage events for sponsors*/}
        {/*  </Text>*/}
        {/*</TouchableOpacity>*/}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  logoContainer: {
    alignItems: 'center',
    paddingTop: 80,
    paddingBottom: 40,
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
  content: {
    flex: 1,
    padding: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    textAlign: 'center',
    color: '#333',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 18,
    textAlign: 'center',
    color: '#666',
    marginBottom: 40,
  },
  roleCard: {
    backgroundColor: '#fff',
    padding: 24,
    borderRadius: 16,
    marginBottom: 20,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  roleIcon: {
    fontSize: 48,
    marginBottom: 16,
  },
  roleTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#333',
  },
  roleDescription: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
    lineHeight: 22,
  },
});