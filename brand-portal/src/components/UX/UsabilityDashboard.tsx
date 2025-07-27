/**
 * Usability Dashboard
 * Brand portal interface for UX metrics and usability testing
 */

import React, { useState, useEffect } from 'react';

interface UsabilityMetrics {
  taskCompletionRate: number;
  averageTaskTime: number;
  errorRate: number;
  userSatisfaction: number;
  navigationEfficiency: number;
}

interface PainPoint {
  screen: string;
  element: string;
  errorCount: number;
  criticalRate: number;
  severity: string;
}

export const UsabilityDashboard: React.FC = () => {
  const [metrics, setMetrics] = useState<UsabilityMetrics>({
    taskCompletionRate: 0,
    averageTaskTime: 0,
    errorRate: 0,
    userSatisfaction: 0,
    navigationEfficiency: 0,
  });
  const [painPoints, setPainPoints] = useState<PainPoint[]>([]);
  const [timeframe, setTimeframe] = useState('7d');

  useEffect(() => {
    loadUsabilityData();
  }, [timeframe]);

  const loadUsabilityData = async () => {
    try {
      // Simulate API calls
      setMetrics({
        taskCompletionRate: 87.5,
        averageTaskTime: 45.2,
        errorRate: 3.8,
        userSatisfaction: 4.2,
        navigationEfficiency: 92.1,
      });

      setPainPoints([
        {
          screen: 'EventCheckIn',
          element: 'check_in_button',
          errorCount: 23,
          criticalRate: 0.15,
          severity: 'high',
        },
        {
          screen: 'ContentCreation',
          element: 'photo_upload',
          errorCount: 18,
          criticalRate: 0.22,
          severity: 'medium',
        },
      ]);
    } catch (error) {
      console.error('Error loading usability data:', error);
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return '#dc3545';
      case 'high': return '#fd7e14';
      case 'medium': return '#ffc107';
      case 'low': return '#28a745';
      default: return '#6c757d';
    }
  };

  const getMetricColor = (value: number, isInverted = false) => {
    const threshold = isInverted ? 5 : 80;
    if (isInverted) {
      return value < threshold ? '#28a745' : value < threshold * 2 ? '#ffc107' : '#dc3545';
    }
    return value > threshold ? '#28a745' : value > threshold * 0.7 ? '#ffc107' : '#dc3545';
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.title}>Usability Dashboard</h1>
        <select
          value={timeframe}
          onChange={(e) => setTimeframe(e.target.value)}
          style={styles.timeframeSelect}
        >
          <option value="24h">Last 24 Hours</option>
          <option value="7d">Last 7 Days</option>
          <option value="30d">Last 30 Days</option>
        </select>
      </div>

      <div style={styles.metricsGrid}>
        <div style={styles.metricCard}>
          <h3 style={styles.metricTitle}>Task Completion Rate</h3>
          <div
            style={{
              ...styles.metricValue,
              color: getMetricColor(metrics.taskCompletionRate),
            }}
          >
            {metrics.taskCompletionRate.toFixed(1)}%
          </div>
          <p style={styles.metricDescription}>
            Percentage of users who complete their intended tasks
          </p>
        </div>

        <div style={styles.metricCard}>
          <h3 style={styles.metricTitle}>Average Task Time</h3>
          <div
            style={{
              ...styles.metricValue,
              color: getMetricColor(metrics.averageTaskTime, true),
            }}
          >
            {metrics.averageTaskTime.toFixed(1)}s
          </div>
          <p style={styles.metricDescription}>
            Average time to complete primary tasks
          </p>
        </div>

        <div style={styles.metricCard}>
          <h3 style={styles.metricTitle}>Error Rate</h3>
          <div
            style={{
              ...styles.metricValue,
              color: getMetricColor(metrics.errorRate, true),
            }}
          >
            {metrics.errorRate.toFixed(1)}%
          </div>
          <p style={styles.metricDescription}>
            Percentage of user sessions with errors
          </p>
        </div>

        <div style={styles.metricCard}>
          <h3 style={styles.metricTitle}>User Satisfaction</h3>
          <div
            style={{
              ...styles.metricValue,
              color: getMetricColor(metrics.userSatisfaction * 20),
            }}
          >
            {metrics.userSatisfaction.toFixed(1)}/5
          </div>
          <p style={styles.metricDescription}>
            Average user satisfaction rating
          </p>
        </div>

        <div style={styles.metricCard}>
          <h3 style={styles.metricTitle}>Navigation Efficiency</h3>
          <div
            style={{
              ...styles.metricValue,
              color: getMetricColor(metrics.navigationEfficiency),
            }}
          >
            {metrics.navigationEfficiency.toFixed(1)}%
          </div>
          <p style={styles.metricDescription}>
            Efficiency of user navigation patterns
          </p>
        </div>
      </div>

      <div style={styles.section}>
        <h2 style={styles.sectionTitle}>Pain Points Analysis</h2>
        <div style={styles.painPointsList}>
          {painPoints.map((point, index) => (
            <div key={index} style={styles.painPointCard}>
              <div style={styles.painPointHeader}>
                <h4 style={styles.painPointTitle}>
                  {point.screen} - {point.element}
                </h4>
                <span
                  style={{
                    ...styles.severityBadge,
                    backgroundColor: getSeverityColor(point.severity),
                  }}
                >
                  {point.severity.toUpperCase()}
                </span>
              </div>
              <div style={styles.painPointStats}>
                <div style={styles.stat}>
                  <span style={styles.statLabel}>Error Count:</span>
                  <span style={styles.statValue}>{point.errorCount}</span>
                </div>
                <div style={styles.stat}>
                  <span style={styles.statLabel}>Critical Rate:</span>
                  <span style={styles.statValue}>{(point.criticalRate * 100).toFixed(1)}%</span>
                </div>
              </div>
              <div style={styles.painPointActions}>
                <button style={styles.actionButton}>View Details</button>
                <button style={styles.actionButton}>Create Issue</button>
              </div>
            </div>
          ))}
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
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '2rem',
  },
  title: {
    fontSize: '2rem',
    fontWeight: 'bold',
    color: '#333',
  },
  timeframeSelect: {
    padding: '0.5rem 1rem',
    border: '1px solid #ddd',
    borderRadius: '4px',
    fontSize: '1rem',
  },
  metricsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
    gap: '1.5rem',
    marginBottom: '2rem',
  },
  metricCard: {
    backgroundColor: '#fff',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    textAlign: 'center' as const,
  },
  metricTitle: {
    fontSize: '1rem',
    color: '#666',
    marginBottom: '0.5rem',
    fontWeight: '500',
  },
  metricValue: {
    fontSize: '2.5rem',
    fontWeight: 'bold',
    marginBottom: '0.5rem',
  },
  metricDescription: {
    fontSize: '0.875rem',
    color: '#888',
    margin: 0,
  },
  section: {
    backgroundColor: '#fff',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    marginBottom: '2rem',
  },
  sectionTitle: {
    fontSize: '1.5rem',
    fontWeight: '600',
    color: '#333',
    marginBottom: '1rem',
  },
  painPointsList: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '1rem',
  },
  painPointCard: {
    border: '1px solid #e0e0e0',
    borderRadius: '8px',
    padding: '1rem',
  },
  painPointHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '0.5rem',
  },
  painPointTitle: {
    fontSize: '1.125rem',
    fontWeight: '600',
    color: '#333',
    margin: 0,
  },
  severityBadge: {
    padding: '0.25rem 0.75rem',
    borderRadius: '12px',
    color: 'white',
    fontSize: '0.75rem',
    fontWeight: 'bold',
  },
  painPointStats: {
    display: 'flex',
    gap: '2rem',
    marginBottom: '1rem',
  },
  stat: {
    display: 'flex',
    flexDirection: 'column' as const,
  },
  statLabel: {
    fontSize: '0.875rem',
    color: '#666',
  },
  statValue: {
    fontSize: '1.125rem',
    fontWeight: '600',
    color: '#333',
  },
  painPointActions: {
    display: 'flex',
    gap: '0.5rem',
  },
  actionButton: {
    padding: '0.5rem 1rem',
    backgroundColor: '#007AFF',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '0.875rem',
  },
};