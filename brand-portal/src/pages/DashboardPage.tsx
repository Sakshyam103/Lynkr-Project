/**
 * Brand Dashboard Page
 * Main dashboard with analytics overview and navigation
 */

import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { AuthService } from '../services/AuthService';

interface DashboardStats {
  totalAttendees: number;
  contentPieces: number;
  engagementRate: number;
  conversionRate: number;
}

interface AttendanceData {
  date: string;
  attendees: number;
}

export const DashboardPage: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats>({
    totalAttendees: 0,
    contentPieces: 0,
    engagementRate: 0,
    conversionRate: 0,
  });
  const [attendanceData, setAttendanceData] = useState<AttendanceData[]>([]);
  const [loading, setLoading] = useState(true);
  const user = AuthService.getCurrentUser();

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      // Simulate API calls - in real implementation, these would be actual API calls
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setStats({
        totalAttendees: 1247,
        contentPieces: 89,
        engagementRate: 23.5,
        conversionRate: 4.2,
      });

      setAttendanceData([
        { date: '2024-01', attendees: 120 },
        { date: '2024-02', attendees: 180 },
        { date: '2024-03', attendees: 250 },
        { date: '2024-04', attendees: 320 },
        { date: '2024-05', attendees: 280 },
        { date: '2024-06', attendees: 380 },
      ]);
    } catch (error) {
      console.error('Error loading dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    AuthService.logout();
    window.location.href = '/login';
  };

  if (loading) {
    return <div style={styles.loading}>Loading dashboard...</div>;
  }

  return (
    <div style={styles.container}>
      <header style={styles.header}>
        <h1 style={styles.title}>Brand Analytics Dashboard</h1>
        <div style={styles.userInfo}>
          <span>Welcome, {user?.name}</span>
          <button onClick={handleLogout} style={styles.logoutButton}>
            Logout
          </button>
        </div>
      </header>

      <div style={styles.statsGrid}>
        <div style={styles.statCard}>
          <h3 style={styles.statTitle}>Total Attendees</h3>
          <p style={styles.statValue}>{stats.totalAttendees.toLocaleString()}</p>
        </div>
        <div style={styles.statCard}>
          <h3 style={styles.statTitle}>Content Pieces</h3>
          <p style={styles.statValue}>{stats.contentPieces}</p>
        </div>
        <div style={styles.statCard}>
          <h3 style={styles.statTitle}>Engagement Rate</h3>
          <p style={styles.statValue}>{stats.engagementRate}%</p>
        </div>
        <div style={styles.statCard}>
          <h3 style={styles.statTitle}>Conversion Rate</h3>
          <p style={styles.statValue}>{stats.conversionRate}%</p>
        </div>
      </div>

      <div style={styles.chartSection}>
        <h2 style={styles.sectionTitle}>Event Attendance Trends</h2>
        <div style={styles.chartContainer}>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={attendanceData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="attendees" fill="#007AFF" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div style={styles.quickActions}>
        <h2 style={styles.sectionTitle}>Quick Actions</h2>
        <div style={styles.actionGrid}>
          <button style={styles.actionButton}>View Content Gallery</button>
          <button style={styles.actionButton}>Create Campaign</button>
          <button style={styles.actionButton}>Export Analytics</button>
          <button style={styles.actionButton}>Manage Events</button>
        </div>
      </div>
    </div>
  );
};

const styles = {
  container: {
    padding: '2rem',
    backgroundColor: '#f8f9fa',
    minHeight: '100vh',
  },
  loading: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '100vh',
    fontSize: '1.2rem',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '2rem',
    backgroundColor: 'white',
    padding: '1rem 2rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  },
  title: {
    fontSize: '1.8rem',
    fontWeight: 'bold',
    color: '#333',
    margin: 0,
  },
  userInfo: {
    display: 'flex',
    alignItems: 'center',
    gap: '1rem',
  },
  logoutButton: {
    padding: '0.5rem 1rem',
    backgroundColor: '#dc3545',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  statsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
    gap: '1rem',
    marginBottom: '2rem',
  },
  statCard: {
    backgroundColor: 'white',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    textAlign: 'center' as const,
  },
  statTitle: {
    fontSize: '0.9rem',
    color: '#666',
    margin: '0 0 0.5rem 0',
    textTransform: 'uppercase' as const,
  },
  statValue: {
    fontSize: '2rem',
    fontWeight: 'bold',
    color: '#007AFF',
    margin: 0,
  },
  chartSection: {
    backgroundColor: 'white',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    marginBottom: '2rem',
  },
  sectionTitle: {
    fontSize: '1.3rem',
    fontWeight: '600',
    color: '#333',
    marginBottom: '1rem',
  },
  chartContainer: {
    height: '300px',
  },
  quickActions: {
    backgroundColor: 'white',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  },
  actionGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
    gap: '1rem',
  },
  actionButton: {
    padding: '1rem',
    backgroundColor: '#f8f9fa',
    border: '1px solid #dee2e6',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '1rem',
    transition: 'background-color 0.2s',
  },
};