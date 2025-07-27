import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, TextInput, Alert, Modal } from 'react-native';
import { apiService } from '../services/api';

export default function BrandSurveyScreen() {
  const [surveys, setSurveys] = useState([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedSurvey, setSelectedSurvey] = useState<any>(null);
  const [showAnalyticsModal, setShowAnalyticsModal] = useState(false);
  const [newSurvey, setNewSurvey] = useState({
    eventId: '',
    title: '',
    description: '',
    rewardPoints: '50',
    questions: [{ question: '', type: 'multiple_choice', options: [''] }]
  });

  useEffect(() => {
    loadSurveys();
  }, []);

  const loadSurveys = async () => {
    try {
      // Mock data since API might not exist yet
      setSurveys([
        {
          id: 'survey_1',
          title: 'Event Experience Survey',
          eventName: 'Tech Conference 2024',
          totalResponses: 45,
          completionRate: 78.5,
          status: 'active',
          createdAt: '2024-01-15T10:00:00Z'
        },
        {
          id: 'survey_2', 
          title: 'Product Feedback Survey',
          eventName: 'Brand Expo 2024',
          totalResponses: 23,
          completionRate: 65.2,
          status: 'active',
          createdAt: '2024-01-10T14:30:00Z'
        }
      ]);
    } catch (error) {
      console.log('Error loading surveys:', error);
    }
  };

  const createSurvey = async () => {
    if (!newSurvey.title || !newSurvey.eventId) {
      Alert.alert('Error', 'Please fill in required fields');
      return;
    }

    try {
      await apiService.scheduleSurveys(newSurvey.eventId, {
        title: newSurvey.title,
        description: newSurvey.description,
        questions: newSurvey.questions,
        rewardPoints: parseInt(newSurvey.rewardPoints)
      });
      
      Alert.alert('Success', 'Survey created successfully!');
      setShowCreateModal(false);
      setNewSurvey({
        eventId: '',
        title: '',
        description: '',
        rewardPoints: '50',
        questions: [{ question: '', type: 'multiple_choice', options: [''] }]
      });
      loadSurveys();
    } catch (error) {
      Alert.alert('Error', 'Failed to create survey');
    }
  };

  const viewSurveyAnalytics = async (surveyId: string) => {
    try {
      const analytics = await apiService.getSurveyAnalytics(surveyId);
      setSelectedSurvey(analytics);
      setShowAnalyticsModal(true);
    } catch (error) {
      // Mock analytics data
      setSelectedSurvey({
        id: surveyId,
        title: 'Event Experience Survey',
        totalResponses: 45,
        completionRate: 78.5,
        avgCompletionTime: 180, // seconds
        responses: [
          { userId: 'user_1', userName: 'John Doe', completedAt: '2024-01-15T11:30:00Z' },
          { userId: 'user_2', userName: 'Jane Smith', completedAt: '2024-01-15T12:15:00Z' },
          { userId: 'user_3', userName: 'Mike Johnson', completedAt: '2024-01-15T13:45:00Z' }
        ]
      });
      setShowAnalyticsModal(true);
    }
  };

  const addQuestion = () => {
    setNewSurvey({
      ...newSurvey,
      questions: [...newSurvey.questions, { question: '', type: 'multiple_choice', options: [''] }]
    });
  };

  const updateQuestion = (index: number, field: string, value: any) => {
    const updatedQuestions = [...newSurvey.questions];
    updatedQuestions[index] = { ...updatedQuestions[index], [field]: value };
    setNewSurvey({ ...newSurvey, questions: updatedQuestions });
  };

  const addOption = (questionIndex: number) => {
    const updatedQuestions = [...newSurvey.questions];
    updatedQuestions[questionIndex].options.push('');
    setNewSurvey({ ...newSurvey, questions: updatedQuestions });
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Survey Management</Text>
        <TouchableOpacity 
          style={styles.createButton}
          onPress={() => setShowCreateModal(true)}
        >
          <Text style={styles.createButtonText}>+ Create Survey</Text>
        </TouchableOpacity>
      </View>

      <ScrollView style={styles.surveysList}>
        {surveys.map((survey: any) => (
          <View key={survey.id} style={styles.surveyCard}>
            <View style={styles.surveyHeader}>
              <Text style={styles.surveyTitle}>{survey.title}</Text>
              <View style={[styles.statusBadge, { backgroundColor: survey.status === 'active' ? '#28a745' : '#6c757d' }]}>
                <Text style={styles.statusText}>{survey.status.toUpperCase()}</Text>
              </View>
            </View>
            
            <Text style={styles.eventName}>📍 {survey.eventName}</Text>
            
            <View style={styles.surveyStats}>
              <View style={styles.statItem}>
                <Text style={styles.statNumber}>{survey.totalResponses}</Text>
                <Text style={styles.statLabel}>Responses</Text>
              </View>
              <View style={styles.statItem}>
                <Text style={styles.statNumber}>{survey.completionRate}%</Text>
                <Text style={styles.statLabel}>Completion Rate</Text>
              </View>
              <View style={styles.statItem}>
                <Text style={styles.statNumber}>{formatDate(survey.createdAt)}</Text>
                <Text style={styles.statLabel}>Created</Text>
              </View>
            </View>
            
            <TouchableOpacity 
              style={styles.analyticsButton}
              onPress={() => viewSurveyAnalytics(survey.id)}
            >
              <Text style={styles.analyticsButtonText}>View Analytics & Responses</Text>
            </TouchableOpacity>
          </View>
        ))}
      </ScrollView>

      {/* Create Survey Modal */}
      <Modal visible={showCreateModal} transparent animationType="slide">
        <View style={styles.modalOverlay}>
          <View style={styles.modal}>
            <ScrollView>
              <Text style={styles.modalTitle}>Create New Survey</Text>
              
              <TextInput
                style={styles.input}
                placeholder="Event ID"
                value={newSurvey.eventId}
                onChangeText={(text) => setNewSurvey({...newSurvey, eventId: text})}
              />
              
              <TextInput
                style={styles.input}
                placeholder="Survey Title"
                value={newSurvey.title}
                onChangeText={(text) => setNewSurvey({...newSurvey, title: text})}
              />
              
              <TextInput
                style={styles.input}
                placeholder="Description"
                value={newSurvey.description}
                onChangeText={(text) => setNewSurvey({...newSurvey, description: text})}
                multiline
              />
              
              <TextInput
                style={styles.input}
                placeholder="Reward Points"
                value={newSurvey.rewardPoints}
                onChangeText={(text) => setNewSurvey({...newSurvey, rewardPoints: text})}
                keyboardType="numeric"
              />
              
              <Text style={styles.sectionTitle}>Questions</Text>
              {newSurvey.questions.map((question, index) => (
                <View key={index} style={styles.questionContainer}>
                  <TextInput
                    style={styles.input}
                    placeholder={`Question ${index + 1}`}
                    value={question.question}
                    onChangeText={(text) => updateQuestion(index, 'question', text)}
                  />
                  
                  {question.type === 'multiple_choice' && (
                    <View>
                      {question.options.map((option, optionIndex) => (
                        <TextInput
                          key={optionIndex}
                          style={styles.optionInput}
                          placeholder={`Option ${optionIndex + 1}`}
                          value={option}
                          onChangeText={(text) => {
                            const newOptions = [...question.options];
                            newOptions[optionIndex] = text;
                            updateQuestion(index, 'options', newOptions);
                          }}
                        />
                      ))}
                      <TouchableOpacity 
                        style={styles.addOptionButton}
                        onPress={() => addOption(index)}
                      >
                        <Text style={styles.addOptionText}>+ Add Option</Text>
                      </TouchableOpacity>
                    </View>
                  )}
                </View>
              ))}
              
              <TouchableOpacity style={styles.addQuestionButton} onPress={addQuestion}>
                <Text style={styles.addQuestionText}>+ Add Question</Text>
              </TouchableOpacity>
              
              <View style={styles.modalButtons}>
                <TouchableOpacity 
                  style={styles.cancelButton}
                  onPress={() => setShowCreateModal(false)}
                >
                  <Text style={styles.cancelText}>Cancel</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.submitButton} onPress={createSurvey}>
                  <Text style={styles.submitText}>Create Survey</Text>
                </TouchableOpacity>
              </View>
            </ScrollView>
          </View>
        </View>
      </Modal>

      {/* Analytics Modal */}
      <Modal visible={showAnalyticsModal} transparent animationType="slide">
        <View style={styles.modalOverlay}>
          <View style={styles.modal}>
            <ScrollView>
              <Text style={styles.modalTitle}>Survey Analytics</Text>
              
              {selectedSurvey && (
                <>
                  <View style={styles.analyticsStats}>
                    <View style={styles.analyticCard}>
                      <Text style={styles.analyticNumber}>{selectedSurvey.totalResponses}</Text>
                      <Text style={styles.analyticLabel}>Total Responses</Text>
                    </View>
                    <View style={styles.analyticCard}>
                      <Text style={styles.analyticNumber}>{selectedSurvey.completionRate}%</Text>
                      <Text style={styles.analyticLabel}>Completion Rate</Text>
                    </View>
                    <View style={styles.analyticCard}>
                      <Text style={styles.analyticNumber}>{Math.floor(selectedSurvey.avgCompletionTime / 60)}m</Text>
                      <Text style={styles.analyticLabel}>Avg Time</Text>
                    </View>
                  </View>
                  
                  <Text style={styles.sectionTitle}>Recent Responses</Text>
                  {selectedSurvey.responses?.map((response: any, index: number) => (
                    <View key={index} style={styles.responseCard}>
                      <Text style={styles.responseName}>{response.userName}</Text>
                      <Text style={styles.responseDate}>{formatDate(response.completedAt)}</Text>
                    </View>
                  ))}
                </>
              )}
              
              <TouchableOpacity 
                style={styles.closeButton}
                onPress={() => setShowAnalyticsModal(false)}
              >
                <Text style={styles.closeButtonText}>Close</Text>
              </TouchableOpacity>
            </ScrollView>
          </View>
        </View>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
  },
  createButton: {
    backgroundColor: '#28a745',
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
  },
  createButtonText: {
    color: '#fff',
    fontWeight: '600',
  },
  surveysList: {
    padding: 20,
  },
  surveyCard: {
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 12,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  surveyHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  surveyTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    flex: 1,
  },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 12,
  },
  statusText: {
    color: '#fff',
    fontSize: 10,
    fontWeight: 'bold',
  },
  eventName: {
    fontSize: 14,
    color: '#666',
    marginBottom: 12,
  },
  surveyStats: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 12,
  },
  statItem: {
    alignItems: 'center',
  },
  statNumber: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#28a745',
  },
  statLabel: {
    fontSize: 12,
    color: '#666',
  },
  analyticsButton: {
    backgroundColor: '#007AFF',
    padding: 10,
    borderRadius: 6,
    alignItems: 'center',
  },
  analyticsButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  modalOverlay: {
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
    maxHeight: '80%',
  },
  modalTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 20,
    textAlign: 'center',
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    marginBottom: 12,
    fontSize: 16,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginTop: 16,
    marginBottom: 12,
  },
  questionContainer: {
    backgroundColor: '#f8f9fa',
    padding: 12,
    borderRadius: 8,
    marginBottom: 12,
  },
  optionInput: {
    borderWidth: 1,
    borderColor: '#e0e0e0',
    borderRadius: 6,
    padding: 8,
    marginBottom: 8,
    fontSize: 14,
  },
  addOptionButton: {
    backgroundColor: '#e9ecef',
    padding: 8,
    borderRadius: 6,
    alignItems: 'center',
  },
  addOptionText: {
    color: '#007AFF',
    fontSize: 14,
  },
  addQuestionButton: {
    backgroundColor: '#007AFF',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
    marginBottom: 20,
  },
  addQuestionText: {
    color: '#fff',
    fontWeight: '600',
  },
  modalButtons: {
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
  },
  submitButton: {
    flex: 1,
    backgroundColor: '#28a745',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  submitText: {
    color: '#fff',
    fontWeight: '600',
  },
  analyticsStats: {
    flexDirection: 'row',
    gap: 12,
    marginBottom: 20,
  },
  analyticCard: {
    flex: 1,
    backgroundColor: '#f8f9fa',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  analyticNumber: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#28a745',
  },
  analyticLabel: {
    fontSize: 12,
    color: '#666',
  },
  responseCard: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: 12,
    backgroundColor: '#f8f9fa',
    borderRadius: 8,
    marginBottom: 8,
  },
  responseName: {
    fontSize: 14,
    fontWeight: '600',
  },
  responseDate: {
    fontSize: 12,
    color: '#666',
  },
  closeButton: {
    backgroundColor: '#007AFF',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
    marginTop: 16,
  },
  closeButtonText: {
    color: '#fff',
    fontWeight: '600',
  },
});