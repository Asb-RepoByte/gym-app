-- init.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Muscle groups as enum (you can add more)
CREATE TYPE muscle_group AS ENUM (
    'chest', 'back', 'shoulders', 'biceps', 'triceps',
    'quads', 'hamstrings', 'glutes', 'calves', 'core', 'forearms'
);

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workout session (commute times)
CREATE TABLE workout_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    departure_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    homecoming_time TIMESTAMP,
    companions_count INTEGER DEFAULT 0,
    overall_mood TEXT,
    body_weight_kg DECIMAL(5,2),
    biometric_data JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Exercise library (reusable, currently global for all users for simplicity but can be made user-specific later)
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    exercise_type VARCHAR(50) CHECK (exercise_type IN ('cardio', 'weight_lifting', 'flexibility')),
    is_machine BOOLEAN DEFAULT false,
    is_volume_based BOOLEAN DEFAULT true,  -- false = time based
    target_muscles muscle_group[],
    image_urls JSONB DEFAULT '[]'::jsonb,   -- array of image URLs/paths
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Exercise log (links a session to an exercise, holds failure & feeling)
CREATE TABLE exercise_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id),
    start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    is_failure BOOLEAN DEFAULT false,
    feeling_score INTEGER CHECK (feeling_score BETWEEN 1 AND 5),
    notes TEXT
);

-- Individual sets (weight, reps, timings, pain)
CREATE TABLE sets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    log_id UUID NOT NULL REFERENCES exercise_logs(id) ON DELETE CASCADE,
    set_number INTEGER NOT NULL,
    weight_kg DECIMAL(6,2),
    reps INTEGER,
    duration_seconds INTEGER,   -- for cardio/time-based exercises
    start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    pain_data JSONB,            -- e.g. {"location": "lower_back", "intensity": 3}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_workout_sessions_user ON workout_sessions(user_id);
CREATE INDEX idx_exercise_logs_session ON exercise_logs(session_id);
CREATE INDEX idx_sets_log ON sets(log_id);
