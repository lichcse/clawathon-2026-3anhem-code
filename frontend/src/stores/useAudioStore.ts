import { create } from 'zustand';

interface AudioStreamData {
  userId: string;
  audioContext?: AudioContext;
  source?: AudioBufferSourceNode;
}

interface AudioStore {
  isMicOn: boolean;
  volume: number;
  activeStreams: Map<string, AudioStreamData>;
  localStream: MediaStream | null;
  audioContext: AudioContext | null;

  toggleMic: () => void;
  setVolume: (level: number) => void;
  addStream: (userId: string) => void;
  removeStream: (userId: string) => void;
  setLocalStream: (stream: MediaStream | null) => void;
  initAudioContext: () => AudioContext;
  clearStreams: () => void;
}

export const useAudioStore = create<AudioStore>((set, get) => ({
  isMicOn: false,
  volume: 0.8,
  activeStreams: new Map(),
  localStream: null,
  audioContext: null,

  toggleMic: () => set((state) => ({ isMicOn: !state.isMicOn })),

  setVolume: (level) => set({ volume: Math.max(0, Math.min(1, level)) }),

  addStream: (userId) => {
    set((state) => {
      const newStreams = new Map(state.activeStreams);
      newStreams.set(userId, { userId });
      return { activeStreams: newStreams };
    });
  },

  removeStream: (userId) => {
    set((state) => {
      const newStreams = new Map(state.activeStreams);
      newStreams.delete(userId);
      return { activeStreams: newStreams };
    });
  },

  setLocalStream: (stream) => set({ localStream: stream }),

  initAudioContext: () => {
    const state = get();
    if (state.audioContext) return state.audioContext;
    const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
    set({ audioContext });
    return audioContext;
  },

  clearStreams: () => set({ activeStreams: new Map(), localStream: null }),
}));
