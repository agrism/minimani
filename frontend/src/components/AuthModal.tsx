import React, { useState } from 'react';
import { X, User, Mail, Lock, Sparkles, CheckCircle2, AlertCircle, ArrowRight } from 'lucide-react';
import { User as UserType } from '../types';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAuthSuccess: (user: UserType, token: string) => void;
  initialMode?: 'login' | 'register';
}

const AVATARS = ['🚗', '🏎️', '🚜', '🏡', '💎', '📱', '🔨', '🍻', '🌲', '👑', '🚀', '🔥'];

export const AuthModal: React.FC<AuthModalProps> = ({
  isOpen,
  onClose,
  onAuthSuccess,
  initialMode = 'register',
}) => {
  const [mode, setMode] = useState<'login' | 'register'>(initialMode);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [nickname, setNickname] = useState('');
  const [avatar, setAvatar] = useState('🚗');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSubmitting(true);

    try {
      const endpoint = mode === 'register' ? '/api/auth/register' : '/api/auth/login';
      const body = mode === 'register'
        ? { email, password, nickname, avatar }
        : { email, password };

      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });

      const data = await res.json();

      if (!res.ok) {
        setError(data.error || 'Radās neparedzēta kļūda');
        setSubmitting(false);
        return;
      }

      // Save token in localStorage
      localStorage.setItem('minimani_auth_token', data.token);
      onAuthSuccess(data.user, data.token);
      onClose();
    } catch (err) {
      console.error('Auth error:', err);
      setError('Neizdevās sazināties ar serveri');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm animate-fade-in">
      <div
        className="relative w-full max-w-md bg-slate-900 border border-slate-750 rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Galvene */}
        <div className="p-5 border-b border-slate-800 flex items-center justify-between bg-slate-850/60">
          <div className="flex items-center gap-2.5">
            <div className="w-9 h-9 rounded-xl bg-amber-500/20 text-amber-400 flex items-center justify-center border border-amber-500/30">
              <User className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-extrabold text-base text-slate-100">
                {mode === 'register' ? 'Izveidot Profilu' : 'Pieslēgties'}
              </h3>
              <p className="text-[11px] text-slate-400">
                {mode === 'register' ? 'Saglabā savu niku un uzvaru statistiku' : 'Laipni lūdzam atpakaļ!'}
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-xl text-slate-400 hover:text-slate-100 hover:bg-slate-800 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Tabu izvēle */}
        <div className="grid grid-cols-2 p-2 bg-slate-950/60 border-b border-slate-800">
          <button
            type="button"
            onClick={() => { setMode('register'); setError(null); }}
            className={`py-2 text-xs font-bold rounded-xl transition-all ${
              mode === 'register'
                ? 'bg-amber-500 text-slate-950 shadow-md shadow-amber-500/20'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Reģistrēties
          </button>
          <button
            type="button"
            onClick={() => { setMode('login'); setError(null); }}
            className={`py-2 text-xs font-bold rounded-xl transition-all ${
              mode === 'login'
                ? 'bg-amber-500 text-slate-950 shadow-md shadow-amber-500/20'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Ienākt
          </button>
        </div>

        {/* Forma */}
        <form onSubmit={handleSubmit} className="p-5 overflow-y-auto space-y-4">
          {error && (
            <div className="p-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs flex items-center gap-2">
              <AlertCircle className="w-4 h-4 shrink-0 text-rose-400" />
              <span>{error}</span>
            </div>
          )}

          {mode === 'register' && (
            <>
              {/* Publiskais segvārds (Niks) */}
              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-300 flex items-center justify-between">
                  <span>Publiskais segvārds (Niks)</span>
                  <span className="text-[10px] text-amber-400 font-normal">Redzams citiem spēlē</span>
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
                    <User className="w-4 h-4" />
                  </div>
                  <input
                    type="text"
                    required
                    value={nickname}
                    onChange={(e) => setNickname(e.target.value)}
                    placeholder="Piem. CenuMeistars"
                    maxLength={30}
                    className="w-full pl-9 pr-3 py-2.5 rounded-xl bg-slate-950 border border-slate-750 text-slate-100 text-sm focus:outline-none focus:border-amber-400"
                  />
                </div>
              </div>

              {/* Avatara izvēle */}
              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-300">Izvēlies savu ikonu</label>
                <div className="grid grid-cols-6 gap-2">
                  {AVATARS.map((av) => (
                    <button
                      key={av}
                      type="button"
                      onClick={() => setAvatar(av)}
                      className={`h-10 text-xl rounded-xl border flex items-center justify-center transition-all ${
                        avatar === av
                          ? 'bg-amber-500/20 border-amber-400 scale-105 shadow'
                          : 'bg-slate-950 border-slate-800 hover:border-slate-700'
                      }`}
                    >
                      {av}
                    </button>
                  ))}
                </div>
              </div>
            </>
          )}

          {/* E-pasts */}
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-300">E-pasts</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
                <Mail className="w-4 h-4" />
              </div>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="tavs.epasts@piemers.lv"
                className="w-full pl-9 pr-3 py-2.5 rounded-xl bg-slate-950 border border-slate-750 text-slate-100 text-sm focus:outline-none focus:border-amber-400"
              />
            </div>
          </div>

          {/* Parole */}
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-300">
              Parole {mode === 'register' && <span className="text-slate-400 font-normal">(vismaz 6 simboli)</span>}
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
                <Lock className="w-4 h-4" />
              </div>
              <input
                type="password"
                required
                minLength={6}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="w-full pl-9 pr-3 py-2.5 rounded-xl bg-slate-950 border border-slate-750 text-slate-100 text-sm focus:outline-none focus:border-amber-400"
              />
            </div>
          </div>

          {/* Galvenā poga */}
          <button
            type="submit"
            disabled={submitting}
            className="w-full py-3 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-black text-sm flex items-center justify-center gap-2 transition-all shadow-lg shadow-amber-500/20 active:scale-98 disabled:opacity-50"
          >
            {submitting ? (
              <div className="w-5 h-5 rounded-full border-2 border-slate-950 border-t-transparent animate-spin" />
            ) : mode === 'register' ? (
              <>
                <Sparkles className="w-4 h-4" />
                Izveidot profilu
              </>
            ) : (
              <>
                <ArrowRight className="w-4 h-4" />
                Ienākt
              </>
            )}
          </button>
        </form>

        {/* Kājene ar viesa atgādinājumu */}
        <div className="p-4 bg-slate-950/60 border-t border-slate-800 text-center space-y-2">
          <p className="text-[11px] text-slate-400">
            Reģistrācija nav obligāta — vari spēlēt arī bez konta.
          </p>
          <button
            type="button"
            onClick={onClose}
            className="text-xs text-amber-400 hover:text-amber-300 font-bold underline"
          >
            Turpināt kā viesim →
          </button>
        </div>
      </div>
    </div>
  );
};
