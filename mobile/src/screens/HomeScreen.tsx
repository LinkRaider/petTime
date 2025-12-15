import React, { useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  FlatList,
  ActivityIndicator,
} from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { RootStackParamList } from '../navigation/RootNavigator';
import { useAuthStore } from '../store/authStore';
import { usePetStore } from '../store/petStore';

type NavigationProp = NativeStackNavigationProp<RootStackParamList>;

export default function HomeScreen() {
  const navigation = useNavigation<NavigationProp>();
  const { user, logout } = useAuthStore();
  const { pets, fetchPets, selectPet, isLoading } = usePetStore();

  useEffect(() => {
    fetchPets();
  }, []);

  const handleLogout = async () => {
    await logout();
  };

  const handlePetPress = (pet: any) => {
    selectPet(pet);
    navigation.navigate('PetDetail', { petId: pet.id });
  };

  const renderPetCard = ({ item }: any) => (
    <TouchableOpacity
      style={styles.petCard}
      onPress={() => handlePetPress(item)}
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
      <View style={styles.moodBadge}>
        <Text style={styles.moodText}>{item.mood}</Text>
      </View>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.welcome}>Welcome, {user?.name}!</Text>
        <TouchableOpacity onPress={handleLogout}>
          <Text style={styles.logoutText}>Logout</Text>
        </TouchableOpacity>
      </View>

      {isLoading ? (
        <ActivityIndicator size="large" color="#4F46E5" style={styles.loader} />
      ) : pets.length === 0 ? (
        <View style={styles.empty}>
          <Text style={styles.emptyText}>No pets yet</Text>
          <TouchableOpacity
            style={styles.addButton}
            onPress={() => navigation.navigate('CreatePet')}
          >
            <Text style={styles.addButtonText}>Add Your First Pet</Text>
          </TouchableOpacity>
        </View>
      ) : (
        <FlatList
          data={pets}
          renderItem={renderPetCard}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.list}
        />
      )}

      {pets.length > 0 && (
        <TouchableOpacity
          style={styles.fab}
          onPress={() => navigation.navigate('CreatePet')}
        >
          <Text style={styles.fabText}>+</Text>
        </TouchableOpacity>
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
    fontWeight: 'bold',
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
    fontWeight: 'bold',
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
    fontWeight: 'bold',
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
    backgroundColor: '#DBEAFE',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 12,
    alignSelf: 'flex-start',
  },
  moodText: {
    color: '#1E40AF',
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
    fontWeight: 'bold',
  },
});
