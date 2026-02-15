import React from 'react';
import { cn } from "@/lib/utils";
import { Card } from "@/components/ui/card";
import { Check, Car, Gauge } from 'lucide-react';

export default function CarSelector({ cars, selectedCarId, onSelect }) {
  if (!cars || !Array.isArray(cars) || cars.length === 0) {
    return (
      <div className="text-center py-12">
        <Car className="w-12 h-12 text-slate-300 mx-auto mb-3" />
        <h3 className="font-semibold text-slate-900 mb-1">У вас нет автомобилей</h3>
        <p className="text-sm text-slate-500">Добавьте автомобиль в личном кабинете</p>
      </div>
    );
  }

  return (
    <div>
      <h2 className="text-xl font-bold text-slate-900 mb-1">Выберите автомобиль</h2>
      <p className="text-sm text-slate-500 mb-5">Выберите авто для записи на обслуживание</p>
      <div className="grid gap-3">
        {cars.map(car => {
          const carId = car.user_car_id || car.id;
          const isSelected = selectedCarId === carId;
          return (
            <button
              key={carId}
              onClick={() => onSelect(carId)}
              className={cn(
                "flex items-center gap-4 p-4 rounded-xl border-2 text-left transition-all w-full",
                isSelected ? "border-slate-900 bg-slate-50" : "border-slate-100 bg-white hover:border-slate-200"
              )}
            >
              {isSelected && (
                <div className="w-6 h-6 rounded-full bg-slate-900 flex items-center justify-center flex-shrink-0">
                  <Check className="w-3.5 h-3.5 text-white" />
                </div>
              )}
              {!isSelected && <div className="w-6 h-6 rounded-full border-2 border-slate-200 flex-shrink-0" />}
              <div className="flex-1 min-w-0">
                <p className="font-semibold text-slate-900">
                  {car.brand_name} {car.model_name} {car.trim_name}
                </p>
                <div className="flex items-center gap-3 mt-1 text-xs text-slate-500">
                  <span>{car.year} г.</span>
                  <span>VIN: {car.vin}</span>
                  {car.current_mileage && (
                    <span className="flex items-center gap-1">
                      <Gauge className="w-3 h-3" />
                      {car.current_mileage.toLocaleString('ru-RU')} км
                    </span>
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