import React, { useState, useEffect } from 'react';
import { Play, Zap, Plus, Users, Lock, Globe, Sparkles } from 'lucide-react';
import { PublicRoomItem } from '../types';

interface HomeScreenProps {
  nickname: string;
  avatar: string;
  onUpdateProfile: (name: string, avatar: string) => void;
  onCreateRoom: (isPublic: boolean, category: string, gameMode: string) => void;
  onJoinRoom: (code: string) => void;
  onQuickJoin: () => void;
}

const AVATARS = ['🚗', '🏎️', '🚜', '🏡', '💎', '📱', '🔨', '🍻', '🌲', '👑', '🚀', '🔥'];

export const HomeScreen: React.FC<HomeScreenProps> = ({
  nickname,
  avatar,
  onUpdateProfile,
  onCreateRoom,
  onJoinRoom,
  onQuickJoin,
}) => {
  const [code, setCode] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isPublic, setIsPublic] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [selectedMode, setSelectedMode] = useState('exact');
  const [publicRooms, setPublicRooms] = useState<PublicRoomItem[]>([]);
  const [isLoadingRooms, setIsLoadingRooms] = useState(false);

  useEffect(() => {
    fetchPublicRooms();
    const interval = setInterval(fetchPublicRooms, 4000);
    return () => clearInterval(interval);
  }, []);

  const fetchPublicRooms = async () => {
    try {
      setIsLoadingRooms(true);
      const res = await fetch('/api/rooms/public');
      if (res.ok) {
        const data = await res.json();
        setPublicRooms(data);
      }
    } catch (e) {
      // ignore
    } finally {
      setIsLoadingRooms(false);
    }
  };

  const handleJoinByCode = (e: React.FormEvent) => {
    e.preventDefault();
    if (code.trim()) {
      onJoinRoom(code.trim().toUpperCase());
    }
  };

  return (
    <div className="w-full max-w-md mx-auto px-4 py-6 space-y-6">
      {/* Profila kartīte */}
      <div className="bg-slate-800/80 backdrop-blur border border-slate-700/80 rounded-2xl p-4 shadow-xl">
        <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider block mb-2">
          Tavs spēlētāja profils
        </label>
        <div className="flex items-center gap-3">
          <div className="relative group">
            <button
              type="button"
              className="w-14 h-14 rounded-2xl bg-slate-700 flex items-center justify-center text-3xl border-2 border-amber-500/50 hover:border-amber-400 transition-all shadow-inner"
              title="Mainīt avataru"
            >
              {avatar}
            </button>
            <div className="absolute top-16 left-0 z-20 hidden group-hover:flex group-focus-within:flex flex-wrap gap-1.5 p-2 bg-slate-900 border border-slate-700 rounded-xl shadow-2xl w-48">
              {AVATARS.map((a) => (
                <button
                  key={a}
                  onClick={() => onUpdateProfile(nickname, a)}
                  className="w-9 h-9 rounded-lg hover:bg-slate-800 text-xl flex items-center justify-center transition-colors"
                >
                  {a}
                </button>
              ))}
            </div>
          </div>
          <div className="flex-1">
            <input
              type="text"
              value={nickname}
              onChange={(e) => onUpdateProfile(e.target.value, avatar)}
              placeholder="Ievadi savu vārdu..."
              maxLength={20}
              className="w-full bg-slate-900 border border-slate-700 rounded-xl px-3.5 py-2.5 text-slate-100 placeholder-slate-500 font-semibold focus:outline-none focus:border-amber-500 transition-colors"
            />
          </div>
        </div>
      </div>

      {/* Galvenās pogas */}
      <div className="grid grid-cols-1 gap-3">
        {/* Ātrā spēle */}
        <button
          onClick={onQuickJoin}
          className="w-full py-4 px-5 rounded-2xl bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-400 hover:to-orange-400 text-slate-950 font-black text-lg flex items-center justify-between shadow-lg shadow-orange-500/25 transition-all transform active:scale-[0.98]"
        >
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-slate-950/20 flex items-center justify-center">
              <Zap className="w-6 h-6 stroke-[2.5]" />
            </div>
            <div className="text-left">
              <div>Ātrā Spēle</div>
              <div className="text-xs font-semibold text-slate-900/80">Ielēkt brīvā atvērtajā istabā</div>
            </div>
          </div>
          <Play className="w-5 h-5 fill-current" />
        </button>

        {/* Izveidot spēli */}
        <button
          onClick={() => setShowCreateModal(true)}
          className="w-full py-3.5 px-5 rounded-2xl bg-slate-800 hover:bg-slate-750 border border-slate-700 font-bold text-slate-200 flex items-center justify-between transition-all active:scale-[0.98]"
        >
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-slate-700 flex items-center justify-center text-amber-400">
              <Plus className="w-5 h-5 stroke-[2.5]" />
            </div>
            <div className="text-left">
              <div className="text-base">Izveidot Spēles Istabu</div>
              <div className="text-xs text-slate-400 font-normal">Privāta draugiem vai atvērta ikvienam</div>
            </div>
          </div>
        </button>
      </div>

      {/* Pievienoties ar kodu */}
      <form onSubmit={handleJoinByCode} className="bg-slate-800/40 border border-slate-700/60 rounded-2xl p-4">
        <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider block mb-2">
          Pievienoties ar istabas kodu
        </label>
        <div className="flex gap-2">
          <input
            type="text"
            value={code}
            onChange={(e) => setCode(e.target.value.toUpperCase())}
            placeholder="Piem. AUTO"
            maxLength={6}
            className="flex-1 bg-slate-900 border border-slate-700 rounded-xl px-4 py-2.5 font-mono text-center tracking-widest uppercase font-bold text-lg text-amber-400 placeholder-slate-600 focus:outline-none focus:border-amber-500"
          />
          <button
            type="submit"
            disabled={!code.trim()}
            className="px-5 rounded-xl bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-slate-100 font-bold transition-colors"
          >
            Ieiet
          </button>
        </div>
      </form>

      {/* Atvērtās publiskās istabas */}
      <div className="space-y-2.5">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1.5">
            <Globe className="w-4 h-4 text-emerald-400" />
            <h3 className="text-sm font-bold text-slate-200">Atvērtās spēles tiešsaistē</h3>
          </div>
          <span className="text-xs text-slate-500 font-medium">
            {publicRooms.length} {publicRooms.length === 1 ? 'istaba' : 'istabas'}
          </span>
        </div>

        {publicRooms.length === 0 ? (
          <div className="text-center py-6 px-4 rounded-2xl border border-dashed border-slate-800 text-slate-500 text-xs">
            Pašlaik nav atvērtu istabu. Esi pirmais un izveido atvērtu spēli!
          </div>
        ) : (
          <div className="space-y-2">
            {publicRooms.map((room) => (
              <div
                key={room.code}
                onClick={() => onJoinRoom(room.code)}
                className="flex items-center justify-between p-3.5 rounded-xl bg-slate-800/80 hover:bg-slate-750 border border-slate-700/80 cursor-pointer transition-all hover:border-amber-500/50"
              >
                <div className="flex items-center gap-3">
                  <div className="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-pulse"></div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="font-mono font-bold text-amber-400">{room.code}</span>
                      <span className="text-xs px-2 py-0.5 rounded-full bg-slate-700 text-slate-300 font-medium">
                        {room.category === 'auto'
                          ? '🚗 Auto'
                          : room.category === 'real_estate'
                          ? '🏢 Īpašumi'
                          : room.category === 'tech'
                          ? '🛠️ Tehnika'
                          : room.category === 'pearls'
                          ? '💎 Pērles'
                          : '🎲 Viss'}
                      </span>
                    </div>
                    <div className="text-xs text-slate-400 mt-0.5">
                      {room.question_count} jautājumi · 4 Varianti
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-slate-300 font-medium flex items-center gap-1 bg-slate-900 px-2.5 py-1 rounded-lg border border-slate-750">
                    <Users className="w-3.5 h-3.5 text-slate-400" />
                    {room.player_count}/{room.max_players}
                  </span>
                  <span className="text-xs font-bold text-amber-400">Pievienoties →</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Izveidošanas modālais logs */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-sm w-full p-5 space-y-4 shadow-2xl">
            <h3 className="font-bold text-lg text-slate-100 flex items-center gap-2">
              <Sparkles className="w-5 h-5 text-amber-400" />
              Jaunas spēles iestatījumi
            </h3>

            {/* Publiska vai privāta */}
            <div className="space-y-1.5">
              <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider">
                Istabas redzamība
              </label>
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={() => setIsPublic(true)}
                  className={`p-3 rounded-xl border text-sm font-semibold flex items-center gap-2 justify-center transition-all ${
                    isPublic
                      ? 'bg-amber-500/10 border-amber-500 text-amber-400'
                      : 'bg-slate-800 border-slate-700 text-slate-400 hover:text-slate-200'
                  }`}
                >
                  <Globe className="w-4 h-4" />
                  Publiska (Ikvienam)
                </button>
                <button
                  type="button"
                  onClick={() => setIsPublic(false)}
                  className={`p-3 rounded-xl border text-sm font-semibold flex items-center gap-2 justify-center transition-all ${
                    !isPublic
                      ? 'bg-amber-500/10 border-amber-500 text-amber-400'
                      : 'bg-slate-800 border-slate-700 text-slate-400 hover:text-slate-200'
                  }`}
                >
                  <Lock className="w-4 h-4" />
                  Privāta (Ar kodu)
                </button>
              </div>
            </div>

            {/* Kategorija */}
            <div className="space-y-1.5">
              <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider">
                Sākotnējā ss.lv kategorija
              </label>
              <select
                value={selectedCategory}
                onChange={(e) => setSelectedCategory(e.target.value)}
                className="w-full bg-slate-800 border border-slate-700 rounded-xl px-3.5 py-2.5 text-slate-100 font-medium focus:outline-none focus:border-amber-500"
              >
                <option value="all">🎲 Visas ss.lv kategorijas miksētas</option>
                <option value="auto">🚗 ss.lv Auto & Transports</option>
                <option value="real_estate">🏢 ss.lv Nekustamie īpašumi</option>
                <option value="tech">🛠️ ss.lv Tehnika & Būvniecība</option>
                <option value="pearls">💎 ss.lv Pērles & Kuriozi</option>
              </select>
            </div>

            <div className="pt-2 flex gap-2">
              <button
                type="button"
                onClick={() => setShowCreateModal(false)}
                className="flex-1 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-sm transition-colors"
              >
                Atcelt
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowCreateModal(false);
                  onCreateRoom(isPublic, selectedCategory, 'multiple_choice');
                }}
                className="flex-1 py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-sm transition-colors shadow-lg shadow-amber-500/20"
              >
                Izveidot istabu
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

