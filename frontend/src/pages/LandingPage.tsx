import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/useAuthStore';
import { useUIStore } from '@/stores/useUIStore';
import { ApiService } from '@/services/api';
import type { Room } from '@/types';

export default function LandingPage() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const addToast = useUIStore((state) => state.addToast);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newRoom, setNewRoom] = useState({ name: '', description: '' });
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadRooms();
  }, []);

  const loadRooms = async () => {
    try {
      const response = await ApiService.listRooms();
      setRooms(response.rooms || []);
    } catch (error) {
      addToast('Failed to load rooms', 'error');
    }
  };

  const handleCreateRoom = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const room = await ApiService.createRoom(newRoom.name, newRoom.description);
      addToast('Room created successfully!', 'success');
      setNewRoom({ name: '', description: '' });
      setShowCreateForm(false);
      navigate(`/rooms/${room.id}`);
    } catch (error) {
      addToast('Failed to create room', 'error');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-12">
          <h1 className="text-4xl font-bold">Voice Chat Rooms</h1>
          <div className="flex gap-4">
            <span className="text-gray-400">Welcome, {user?.username}</span>
            <button
              onClick={() => {
                logout();
                navigate('/login');
              }}
              className="px-4 py-2 bg-red-600 hover:bg-red-700 rounded"
            >
              Logout
            </button>
          </div>
        </div>

        {!showCreateForm ? (
          <button
            onClick={() => setShowCreateForm(true)}
            className="mb-8 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-semibold"
          >
            Create New Room
          </button>
        ) : (
          <form onSubmit={handleCreateRoom} className="mb-8 bg-gray-800 p-6 rounded-lg">
            <h2 className="text-2xl font-bold mb-4">Create a New Room</h2>
            <input
              type="text"
              placeholder="Room Name"
              value={newRoom.name}
              onChange={(e) => setNewRoom({ ...newRoom, name: e.target.value })}
              className="w-full px-4 py-2 mb-4 bg-gray-700 rounded border border-gray-600 focus:border-blue-500 outline-none"
              required
            />
            <textarea
              placeholder="Description (optional)"
              value={newRoom.description}
              onChange={(e) => setNewRoom({ ...newRoom, description: e.target.value })}
              className="w-full px-4 py-2 mb-4 bg-gray-700 rounded border border-gray-600 focus:border-blue-500 outline-none"
              rows={3}
            />
            <div className="flex gap-2">
              <button
                type="submit"
                disabled={isLoading}
                className="px-6 py-2 bg-green-600 hover:bg-green-700 rounded disabled:opacity-50"
              >
                {isLoading ? 'Creating...' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => setShowCreateForm(false)}
                className="px-6 py-2 bg-gray-600 hover:bg-gray-700 rounded"
              >
                Cancel
              </button>
            </div>
          </form>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {rooms.map((room) => (
            <div
              key={room.id}
              className="bg-gray-800 rounded-lg p-6 hover:bg-gray-750 cursor-pointer transition"
              onClick={() => navigate(`/rooms/${room.id}`)}
            >
              <h3 className="text-xl font-bold mb-2">{room.name}</h3>
              <p className="text-gray-400 mb-4">{room.description}</p>
              <div className="flex justify-between text-sm text-gray-500">
                <span>Seats: {room.members?.length || 0}/18</span>
                <span>{room.is_public ? 'Public' : 'Private'}</span>
              </div>
            </div>
          ))}
        </div>

        {rooms.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-400 text-lg">No rooms available. Create one to get started!</p>
          </div>
        )}
      </div>
    </div>
  );
}
