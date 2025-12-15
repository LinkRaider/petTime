import React, { useState, useEffect } from 'react';
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
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { RootStackParamList } from '../../navigation/RootNavigator';
import { usePetStore } from '../../store/petStore';

type NavigationProp = NativeStackNavigationProp<RootStackParamList>;

export default function CreatePetScreen() {
  const navigation = useNavigation<NavigationProp>();
  const { petTypes, fetchPetTypes, createPet, isLoading } = usePetStore();
  const [name, setName] = useState('');
  const [breed, setBreed] = useState('');
  const [selectedType, setSelectedType] = useState<string | null>(null);

  useEffect(() => {
    fetchPetTypes();
  }, []);

  const handleCreate = async () => {
    if (!name || !selectedType) {
      Alert.alert('Error', 'Please enter a name and select a pet type');
      return;
    }

    try {
      await createPet({
        pet_type_id: selectedType,
        name,
        breed: breed || undefined,
      });
      Alert.alert('Success', `${name} has been added!`);
      navigation.goBack();
    } catch (error) {
      Alert.alert('Error', 'Failed to create pet. Please try again.');
    }
  };

  return (
    <ScrollView style={styles.container}>
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
    fontWeight: 'bold',
  },
});
