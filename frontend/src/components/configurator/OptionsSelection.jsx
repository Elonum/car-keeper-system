import React from 'react';
import { Card } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import PriceDisplay from '../common/PriceDisplay';
import { Badge } from "@/components/ui/badge";

const categoryLabels = {
  safety: "Безопасность",
  comfort: "Комфорт",
  exterior: "Экстерьер",
  interior: "Интерьер",
  technology: "Технологии",
  performance: "Производительность",
  other: "Прочее",
  others: "Прочее",
};

const categoryColors = {
  safety: "bg-red-50 text-red-700 border-red-200",
  comfort: "bg-blue-50 text-blue-700 border-blue-200",
  exterior: "bg-purple-50 text-purple-700 border-purple-200",
  interior: "bg-amber-50 text-amber-700 border-amber-200",
  technology: "bg-green-50 text-green-700 border-green-200",
  performance: "bg-orange-50 text-orange-700 border-orange-200",
};

export default function OptionsSelection({ options, selectedOptionIds, onToggleOption }) {
  if (!options || options.length === 0) {
    return (
      <div className="text-center py-12 text-slate-500">
        <p>Нет доступных опций</p>
      </div>
    );
  }

  // Group by category
  const grouped = options.reduce((acc, opt) => {
    const cat = opt.category || 'other';
    if (!acc[cat]) acc[cat] = [];
    acc[cat].push(opt);
    return acc;
  }, {});

  return (
    <div className="space-y-6">
      {Object.entries(grouped).map(([category, opts]) => (
        <div key={category}>
          <div className="flex items-center gap-2 mb-3">
            <Badge variant="outline" className={`${categoryColors[category] || 'bg-slate-50 text-slate-700'} border text-xs font-medium px-2.5 py-1`}>
              {categoryLabels[category] || category}
            </Badge>
            <div className="h-px flex-1 bg-slate-200" />
          </div>

          <div className="grid gap-3">
            {opts.map(option => {
              const optionId = option.option_id || option.id;
              const isSelected = selectedOptionIds.includes(optionId);
              const isAvailable = option.is_available !== false;

              return (
                <Card 
                  key={optionId}
                  className={`p-4 transition-all cursor-pointer ${
                    isSelected ? 'border-blue-600 bg-blue-50/30 shadow-md' : 'hover:border-slate-300 hover:shadow-sm'
                  } ${!isAvailable ? 'opacity-50' : ''}`}
                  onClick={() => isAvailable && onToggleOption(optionId)}
                >
                  <div className="flex items-start gap-4">
                    <Checkbox 
                      checked={isSelected}
                      disabled={!isAvailable}
                      className="mt-1"
                    />
                    
                    <div className="flex-1 min-w-0">
                      <div className="flex items-start justify-between gap-3 mb-1">
                        <h4 className="font-semibold text-slate-900">{option.name}</h4>
                        <PriceDisplay price={option.price} size="sm" className="flex-shrink-0" prefix="+" />
                      </div>
                      
                      {option.description && (
                        <p className="text-sm text-slate-500 leading-relaxed">{option.description}</p>
                      )}

                      {!isAvailable && (
                        <p className="text-xs text-red-500 mt-1 font-medium">Недоступно</p>
                      )}
                    </div>
                  </div>
                </Card>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}