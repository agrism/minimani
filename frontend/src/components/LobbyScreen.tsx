import React, { useState } from 'react';
import { Copy, Check, Users, Play, Crown, Sparkles, CheckCircle2, LogOut } from 'lucide-react';
import { Room, Player } from '../types';

interface LobbyScreenProps {
  room: Room;
  currentPlayerId: string;
  onVote: (questionCount: number, gameMode: string, category: string) => void;
  onSetReady: (ready: boolean) => void;
  onStartGame: () => void;
  onLeave?: () => void;
}

export const LobbyScreen: React.FC<LobbyScreenProps> = ({
  room,
  currentPlayerId,
  onVote,
  onSetReady,
  onStartGame,
  onLeave,
}) => {
  const [copied, setCopied] = useState(false);

  const playersList: Player[] = Object.values(room.players || {});
  const currentPlayer = room.players[currentPlayerId];
  const isHost = room.host_id === currentPlayerId || playersList[0]?.id === currentPlayerId;
  const totalPlayers = Math.max(1, playersList.length);

  const copyRoomLink = () => {
    const url = window.location.origin + '?room=' + room.code;
    navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  // Calculate vote percentages for Question Count
  const countVotes = room.votes_count || {};
  const countStats: Record<number, number> = { 5: 0, 10: 0, 15: 0, 20: 0 };
  Object.values(countVotes).forEach((v) => {
    if (countStats[v] !== undefined) countStats[v]++;
  });

  // Calculate vote percentages for Game Mode
  const modeVotes = room.votes_mode || {};
  const modeStats: Record<string, number> = { exact: 0, multiple_choice: 0, mix: 0 };
  Object.values(modeVotes).forEach((m) => {
    if (modeStats[m] !== undefined) modeStats[m]++;
  });

  const myCountVote = countVotes[currentPlayerId] || room.question_count;
  const myModeVote = modeVotes[currentPlayerId] || room.game_mode;
  const myCatVote = (room.votes_category && room.votes_category[currentPlayerId]) || room.category;

  return (
    <div className="w-full max-w-md mx-auto px-4 py-5 space-y-5">
      {/* Istabas kods & dalīšanās */}
      <div className="bg-gradient-to-b from-slate-800 to-slate-850 border border-slate-700/80 rounded-2xl p-4 text-center shadow-xl space-y-2">
        <span className="text-xs font-semibold text-slate-400 uppercase tracking-widest block">
          Istabas Kods
        </span>
        <div className="flex items-center justify-center gap-2">
          <span className="font-mono text-4xl font-black tracking-widest text-amber-400">
            {room.code}
          </span>
          <button
            onClick={copyRoomLink}
            className="p-2 rounded-xl bg-slate-700/80 hover:bg-slate-750 text-slate-300 hover:text-amber-400 transition-colors"
            title="Kopēt saiti"
          >
            {copied ? <Check className="w-5 h-5 text-emerald-400" /> : <Copy className="w-5 h-5" />}
          </button>
        </div>
        <p className="text-xs text-slate-400">
          Uzaicini draugus! Atveriet spēli telefonā un ievadiet šo 4-burtu kodu.
        </p>
      </div>

      {/* Dalībnieku saraksts */}
      <div className="bg-slate-800/80 border border-slate-700/80 rounded-2xl p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Users className="w-4 h-4 text-amber-400" />
            <h3 className="text-sm font-bold text-slate-200">
              Dalībnieki ({playersList.length}/{room.max_players})
            </h3>
          </div>
          <span className="text-xs text-slate-400">
            {playersList.filter((p) => p.is_ready).length} gatavi
          </span>
        </div>

        <div className="grid grid-cols-2 gap-2">
          {playersList.map((p, idx) => {
            const isMe = p.id === currentPlayerId;
            return (
              <div
                key={p.id}
                className={`flex items-center gap-2.5 p-2.5 rounded-xl border transition-all ${
                  p.is_ready
                    ? 'bg-emerald-950/20 border-emerald-500/40 text-slate-100'
                    : 'bg-slate-900/60 border-slate-800 text-slate-300'
                }`}
              >
                <div className="w-9 h-9 rounded-xl bg-slate-800 flex items-center justify-center text-xl shrink-0 border border-slate-700">
                  {p.avatar}
                </div>
                <div className="overflow-hidden flex-1 min-w-0">
                  <div className="flex items-center gap-1 font-semibold text-xs truncate">
                    <span className="truncate">{p.nickname}</span>
                    {idx === 0 && <Crown className="w-3.5 h-3.5 text-amber-400 shrink-0" />}
                  </div>
                  <div className="flex items-center gap-1 mt-0.5">
                    {p.is_ready ? (
                      <span className="text-[10px] text-emerald-400 font-bold flex items-center gap-0.5">
                        <CheckCircle2 className="w-2.5 h-2.5" /> Gatavs
                      </span>
                    ) : (
                      <span className="text-[10px] text-slate-500">Gatavojas...</span>
                    )}
                    {isMe && <span className="text-[10px] text-amber-400/80 ml-auto font-mono">(Tu)</span>}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Demokrātiskā balsošana */}
      <div className="bg-slate-800/80 border border-slate-700/80 rounded-2xl p-4 space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1.5">
            <Sparkles className="w-4 h-4 text-amber-400" />
            <h3 className="text-sm font-bold text-slate-200">Balsošana par noteikumiem</h3>
          </div>
          <span className="text-[11px] text-slate-400">&gt;50% balsu apstiprina izvēli</span>
        </div>

        {/* 1. Jautājumu skaits */}
        <div className="space-y-1.5">
          <div className="text-xs font-semibold text-slate-400 flex items-center justify-between">
            <span>JAUTĀJUMU SKAITS:</span>
            <span className="text-amber-400 font-bold">{room.question_count} kārtas izvēlētas</span>
          </div>
          <div className="grid grid-cols-4 gap-2">
            {[5, 10, 15, 20].map((num) => {
              const count = countStats[num] || 0;
              const pct = Math.round((count / totalPlayers) * 100);
              const isSelected = room.question_count === num;
              const isMyVote = myCountVote === num;

              return (
                <button
                  key={num}
                  type="button"
                  onClick={() => onVote(num, myModeVote, myCatVote)}
                  className={`relative overflow-hidden p-2.5 rounded-xl border text-center transition-all ${
                    isMyVote
                      ? 'border-amber-500 bg-amber-500/10 text-amber-300'
                      : isSelected
                      ? 'border-emerald-500/60 bg-emerald-950/20 text-slate-200'
                      : 'border-slate-700 bg-slate-900 text-slate-400 hover:text-slate-200'
                  }`}
                >
                  {/* Procentu fona josla */}
                  <div
                    className="absolute bottom-0 left-0 top-0 bg-amber-500/20 -z-0 transition-all duration-300"
                    style={{ width: `${pct}%` }}
                  />
                  <div className="relative z-10">
                    <div className="text-base font-extrabold">{num}</div>
                    <div className="text-[10px] font-mono opacity-80">{pct}% ({count})</div>
                  </div>
                </button>
              );
            })}
          </div>
        </div>
      </div>

      {/* Apakšējās pogas */}
      <div className="space-y-2">
        <button
          type="button"
          onClick={() => onSetReady(!currentPlayer?.is_ready)}
          className={`w-full py-4 rounded-2xl font-black text-lg transition-all transform active:scale-[0.98] shadow-lg flex items-center justify-center gap-2 ${
            currentPlayer?.is_ready
              ? 'bg-slate-800 text-emerald-400 border border-emerald-500/50 hover:bg-slate-750'
              : 'bg-emerald-500 hover:bg-emerald-400 text-slate-950 shadow-emerald-500/20'
          }`}
        >
          {currentPlayer?.is_ready ? (
            <>
              <Check className="w-6 h-6 stroke-[3]" /> Esmu gatavs! (Atcelt)
            </>
          ) : (
            'Nospied, kad esi gatavs!'
          )}
        </button>

        {isHost && (
          <button
            type="button"
            onClick={onStartGame}
            className="w-full py-3 rounded-xl bg-amber-500/10 hover:bg-amber-500/20 border border-amber-500/40 text-amber-400 font-bold text-sm flex items-center justify-center gap-2 transition-colors"
          >
            <Play className="w-4 h-4 fill-current" />
            Kā administrators sākt spēli tagad
          </button>
        )}

        {onLeave && (
          <button
            type="button"
            onClick={onLeave}
            className="w-full py-2.5 rounded-xl border border-slate-800 bg-slate-900/60 hover:bg-rose-950/20 hover:border-rose-500/40 text-slate-400 hover:text-rose-300 text-xs font-semibold flex items-center justify-center gap-2 transition-all mt-1"
          >
            <LogOut className="w-3.5 h-3.5" />
            Pamest istabu un izvēlēties citu
          </button>
        )}
      </div>
    </div>
  );
};

