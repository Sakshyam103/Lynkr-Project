/**
 * Lynkr Mobile App - Role-based Navigation with Endpoint Testing
 */

import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { createStackNavigator } from '@react-navigation/stack';
import { Provider, useSelector } from 'react-redux';
import { StatusBar } from 'expo-status-bar';
import { store, RootState } from './src/store/store';
import  ContentResultScreen from './src/screens/CreateContentScreen';



// Auth Screens
import RoleSelectionScreen from './src/screens/RoleSelectionScreen';
import LoginScreen from './src/screens/LoginScreen';
import RegisterScreen from './src/screens/RegisterScreen';

// User Screens
import HomeScreen from './src/screens/HomeScreen';
import EventsScreen from './src/screens/EventsScreen';
import EventDetailsScreen from './src/screens/EventDetailsScreen';
import ContentScreen from './src/screens/ContentScreen';
import ProfileScreen from './src/screens/ProfileScreen';
import SurveyScreen from './src/screens/SurveyScreen';
import UserRewardsScreen from './src/screens/UserRewardsScreen';

// Brand Screens
import BrandDashboardScreen from './src/screens/BrandDashboardScreen';
import CreateEventScreen from './src/screens/CreateEventScreen';
import BrandAnalyticsScreen from './src/screens/BrandAnalyticsScreen';
import DiscountCodesScreen from './src/screens/DiscountCodesScreen';
import BrandSurveyScreen from './src/screens/BrandSurveyScreen';
import CreateContentScreen from './src/screens/CreateContentScreen';
import CreateSurveyScreen from './src/screens/CreateSurveyScreen';



// Testing Screen
import EndpointTestScreen from './src/screens/EndpointTestScreen';

const Tab = createBottomTabNavigator();
const Stack = createStackNavigator();

// Event Stack (shared by all roles)
function EventStack() {
  return (
    <Stack.Navigator>
      <Stack.Screen name="EventsList" component={EventsScreen} options={{ title: 'Events' }} />
      <Stack.Screen name="EventDetails" component={EventDetailsScreen} options={{ title: 'Event Details' }} />
      <Stack.Screen name="CreateContent" component={CreateContentScreen} options={{ title: 'Create Content' }} />
        <Stack.Screen name="ContentResult" component={ContentResultScreen} options={{ title: 'Upload Result' }} />
    </Stack.Navigator>
  );
}

// Brand Event Stack (includes create event)
function BrandEventStack() {
  return (
    <Stack.Navigator>
      <Stack.Screen name="EventsList" component={EventsScreen} options={{ title: 'My Events' }} />
      <Stack.Screen name="EventDetails" component={EventDetailsScreen} options={{ title: 'Event Details' }} />
      {/*<Stack.Screen name="CreateEvent" component={CreateEventScreen} options={{ title: 'Create Event' }} />*/}
      <Stack.Screen name="Analytics" component={BrandAnalyticsScreen} options={{ title: 'Analytics' }} />
      <Stack.Screen name="CreateSurvey" component={CreateSurveyScreen} options={{ title: 'Create Survey' }} />

    </Stack.Navigator>
  );
}

// Auth Stack
function AuthStack() {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false }}>
      <Stack.Screen name="RoleSelection" component={RoleSelectionScreen} />
      <Stack.Screen name="Login" component={LoginScreen} />
      <Stack.Screen name="Register" component={RegisterScreen} />
    </Stack.Navigator>
  );
}

// USER NAVIGATION - Event attendees
function UserTabs() {
  return (
    <Tab.Navigator
      screenOptions={{
        tabBarActiveTintColor: '#007AFF',
        tabBarInactiveTintColor: '#666',
      }}
    >
      <Tab.Screen 
        name="Home" 
        component={HomeScreen}
        options={{
          tabBarLabel: 'Home',
          tabBarIcon: () => '🏠',
        }}
      />
      <Tab.Screen 
        name="Events" 
        component={EventStack}
        options={{
          tabBarLabel: 'Events',
          tabBarIcon: () => '📅',
          headerShown: false,
        }}
      />
      <Tab.Screen
        name="Content"
        component={ContentScreen}
        options={{
          tabBarLabel: 'Content',
          tabBarIcon: () => '📸',
        }}
      />
      <Tab.Screen 
        name="Surveys" 
        component={SurveyScreen}
        options={{
          tabBarLabel: 'Surveys',
          tabBarIcon: () => '📝',
        }}
      />
      <Tab.Screen 
        name="Rewards" 
        component={UserRewardsScreen}
        options={{
          tabBarLabel: 'Rewards',
          tabBarIcon: () => '🏆',
        }}
      />
      <Tab.Screen 
        name="Profile" 
        component={ProfileScreen}
        options={{
          tabBarLabel: 'Profile',
          tabBarIcon: () => '👤',
        }}
      />
      <Tab.Screen 
        name="Test" 
        component={EndpointTestScreen}
        options={{
          tabBarLabel: 'Test',
          tabBarIcon: () => '🧪',
        }}
      />
    </Tab.Navigator>
  );
}

// BRAND NAVIGATION - Event sponsors
function BrandTabs() {
  return (
    <Tab.Navigator
      screenOptions={{
        tabBarActiveTintColor: '#28a745',
        tabBarInactiveTintColor: '#666',
      }}
    >
      <Tab.Screen 
        name="Dashboard" 
        component={BrandDashboardScreen}
        options={{
          tabBarLabel: 'Dashboard',
          tabBarIcon: () => '📊',
        }}
      />
      <Tab.Screen 
        name="Events" 
        component={BrandEventStack}
        options={{
          tabBarLabel: 'Events',
          tabBarIcon: () => '📅',
          headerShown: false,
        }}
      />
      <Tab.Screen
        name="Content"
        component={ContentScreen}
        options={{
          tabBarLabel: 'Content',
          tabBarIcon: () => '📸',
        }}
      />
      <Tab.Screen 
        name="Analytics" 
        component={BrandAnalyticsScreen}
        options={{
          tabBarLabel: 'Analytics',
          tabBarIcon: () => '📈',
        }}
      />
      <Tab.Screen 
  name="Surveys" 
  component={BrandSurveyScreen}
  options={{
    tabBarLabel: 'Surveys',
    tabBarIcon: () => '📝',
  }}
/>

      <Tab.Screen 
        name="Discounts" 
        component={DiscountCodesScreen}
        options={{
          tabBarLabel: 'Discounts',
          tabBarIcon: () => '🎫',
        }}
      />


      <Tab.Screen 
        name="Profile" 
        component={ProfileScreen}
        options={{
          tabBarLabel: 'Profile',
          tabBarIcon: () => '👤',
        }}
      />
      <Tab.Screen 
        name="Test" 
        component={EndpointTestScreen}
        options={{
          tabBarLabel: 'Test',
          tabBarIcon: () => '🧪',
        }}
      />
    </Tab.Navigator>
  );
}

function AppContent() {
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);

  const getMainComponent = () => {
    if (!isAuthenticated) return <AuthStack />;
    
    // Route based on user role
    switch (user?.role) {
      case 'brand':
        return <BrandTabs />;
      // case 'organization':
      //   return <BrandTabs />; // Organizations get same features as brands for now
      default:
        return <UserTabs />;
    }
  };

  return (
    <NavigationContainer>
      <StatusBar style="dark" />
      {getMainComponent()}
    </NavigationContainer>
  );
}

export default function App() {
  return (
    <Provider store={store}>
      <AppContent />
    </Provider>
  );
}