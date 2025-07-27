import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, Alert, Modal } from 'react-native';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../store/store';
import { logout } from '../store/slices/authSlice';
import { apiService } from '../services/api';

export default function ProfileScreen() {
  const dispatch = useDispatch();
  const { user } = useSelector((state: RootState) => state.auth);
  const [showDataModal, setShowDataModal] = useState(false);

  const isUser = user?.role === 'user' || !user?.role;
  const isBrand = user?.role === 'brand';

  const handleLogout = () => {
    Alert.alert(
      'Logout',
      'Are you sure you want to logout?',
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Logout', 
          style: 'destructive',
          onPress: () => {
            apiService.logout();
            dispatch(logout());
          }
        }
      ]
    );
  };

  const handleDataDeletion = async () => {
    Alert.alert(
      'Delete Account',
      'This will permanently delete all your data. This action cannot be undone.',
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Delete', 
          style: 'destructive',
          onPress: async () => {
            try {
              await apiService.requestDataDeletion();
              Alert.alert('Success', 'Data deletion request submitted. Your account will be deleted within 30 days.');
            } catch (error) {
              Alert.alert('Error', 'Failed to submit deletion request');
            }
          }
        }
      ]
    );
  };

  const handleDataExport = async () => {
    try {
      const data = await apiService.exportUserData();
      Alert.alert('Success', 'Your data export has been prepared. Check your email for download link.');
    } catch (error) {
      Alert.alert('Error', 'Failed to export data');
    }
  };

  const handleDataAnonymization = async () => {
    Alert.alert(
      'Anonymize Data',
      'This will remove all personally identifiable information from your account.',
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Anonymize', 
          onPress: async () => {
            try {
              await apiService.anonymizeUserData();
              Alert.alert('Success', 'Your data has been anonymized');
            } catch (error) {
              Alert.alert('Error', 'Failed to anonymize data');
            }
          }
        }
      ]
    );
  };

  const getRoleColor = () => {
    switch (user?.role) {
      case 'brand': return '#28a745';
      case 'organization': return '#ffc107';
      default: return '#007AFF';
    }
  };

  const getRoleIcon = () => {
    switch (user?.role) {
      case 'brand': return '🏢';
      case 'organization': return '🎪';
      default: return '👤';
    }
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <View style={[styles.avatar, { backgroundColor: getRoleColor() }]}>
          <Text style={styles.avatarText}>{getRoleIcon()}</Text>
        </View>
        <Text style={styles.userName}>{user?.name || 'User'}</Text>
        <Text style={styles.userEmail}>{user?.email || 'user@example.com'}</Text>
        <View style={[styles.roleBadge, { backgroundColor: getRoleColor() }]}>
          <Text style={styles.roleText}>{(user?.role || 'user').toUpperCase()}</Text>
        </View>
      </View>

      {/* USER ONLY STATS */}
      {isUser && (
        <View style={styles.statsSection}>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>12</Text>
            <Text style={styles.statLabel}>Events</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>450</Text>
            <Text style={styles.statLabel}>Points</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>8</Text>
            <Text style={styles.statLabel}>Content</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>5</Text>
            <Text style={styles.statLabel}>Badges</Text>
          </View>
        </View>
      )}

      {/* BRAND ONLY STATS */}
      {isBrand && (
        <View style={styles.statsSection}>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>8</Text>
            <Text style={styles.statLabel}>Campaigns</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>25K</Text>
            <Text style={styles.statLabel}>Reach</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>$15K</Text>
            <Text style={styles.statLabel}>ROI</Text>
          </View>
          <View style={styles.statItem}>
            <Text style={styles.statNumber}>92%</Text>
            <Text style={styles.statLabel}>Engagement</Text>
          </View>
        </View>
      )}

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Account Settings</Text>
        
        <TouchableOpacity style={styles.settingItem}>
          <Text style={styles.settingText}>🔔 Notifications</Text>
          <Text style={styles.settingArrow}>›</Text>
        </TouchableOpacity>

        <TouchableOpacity 
          style={styles.settingItem}
          onPress={() => {
            // Navigate to privacy settings or call API
            apiService.updatePrivacySettings({
              analyticsConsent: true,
              marketingConsent: false,
              locationTracking: true
            });
          }}
        >
          <Text style={styles.settingText}>🔒 Privacy Settings</Text>
          <Text style={styles.settingArrow}>›</Text>
        </TouchableOpacity>

        <TouchableOpacity style={styles.settingItem}>
          <Text style={styles.settingText}>📍 Location Preferences</Text>
          <Text style={styles.settingArrow}>›</Text>
        </TouchableOpacity>

        {isBrand && (
          <TouchableOpacity style={styles.settingItem}>
            <Text style={styles.settingText}>💳 Billing & Payments</Text>
            <Text style={styles.settingArrow}>›</Text>
          </TouchableOpacity>
        )}

        <TouchableOpacity style={styles.settingItem}>
          <Text style={styles.settingText}>❓ Help & Support</Text>
          <Text style={styles.settingArrow}>›</Text>
        </TouchableOpacity>

        <TouchableOpacity style={styles.settingItem}>
          <Text style={styles.settingText}>📄 Terms & Privacy</Text>
          <Text style={styles.settingArrow}>›</Text>
        </TouchableOpacity>
      </View>

      {/* DATA MANAGEMENT SECTION - USER ONLY */}
      {isUser && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Data Management</Text>
          
          <TouchableOpacity 
            style={styles.settingItem}
            onPress={() => setShowDataModal(true)}
          >
            <Text style={styles.settingText}>📊 My Data</Text>
            <Text style={styles.settingArrow}>›</Text>
          </TouchableOpacity>
        </View>
      )}

      <View style={styles.section}>
        <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
          <Text style={styles.logoutText}>Logout</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.footer}>
        <Text style={styles.footerText}>Lynkr v1.0.0</Text>
        <Text style={styles.footerText}>© 2024 Brand Activations Platform</Text>
      </View>

      {/* DATA MANAGEMENT MODAL */}
      <Modal visible={showDataModal} transparent animationType="slide">
        <View style={styles.modalOverlay}>
          <View style={styles.modal}>
            <Text style={styles.modalTitle}>Data Management</Text>
            <Text style={styles.modalSubtitle}>Manage your personal data and privacy</Text>
            
            <TouchableOpacity style={styles.dataButton} onPress={handleDataExport}>
              <Text style={styles.dataButtonText}>📥 Export My Data</Text>
              <Text style={styles.dataButtonDesc}>Download all your data</Text>
            </TouchableOpacity>
            
            <TouchableOpacity style={styles.dataButton} onPress={handleDataAnonymization}>
              <Text style={styles.dataButtonText}>🔒 Anonymize Data</Text>
              <Text style={styles.dataButtonDesc}>Remove personal identifiers</Text>
            </TouchableOpacity>
            
            <TouchableOpacity style={[styles.dataButton, styles.dangerButton]} onPress={handleDataDeletion}>
              <Text style={[styles.dataButtonText, styles.dangerText]}>🗑️ Delete Account</Text>
              <Text style={[styles.dataButtonDesc, styles.dangerText]}>Permanently delete all data</Text>
            </TouchableOpacity>
            
            <TouchableOpacity 
              style={styles.closeButton}
              onPress={() => setShowDataModal(false)}
            >
              <Text style={styles.closeButtonText}>Close</Text>
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
    padding: 30,
    alignItems: 'center',
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 15,
  },
  avatarText: {
    fontSize: 32,
    color: '#fff',
  },
  userName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 5,
  },
  userEmail: {
    fontSize: 16,
    color: '#666',
    marginBottom: 10,
  },
  roleBadge: {
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 12,
  },
  roleText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: 'bold',
  },
  statsSection: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    marginTop: 10,
    paddingVertical: 20,
  },
  statItem: {
    flex: 1,
    alignItems: 'center',
  },
  statNumber: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#007AFF',
  },
  statLabel: {
    fontSize: 12,
    color: '#666',
    marginTop: 4,
  },
  section: {
    backgroundColor: '#fff',
    marginTop: 10,
    padding: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#333',
    marginBottom: 15,
  },
  settingItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 15,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  settingText: {
    fontSize: 16,
    color: '#333',
  },
  settingArrow: {
    fontSize: 18,
    color: '#ccc',
  },
  logoutButton: {
    backgroundColor: '#dc3545',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
  },
  logoutText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  footer: {
    padding: 20,
    alignItems: 'center',
  },
  footerText: {
    fontSize: 12,
    color: '#666',
    marginBottom: 2,
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
    marginBottom: 8,
    textAlign: 'center',
  },
  modalSubtitle: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
    marginBottom: 20,
  },
  dataButton: {
    padding: 16,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#e0e0e0',
    marginBottom: 12,
  },
  dataButtonText: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 4,
  },
  dataButtonDesc: {
    fontSize: 14,
    color: '#666',
  },
  dangerButton: {
    borderColor: '#dc3545',
    backgroundColor: '#fff5f5',
  },
  dangerText: {
    color: '#dc3545',
  },
  closeButton: {
    backgroundColor: '#007AFF',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
    marginTop: 8,
  },
  closeButtonText: {
    color: '#fff',
    fontWeight: '600',
  },
});