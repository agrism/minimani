import React, { useState, useEffect } from 'react';
import { X, HelpCircle, Zap, Target, Trophy, Sparkles, Clock, CheckCircle2, Calculator } from 'lucide-react';
import { ScoringRulesResponse } from '../types';

interface RulesModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const RulesModal: React.FC<RulesModalProps> = ({ isOpen, onClose }) => {
  const [rules, setRules] = useState<ScoringRulesResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [simSeconds, setSimSeconds] = useState(5);
  const totalLimitSec = 45;

  useEffect(() => {
    if (!isOpen) return;

    fetch('/api/rules')
      .then((res) => res.json())
      .then((data: ScoringRulesResponse) => {
        setRules(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to load scoring rules:', err);
        setLoading(false);
      });
  }, [isOpen]);

  if (!isOpen) return null;

  const currentRules = rules?.current_rules;
  const basePoints = currentRules?.base_points || 500;
  const maxSpeed = currentRules?.max_speed_bonus || 500;

  // Live calculator calculation
  const remaining = Math.max(0, totalLimitSec - simSeconds);
  const calculatedSpeedBonus = Math.round((remaining / totalLimitSec) * maxSpeed);
  const calculatedTotal = basePoints + calculatedSpeedBonus;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm animate-fade-in">
      <div
        className="relative w-full max-w-lg bg-slate-900 border border-slate-750 rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Galvene */}
        <div className="p-5 border-b border-slate-800 flex items-center justify-between bg-slate-850/60">
          <div className="flex items-center gap-2.5">
            <div className="w-9 h-9 rounded-xl bg-amber-500/20 text-amber-400 flex items-center justify-center border border-amber-500/30">
              <HelpCircle className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-extrabold text-base text-slate-100">Punktu Aprēķina Noteikumi</h3>
              <p className="text-[11px] text-slate-400">Dinamiski sasaistīti ar spēles reālo loģiku</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-xl text-slate-400 hover:text-slate-100 hover:bg-slate-800 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Saturs */}
        <div className="p-5 overflow-y-auto space-y-5 text-sm">
          {loading && (
            <div className="flex flex-col items-center justify-center py-12 space-y-3">
              <div className="w-8 h-8 rounded-full border-2 border-amber-400 border-t-transparent animate-spin" />
              <p className="text-xs text-slate-400">Ielādējam punktu sistēmas parametrus...</p>
            </div>
          )}

