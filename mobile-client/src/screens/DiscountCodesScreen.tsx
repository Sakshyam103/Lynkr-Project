import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, TextInput, Alert, Modal } from 'react-native';
import { apiService } from '../services/api';

export default function DiscountCodesScreen() {
  const [codes, setCodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newCode, setNewCode] = useState({
    eventId: '',
    discountPct: '',
    maxUses: '',
    expiresIn: '30',
  });

  useEffect(() => {
    loadCodes();
  }, []);

  const loadCodes = async () => {
    setLoading(true);
    try {
      const response = await apiService.getBrandCodes();
      // Handle both response.codes and direct array response
      setCodes(response.codes || response || []);
    } catch (error) {
      console.log('Error loading codes:', error);
      Alert.alert('Error', 'Failed to load discount codes');
    } finally {
      setLoading(false);
    }
  };

  const createCode = async () => {
    if (!newCode.eventId || !newCode.discountPct || !newCode.maxUses) {
      Alert.alert('Error', 'Please fill in all required fields');
      return;
    }

    try {
      await apiService.generateDiscountCode({
        eventId: newCode.eventId,
        discountPct: parseFloat(newCode.discountPct),
        maxUses: parseInt(newCode.maxUses),
        expiresIn: parseInt(newCode.expiresIn),
      });
      Alert.alert('Success', 'Discount code created successfully!');
      setShowCreateModal(false);
      setNewCode({ eventId: '', discountPct: '', maxUses: '', expiresIn: '30' });
      loadCodes();
    } catch (error) {
      Alert.alert('Error', 'Failed to create discount code');
    }
  };

  const validateCode = async (code: string) => {
    try {
      const result = await apiService.validateDiscountCode(code);
      Alert.alert('Validation Result', `Code is ${result.valid ? 'valid' : 'invalid'}`);
    } catch (error) {
      Alert.alert('Error', 'Failed to validate code');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const getUsagePercentage = (used: number, max: number) => {
    return max > 0 ? (used / max) * 100 : 0;
  };

  const getStatusColor = (code: any) => {
    const now = new Date();
    const expires = new Date(code.expiresAt);
    
    if (expires < now) return '#dc3545'; // Expired - red
    if (code.usedCount >= code.maxUses) return '#6c757d'; // Used up - gray
    return '#28a745'; // Active - green
  };

  const getStatusText = (code: any) => {
    const now = new Date();
    const expires = new Date(code.expiresAt);
    
    if (expires < now) return 'Expired';
    if (code.usedCount >= code.maxUses) return 'Used Up';
    return 'Active';
  };

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <Text style={styles.loadingText}>Loading discount codes...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Discount Codes</Text>
        <TouchableOpacity 
          style={styles.createButton}
          onPress={() => setShowCreateModal(true)}
        >
          <Text style={styles.createButtonText}>+ Create Code</Text>
        </TouchableOpacity>
      </View>

      <ScrollView style={styles.codesList}>
        {codes.length === 0 ? (
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No discount codes yet</Text>
            <Text style={styles.emptySubtext}>Create your first discount code to get started</Text>
          </View>
        ) : (
          codes.map((code: any) => (
            <View key={code.id} style={styles.codeCard}>
              <View style={styles.codeHeader}>
                <Text style={styles.codeText}>{code.code}</Text>
                <View style={[styles.statusBadge, { backgroundColor: getStatusColor(code) }]}>
                  <Text style={styles.statusText}>{getStatusText(code)}</Text>
                </View>
              </View>
              
              <View style={styles.codeValue}>
                <Text style={styles.discountText}>{code.discountPct}% OFF</Text>
              </View>
              
              <View style={styles.usageContainer}>
                <View style={styles.usageBar}>
                  <View 
                    style={[
                      styles.usageProgress, 
                      { width: `${getUsagePercentage(code.usedCount, code.maxUses)}%` }
                    ]} 
                  />
                </View>
                <Text style={styles.usageText}>
                  {code.usedCount}/{code.maxUses} uses ({getUsagePercentage(code.usedCount, code.maxUses).toFixed(0)}%)
                </Text>
              </View>
              
              <Text style={styles.expiryText}>
                Expires: {formatDate(code.expiresAt)}
              </Text>
              
              <View style={styles.codeActions}>
                <TouchableOpacity 
                  style={styles.validateButton}
                  onPress={() => validateCode(code.code)}
                >
                  <Text style={styles.validateText}>Validate</Text>
                </TouchableOpacity>
                
                <TouchableOpacity 
                  style={styles.copyButton}
                  onPress={() => {
                    // In a real app, you'd copy to clipboard
                    Alert.alert('Copied', `Code ${code.code} copied to clipboard`);
                  }}
                >
                  <Text style={styles.copyText}>Copy Code</Text>
                </TouchableOpacity>
              </View>
            </View>
          ))
        )}
      </ScrollView>

      {/* Create Code Modal */}
      <Modal visible={showCreateModal} transparent animationType="slide">
        <View style={styles.modalOverlay}>
          <View style={styles.modal}>
            <Text style={styles.modalTitle}>Create Discount Code</Text>
            
            <TextInput
              style={styles.input}
              placeholder="Event ID"
              value={newCode.eventId}
              onChangeText={(text) => setNewCode({...newCode, eventId: text})}
            />
            
            <TextInput
              style={styles.input}
              placeholder="Discount Percentage (e.g., 20)"
              value={newCode.discountPct}
              onChangeText={(text) => setNewCode({...newCode, discountPct: text})}
              keyboardType="numeric"
            />
            
            <TextInput
              style={styles.input}
              placeholder="Max Uses"
              value={newCode.maxUses}
              onChangeText={(text) => setNewCode({...newCode, maxUses: text})}
              keyboardType="numeric"
            />
            
            <TextInput
              style={styles.input}
              placeholder="Expires in (days)"
              value={newCode.expiresIn}
              onChangeText={(text) => setNewCode({...newCode, expiresIn: text})}
              keyboardType="numeric"
            />
            
            <View style={styles.modalButtons}>
              <TouchableOpacity 
                style={styles.cancelButton}
                onPress={() => setShowCreateModal(false)}
              >
                <Text style={styles.cancelText}>Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity style={styles.submitButton} onPress={createCode}>
                <Text style={styles.submitText}>Create</Text>
              </TouchableOpacity>
            </View>
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
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    fontSize: 16,
    color: '#666',
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
  codesList: {
    padding: 20,
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
    textAlign: 'center',
  },
  codeCard: {
    backgroundColor: '#fff',
    padding: 20,
    borderRadius: 12,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  codeHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  codeText: {
    fontSize: 20,
    fontWeight: 'bold',
    fontFamily: 'monospace',
    color: '#333',
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
  codeValue: {
    alignItems: 'center',
    marginBottom: 16,
  },
  discountText: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#28a745',
  },
  usageContainer: {
    marginBottom: 12,
  },
  usageBar: {
    height: 8,
    backgroundColor: '#e9ecef',
    borderRadius: 4,
    marginBottom: 8,
  },
  usageProgress: {
    height: '100%',
    backgroundColor: '#28a745',
    borderRadius: 4,
  },
  usageText: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
  },
  expiryText: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
    marginBottom: 16,
  },
  codeActions: {
    flexDirection: 'row',
    gap: 12,
  },
  validateButton: {
    flex: 1,
    backgroundColor: '#007AFF',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  validateText: {
    color: '#fff',
    fontWeight: '600',
  },
  copyButton: {
    flex: 1,
    backgroundColor: '#f8f9fa',
    borderWidth: 1,
    borderColor: '#28a745',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  copyText: {
    color: '#28a745',
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
    maxWidth: 400,
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
    marginBottom: 16,
    fontSize: 16,
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
});