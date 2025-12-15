# PetTime - Gamified Pet Activities App

A mobile application that turns pet activities (walks, play sessions) into an engaging RPG-style experience with XP progression, streaks, collectibles, and mood tracking. Built with Go backend and React Native frontend.

## Features

- **RPG Progression System**: Pets earn XP and level up through activities
- **Multiple Game Types**: Walk tracking, Fetch mini-game (coming soon)
- **Pet Mood System**: Dynamic mood changes based on activity frequency
- **Streak Tracking**: Consecutive day activity tracking
- **Multiple Pet Support**: Manage multiple pets (dogs, cats)
- **Offline-First**: Activities tracked locally, synced when online
- **Cross-Platform**: iOS and Android support via Expo

## Tech Stack

### Backend
- **Go 1.23** with clean architecture
- **PostgreSQL 16** with JSONB for flexible game data
- **Chi Router** for HTTP routing
- **JWT Authentication** with refresh tokens
- **Docker** for containerization

### Mobile
- **React Native** via Expo
- **TypeScript** for type safety
- **Zustand** for state management
- **React Navigation** for routing
- **AsyncStorage** for local persistence
- **Axios** for API communication

## Project Structure

```
petTime/
├── backend/
│   ├── cmd/api/main.go           # Application entry point
│   ├── internal/
│   │   ├── config/               # Configuration
│   │   ├── handlers/             # HTTP handlers
│   │   ├── services/             # Business logic
│   │   ├── repositories/         # Data access
│   │   ├── models/               # Domain models
│   │   └── middleware/           # Auth, CORS, etc.
│   ├── migrations/               # Database migrations
│   ├── Dockerfile
│   └── go.mod
├── mobile/
│   ├── src/
│   │   ├── api/                  # API client
│   │   ├── screens/              # Screen components
│   │   ├── store/                # Zustand stores
│   │   ├── navigation/           # Navigation setup
│   │   └── types/                # TypeScript types
│   ├── App.tsx
│   └── package.json
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ and npm
- Expo CLI: `npm install -g expo-cli`
- iOS Simulator (for Mac) or Android Studio (for Android development)

### Backend Setup

1. **Start the backend services**:
```bash
cd backend
docker-compose up -d
```

This will start:
- PostgreSQL on port 5432
- Go API on port 8080

2. **Verify backend is running**:
```bash
curl http://localhost:8080/health
# Should return: {"status":"healthy","timestamp":"..."}
```

3. **View logs**:
```bash
docker-compose logs -f api
```

### Mobile Setup

1. **Install dependencies**:
```bash
cd mobile
npm install
```

2. **Configure API URL** (important for physical devices):

Edit `mobile/src/api/client.ts`:

```typescript
// For iOS Simulator or Android Emulator:
const API_URL = 'http://localhost:8080/api/v1';

// For physical device on same network:
const API_URL = 'http://YOUR_LOCAL_IP:8080/api/v1';
// Example: 'http://192.168.1.100:8080/api/v1'
```

To find your local IP:
- **Mac/Linux**: `ifconfig | grep inet`
- **Windows**: `ipconfig`

3. **Start the app**:
```bash
npx expo start
```

4. **Run on device**:
- Press `i` for iOS Simulator
- Press `a` for Android Emulator
- Scan QR code with Expo Go app for physical device

## API Documentation

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}

Response: 201 Created
{
  "user": { "id": "...", "email": "...", "name": "..." },
  "tokens": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response: 200 OK
{
  "user": { "id": "...", "email": "...", "name": "..." },
  "tokens": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_at": "2024-01-01T00:00:00Z"
  }
}
```

### Pets

#### Create Pet
```http
POST /api/v1/pets
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "pet_type_id": "dog",
  "name": "Max",
  "breed": "Golden Retriever"
}

Response: 201 Created
{
  "id": "...",
  "name": "Max",
  "pet_type_id": "dog",
  "level": 1,
  "total_xp": 0,
  "mood": "happy",
  "streak_days": 0
}
```

#### List Pets
```http
GET /api/v1/pets
Authorization: Bearer {access_token}

Response: 200 OK
[
  {
    "id": "...",
    "name": "Max",
    "pet_type": { "id": "dog", "name": "Dog" },
    "level": 5,
    "total_xp": 1250,
    "mood": "happy",
    "streak_days": 7
  }
]
```

#### Get Pet Stats
```http
GET /api/v1/pets/{pet_id}/stats
Authorization: Bearer {access_token}

Response: 200 OK
{
  "total_activities": 45,
  "total_duration_seconds": 27000,
  "total_distance_meters": 15000,
  "current_streak": 7,
  "best_streak": 12,
  "level_progress": 0.65,
  "xp_to_next_level": 350
}
```

### Activities

#### Record Activity
```http
POST /api/v1/activities
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "pet_id": "...",
  "game_type_id": "walk",
  "started_at": "2024-01-01T10:00:00Z",
  "ended_at": "2024-01-01T10:30:00Z",
  "duration_seconds": 1800,
  "game_data": {
    "distance_meters": 2500,
    "avg_speed_kmh": 4.2
  }
}

Response: 201 Created
{
  "id": "...",
  "xp_earned": 85,
  "pet_level_before": 5,
  "pet_level_after": 5,
  "level_up": false
}
```

## Gamification System

