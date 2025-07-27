import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, RefreshControl, TextInput } from 'react-native';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../store/store';
import { setEvents, setLoading } from '../store/slices/eventsSlice';
import { apiService } from '../services/api';

export default function EventsScreen({ navigation }: any) {
  const dispatch = useDispatch();
  const { events, loading } = useSelector((state: RootState) => state.events);
  const { user } = useSelector((state: RootState) => state.auth);
  const [searchQuery, setSearchQuery] = useState('');
  const [filteredEvents, setFilteredEvents] = useState<any[]>([]);

  const isUser = user?.role === 'user' || !user?.role;
  const isBrand = user?.role === 'brand';

  useEffect(() => {
    // Set user role in apiService
    if (user?.role) {
      apiService.setUserRole(user.role);
    }
    loadEvents();
  }, [user]);

  useEffect(() => {
    filterEvents();
  }, [events, searchQuery]);

  const loadEvents = async () => {
    dispatch(setLoading(true));
    try {
      console.log('Loading events for user role:', user?.role);
      console.log('API base URL will be:', user?.role === 'brand' ? 'brand/v1' : 'user/v1');
      const eventsData = await apiService.getEvents();
      console.log('Events loaded:', eventsData);
      dispatch(setEvents(eventsData));
    } catch (error) {
      console.error('Error loading events:', error);
      dispatch(setEvents([])); // Show empty instead of mock data
    } finally {
      dispatch(setLoading(false));
    }
  };



  const filterEvents = () => {
    if (!searchQuery) {
      setFilteredEvents(events);
    } else {
      const filtered = events.filter(event =>
        event.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        event.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
        event.location.address.toLowerCase().includes(searchQuery.toLowerCase())
      );
      setFilteredEvents(filtered);
    }
  };

  const handleCheckIn = async (eventId: string) => {
    try {
      await apiService.checkIn(eventId, { latitude: 37.7749, longitude: -122.4194 });
      // Show success feedback
    } catch (error) {
      console.log('Check-in error:', error);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'live': return '#28a745';
      case 'upcoming': return '#007AFF';
      case 'ended': return '#6c757d';
      default: return '#007AFF';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'live': return '🔴 LIVE';
      case 'upcoming': return '📅 Upcoming';
      case 'ended': return '✅ Ended';
      default: return '📅 Event';
    }
  };

  const renderEvent = ({ item }: any) => (
    <TouchableOpacity 
      style={styles.eventCard}
      onPress={() => navigation.navigate('EventDetails', { event: item })}
    >
      <View style={styles.eventHeader}>
        <View style={styles.eventTitleContainer}>
          <Text style={styles.eventName}>{item.name}</Text>
          <View style={[styles.statusBadge, { backgroundColor: getStatusColor(item.status) }]}>
            <Text style={styles.statusText}>{getStatusText(item.status)}</Text>
          </View>
        </View>
        <Text style={styles.eventDate}>
          {new Date(item.startDate).toLocaleDateString()}
        </Text>
      </View>
      
      <Text style={styles.eventDescription} numberOfLines={2}>{item.description}</Text>
      
      <View style={styles.eventMeta}>
        <Text style={styles.eventLocation}>📍 {item.location.address}</Text>
        <Text style={styles.attendeeCount}>👥 {item.attendeeCount} attending</Text>
      </View>
      
      <View style={styles.eventActions}>
        {/* USER ACTIONS - Can check in and view details */}
        {isUser && (
          <>
            {/*<TouchableOpacity*/}
            {/*  style={[styles.actionButton, styles.checkInButton]}*/}
            {/*  onPress={() => handleCheckIn(item.id)}*/}
            {/*>*/}
            {/*  <Text style={styles.checkInText}>Check In</Text>*/}
            {/*</TouchableOpacity>*/}
            
            <TouchableOpacity 
              style={[styles.actionButton, styles.detailsButton]}
              onPress={() => navigation.navigate('EventDetails', { event: item })}
            >
              <Text style={styles.detailsText}>View Details</Text>
            </TouchableOpacity>
          </>
        )}

        {/* BRAND ACTIONS - Can ONLY view analytics, NO check-in or content creation */}
        {/* {isBrand && (
          <>
            <TouchableOpacity 
              style={[styles.actionButton, styles.analyticsButton]}
              onPress={() => navigation.navigate('Analytics', { eventId: item.id })}
            >
              <Text style={styles.analyticsText}>View Analytics</Text>
            </TouchableOpacity>
            
            <TouchableOpacity 
              style={[styles.actionButton, styles.detailsButton]}
              onPress={() => navigation.navigate('EventDetails', { event: item, viewOnly: true })}
            >
              <Text style={styles.detailsText}>Event Details</Text>
            </TouchableOpacity>
          </>
        )} */}
        {/* BRAND ACTIONS - Can ONLY view analytics, NO check-in or content creation */}
{isBrand && (
  <>
    <TouchableOpacity 
      style={[styles.actionButton, styles.analyticsButton]}
      onPress={() => navigation.navigate('Analytics', { eventId: item.id })}
    >
      <Text style={styles.analyticsButtonText}>📊 View Analytics</Text>
    </TouchableOpacity>
    
    <TouchableOpacity 
      style={[styles.actionButton, styles.surveyButton]}
      onPress={() => navigation.navigate('CreateSurvey', { eventId: item.id, eventName: item.name })}
    >
      <Text style={styles.surveyButtonText}>📝 Create Survey</Text>
    </TouchableOpacity>
    
    <TouchableOpacity 
      style={[styles.actionButton, styles.detailsButton]}
      onPress={() => navigation.navigate('EventDetails', { event: item, viewOnly: true })}
    >
      <Text style={styles.detailsText}>Event Details</Text>
    </TouchableOpacity>
  </>
)}

      </View>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>
          {isBrand ? 'Sponsored Events Analytics' : 'Discover Events'}
        </Text>
        <Text style={styles.subtitle}>
          {isBrand ? 'View analytics for your sponsored events' : 'Find sponsored events near you'}
        </Text>
        
        <View style={styles.searchContainer}>
          <TextInput
            style={styles.searchInput}
            placeholder="Search events, locations..."
            value={searchQuery}
            onChangeText={setSearchQuery}
            placeholderTextColor="#999"
          />
          <Text style={styles.searchIcon}>🔍</Text>
        </View>

        {/* BRAND ONLY - Create Event Button (if brands can create events) */}
        {isBrand && (
          <TouchableOpacity 
            style={styles.createEventButton}
            onPress={() => navigation.navigate('CreateEvent')}
          >
            <Text style={styles.createEventText}>+ Create New Event</Text>
          </TouchableOpacity>
        )}
      </View>

      <FlatList
        data={events}
        renderItem={renderEvent}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContainer}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl refreshing={loading} onRefresh={loadEvents} />
        }
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No events found</Text>
            <Text style={styles.emptySubtext}>
              {isBrand ? 'No sponsored events available for analytics' : 'Try adjusting your search or check back later'}
            </Text>
          </View>
        }
      />
    </View>
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
    borderBottomLeftRadius: 20,
    borderBottomRightRadius: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#333',
  },
  subtitle: {
    fontSize: 16,
    color: '#666',
    marginTop: 4,
    marginBottom: 20,
  },
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f8f9fa',
    borderRadius: 12,
    paddingHorizontal: 16,
    marginBottom: 16,
  },
  searchInput: {
    flex: 1,
    paddingVertical: 12,
    fontSize: 16,
    color: '#333',
  },
  searchIcon: {
    fontSize: 16,
    marginLeft: 8,
  },
  createEventButton: {
    backgroundColor: '#28a745',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  createEventText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  listContainer: {
    padding: 20,
  },
  eventCard: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 20,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 12,
    elevation: 4,
  },
  eventHeader: {
    marginBottom: 12,
  },
  eventTitleContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 8,
  },
  eventName: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    flex: 1,
    marginRight: 12,
  },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 12,
  },
  statusText: {
    color: '#fff',
    fontSize: 10,
    fontWeight: 'bold',
  },
  eventDate: {
    fontSize: 14,
    color: '#007AFF',
    fontWeight: '600',
  },
  eventDescription: {
    fontSize: 14,
    color: '#666',
    lineHeight: 20,
    marginBottom: 12,
  },
  eventMeta: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 16,
  },
  eventLocation: {
    fontSize: 14,
    color: '#666',
    flex: 1,
  },
  attendeeCount: {
    fontSize: 14,
    color: '#666',
    fontWeight: '500',
  },
  eventActions: {
    flexDirection: 'row',
    gap: 12,
  },
  actionButton: {
    flex: 1,
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderRadius: 12,
    alignItems: 'center',
  },
  checkInButton: {
    backgroundColor: '#007AFF',
  },
  checkInText: {
    color: '#fff',
    fontWeight: '600',
    fontSize: 14,
  },
  analyticsButton: {
    backgroundColor: '#28a745',
  },
  analyticsText: {
    color: '#fff',
    fontWeight: '600',
    fontSize: 14,
  },
  detailsButton: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: '#007AFF',
  },
  detailsText: {
    color: '#007AFF',
    fontWeight: '600',
    fontSize: 14,
  },
  emptyContainer: {
    alignItems: 'center',
    paddingTop: 60,
  },
  emptyText: {
    fontSize: 18,
    fontWeight: '600',
    color: '#666',
    marginBottom: 8,
  },
  emptySubtext: {
    fontSize: 14,
    color: '#999',
    textAlign: 'center',
  },
  surveyButton: {
    backgroundColor: '#6f42c1',
  },
  surveyButtonText: {
    color: '#fff',
    fontWeight: '600',
    fontSize: 14,
  },
});