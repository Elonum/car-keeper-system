import React from 'react';
import { Link } from 'react-router-dom';
import { createPageUrl } from '../../utils';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import EmptyState from '../common/EmptyState';
import { Car, Calendar, Gauge, Wrench, Trash2 } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export default function UserCarsList({ cars, isLoading, onRequestDelete, deletingId }) {
  if (isLoading) {
    return (
      <div className="grid md:grid-cols-2 gap-4">
        {[1, 2].map((i) => (
          <Card key={i} className="p-6 animate-pulse">
            <div className="h-6 bg-slate-200 rounded w-2/3 mb-4" />
            <div className="h-4 bg-slate-100 rounded w-1/2" />
          </Card>
        ))}
      </div>
    );
  }

  if (!cars || cars.length === 0) {
    return (
      <EmptyState
        icon={Car}
        title="Нет автомобилей"
        description="Добавьте автомобиль — он появится здесь и будет доступен при записи на ТО"
      />
    );
  }

  return (
    <div className="grid md:grid-cols-2 gap-6">
      {cars.map((car) => {
        const carId = car.user_car_id || car.id;
        if (!carId) return null;
        const isDeleting = deletingId === carId;
        return (
          <Card key={carId} className="p-6 hover:shadow-md transition-shadow flex flex-col">
            {car.image_url && (
              <div className="aspect-video rounded-lg overflow-hidden mb-4 bg-slate-100">
                <img src={car.image_url} alt="" className="w-full h-full object-cover" />
              </div>
            )}

            <div className="flex items-start gap-3 mb-4">
              <div className="w-10 h-10 rounded-lg bg-slate-900 flex items-center justify-center flex-shrink-0">
                <Car className="w-5 h-5 text-white" />
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="text-lg font-bold text-slate-900 mb-1">
                  {car.brand_name} {car.model_name}
                </h3>
                <p className="text-sm text-slate-500 truncate">{car.trim_name}</p>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3 mb-4 text-sm">
              <div>
                <p className="text-slate-500 text-xs mb-0.5">VIN</p>
                <p className="font-mono font-semibold text-slate-900 break-all">{car.vin}</p>
              </div>
              <div>
                <p className="text-slate-500 text-xs mb-0.5">Год</p>
                <p className="font-semibold text-slate-900">{car.year}</p>
              </div>
              {car.current_mileage != null && car.current_mileage !== '' && (
                <div>
                  <p className="text-slate-500 text-xs mb-0.5 flex items-center gap-1">
                    <Gauge className="w-3 h-3" /> Пробег
                  </p>
                  <p className="font-semibold text-slate-900">
                    {Number(car.current_mileage).toLocaleString('ru-RU')} км
                  </p>
                </div>
              )}
              {car.purchase_date && (
                <div>
                  <p className="text-slate-500 text-xs mb-0.5 flex items-center gap-1">
                    <Calendar className="w-3 h-3" /> Куплен
                  </p>
                  <p className="font-semibold text-slate-900">
                    {format(new Date(car.purchase_date), 'd MMM yyyy', { locale: ru })}
                  </p>
                </div>
              )}
            </div>

            {car.color_name && (
              <p className="text-sm text-slate-600 mb-4">
                <span className="font-medium">Цвет:</span>{' '}
                {car.color_hex && (
                  <span
                    className="inline-block h-3 w-3 rounded-full border border-slate-200 align-middle mr-1.5"
                    style={{ backgroundColor: car.color_hex }}
                    aria-hidden
                  />
                )}
                {car.color_name}
              </p>
            )}

            <div className="mt-auto flex flex-col gap-2">
              <Link to={`${createPageUrl('ServiceAppointment')}?car_id=${carId}`}>
                <Button className="w-full gap-2" variant="outline">
                  <Wrench className="w-4 h-4" />
                  Записаться на ТО
                </Button>
              </Link>
              {onRequestDelete && (
                <Button
                  type="button"
                  variant="ghost"
                  className="w-full gap-2 text-red-600 hover:text-red-700 hover:bg-red-50"
                  disabled={isDeleting}
                  onClick={() => onRequestDelete(car)}
                >
                  <Trash2 className="w-4 h-4" />
                  {isDeleting ? 'Удаление…' : 'Удалить из гаража'}
                </Button>
              )}
            </div>
          </Card>
        );
      })}
    </div>
  );
}