### XP Calculation

**Walk Game:**
- Base: 2 XP per minute
- Distance bonus: 10 XP per kilometer
- Example: 30-minute walk covering 2.5km = 60 + 25 = 85 XP

**Fetch Game:**
- Base: 1 XP per throw
- Combo bonus: +5 XP per 5 consecutive successful catches
- Frenzy mode: 1.5x multiplier

### Level Progression

Formula: `XP Required = (Level - 1)² × 100`

| Level | Total XP Required | XP for This Level |
|-------|------------------|-------------------|
| 1     | 0                | -                 |
| 2     | 100              | 100               |
| 3     | 400              | 300               |
| 4     | 900              | 500               |
| 5     | 1600             | 700               |

### Pet Mood System

Mood degrades based on time since last activity:

| Time Since Activity | Mood     | Color   |
|--------------------|----------|---------|
| < 6 hours          | Happy    | Green   |
| 6-12 hours         | Content  | Blue    |
| 12-24 hours        | Tired    | Yellow  |
| 24-48 hours        | Sad      | Red     |
| > 48 hours         | Bored    | Gray    |

### Streak Tracking

- Consecutive days with at least one activity
- Breaks if no activity for 24 hours
- Contributes to achievements (coming soon)

## Testing

### Backend Tests

Run unit tests:
```bash
cd backend
go test ./... -v
```

Run with coverage:
```bash
go test ./... -cover
go test ./internal/models -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Current coverage:
- `models/`: 60%
- `services/`: 10.7%
- Overall: 3.7%

### Manual API Testing

Example workflow:

1. **Register a user**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

2. **Create a pet**:
```bash
curl -X POST http://localhost:8080/api/v1/pets \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "pet_type_id": "dog",
    "name": "Max",
    "breed": "Golden Retriever"
  }'
```

3. **Record a walk**:
```bash
curl -X POST http://localhost:8080/api/v1/activities \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "pet_id": "YOUR_PET_ID",
    "game_type_id": "walk",
    "started_at": "2024-01-01T10:00:00Z",
    "ended_at": "2024-01-01T10:30:00Z",
    "duration_seconds": 1800,
    "game_data": {
      "distance_meters": 2500
    }
  }'
```

### Mobile Testing

1. Launch the app in Expo
2. Register a new account
3. Create your first pet
4. View pet details and stats
5. Verify mood badges and level display

## Architecture Highlights

### Clean Architecture (Backend)

```
User Request → Handler → Service → Repository → Database
                  ↓
            Middleware (Auth, CORS)
```

- **Handlers**: HTTP request/response handling
- **Services**: Business logic and XP calculations
- **Repositories**: Database queries
- **Models**: Domain entities and validation

### State Management (Mobile)

Using Zustand for lightweight state management:

```typescript
// Example: Pet Store
const { pets, fetchPets, createPet } = usePetStore();

// Example: Auth Store
const { user, login, logout } = useAuthStore();
```

### Extensibility

**Adding new pet types:**
```sql
INSERT INTO pet_types (id, name, icon, config)
VALUES ('bird', 'Bird', '🐦', '{}');
```

**Adding new game types:**
```sql
INSERT INTO game_types (id, name, xp_config, supported_pet_types)
VALUES (
  'agility',
  'Agility Course',
  '{"base_xp_per_obstacle": 5, "time_bonus": true}',
  ARRAY['dog']
);
```

The system automatically supports new types without code changes.

## Troubleshooting

### Backend Issues

**Container won't start:**
```bash
# Check logs
docker-compose logs api

# Rebuild from scratch
docker-compose down
docker-compose build --no-cache
docker-compose up
```

**Database connection error:**
```bash
# Verify PostgreSQL is running
docker-compose ps

# Check database logs
docker-compose logs db

# Reset database
docker-compose down -v
docker-compose up -d
```

### Mobile Issues

**"Network request failed":**
- Ensure backend is running: `curl http://localhost:8080/health`
- Check API_URL in `mobile/src/api/client.ts`
- For physical devices, use local IP instead of localhost
- Verify device and computer are on same Wi-Fi network

**Expo won't start:**
```bash
# Clear cache
npx expo start -c

# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install
```

**Build errors:**
```bash
# Clear watchman
watchman watch-del-all

# Clear Metro bundler cache
npx expo start -c
```

## Current Features Status

✅ **Implemented:**
- User authentication (register, login, JWT)
- Pet CRUD operations
- Activity tracking with XP calculation
- Level progression
- Streak tracking
- Mood system
- Pet statistics
- Mobile UI for auth, pet list, pet details, pet creation

⏳ **Coming Soon (Phase 2):**
- Real-time GPS walk tracking
- Fetch mini-game implementation
- Activity history view
- Daily missions
- Achievements/badges
- Card collection system

## Database Schema

Key tables:

- `users`: User accounts and authentication
- `pets`: Pet profiles with XP and level
- `pet_types`: Extensible pet type registry (dog, cat, etc.)
- `activities`: Activity records with JSONB game data
- `game_types`: Game type registry with XP formulas
- `achievements`: Achievement definitions
- `cards`: Collectible card definitions
- `missions`: Daily mission tracking

Full schema: `backend/migrations/001_initial.up.sql`

## Contributing

This is a personal project, but suggestions are welcome.

## License

MIT
