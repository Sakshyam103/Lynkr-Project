/**
 * Campaign Management Page
 * Interface for creating and managing brand campaigns
 */

import React, { useState, useEffect } from 'react';

interface Campaign {
  id: string;
  name: string;
  status: 'draft' | 'active' | 'completed';
  startDate: string;
  endDate: string;
  budget: number;
  attendees: number;
}

export const CampaignManagementPage: React.FC = () => {
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newCampaign, setNewCampaign] = useState({
    name: '',
    startDate: '',
    endDate: '',
    budget: 0,
  });

  useEffect(() => {
    loadCampaigns();
  }, []);

  const loadCampaigns = async () => {
    // Simulate API call
    setCampaigns([
      {
        id: '1',
        name: 'Tech Conference 2024',
        status: 'active',
        startDate: '2024-03-15',
        endDate: '2024-03-17',
        budget: 50000,
        attendees: 1247,
      },
      {
        id: '2',
        name: 'Product Launch Event',
        status: 'completed',
        startDate: '2024-02-10',
        endDate: '2024-02-12',
        budget: 75000,
        attendees: 890,
      },
    ]);
  };

  const handleCreateCampaign = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const campaign: Campaign = {
      id: Date.now().toString(),
      ...newCampaign,
      status: 'draft',
      attendees: 0,
    };

    setCampaigns([campaign, ...campaigns]);
    setNewCampaign({ name: '', startDate: '', endDate: '', budget: 0 });
    setShowCreateForm(false);
  };

  const getStatusColor = (status: Campaign['status']) => {
    switch (status) {
      case 'active': return '#28a745';
      case 'completed': return '#6c757d';
      case 'draft': return '#ffc107';
      default: return '#6c757d';
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.title}>Campaign Management</h1>
        <button
          onClick={() => setShowCreateForm(true)}
          style={styles.createButton}
        >
          Create Campaign
        </button>
      </div>

      {showCreateForm && (
        <div style={styles.modal}>
          <div style={styles.modalContent}>
            <h2>Create New Campaign</h2>
            <form onSubmit={handleCreateCampaign} style={styles.form}>
              <div style={styles.field}>
                <label>Campaign Name</label>
                <input
                  type="text"
                  value={newCampaign.name}
                  onChange={(e) => setNewCampaign({...newCampaign, name: e.target.value})}
                  style={styles.input}
                  required
                />
              </div>
              <div style={styles.field}>
                <label>Start Date</label>
                <input
                  type="date"
                  value={newCampaign.startDate}
                  onChange={(e) => setNewCampaign({...newCampaign, startDate: e.target.value})}
                  style={styles.input}
                  required
                />
              </div>
              <div style={styles.field}>
                <label>End Date</label>
                <input
                  type="date"
                  value={newCampaign.endDate}
                  onChange={(e) => setNewCampaign({...newCampaign, endDate: e.target.value})}
                  style={styles.input}
                  required
                />
              </div>
              <div style={styles.field}>
                <label>Budget ($)</label>
                <input
                  type="number"
                  value={newCampaign.budget}
                  onChange={(e) => setNewCampaign({...newCampaign, budget: Number(e.target.value)})}
                  style={styles.input}
                  required
                />
              </div>
              <div style={styles.buttonGroup}>
                <button type="submit" style={styles.submitButton}>Create</button>
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  style={styles.cancelButton}
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div style={styles.campaignGrid}>
        {campaigns.map((campaign) => (
          <div key={campaign.id} style={styles.campaignCard}>
            <div style={styles.campaignHeader}>
              <h3 style={styles.campaignName}>{campaign.name}</h3>
              <span
                style={{
                  ...styles.statusBadge,
                  backgroundColor: getStatusColor(campaign.status),
                }}
              >
                {campaign.status.toUpperCase()}
              </span>
            </div>
            <div style={styles.campaignDetails}>
              <p><strong>Duration:</strong> {campaign.startDate} to {campaign.endDate}</p>
              <p><strong>Budget:</strong> ${campaign.budget.toLocaleString()}</p>
              <p><strong>Attendees:</strong> {campaign.attendees.toLocaleString()}</p>
            </div>
            <div style={styles.campaignActions}>
              <button style={styles.actionButton}>View Analytics</button>
              <button style={styles.actionButton}>Edit</button>
            </div>
          </div>
        ))}
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
  createButton: {
    padding: '0.75rem 1.5rem',
    backgroundColor: '#007AFF',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '1rem',
  },
  modal: {
    position: 'fixed' as const,
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.5)',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1000,
  },
  modalContent: {
    backgroundColor: 'white',
    padding: '2rem',
    borderRadius: '8px',
    width: '90%',
    maxWidth: '500px',
  },
  form: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '1rem',
  },
  field: {
    display: 'flex',
    flexDirection: 'column' as const,
  },
  input: {
    padding: '0.75rem',
    border: '1px solid #ddd',
    borderRadius: '4px',
    marginTop: '0.25rem',
  },
  buttonGroup: {
    display: 'flex',
    gap: '1rem',
    marginTop: '1rem',
  },
  submitButton: {
    padding: '0.75rem 1.5rem',
    backgroundColor: '#28a745',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  cancelButton: {
    padding: '0.75rem 1.5rem',
    backgroundColor: '#6c757d',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  campaignGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(350px, 1fr))',
    gap: '1.5rem',
  },
  campaignCard: {
    backgroundColor: 'white',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  },
  campaignHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '1rem',
  },
  campaignName: {
    fontSize: '1.2rem',
    fontWeight: '600',
    margin: 0,
  },
  statusBadge: {
    padding: '0.25rem 0.75rem',
    borderRadius: '12px',
    color: 'white',
    fontSize: '0.75rem',
    fontWeight: 'bold',
  },
  campaignDetails: {
    marginBottom: '1rem',
  },
  campaignActions: {
    display: 'flex',
    gap: '0.5rem',
  },
  actionButton: {
    padding: '0.5rem 1rem',
    backgroundColor: '#f8f9fa',
    border: '1px solid #dee2e6',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '0.875rem',
  },
};