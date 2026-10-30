import React, { useEffect } from 'react';
import confetti from 'canvas-confetti';
import { Trophy, RotateCcw, Home } from 'lucide-react';
import { Room } from '../types';

interface PodiumScreenProps {
  room: Room;
  currentPlayerId: string;
  onPlayAgain: () => void;
  onGoHome: () => void;
}

export const PodiumScreen: React.FC<PodiumScreenProps> = ({
  room,
  currentPlayerId,
  onPlayAgain,
  onGoHome,
}) => {
  const ranked = Object.values(room.players || {}).sort((a, b) => b.score - a.score);

  const first = ranked[0];
  const second = ranked[1];
  const third = ranked[2];

  useEffect(() => {
    // Launch celebratory confetti
    confetti({
      particleCount: 100,
      spread: 70,
      origin: { y: 0.6 },
    });

    const timer = setTimeout(() => {
      confetti({
        particleCount: 60,
        angle: 60,
        spread: 55,
        origin: { x: 0 },
      });
      confetti({
        particleCount: 60,
        angle: 120,
        spread: 55,
        origin: { x: 1 },
      });
    }, 400);

    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="w-full max-w-md mx-auto px-4 py-6 space-y-6">
      <div className="text-center space-y-1">
        <span className="text-xs font-bold text-amber-400 uppercase tracking-widest">
          Spēle Noslēgusies
        </span>
        <h2 className="text-3xl font-black text-slate-100 flex items-center justify-center gap-2">
          <Trophy className="w-8 h-8 text-amber-400" />
          Uzvarētāju Pjedestāls
        </h2>
      </div>

      {/* Pjedestāls (1., 2., 3. vieta) */}
      <div className="flex items-end justify-center gap-2 pt-8 pb-4">
        {/* 2. Vieta */}
        {second && (
          <div className="flex-1 text-center space-y-2">
            <div className="text-3xl mb-1">{second.avatar}</div>
            <div className="text-xs font-bold text-slate-200 truncate">{second.nickname}</div>
            <div className="h-28 rounded-t-2xl bg-gradient-to-t from-slate-700 to-slate-600 flex flex-col items-center justify-center text-slate-200 border-t-2 border-slate-400 shadow-lg">
              <span className="text-2xl font-black font-mono">2</span>
              <span className="text-[11px] font-mono text-slate-300 font-bold">{second.score} pt</span>
            </div>
          </div>
        )}

        {/* 1. Vieta */}
        {first && (
          <div className="flex-1 text-center space-y-2 -mt-6">
            <div className="relative inline-block">
              <div className="text-4xl mb-1">{first.avatar}</div>
              <span className="absolute -top-3.5 -right-2 text-xl">👑</span>
            </div>
            <div className="text-xs font-black text-amber-400 truncate">{first.nickname}</div>
            <div className="h-36 rounded-t-2xl bg-gradient-to-t from-amber-600 to-amber-500 flex flex-col items-center justify-center text-slate-950 border-t-2 border-amber-300 shadow-2xl shadow-amber-500/30">
              <span className="text-3xl font-black font-mono">1</span>
              <span className="text-xs font-mono font-black">{first.score} pt</span>
            </div>
          </div>
        )}

        {/* 3. Vieta */}
        {third && (
          <div className="flex-1 text-center space-y-2">
            <div className="text-3xl mb-1">{third.avatar}</div>
            <div className="text-xs font-bold text-slate-200 truncate">{third.nickname}</div>
            <div className="h-20 rounded-t-2xl bg-gradient-to-t from-amber-900/60 to-amber-800/60 flex flex-col items-center justify-center text-amber-300 border-t-2 border-amber-700 shadow-md">
              <span className="text-xl font-black font-mono">3</span>
              <span className="text-[10px] font-mono text-amber-200/80 font-bold">{third.score} pt</span>
            </div>
          </div>
        )}
      </div>

      {/* Pilna tabula */}
      <div className="bg-slate-800/80 border border-slate-700/80 rounded-2xl p-4 space-y-2">
        <h4 className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">
          Visi dalībnieki
        </h4>
        <div className="space-y-1.5">
          {ranked.map((p, rank) => {
            const isMe = p.id === currentPlayerId;
            return (
              <div
                key={p.id}
                className={`p-2.5 rounded-xl border flex items-center justify-between text-xs font-semibold ${
                  isMe ? 'bg-amber-500/10 border-amber-500/40 text-amber-300' : 'bg-slate-900/60 border-slate-800 text-slate-200'
                }`}
              >
                <div className="flex items-center gap-2.5">
                  <span className="font-mono font-bold w-5 text-slate-500">{rank + 1}.</span>
                  <span className="text-base">{p.avatar}</span>
                  <span>{p.nickname}</span>
                  {isMe && <span className="text-[10px] text-amber-400 font-mono">(Tu)</span>}
                </div>
                <span className="font-mono font-bold text-amber-400">{p.score} punkti</span>
              </div>
            );
          })}
        </div>
      </div>

      {/* Pogas */}
      <div className="space-y-2.5">
        <button
          onClick={onPlayAgain}
          className="w-full py-4 rounded-2xl bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-400 hover:to-orange-400 text-slate-950 font-black text-base flex items-center justify-center gap-2 shadow-lg shadow-orange-500/25 transition-all transform active:scale-[0.98]"
        >
          <RotateCcw className="w-5 h-5 stroke-[2.5]" />
          Spēlēt Vēlreiz
        </button>
        <button
          onClick={onGoHome}
          className="w-full py-3 rounded-xl bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-300 font-bold text-sm flex items-center justify-center gap-2 transition-colors"
        >
          <Home className="w-4 h-4" />
          Doties uz sākumlapu
        </button>
      </div>
    </div>
  );
};

