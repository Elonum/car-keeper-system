import React from 'react';
import { Card } from "@/components/ui/card";
import PriceDisplay from '../common/PriceDisplay';
import { Check, Palette, Package } from 'lucide-react';

export default function ConfigurationSummary({ trim, color, selectedOptions, totalPrice }) {
  return (
    <div className="space-y-6">
      {/* Vehicle */}
      <Card className="p-6 bg-gradient-to-br from-slate-50 to-white">
        <h3 className="text-sm font-semibold text-slate-500 uppercase tracking-wider mb-3">
          Автомобиль
        </h3>
        <div className="space-y-2">
          <p className="text-lg font-bold text-slate-900">
            {trim?.brand_name} {trim?.model_name}
          </p>
          {trim?.generation_name && (
            <p className="text-sm text-slate-500">{trim.generation_name}</p>
          )}
          <p className="text-base font-semibold text-slate-700">
            Комплектация: {trim?.name}
          </p>
          <div className="pt-2 border-t border-slate-200">
            <PriceDisplay price={trim?.base_price} size="md" />
          </div>
        </div>
      </Card>

      {/* Color */}
      {color && (
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-3">
            <Palette className="w-4 h-4 text-slate-500" />
            <h3 className="text-sm font-semibold text-slate-500 uppercase tracking-wider">
              Цвет
            </h3>
          </div>
          <div className="flex items-center gap-3">
            <div 
              className="w-10 h-10 rounded-full border-2 border-slate-200 shadow-sm flex-shrink-0"
              style={{ backgroundColor: color.hex_code }}
            />
            <div className="flex-1">
              <p className="font-semibold text-slate-900">{color.name}</p>
              {color.price_delta > 0 ? (
                <PriceDisplay price={color.price_delta} size="sm" className="text-sm" prefix="+" />
              ) : (
                <p className="text-xs text-green-600 font-medium">Включено</p>
              )}
            </div>
          </div>
        </Card>
      )}

      {/* Options */}
      {selectedOptions && selectedOptions.length > 0 && (
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-3">
            <Package className="w-4 h-4 text-slate-500" />
            <h3 className="text-sm font-semibold text-slate-500 uppercase tracking-wider">
              Опции ({selectedOptions.length})
            </h3>
          </div>
          <div className="space-y-2">
            {selectedOptions.map(opt => (
              <div key={opt.id} className="flex items-start justify-between gap-3 py-2 border-b border-slate-100 last:border-0">
                <div className="flex items-start gap-2 flex-1">
                  <Check className="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
                  <p className="text-sm text-slate-700">{opt.name}</p>
                </div>
                <PriceDisplay price={opt.price} size="sm" className="text-sm flex-shrink-0" prefix="+" />
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* Total */}
      <Card className="p-6 bg-slate-900 text-white">
        <div className="flex items-end justify-between">
          <div>
            <p className="text-sm text-slate-300 mb-1">Итоговая цена</p>
            <PriceDisplay price={totalPrice} size="xl" className="text-white" />
          </div>
        </div>
      </Card>
    </div>
  );
}