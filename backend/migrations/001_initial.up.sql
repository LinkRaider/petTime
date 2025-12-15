-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    auth_provider VARCHAR(50) DEFAULT 'email',
    auth_provider_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_auth_provider ON users(auth_provider, auth_provider_id);

-- Pet types (extensible)
CREATE TABLE pet_types (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    config JSONB DEFAULT '{}'
);

-- Insert default pet types
INSERT INTO pet_types (id, name, icon, config) VALUES
    ('dog', 'Dog', 'dog', '{"default_activities": ["walk", "fetch"]}'),
    ('cat', 'Cat', 'cat', '{"default_activities": ["walk", "play"]}');

-- Pets table
CREATE TABLE pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    pet_type_id VARCHAR(50) REFERENCES pet_types(id),
    name VARCHAR(100) NOT NULL,
    breed VARCHAR(100),
    avatar_url TEXT,
    birth_date DATE,
    total_xp INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    mood VARCHAR(50) DEFAULT 'happy',
    streak_days INTEGER DEFAULT 0,
    last_activity_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pets_user_id ON pets(user_id);

-- Game types registry
CREATE TABLE game_types (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    xp_config JSONB DEFAULT '{}',
    supported_pet_types VARCHAR(50)[] DEFAULT '{}',
    enabled BOOLEAN DEFAULT true
);

-- Insert default game types
INSERT INTO game_types (id, name, description, icon, xp_config, supported_pet_types) VALUES
    ('walk', 'Walk', 'Track your walks and earn XP', 'walking',
     '{"base_xp_per_minute": 2, "distance_bonus_per_km": 10, "streak_multiplier": 1.5}',
     ARRAY['dog', 'cat']),
    ('fetch', 'Fetch', 'Play fetch and track throws', 'ball',
     '{"xp_per_throw": 1, "combo_bonus": 5, "frenzy_multiplier": 2}',
     ARRAY['dog']);

-- Activities table
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID REFERENCES pets(id) ON DELETE CASCADE,
    game_type_id VARCHAR(50) REFERENCES game_types(id),
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    duration_seconds INTEGER,
    xp_earned INTEGER DEFAULT 0,
    game_data JSONB DEFAULT '{}',
    client_id UUID,
    synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_activities_pet_id ON activities(pet_id);
CREATE INDEX idx_activities_started_at ON activities(started_at);
CREATE INDEX idx_activities_client_id ON activities(client_id);

-- Achievements
CREATE TABLE achievements (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    category VARCHAR(50),
    criteria JSONB NOT NULL,
    xp_reward INTEGER DEFAULT 0
);

-- Insert some default achievements
INSERT INTO achievements (id, name, description, icon, category, criteria, xp_reward) VALUES
    ('first_walk', 'First Steps', 'Complete your first walk', 'footprints', 'milestone',
     '{"type": "activity_count", "game_type": "walk", "count": 1}', 50),
    ('streak_7', 'Week Warrior', 'Maintain a 7-day streak', 'fire', 'streak',
     '{"type": "streak_days", "days": 7}', 100),
    ('streak_30', 'Monthly Champion', 'Maintain a 30-day streak', 'trophy', 'streak',
     '{"type": "streak_days", "days": 30}', 500),
    ('distance_10k', 'Explorer', 'Walk a total of 10km', 'map', 'distance',
     '{"type": "total_distance", "meters": 10000}', 200);

-- User achievements
CREATE TABLE user_achievements (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id VARCHAR(100) REFERENCES achievements(id),
    pet_id UUID REFERENCES pets(id) ON DELETE CASCADE,
    unlocked_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, achievement_id, pet_id)
);

-- Cards
CREATE TABLE cards (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url TEXT,
    rarity VARCHAR(20) NOT NULL,
    category VARCHAR(50),
    drop_config JSONB DEFAULT '{}'
);

-- Insert some default cards
INSERT INTO cards (id, name, description, rarity, category, drop_config) VALUES
    ('sunny_walk', 'Sunny Day Walk', 'A beautiful sunny day for a walk', 'common', 'weather',
     '{"min_duration_minutes": 10}'),
    ('rainy_walk', 'Rainy Adventure', 'Walking in the rain', 'rare', 'weather',
     '{"min_duration_minutes": 15, "weather": "rain"}'),
    ('night_owl', 'Night Owl', 'A walk under the stars', 'rare', 'time',
     '{"time_range": {"start": 21, "end": 5}}'),
    ('marathon', 'Marathon Runner', 'Complete an extra long walk', 'epic', 'achievement',
     '{"min_distance_meters": 5000}'),
    ('first_friend', 'Best Friends', 'Your first activity together', 'legendary', 'milestone',
     '{"first_activity": true}');

-- User cards
CREATE TABLE user_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    card_id VARCHAR(100) REFERENCES cards(id),
    obtained_at TIMESTAMPTZ DEFAULT NOW(),
    activity_id UUID REFERENCES activities(id)
);

CREATE INDEX idx_user_cards_user_id ON user_cards(user_id);

-- Missions
CREATE TABLE missions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    mission_type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    target_value INTEGER NOT NULL,
    current_value INTEGER DEFAULT 0,
    xp_reward INTEGER NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_missions_user_id ON missions(user_id);
CREATE INDEX idx_missions_expires_at ON missions(expires_at);

-- Refresh tokens for JWT
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
