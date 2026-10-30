import React, { useState, useEffect } from 'react';
import { Clock, CheckCircle2, ZoomIn, MapPin } from 'lucide-react';
import { Room } from '../types';
import { playTick } from '../utils/audio';

interface GameScreenProps {
  room: Room;
  currentPlayerId: string;
  onSubmitGuess: (price: number) => void;
}

export const GameScreen: React.FC<GameScreenProps> = ({
  room,
  currentPlayerId,
  onSubmitGuess,
}) => {
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [isZoomed, setIsZoomed] = useState(false);
  const [timeLeft, setTimeLeft] = useState(room.time_limit_sec || 45);

  const listing = room.current_listing;
  const currentPlayer = room.players[currentPlayerId];
  const hasGuessed = currentPlayer?.has_guess || false;

  const totalPlayers = Object.keys(room.players || {}).length;
  const guessedPlayersCount = Object.values(room.players || {}).filter((p) => p.has_guess).length;

  // Local timer synced with round_started_at
  useEffect(() => {
    const updateTimer = () => {
      if (!room.round_started_at) return;
      const started = new Date(room.round_started_at).getTime();
      const elapsed = Math.floor((Date.now() - started) / 1000);
      const remaining = Math.max(0, (room.time_limit_sec || 45) - elapsed);
      setTimeLeft(remaining);
      if (remaining <= 5 && remaining > 0 && remaining !== timeLeft) {
        playTick();
      }
    };

    updateTimer();
    const interval = setInterval(updateTimer, 250);
    return () => clearInterval(interval);
  }, [room.round_started_at, room.time_limit_sec]);

  // Reset local input on round change
  useEffect(() => {
    setSelectedOption(null);
  }, [room.current_round]);

  const handleOptionSelect = (optionPrice: number) => {
    if (!hasGuessed) {
      setSelectedOption(optionPrice);
      onSubmitGuess(optionPrice);
    }
  };

  const timerPercent = Math.max(0, Math.min(100, (timeLeft / (room.time_limit_sec || 45)) * 100));
  const isUrgent = timeLeft <= 10;

  if (!listing) {
    return (
      <div className="flex items-center justify-center min-h-[50vh] text-slate-400">
        Ielādējam nākamo ss.lv sludinājumu...
      </div>
    );
  }

  const isMultipleChoice = listing.options && listing.options.length > 0;

  return (
    <div className="w-full max-w-md mx-auto px-4 py-4 space-y-4">
      {/* Taimeris un Raunds */}
      <div className="space-y-1.5">
        <div className="flex items-center justify-between text-xs font-bold text-slate-300">
          <span className="uppercase tracking-wider">
            Raunds {room.current_round} / {room.question_count}
          </span>
          <span className={`flex items-center gap-1 font-mono text-sm ${isUrgent ? 'text-rose-400 animate-pulse' : 'text-amber-400'}`}>
            <Clock className="w-4 h-4" />
            {timeLeft}s
          </span>
        </div>
        <div className="w-full h-2 rounded-full bg-slate-800 overflow-hidden">
          <div
            className={`h-full transition-all duration-300 rounded-full ${
              isUrgent ? 'bg-rose-500' : timeLeft <= 20 ? 'bg-amber-400' : 'bg-emerald-400'
            }`}
            style={{ width: `${timerPercent}%` }}
          />
        </div>
      </div>

      {/* Sludinājuma kartīte */}
      <div className="bg-slate-800/90 border border-slate-700/80 rounded-2xl overflow-hidden shadow-xl">
        {/* Attēls */}
        <div className="relative aspect-[4/3] bg-slate-950 overflow-hidden group cursor-pointer" onClick={() => setIsZoomed(true)}>
          <img
            src={listing.image_path}
            alt={listing.title}
            referrerPolicy="no-referrer"
            className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
          <div className="absolute top-3 left-3 flex gap-2">
            <span className="px-2.5 py-1 rounded-lg bg-slate-950/80 backdrop-blur text-[11px] font-bold text-amber-400 border border-amber-500/30">
              ss.lv
            </span>
            {listing.region && (
              <span className="px-2.5 py-1 rounded-lg bg-slate-950/80 backdrop-blur text-[11px] font-semibold text-slate-200 border border-slate-700/60 flex items-center gap-1">
                <MapPin className="w-3 h-3 text-emerald-400" />
                {listing.region}
              </span>
            )}
          </div>
          <button
            type="button"
            className="absolute bottom-3 right-3 p-2 rounded-xl bg-slate-950/80 backdrop-blur text-slate-200 border border-slate-700 hover:text-amber-400"
            title="Palielināt attēlu"
          >
            <ZoomIn className="w-4 h-4" />
          </button>
        </div>

        {/* Teksts un Parametri */}
        <div className="p-4 space-y-3">
          <h2 className="font-extrabold text-slate-100 text-sm sm:text-base leading-snug">
            {listing.title}
          </h2>

          {/* Sludinājuma parametru tabula no ss.lv */}
          {listing.parameters && listing.parameters.length > 0 ? (
            <div className="bg-slate-900/90 rounded-xl p-3 border border-slate-750 divide-y divide-slate-800 shadow-inner">
              {listing.parameters.map((param, idx) => (
                <div key={idx} className="py-1.5 flex items-start justify-between text-xs gap-3">
                  <span className="text-slate-400 font-medium shrink-0">{param.label}:</span>
                  <span className="text-slate-100 font-bold text-right">{param.value}</span>
                </div>
              ))}
            </div>
          ) : (
            listing.description && (
              <p className="text-xs text-slate-300 leading-relaxed line-clamp-3">
                {listing.description}
              </p>
            )
          )}

          {listing.parameters && listing.parameters.length > 0 && listing.description && (
            <p className="text-[11px] text-slate-400 leading-relaxed line-clamp-2 italic">
              "{listing.description}"
            </p>
          )}
        </div>
      </div>

      {/* Minēšanas lauks */}
      {hasGuessed ? (
        <div className="p-5 rounded-2xl bg-emerald-950/20 border border-emerald-500/40 text-center space-y-2">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-emerald-500/20 text-emerald-400">
            <CheckCircle2 className="w-7 h-7" />
          </div>
          <h4 className="font-bold text-slate-100 text-base">Tavs minējums ir pieņemts!</h4>
          <p className="text-xs text-slate-400">
            Gaidām pārējos spēlētājus ({guessedPlayersCount}/{totalPlayers})...
          </p>
          <div className="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mt-3">
            <div
              className="bg-emerald-400 h-full transition-all duration-300"
              style={{ width: `${(guessedPlayersCount / totalPlayers) * 100}%` }}
            />
          </div>
        </div>
      ) : (
        /* Režīms: Vienmēr 4 Atbilžu varianti */
        <div className="space-y-2.5">
          <span className="text-xs font-semibold text-slate-400 uppercase tracking-wider block text-center">
            Izvēlies pareizo ss.lv cenu:
          </span>
          <div className="grid grid-cols-2 gap-2.5">
            {(listing.options && listing.options.length >= 4 ? listing.options : [500, 1500, 4500, 10000]).map((opt, i) => {
              const colors = [
                'bg-red-500/10 hover:bg-red-500/20 border-red-500/50 text-red-300',
                'bg-blue-500/10 hover:bg-blue-500/20 border-blue-500/50 text-blue-300',
                'bg-amber-500/10 hover:bg-amber-500/20 border-amber-500/50 text-amber-300',
                'bg-emerald-500/10 hover:bg-emerald-500/20 border-emerald-500/50 text-emerald-300',
              ];
              const letters = ['A', 'B', 'C', 'D'];

              return (
                <button
                  key={opt}
                  type="button"
                  onClick={() => handleOptionSelect(opt)}
                  className={`p-4 rounded-2xl border text-center font-black text-xl transition-all transform active:scale-95 shadow-md flex items-center justify-center gap-2 ${colors[i % 4]}`}
                >
                  <span className="text-xs opacity-60 font-mono">{letters[i]}:</span>
                  <span>{opt.toLocaleString('lv-LV')} €</span>
                </button>
              );
            })}
          </div>
        </div>
      )}

      {/* Palielinātais attēls (Lightbox) */}
      {isZoomed && (
        <div
          className="fixed inset-0 z-50 bg-slate-950/90 backdrop-blur-md flex items-center justify-center p-4 cursor-pointer"
          onClick={() => setIsZoomed(false)}
        >
          <img
            src={listing.image_path}
            alt={listing.title}
            referrerPolicy="no-referrer"
            className="max-w-full max-h-[85vh] rounded-2xl object-contain shadow-2xl"
          />
        </div>
      )}
    </div>
  );
};

