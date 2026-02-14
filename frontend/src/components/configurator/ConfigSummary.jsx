import React from 'react';
import PriceDisplay from '../common/PriceDisplay';
import { Card } from "@/components/ui/card";
import { Check, Palette, Settings, Car } from 'lucide-react';

export default function ConfigSummary({ trim, color, selectedOptions, totalPrice }) {
  return (
    <div>
      <h2 className="text-xl font-bold text-slate-900 mb-1">Ваша конфигурация</h2>
      <p className="text-sm text-slate-500 mb-6">Проверьте все параметры перед оформлением</p>

      <div className="space-y-4">
        {/* Car info */}
        <Card className="p-5 border-0 shadow-sm">
          <div className="flex items-start gap-3">
            <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
              <Car className="w-5 h-5 text-blue-600" />
            </div>
            <div>
              <p className="text-xs text-slate-500 font-medium uppercase tracking-wider">Автомобиль</p>
              <p className="text-lg font-bold text-slate-900 mt-0.5">
                {trim?.brand_name} {trim?.model_name}
              </p>
              <p className="text-sm text-slate-500">
                {trim?.generation_name} • {trim?.name}
              </p>
              <div className="mt-2">
                <PriceDisplay price={trim?.base_price} size="sm" className="text-slate-600" />
              </div>
            </div>
          </div>
        </Card>

        {/* Color */}
        {color && (
          <Card className="p-5 border-0 shadow-sm">
            <div className="flex items-start gap-3">
              <div className="w-10 h-10 rounded-xl bg-purple-50 flex items-center justify-center flex-shrink-0">
                <Palette className="w-5 h-5 text-purple-600" />
              </div>
              <div className="flex-1">
                <p className="text-xs text-slate-500 font-medium uppercase tracking-wider">Цвет</p>
                <div className="flex items-center gap-2 mt-1">
                  <div className="w-5 h-5 rounded-full border border-black/10" style={{ backgroundColor: color.hex_code }} />
                  <p className="font-semibold text-slate-900">{color.name}</p>
                </div>
                {color.price_delta > 0 && (
                  <PriceDisplay price={color.price_delta} size="sm" prefix="+" className="text-blue-600 mt-1" />
                )}
              </div>
            </div>
          </Card>
        )}

        {/* Options */}
        {selectedOptions.length > 0 && (
          <Card className="p-5 border-0 shadow-sm">
            <div className="flex items-start gap-3">
              <div className="w-10 h-10 rounded-xl bg-emerald-50 flex items-center justify-center flex-shrink-0">
                <Settings className="w-5 h-5 text-emerald-600" />
              </div>
              <div className="flex-1">
                <p className="text-xs text-slate-500 font-medium uppercase tracking-wider mb-2">
                  Опции ({selectedOptions.length})
                </p>
                <div className="space-y-2">
                  {selectedOptions.map(opt => (
                    <div key={opt.id} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Check className="w-3.5 h-3.5 text-green-500" />
                        <span className="text-sm text-slate-700">{opt.name}</span>
                      </div>
                      <PriceDisplay price={opt.price} size="sm" prefix="+" className="text-slate-500" />
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </Card>
        )}

        {/* Total */}
        <Card className="p-6 border-0 bg-slate-900 text-white">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-slate-300">Итого</p>
            <PriceDisplay price={totalPrice} size="xl" className="text-white" />
          </div>
        </Card>
      </div>
    </div>
  );
}