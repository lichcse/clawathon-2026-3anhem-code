import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useRoomStore } from '@/stores/useRoomStore';
import { useAuthStore } from '@/stores/useAuthStore';
import { useAudioStore } from '@/stores/useAudioStore';
import { useUIStore } from '@/stores/useUIStore';
import { ApiService } from '@/services/api';
import { WebSocketService } from '@/services/websocket';
import { audioService } from '@/services/audio';

export default function RoomPage() {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const token = useAuthStore((state) => state.token);
  const room = useRoomStore((state) => state.room);
  const setRoom = useRoomStore((state) => state.setRoom);
  const seats = useRoomStore((state) => state.seats);
  const members = useRoomStore((state) => state.members);
  const updateSeat = useRoomStore((state) => state.updateSeat);
  const setSnapshot = useRoomStore((state) => state.setSnapshot);
  const addMember = useRoomStore((state) => state.addMember);
  const removeMember = useRoomStore((state) => state.removeMember);
  const isMicOn = useAudioStore((state) => state.isMicOn);
  const toggleMic = useAudioStore((state) => state.toggleMic);
  const addToast = useUIStore((state) => state.addToast);

  const [isLoading, setIsLoading] = useState(true);
  const [ws, setWs] = useState<WebSocketService | null>(null);

  useEffect(() => {
    loadRoom();
  }, [roomId]);

  useEffect(() => {
    if (room && token && user) {
      connectWebSocket();
      setupAudio();
    }

    return () => {
      if (ws) {
        ws.disconnect();
      }
      audioService.stopCapture();
    };
  }, [room, token]);

  const loadRoom = async () => {
    try {
      if (!roomId) return;
      const roomData = await ApiService.getRoom(roomId);
      setRoom(roomData);
      setIsLoading(false);
    } catch (error) {
      addToast('Failed to load room', 'error');
      navigate('/');
    }
  };

  const connectWebSocket = async () => {
    const wsService = WebSocketService.getInstance();
    setWs(wsService);

    try {
      await wsService.connect(token!, roomId!);

      wsService.on('room_snapshot', ({ seats, members: snapMembers }: any) => {
        setSnapshot(seats ?? [], snapMembers ?? []);
      });

      wsService.on('seat_occupied', ({ seat_id, user_id, username }: any) => {
        updateSeat(seat_id, user_id, username);
        if (user_id !== user?.id) {
          addToast(`${username || 'Someone'} sat down`, 'info');
        }
      });

      wsService.on('seat_vacated', ({ seat_id }: any) => {
        updateSeat(seat_id, null, null);
      });

      wsService.on('user_joined', ({ user_id, username }: any) => {
        addMember({ user_id, username: username || user_id, is_muted: false });
        if (user_id !== user?.id) {
          addToast(`${username || 'Someone'} joined the room`, 'info');
        }
      });

      wsService.on('user_left', ({ user_id }: any) => {
        audioService.removeRemoteStream(user_id);
        const m = members.get(user_id);
        removeMember(user_id);
        if (user_id !== user?.id) {
          addToast(`${m?.username || 'Someone'} left the room`, 'info');
        }
      });

      wsService.on('audio', ({ user_id, data }: any) => {
        if (user_id !== user?.id && data) {
          audioService.playRemoteAudio(user_id, data);
        }
      });
    } catch (error) {
      addToast('Failed to connect to room', 'error');
      navigate('/');
    }
  };

  const setupAudio = async () => {
    try {
      await audioService.startCapture();
      const wsService = WebSocketService.getInstance();
      audioService.onAudioData = (data) => {
        const { seats } = useRoomStore.getState();
        const isSitting = seats.some((s) => s.occupied_by_id === user?.id);
        if (!isSitting) return;
        wsService.emit('audio', { data });
      };
      addToast('Microphone enabled', 'success');
    } catch (error) {
      addToast('Failed to access microphone', 'error');
    }
  };

  const handleSitDown = async (seatId: string) => {
    try {
      if (roomId) {
        await ApiService.occupySeat(roomId, seatId);
        ws?.emit('sit_down', { seat_id: seatId });
      }
    } catch (error) {
      addToast(error instanceof Error ? error.message : 'Failed to sit down', 'error');
    }
  };

  const handleStandUp = async (seatId: string) => {
    try {
      if (roomId) {
        await ApiService.vacateSeat(roomId, seatId);
        ws?.emit('stand_up', { seat_id: seatId });
      }
    } catch (error) {
      addToast('Failed to stand up', 'error');
    }
  };

  const handleLeaveRoom = async () => {
    ws?.emit('leave_room', {});
    navigate('/', { state: { justLeft: true } });
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <p className="text-xl text-gray-400">Loading room...</p>
      </div>
    );
  }

  if (!room) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <p className="text-xl text-red-400">Room not found</p>
      </div>
    );
  }

  const memberList = Array.from(members.values());

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-7xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-bold">{room.name}</h1>
            <p className="text-gray-400">{room.description}</p>
          </div>
          <div className="flex gap-4">
            <button
              onClick={() => {
                toggleMic();
                audioService.setMicMuted(!isMicOn);
                ws?.emit('mic_toggle', { is_muted: isMicOn });
              }}
              className={`px-6 py-3 rounded font-semibold ${
                isMicOn ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'
              }`}
            >
              {isMicOn ? '🎤 On' : '🎙️ Off'}
            </button>
            <button
              onClick={handleLeaveRoom}
              className="px-6 py-3 bg-red-600 hover:bg-red-700 rounded font-semibold"
            >
              Leave Room
            </button>
          </div>
        </div>

        <div className="space-y-4 mb-8">
          {[1, 2, 3].map((row) => (
            <div key={row} className="flex gap-4">
              {[1, 2, 3, 4, 5, 6].map((col) => {
                const seat = seats.find((s) => s.row === row && s.col === col);
                if (!seat) return null;
                const isOccupied = seat.occupied_by_id !== null;
                const isUserSeat = seat.occupied_by_id === user?.id;
                const occupantName = isUserSeat
                  ? 'You'
                  : seat.username || members.get(seat.occupied_by_id ?? '')?.username || null;

                return (
                  <button
                    key={seat.id}
                    onClick={() => {
                      if (isUserSeat) {
                        handleStandUp(seat.id);
                      } else if (!isOccupied) {
                        handleSitDown(seat.id);
                      }
                    }}
                    disabled={isOccupied && !isUserSeat}
                    className={`flex-1 py-8 rounded text-center font-semibold transition ${
                      isUserSeat
                        ? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
                        : isOccupied
                        ? 'bg-gray-600 opacity-60 cursor-not-allowed'
                        : 'bg-gray-700 hover:bg-gray-600 cursor-pointer'
                    }`}
                  >
                    <div className="text-sm text-gray-300">Seat {row}-{col}</div>
                    {occupantName ? (
                      <div className="mt-1 font-bold truncate px-2">{occupantName}</div>
                    ) : (
                      <div className="mt-1 text-gray-500 text-sm">Empty</div>
                    )}
                  </button>
                );
              })}
            </div>
          ))}
        </div>

        <div className="bg-gray-800 rounded-lg p-8">
          <h2 className="text-2xl font-bold mb-4">Members ({memberList.length})</h2>
          <div className="space-y-2">
            {memberList.map((member) => (
              <div key={member.user_id} className="flex justify-between items-center p-3 bg-gray-700 rounded">
                <span className="font-semibold">
                  {member.user_id === user?.id ? `${member.username} (You)` : member.username}
                </span>
                <span className="text-gray-400">{member.is_muted ? '🔇 Muted' : '🔊 Unmuted'}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
