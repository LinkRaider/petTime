package services

import (
	"testing"
	"time"

	"github.com/joaosantos/pettime/internal/models"
)

func TestCalculateMood(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name             string
		lastActivityAt   *time.Time
		expectedMood     models.Mood
	}{
		{
			name:           "Never walked - bored",
			lastActivityAt: nil,
			expectedMood:   models.MoodBored,
		},
		{
			name:           "Walked 2 hours ago - happy",
			lastActivityAt: timePtr(now.Add(-2 * time.Hour)),
			expectedMood:   models.MoodHappy,
		},
		{
			name:           "Walked 5 hours ago - happy",
			lastActivityAt: timePtr(now.Add(-5 * time.Hour)),
			expectedMood:   models.MoodHappy,
		},
		{
			name:           "Walked 8 hours ago - content",
			lastActivityAt: timePtr(now.Add(-8 * time.Hour)),
			expectedMood:   models.MoodContent,
		},
		{
			name:           "Walked 11 hours ago - content",
			lastActivityAt: timePtr(now.Add(-11 * time.Hour)),
			expectedMood:   models.MoodContent,
		},
		{
			name:           "Walked 18 hours ago - tired",
			lastActivityAt: timePtr(now.Add(-18 * time.Hour)),
			expectedMood:   models.MoodTired,
		},
		{
			name:           "Walked 23 hours ago - tired",
			lastActivityAt: timePtr(now.Add(-23 * time.Hour)),
			expectedMood:   models.MoodTired,
		},
		{
			name:           "Walked 30 hours ago - sad",
			lastActivityAt: timePtr(now.Add(-30 * time.Hour)),
			expectedMood:   models.MoodSad,
		},
		{
			name:           "Walked 47 hours ago - sad",
			lastActivityAt: timePtr(now.Add(-47 * time.Hour)),
			expectedMood:   models.MoodSad,
		},
		{
			name:           "Walked 3 days ago - bored",
			lastActivityAt: timePtr(now.Add(-72 * time.Hour)),
			expectedMood:   models.MoodBored,
		},
		{
			name:           "Walked 1 week ago - bored",
			lastActivityAt: timePtr(now.Add(-168 * time.Hour)),
			expectedMood:   models.MoodBored,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pet := &models.Pet{
				LastActivityAt: tt.lastActivityAt,
			}

			mood := calculateMood(pet)

			if mood != tt.expectedMood {
				t.Errorf("calculateMood() = %v, want %v", mood, tt.expectedMood)
			}
		})
	}
}

func TestMoodTransitions(t *testing.T) {
	// Test that mood degrades over time as expected
	now := time.Now()

	transitions := []struct {
		hoursAgo int
		mood     models.Mood
	}{
		{0, models.MoodHappy},
		{5, models.MoodHappy},
		{6, models.MoodContent},
		{11, models.MoodContent},
		{12, models.MoodTired},
		{23, models.MoodTired},
		{24, models.MoodSad},
		{47, models.MoodSad},
		{48, models.MoodBored},
	}

	for _, tr := range transitions {
		lastActivity := now.Add(-time.Duration(tr.hoursAgo) * time.Hour)
		pet := &models.Pet{
			LastActivityAt: &lastActivity,
		}

		mood := calculateMood(pet)
		if mood != tr.mood {
			t.Errorf("After %d hours: mood = %v, want %v", tr.hoursAgo, mood, tr.mood)
		}
	}
}

func TestMoodBoundaries(t *testing.T) {
	// Test exact boundary conditions
	now := time.Now()

	boundaries := []struct {
		hours        float64
		expectedMood models.Mood
		description  string
	}{
		{5.99, models.MoodHappy, "Just before 6 hours"},
		{6.0, models.MoodContent, "Exactly 6 hours"},
		{11.99, models.MoodContent, "Just before 12 hours"},
		{12.0, models.MoodTired, "Exactly 12 hours"},
		{23.99, models.MoodTired, "Just before 24 hours"},
		{24.0, models.MoodSad, "Exactly 24 hours"},
		{47.99, models.MoodSad, "Just before 48 hours"},
		{48.0, models.MoodBored, "Exactly 48 hours"},
	}

	for _, b := range boundaries {
		lastActivity := now.Add(-time.Duration(b.hours * float64(time.Hour)))
		pet := &models.Pet{
			LastActivityAt: &lastActivity,
		}

		mood := calculateMood(pet)
		if mood != b.expectedMood {
			t.Errorf("%s: mood = %v, want %v", b.description, mood, b.expectedMood)
		}
	}
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}
