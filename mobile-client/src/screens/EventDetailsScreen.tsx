import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, Alert, Modal } from 'react-native';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { apiService } from '../services/api';
import * as Location from 'expo-location';



  export default function EventDetailsScreen({route, navigation}: any) {
    const {event} = route.params;
    const {user} = useSelector((state: RootState) => state.auth);
    const [loading, setLoading] = useState(false);
    const [checkedIn, setCheckedIn] = useState(false);
    const [errorModal, setErrorModal] = useState({visible: false, message: ''});

    const isUser = user?.role === 'user' || !user?.role;

    const handleCheckInOut = async () => {
      if (!isUser) {
        Alert.alert('Access Denied', 'Only users can check into events');
        return;
      }

      setLoading(true);
      try {
        if (checkedIn) {
          await apiService.checkOut(event.id);
          setCheckedIn(false);
          navigation.navigate('CheckInSuccess', {event: event, isCheckOut: true});
        } else {
          const {status} = await Location.requestForegroundPermissionsAsync();
          if (status !== 'granted') {
            Alert.alert('Permission denied', 'Location permission is required to check in');
            return;
          }


          const locationEnabled = await Location.hasServicesEnabledAsync();
          if (!locationEnabled) {
            setErrorModal({ visible: true, message: 'Please enable location services in your device settings.' });
            return;
          }


          const location = await Location.getCurrentPositionAsync({
            accuracy: Location.Accuracy.Balanced, // Less strict than High
          }).catch(async () => {
            // Fallback to last known location
            return await Location.getLastKnownPositionAsync({
              maxAge: 300000, // 5 minutes
            });
          });

          if (!location) {
            setErrorModal({ visible: true, message: 'Unable to get your location. Please enable GPS and try again.' });
            return;
          }

          await apiService.checkIn(event.id, {
            latitude: location.coords.latitude,
            longitude: location.coords.longitude
          });
          setCheckedIn(true);
          navigation.navigate('CheckInSuccess', {event: event});
        }
      } catch (error: any) {
        let errorMessage = `Failed to ${checkedIn ? 'check out' : 'check in'}`;

        if (error.message) {
          if (error.message.includes('400')) {
            errorMessage = 'You must be at the event location to check in';
          } else if (error.message.includes('500')) {
            errorMessage = 'Server error occurred. Please try again later.';
          } else if (error.message.includes('event has not started yet')) {
            errorMessage = 'Event has not started yet';
          } else if (error.message.includes('event has already ended')) {
            errorMessage = 'Event has already ended';
          }
        }

        setErrorModal({visible: true, message: errorMessage});
      } finally {
        setLoading(false);
      }
    };

    const formatDate = (dateString: string) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      });
    };

    return (
        <ScrollView style={styles.container}>
          <View style={styles.header}>
            <Text style={styles.eventName}>{event.name}</Text>
            <Text style={styles.eventDescription}>{event.description}</Text>

            <View style={styles.eventMeta}>
              <Text style={styles.metaItem}>📅 {formatDate(event.startDate)}</Text>
              <Text style={styles.metaItem}>📍 {event.location.address}</Text>
              <Text style={styles.metaItem}>👥 {event.attendeeCount} attending</Text>
            </View>
          </View>

          <View style={styles.actionsContainer}>
            {isUser && (
                <>
                  <TouchableOpacity
                      style={[styles.actionButton, checkedIn ? styles.checkOutButton : styles.primaryButton]}
                      onPress={handleCheckInOut}
                      disabled={loading}
                  >
                    <Text style={styles.primaryButtonText}>
                      {loading ? (checkedIn ? 'Checking out...' : 'Checking in...') : checkedIn ? '🚪 Check Out' : '📍 Check In'}
                    </Text>
                  </TouchableOpacity>

                  <TouchableOpacity
                      style={[styles.actionButton, styles.secondaryButton]}
                      onPress={() => navigation.navigate('CreateContent', {eventId: event.id})}
                  >
                    <Text style={styles.secondaryButtonText}>📸 Create Content</Text>
                  </TouchableOpacity>
                </>
            )}
          </View>

          <Modal
              visible={errorModal.visible}
              transparent={true}
              animationType="fade"
          >
            <View style={styles.modalOverlay}>
              <View style={styles.modalContent}>
                <Text style={styles.modalTitle}>Check-in Failed</Text>
                <Text style={styles.modalMessage}>{errorModal.message}</Text>
                <TouchableOpacity
                    style={styles.modalButton}
                    onPress={() => setErrorModal({visible: false, message: ''})}
                >
                  <Text style={styles.modalButtonText}>OK</Text>
                </TouchableOpacity>
              </View>
            </View>
          </Modal>
        </ScrollView>
    );
  }

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: '#f8f9fa',
    },
    header: {
      backgroundColor: '#fff',
      padding: 20,
    },
    eventName: {
      fontSize: 24,
      fontWeight: 'bold',
      color: '#333',
      marginBottom: 8,
    },
    eventDescription: {
      fontSize: 16,
      color: '#666',
      lineHeight: 22,
      marginBottom: 16,
    },
    eventMeta: {
      gap: 8,
    },
    metaItem: {
      fontSize: 14,
      color: '#666',
      marginBottom: 4,
    },
    actionsContainer: {
      flexDirection: 'row',
      padding: 20,
      gap: 12,
    },
    actionButton: {
      flex: 1,
      padding: 16,
      borderRadius: 12,
      alignItems: 'center',
    },
    primaryButton: {
      backgroundColor: '#007AFF',
    },
    primaryButtonText: {
      color: '#fff',
      fontSize: 16,
      fontWeight: '600',
    },
    secondaryButton: {
      backgroundColor: 'transparent',
      borderWidth: 1,
      borderColor: '#007AFF',
    },
    secondaryButtonText: {
      color: '#007AFF',
      fontSize: 16,
      fontWeight: '600',
    },
    checkOutButton: {
      backgroundColor: '#dc3545',
    },
    modalOverlay: {
      flex: 1,
      backgroundColor: 'rgba(0, 0, 0, 0.5)',
      justifyContent: 'center',
      alignItems: 'center',
    },
    modalContent: {
      backgroundColor: '#fff',
      borderRadius: 12,
      padding: 20,
      margin: 20,
      alignItems: 'center',
      minWidth: 280,
    },
    modalTitle: {
      fontSize: 18,
      fontWeight: 'bold',
      color: '#dc3545',
      marginBottom: 12,
    },
    modalMessage: {
      fontSize: 16,
      color: '#333',
      textAlign: 'center',
      marginBottom: 20,
      lineHeight: 22,
    },
    modalButton: {
      backgroundColor: '#007AFF',
      paddingHorizontal: 24,
      paddingVertical: 12,
      borderRadius: 8,
    },
    modalButtonText: {
      color: '#fff',
      fontSize: 16,
      fontWeight: '600',
    },
  });
