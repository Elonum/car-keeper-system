import React from 'react';
import { Link } from 'react-router-dom';
import { createPageUrl } from '../../utils';
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import EmptyState from '../common/EmptyState';
import { Car, Wrench, Gauge, Calendar } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export default function MyCars({ cars }) {
  if (!cars || !Array.isArray(cars) || cars.length === 0) {
    return <EmptyState icon={Car} title="Нет автомобилей" description="У вас пока нет добавленных автомобилей" />;
  }

  return (
    <div className="grid gap-4">
      {cars.map(car => (
        <Card key={car.id} className="p-5 border-0 shadow-sm hover:shadow-md transition-shadow">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 rounded-xl bg-slate-100 flex items-center justify-center flex-shrink-0">
                <Car className="w-6 h-6 text-slate-500" />
              </div>
              <div>
                <h3 className="font-bold text-slate-900">
                  {car.brand_name} {car.model_name}
                </h3>
                <p className="text-sm text-slate-500">{car.trim_name} • {car.color_name}</p>
                <div className="flex flex-wrap gap-3 mt-2 text-xs text-slate-400">
                  <span>VIN: {car.vin}</span>
                  <span>{car.year} г.</span>
                  {car.current_mileage && (
                    <span className="flex items-center gap-1">
                      <Gauge className="w-3 h-3" /> {car.current_mileage.toLocaleString('ru-RU')} км
                    </span>
                  )}
                  {car.purchase_date && (
                    <span className="flex items-center gap-1">
                      <Calendar className="w-3 h-3" />
                      {format(new Date(car.purchase_date), 'd MMM yyyy', { locale: ru })}
                    </span>
                  )}
                </div>
              </div>
            </div>
            <Link to={createPageUrl("ServiceAppointment")}>
              <Button variant="outline" size="sm" className="gap-1.5 rounded-xl">
                <Wrench className="w-3.5 h-3.5" /> Записать на ТО
              </Button>
            </Link>
          </div>
        </Card>
      ))}
    </div>
  );
}