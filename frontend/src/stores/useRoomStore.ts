import { create } from 'zustand';
import { Room, Seat, Member } from '@/types';

interface RoomStore {
  room: Room | null;
  currentRoomId: string | null;
  seats: Seat[];
  members: Map<string, Member>;
  error: string | null;

  setRoom: (room: Room) => void;
  setSnapshot: (seats: any[], members: any[]) => void;
  setCurrentRoomId: (roomId: string) => void;
  updateSeat: (seatId: string, occupiedById: string | null, username?: string | null) => void;
  addMember: (member: Member) => void;
  removeMember: (userId: string) => void;
  updateMember: (userId: string, isMuted: boolean) => void;
  clearRoom: () => void;
  setError: (error: string | null) => void;
}

export const useRoomStore = create<RoomStore>((set) => ({
  room: null,
  currentRoomId: null,
  seats: [],
  members: new Map(),
  error: null,

  setRoom: (room) => {
    const memberMap = new Map(room.members.map(m => [m.user_id, m]));
    // Enrich seats with occupant usernames
    const enrichedSeats = room.seats.map(seat => ({
      ...seat,
      username: seat.occupied_by_id ? (memberMap.get(seat.occupied_by_id)?.username ?? null) : null,
    }));
    set({ room, seats: enrichedSeats, members: memberMap });
  },

  // Called when room_snapshot WS event arrives — replaces seat+member state with fresh data
  setSnapshot: (snapshotSeats, snapshotMembers) => {
    set((state) => {
      const newMembers = new Map(state.members);
      // Add/update members from snapshot (keeps any that may have already been added)
      snapshotMembers.forEach((m: any) => {
        newMembers.set(m.user_id, { user_id: m.user_id, username: m.username || m.user_id, is_muted: m.is_muted });
      });
      // Update seats with fresh occupancy + username
      const newSeats = state.seats.map(seat => {
        const snap = snapshotSeats.find((s: any) => s.id === seat.id);
        if (!snap) return seat;
        return { ...seat, occupied_by_id: snap.occupied_by_id ?? null, username: snap.username ?? null };
      });
      return { seats: newSeats, members: newMembers };
    });
  },

  setCurrentRoomId: (roomId) => set({ currentRoomId: roomId }),

  updateSeat: (seatId, occupiedById, username) => {
    set((state) => ({
      seats: state.seats.map((seat) =>
        seat.id === seatId
          ? { ...seat, occupied_by_id: occupiedById, username: username ?? null }
          : seat
      ),
    }));
  },

  addMember: (member) => {
    set((state) => {
      const newMembers = new Map(state.members);
      newMembers.set(member.user_id, member);
      return { members: newMembers };
    });
  },

  removeMember: (userId) => {
    set((state) => {
      const newMembers = new Map(state.members);
      newMembers.delete(userId);
      return { members: newMembers };
    });
  },

  updateMember: (userId, isMuted) => {
    set((state) => {
      const newMembers = new Map(state.members);
      const member = newMembers.get(userId);
      if (member) {
        newMembers.set(userId, { ...member, is_muted: isMuted });
      }
      return { members: newMembers };
    });
  },

  clearRoom: () => set({ room: null, currentRoomId: null, seats: [], members: new Map() }),
  setError: (error) => set({ error }),
}));
