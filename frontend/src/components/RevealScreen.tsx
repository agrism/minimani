import React from 'react';
import { Sparkles, Trophy, Target, ArrowUp, ArrowDown } from 'lucide-react';
import { Room, PlayerGuessResult } from '../types';

interface RevealScreenProps {
  room: Room;
  currentPlayerId: string;
}

export const RevealScreen: React.FC<RevealScreenProps> = ({
  room,
  currentPlayerId,
}) => {
  const listing = room.current_listing;
  const actualPrice = listing?.actual_price || 0;
  const guessesMap = room.round_guesses || {};
  const guessesList: PlayerGuessResult[] = Object.values(guessesMap);

  // Sort by points descending
  guessesList.sort((a, b) => b.score_points - a.score_points);

  // Cumulative player rankings
  const rankedPlayers = Object.values(room.players || {}).sort((a, b) => b.score - a.score);

  return (
    <div className="w-full max-w-md mx-auto px-4 py-4 space-y-4">
      {/* Patiesās cenas atklāšana */}
      <div className="bg-gradient-to-br from-amber-500 to-orange-500 rounded-3xl p-6 text-center text-slate-950 shadow-2xl shadow-orange-500/30 transform animate-bounce-short">
        <span className="text-xs font-black tracking-widest uppercase opacity-80 block mb-1">
          ss.lv Patiesā Cena
        </span>
        <div className="text-4xl font-black font-mono tracking-tight">
          {actualPrice.toLocaleString('lv-LV')} €
        </div>
        <p className="text-xs font-bold mt-1 text-slate-900/80">
          {listing?.title}
        </p>
      </div>

      {/* Šī raunda minējumi */}
      <div className="bg-slate-800/80 border border-slate-700/80 rounded-2xl p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1.5">
            <Target className="w-4 h-4 text-amber-400" />
            <h3 className="text-sm font-bold text-slate-200">Raunda minējumi</h3>
          </div>
          <span className="text-xs text-slate-400">Raunds {room.current_round}</span>
        </div>

        <div className="space-y-2">
          {guessesList.map((g, idx) => {
            const isMe = g.player_id === currentPlayerId;
            const diff = g.guess_value - actualPrice;

            return (
              <div
                key={g.player_id}
                className={`p-3 rounded-xl border flex items-center justify-between transition-all ${
                  isMe
                    ? 'bg-amber-500/10 border-amber-500/50'
                    : 'bg-slate-900/60 border-slate-750'
                }`}
              >
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center text-lg border border-slate-750">
                    {g.avatar}
                  </div>
                  <div>
                    <div className="flex items-center gap-1.5 font-bold text-xs text-slate-100">
                      <span>{g.nickname}</span>
                      {isMe && <span className="text-amber-400 text-[10px] font-mono">(Tu)</span>}
                      {g.is_bullseye && (
                        <span className="text-[10px] px-1.5 py-0.2 rounded bg-amber-400 text-slate-950 font-black flex items-center gap-0.5">
                          <Sparkles className="w-2.5 h-2.5" /> Bullseye!
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-2 text-xs font-mono text-slate-400 mt-0.5">
                      <span>{g.guess_value.toLocaleString('lv-LV')} €</span>
                      <span className={`flex items-center text-[10px] ${diff > 0 ? 'text-rose-400' : diff < 0 ? 'text-sky-400' : 'text-emerald-400'}`}>
                        {diff > 0 ? <ArrowUp className="w-3 h-3" /> : <ArrowDown className="w-3 h-3" />}
                        {Math.abs(diff).toLocaleString('lv-LV')} €
                      </span>
                    </div>
                  </div>
                </div>

                <div className="text-right">
                  <div className="font-extrabold text-sm text-emerald-400">
                    +{g.score_points} pt
                  </div>
                  <div className="text-[10px] text-slate-400 font-mono">
                    {g.accuracy_pct}% precizitāte
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Kopējā rangu tabula */}
      <div className="bg-slate-800/80 border border-slate-700/80 rounded-2xl p-4 space-y-2.5">
        <div className="flex items-center gap-1.5">
          <Trophy className="w-4 h-4 text-amber-400" />
          <h3 className="text-sm font-bold text-slate-200">Kopvērtējuma līderi</h3>
        </div>

        <div className="space-y-1.5">
          {rankedPlayers.slice(0, 5).map((p, rank) => (
            <div
              key={p.id}
              className="flex items-center justify-between px-3 py-2 rounded-lg bg-slate-900/60 text-xs font-semibold"
            >
              <div className="flex items-center gap-2">
                <span className={`font-mono font-bold w-4 ${rank === 0 ? 'text-amber-400' : rank === 1 ? 'text-slate-300' : rank === 2 ? 'text-amber-600' : 'text-slate-500'}`}>
                  {rank + 1}.
                </span>
                <span>{p.avatar}</span>
                <span className="text-slate-200">{p.nickname}</span>
              </div>
              <span className="font-mono text-amber-400 font-bold">{p.score} pt</span>
            </div>
          ))}
        </div>
      </div>

      <div className="text-center text-xs text-slate-500 font-medium animate-pulse py-1">
        Nākamais ss.lv sludinājums tūlīt startēs...
      </div>
    </div>
  );
};

