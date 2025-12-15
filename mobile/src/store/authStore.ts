import { create } from 'zustand';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { User, AuthTokens, LoginRequest, RegisterRequest } from '../types';
import { authApi } from '../api';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  loadUser: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,

  login: async (data: LoginRequest) => {
    try {
      set({ isLoading: true, error: null });
      const response = await authApi.login(data);

      await AsyncStorage.setItem('user', JSON.stringify(response.user));
      await AsyncStorage.setItem('access_token', response.tokens.access_token);
      await AsyncStorage.setItem('refresh_token', response.tokens.refresh_token);

      set({ user: response.user, isAuthenticated: true, isLoading: false });
    } catch (error: any) {
      set({
        error: error.response?.data?.message || 'Login failed',
        isLoading: false,
      });
      throw error;
    }
  },

  register: async (data: RegisterRequest) => {
    try {
      set({ isLoading: true, error: null });
      const response = await authApi.register(data);

      await AsyncStorage.setItem('user', JSON.stringify(response.user));
      await AsyncStorage.setItem('access_token', response.tokens.access_token);
      await AsyncStorage.setItem('refresh_token', response.tokens.refresh_token);

      set({ user: response.user, isAuthenticated: true, isLoading: false });
    } catch (error: any) {
      set({
        error: error.response?.data?.message || 'Registration failed',
        isLoading: false,
      });
      throw error;
    }
  },

  logout: async () => {
    try {
      const refreshToken = await AsyncStorage.getItem('refresh_token');
      if (refreshToken) {
        await authApi.logout(refreshToken);
      }
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      await AsyncStorage.multiRemove(['user', 'access_token', 'refresh_token']);
      set({ user: null, isAuthenticated: false });
    }
  },

  loadUser: async () => {
    try {
      const userStr = await AsyncStorage.getItem('user');
      const token = await AsyncStorage.getItem('access_token');

      if (userStr && token) {
        const user = JSON.parse(userStr);
        set({ user, isAuthenticated: true });
      } else {
        set({ isAuthenticated: false });
      }
    } catch (error) {
      console.error('Load user error:', error);
      set({ isAuthenticated: false });
    }
  },

  clearError: () => set({ error: null }),
}));
