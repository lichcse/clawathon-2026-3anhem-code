import { create } from 'zustand';

export interface Toast {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info' | 'warning';
}

interface UIStore {
  toasts: Toast[];
  isConnecting: boolean;
  isConnected: boolean;

  addToast: (message: string, type?: Toast['type']) => void;
  removeToast: (id: string) => void;
  setConnecting: (connecting: boolean) => void;
  setConnected: (connected: boolean) => void;
}

export const useUIStore = create<UIStore>((set) => ({
  toasts: [],
  isConnecting: false,
  isConnected: false,

  addToast: (message, type = 'info') => {
    const id = Math.random().toString(36).substr(2, 9);
    set((state) => ({
      toasts: [...state.toasts, { id, message, type }],
    }));

    // Auto-remove after 5 seconds
    setTimeout(() => {
      set((state) => ({
        toasts: state.toasts.filter((t) => t.id !== id),
      }));
    }, 5000);
  },

  removeToast: (id) => {
    set((state) => ({
      toasts: state.toasts.filter((t) => t.id !== id),
    }));
  },

  setConnecting: (connecting) => set({ isConnecting: connecting }),
  setConnected: (connected) => set({ isConnected: connected }),
}));
