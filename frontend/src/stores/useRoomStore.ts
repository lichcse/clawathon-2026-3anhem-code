import { create } from 'zustand';
import { Room, Seat, Member } from '@/types';

interface RoomStore {
  room: Room | null;
  currentRoomId: string | null;
  seats: Seat[];
  members: Map<string, Member>;
  error: string | null;

  setRoom: (room: Room) => void;
  setCurrentRoomId: (roomId: string) => void;
  updateSeat: (seatId: string, occupiedById: string | null) => void;
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
    set({ room, seats: room.seats, members: memberMap });
  },

  setCurrentRoomId: (roomId) => set({ currentRoomId: roomId }),

  updateSeat: (seatId, occupiedById) => {
    set((state) => ({
      seats: state.seats.map((seat) =>
        seat.id === seatId ? { ...seat, occupied_by_id: occupiedById } : seat
      ),
    }));
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
