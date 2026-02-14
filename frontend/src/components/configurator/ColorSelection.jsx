import React from 'react';
import { Card } from "@/components/ui/card";
import { Check } from 'lucide-react';
import { cn } from "@/lib/utils";
import PriceDisplay from '../common/PriceDisplay';

export default function ColorSelection({ colors, selectedColorId, onSelectColor }) {
  if (!colors || colors.length === 0) {
    return (
      <div className="text-center py-12 text-slate-500">
        <p>Нет доступных цветов</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
      {colors.map(color => {
        const isSelected = selectedColorId === color.id;
        const isAvailable = color.is_available !== false;

        return (
          <button
            key={color.id}
            onClick={() => isAvailable && onSelectColor(color)}
            disabled={!isAvailable}
            className={cn(
              "group relative rounded-xl p-4 border-2 transition-all",
              isSelected 
                ? "border-blue-600 shadow-lg scale-105" 
                : "border-slate-200 hover:border-slate-300 hover:shadow-md",
              !isAvailable && "opacity-50 cursor-not-allowed"
            )}
          >
            {/* Color circle */}
            <div className="relative mb-3">
              <div 
                className="w-16 h-16 mx-auto rounded-full border-2 border-slate-200 shadow-inner"
                style={{ backgroundColor: color.hex_code }}
              />
              {isSelected && (
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="w-7 h-7 rounded-full bg-blue-600 flex items-center justify-center shadow-lg">
                    <Check className="w-4 h-4 text-white" strokeWidth={3} />
                  </div>
                </div>
              )}
            </div>

            {/* Name */}
            <p className="text-sm font-medium text-slate-900 mb-1 text-center">
              {color.name}
            </p>

            {/* Price delta */}
            {color.price_delta > 0 ? (
              <p className="text-xs text-slate-500 text-center">
                +<PriceDisplay price={color.price_delta} size="sm" className="text-xs" />
              </p>
            ) : (
              <p className="text-xs text-green-600 text-center font-medium">
                Включено
              </p>
            )}

            {!isAvailable && (
              <div className="absolute inset-0 flex items-center justify-center bg-white/80 rounded-xl">
                <span className="text-xs font-medium text-slate-400">Недоступен</span>
              </div>
            )}
          </button>
        );
      })}
    </div>
  );
}