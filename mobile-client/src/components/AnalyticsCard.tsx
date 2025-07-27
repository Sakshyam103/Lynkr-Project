import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ActivityIndicator } from 'react-native';
import { apiService } from '../services/api';

interface AnalyticsCardProps {
  eventId: string;
  type: 'engagement' | 'attendance' | 'content' | 'realtime';
  title: string;
}

export default function AnalyticsCard({ eventId, type, title }: AnalyticsCardProps) {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const loadAnalytics = async () => {
    setLoading(true);
    try {
      let result;
      switch (type) {
        case 'engagement':
          result = await apiService.getEngagementMetrics(eventId);
          break;
        case 'attendance':
          result = await apiService.getAttendanceAnalytics(eventId);
          break;
        case 'realtime':
          result = await apiService.getRealtimeStats(eventId);
          break;
        default:
          result = {};
      }
      setData(result);
    } catch (error) {
      console.log('Analytics error:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAnalytics();
  }, [eventId, type]);

  return (
    <View style={styles.card}>
      <Text style={styles.title}>{title}</Text>
      {loading ? (
        <ActivityIndicator />
      ) : (
        <View style={styles.content}>
          {type === 'engagement' && data && (
            <>
              <Text style={styles.metric}>Rate: {Math.round((data.engagementRate || 0) * 100)}%</Text>
              <Text style={styles.metric}>Interactions: {data.totalInteractions || 0}</Text>
            </>
          )}
          {type === 'attendance' && data && (
            <>
              <Text style={styles.metric}>Total: {data.totalAttendees || 0}</Text>
              <Text style={styles.metric}>Check-ins: {data.checkins || 0}</Text>
            </>
          )}
          {type === 'realtime' && data && (
            <>
              <Text style={styles.metric}>Live: {data.currentAttendees || 0}</Text>
              <Text style={styles.metric}>Active: {data.activeUsers || 0}</Text>
            </>
          )}
        </View>
      )}
      <TouchableOpacity style={styles.refreshButton} onPress={loadAnalytics}>
        <Text style={styles.refreshText}>Refresh</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  title: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  content: {
    marginBottom: 12,
  },
  metric: {
    fontSize: 14,
    color: '#666',
    marginBottom: 4,
  },
  refreshButton: {
    backgroundColor: '#007AFF',
    padding: 8,
    borderRadius: 6,
    alignItems: 'center',
  },
  refreshText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
  },
});