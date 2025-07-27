/**
 * Discount Codes Management Page
 * Brand portal interface for managing discount codes
 */

import React, { useState, useEffect } from 'react';

interface DiscountCode {
  id: string;
  code: string;
  discountPct: number;
  usedCount: number;
  maxUses: number;
  expiresAt: string;
}

export const DiscountCodesPage: React.FC = () => {
  const [codes, setCodes] = useState<DiscountCode[]>([]);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newCode, setNewCode] = useState({
    eventId: '',
    discountPct: 10,
    maxUses: 100,
    expiresIn: 30,
  });

  useEffect(() => {
    loadCodes();
  }, []);

  const loadCodes = async () => {
    // Simulate API call
    setCodes([
      {
        id: 'discount_1',
        code: 'EVENT20',
        discountPct: 20,
        usedCount: 45,
        maxUses: 100,
        expiresAt: '2024-12-31',
      },
      {
        id: 'discount_2',
        code: 'SPECIAL15',
        discountPct: 15,
        usedCount: 23,
        maxUses: 50,
        expiresAt: '2024-11-30',
      },
    ]);
  };

  const handleCreateCode = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const response = await fetch('/api/v1/discount/generate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newCode),
      });
      
      if (response.ok) {
        const code = await response.json();
        setCodes([code, ...codes]);
        setNewCode({ eventId: '', discountPct: 10, maxUses: 100, expiresIn: 30 });
        setShowCreateForm(false);
      }
    } catch (error) {
      console.error('Error creating code:', error);
    }
  };

  const getUsagePercentage = (used: number, max: number) => {
    return max > 0 ? (used / max) * 100 : 0;
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.title}>Discount Codes</h1>
        <button
          onClick={() => setShowCreateForm(true)}
          style={styles.createButton}
        >
          Create Code
        </button>
      </div>

      {showCreateForm && (
        <div style={styles.modal}>
          <div style={styles.modalContent}>
            <h2>Create Discount Code</h2>
            <form onSubmit={handleCreateCode} style={styles.form}>
              <div style={styles.field}>
                <label>Event ID</label>
                <input
                  type="text"
                  value={newCode.eventId}
                  onChange={(e) => setNewCode({...newCode, eventId: e.target.value})}
                  style={styles.input}
                  required
                />
              </div>
              <div style={styles.field}>
                <label>Discount Percentage</label>
                <input
                  type="number"
                  value={newCode.discountPct}
                  onChange={(e) => setNewCode({...newCode, discountPct: Number(e.target.value)})}
                  style={styles.input}
                  min="1"
                  max="100"
                  required
                />
              </div>
              <div style={styles.field}>
                <label>Max Uses</label>
                <input
                  type="number"
                  value={newCode.maxUses}
                  onChange={(e) => setNewCode({...newCode, maxUses: Number(e.target.value)})}
                  style={styles.input}
                  min="1"
                  required
                />
              </div>
              <div style={styles.field}>
                <label>Expires In (days)</label>
                <input
                  type="number"
                  value={newCode.expiresIn}
                  onChange={(e) => setNewCode({...newCode, expiresIn: Number(e.target.value)})}
                  style={styles.input}
                  min="1"
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

      <div style={styles.codesGrid}>
        {codes.map((code) => (
          <div key={code.id} style={styles.codeCard}>
            <div style={styles.codeHeader}>
              <h3 style={styles.codeName}>{code.code}</h3>
              <span style={styles.discount}>{code.discountPct}% OFF</span>
            </div>
            
            <div style={styles.codeStats}>
              <div style={styles.stat}>
                <span style={styles.statLabel}>Used:</span>
                <span style={styles.statValue}>{code.usedCount}/{code.maxUses}</span>
              </div>
              <div style={styles.stat}>
                <span style={styles.statLabel}>Expires:</span>
                <span style={styles.statValue}>{code.expiresAt}</span>
              </div>
            </div>

            <div style={styles.progressBar}>
              <div
                style={{
                  ...styles.progress,
                  width: `${getUsagePercentage(code.usedCount, code.maxUses)}%`,
                }}
              />
            </div>

            <div style={styles.codeActions}>
              <button style={styles.actionButton}>View Analytics</button>
              <button style={styles.actionButton}>Copy Code</button>
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
  codesGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
    gap: '1.5rem',
  },
  codeCard: {
    backgroundColor: 'white',
    padding: '1.5rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  },
  codeHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '1rem',
  },
  codeName: {
    fontSize: '1.5rem',
    fontWeight: 'bold',
    margin: 0,
    fontFamily: 'monospace',
  },
  discount: {
    backgroundColor: '#28a745',
    color: 'white',
    padding: '0.25rem 0.75rem',
    borderRadius: '12px',
    fontSize: '0.875rem',
    fontWeight: 'bold',
  },
  codeStats: {
    display: 'flex',
    justifyContent: 'space-between',
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
    fontSize: '1rem',
    fontWeight: '600',
    color: '#333',
  },
  progressBar: {
    width: '100%',
    height: '8px',
    backgroundColor: '#e9ecef',
    borderRadius: '4px',
    overflow: 'hidden',
    marginBottom: '1rem',
  },
  progress: {
    height: '100%',
    backgroundColor: '#007AFF',
    transition: 'width 0.3s ease',
  },
  codeActions: {
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