import React, { useState } from 'react';
import { View, Text, ScrollView, TextInput, TouchableOpacity, StyleSheet, Alert } from 'react-native';
import { apiService } from '../services/api';

export default function CreateEventScreen({ navigation }: any) {
  const [eventData, setEventData] = useState({
    name: '',
    description: '',
    startDate: '',
    endDate: '',
    address: '',
    latitude: '',
    longitude: '',
  });
  const [loading, setLoading] = useState(false);

  const createEvent = async () => {
    if (!eventData.name || !eventData.description || !eventData.address) {
      Alert.alert('Error', 'Please fill in all required fields');
      return;
    }

    setLoading(true);
    try {
      await apiService.createEvent({
        ...eventData,
        location: {
          latitude: parseFloat(eventData.latitude) || 37.7749,
          longitude: parseFloat(eventData.longitude) || -122.4194,
          address: eventData.address,
        },
      });
      
      Alert.alert('Success', 'Event created successfully!', [
        { text: 'OK', onPress: () => navigation.goBack() }
      ]);
    } catch (error) {
      Alert.alert('Error', 'Failed to create event');
    } finally {
      setLoading(false);
    }
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Create New Event</Text>
        <Text style={styles.subtitle}>Set up your sponsored event</Text>
      </View>

      <View style={styles.form}>
        <Text style={styles.label}>Event Name *</Text>
        <TextInput
          style={styles.input}
          placeholder="Enter event name"
          value={eventData.name}
          onChangeText={(text) => setEventData({...eventData, name: text})}
        />

        <Text style={styles.label}>Description *</Text>
        <TextInput
          style={[styles.input, styles.textArea]}
          placeholder="Describe your event"
          value={eventData.description}
          onChangeText={(text) => setEventData({...eventData, description: text})}
          multiline
          numberOfLines={4}
        />

        <Text style={styles.label}>Start Date & Time</Text>
        <TextInput
          style={styles.input}
          placeholder="2024-04-15T10:00:00Z"
          value={eventData.startDate}
          onChangeText={(text) => setEventData({...eventData, startDate: text})}
        />

        <Text style={styles.label}>End Date & Time</Text>
        <TextInput
          style={styles.input}
          placeholder="2024-04-15T18:00:00Z"
          value={eventData.endDate}
          onChangeText={(text) => setEventData({...eventData, endDate: text})}
        />

        <Text style={styles.label}>Location Address *</Text>
        <TextInput
          style={styles.input}
          placeholder="123 Main St, San Francisco, CA"
          value={eventData.address}
          onChangeText={(text) => setEventData({...eventData, address: text})}
        />

        <View style={styles.row}>
          <View style={styles.halfInput}>
            <Text style={styles.label}>Latitude</Text>
            <TextInput
              style={styles.input}
              placeholder="37.7749"
              value={eventData.latitude}
              onChangeText={(text) => setEventData({...eventData, latitude: text})}
              keyboardType="numeric"
            />
          </View>
          <View style={styles.halfInput}>
            <Text style={styles.label}>Longitude</Text>
            <TextInput
              style={styles.input}
              placeholder="-122.4194"
              value={eventData.longitude}
              onChangeText={(text) => setEventData({...eventData, longitude: text})}
              keyboardType="numeric"
            />
          </View>
        </View>

        <TouchableOpacity 
          style={[styles.createButton, loading && styles.disabled]}
          onPress={createEvent}
          disabled={loading}
        >
          <Text style={styles.createButtonText}>
            {loading ? 'Creating...' : 'Create Event'}
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#333',
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginTop: 4,
  },
  form: {
    padding: 20,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 8,
    marginTop: 16,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    backgroundColor: '#fff',
  },
  textArea: {
    height: 100,
    textAlignVertical: 'top',
  },
  row: {
    flexDirection: 'row',
    gap: 12,
  },
  halfInput: {
    flex: 1,
  },
  createButton: {
    backgroundColor: '#28a745',
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
    marginTop: 24,
  },
  createButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  disabled: {
    opacity: 0.6,
  },
});