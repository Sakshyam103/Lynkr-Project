/**
 * Conversion Funnel Chart
 * Visualization component for conversion funnel analytics
 */

import React from 'react';

interface FunnelStage {
  stage: string;
  users: number;
  conversions: number;
  rate: number;
}

interface ConversionFunnelChartProps {
  stages: FunnelStage[];
  totalUsers: number;
  revenue: number;
  roi: number;
}

export const ConversionFunnelChart: React.FC<ConversionFunnelChartProps> = ({
  stages,
  totalUsers,
  revenue,
  roi,
}) => {
  const getStageLabel = (stage: string) => {
    const labels: Record<string, string> = {
      attendance: 'Event Attendance',
      content_view: 'Content Views',
      engagement: 'User Engagement',
      website_visit: 'Website Visits',
      purchase: 'Purchases',
    };
    return labels[stage] || stage;
  };

  const getStageWidth = (users: number) => {
    return totalUsers > 0 ? (users / totalUsers) * 100 : 0;
  };

  const getStageColor = (index: number) => {
    const colors = ['#4CAF50', '#2196F3', '#FF9800', '#9C27B0', '#F44336'];
    return colors[index] || '#757575';
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h2 style={styles.title}>Conversion Funnel</h2>
        <div style={styles.metrics}>
          <div style={styles.metric}>
            <span style={styles.metricLabel}>Total Revenue:</span>
            <span style={styles.metricValue}>${revenue.toLocaleString()}</span>
          </div>
          <div style={styles.metric}>
            <span style={styles.metricLabel}>ROI:</span>
            <span style={[styles.metricValue, { color: roi >= 0 ? '#4CAF50' : '#F44336' }]}>
              {roi.toFixed(1)}%
            </span>
          </div>
        </div>
      </div>

      <div style={styles.funnel}>
        {stages.map((stage, index) => (
          <div key={stage.stage} style={styles.stageContainer}>
            <div
              style={{
                ...styles.stage,
                width: `${getStageWidth(stage.users)}%`,
                backgroundColor: getStageColor(index),
              }}
            >
              <div style={styles.stageContent}>
                <div style={styles.stageLabel}>{getStageLabel(stage.stage)}</div>
                <div style={styles.stageStats}>
                  <span style={styles.stageUsers}>{stage.users.toLocaleString()}</span>
                  <span style={styles.stageRate}>{stage.rate.toFixed(1)}%</span>
                </div>
              </div>
            </div>
            {index < stages.length - 1 && (
              <div style={styles.dropoff}>
                <span style={styles.dropoffText}>
                  -{((stages[index].users - stages[index + 1].users) / stages[index].users * 100).toFixed(1)}%
                </span>
              </div>
            )}
          </div>
        ))}
      </div>

      <div style={styles.summary}>
        <div style={styles.summaryItem}>
          <span style={styles.summaryLabel}>Total Users:</span>
          <span style={styles.summaryValue}>{totalUsers.toLocaleString()}</span>
        </div>
        <div style={styles.summaryItem}>
          <span style={styles.summaryLabel}>Final Conversion Rate:</span>
          <span style={styles.summaryValue}>
            {stages.length > 0 ? stages[stages.length - 1].rate.toFixed(2) : 0}%
          </span>
        </div>
        <div style={styles.summaryItem}>
          <span style={styles.summaryLabel}>Revenue per User:</span>
          <span style={styles.summaryValue}>
            ${totalUsers > 0 ? (revenue / totalUsers).toFixed(2) : '0.00'}
          </span>
        </div>
      </div>
    </div>
  );
};

const styles = {
  container: {
    backgroundColor: 'white',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '2rem',
  },
  title: {
    fontSize: '1.5rem',
    fontWeight: '600',
    color: '#333',
    margin: 0,
  },
  metrics: {
    display: 'flex',
    gap: '2rem',
  },
  metric: {
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'flex-end',
  },
  metricLabel: {
    fontSize: '0.875rem',
    color: '#666',
  },
  metricValue: {
    fontSize: '1.25rem',
    fontWeight: 'bold',
    color: '#333',
  },
  funnel: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '0.5rem',
    marginBottom: '2rem',
  },
  stageContainer: {
    display: 'flex',
    alignItems: 'center',
    gap: '1rem',
  },
  stage: {
    minHeight: '60px',
    display: 'flex',
    alignItems: 'center',
    borderRadius: '4px',
    transition: 'all 0.3s ease',
    position: 'relative' as const,
  },
  stageContent: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    width: '100%',
    padding: '0 1rem',
    color: 'white',
  },
  stageLabel: {
    fontSize: '1rem',
    fontWeight: '500',
  },
  stageStats: {
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'flex-end',
  },
  stageUsers: {
    fontSize: '1.125rem',
    fontWeight: 'bold',
  },
  stageRate: {
    fontSize: '0.875rem',
    opacity: 0.9,
  },
  dropoff: {
    display: 'flex',
    alignItems: 'center',
    minWidth: '60px',
  },
  dropoffText: {
    fontSize: '0.75rem',
    color: '#F44336',
    fontWeight: '500',
  },
  summary: {
    display: 'flex',
    justifyContent: 'space-around',
    padding: '1rem',
    backgroundColor: '#f8f9fa',
    borderRadius: '4px',
  },
  summaryItem: {
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'center',
  },
  summaryLabel: {
    fontSize: '0.875rem',
    color: '#666',
    marginBottom: '0.25rem',
  },
  summaryValue: {
    fontSize: '1.125rem',
    fontWeight: 'bold',
    color: '#333',
  },
};