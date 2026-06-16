const API_BASE = import.meta.env.DEV ? 'http://localhost:8080/api/v1' : '/api/v1';

export class ApiService {
  private static getAuthHeader(): Record<string, string> {
    const token = localStorage.getItem('token');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  static async register(username: string, email: string, password: string) {
    const res = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, email, password }),
    });
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new Error(body.error || 'Register failed');
    }
    return res.json();
  }

  static async login(email: string, password: string) {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new Error(body.error || 'Login failed');
    }
    return res.json();
  }

  static async getMe() {
    const res = await fetch(`${API_BASE}/auth/me`, {
      headers: this.getAuthHeader(),
    });
    if (!res.ok) throw new Error('Get me failed');
    return res.json();
  }

  static async listRooms(limit = 20, offset = 0) {
    const res = await fetch(`${API_BASE}/rooms?limit=${limit}&offset=${offset}`);
    if (!res.ok) throw new Error('List rooms failed');
    return res.json();
  }

  static async createRoom(name: string, description: string) {
    const res = await fetch(`${API_BASE}/rooms`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...this.getAuthHeader(),
      },
      body: JSON.stringify({ name, description }),
    });
    if (!res.ok) throw new Error('Create room failed');
    return res.json();
  }

  static async getRoom(roomId: string) {
    const res = await fetch(`${API_BASE}/rooms/${roomId}`);
    if (!res.ok) throw new Error('Get room failed');
    return res.json();
  }

  static async deleteRoom(roomId: string) {
    const res = await fetch(`${API_BASE}/rooms/${roomId}`, {
      method: 'DELETE',
      headers: this.getAuthHeader(),
    });
    if (!res.ok) throw new Error('Delete room failed');
    return res.json();
  }

  static async occupySeat(roomId: string, seatId: string) {
    const res = await fetch(`${API_BASE}/rooms/${roomId}/seats/${seatId}/occupy`, {
      method: 'POST',
      headers: this.getAuthHeader(),
    });
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new Error(body.error || 'Occupy seat failed');
    }
    return res.json();
  }

  static async vacateSeat(roomId: string, seatId: string) {
    const res = await fetch(`${API_BASE}/rooms/${roomId}/seats/${seatId}`, {
      method: 'DELETE',
      headers: this.getAuthHeader(),
    });
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new Error(body.error || 'Vacate seat failed');
    }
    return res.json();
  }
}
