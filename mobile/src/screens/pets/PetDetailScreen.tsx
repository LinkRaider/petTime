import React, { useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { useRoute, useNavigation, RouteProp } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { RootStackParamList } from '../../navigation/RootNavigator';
import { usePetStore } from '../../store/petStore';

type PetDetailRouteProp = RouteProp<RootStackParamList, 'PetDetail'>;
type NavigationProp = NativeStackNavigationProp<RootStackParamList>;

export default function PetDetailScreen() {
  const route = useRoute<PetDetailRouteProp>();
  const navigation = useNavigation<NavigationProp>();
  const { selectedPet, stats, fetchPetStats, isLoading } = usePetStore();

  useEffect(() => {
    if (route.params?.petId) {
      fetchPetStats(route.params.petId);
    }
  }, [route.params?.petId]);

  if (isLoading || !selectedPet) {
    return (
      <View style={styles.loader}>
        <ActivityIndicator size="large" color="#4F46E5" />
      </View>
    );
  }

  const getMoodEmoji = (mood: string) => {
    const emojis: Record<string, string> = {
      happy: '😊',
      content: '😌',
      tired: '😴',
      sad: '😢',
      bored: '😐',
    };
    return emojis[mood] || '🐾';
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

  return (
    <ScrollView style={styles.container}>
      {/* Pet Header */}
      <View style={styles.header}>
        <View style={styles.petInfo}>
          <Text style={styles.petName}>{selectedPet.name}</Text>
          <Text style={styles.petBreed}>
            {selectedPet.breed || selectedPet.pet_type?.name}
          </Text>
        </View>
        <View
          style={[
            styles.moodBadge,
            { backgroundColor: getMoodColor(selectedPet.mood) },
          ]}
        >
          <Text style={styles.moodEmoji}>{getMoodEmoji(selectedPet.mood)}</Text>
          <Text style={styles.moodText}>{selectedPet.mood}</Text>
        </View>
      </View>

      {/* Level & XP */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Level & Experience</Text>
        <View style={styles.levelContainer}>
          <View style={styles.levelBadge}>
            <Text style={styles.levelNumber}>{selectedPet.level}</Text>
            <Text style={styles.levelLabel}>Level</Text>
          </View>
          <View style={styles.xpInfo}>
            <Text style={styles.xpText}>{selectedPet.total_xp} XP</Text>
            {stats && (
              <>
                <View style={styles.progressBar}>
                  <View
                    style={[
                      styles.progressFill,
                      { width: `${stats.level_progress * 100}%` },
                    ]}
                  />
                </View>
                <Text style={styles.xpToNext}>
                  {stats.xp_to_next_level} XP to next level
                </Text>
              </>
            )}
          </View>
        </View>
      </View>

      {/* Stats */}
      {stats && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Statistics</Text>
          <View style={styles.statsGrid}>
            <View style={styles.statCard}>
              <Text style={styles.statValue}>{stats.total_activities}</Text>
              <Text style={styles.statLabel}>Activities</Text>
            </View>
            <View style={styles.statCard}>
              <Text style={styles.statValue}>{stats.current_streak}</Text>
              <Text style={styles.statLabel}>Day Streak</Text>
            </View>
            <View style={styles.statCard}>
              <Text style={styles.statValue}>
                {Math.round(stats.total_duration_seconds / 60)}
              </Text>
              <Text style={styles.statLabel}>Minutes</Text>
            </View>
            <View style={styles.statCard}>
              <Text style={styles.statValue}>
                {(stats.total_distance_meters / 1000).toFixed(1)}
              </Text>
              <Text style={styles.statLabel}>Kilometers</Text>
            </View>
          </View>
        </View>
      )}

      {/* Action Buttons */}
      <View style={styles.actions}>
        <TouchableOpacity style={styles.actionButton}>
          <Text style={styles.actionButtonText}>Start Walk</Text>
        </TouchableOpacity>
        <TouchableOpacity style={[styles.actionButton, styles.secondaryButton]}>
          <Text style={[styles.actionButtonText, styles.secondaryButtonText]}>
            Play Fetch
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F9FAFB',
  },
  loader: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  header: {
    backgroundColor: '#fff',
    padding: 20,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    borderBottomWidth: 1,
    borderBottomColor: '#E5E7EB',
  },
  petInfo: {
    flex: 1,
  },
  petName: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#1F2937',
  },
  petBreed: {
    fontSize: 16,
    color: '#6B7280',
    marginTop: 4,
  },
  moodBadge: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 20,
    flexDirection: 'row',
    alignItems: 'center',
  },
  moodEmoji: {
    fontSize: 20,
    marginRight: 6,
  },
  moodText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
    textTransform: 'capitalize',
  },
  section: {
    backgroundColor: '#fff',
    padding: 20,
    marginTop: 12,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#1F2937',
    marginBottom: 16,
  },
  levelContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  levelBadge: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: '#4F46E5',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 20,
  },
  levelNumber: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#fff',
  },
  levelLabel: {
    fontSize: 12,
    color: '#E0E7FF',
  },
  xpInfo: {
    flex: 1,
  },
  xpText: {
    fontSize: 20,
    fontWeight: '600',
    color: '#1F2937',
    marginBottom: 8,
  },
  progressBar: {
    height: 8,
    backgroundColor: '#E5E7EB',
    borderRadius: 4,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: '#4F46E5',
  },
  xpToNext: {
    fontSize: 12,
    color: '#6B7280',
    marginTop: 4,
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginHorizontal: -6,
  },
  statCard: {
    flex: 1,
    minWidth: '45%',
    backgroundColor: '#F9FAFB',
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
    margin: 6,
  },
  statValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#4F46E5',
  },
  statLabel: {
    fontSize: 12,
    color: '#6B7280',
    marginTop: 4,
  },
  actions: {
    padding: 20,
  },
  actionButton: {
    backgroundColor: '#4F46E5',
    borderRadius: 8,
    padding: 16,
    alignItems: 'center',
    marginBottom: 12,
  },
  secondaryButton: {
    backgroundColor: '#fff',
    borderWidth: 2,
    borderColor: '#4F46E5',
  },
  actionButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  secondaryButtonText: {
    color: '#4F46E5',
  },
});
