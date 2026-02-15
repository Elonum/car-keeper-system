import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { serviceService } from '@/services/serviceService';
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import StatusBadge from '../common/StatusBadge';
import EmptyState from '../common/EmptyState';
import { Wrench, MapPin, Calendar, X } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { toast } from 'sonner';

export default function ServiceAppointmentsList({ appointments, isLoading }) {
  const queryClient = useQueryClient();

  const cancelMutation = useMutation({
    mutationFn: (id) => serviceService.cancelAppointment(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-appointments'] });
      toast.success('Запись отменена');
    },
    onError: () => {
      toast.error('Ошибка при отмене записи');
    },
  });

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2].map(i => (
          <Card key={i} className="p-6 animate-pulse">
            <div className="h-6 bg-slate-200 rounded w-1/3 mb-4" />
            <div className="h-4 bg-slate-100 rounded w-1/2" />
          </Card>
        ))}
      </div>
    );
  }

  if (!appointments || appointments.length === 0) {
    return (
      <EmptyState 
        icon={Wrench}
        title="Нет записей на ТО"
        description="У вас пока нет записей на сервисное обслуживание"
      />
    );
  }

  return (
    <div className="space-y-4">
      {appointments.map(appt => (
        <Card key={appt.service_appointment_id || appt.id} className="p-6 hover:shadow-md transition-shadow">
          <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
            <div className="flex-1">
              <div className="flex items-start gap-3 mb-3">
                <div className="w-10 h-10 rounded-lg bg-blue-600 flex items-center justify-center flex-shrink-0">
                  <Wrench className="w-5 h-5 text-white" />
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-bold text-slate-900 mb-1">
                    {appt.car_display || 'Сервисная запись'}
                  </h3>
                  {appt.appointment_date && (
                    <div className="flex items-center gap-2 text-sm text-slate-500 mb-2">
                      <Calendar className="w-3.5 h-3.5" />
                      <span>
                        {format(new Date(appt.appointment_date), 'd MMMM yyyy, HH:mm', { locale: ru })}
                      </span>
                    </div>
                  )}
                  <StatusBadge status={appt.status} />
                </div>
              </div>

              {appt.branch_display && (
                <div className="flex items-start gap-2 text-sm text-slate-600 mb-2">
                  <MapPin className="w-4 h-4 mt-0.5 flex-shrink-0" />
                  <span>{appt.branch_display}</span>
                </div>
              )}

              {appt.description && (
                <p className="text-sm text-slate-600 mt-3 p-3 bg-slate-50 rounded-lg">
                  {appt.description}
                </p>
              )}
            </div>

            {appt.status === 'scheduled' && (
              <Button 
                size="sm" 
                variant="ghost" 
                onClick={() => cancelMutation.mutate(appt.service_appointment_id || appt.id)}
                disabled={cancelMutation.isPending}
                className="text-red-600 hover:text-red-700 hover:bg-red-50 gap-1.5"
              >
                <X className="w-3.5 h-3.5" />
                Отменить
              </Button>
            )}
          </div>
        </Card>
      ))}
    </div>
  );
}
