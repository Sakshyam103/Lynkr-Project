import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, RefreshControl } from 'react-native';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { apiService } from '../services/api';

export default function HomeScreen({ navigation }: any) {
  const user = useSelector((state: RootState) => state.auth.user);
  const [stats, setStats] = useState({
    eventsAttended: 0,
    pointsEarned: 0,
    contentCreated: 0,
    badges: 0,
  });
  const [recentActivity, setRecentActivity] = useState([]);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const [rewards, badges] = await Promise.all([
        apiService.getUserRewards(),
        apiService.getUserBadges(),
      ]);
      
      setStats({
        eventsAttended: rewards.eventsAttended || 12,
        pointsEarned: rewards.totalPoints || 450,
        contentCreated: rewards.contentCount || 8,
        badges: badges.length || 5,
      });
    } catch (error) {
      console.log('Error loading dashboard data:', error);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadDashboardData();
    setRefreshing(false);
  };

  const StatCard = ({ number, label, color = '#007AFF' }: any) => (
    <View style={[styles.statCard, { borderTopColor: color }]}>
      <Text style={[styles.statNumber, { color }]}>{number}</Text>
      <Text style={styles.statLabel}>{label}</Text>
    </View>
  );

  const ActionCard = ({ icon, title, description, onPress, color = '#007AFF' }: any) => (
    <TouchableOpacity style={styles.actionCard} onPress={onPress}>
      <View style={[styles.actionIcon, { backgroundColor: color + '20' }]}>
        <Text style={styles.actionIconText}>{icon}</Text>
      </View>
      <View style={styles.actionContent}>
        <Text style={styles.actionTitle}>{title}</Text>
        <Text style={styles.actionDescription}>{description}</Text>
      </View>
      <Text style={styles.actionArrow}>›</Text>
    </TouchableOpacity>
  );

  return (
    <ScrollView 
      style={styles.container}
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
    >
      <View style={styles.header}>
        <Text style={styles.greeting}>Welcome back,</Text>
        <Text style={styles.userName}>{user?.name || 'User'}!</Text>
        <Text style={styles.subtitle}>Ready to discover amazing events?</Text>
      </View>

      <View style={styles.statsContainer}>
        <StatCard number={stats.eventsAttended} label="Events Attended" color="#28a745" />
        <StatCard number={stats.pointsEarned} label="Points Earned" color="#ffc107" />
        <StatCard number={stats.contentCreated} label="Content Created" color="#dc3545" />
        <StatCard number={stats.badges} label="Badges Earned" color="#6f42c1" />
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Quick Actions</Text>
        
        <ActionCard
          icon="📅"
          title="Discover Events"
          description="Find sponsored events near you"
          onPress={() => navigation.navigate('Events')}
          color="#007AFF"
        />

        <ActionCard
          icon="📸"
          title="Create Content"
          description="Share photos and videos from events"
          onPress={() => navigation.navigate('Content')}
          color="#28a745"
        />

        <ActionCard
          icon="🎁"
          title="View Rewards"
          description="Check your points and badges"
          onPress={() => navigation.navigate('Profile')}
          color="#ffc107"
        />

        <ActionCard
          icon="📊"
          title="Analytics"
          description="View your engagement statistics"
          onPress={() => {}}
          color="#6f42c1"
        />
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Recent Activity</Text>
        <View style={styles.activityContainer}>
          <View style={styles.activityItem}>
            <View style={styles.activityIcon}>
              <Text>✅</Text>
            </View>
            <View style={styles.activityContent}>
              <Text style={styles.activityText}>Checked in to Tech Conference 2024</Text>
              <Text style={styles.activityTime}>2 hours ago</Text>
            </View>
          </View>
          
          <View style={styles.activityItem}>
            <View style={styles.activityIcon}>
              <Text>📸</Text>
            </View>
            <View style={styles.activityContent}>
              <Text style={styles.activityText}>Posted content at Brand Expo</Text>
              <Text style={styles.activityTime}>1 day ago</Text>
            </View>
          </View>
          
          <View style={styles.activityItem}>
            <View style={styles.activityIcon}>
              <Text>🏆</Text>
            </View>
            <View style={styles.activityContent}>
              <Text style={styles.activityText}>Earned 25 points for quality content</Text>
              <Text style={styles.activityTime}>2 days ago</Text>
            </View>
          </View>
        </View>
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
    padding: 24,
    backgroundColor: '#fff',
    borderBottomLeftRadius: 20,
    borderBottomRightRadius: 20,
  },
  greeting: {
    fontSize: 18,
    color: '#666',
  },
  userName: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 4,
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
  },
  statsContainer: {
    flexDirection: 'row',
    padding: 20,
    gap: 12,
  },
  statCard: {
    flex: 1,
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 16,
    alignItems: 'center',
    borderTopWidth: 3,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  statNumber: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  statLabel: {
    fontSize: 12,
    color: '#666',
    textAlign: 'center',
    fontWeight: '500',
  },
  section: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 22,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 16,
  },
  actionCard: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 16,
    marginBottom: 12,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  actionIcon: {
    width: 48,
    height: 48,
    borderRadius: 24,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 16,
  },
  actionIconText: {
    fontSize: 20,
  },
  actionContent: {
    flex: 1,
  },
  actionTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 2,
  },
  actionDescription: {
    fontSize: 14,
    color: '#666',
  },
  actionArrow: {
    fontSize: 20,
    color: '#ccc',
  },
  activityContainer: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 16,
  },
  activityItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  activityIcon: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: '#f8f9fa',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  activityContent: {
    flex: 1,
  },
  activityText: {
    fontSize: 14,
    color: '#333',
    fontWeight: '500',
  },
  activityTime: {
    fontSize: 12,
    color: '#666',
    marginTop: 2,
  },
});