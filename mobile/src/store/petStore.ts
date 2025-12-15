import { create } from 'zustand';
import { Pet, PetType, CreatePetRequest, PetStats } from '../types';
import { petApi } from '../api';

interface PetState {
  pets: Pet[];
  selectedPet: Pet | null;
  petTypes: PetType[];
  stats: PetStats | null;
  isLoading: boolean;
  error: string | null;

  fetchPets: () => Promise<void>;
  fetchPetTypes: () => Promise<void>;
  fetchPetStats: (petId: string) => Promise<void>;
  selectPet: (pet: Pet | null) => void;
  createPet: (data: CreatePetRequest) => Promise<Pet>;
  updatePet: (id: string, data: Partial<CreatePetRequest>) => Promise<void>;
  deletePet: (id: string) => Promise<void>;
  clearError: () => void;
}

export const usePetStore = create<PetState>((set, get) => ({
  pets: [],
  selectedPet: null,
  petTypes: [],
  stats: null,
  isLoading: false,
  error: null,

  fetchPets: async () => {
    try {
      set({ isLoading: true, error: null });
      const pets = await petApi.getAll();
      set({ pets, isLoading: false });
    } catch (error: any) {
      set({
        error: error.response?.data?.message || 'Failed to fetch pets',
        isLoading: false,
      });
    }
  },

  fetchPetTypes: async () => {
    try {
      const petTypes = await petApi.getPetTypes();
      set({ petTypes });
    } catch (error: any) {
      set({ error: error.response?.data?.message || 'Failed to fetch pet types' });
    }
  },

  fetchPetStats: async (petId: string) => {
    try {
      set({ isLoading: true, error: null });
      const data = await petApi.getStats(petId);
      set({ stats: data.stats, isLoading: false });
    } catch (error: any) {
      set({
        error: error.response?.data?.message || 'Failed to fetch stats',
        isLoading: false,
      });
    }
  },

  selectPet: (pet: Pet | null) => {
    set({ selectedPet: pet, stats: null });
  },

  createPet: async (data: CreatePetRequest) => {
    try {
      set({ isLoading: true, error: null });
      const pet = await petApi.create(data);
      set((state) => ({
        pets: [...state.pets, pet],
        selectedPet: pet,
        isLoading: false,
      }));
      return pet;
    } catch (error: any) {
      set({
        error: error.response?.data?.message || 'Failed to create pet',
        isLoading: false,
      });
      throw error;
    }
  },

  updatePet: async (id: string, data: Partial<CreatePetRequest>) => {
    try {
      set({ isLoading: true, error: null });
      const updatedPet = await petApi.update(id, data);
      set((state) => ({
        pets: state.pets.map((pet) => (pet.id === id ? updatedPet : pet)),
        selectedPet:
          state.selectedPet?.id === id ? updatedPet : state.selectedPet,
        isLoading: false,
      }));
    } catch (error: any) {
      set({
        error: error.response?.data?.message || 'Failed to update pet',
        isLoading: false,
      });
      throw error;
    }
  },

  deletePet: async (id: string) => {
    try {
      set({ isLoading: true, error: null });
      await petApi.delete(id);
      set((state) => ({
        pets: state.pets.filter((pet) => pet.id !== id),
        selectedPet: state.selectedPet?.id === id ? null : state.selectedPet,
        isLoading: false,
      }));
    } catch (error: any) {
      set({
        error: error.response?.data?.message || 'Failed to delete pet',
        isLoading: false,
      });
      throw error;
    }
  },

  clearError: () => set({ error: null }),
}));