          {!loading && currentRules && (
            <>
              {/* Aktīvā spēles režīma kartīte */}
              <div className="bg-gradient-to-br from-amber-500/10 via-orange-500/5 to-slate-900 border border-amber-500/30 rounded-2xl p-4 space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-[11px] font-bold uppercase tracking-wider text-amber-400 bg-amber-500/20 px-2.5 py-0.5 rounded-full border border-amber-500/40">
                    Aktīvais režīms
                  </span>
                  <span className="text-xs font-mono font-bold text-slate-300">
                    Maks. {currentRules.max_total_points} pt / kārtā
                  </span>
                </div>
                <h4 className="font-bold text-slate-100 text-base">{currentRules.title}</h4>
                <p className="text-xs text-slate-300 leading-relaxed">{currentRules.description}</p>
              </div>

              {/* Punktu sadalījums (Cards) */}
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-slate-850 border border-slate-750 p-3.5 rounded-2xl space-y-1">
                  <div className="flex items-center gap-1.5 text-xs font-bold text-emerald-400">
                    <Target className="w-4 h-4" />
                    <span>Bāzes punkti</span>
                  </div>
                  <div className="text-2xl font-black text-slate-100 font-mono">
                    +{currentRules.base_points}
                  </div>
                  <p className="text-[11px] text-slate-400 leading-tight">
                    Garantēti par katru pareizi izvēlēto cenu
                  </p>
                </div>

                <div className="bg-slate-850 border border-slate-750 p-3.5 rounded-2xl space-y-1">
                  <div className="flex items-center gap-1.5 text-xs font-bold text-amber-400">
                    <Zap className="w-4 h-4" />
                    <span>Ātruma bonuss</span>
                  </div>
                  <div className="text-2xl font-black text-slate-100 font-mono">
                    līdz +{currentRules.max_speed_bonus}
                  </div>
                  <p className="text-[11px] text-slate-400 leading-tight">
                    Jo ātrāk atbildi, jo lielāku bonusu saņem
                  </p>
                </div>
              </div>

              {/* Interaktīvais kalkulators */}
              <div className="bg-slate-850/80 border border-slate-750 rounded-2xl p-4 space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Calculator className="w-4 h-4 text-amber-400" />
                    <h5 className="font-bold text-xs text-slate-200">Interaktīvs punktu aprēķins</h5>
                  </div>
                  <span className="text-xs font-mono font-bold text-amber-400 bg-slate-900 px-2 py-0.5 rounded-lg border border-slate-700">
                    Atbilde {simSeconds}. sekundē
                  </span>
                </div>

                <input
                  type="range"
                  min="0"
                  max={totalLimitSec}
                  value={simSeconds}
                  onChange={(e) => setSimSeconds(Number(e.target.value))}
                  className="w-full accent-amber-500 cursor-pointer"
                />

                <div className="flex justify-between text-[10px] text-slate-400 font-mono">
                  <span>0s (Uzreiz)</span>
                  <span>{totalLimitSec / 2}s (Puse laika)</span>
                  <span>{totalLimitSec}s (Pēdējā sekunde)</span>
                </div>

                {/* Rezultāts */}
                <div className="p-3 rounded-xl bg-slate-900 border border-slate-750 flex items-center justify-between">
                  <div className="text-xs text-slate-300">
                    <span>{basePoints} bāze</span>
                    <span className="text-slate-500 mx-1.5">+</span>
                    <span className="text-amber-400 font-semibold">{calculatedSpeedBonus} ātrums</span>
                  </div>
                  <div className="text-right">
                    <span className="text-base font-black text-emerald-400 font-mono">
                      = {calculatedTotal} pt
                    </span>
                  </div>
                </div>
              </div>

              {/* Formulas paskaidrojums */}
              <div className="space-y-1.5">
                <h5 className="font-bold text-xs text-slate-300">Precīza aprēķina formula:</h5>
                <div className="p-3 rounded-xl bg-slate-950 border border-slate-800 text-xs font-mono text-slate-300 leading-relaxed">
                  {currentRules.formula_explanation}
                </div>
              </div>

              {/* Piemēri */}
              <div className="space-y-1.5">
                <h5 className="font-bold text-xs text-slate-300">Piemēri:</h5>
                <div className="space-y-1">
                  {currentRules.examples.map((ex, i) => (
                    <div key={i} className="flex items-start gap-2 text-xs text-slate-300">
                      <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400 shrink-0 mt-0.5" />
                      <span>{ex}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Padomi uzvarai */}
              {rules.general_tips && rules.general_tips.length > 0 && (
                <div className="p-3.5 rounded-2xl bg-amber-500/5 border border-amber-500/20 space-y-1.5">
                  <div className="flex items-center gap-1.5 text-xs font-bold text-amber-400">
                    <Sparkles className="w-3.5 h-3.5" />
                    <span>Padomi labākam rezultātam:</span>
                  </div>
                  <ul className="list-disc list-inside text-xs text-slate-300 space-y-1">
                    {rules.general_tips.map((tip, idx) => (
                      <li key={idx} className="leading-snug">{tip}</li>
                    ))}
                  </ul>
                </div>
              )}
            </>
          )}
        </div>

        {/* Kājene */}
        <div className="p-4 border-t border-slate-800 bg-slate-850/60 flex justify-end">
          <button
            onClick={onClose}
            className="w-full py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs transition-colors"
          >
            Sapratu, aizvērt
          </button>
        </div>
      </div>
    </div>
  );
};
