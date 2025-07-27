import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, TextInput, Alert } from 'react-native';
import { apiService } from '../services/api';

export default function SurveyScreen({ navigation }: any) {
  const [availableSurveys, setAvailableSurveys] = useState<any[]>([]);
  const [completedSurveys, setCompletedSurveys] = useState(new Set());
  const [currentSurvey, setCurrentSurvey] = useState<any>(null);
  const [responses, setResponses] = useState<any>({});

  useEffect(() => {
    loadAvailableSurveys();
  }, []);

  const loadAvailableSurveys = async () => {
    try {
      const surveys = await apiService.getAvailableSurveys();
      setAvailableSurveys(surveys);
    } catch (error) {
      console.log('Error loading surveys:', error);
      // Mock data fallback
      setAvailableSurveys([
        {
          id: 'survey_1',
          title: 'Event Experience Survey',
          description: 'Tell us about your experience at the event',
          eventName: 'Tech Conference 2024',
          reward: 50,
          questions: [
            {
              id: 'q1',
              type: 'multiple_choice',
              question: 'How would you rate the overall event?',
              options: ['Excellent', 'Good', 'Average', 'Poor']
            },
            {
              id: 'q2',
              type: 'text',
              question: 'What did you like most about the event?'
            }
          ]
        }
      ]);
    }
  };

  const startSurvey = (survey: any) => {
    setCurrentSurvey(survey);
    setResponses({});
  };

  const updateResponse = (questionId: string, answer: string) => {
    setResponses({
      ...responses,
      [questionId]: answer
    });
  };


  const submitSurvey = async () => {
    try {
      // Validate all questions are answered
      const unansweredQuestions = currentSurvey.questions.filter(
        (q: any) => !responses[q.id] || responses[q.id].trim() === ''
      );
  
      if (unansweredQuestions.length > 0) {
        Alert.alert('Error', 'Please answer all questions');
        return;
      }
  
      // Format responses as object (not array) to match backend expectation
      const requestData = {
        surveyId: currentSurvey.id,
        responses: responses  // Send responses object directly
      };
  
      console.log('Sending survey data:', JSON.stringify(requestData, null, 2));
  
      await apiService.submitSurveyResponse(requestData);
      
      // Mark survey as completed
      setCompletedSurveys(prev => new Set([...prev, currentSurvey.id]));
      
      // Remove completed survey from available surveys
      setAvailableSurveys(prev => prev.filter((survey: any) => survey.id !== currentSurvey.id));
      
      Alert.alert('Success', `Survey submitted! You earned ${currentSurvey.reward} points.`, [
        { text: 'OK', onPress: () => {
          setCurrentSurvey(null);
          setResponses({});
        }}
      ]);
    } catch (error) {
      console.log('Survey error:', error);
      Alert.alert('Error', 'Failed to submit survey');
    }
  };
  
  

  if (currentSurvey) {
    return (
      <ScrollView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity 
            style={styles.backButton}
            onPress={() => setCurrentSurvey(null)}
          >
            <Text style={styles.backText}>← Back</Text>
          </TouchableOpacity>
          <Text style={styles.title}>{currentSurvey.title}</Text>
          <Text style={styles.subtitle}>{currentSurvey.eventName}</Text>
        </View>

        <View style={styles.questionsContainer}>
          {currentSurvey.questions.map((question: any, index: number) => (
            <View key={question.id} style={styles.questionCard}>
              <Text style={styles.questionNumber}>Question {index + 1}</Text>
              <Text style={styles.questionText}>{question.question}</Text>

              {question.type === 'multiple_choice' && (
                <View style={styles.optionsContainer}>
                  {question.options.map((option: string) => (
                    <TouchableOpacity
                      key={option}
                      style={[
                        styles.optionButton,
                        responses[question.id] === option && styles.optionButtonSelected
                      ]}
                      onPress={() => updateResponse(question.id, option)}
                    >
                      <Text style={[
                        styles.optionText,
                        responses[question.id] === option && styles.optionTextSelected
                      ]}>
                        {option}
                      </Text>
                    </TouchableOpacity>
                  ))}
                </View>
              )}

              {question.type === 'text' && (
                <TextInput
                  style={styles.textInput}
                  placeholder="Type your answer here..."
                  value={responses[question.id] || ''}
                  onChangeText={(text) => updateResponse(question.id, text)}
                  multiline
                  numberOfLines={3}
                />
              )}
            </View>
          ))}
        </View>

        <View style={styles.submitContainer}>
          <TouchableOpacity style={styles.submitButton} onPress={submitSurvey}>
            <Text style={styles.submitText}>Submit Survey (+{currentSurvey.reward} points)</Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    );
  }

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Available Surveys</Text>
        <Text style={styles.subtitle}>Complete surveys to earn rewards</Text>
      </View>

      <View style={styles.surveysContainer}>
        {availableSurveys.map((survey: any) => (
          <View key={survey.id} style={styles.surveyCard}>
            <Text style={styles.surveyTitle}>{survey.title}</Text>
            <Text style={styles.surveyDescription}>{survey.description}</Text>
            <Text style={styles.surveyEvent}>📍 {survey.eventName}</Text>
            
            <View style={styles.surveyFooter}>
              <Text style={styles.surveyReward}>🎁 {survey.reward} points</Text>
              {completedSurveys.has(survey.id) ? (
                <View style={styles.completedButton}>
                  <Text style={styles.completedText}>✓ Completed</Text>
                </View>
              ) : (
                <TouchableOpacity 
                  style={styles.startButton}
                  onPress={() => startSurvey(survey)}
                >
                  <Text style={styles.startText}>Start Survey</Text>
                </TouchableOpacity>
              )}
            </View>
          </View>
        ))}

        {availableSurveys.length === 0 && (
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No surveys available</Text>
            <Text style={styles.emptySubtext}>Check back later for new surveys</Text>
          </View>
        )}
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
  backButton: {
    marginBottom: 10,
  },
  backText: {
    color: '#007AFF',
    fontSize: 16,
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
  surveysContainer: {
    padding: 20,
  },
  surveyCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 20,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  surveyTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  surveyDescription: {
    fontSize: 14,
    color: '#666',
    marginBottom: 8,
    lineHeight: 20,
  },
  surveyEvent: {
    fontSize: 14,
    color: '#007AFF',
    marginBottom: 16,
  },
  surveyFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  surveyReward: {
    fontSize: 14,
    color: '#28a745',
    fontWeight: '600',
  },
  startButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
  },
  startText: {
    color: '#fff',
    fontWeight: '600',
  },
  completedButton: {
    backgroundColor: '#28a745',
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
  },
  completedText: {
    color: '#fff',
    fontWeight: '600',
  },
  questionsContainer: {
    padding: 20,
  },
  questionCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 20,
    marginBottom: 16,
  },
  questionNumber: {
    fontSize: 12,
    color: '#007AFF',
    fontWeight: '600',
    marginBottom: 8,
  },
  questionText: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 16,
    lineHeight: 22,
  },
  optionsContainer: {
    gap: 8,
  },
  optionButton: {
    padding: 12,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#ddd',
    backgroundColor: '#f8f9fa',
  },
  optionButtonSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  optionText: {
    fontSize: 14,
    color: '#333',
  },
  optionTextSelected: {
    color: '#fff',
  },
  textInput: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    textAlignVertical: 'top',
  },
  submitContainer: {
    padding: 20,
  },
  submitButton: {
    backgroundColor: '#28a745',
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
  },
  submitText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  emptyContainer: {
    alignItems: 'center',
    paddingTop: 60,
  },
  emptyText: {
    fontSize: 18,
    fontWeight: '600',
    color: '#666',
    marginBottom: 8,
  },
  emptySubtext: {
    fontSize: 14,
    color: '#999',
  },
});