package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

// JSONB custom type
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &j)
}

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type WorkoutSession struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	DepartureTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CheckInTime    *time.Time
	CheckOutTime   *time.Time
	HomecomingTime *time.Time
	CompanionsCount int       `gorm:"default:0"`
	OverallMood    string
	BodyWeightKg   float64   `gorm:"type:decimal(5,2)"`
	BiometricData  JSONB     `gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

type Exercise struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name           string    `gorm:"uniqueIndex;not null"`
	ExerciseType   string
	IsMachine      bool      `gorm:"default:false"`
	IsVolumeBased  bool      `gorm:"default:true"`
	// TargetMuscles omit or map to string array physically, postgres array needs special driver like pq.StringArray. We'll use simple jsonb or ignore it for strict ORM here, or pq.StringArray.
	// For simplicity, we can store JSONB or a simple comma-separated string in Go if we don't strictly bind it, or use github.com/lib/pq pq.StringArray. Let's use generic JSONB.
	TargetMuscles  JSONB     `gorm:"type:jsonb;default:'[]'"` 
	ImageUrls      JSONB     `gorm:"type:jsonb;default:'[]'"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

type ExerciseLog struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SessionID    uuid.UUID `gorm:"type:uuid;not null"`
	ExerciseID   uuid.UUID `gorm:"type:uuid;not null"`
	StartTime    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	EndTime      *time.Time
	IsFailure    bool      `gorm:"default:false"`
	FeelingScore int
	Notes        string
}

type Set struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	LogID           uuid.UUID `gorm:"type:uuid;not null"`
	SetNumber       int       `gorm:"not null"`
	WeightKg        float64   `gorm:"type:decimal(6,2)"`
	Reps            int
	DurationSeconds int
	StartTime       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	EndTime         *time.Time
	PainData        JSONB     `gorm:"type:jsonb"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}
