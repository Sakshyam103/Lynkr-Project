import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, Alert } from 'react-native';
import { apiService } from '../services/api';

export default function CreateSurveyScreen({ route, navigation }: any) {
  const { eventId, eventName } = route.params;
  const [surveyData, setSurveyData] = useState({
    title: '',
    description: '',
    rewardPoints: '50',
    questions: [{ question: '', type: 'multiple_choice', options: [''] }]
  });

  const createSurvey = async () => {
    try {
      await apiService.scheduleSurveys(eventId, {
        title: surveyData.title,
        description: surveyData.description,
        questions: surveyData.questions,
        rewardPoints: parseInt(surveyData.rewardPoints)
      });
      
      Alert.alert('Success', 'Survey created successfully!', [
        { text: 'OK', onPress: () => navigation.goBack() }
      ]);
    } catch (error) {
      Alert.alert('Error', 'Failed to create survey');
    }
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Create Survey</Text>
        <Text style={styles.subtitle}>For: {eventName}</Text>
        <Text style={styles.eventId}>Event ID: {eventId}</Text>
      </View>
      
      {/* Your survey creation form here */}
      
      <TouchableOpacity style={styles.createButton} onPress={createSurvey}>
        <Text style={styles.createButtonText}>Create Survey</Text>
      </TouchableOpacity>
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
    marginBottom: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginBottom: 4,
  },
  eventId: {
    fontSize: 14,
    color: '#999',
  },
  createButton: {
    backgroundColor: '#6f42c1',
    margin: 20,
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
  },
  createButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
