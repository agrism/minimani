import React, { useState, useEffect, useRef } from 'react';
import { Header } from './components/Header';
import { HomeScreen } from './components/HomeScreen';
import { LobbyScreen } from './components/LobbyScreen';
import { GameScreen } from './components/GameScreen';
import { RevealScreen } from './components/RevealScreen';
import { PodiumScreen } from './components/PodiumScreen';
import { Room } from './types';

export const App: React.FC = () => {
  const [playerId] = useState<string>(() => {
    let id = localStorage.getItem('minimani_player_id');
    if (!id) {
      id = 'p_' + Math.random().toString(36).substring(2, 9);
      localStorage.setItem('minimani_player_id', id);
    }
    return id;
  });

  const [nickname, setNickname] = useState<string>(() => {
    return localStorage.getItem('minimani_nickname') || `Minētājs #${Math.floor(Math.random() * 900) + 100}`;
  });

  const [avatar, setAvatar] = useState<string>(() => {
    return localStorage.getItem('minimani_avatar') || '🚗';
  });

  const [room, setRoom] = useState<Room | null>(null);
  const [isStarting, setIsStarting] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  // Update profile
  const handleUpdateProfile = (name: string, av: string) => {
    setNickname(name);
    setAvatar(av);
    localStorage.setItem('minimani_nickname', name);
    localStorage.setItem('minimani_avatar', av);

    // If currently in room, broadcast update
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          action: 'join',
          player_id: playerId,
          nickname: name,
          avatar: av,
        })
      );
    }
  };

  // Connect to room WebSocket
  const connectToRoom = (roomCode: string) => {
    if (socketRef.current) {
      socketRef.current.close();
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws/rooms/${roomCode}`;

    const ws = new WebSocket(wsUrl);
    socketRef.current = ws;

    ws.onopen = () => {
      // Send join message
      ws.send(
        JSON.stringify({
          action: 'join',
          player_id: playerId,
          nickname: nickname,
          avatar: avatar,
        })
      );
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'room_state' && data.room) {
          setRoom(data.room);
          if (data.room.status === 'starting') {
            setIsStarting(true);
          } else {
            setIsStarting(false);
          }
        }
      } catch (err) {
        console.error('Failed to parse WS message:', err);
      }
    };

    ws.onclose = () => {
      console.log('WS disconnected');
    };
  };

  // Auto-join from URL param
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const roomParam = params.get('room');
    if (roomParam) {
      connectToRoom(roomParam.toUpperCase());
    }
  }, []);

  // API Actions
  const handleCreateRoom = async (isPublic: boolean, category: string, gameMode: string) => {
    try {
      const res = await fetch('/api/rooms/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          host_id: playerId,
          is_public: isPublic,
          category: category,
          game_mode: gameMode,
        }),
      });
      if (res.ok) {
        const data = await res.json();
        connectToRoom(data.code);
      }
    } catch (e) {
      console.error('Failed to create room:', e);
    }
  };

  const handleQuickJoin = async () => {
    try {
      const res = await fetch('/api/rooms/quick-join', {
        method: 'POST',
      });
      if (res.ok) {
        const data = await res.json();
        connectToRoom(data.code);
      }
    } catch (e) {
      console.error('Failed to quick join:', e);
    }
  };

  const handleJoinRoom = (code: string) => {
    connectToRoom(code);
  };

  const handleLeaveRoom = () => {
    if (socketRef.current) {
      socketRef.current.close();
      socketRef.current = null;
    }
    setRoom(null);
    window.history.replaceState({}, '', window.location.pathname);
  };

  // WS In-Game Actions
  const handleVote = (questionCount: number, gameMode: string, category: string) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          action: 'vote',
          player_id: playerId,
          question_count: questionCount,
          game_mode: gameMode,
          category: category,
        })
      );
    }
  };

  const handleSetReady = (ready: boolean) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          action: 'ready',
          player_id: playerId,
          is_ready: ready,
        })
      );
    }
  };

  const handleStartGame = () => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          action: 'start',
        })
      );
    }
  };

  const handleSubmitGuess = (price: number) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          action: 'guess',
          player_id: playerId,
          guess_price: price,
        })
      );
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
      <Header
        roomCode={room?.code}
        playerCount={room ? Object.keys(room.players || {}).length : undefined}
        onLeave={room ? handleLeaveRoom : undefined}
      />

      <main className="flex-1 flex flex-col justify-center">
        {/* Sākuma ekrāns */}
        {!room && (
          <HomeScreen
            nickname={nickname}
            avatar={avatar}
            onUpdateProfile={handleUpdateProfile}
            onCreateRoom={handleCreateRoom}
            onJoinRoom={handleJoinRoom}
            onQuickJoin={handleQuickJoin}
          />
        )}

        {/* 3 sekunžu starta atskaite */}
        {room && isStarting && (
          <div className="flex flex-col items-center justify-center p-8 space-y-4 text-center">
            <span className="text-sm font-bold text-amber-400 uppercase tracking-widest animate-pulse">
              Visi ir gatavi!
            </span>
            <h1 className="text-5xl font-black text-slate-100 animate-bounce">
              Spēle Tūlīt Sāksies!
            </h1>
            <p className="text-slate-400 text-sm">
              Sagatavojies minēt ss.lv sludinājumu cenas...
            </p>
          </div>
        )}

        {/* Vestibils (Lobby) */}
        {room && room.status === 'lobby' && !isStarting && (
          <LobbyScreen
            room={room}
            currentPlayerId={playerId}
            onVote={handleVote}
            onSetReady={handleSetReady}
            onStartGame={handleStartGame}
          />
        )}

        {/* Spēles ekrāns */}
        {room && room.status === 'round_active' && (
          <GameScreen
            room={room}
            currentPlayerId={playerId}
            onSubmitGuess={handleSubmitGuess}
          />
        )}

        {/* Atklāšanas ekrāns */}
        {room && room.status === 'round_reveal' && (
          <RevealScreen
            room={room}
            currentPlayerId={playerId}
          />
        )}

        {/* Pjedestāls */}
        {room && room.status === 'game_over' && (
          <PodiumScreen
            room={room}
            currentPlayerId={playerId}
            onPlayAgain={() => handleVote(room.question_count, room.game_mode, room.category)}
            onGoHome={handleLeaveRoom}
          />
        )}
      </main>

      <footer className="py-3 text-center text-[11px] text-slate-600 border-t border-slate-900">
        MiniCenu spēle · Sludinājumu dati no ss.lv · Ražots Latvijā
      </footer>
    </div>
  );
};

