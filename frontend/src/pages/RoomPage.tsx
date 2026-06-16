import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useRoomStore } from '@/stores/useRoomStore';
import { useAuthStore } from '@/stores/useAuthStore';
import { useAudioStore } from '@/stores/useAudioStore';
import { useUIStore } from '@/stores/useUIStore';
import { ApiService } from '@/services/api';
import { WebSocketService } from '@/services/websocket';
import { audioService } from '@/services/audio';
import { SeatGridSkeleton, MemberListSkeleton } from '../components/common/Skeleton';
import { AudioMeter } from '../components/room/AudioMeter';
import { PresenceIndicator } from '../components/room/PresenceIndicator';

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
  const localStream = useAudioStore((state) => state.localStream);
  const setLocalStream = useAudioStore((state) => state.setLocalStream);
  const addToast = useUIStore((state) => state.addToast);

  const [isLoadingRoom, setIsLoadingRoom] = useState(true);
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
      setLocalStream(null);
    };
  }, [room, token]);

  const loadRoom = async () => {
    try {
      if (!roomId) return;
      const roomData = await ApiService.getRoom(roomId);
      setRoom(roomData);
      setIsLoadingRoom(false);
    } catch (error) {
      addToast('Failed to load room', 'error');
      navigate('/');
    }
  };

  const connectWebSocket = async () => {
    const wsService = WebSocketService.getInstance();
    setWs(wsService);

    wsService.setCallbacks({
      onReconnecting: (attempt, max) => {
        addToast(`Reconnecting... (attempt ${attempt}/${max})`, 'warning');
      },
      onReconnected: () => {
        addToast('Reconnected!', 'success');
      },
      onReconnectFailed: () => {
        addToast('Connection lost. Please refresh.', 'error');
      },
    });

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
      const stream = await audioService.startCapture();
      setLocalStream(stream);
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

  if (!room && !isLoadingRoom) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <p className="text-xl text-red-400">Room not found</p>
      </div>
    );
  }

  const memberList = Array.from(members.values());

  return (
    <div className="min-h-screen bg-gray-900 pb-24 sm:pb-8">
      {/* Header */}
      <div className="bg-gray-800 border-b border-gray-700 px-4 sm:px-8 py-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center gap-4">
          <div className="min-w-0">
            {isLoadingRoom ? (
              <div className="space-y-1">
                <div className="h-7 w-48 bg-gray-700 rounded animate-pulse" />
                <div className="h-4 w-32 bg-gray-700 rounded animate-pulse" />
              </div>
            ) : (
              <>
                <h1 className="text-2xl sm:text-4xl font-bold truncate">{room?.name}</h1>
                <p className="text-gray-400 text-sm truncate">{room?.description}</p>
              </>
            )}
          </div>
          {/* Desktop controls */}
          <div className="hidden sm:flex gap-4 items-center shrink-0">
            {localStream && (
              <AudioMeter stream={localStream} className="self-center" />
            )}
            <button
              onClick={() => {
                toggleMic();
                audioService.setMicMuted(!isMicOn);
                ws?.emit('mic_toggle', { is_muted: isMicOn });
              }}
              className={`px-6 py-3 rounded font-semibold transition-colors ${
                isMicOn ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'
              }`}
            >
              {isMicOn ? '🎤 On' : '🎙️ Off'}
            </button>
            <button
              onClick={handleLeaveRoom}
              className="px-6 py-3 bg-red-600 hover:bg-red-700 rounded font-semibold transition-colors"
            >
              Leave Room
            </button>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-8 py-6">
        <div className="flex flex-col lg:flex-row gap-6">
          {/* Seat grid */}
          <div className="flex-1">
            <h2 className="text-lg font-semibold text-gray-300 mb-4">Seats</h2>
            {isLoadingRoom ? (
              <SeatGridSkeleton />
            ) : (
              <div className="space-y-3">
                {[1, 2, 3].map((row) => (
                  <div key={row} className="grid grid-cols-3 sm:grid-cols-6 gap-2 sm:gap-3">
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
                          className={`py-6 sm:py-8 rounded text-center font-semibold transition ${
                            isUserSeat
                              ? 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
                              : isOccupied
                              ? 'bg-gray-600 opacity-60 cursor-not-allowed'
                              : 'bg-gray-700 hover:bg-gray-600 cursor-pointer'
                          }`}
                        >
                          <div className="text-xs sm:text-sm text-gray-300">Seat {row}-{col}</div>
                          {occupantName ? (
                            <div className="mt-1 font-bold truncate px-1 text-sm">{occupantName}</div>
                          ) : (
                            <div className="mt-1 text-gray-500 text-xs sm:text-sm">Empty</div>
                          )}
                        </button>
                      );
                    })}
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Members panel */}
          <div className="lg:w-72 bg-gray-800 rounded-lg p-4 sm:p-6 self-start">
            <h2 className="text-xl font-bold mb-4">Members ({memberList.length})</h2>
            {isLoadingRoom ? (
              <MemberListSkeleton />
            ) : (
              <div className="space-y-2">
                {memberList.map((member) => (
                  <div key={member.user_id} className="flex justify-between items-center p-3 bg-gray-700 rounded gap-2">
                    <div className="flex items-center gap-2 min-w-0">
                      <PresenceIndicator isOnline={true} size="sm" />
                      <span className="font-semibold truncate text-sm">
                        {member.user_id === user?.id ? `${member.username} (You)` : member.username}
                      </span>
                    </div>
                    <span className="text-gray-400 text-sm shrink-0">
                      {member.is_muted ? '🔇' : '🔊'}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Mobile fixed bottom controls */}
      <div className="sm:hidden fixed bottom-0 left-0 right-0 bg-gray-800 border-t border-gray-700 px-4 py-3 flex gap-3 items-center z-40">
        {localStream && (
          <AudioMeter stream={localStream} className="flex-shrink-0" />
        )}
        <button
          onClick={() => {
            toggleMic();
            audioService.setMicMuted(!isMicOn);
            ws?.emit('mic_toggle', { is_muted: isMicOn });
          }}
          className={`flex-1 py-2.5 rounded font-semibold text-sm transition-colors ${
            isMicOn ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'
          }`}
        >
          {isMicOn ? '🎤 On' : '🎙️ Off'}
        </button>
        <button
          onClick={handleLeaveRoom}
          className="flex-1 py-2.5 bg-red-600 hover:bg-red-700 rounded font-semibold text-sm transition-colors"
        >
          Leave Room
        </button>
      </div>
    </div>
  );
}
