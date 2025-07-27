/**
 * Export Page
 * Brand portal interface for data export and CRM integration management
 */

import React, { useState, useEffect } from 'react';

interface ExportRequest {
  id: string;
  eventId: string;
  dataType: string;
  format: string;
  status: string;
  fileUrl?: string;
  createdAt: string;
}

interface CRMIntegration {
  id: string;
  crmType: string;
  status: string;
  syncInterval: number;
}

export const ExportPage: React.FC = () => {
  const [exportRequests, setExportRequests] = useState<ExportRequest[]>([]);
  const [crmIntegrations, setCrmIntegrations] = useState<CRMIntegration[]>([]);
  const [showExportForm, setShowExportForm] = useState(false);
  const [showCRMForm, setShowCRMForm] = useState(false);
  const [newExport, setNewExport] = useState({
    eventId: '',
    dataType: 'attendance',
    format: 'csv',
  });
  const [newCRM, setNewCRM] = useState({
    crmType: 'salesforce',
    apiKey: '',
    apiSecret: '',
    webhookUrl: '',
    syncInterval: 60,
  });

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    // Simulate API calls
    setExportRequests([
      {
        id: 'export_1',
        eventId: 'event_1',
        dataType: 'attendance',
        format: 'csv',
        status: 'completed',
        fileUrl: 'https://exports.lynkr.com/export_1.csv',
        createdAt: '2024-03-15T10:30:00Z',
      },
    ]);

    setCrmIntegrations([
      {
        id: 'crm_1',
        crmType: 'salesforce',
        status: 'active',
        syncInterval: 60,
      },
    ]);
  };

  const handleCreateExport = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const response = await fetch('/api/v1/export/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newExport),
      });
      
      if (response.ok) {
        const exportReq = await response.json();
        setExportRequests([exportReq, ...exportRequests]);
        setNewExport({ eventId: '', dataType: 'attendance', format: 'csv' });
        setShowExportForm(false);
      }
    } catch (error) {
      console.error('Error creating export:', error);
    }
  };

  const handleCreateCRM = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const response = await fetch('/api/v1/crm/integrations', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newCRM),
      });
      
      if (response.ok) {
        const integration = await response.json();
        setCrmIntegrations([integration, ...crmIntegrations]);
        setNewCRM({ crmType: 'salesforce', apiKey: '', apiSecret: '', webhookUrl: '', syncInterval: 60 });
        setShowCRMForm(false);
      }
    } catch (error) {
      console.error('Error creating CRM integration:', error);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return '#28a745';
      case 'processing': return '#ffc107';
      case 'failed': return '#dc3545';
      default: return '#6c757d';
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.title}>Data Export & Integrations</h1>
      </div>

      <div style={styles.section}>
        <div style={styles.sectionHeader}>
          <h2 style={styles.sectionTitle}>Data Exports</h2>
          <button
            onClick={() => setShowExportForm(true)}
            style={styles.createButton}
          >
            Create Export
          </button>
        </div>

        {showExportForm && (
          <div style={styles.modal}>
            <div style={styles.modalContent}>
              <h3>Create Data Export</h3>
              <form onSubmit={handleCreateExport} style={styles.form}>
                <div style={styles.field}>
                  <label>Event ID</label>
                  <input
                    type="text"
                    value={newExport.eventId}
                    onChange={(e) => setNewExport({...newExport, eventId: e.target.value})}
                    style={styles.input}
                    required
                  />
                </div>
                <div style={styles.field}>
                  <label>Data Type</label>
                  <select
                    value={newExport.dataType}
                    onChange={(e) => setNewExport({...newExport, dataType: e.target.value})}
                    style={styles.input}
                  >
                    <option value="attendance">Attendance Data</option>
                    <option value="content">Content Data</option>
                    <option value="analytics">Analytics Data</option>
                    <option value="feedback">Feedback Data</option>
                  </select>
                </div>
                <div style={styles.field}>
                  <label>Format</label>
                  <select
                    value={newExport.format}
                    onChange={(e) => setNewExport({...newExport, format: e.target.value})}
                    style={styles.input}
                  >
                    <option value="csv">CSV</option>
                    <option value="json">JSON</option>
                  </select>
                </div>
                <div style={styles.buttonGroup}>
                  <button type="submit" style={styles.submitButton}>Create Export</button>
                  <button
                    type="button"
                    onClick={() => setShowExportForm(false)}
                    style={styles.cancelButton}
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        <div style={styles.exportList}>
          {exportRequests.map((req) => (
            <div key={req.id} style={styles.exportItem}>
              <div style={styles.exportInfo}>
                <h4 style={styles.exportTitle}>
                  {req.dataType.charAt(0).toUpperCase() + req.dataType.slice(1)} Export
                </h4>
                <p style={styles.exportDetails}>
                  Format: {req.format.toUpperCase()} | Event: {req.eventId}
                </p>
                <p style={styles.exportDate}>
                  Created: {new Date(req.createdAt).toLocaleDateString()}
                </p>
              </div>
              <div style={styles.exportActions}>
                <span
                  style={{
                    ...styles.statusBadge,
                    backgroundColor: getStatusColor(req.status),
                  }}
                >
                  {req.status.toUpperCase()}
                </span>
                {req.status === 'completed' && req.fileUrl && (
                  <a
                    href={req.fileUrl}
                    style={styles.downloadButton}
                    download
                  >
                    Download
                  </a>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div style={styles.section}>
        <div style={styles.sectionHeader}>
          <h2 style={styles.sectionTitle}>CRM Integrations</h2>
          <button
            onClick={() => setShowCRMForm(true)}
            style={styles.createButton}
          >
            Add Integration
          </button>
        </div>

        {showCRMForm && (
          <div style={styles.modal}>
            <div style={styles.modalContent}>
              <h3>Add CRM Integration</h3>
              <form onSubmit={handleCreateCRM} style={styles.form}>
                <div style={styles.field}>
                  <label>CRM Type</label>
                  <select
                    value={newCRM.crmType}
                    onChange={(e) => setNewCRM({...newCRM, crmType: e.target.value})}
                    style={styles.input}
                  >
                    <option value="salesforce">Salesforce</option>
                    <option value="hubspot">HubSpot</option>
                    <option value="mailchimp">Mailchimp</option>
                  </select>
                </div>
                <div style={styles.field}>
                  <label>API Key</label>
                  <input
                    type="text"
                    value={newCRM.apiKey}
                    onChange={(e) => setNewCRM({...newCRM, apiKey: e.target.value})}
                    style={styles.input}
                    required
                  />
                </div>
                <div style={styles.field}>
                  <label>API Secret</label>
                  <input
                    type="password"
                    value={newCRM.apiSecret}
                    onChange={(e) => setNewCRM({...newCRM, apiSecret: e.target.value})}
                    style={styles.input}
                    required
                  />
                </div>
                <div style={styles.field}>
                  <label>Webhook URL</label>
                  <input
                    type="url"
                    value={newCRM.webhookUrl}
                    onChange={(e) => setNewCRM({...newCRM, webhookUrl: e.target.value})}
                    style={styles.input}
                  />
                </div>
                <div style={styles.field}>
                  <label>Sync Interval (minutes)</label>
                  <input
                    type="number"
                    value={newCRM.syncInterval}
                    onChange={(e) => setNewCRM({...newCRM, syncInterval: Number(e.target.value)})}
                    style={styles.input}
                    min="15"
                  />
                </div>
                <div style={styles.buttonGroup}>
                  <button type="submit" style={styles.submitButton}>Add Integration</button>
                  <button
                    type="button"
                    onClick={() => setShowCRMForm(false)}
                    style={styles.cancelButton}
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        <div style={styles.integrationList}>
          {crmIntegrations.map((integration) => (
            <div key={integration.id} style={styles.integrationItem}>
              <div style={styles.integrationInfo}>
                <h4 style={styles.integrationTitle}>
                  {integration.crmType.charAt(0).toUpperCase() + integration.crmType.slice(1)}
                </h4>
                <p style={styles.integrationDetails}>
                  Sync Interval: {integration.syncInterval} minutes
                </p>
              </div>
              <div style={styles.integrationActions}>
                <span
                  style={{
                    ...styles.statusBadge,
                    backgroundColor: getStatusColor(integration.status),
                  }}
                >
                  {integration.status.toUpperCase()}
                </span>
                <button style={styles.actionButton}>Configure</button>
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
    marginBottom: '2rem',
  },
  title: {
    fontSize: '2rem',
    fontWeight: 'bold',
    color: '#333',
  },
  section: {
    backgroundColor: '#fff',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    marginBottom: '2rem',
  },
  sectionHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '1.5rem',
  },
  sectionTitle: {
    fontSize: '1.5rem',
    fontWeight: '600',
    color: '#333',
    margin: 0,
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
  exportList: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '1rem',
  },
  exportItem: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '1rem',
    border: '1px solid #e0e0e0',
    borderRadius: '4px',
  },
  exportInfo: {
    flex: 1,
  },
  exportTitle: {
    fontSize: '1.125rem',
    fontWeight: '600',
    margin: '0 0 0.25rem 0',
  },
  exportDetails: {
    fontSize: '0.875rem',
    color: '#666',
    margin: '0 0 0.25rem 0',
  },
  exportDate: {
    fontSize: '0.75rem',
    color: '#999',
    margin: 0,
  },
  exportActions: {
    display: 'flex',
    alignItems: 'center',
    gap: '1rem',
  },
  statusBadge: {
    padding: '0.25rem 0.75rem',
    borderRadius: '12px',
    color: 'white',
    fontSize: '0.75rem',
    fontWeight: 'bold',
  },
  downloadButton: {
    padding: '0.5rem 1rem',
    backgroundColor: '#007AFF',
    color: 'white',
    textDecoration: 'none',
    borderRadius: '4px',
    fontSize: '0.875rem',
  },
  integrationList: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '1rem',
  },
  integrationItem: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '1rem',
    border: '1px solid #e0e0e0',
    borderRadius: '4px',
  },
  integrationInfo: {
    flex: 1,
  },
  integrationTitle: {
    fontSize: '1.125rem',
    fontWeight: '600',
    margin: '0 0 0.25rem 0',
  },
  integrationDetails: {
    fontSize: '0.875rem',
    color: '#666',
    margin: 0,
  },
  integrationActions: {
    display: 'flex',
    alignItems: 'center',
    gap: '1rem',
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