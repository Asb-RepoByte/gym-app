from fastapi import FastAPI, HTTPException, Depends
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional, Dict, Any
from datetime import datetime
import psycopg2
from psycopg2.extras import RealDictCursor
import json

app = FastAPI()
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])

def get_db():
    conn = psycopg2.connect(
        "host=db dbname=gym user=abs password=yourpassword",
        cursor_factory=RealDictCursor
    )
    try:
        yield conn
    finally:
        conn.close()

# ---------- Pydantic models ----------
class SessionStart(BaseModel):
    departure_time: Optional[datetime] = None  # None = use current time

class SessionUpdate(BaseModel):
    check_in_time: Optional[datetime] = None
    check_out_time: Optional[datetime] = None
    homecoming_time: Optional[datetime] = None
    companions_count: Optional[int] = None
    overall_mood: Optional[str] = None

class ExerciseCreate(BaseModel):
    name: str
    exercise_type: str
    is_machine: bool = False
    is_volume_based: bool = True
    target_muscles: List[str]
    image_urls: List[str] = []

class ExerciseLogCreate(BaseModel):
    session_id: str
    exercise_id: str
    start_time: Optional[datetime] = None
    is_failure: bool = False

class SetCreate(BaseModel):
    log_id: str
    set_number: int
    weight_kg: Optional[float] = None
    reps: Optional[int] = None
    duration_seconds: Optional[int] = None
    start_time: Optional[datetime] = None
    end_time: Optional[datetime] = None
    pain_data: Optional[Dict[str, Any]] = None

# ---------- Endpoints ----------
@app.post("/sessions")
def start_session(data: SessionStart, db=Depends(get_db)):
    """Start a new workout session (departure time defaults to now)"""
    dep_time = data.departure_time or datetime.now()
    with db.cursor() as cur:
        cur.execute(
            "INSERT INTO workout_sessions (departure_time) VALUES (%s) RETURNING id",
            (dep_time,)
        )
        db.commit()
        return {"session_id": cur.fetchone()["id"]}

@app.patch("/sessions/{session_id}")
def update_session(session_id: str, data: SessionUpdate, db=Depends(get_db)):
    """Update check‑in, check‑out, homecoming, companions, mood"""
    updates = []
    params = []
    for field, value in data.dict(exclude_unset=True).items():
        if value is not None:
            updates.append(f"{field} = %s")
            params.append(value)
    if not updates:
        raise HTTPException(400, "No fields to update")
    params.append(session_id)
    with db.cursor() as cur:
        cur.execute(f"UPDATE workout_sessions SET {', '.join(updates)} WHERE id = %s", params)
        db.commit()
        return {"status": "updated"}

@app.get("/sessions/{session_id}")
def get_session(session_id: str, db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT * FROM workout_sessions WHERE id = %s", (session_id,))
        session = cur.fetchone()
        if not session:
            raise HTTPException(404, "Session not found")
        return session

# Exercises (library)
@app.post("/exercises")
def create_exercise(ex: ExerciseCreate, db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("""
            INSERT INTO exercises (name, exercise_type, is_machine, is_volume_based, target_muscles, image_urls)
            VALUES (%s, %s, %s, %s, %s, %s) RETURNING id
        """, (ex.name, ex.exercise_type, ex.is_machine, ex.is_volume_based, ex.target_muscles, json.dumps(ex.image_urls)))
        db.commit()
        return {"exercise_id": cur.fetchone()["id"]}

@app.get("/exercises")
def list_exercises(db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT id, name, exercise_type, target_muscles, image_urls FROM exercises ORDER BY name")
        return cur.fetchall()

# Exercise logs (failure is here)
@app.post("/exercise_logs")
def start_exercise_log(data: ExerciseLogCreate, db=Depends(get_db)):
    start = data.start_time or datetime.now()
    with db.cursor() as cur:
        cur.execute("""
            INSERT INTO exercise_logs (session_id, exercise_id, start_time, is_failure)
            VALUES (%s, %s, %s, %s) RETURNING id
        """, (data.session_id, data.exercise_id, start, data.is_failure))
        db.commit()
        return {"log_id": cur.fetchone()["id"]}

@app.patch("/exercise_logs/{log_id}/end")
def end_exercise_log(log_id: str, end_time: Optional[datetime] = None, feeling_score: Optional[int] = None, notes: Optional[str] = None, db=Depends(get_db)):
    end = end_time or datetime.now()
    with db.cursor() as cur:
        cur.execute("UPDATE exercise_logs SET end_time = %s, feeling_score = COALESCE(%s, feeling_score), notes = COALESCE(%s, notes) WHERE id = %s",
                    (end, feeling_score, notes, log_id))
        db.commit()
        return {"status": "ended"}

# Sets
@app.post("/sets")
def add_set(data: SetCreate, db=Depends(get_db)):
    start = data.start_time or datetime.now()
    with db.cursor() as cur:
        cur.execute("""
            INSERT INTO sets (log_id, set_number, weight_kg, reps, duration_seconds, start_time, end_time, pain_data)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s) RETURNING id
        """, (data.log_id, data.set_number, data.weight_kg, data.reps, data.duration_seconds,
              start, data.end_time, json.dumps(data.pain_data) if data.pain_data else None))
        db.commit()
        return {"set_id": cur.fetchone()["id"]}

@app.get("/exercise_logs/{log_id}/sets")
def get_sets(log_id: str, db=Depends(get_db)):
    with db.cursor() as cur:
        cur.execute("SELECT * FROM sets WHERE log_id = %s ORDER BY set_number", (log_id,))
        return cur.fetchall()

# ---------- Export / Import (portability) ----------
@app.get("/export")
def export_all(db=Depends(get_db)):
    """Export all data as a single JSON object"""
    with db.cursor() as cur:
        cur.execute("SELECT * FROM workout_sessions ORDER BY created_at")
        sessions = cur.fetchall()
        cur.execute("SELECT * FROM exercises")
        exercises = cur.fetchall()
        cur.execute("SELECT * FROM exercise_logs")
        logs = cur.fetchall()
        cur.execute("SELECT * FROM sets")
        sets = cur.fetchall()
    return {
        "sessions": sessions,
        "exercises": exercises,
        "exercise_logs": logs,
        "sets": sets
    }

@app.post("/import")
def import_all(data: dict, db=Depends(get_db)):
    """Import previously exported data (clears existing tables)"""
    with db.cursor() as cur:
        # Clear in correct order (due to foreign keys)
        cur.execute("TRUNCATE sets, exercise_logs, exercises, workout_sessions RESTART IDENTITY CASCADE")
        for table in ["sessions", "exercises", "exercise_logs", "sets"]:
            rows = data.get(table, [])
            if not rows:
                continue
            for row in rows:
                # Remove 'id' if present to let DB generate new UUIDs (or keep if you want exact)
                # Here we keep the original id – useful for exact restore
                cols = list(row.keys())
                placeholders = ','.join(['%s'] * len(cols))
                cur.execute(f"INSERT INTO {table} ({','.join(cols)}) VALUES ({placeholders})", list(row.values()))
        db.commit()
    return {"status": "imported"}
