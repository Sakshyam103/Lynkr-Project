/**
 * Create Content Screen - Multipart form for uploading content with image
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ScrollView,
  Image,
} from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import { RouteProp } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { useSelector } from 'react-redux';
import { RootState } from '../store/store';
import { apiService } from '../services/api';
import * as FileSystem from 'expo-file-system';

async function base64ToFile(base64String: any) {
  const filename = `${FileSystem.cacheDirectory}upload.jpg`;
  await FileSystem.writeAsStringAsync(filename, base64String, { encoding: FileSystem.EncodingType.Base64 });
  return filename; // This will be a file://... URI
}


type RootStackParamList = {
  CreateContent: { eventId: string };
};

type CreateContentScreenRouteProp = RouteProp<RootStackParamList, 'CreateContent'>;
type CreateContentScreenNavigationProp = StackNavigationProp<RootStackParamList, 'CreateContent'>;

interface Props {
  route: CreateContentScreenRouteProp;
  navigation: CreateContentScreenNavigationProp;
}

export default function CreateContentScreen({ route, navigation }: Props) {
  const { eventId } = route.params;
  const { token } = useSelector((state: RootState) => state.auth);
  const [caption, setCaption] = useState('');
  const [tags, setTags] = useState('');
  const [allowBrandAccess, setAllowBrandAccess] = useState(true);
  const [allowCommercialUse, setAllowCommercialUse] = useState(false);
  const [allowModification, setAllowModification] = useState(true);
  const [allowSocialSharing, setAllowSocialSharing] = useState(true);
  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const pickImage = async () => {
    try {
      const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
      if (status !== 'granted') {
        Alert.alert('Permission needed', 'Camera roll permission is required');
        return;
      }

      const result = await ImagePicker.launchImageLibraryAsync({
        mediaTypes: ImagePicker.MediaTypeOptions.Images,
        allowsEditing: true,
        aspect: [4, 3],
        quality: 0.8,
      });

      if (!result.canceled && result.assets && result.assets[0]) {
        setSelectedImage(result.assets[0].uri);
      }
    } catch (error) {
      console.error('Image picker error:', error);
      Alert.alert('Error', 'Failed to open image picker. Please try again.');
    }
  };

  const submitContent = async () => {
    if (!selectedImage) {
      Alert.alert('Error', 'Please select an image');
      return;
    }

    if (!caption.trim()) {
      Alert.alert('Error', 'Please add a caption');
      return;
    }

    if (!token) {
      Alert.alert('Error', 'Authentication token not found. Please login again.');
      return;
    }

    setLoading(true);

    try {
      const formData = new FormData();
      formData.append('eventId', eventId);
      formData.append('caption', caption);
      
      // Format tags as JSON array
      const tagsArray = tags.split(',').map((tag, index) => ({
        id: `tag_${index}`,
        name: tag.trim(),
        type: 'category',
        brandId: '1',
        eventId: eventId.toString()
      }));
      formData.append('tags', JSON.stringify(tagsArray));
      
      // Format permissions as JSON object
      const permissions = {
        allowBrandAccess,
        allowCommercialUse,
        allowModification,
        allowSocialSharing,
        expirationDays: 30
      };
      formData.append('permissions', JSON.stringify(permissions));

      let uploadUri = selectedImage;
      if (uploadUri.startsWith('data:')) {
        const base64Data = uploadUri.split(',')[1];
        uploadUri = await base64ToFile(base64Data);
      }

      // Append image to formData
      formData.append('media', {
        uri: selectedImage,
        type: 'image/jpeg',
        name: 'content.jpg',
      } as any);
      console.log('Selected image URI:', selectedImage);

      await apiService.uploadContent(formData);
      
      // Alert.alert('Success', 'Content uploaded successfully', [
      //   { text: 'OK', onPress: () => navigation.goBack() }
      // ]);
      (navigation as any).navigate('ContentResult', {
        success: true,
        message: 'Your content has been uploaded successfully!',
        eventId: eventId
      });
    } catch (error) {
      console.error('Content upload error:', error);
      (navigation as any).navigate('ContentResult', {
        success: false,
        message: 'Failed to upload content. Please try again.',
        eventId: eventId
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.form}>
        <Text style={styles.label}>Select Image</Text>
        <TouchableOpacity style={styles.imageButton} onPress={pickImage}>
          {selectedImage ? (
            <Image source={{ uri: selectedImage }} style={styles.selectedImage} />
          ) : (
            <Text style={styles.imageButtonText}>📷 Choose Image</Text>
          )}
        </TouchableOpacity>

        <Text style={styles.label}>Caption</Text>
        <TextInput
          style={styles.textArea}
          value={caption}
          onChangeText={setCaption}
          placeholder="Write a caption..."
          multiline
          numberOfLines={4}
        />

        <Text style={styles.label}>Tags (comma separated)</Text>
        <TextInput
          style={styles.input}
          value={tags}
          onChangeText={setTags}
          placeholder="tag1, tag2, tag3"
        />

        <Text style={styles.label}>Permissions</Text>
        
        <View style={styles.permissionItem}>
          <TouchableOpacity
            style={[styles.checkbox, allowBrandAccess && styles.checkedBox]}
            onPress={() => setAllowBrandAccess(!allowBrandAccess)}
          >
            <Text style={styles.checkmark}>{allowBrandAccess ? '✓' : ''}</Text>
          </TouchableOpacity>
          <Text style={styles.permissionLabel}>Allow brand access</Text>
        </View>
        
        <View style={styles.permissionItem}>
          <TouchableOpacity
            style={[styles.checkbox, allowCommercialUse && styles.checkedBox]}
            onPress={() => setAllowCommercialUse(!allowCommercialUse)}
          >
            <Text style={styles.checkmark}>{allowCommercialUse ? '✓' : ''}</Text>
          </TouchableOpacity>
          <Text style={styles.permissionLabel}>Allow commercial use</Text>
        </View>
        
        <View style={styles.permissionItem}>
          <TouchableOpacity
            style={[styles.checkbox, allowModification && styles.checkedBox]}
            onPress={() => setAllowModification(!allowModification)}
          >
            <Text style={styles.checkmark}>{allowModification ? '✓' : ''}</Text>
          </TouchableOpacity>
          <Text style={styles.permissionLabel}>Allow modification</Text>
        </View>
        
        <View style={styles.permissionItem}>
          <TouchableOpacity
            style={[styles.checkbox, allowSocialSharing && styles.checkedBox]}
            onPress={() => setAllowSocialSharing(!allowSocialSharing)}
          >
            <Text style={styles.checkmark}>{allowSocialSharing ? '✓' : ''}</Text>
          </TouchableOpacity>
          <Text style={styles.permissionLabel}>Allow social sharing</Text>
        </View>

        <TouchableOpacity
          style={[styles.submitButton, loading && styles.disabledButton]}
          onPress={submitContent}
          disabled={loading}
        >
          <Text style={styles.submitButtonText}>
            {loading ? 'Uploading...' : 'Upload Content'}
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  form: {
    padding: 20,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 8,
    color: '#333',
  },
  input: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 12,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#ddd',
  },
  textArea: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 12,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#ddd',
    minHeight: 100,
    textAlignVertical: 'top',
  },
  imageButton: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 20,
    marginBottom: 16,
    borderWidth: 2,
    borderColor: '#ddd',
    borderStyle: 'dashed',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 150,
  },
  imageButtonText: {
    fontSize: 16,
    color: '#666',
  },
  selectedImage: {
    width: '100%',
    height: 150,
    borderRadius: 8,
  },
  permissionItem: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 12,
  },
  checkbox: {
    width: 24,
    height: 24,
    borderRadius: 4,
    borderWidth: 2,
    borderColor: '#007AFF',
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: 12,
    backgroundColor: '#fff',
  },
  checkedBox: {
    backgroundColor: '#007AFF',
  },
  checkmark: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  permissionLabel: {
    fontSize: 16,
    color: '#333',
  },
  submitButton: {
    backgroundColor: '#007AFF',
    borderRadius: 8,
    padding: 16,
    alignItems: 'center',
    marginTop: 10,
  },
  disabledButton: {
    backgroundColor: '#ccc',
  },
  submitButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});