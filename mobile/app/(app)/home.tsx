import { useEffect, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  FlatList,
  ActivityIndicator,
} from 'react-native';
import { Link, router } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { petsApi } from '../../src/api';
import { Pet, User } from '../../src/types';

export default function HomeScreen() {
  const [user, setUser] = useState<User | null>(null);
  const [pets, setPets] = useState<Pet[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const userStr = await AsyncStorage.getItem('user');
      if (userStr) {
        setUser(JSON.parse(userStr));
      }
      const petsData = await petsApi.list();
      setPets(petsData);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = async () => {
    await AsyncStorage.multiRemove(['user', 'access_token', 'refresh_token']);
    router.replace('/(auth)/login');
  };

  const getMoodColor = (mood: string) => {
    const colors: Record<string, string> = {
      happy: '#10B981',
      content: '#3B82F6',
      tired: '#F59E0B',
      sad: '#EF4444',
      bored: '#6B7280',
    };
    return colors[mood] || '#6B7280';
  };

  const renderPetCard = ({ item }: { item: Pet }) => (
    <TouchableOpacity
      style={styles.petCard}
      onPress={() => router.push(`/(app)/pet/${item.id}`)}
    >
      <View style={styles.petInfo}>
        <Text style={styles.petName}>{item.name}</Text>
        <Text style={styles.petType}>{item.pet_type?.name}</Text>
        <View style={styles.stats}>
          <Text style={styles.statText}>Level {item.level}</Text>
          <Text style={styles.statText}>{item.total_xp} XP</Text>
          <Text style={styles.statText}>{item.streak_days} day streak</Text>
        </View>
      </View>
      <View style={[styles.moodBadge, { backgroundColor: getMoodColor(item.mood) }]}>
        <Text style={styles.moodText}>{item.mood}</Text>
      </View>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.welcome}>Welcome, {user?.name || 'User'}!</Text>
        <TouchableOpacity onPress={handleLogout}>
          <Text style={styles.logoutText}>Logout</Text>
        </TouchableOpacity>
      </View>

      {isLoading ? (
        <ActivityIndicator size="large" color="#4F46E5" style={styles.loader} />
      ) : pets.length === 0 ? (
        <View style={styles.empty}>
          <Text style={styles.emptyText}>No pets yet</Text>
          <Link href="/(app)/create-pet" asChild>
            <TouchableOpacity style={styles.addButton}>
              <Text style={styles.addButtonText}>Add Your First Pet</Text>
            </TouchableOpacity>
          </Link>
        </View>
      ) : (
        <FlatList
          data={pets}
          renderItem={renderPetCard}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.list}
          onRefresh={loadData}
          refreshing={isLoading}
        />
      )}

      {pets.length > 0 && (
        <Link href="/(app)/create-pet" asChild>
          <TouchableOpacity style={styles.fab}>
            <Text style={styles.fabText}>+</Text>
          </TouchableOpacity>
        </Link>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F9FAFB',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#E5E7EB',
  },
  welcome: {
    fontSize: 20,
    fontWeight: '700',
    color: '#1F2937',
  },
  logoutText: {
    color: '#4F46E5',
    fontSize: 14,
  },
  loader: {
    marginTop: 50,
  },
  empty: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  emptyText: {
    fontSize: 18,
    color: '#6B7280',
    marginBottom: 20,
  },
  addButton: {
    backgroundColor: '#4F46E5',
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  addButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '700',
  },
  list: {
    padding: 16,
  },
  petCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    flexDirection: 'row',
    justifyContent: 'space-between',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  petInfo: {
    flex: 1,
  },
  petName: {
    fontSize: 20,
    fontWeight: '700',
    color: '#1F2937',
    marginBottom: 4,
  },
  petType: {
    fontSize: 14,
    color: '#6B7280',
    marginBottom: 8,
  },
  stats: {
    flexDirection: 'row',
  },
  statText: {
    fontSize: 12,
    color: '#4F46E5',
    fontWeight: '600',
    marginRight: 12,
  },
  moodBadge: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 12,
    alignSelf: 'flex-start',
  },
  moodText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
    textTransform: 'capitalize',
  },
  fab: {
    position: 'absolute',
    right: 20,
    bottom: 20,
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: '#4F46E5',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 4,
    elevation: 8,
  },
  fabText: {
    color: '#fff',
    fontSize: 32,
    fontWeight: '700',
  },
});
