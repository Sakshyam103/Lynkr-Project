import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, RefreshControl } from 'react-native';
import { apiService } from '../services/api';

export default function UserRewardsScreen() {
  const [rewards, setRewards] = useState<any>(null);
  const [badges, setBadges] = useState([]);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    loadRewardsData();
  }, []);

  const loadRewardsData = async () => {
    try {
      const [rewardsData, badgesData] = await Promise.all([
        apiService.getUserRewards(),
        apiService.getUserBadges(),
      ]);
      setRewards(rewardsData);
      setBadges(badgesData);
    } catch (error) {
      console.log('Error loading rewards:', error);
      // Mock data
      setRewards({
        totalPoints: 1250,
        level: 'Gold',
        nextLevelPoints: 1500,
        eventsAttended: 15,
        contentCreated: 23,
        surveysCompleted: 8,
      });
      setBadges([
        { id: 1, name: 'Event Explorer', description: 'Attended 10+ events', icon: '🏆', earned: true },
        { id: 2, name: 'Content Creator', description: 'Posted 20+ content items', icon: '📸', earned: true },
        { id: 3, name: 'Survey Master', description: 'Completed 5+ surveys', icon: '📝', earned: true },
        { id: 4, name: 'Early Adopter', description: 'One of first 1000 users', icon: '⭐', earned: true },
        { id: 5, name: 'Social Butterfly', description: 'Attend 25+ events', icon: '🦋', earned: false },
      ]);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadRewardsData();
    setRefreshing(false);
  };

  const getProgressPercentage = () => {
    if (!rewards) return 0;
    return (rewards.totalPoints / rewards.nextLevelPoints) * 100;
  };

  return (
    <ScrollView 
      style={styles.container}
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
    >
      <View style={styles.header}>
        <Text style={styles.title}>My Rewards</Text>
        <Text style={styles.subtitle}>Points, badges, and achievements</Text>
      </View>

      {rewards && (
        <>
          <View style={styles.pointsCard}>
            <View style={styles.pointsHeader}>
              <Text style={styles.pointsTitle}>Total Points</Text>
              <Text style={styles.level}>{rewards.level} Level</Text>
            </View>
            <Text style={styles.pointsNumber}>{rewards.totalPoints}</Text>
            
            <View style={styles.progressContainer}>
              <View style={styles.progressBar}>
                <View 
                  style={[styles.progressFill, { width: `${getProgressPercentage()}%` }]} 
                />
              </View>
              <Text style={styles.progressText}>
                {rewards.nextLevelPoints - rewards.totalPoints} points to next level
              </Text>
            </View>
          </View>

          <View style={styles.statsContainer}>
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{rewards.eventsAttended}</Text>
              <Text style={styles.statLabel}>Events Attended</Text>
            </View>
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{rewards.contentCreated}</Text>
              <Text style={styles.statLabel}>Content Created</Text>
            </View>
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{rewards.surveysCompleted}</Text>
              <Text style={styles.statLabel}>Surveys Completed</Text>
            </View>
          </View>
        </>
      )}

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Badges & Achievements</Text>
        
        <View style={styles.badgesGrid}>
          {badges.map((badge: any) => (
            <View 
              key={badge.id} 
              style={[styles.badgeCard, !badge.earned && styles.badgeCardLocked]}
            >
              <Text style={[styles.badgeIcon, !badge.earned && styles.badgeIconLocked]}>
                {badge.earned ? badge.icon : '🔒'}
              </Text>
              <Text style={[styles.badgeName, !badge.earned && styles.badgeNameLocked]}>
                {badge.name}
              </Text>
              <Text style={[styles.badgeDescription, !badge.earned && styles.badgeDescriptionLocked]}>
                {badge.description}
              </Text>
              {badge.earned && (
                <View style={styles.earnedBadge}>
                  <Text style={styles.earnedText}>✅ Earned</Text>
                </View>
              )}
            </View>
          ))}
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Earn More Points</Text>
        
        <TouchableOpacity style={styles.actionCard}>
          <Text style={styles.actionIcon}>📅</Text>
          <View style={styles.actionContent}>
            <Text style={styles.actionTitle}>Attend Events</Text>
            <Text style={styles.actionDescription}>+50 points per check-in</Text>
          </View>
        </TouchableOpacity>

        <TouchableOpacity style={styles.actionCard}>
          <Text style={styles.actionIcon}>📸</Text>
          <View style={styles.actionContent}>
            <Text style={styles.actionTitle}>Create Content</Text>
            <Text style={styles.actionDescription}>+25 points per post</Text>
          </View>
        </TouchableOpacity>

        <TouchableOpacity style={styles.actionCard}>
          <Text style={styles.actionIcon}>📝</Text>
          <View style={styles.actionContent}>
            <Text style={styles.actionTitle}>Complete Surveys</Text>
            <Text style={styles.actionDescription}>+10-100 points per survey</Text>
          </View>
        </TouchableOpacity>

        <TouchableOpacity style={styles.actionCard}>
          <Text style={styles.actionIcon}>🛒</Text>
          <View style={styles.actionContent}>
            <Text style={styles.actionTitle}>Make Purchases</Text>
            <Text style={styles.actionDescription}>+5 points per $1 spent</Text>
          </View>
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
  pointsCard: {
    backgroundColor: '#fff',
    margin: 20,
    padding: 24,
    borderRadius: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  pointsHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  pointsTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#333',
  },
  level: {
    backgroundColor: '#ffc107',
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 12,
    fontSize: 12,
    fontWeight: 'bold',
    color: '#fff',
  },
  pointsNumber: {
    fontSize: 48,
    fontWeight: 'bold',
    color: '#007AFF',
    textAlign: 'center',
    marginBottom: 20,
  },
  progressContainer: {
    alignItems: 'center',
  },
  progressBar: {
    width: '100%',
    height: 8,
    backgroundColor: '#e0e0e0',
    borderRadius: 4,
    marginBottom: 8,
  },
  progressFill: {
    height: '100%',
    backgroundColor: '#007AFF',
    borderRadius: 4,
  },
  progressText: {
    fontSize: 14,
    color: '#666',
  },
  statsContainer: {
    flexDirection: 'row',
    paddingHorizontal: 20,
    gap: 12,
  },
  statItem: {
    flex: 1,
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
  },
  statNumber: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#28a745',
  },
  statLabel: {
    fontSize: 12,
    color: '#666',
    textAlign: 'center',
    marginTop: 4,
  },
  section: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  badgesGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },
  badgeCard: {
    backgroundColor: '#fff',
    width: '48%',
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  badgeCardLocked: {
    backgroundColor: '#f8f9fa',
    opacity: 0.6,
  },
  badgeIcon: {
    fontSize: 32,
    marginBottom: 8,
  },
  badgeIconLocked: {
    opacity: 0.5,
  },
  badgeName: {
    fontSize: 14,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 4,
  },
  badgeNameLocked: {
    color: '#999',
  },
  badgeDescription: {
    fontSize: 12,
    color: '#666',
    textAlign: 'center',
    marginBottom: 8,
  },
  badgeDescriptionLocked: {
    color: '#999',
  },
  earnedBadge: {
    backgroundColor: '#28a745',
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 8,
  },
  earnedText: {
    fontSize: 10,
    color: '#fff',
    fontWeight: 'bold',
  },
  actionCard: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 12,
    marginBottom: 12,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  actionIcon: {
    fontSize: 24,
    marginRight: 16,
  },
  actionContent: {
    flex: 1,
  },
  actionTitle: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 2,
  },
  actionDescription: {
    fontSize: 14,
    color: '#28a745',
    fontWeight: '500',
  },
});