import { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { router } from 'expo-router';
import { petsApi } from '../../src/api';
import { PetType } from '../../src/types';

export default function CreatePetScreen() {
  const [petTypes, setPetTypes] = useState<PetType[]>([]);
  const [name, setName] = useState('');
  const [breed, setBreed] = useState('');
  const [selectedType, setSelectedType] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadPetTypes();
  }, []);

  const loadPetTypes = async () => {
    try {
      const types = await petsApi.getTypes();
      setPetTypes(types);
    } catch (error) {
      console.error('Error loading pet types:', error);
    }
  };

  const handleCreate = async () => {
    if (!name || !selectedType) {
      Alert.alert('Error', 'Please enter a name and select a pet type');
      return;
    }

    setIsLoading(true);
    try {
      await petsApi.create({
        pet_type_id: selectedType,
        name,
        breed: breed || undefined,
      });
      Alert.alert('Success', `${name} has been added!`);
      router.back();
    } catch (error) {
      Alert.alert('Error', 'Failed to create pet. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={() => router.back()}>
          <Text style={styles.backText}>← Back</Text>
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Add Pet</Text>
        <View style={{ width: 50 }} />
      </View>

      <View style={styles.content}>
        <Text style={styles.label}>Pet Name *</Text>
        <TextInput
          style={styles.input}
          placeholder="Enter pet name"
          value={name}
          onChangeText={setName}
          editable={!isLoading}
        />

        <Text style={styles.label}>Pet Type *</Text>
        <View style={styles.typeContainer}>
          {petTypes.map((type) => (
            <TouchableOpacity
              key={type.id}
              style={[
                styles.typeButton,
                selectedType === type.id && styles.typeButtonSelected,
              ]}
              onPress={() => setSelectedType(type.id)}
              disabled={isLoading}
            >
              <Text
                style={[
                  styles.typeText,
                  selectedType === type.id && styles.typeTextSelected,
                ]}
              >
                {type.name}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        <Text style={styles.label}>Breed (Optional)</Text>
        <TextInput
          style={styles.input}
          placeholder="Enter breed"
          value={breed}
          onChangeText={setBreed}
          editable={!isLoading}
        />

        <TouchableOpacity
          style={[styles.createButton, isLoading && styles.buttonDisabled]}
          onPress={handleCreate}
          disabled={isLoading}
        >
          {isLoading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.createButtonText}>Add Pet</Text>
          )}
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
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#E5E7EB',
  },
  backText: {
    color: '#4F46E5',
    fontSize: 16,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#1F2937',
  },
  content: {
    padding: 20,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
    color: '#1F2937',
    marginBottom: 8,
    marginTop: 16,
  },
  input: {
    backgroundColor: '#fff',
    borderWidth: 1,
    borderColor: '#D1D5DB',
    borderRadius: 8,
    padding: 16,
    fontSize: 16,
  },
  typeContainer: {
    flexDirection: 'row',
  },
  typeButton: {
    flex: 1,
    backgroundColor: '#fff',
    borderWidth: 2,
    borderColor: '#D1D5DB',
    borderRadius: 8,
    padding: 16,
    alignItems: 'center',
    marginRight: 12,
  },
  typeButtonSelected: {
    borderColor: '#4F46E5',
    backgroundColor: '#EEF2FF',
  },
  typeText: {
    fontSize: 16,
    color: '#6B7280',
    fontWeight: '600',
  },
  typeTextSelected: {
    color: '#4F46E5',
  },
  createButton: {
    backgroundColor: '#4F46E5',
    borderRadius: 8,
    padding: 16,
    alignItems: 'center',
    marginTop: 32,
  },
  buttonDisabled: {
    opacity: 0.6,
  },
  createButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '700',
  },
});
