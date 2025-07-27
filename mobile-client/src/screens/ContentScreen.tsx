import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, Alert, TextInput } from 'react-native';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';

export default function ContentScreen({ navigation }: any) {
  const { user } = useSelector((state: RootState) => state.auth);
  const [caption, setCaption] = useState('');
  const [tags, setTags] = useState('');
  const [selectedImage, setSelectedImage] = useState<string | null>(null);

  const isBrand = user?.role === 'brand';
  const isUser = user?.role === 'user' || !user?.role;

  const pickImage = async () => {
    Alert.alert('Image Picker', 'Image picker functionality will be available in the native app');
    setSelectedImage('mock-image-selected');
  };

  const takePhoto = async () => {
    Alert.alert('Camera', 'Camera functionality will be available in the native app');
    setSelectedImage('mock-photo-taken');
  };

  const handlePost = () => {
    if (!selectedImage) {
      Alert.alert('Error', 'Please select an image first');
      return;
    }

    Alert.alert('Success', 'Content posted successfully!', [
      { text: 'OK', onPress: () => {
        setSelectedImage(null);
        setCaption('');
        setTags('');
      }}
    ]);
  };

  const mockContent = [
    {
      id: '1',
      user: 'Sarah M.',
      caption: 'Amazing keynote at Tech Conference! #innovation #tech',
      time: '2 hours ago',
      likes: 24,
      event: 'Tech Conference 2024'
    },
    {
      id: '2',
      user: 'Mike R.',
      caption: 'Great networking opportunities here! #networking #business',
      time: '4 hours ago',
      likes: 18,
      event: 'Brand Expo 2024'
    }
  ];

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>
          {isBrand ? 'User Content' : 'Create Content'}
        </Text>
        <Text style={styles.subtitle}>
          {isBrand ? 'Content created by event attendees' : 'Share your event experience'}
        </Text>
      </View>

      {/* USER ONLY - Content Creation Section */}
      {isUser && (
        <View style={styles.createSection}>
          <View style={styles.mediaButtons}>
            <TouchableOpacity style={styles.mediaButton} onPress={takePhoto}>
              <Text style={styles.mediaButtonText}>📷 Take Photo</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.mediaButton} onPress={pickImage}>
              <Text style={styles.mediaButtonText}>🖼️ Choose Image</Text>
            </TouchableOpacity>
          </View>

          {selectedImage && (
            <View style={styles.selectedImageContainer}>
              <Text style={styles.selectedImageText}>✅ Image selected</Text>
            </View>
          )}

          <TextInput
            style={styles.captionInput}
            placeholder="Write a caption..."
            value={caption}
            onChangeText={setCaption}
            multiline
            numberOfLines={3}
          />

          <TextInput
            style={styles.tagsInput}
            placeholder="Add tags (e.g., #tech #innovation)"
            value={tags}
            onChangeText={setTags}
          />

          <View style={styles.permissionsSection}>
            <Text style={styles.permissionsTitle}>Content Permissions</Text>
            <TouchableOpacity style={styles.permissionOption}>
              <Text style={styles.permissionText}>✅ Allow brands to use this content</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.permissionOption}>
              <Text style={styles.permissionText}>✅ Show in public feed</Text>
            </TouchableOpacity>
          </View>

          <TouchableOpacity style={styles.postButton} onPress={handlePost}>
            <Text style={styles.postButtonText}>Post Content</Text>
          </TouchableOpacity>
        </View>
      )}

      {/* Content Feed - Visible to both users and brands */}
      <View style={styles.feedSection}>
        <Text style={styles.feedTitle}>
          {isBrand ? 'User Generated Content' : 'Recent Content'}
        </Text>
        
        {mockContent.map((item) => (
          <View key={item.id} style={styles.contentCard}>
            <View style={styles.contentHeader}>
              <Text style={styles.userName}>{item.user}</Text>
              <Text style={styles.eventName}>{item.event}</Text>
            </View>
            <Text style={styles.contentCaption}>{item.caption}</Text>
            <View style={styles.contentFooter}>
              <Text style={styles.contentTime}>{item.time}</Text>
              <Text style={styles.contentLikes}>❤️ {item.likes}</Text>
            </View>
            
            {/* BRAND ONLY - Additional Actions */}
            {isBrand && (
              <View style={styles.brandActions}>
                <TouchableOpacity style={styles.brandActionButton}>
                  <Text style={styles.brandActionText}>Request Usage Rights</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.brandActionButton}>
                  <Text style={styles.brandActionText}>Analyze Sentiment</Text>
                </TouchableOpacity>
              </View>
            )}
          </View>
        ))}
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#333',
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginTop: 4,
  },
  createSection: {
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
  mediaButtons: {
    flexDirection: 'row',
    gap: 10,
    marginBottom: 15,
  },
  mediaButton: {
    flex: 1,
    backgroundColor: '#f8f9fa',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  mediaButtonText: {
    fontSize: 14,
    color: '#333',
  },
  selectedImageContainer: {
    backgroundColor: '#e8f5e8',
    padding: 10,
    borderRadius: 8,
    marginBottom: 15,
  },
  selectedImageText: {
    color: '#28a745',
    textAlign: 'center',
  },
  captionInput: {
    borderWidth: 1,
    borderColor: '#e0e0e0',
    borderRadius: 8,
    padding: 15,
    marginBottom: 15,
    fontSize: 16,
    textAlignVertical: 'top',
  },
  tagsInput: {
    borderWidth: 1,
    borderColor: '#e0e0e0',
    borderRadius: 8,
    padding: 15,
    marginBottom: 15,
    fontSize: 16,
  },
  permissionsSection: {
    marginBottom: 20,
  },
  permissionsTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 10,
  },
  permissionOption: {
    paddingVertical: 8,
  },
  permissionText: {
    fontSize: 14,
    color: '#333',
  },
  postButton: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
  },
  postButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  feedSection: {
    padding: 20,
  },
  feedTitle: {
    fontSize: 20,
    fontWeight: '600',
    color: '#333',
    marginBottom: 15,
  },
  contentCard: {
    backgroundColor: '#fff',
    padding: 15,
    borderRadius: 8,
    marginBottom: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  contentHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 8,
  },
  userName: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
  },
  eventName: {
    fontSize: 12,
    color: '#007AFF',
  },
  contentCaption: {
    fontSize: 14,
    color: '#333',
    marginBottom: 8,
    lineHeight: 20,
  },
  contentFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  contentTime: {
    fontSize: 12,
    color: '#666',
  },
  contentLikes: {
    fontSize: 12,
    color: '#666',
  },
  brandActions: {
    flexDirection: 'row',
    marginTop: 12,
    gap: 8,
  },
  brandActionButton: {
    flex: 1,
    backgroundColor: '#28a745',
    padding: 8,
    borderRadius: 6,
    alignItems: 'center',
  },
  brandActionText: {
    fontSize: 12,
    color: '#fff',
    fontWeight: '600',
  },
});
