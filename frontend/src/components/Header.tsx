import React from 'react';
import { Tag, Users, LogOut, HelpCircle, User as UserIcon } from 'lucide-react';
import { User } from '../types';

interface HeaderProps {
  roomCode?: string;
  playerCount?: number;
  currentUser?: User | null;
  onLeave?: () => void;
  onOpenRules?: () => void;
  onOpenAuth?: () => void;
  onLogout?: () => void;
}

export const Header: React.FC<HeaderProps> = ({
  roomCode,
  playerCount,
  currentUser,
  onLeave,
  onOpenRules,
  onOpenAuth,
  onLogout,
}) => {
  return (
    <header className="w-full max-w-2xl mx-auto px-4 py-3 flex items-center justify-between border-b border-slate-800">
      <div
        className={`flex items-center gap-2 ${onLeave ? 'cursor-pointer hover:opacity-85 transition-opacity' : ''}`}
        onClick={() => onLeave && onLeave()}
        title={onLeave ? 'Doties uz sākumlapu' : undefined}
      >
        <div className="w-10 h-10 rounded-xl bg-amber-500 flex items-center justify-center font-black text-slate-950 text-xl shadow-lg shadow-amber-500/20">
          <Tag className="w-6 h-6 stroke-[2.5]" />
        </div>
        <div>
          <span className="font-extrabold text-lg tracking-tight bg-gradient-to-r from-amber-400 to-orange-400 bg-clip-text text-transparent">
            MiniMani
          </span>
          <span className="text-xs ml-1.5 px-1.5 py-0.5 rounded bg-slate-800 text-slate-400 font-medium">
            ss.lv
          </span>
        </div>
      </div>

      <div className="flex items-center gap-2 sm:gap-2.5">
        {onOpenRules && (
          <button
            onClick={onOpenRules}
            className="flex items-center gap-1.5 text-xs text-slate-300 hover:text-amber-400 bg-slate-850 hover:bg-slate-800 border border-slate-750 px-2.5 py-1.5 rounded-xl font-bold transition-all shadow-sm active:scale-95"
            title="Kā tiek aprēķināti punkti?"
          >
            <HelpCircle className="w-3.5 h-3.5 text-amber-400" />
            <span className="hidden sm:inline">Noteikumi</span>
          </button>
        )}

        {/* Lietotāja profila vai Ienākšanas poga */}
        {currentUser ? (
          <div className="flex items-center gap-1.5 bg-slate-850 border border-slate-750 px-2.5 py-1.5 rounded-xl shadow-sm">
            <span className="text-base">{currentUser.avatar}</span>
            <span className="text-xs font-bold text-slate-200 hidden sm:inline max-w-[100px] truncate">
              {currentUser.nickname}
            </span>
            {onLogout && (
              <button
                onClick={onLogout}
                className="text-[10px] text-slate-400 hover:text-rose-400 font-bold ml-1 border-l border-slate-700 pl-1.5 transition-colors"
                title="Iziet no profila (spēlēt kā viesim)"
              >
                Iziet
              </button>
            )}
          </div>
        ) : (
          onOpenAuth && (
            <button
              onClick={onOpenAuth}
              className="flex items-center gap-1.5 text-xs text-amber-400 hover:text-amber-300 bg-amber-500/10 hover:bg-amber-500/20 border border-amber-500/30 px-2.5 py-1.5 rounded-xl font-bold transition-all shadow-sm active:scale-95"
            >
              <UserIcon className="w-3.5 h-3.5" />
              <span>Ienākt</span>
            </button>
          )
        )}

        {roomCode && (
          <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-slate-800 border border-slate-700 shadow-sm">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span className="font-mono font-bold text-amber-400 text-sm tracking-wider">{roomCode}</span>
            {playerCount !== undefined && (
              <span className="text-xs text-slate-400 flex items-center gap-1 ml-1 border-l border-slate-700 pl-2 font-medium">
                <Users className="w-3.5 h-3.5" />
                {playerCount}
              </span>
            )}
          </div>
        )}

        {roomCode && onLeave && (
          <button
            onClick={onLeave}
            className="flex items-center gap-1.5 text-xs text-rose-300 hover:text-white bg-rose-950/40 hover:bg-rose-900/60 border border-rose-800/60 hover:border-rose-500 transition-all px-3 py-1.5 rounded-xl font-bold shadow-sm active:scale-95"
            title="Pamest istabu un izvēlēties citu"
          >
            <LogOut className="w-3.5 h-3.5" />
            <span>Pamest</span>
          </button>
        )}
      </div>
    </header>
  );
};

