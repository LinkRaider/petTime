package models

import "testing"

func TestCalculateLevel(t *testing.T) {
	tests := []struct {
		name     string
		xp       int
		expected int
	}{
		{"Level 1 - 0 XP", 0, 1},
		{"Level 1 - 99 XP", 99, 1},
		{"Level 2 - 100 XP", 100, 2},
		{"Level 2 - 399 XP", 399, 2},
		{"Level 3 - 400 XP", 400, 3},
		{"Level 3 - 899 XP", 899, 3},
		{"Level 4 - 900 XP", 900, 4},
		{"Level 5 - 1600 XP", 1600, 5},
		{"Level 10 - 8100 XP", 8100, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateLevel(tt.xp)
			if result != tt.expected {
				t.Errorf("CalculateLevel(%d) = %d, want %d", tt.xp, result, tt.expected)
			}
		})
	}
}

func TestXPForLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"Level 1", 1, 0},
		{"Level 2", 2, 100},
		{"Level 3", 3, 400},
		{"Level 4", 4, 900},
		{"Level 5", 5, 1600},
		{"Level 10", 10, 8100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XPForLevel(tt.level)
			if result != tt.expected {
				t.Errorf("XPForLevel(%d) = %d, want %d", tt.level, result, tt.expected)
			}
		})
	}
}

func TestXPToNextLevel(t *testing.T) {
	tests := []struct {
		name       string
		currentXP  int
		expected   int
	}{
		{"0 XP needs 100 for level 2", 0, 100},
		{"50 XP needs 50 for level 2", 50, 50},
		{"99 XP needs 1 for level 2", 99, 1},
		{"100 XP needs 300 for level 3", 100, 300},
		{"250 XP needs 150 for level 3", 250, 150},
		{"399 XP needs 1 for level 3", 399, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XPToNextLevel(tt.currentXP)
			if result != tt.expected {
				t.Errorf("XPToNextLevel(%d) = %d, want %d", tt.currentXP, result, tt.expected)
			}
		})
	}
}

func TestLevelProgression(t *testing.T) {
	// Verify the level progression is consistent
	for level := 1; level <= 10; level++ {
		xpForCurrentLevel := XPForLevel(level)
		xpForNextLevel := XPForLevel(level + 1)

		// XP should always increase for higher levels
		if level > 1 && xpForNextLevel <= xpForCurrentLevel {
			t.Errorf("Level %d: XP progression broken. Current: %d, Next: %d",
				level, xpForCurrentLevel, xpForNextLevel)
		}

		// Verify CalculateLevel works for the exact XP threshold
		calculatedLevel := CalculateLevel(xpForCurrentLevel)
		if calculatedLevel != level {
			t.Errorf("Level %d: CalculateLevel(%d) = %d, want %d",
				level, xpForCurrentLevel, calculatedLevel, level)
		}
	}
}
