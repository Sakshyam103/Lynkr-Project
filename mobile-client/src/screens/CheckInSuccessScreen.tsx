/**
 * Check-in Success Screen - Shows success confirmation
 */

import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';

export default function CheckInSuccessScreen({ route, navigation }: any) {
  const { event, isCheckOut } = route.params;

  return (
    <View style={styles.container}>
      <View style={styles.successCard}>
        <Text style={styles.successIcon}>
          {isCheckOut ? '🚪' : '✅'}
        </Text>
        
        <Text style={styles.successTitle}>
          {isCheckOut ? 'Checked Out Successfully!' : 'Checked In Successfully!'}
        </Text>
        
        <Text style={styles.eventName}>{event.name}</Text>
        
        <Text style={styles.successMessage}>
          {isCheckOut 
            ? 'Thank you for attending the event!' 
            : 'Welcome to the event! Enjoy your experience.'
          }
        </Text>
        
        <TouchableOpacity 
          style={styles.continueButton}
          onPress={() => navigation.navigate('Events')}
        >
          <Text style={styles.continueText}>Continue</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  successCard: {
    backgroundColor: '#fff',
    borderRadius: 20,
    padding: 40,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 12,
    elevation: 4,
    width: '100%',
    maxWidth: 350,
  },
  successIcon: {
    fontSize: 80,
    marginBottom: 20,
  },
  successTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#28a745',
    textAlign: 'center',
    marginBottom: 16,
  },
  eventName: {
    fontSize: 18,
    fontWeight: '600',
    color: '#333',
    textAlign: 'center',
    marginBottom: 16,
  },
  successMessage: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
    lineHeight: 22,
    marginBottom: 30,
  },
  continueButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 40,
    paddingVertical: 16,
    borderRadius: 12,
    width: '100%',
  },
  continueText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
    textAlign: 'center',
  },
});