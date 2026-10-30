import React from 'react';
import { Tag, Users } from 'lucide-react';

interface HeaderProps {
  roomCode?: string;
  playerCount?: number;
  onLeave?: () => void;
}

export const Header: React.FC<HeaderProps> = ({ roomCode, playerCount, onLeave }) => {
  return (
    <header className="w-full max-w-2xl mx-auto px-4 py-3 flex items-center justify-between border-b border-slate-800">
      <div className="flex items-center gap-2 cursor-pointer" onClick={() => onLeave && onLeave()}>
        <div className="w-10 h-10 rounded-xl bg-amber-500 flex items-center justify-center font-black text-slate-950 text-xl shadow-lg shadow-amber-500/20">
          <Tag className="w-6 h-6 stroke-[2.5]" />
        </div>
        <div>
          <span className="font-extrabold text-lg tracking-tight bg-gradient-to-r from-amber-400 to-orange-400 bg-clip-text text-transparent">
            MiniCenu
          </span>
          <span className="text-xs ml-1.5 px-1.5 py-0.5 rounded bg-slate-800 text-slate-400 font-medium">
            ss.lv
          </span>
        </div>
      </div>

      {roomCode && (
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5 px-3 py-1 rounded-full bg-slate-800 border border-slate-700">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span className="font-mono font-bold text-amber-400 text-sm tracking-wider">{roomCode}</span>
            {playerCount !== undefined && (
              <span className="text-xs text-slate-400 flex items-center gap-1 ml-1 border-l border-slate-700 pl-2">
                <Users className="w-3.5 h-3.5" />
                {playerCount}
              </span>
            )}
          </div>
          {onLeave && (
            <button
              onClick={onLeave}
              className="text-xs text-slate-400 hover:text-rose-400 transition-colors px-2 py-1"
            >
              Iziet
            </button>
          )}
        </div>
      )}
    </header>
  );
};

