import React, { useState } from 'react';
import {
  View,
  Text,
  Modal,
  TouchableOpacity,
  StyleSheet,
  Alert
} from 'react-native';
import { apiService } from '../services/api';

interface FeedbackModalProps {
  visible: boolean;
  onClose: () => void;
  eventId: string;
}

export default function FeedbackModal({ visible, onClose, eventId }: FeedbackModalProps) {
  const [rating, setRating] = useState(3);
  const [submitting, setSubmitting] = useState(false);

  const submitFeedback = async () => {
    setSubmitting(true);
    try {
      await apiService.submitSliderFeedback({
        eventId,
        category: 'overall_experience',
        rating,
        maxRating: 5,
      });
      Alert.alert('Success', 'Thank you for your feedback!');
      onClose();
    } catch (error) {
      Alert.alert('Error', 'Failed to submit feedback');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal visible={visible} transparent animationType="slide">
      <View style={styles.overlay}>
        <View style={styles.modal}>
          <Text style={styles.title}>Rate Your Experience</Text>
          <Text style={styles.subtitle}>How would you rate this event?</Text>
          
          <View style={styles.sliderContainer}>
            <Text style={styles.ratingText}>{rating.toFixed(1)} / 5.0</Text>
            {/*<Slider*/}
            {/*  style={styles.slider}*/}
            {/*  minimumValue={1}*/}
            {/*  maximumValue={5}*/}
            {/*  value={rating}*/}
            {/*  onValueChange={setRating}*/}
            {/*  step={0.1}*/}
            {/*  minimumTrackTintColor="#007AFF"*/}
            {/*  maximumTrackTintColor="#ddd"*/}
            {/*  thumbStyle={{ backgroundColor: '#007AFF' }}*/}
            {/*/>*/}
            <View style={styles.labels}>
              <Text style={styles.label}>Poor</Text>
              <Text style={styles.label}>Excellent</Text>
            </View>
          </View>

          <View style={styles.buttons}>
            <TouchableOpacity style={styles.cancelButton} onPress={onClose}>
              <Text style={styles.cancelText}>Cancel</Text>
            </TouchableOpacity>
            <TouchableOpacity 
              style={[styles.submitButton, submitting && styles.disabled]} 
              onPress={submitFeedback}
              disabled={submitting}
            >
              <Text style={styles.submitText}>
                {submitting ? 'Submitting...' : 'Submit'}
              </Text>
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  modal: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 24,
    width: '90%',
    maxWidth: 400,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
    marginBottom: 24,
  },
  sliderContainer: {
    marginBottom: 24,
  },
  ratingText: {
    fontSize: 24,
    fontWeight: 'bold',
    textAlign: 'center',
    color: '#007AFF',
    marginBottom: 16,
  },
  slider: {
    width: '100%',
    height: 40,
  },
  labels: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: 8,
  },
  label: {
    fontSize: 12,
    color: '#666',
  },
  buttons: {
    flexDirection: 'row',
    gap: 12,
  },
  cancelButton: {
    flex: 1,
    padding: 12,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#ddd',
    alignItems: 'center',
  },
  cancelText: {
    color: '#666',
    fontWeight: '600',
  },
  submitButton: {
    flex: 1,
    backgroundColor: '#007AFF',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  submitText: {
    color: '#fff',
    fontWeight: '600',
  },
  disabled: {
    opacity: 0.6,
  },
});