/**
 * Sentiment Visualization
 * Charts and visualizations for sentiment analysis data
 */

import React from 'react';
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface SentimentData {
  totalAnalyses: number;
  averageScore: number;
  distribution: {
    positive: number;
    negative: number;
    neutral: number;
  };
}

interface SentimentVisualizationProps {
  data: SentimentData;
}

export const SentimentVisualization: React.FC<SentimentVisualizationProps> = ({ data }) => {
  const pieData = [
    { name: 'Positive', value: data.distribution.positive, color: '#2ed573' },
    { name: 'Neutral', value: data.distribution.neutral, color: '#ffa502' },
    { name: 'Negative', value: data.distribution.negative, color: '#ff4757' },
  ];

  const barData = [
    { sentiment: 'Positive', count: data.distribution.positive, color: '#2ed573' },
    { sentiment: 'Neutral', count: data.distribution.neutral, color: '#ffa502' },
    { sentiment: 'Negative', count: data.distribution.negative, color: '#ff4757' },
  ];

  const getSentimentLabel = (score: number) => {
    if (score > 0.1) return 'Positive';
    if (score < -0.1) return 'Negative';
    return 'Neutral';
  };

  const getSentimentColor = (score: number) => {
    if (score > 0.1) return '#2ed573';
    if (score < -0.1) return '#ff4757';
    return '#ffa502';
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h2 style={styles.title}>Sentiment Analysis</h2>
        <div style={styles.summary}>
          <div style={styles.summaryItem}>
            <span style={styles.summaryLabel}>Total Analyses:</span>
            <span style={styles.summaryValue}>{data.totalAnalyses}</span>
          </div>
          <div style={styles.summaryItem}>
            <span style={styles.summaryLabel}>Overall Sentiment:</span>
            <span 
              style={{
                ...styles.summaryValue,
                color: getSentimentColor(data.averageScore),
                fontWeight: 'bold',
              }}
            >
              {getSentimentLabel(data.averageScore)}
            </span>
          </div>
          <div style={styles.summaryItem}>
            <span style={styles.summaryLabel}>Average Score:</span>
            <span style={styles.summaryValue}>
              {data.averageScore.toFixed(2)}
            </span>
          </div>
        </div>
      </div>

      <div style={styles.chartsContainer}>
        <div style={styles.chartSection}>
          <h3 style={styles.chartTitle}>Sentiment Distribution</h3>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                innerRadius={60}
                outerRadius={100}
                paddingAngle={5}
                dataKey="value"
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
              >
                {pieData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div style={styles.chartSection}>
          <h3 style={styles.chartTitle}>Sentiment Counts</h3>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={barData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="sentiment" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="count" fill="#8884d8">
                {barData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div style={styles.insights}>
        <h3 style={styles.insightsTitle}>Key Insights</h3>
        <ul style={styles.insightsList}>
          <li>
            {data.distribution.positive > data.distribution.negative 
              ? '✅ Overall positive sentiment detected'
              : '⚠️ Mixed or negative sentiment detected'
            }
          </li>
          <li>
            Most common sentiment: {
              Object.entries(data.distribution)
                .sort(([,a], [,b]) => b - a)[0][0]
                .charAt(0).toUpperCase() + 
              Object.entries(data.distribution)
                .sort(([,a], [,b]) => b - a)[0][0]
                .slice(1)
            }
          </li>
          <li>
            Engagement level: {
              data.totalAnalyses > 50 ? 'High' : 
              data.totalAnalyses > 20 ? 'Medium' : 'Low'
            } ({data.totalAnalyses} responses)
          </li>
        </ul>
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
    marginBottom: '1.5rem',
  },
  title: {
    fontSize: '1.5rem',
    fontWeight: '600',
    color: '#333',
    marginBottom: '1rem',
  },
  summary: {
    display: 'flex',
    gap: '2rem',
    flexWrap: 'wrap' as const,
  },
  summaryItem: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '0.25rem',
  },
  summaryLabel: {
    fontSize: '0.875rem',
    color: '#666',
  },
  summaryValue: {
    fontSize: '1.125rem',
    fontWeight: '600',
    color: '#333',
  },
  chartsContainer: {
    display: 'grid',
    gridTemplateColumns: '1fr 1fr',
    gap: '2rem',
    marginBottom: '1.5rem',
  },
  chartSection: {
    textAlign: 'center' as const,
  },
  chartTitle: {
    fontSize: '1.125rem',
    fontWeight: '500',
    color: '#333',
    marginBottom: '1rem',
  },
  insights: {
    backgroundColor: '#f8f9fa',
    padding: '1rem',
    borderRadius: '6px',
  },
  insightsTitle: {
    fontSize: '1rem',
    fontWeight: '600',
    color: '#333',
    marginBottom: '0.75rem',
  },
  insightsList: {
    margin: 0,
    paddingLeft: '1.25rem',
    color: '#555',
  },
};