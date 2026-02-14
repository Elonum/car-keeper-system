import React from 'react';
import { cn } from "@/lib/utils";
import { Check, MapPin, Phone, Mail } from 'lucide-react';

export default function BranchSelector({ branches, selectedBranchId, onSelect }) {
  return (
    <div>
      <h2 className="text-xl font-bold text-slate-900 mb-1">Выберите филиал</h2>
      <p className="text-sm text-slate-500 mb-5">Выберите удобный для вас сервисный центр</p>
      <div className="grid gap-3">
        {branches.map(branch => {
          const isSelected = selectedBranchId === branch.id;
          return (
            <button
              key={branch.id}
              onClick={() => onSelect(branch.id)}
              className={cn(
                "flex items-start gap-4 p-5 rounded-xl border-2 text-left transition-all w-full",
                isSelected ? "border-slate-900 bg-slate-50" : "border-slate-100 bg-white hover:border-slate-200"
              )}
            >
              {isSelected ? (
                <div className="w-6 h-6 rounded-full bg-slate-900 flex items-center justify-center flex-shrink-0 mt-0.5">
                  <Check className="w-3.5 h-3.5 text-white" />
                </div>
              ) : (
                <div className="w-6 h-6 rounded-full border-2 border-slate-200 flex-shrink-0 mt-0.5" />
              )}
              <div className="flex-1">
                <p className="font-semibold text-slate-900">{branch.name}</p>
                <div className="flex items-center gap-1.5 text-sm text-slate-500 mt-1">
                  <MapPin className="w-3.5 h-3.5 flex-shrink-0" />
                  <span>{branch.address}</span>
                </div>
                <div className="flex flex-wrap gap-4 mt-2">
                  {branch.phone && (
                    <div className="flex items-center gap-1.5 text-xs text-slate-400">
                      <Phone className="w-3 h-3" /> {branch.phone}
                    </div>
                  )}
                  {branch.email && (
                    <div className="flex items-center gap-1.5 text-xs text-slate-400">
                      <Mail className="w-3 h-3" /> {branch.email}
                    </div>
                  )}
                </div>
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}