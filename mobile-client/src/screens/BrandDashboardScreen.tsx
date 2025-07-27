import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, RefreshControl, Dimensions } from 'react-native';
import { apiService } from '../services/api';

const { width } = Dimensions.get('window');

export default function BrandDashboardScreen({ navigation }: any) {
  const [dashboardData, setDashboardData] = useState<any>(null);
  const [campaigns, setCampaigns] = useState([]);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const [dashboard, campaignData] = await Promise.all([
        apiService.getBrandDashboard(),
        apiService.getBrandCampaigns(),
      ]);
      setDashboardData(dashboard);
      setCampaigns(campaignData.campaigns || []);
    } catch (error) {
      console.log('Dashboard error:', error);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadDashboardData();
    setRefreshing(false);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return '#28a745';
      case 'completed': return '#007AFF';
      case 'paused': return '#ffc107';
      default: return '#6c757d';
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    });
  };

  if (!dashboardData) {
    return (
      <View style={styles.loadingContainer}>
        <Text style={styles.loadingText}>Loading dashboard...</Text>
      </View>
    );
  }

  return (
    <ScrollView 
      style={styles.container}
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
    >
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Brand Dashboard</Text>
        <Text style={styles.subtitle}>Monitor your sponsorship performance</Text>
      </View>

      {/* Key Metrics */}
      <View style={styles.metricsContainer}>
        <View style={styles.metricCard}>
          <Text style={styles.metricNumber}>{dashboardData.totalAttendees?.toLocaleString()}</Text>
          <Text style={styles.metricLabel}>Total Attendees</Text>
          <Text style={styles.metricChange}>+12% vs last month</Text>
        </View>
        
        <View style={styles.metricCard}>
          <Text style={styles.metricNumber}>{dashboardData.engagementRate}%</Text>
          <Text style={styles.metricLabel}>Engagement Rate</Text>
          <Text style={styles.metricChange}>+3.2% vs last month</Text>
        </View>
        
        <View style={styles.metricCard}>
          <Text style={styles.metricNumber}>{dashboardData.conversionRate}%</Text>
          <Text style={styles.metricLabel}>Conversion Rate</Text>
          <Text style={styles.metricChange}>+0.8% vs last month</Text>
        </View>
        
        <View style={styles.metricCard}>
          <Text style={styles.metricNumber}>{dashboardData.contentPieces}</Text>
          <Text style={styles.metricLabel}>Content Pieces</Text>
          <Text style={styles.metricChange}>+15 this week</Text>
        </View>
      </View>

      {/* Attendance Chart */}
      <View style={styles.chartSection}>
        <Text style={styles.sectionTitle}>Monthly Attendance Trend</Text>
        <View style={styles.chartContainer}>
          {dashboardData.attendanceData?.map((item: any, index: number) => {
            const maxAttendees = Math.max(...dashboardData.attendanceData.map((d: any) => d.attendees));
            const height = (item.attendees / maxAttendees) * 120;
            
            return (
              <View key={index} style={styles.chartBar}>
                <View style={[styles.bar, { height }]} />
                <Text style={styles.barValue}>{item.attendees}</Text>
                <Text style={styles.barLabel}>{item.date.split('-')[1]}</Text>
              </View>
            );
          })}
        </View>
      </View>

      {/* Quick Actions */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Quick Actions</Text>
        
        <View style={styles.actionsGrid}>
          <TouchableOpacity 
            style={styles.actionCard}
            onPress={() => navigation.navigate('CreateEvent')}
          >
            <Text style={styles.actionIcon}>📅</Text>
            <Text style={styles.actionTitle}>Create Event</Text>
          </TouchableOpacity>

          <TouchableOpacity 
            style={styles.actionCard}
            onPress={() => navigation.navigate('Analytics')}
          >
            <Text style={styles.actionIcon}>📊</Text>
            <Text style={styles.actionTitle}>View Analytics</Text>
          </TouchableOpacity>

          <TouchableOpacity 
            style={styles.actionCard}
            onPress={() => navigation.navigate('Discounts')}
          >
            <Text style={styles.actionIcon}>🎫</Text>
            <Text style={styles.actionTitle}>Discount Codes</Text>
          </TouchableOpacity>

          <TouchableOpacity 
            style={styles.actionCard}
            onPress={() => navigation.navigate('Content')}
          >
            <Text style={styles.actionIcon}>📸</Text>
            <Text style={styles.actionTitle}>User Content</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Recent Campaigns */}
      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Recent Campaigns</Text>
          <TouchableOpacity onPress={() => {}}>
            <Text style={styles.viewAllText}>View All</Text>
          </TouchableOpacity>
        </View>
        
        {campaigns.slice(0, 3).map((campaign: any) => (
          <View key={campaign.id} style={styles.campaignCard}>
            <View style={styles.campaignHeader}>
              <Text style={styles.campaignName}>{campaign.name}</Text>
              <View style={[styles.statusBadge, { backgroundColor: getStatusColor(campaign.status) }]}>
                <Text style={styles.statusText}>{campaign.status.toUpperCase()}</Text>
              </View>
            </View>
            
            <View style={styles.campaignStats}>
              <View style={styles.campaignStat}>
                <Text style={styles.campaignStatNumber}>{campaign.attendees?.toLocaleString()}</Text>
                <Text style={styles.campaignStatLabel}>Attendees</Text>
              </View>
              <View style={styles.campaignStat}>
                <Text style={styles.campaignStatNumber}>{formatCurrency(campaign.budget)}</Text>
                <Text style={styles.campaignStatLabel}>Budget</Text>
              </View>
              <View style={styles.campaignStat}>
                <Text style={styles.campaignStatNumber}>
                  {formatDate(campaign.startDate)} - {formatDate(campaign.endDate)}
                </Text>
                <Text style={styles.campaignStatLabel}>Duration</Text>
              </View>
            </View>
            
            <View style={styles.campaignActions}>
              <TouchableOpacity style={styles.campaignButton}>
                <Text style={styles.campaignButtonText}>View Details</Text>
              </TouchableOpacity>
              <TouchableOpacity style={[styles.campaignButton, styles.primaryButton]}>
                <Text style={[styles.campaignButtonText, styles.primaryButtonText]}>Analytics</Text>
              </TouchableOpacity>
            </View>
          </View>
        ))}
      </View>

      {/* Performance Insights */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Performance Insights</Text>
        
        <View style={styles.insightCard}>
          <Text style={styles.insightIcon}>🎯</Text>
          <View style={styles.insightContent}>
            <Text style={styles.insightTitle}>Top Performing Event</Text>
            <Text style={styles.insightDescription}>
              Tech Conference 2024 had the highest engagement rate at 28.5%
            </Text>
          </View>
        </View>
        
        <View style={styles.insightCard}>
          <Text style={styles.insightIcon}>📈</Text>
          <View style={styles.insightContent}>
            <Text style={styles.insightTitle}>Growth Opportunity</Text>
            <Text style={styles.insightDescription}>
              Content creation increased 45% - consider more visual campaigns
            </Text>
          </View>
        </View>
        
        <View style={styles.insightCard}>
          <Text style={styles.insightIcon}>💰</Text>
          <View style={styles.insightContent}>
            <Text style={styles.insightTitle}>ROI Highlight</Text>
            <Text style={styles.insightDescription}>
              Average ROI of 3.2x across all sponsored events this quarter
            </Text>
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
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#333',
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginTop: 4,
  },
  metricsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    padding: 20,
    gap: 12,
  },
  metricCard: {
    backgroundColor: '#fff',
    padding: 16,
    borderRadius: 12,
    width: (width - 52) / 2,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  metricNumber: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#28a745',
    marginBottom: 4,
  },
  metricLabel: {
    fontSize: 12,
    color: '#666',
    marginBottom: 4,
  },
  metricChange: {
    fontSize: 11,
    color: '#28a745',
    fontWeight: '500',
  },
  chartSection: {
    backgroundColor: '#fff',
    margin: 20,
    padding: 20,
    borderRadius: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  chartContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-end',
    height: 160,
    marginTop: 16,
  },
  chartBar: {
    alignItems: 'center',
    flex: 1,
  },
  bar: {
    backgroundColor: '#28a745',
    width: 20,
    borderRadius: 2,
    marginBottom: 8,
  },
  barValue: {
    fontSize: 10,
    color: '#333',
    fontWeight: '600',
    marginBottom: 4,
  },
  barLabel: {
    fontSize: 10,
    color: '#666',
  },
  section: {
    padding: 20,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#333',
  },
  viewAllText: {
    fontSize: 14,
    color: '#28a745',
    fontWeight: '600',
  },
  actionsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },
  actionCard: {
    backgroundColor: '#fff',
    padding: 20,
    borderRadius: 12,
    width: (width - 52) / 2,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  actionIcon: {
    fontSize: 32,
    marginBottom: 8,
  },
  actionTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
    textAlign: 'center',
  },
  campaignCard: {
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
  campaignHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  campaignName: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#333',
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
  campaignStats: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 16,
  },
  campaignStat: {
    flex: 1,
    alignItems: 'center',
  },
  campaignStatNumber: {
    fontSize: 14,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 2,
  },
  campaignStatLabel: {
    fontSize: 11,
    color: '#666',
  },
  campaignActions: {
    flexDirection: 'row',
    gap: 8,
  },
  campaignButton: {
    flex: 1,
    padding: 8,
    borderRadius: 6,
    borderWidth: 1,
    borderColor: '#28a745',
    alignItems: 'center',
  },
  primaryButton: {
    backgroundColor: '#28a745',
  },
  campaignButtonText: {
    fontSize: 12,
    color: '#28a745',
    fontWeight: '600',
  },
  primaryButtonText: {
    color: '#fff',
  },
  insightCard: {
    flexDirection: 'row',
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
  insightIcon: {
    fontSize: 24,
    marginRight: 12,
  },
  insightContent: {
    flex: 1,
  },
  insightTitle: {
    fontSize: 14,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 4,
  },
  insightDescription: {
    fontSize: 13,
    color: '#666',
    lineHeight: 18,
  },
});