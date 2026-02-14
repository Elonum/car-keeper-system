import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { serviceService } from '@/services/serviceService';
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import PageLoader from '../common/PageLoader';
import PriceDisplay from '../common/PriceDisplay';
import { Calendar, MapPin, Wrench, FileText, CheckCircle } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export default function CarServiceHistory({ carId }) {
  const { data: appointments = [], isLoading } = useQuery({
    queryKey: ['carServiceHistory', carId],
    queryFn: () => serviceService.getAppointments({ user_car_id: carId }),
    enabled: !!carId,
  });

  if (isLoading) return <PageLoader />;

  if (appointments.length === 0) {
    return (
      <div className="text-center py-12">
        <Wrench className="w-12 h-12 text-slate-300 mx-auto mb-3" />
        <p className="text-slate-500">Нет истории обслуживания</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h3 className="font-semibold text-slate-900 mb-4">История обслуживания</h3>
      
      {appointments.map(appt => (
        <Card key={appt.service_appointment_id || appt.id} className="p-5 border-0 shadow-sm">
          <div className="flex items-start justify-between gap-4 mb-3">
            <div className="flex items-center gap-2">
              {appt.status === 'completed' ? (
                <div className="w-8 h-8 rounded-lg bg-green-50 flex items-center justify-center">
                  <CheckCircle className="w-4 h-4 text-green-600" />
                </div>
              ) : appt.status === 'cancelled' ? (
                <div className="w-8 h-8 rounded-lg bg-red-50 flex items-center justify-center">
                  <Wrench className="w-4 h-4 text-red-600" />
                </div>
              ) : (
                <div className="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center">
                  <Wrench className="w-4 h-4 text-blue-600" />
                </div>
              )}
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-slate-900">
                    {appt.appointment_date && format(new Date(appt.appointment_date), 'd MMM yyyy', { locale: ru })}
                  </span>
                  {appt.status === 'completed' && (
                    <Badge className="bg-green-100 text-green-700 border-green-200 border text-xs">
                      Завершено
                    </Badge>
                  )}
                  {appt.status === 'cancelled' && (
                    <Badge className="bg-red-100 text-red-700 border-red-200 border text-xs">
                      Отменено
                    </Badge>
                  )}
                  {appt.status === 'scheduled' && (
                    <Badge className="bg-blue-100 text-blue-700 border-blue-200 border text-xs">
                      Запланировано
                    </Badge>
                  )}
                </div>
                <p className="text-xs text-slate-500 mt-0.5 flex items-center gap-1">
                  <MapPin className="w-3 h-3" /> {appt.branch_display || appt.branch_name || '—'}
                </p>
              </div>
            </div>
            
            {appt.total_cost > 0 && (
              <PriceDisplay price={appt.total_cost} size="md" />
            )}
          </div>

          {appt.service_names && appt.service_names.length > 0 && (
            <div className="mb-3">
              <p className="text-xs text-slate-400 mb-1.5">Услуги:</p>
              <div className="flex flex-wrap gap-1.5">
                {appt.service_names.map((name, idx) => (
                  <Badge key={idx} variant="outline" className="text-xs">
                    {name}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          {appt.description && (
            <p className="text-sm text-slate-600 mb-2">
              <span className="font-medium">Описание: </span>
              {appt.description}
            </p>
          )}

          {appt.notes && appt.status === 'completed' && (
            <div className="mt-3 pt-3 border-t border-slate-100">
              <div className="flex items-start gap-2">
                <FileText className="w-4 h-4 text-slate-400 mt-0.5" />
                <div>
                  <p className="text-xs text-slate-400 mb-1">Отчёт мастера:</p>
                  <p className="text-sm text-slate-700">{appt.notes}</p>
                </div>
              </div>
            </div>
          )}

          {appt.completed_date && (
            <p className="text-xs text-slate-400 mt-2">
              Завершено: {format(new Date(appt.completed_date), 'd MMM yyyy, HH:mm', { locale: ru })}
            </p>
          )}
        </Card>
      ))}
    </div>
  );
}
