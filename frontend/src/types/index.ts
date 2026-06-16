export interface User {
  id: string;
  username: string;
  email: string;
}

export interface Room {
  id: string;
  name: string;
  description: string;
  owner_id: string;
  max_users: number;
  is_public: boolean;
  member_count: number;
  seats: Seat[];
  members: Member[];
}

export interface Seat {
  id: string;
  row: number;
  col: number;
  occupied_by_id: string | null;
  username?: string | null;
}

export interface Member {
  user_id: string;
  username: string;
  is_muted: boolean;
}

export interface AudioStream {
  userId: string;
  stream: MediaStream | null;
}

export interface Message {
  event: string;
  payload: Record<string, any>;
}
