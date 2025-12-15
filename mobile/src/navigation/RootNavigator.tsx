import React, { useEffect, useState } from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { useAuthStore } from '../store/authStore';

// Screens
import LoginScreen from '../screens/auth/LoginScreen';
import RegisterScreen from '../screens/auth/RegisterScreen';
import HomeScreen from '../screens/HomeScreen';
import PetListScreen from '../screens/pets/PetListScreen';
import PetDetailScreen from '../screens/pets/PetDetailScreen';
import CreatePetScreen from '../screens/pets/CreatePetScreen';

export type RootStackParamList = {
  Login: undefined;
  Register: undefined;
  Home: undefined;
  PetList: undefined;
  PetDetail: { petId: string };
  CreatePet: undefined;
};

const Stack = createNativeStackNavigator<RootStackParamList>();

export default function RootNavigator() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const loadUser = useAuthStore((state) => state.loadUser);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    loadUser().finally(() => setIsReady(true));
  }, [loadUser]);

  if (!isReady) {
    return null;
  }

  // Ensure boolean type
  const authenticated = Boolean(isAuthenticated);

  return (
    <NavigationContainer>
      <Stack.Navigator
        initialRouteName={authenticated ? 'Home' : 'Login'}
        screenOptions={{
          headerStyle: {
            backgroundColor: '#4F46E5',
          },
          headerTintColor: '#fff',
          headerTitleStyle: {
            fontWeight: 'bold',
          },
        }}
      >
        <Stack.Screen
          name="Login"
          component={LoginScreen}
          options={{ headerShown: false }}
        />
        <Stack.Screen
          name="Register"
          component={RegisterScreen}
          options={{ title: 'Sign Up' }}
        />
        <Stack.Screen
          name="Home"
          component={HomeScreen}
          options={{ title: 'PetTime' }}
        />
        <Stack.Screen
          name="PetList"
          component={PetListScreen}
          options={{ title: 'My Pets' }}
        />
        <Stack.Screen
          name="PetDetail"
          component={PetDetailScreen}
          options={{ title: 'Pet Details' }}
        />
        <Stack.Screen
          name="CreatePet"
          component={CreatePetScreen}
          options={{ title: 'Add Pet' }}
        />
      </Stack.Navigator>
    </NavigationContainer>
  );
}
