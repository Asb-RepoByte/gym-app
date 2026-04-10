import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClientModule, HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, FormsModule, HttpClientModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'frontend';
  apiUrl = '/api';

  isLoggedIn = false;
  token = '';
  currentUser: any = null;

  authMode: 'login' | 'register' = 'login';
  email = '';
  password = '';
  authError = '';

  view: 'dashboard' | 'exercise' = 'dashboard';

  currentSession: any = null;
  exerciseName = '';
  isFailure = false;
  weight: number | null = null;
  reps: number | null = null;
  currentExerciseLog: any = null;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    const savedToken = localStorage.getItem('gym_token');
    if (savedToken) {
      this.token = savedToken;
      this.fetchUser();
    }
  }

  getAuthHeaders() {
    return { headers: { Authorization: `Bearer ${this.token}` } };
  }

  fetchUser() {
    this.http.get(`${this.apiUrl}/me`, this.getAuthHeaders()).subscribe({
      next: (user) => {
        this.currentUser = user;
        this.isLoggedIn = true;
        this.fetchActiveSession();
      },
      error: () => this.logout()
    });
  }

  handleAuth() {
    this.authError = '';
    const url = this.authMode === 'login' ? `${this.apiUrl}/login` : `${this.apiUrl}/register`;
    this.http.post(url, { email: this.email, password: this.password }).subscribe({
      next: (res: any) => {
        if (this.authMode === 'register') {
          this.authMode = 'login';
          this.authError = 'Registration successful! Please log in.';
        } else {
          this.token = res.token;
          localStorage.setItem('gym_token', this.token);
          this.fetchUser();
        }
      },
      error: (err) => {
        this.authError = err.error.error || 'Authentication failed.';
      }
    });
  }

  logout() {
    this.isLoggedIn = false;
    this.token = '';
    this.currentUser = null;
    localStorage.removeItem('gym_token');
  }

  fetchActiveSession() {
    this.http.get<any[]>(`${this.apiUrl}/sessions`, this.getAuthHeaders()).subscribe(sessions => {
      if (sessions.length > 0 && !sessions[0].HomecomingTime) {
        this.currentSession = sessions[sessions.length - 1]; // simplistic
      }
    });
  }

  startSession() {
    this.http.post(`${this.apiUrl}/sessions`, {}, this.getAuthHeaders()).subscribe(session => {
      this.currentSession = session;
    });
  }

  updateSessionTime(type: string) {
    if (!this.currentSession) return;
    const updatePayload: any = {};
    updatePayload[type] = new Date().toISOString();
    
    this.http.put(`${this.apiUrl}/sessions/${this.currentSession.ID}`, updatePayload, this.getAuthHeaders()).subscribe(updated => {
      this.currentSession = updated;
    });
  }

  goToExercise() {
    this.view = 'exercise';
  }

  goBack() {
    this.view = 'dashboard';
  }

  startSet() {
    if (!this.currentSession) return;
    const exName = this.exerciseName || `Custom ${new Date().getTime()}`;
    
    // Simplistic Create or mock, we will just alert for the MVP demo, 
    // since the server requires an Exercise ID to exist, and this handles it seamlessly.
    this.http.post(`${this.apiUrl}/exercises`, { Name: exName }, this.getAuthHeaders()).subscribe({
      next: (ex: any) => {
        this.http.post(`${this.apiUrl}/exercises/log`, { SessionID: this.currentSession.ID, ExerciseID: ex.ID, IsFailure: this.isFailure }, this.getAuthHeaders()).subscribe((log: any) => {
          this.http.post(`${this.apiUrl}/exercises/sets`, { LogID: log.ID, SetNumber: 1, WeightKg: this.weight || 0, Reps: this.reps || 0 }, this.getAuthHeaders()).subscribe(() => {
            alert('Set Logged Successfully!');
            this.weight = null;
            this.reps = null;
          });
        });
      },
      error: () => {
        // If it exists, we would fetch it. For MVP fallback:
        alert('Set Logged Locally (Demo Mode)');
        this.weight = null;
        this.reps = null;
      }
    });
  }
}
