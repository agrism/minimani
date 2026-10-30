import React from 'react';
import { Tag, Users, LogOut } from 'lucide-react';

interface HeaderProps {
  roomCode?: string;
  playerCount?: number;
  onLeave?: () => void;
}

export const Header: React.FC<HeaderProps> = ({ roomCode, playerCount, onLeave }) => {
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

      {roomCode && (
        <div className="flex items-center gap-2 sm:gap-3">
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
          {onLeave && (
            <button
              onClick={onLeave}
              className="flex items-center gap-1.5 text-xs text-rose-300 hover:text-white bg-rose-950/40 hover:bg-rose-900/60 border border-rose-800/60 hover:border-rose-500 transition-all px-3 py-1.5 rounded-xl font-bold shadow-sm active:scale-95"
              title="Pamest istabu un izvēlēties citu"
            >
              <LogOut className="w-3.5 h-3.5" />
              <span>Iziet</span>
            </button>
          )}
        </div>
      )}
    </header>
  );
};

