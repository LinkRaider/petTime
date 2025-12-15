export interface User {
  id: string;
  email: string;
  name: string;
  avatar_url?: string;
  auth_provider: string;
  created_at: string;
  updated_at: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface PetType {
  id: string;
  name: string;
  icon?: string;
  config?: any;
}

export interface Pet {
  id: string;
  user_id: string;
  pet_type_id: string;
  pet_type?: PetType;
  name: string;
  breed?: string;
  avatar_url?: string;
  birth_date?: string;
  total_xp: number;
  level: number;
  mood: 'happy' | 'content' | 'tired' | 'sad' | 'bored';
  streak_days: number;
  last_activity_at?: string;
  created_at: string;
  updated_at: string;
}

export interface GameType {
  id: string;
  name: string;
  description?: string;
  icon?: string;
  xp_config?: any;
  supported_pet_types: string[];
  enabled: boolean;
}

export interface Activity {
  id: string;
  pet_id: string;
  game_type_id: string;
  game_type?: GameType;
  started_at: string;
  ended_at?: string;
  duration_seconds?: number;
  xp_earned: number;
  game_data?: any;
  client_id?: string;
  synced_at?: string;
  created_at: string;
}

export interface PetStats {
  total_activities: number;
  total_duration_seconds: number;
  total_distance_meters: number;
  current_streak: number;
  longest_streak: number;
  xp_to_next_level: number;
  level_progress: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface CreatePetRequest {
  pet_type_id: string;
  name: string;
  breed?: string;
  avatar_url?: string;
  birth_date?: string;
}

export interface CreateActivityRequest {
  pet_id: string;
  game_type_id: string;
  started_at: string;
  ended_at?: string;
  game_data?: any;
  client_id?: string;
}
