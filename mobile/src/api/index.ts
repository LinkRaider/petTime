import apiClient from './client';
import {
  User,
  AuthTokens,
  LoginRequest,
  RegisterRequest,
  Pet,
  PetType,
  CreatePetRequest,
  Activity,
  CreateActivityRequest,
  GameType,
  PetStats,
} from '../types';

export const authApi = {
  register: async (data: RegisterRequest) => {
    const response = await apiClient.post<{ user: User; tokens: AuthTokens }>(
      '/auth/register',
      data
    );
    return response.data;
  },

  login: async (data: LoginRequest) => {
    const response = await apiClient.post<{ user: User; tokens: AuthTokens }>(
      '/auth/login',
      data
    );
    return response.data;
  },

  logout: async (refreshToken: string) => {
    await apiClient.post('/auth/logout', { refresh_token: refreshToken });
  },
};

export const petsApi = {
  list: async () => {
    const response = await apiClient.get<Pet[]>('/pets');
    return response.data;
  },

  get: async (id: string) => {
    const response = await apiClient.get<Pet>(`/pets/${id}`);
    return response.data;
  },

  create: async (data: CreatePetRequest) => {
    const response = await apiClient.post<Pet>('/pets', data);
    return response.data;
  },

  update: async (id: string, data: Partial<CreatePetRequest>) => {
    const response = await apiClient.put<Pet>(`/pets/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    await apiClient.delete(`/pets/${id}`);
  },

  getStats: async (id: string) => {
    const response = await apiClient.get<PetStats>(`/pets/${id}/stats`);
    return response.data;
  },

  getTypes: async () => {
    const response = await apiClient.get<PetType[]>('/pet-types');
    return response.data;
  },
};

// Alias for backwards compatibility
export const petApi = petsApi;

export const activityApi = {
  getAll: async (petId?: string) => {
    const params = petId ? { pet_id: petId } : {};
    const response = await apiClient.get<Activity[]>('/activities', { params });
    return response.data;
  },

  getById: async (id: string) => {
    const response = await apiClient.get<Activity>(`/activities/${id}`);
    return response.data;
  },

  create: async (data: CreateActivityRequest) => {
    const response = await apiClient.post<Activity>('/activities', data);
    return response.data;
  },

  update: async (id: string, data: Partial<CreateActivityRequest>) => {
    const response = await apiClient.put<Activity>(`/activities/${id}`, data);
    return response.data;
  },

  sync: async (activities: CreateActivityRequest[]) => {
    const response = await apiClient.post<Activity[]>('/activities/sync', {
      activities,
    });
    return response.data;
  },

  getGameTypes: async () => {
    const response = await apiClient.get<GameType[]>('/game-types');
    return response.data;
  },
};
