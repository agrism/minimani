import React, { useState, useEffect } from 'react';
import { Clock, Send, CheckCircle2, ZoomIn, MapPin } from 'lucide-react';
import { Room } from '../types';

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
  const [guessInput, setGuessInput] = useState<string>('');
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
    };

    updateTimer();
    const interval = setInterval(updateTimer, 250);
    return () => clearInterval(interval);
  }, [room.round_started_at, room.time_limit_sec]);

  // Reset local input on round change
  useEffect(() => {
    setGuessInput('');
    setSelectedOption(null);
  }, [room.current_round]);

  const handleQuickAdd = (amount: number) => {
    const current = parseFloat(guessInput.replace(/\s+/g, '')) || 0;
    setGuessInput(String(current + amount));
  };

  const handleExactSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const price = parseFloat(guessInput.replace(/\s+/g, ''));
    if (!isNaN(price) && price > 0 && !hasGuessed) {
      onSubmitGuess(price);
    }
  };

  const handleOptionSelect = (optionPrice: number) => {
    if (!hasGuessed) {
      setSelectedOption(optionPrice);
      onSubmitGuess(optionPrice);
    }
  };

  const formatPriceInput = (val: string) => {
    const num = val.replace(/\D/g, '');
    if (!num) return '';
    return Number(num).toLocaleString('lv-LV');
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

        {/* Teksts */}
        <div className="p-4 space-y-2">
          <h2 className="font-extrabold text-slate-100 text-base leading-snug">
            {listing.title}
          </h2>
          <p className="text-xs text-slate-300 leading-relaxed line-clamp-3">
            {listing.description}
          </p>
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
      ) : isMultipleChoice ? (
        /* Režīms: 4 Atbilžu varianti */
        <div className="space-y-2.5">
          <span className="text-xs font-semibold text-slate-400 uppercase tracking-wider block text-center">
            Izvēlies pareizo ss.lv cenu:
          </span>
          <div className="grid grid-cols-2 gap-2.5">
            {listing.options?.map((opt, i) => {
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
      ) : (
        /* Režīms: Ierakstīt cenu */
        <form onSubmit={handleExactSubmit} className="space-y-3">
          <div className="relative">
            <input
              type="text"
              inputMode="numeric"
              value={formatPriceInput(guessInput)}
              onChange={(e) => setGuessInput(e.target.value)}
              placeholder="0"
              className="w-full bg-slate-900 border-2 border-slate-700 rounded-2xl px-4 py-3.5 text-center font-mono font-black text-3xl text-amber-400 placeholder-slate-700 focus:outline-none focus:border-amber-500 shadow-inner"
              autoFocus
            />
            <span className="absolute right-5 top-1/2 -translate-y-1/2 text-xl font-bold text-slate-500">
              €
            </span>
          </div>

          {/* Ātrās pogas */}
          <div className="grid grid-cols-4 gap-2">
            {[100, 1000, 5000, 10000].map((amt) => (
              <button
                key={amt}
                type="button"
                onClick={() => handleQuickAdd(amt)}
                className="py-2 rounded-xl bg-slate-800 hover:bg-slate-700 border border-slate-700 text-xs font-bold text-slate-300 transition-colors"
              >
                +{amt >= 1000 ? `${amt / 1000}k` : amt}
              </button>
            ))}
          </div>

          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => setGuessInput('')}
              className="py-3 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-bold text-slate-400"
            >
              Dzēst
            </button>
            <button
              type="submit"
              disabled={!guessInput.trim()}
              className="flex-1 py-3.5 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-400 hover:to-orange-400 text-slate-950 font-black text-base flex items-center justify-center gap-2 shadow-lg shadow-orange-500/20 disabled:opacity-50 transition-all active:scale-[0.98]"
            >
              <Send className="w-4 h-4 stroke-[2.5]" />
              Apstiprināt Cenu
            </button>
          </div>
        </form>
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
            className="max-w-full max-h-[85vh] rounded-2xl object-contain shadow-2xl"
          />
        </div>
      )}
    </div>
  );
};

