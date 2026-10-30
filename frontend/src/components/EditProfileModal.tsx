import React, { useState, useEffect } from 'react';
import { X, User, Lock, Mail, Sparkles, CheckCircle2, AlertCircle, Trophy, Gamepad2, Award } from 'lucide-react';
import { User as UserType } from '../types';

interface EditProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
  currentUser: UserType | null;
  onProfileUpdated: (updatedUser: UserType) => void;
}

const AVATARS = ['🚗', '🏎️', '🚜', '🏡', '💎', '📱', '🔨', '🍻', '🌲', '👑', '🚀', '🔥'];

export const EditProfileModal: React.FC<EditProfileModalProps> = ({
  isOpen,
  onClose,
  currentUser,
  onProfileUpdated,
}) => {
  const [nickname, setNickname] = useState('');
  const [avatar, setAvatar] = useState('🚗');
  const [newPassword, setNewPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (currentUser) {
      setNickname(currentUser.nickname || '');
      setAvatar(currentUser.avatar || '🚗');
      setNewPassword('');
      setError(null);
      setSuccess(false);
    }
  }, [currentUser, isOpen]);

  if (!isOpen || !currentUser) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    setSubmitting(true);

    const token = localStorage.getItem('minimani_auth_token');
    if (!token) {
      setError('Nav autorizācijas sesijas');
      setSubmitting(false);
      return;
    }

    try {
      const res = await fetch('/api/auth/profile', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          nickname,
          avatar,
          new_password: newPassword.trim() || undefined,
        }),
      });

      const data = await res.json();

      if (!res.ok) {
        setError(data.error || 'Neizdevās saglabāt izmaiņas');
        setSubmitting(false);
        return;
      }

      setSuccess(true);
      onProfileUpdated(data.user);
      setTimeout(() => {
        onClose();
      }, 800);
    } catch (err) {
      console.error('Update profile error:', err);
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
            <div className="w-10 h-10 rounded-xl bg-amber-500/20 text-amber-400 flex items-center justify-center border border-amber-500/30 text-2xl">
              {avatar}
            </div>
            <div>
              <h3 className="font-extrabold text-base text-slate-100">Rediģēt Profilu</h3>
              <p className="text-[11px] text-slate-400">{currentUser.email}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-xl text-slate-400 hover:text-slate-100 hover:bg-slate-800 transition-colors"
          >
            <X className="w-5 h-5" />
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

          {success && (
            <div className="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-xs flex items-center gap-2">
              <CheckCircle2 className="w-4 h-4 shrink-0 text-emerald-400" />
              <span>Profils veiksmīgi atjaunināts!</span>
            </div>
          )}

          {/* Spēlētāja statistika */}
          <div className="grid grid-cols-3 gap-2 p-3 bg-slate-950/60 border border-slate-800 rounded-2xl text-center">
            <div>
              <div className="text-[10px] text-slate-400 flex items-center justify-center gap-1 font-medium">
                <Gamepad2 className="w-3 h-3 text-slate-500" /> Spēles
              </div>
              <div className="text-base font-black text-slate-200 mt-0.5">{currentUser.games_played || 0}</div>
            </div>
            <div>
              <div className="text-[10px] text-slate-400 flex items-center justify-center gap-1 font-medium">
                <Trophy className="w-3 h-3 text-amber-400" /> Uzvaras
              </div>
              <div className="text-base font-black text-amber-400 mt-0.5">{currentUser.games_won || 0}</div>
            </div>
            <div>
              <div className="text-[10px] text-slate-400 flex items-center justify-center gap-1 font-medium">
                <Award className="w-3 h-3 text-emerald-400" /> Punkti
              </div>
              <div className="text-base font-black text-emerald-400 mt-0.5">{currentUser.total_score || 0}</div>
            </div>
          </div>

          {/* Publiskais segvārds */}
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-300 flex items-center justify-between">
              <span>Publiskais segvārds (Niks)</span>
              <span className="text-[10px] text-slate-500 font-normal">2-30 simboli</span>
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
                placeholder="Tavs spēlētāja niks"
                maxLength={30}
                className="w-full pl-9 pr-3 py-2.5 rounded-xl bg-slate-950 border border-slate-750 text-slate-100 text-sm focus:outline-none focus:border-amber-400 font-bold"
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
                      ? 'bg-amber-500/20 border-amber-400 scale-105 shadow-md shadow-amber-500/20'
                      : 'bg-slate-950 border-slate-800 hover:border-slate-700'
                  }`}
                >
                  {av}
                </button>
              ))}
            </div>
          </div>

          {/* Jauna parole (neobligāta) */}
          <div className="space-y-1.5 pt-1">
            <label className="text-xs font-bold text-slate-300 flex items-center justify-between">
              <span>Mainīt paroli</span>
              <span className="text-[10px] text-slate-500 font-normal">Atstāj tukšu, ja nemaini</span>
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
                <Lock className="w-4 h-4" />
              </div>
              <input
                type="password"
                minLength={6}
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="Jaunā parole (vismaz 6 simboli)"
                className="w-full pl-9 pr-3 py-2.5 rounded-xl bg-slate-950 border border-slate-750 text-slate-100 text-sm focus:outline-none focus:border-amber-400"
              />
            </div>
          </div>

          {/* Pogas */}
          <div className="pt-2 flex gap-2">
            <button
              type="button"
              onClick={onClose}
              className="w-1/3 py-2.5 rounded-xl border border-slate-750 bg-slate-850 hover:bg-slate-800 text-slate-300 font-bold text-xs transition-colors"
            >
              Atcelt
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="w-2/3 py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-black text-xs flex items-center justify-center gap-2 transition-all shadow-lg shadow-amber-500/20 active:scale-98 disabled:opacity-50"
            >
              {submitting ? (
                <div className="w-4 h-4 rounded-full border-2 border-slate-950 border-t-transparent animate-spin" />
              ) : (
                'Saglabāt izmaiņas'
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
