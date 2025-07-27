/**
 * Content Gallery Page
 * Displays user-generated content accessible to brands
 */

import React, { useState, useEffect } from 'react';

interface ContentItem {
  id: string;
  mediaUrl: string;
  mediaType: 'photo' | 'video';
  caption: string;
  tags: string[];
  createdAt: string;
  eventName: string;
  engagement: {
    views: number;
    shares: number;
    likes: number;
  };
}

export const ContentGalleryPage: React.FC = () => {
  const [content, setContent] = useState<ContentItem[]>([]);
  const [selectedContent, setSelectedContent] = useState<ContentItem | null>(null);
  const [filter, setFilter] = useState('all');

  useEffect(() => {
    loadContent();
  }, []);

  const loadContent = async () => {
    // Simulate API call
    setContent([
      {
        id: '1',
        mediaUrl: 'https://via.placeholder.com/300x300',
        mediaType: 'photo',
        caption: 'Amazing product demo at the tech conference!',
        tags: ['tech', 'product-demo', 'innovation'],
        createdAt: '2024-03-15T10:30:00Z',
        eventName: 'Tech Conference 2024',
        engagement: { views: 1250, shares: 45, likes: 189 },
      },
      {
        id: '2',
        mediaUrl: 'https://via.placeholder.com/300x300',
        mediaType: 'photo',
        caption: 'Great networking session with industry leaders',
        tags: ['networking', 'conference', 'business'],
        createdAt: '2024-03-15T14:20:00Z',
        eventName: 'Tech Conference 2024',
        engagement: { views: 890, shares: 23, likes: 156 },
      },
    ]);
  };

  const filteredContent = content.filter(item => {
    if (filter === 'all') return true;
    return item.mediaType === filter;
  });

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.title}>Content Gallery</h1>
        <div style={styles.filters}>
          <button
            onClick={() => setFilter('all')}
            style={filter === 'all' ? styles.activeFilter : styles.filter}
          >
            All
          </button>
          <button
            onClick={() => setFilter('photo')}
            style={filter === 'photo' ? styles.activeFilter : styles.filter}
          >
            Photos
          </button>
          <button
            onClick={() => setFilter('video')}
            style={filter === 'video' ? styles.activeFilter : styles.filter}
          >
            Videos
          </button>
        </div>
      </div>

      <div style={styles.gallery}>
        {filteredContent.map((item) => (
          <div
            key={item.id}
            style={styles.contentCard}
            onClick={() => setSelectedContent(item)}
          >
            <div style={styles.mediaContainer}>
              <img
                src={item.mediaUrl}
                alt={item.caption}
                style={styles.media}
              />
              {item.mediaType === 'video' && (
                <div style={styles.playIcon}>▶</div>
              )}
            </div>
            <div style={styles.contentInfo}>
              <p style={styles.caption}>{item.caption}</p>
              <div style={styles.tags}>
                {item.tags.map((tag) => (
                  <span key={tag} style={styles.tag}>#{tag}</span>
                ))}
              </div>
              <div style={styles.engagement}>
                <span>👁 {item.engagement.views}</span>
                <span>❤️ {item.engagement.likes}</span>
                <span>📤 {item.engagement.shares}</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {selectedContent && (
        <div style={styles.modal} onClick={() => setSelectedContent(null)}>
          <div style={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <div style={styles.modalHeader}>
              <h2>Content Details</h2>
              <button
                onClick={() => setSelectedContent(null)}
                style={styles.closeButton}
              >
                ×
              </button>
            </div>
            <div style={styles.modalBody}>
              <img
                src={selectedContent.mediaUrl}
                alt={selectedContent.caption}
                style={styles.modalImage}
              />
              <div style={styles.modalInfo}>
                <p><strong>Caption:</strong> {selectedContent.caption}</p>
                <p><strong>Event:</strong> {selectedContent.eventName}</p>
                <p><strong>Created:</strong> {new Date(selectedContent.createdAt).toLocaleDateString()}</p>
                <div style={styles.modalTags}>
                  {selectedContent.tags.map((tag) => (
                    <span key={tag} style={styles.tag}>#{tag}</span>
                  ))}
                </div>
                <div style={styles.modalEngagement}>
                  <div>Views: {selectedContent.engagement.views}</div>
                  <div>Likes: {selectedContent.engagement.likes}</div>
                  <div>Shares: {selectedContent.engagement.shares}</div>
                </div>
                <div style={styles.modalActions}>
                  <button style={styles.actionButton}>Request Usage Rights</button>
                  <button style={styles.actionButton}>Download</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
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
  filters: {
    display: 'flex',
    gap: '0.5rem',
  },
  filter: {
    padding: '0.5rem 1rem',
    backgroundColor: 'white',
    border: '1px solid #dee2e6',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  activeFilter: {
    padding: '0.5rem 1rem',
    backgroundColor: '#007AFF',
    color: 'white',
    border: '1px solid #007AFF',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  gallery: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
    gap: '1.5rem',
  },
  contentCard: {
    backgroundColor: 'white',
    borderRadius: '8px',
    overflow: 'hidden',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    cursor: 'pointer',
    transition: 'transform 0.2s',
  },
  mediaContainer: {
    position: 'relative' as const,
    height: '200px',
    overflow: 'hidden',
  },
  media: {
    width: '100%',
    height: '100%',
    objectFit: 'cover' as const,
  },
  playIcon: {
    position: 'absolute' as const,
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    backgroundColor: 'rgba(0,0,0,0.7)',
    color: 'white',
    borderRadius: '50%',
    width: '50px',
    height: '50px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    fontSize: '1.2rem',
  },
  contentInfo: {
    padding: '1rem',
  },
  caption: {
    fontSize: '0.9rem',
    color: '#333',
    marginBottom: '0.5rem',
    lineHeight: '1.4',
  },
  tags: {
    display: 'flex',
    flexWrap: 'wrap' as const,
    gap: '0.25rem',
    marginBottom: '0.5rem',
  },
  tag: {
    fontSize: '0.75rem',
    backgroundColor: '#e9ecef',
    color: '#495057',
    padding: '0.25rem 0.5rem',
    borderRadius: '12px',
  },
  engagement: {
    display: 'flex',
    gap: '1rem',
    fontSize: '0.8rem',
    color: '#666',
  },
  modal: {
    position: 'fixed' as const,
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.8)',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1000,
  },
  modalContent: {
    backgroundColor: 'white',
    borderRadius: '8px',
    width: '90%',
    maxWidth: '800px',
    maxHeight: '90vh',
    overflow: 'auto',
  },
  modalHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '1rem',
    borderBottom: '1px solid #dee2e6',
  },
  closeButton: {
    background: 'none',
    border: 'none',
    fontSize: '1.5rem',
    cursor: 'pointer',
  },
  modalBody: {
    display: 'flex',
    gap: '1rem',
    padding: '1rem',
  },
  modalImage: {
    width: '50%',
    height: 'auto',
    borderRadius: '4px',
  },
  modalInfo: {
    flex: 1,
  },
  modalTags: {
    display: 'flex',
    flexWrap: 'wrap' as const,
    gap: '0.25rem',
    margin: '1rem 0',
  },
  modalEngagement: {
    display: 'grid',
    gridTemplateColumns: 'repeat(3, 1fr)',
    gap: '1rem',
    margin: '1rem 0',
    padding: '1rem',
    backgroundColor: '#f8f9fa',
    borderRadius: '4px',
  },
  modalActions: {
    display: 'flex',
    gap: '0.5rem',
    marginTop: '1rem',
  },
  actionButton: {
    padding: '0.5rem 1rem',
    backgroundColor: '#007AFF',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
};