import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { serviceService } from '@/services/serviceService';
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import StatusBadge from '../common/StatusBadge';
import EmptyState from '../common/EmptyState';
import { Wrench, Calendar, MapPin, X } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { toast } from "sonner";

export default function MyServiceAppointments({ appointments, refetch }) {
  const queryClient = useQueryClient();

  const cancelMutation = useMutation({
    mutationFn: (id) => serviceService.cancelAppointment(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-appointments'] });
      toast.success("Запись отменена");
      if (refetch) refetch();
    },
    onError: () => {
      toast.error('Ошибка при отмене записи');
    },
  });

  const handleCancel = (appt) => {
    cancelMutation.mutate(appt.service_appointment_id || appt.id);
  };

  if (!appointments || appointments.length === 0) {
    return <EmptyState icon={Wrench} title="Нет записей на ТО" description="Запишитесь на обслуживание в разделе Сервис" />;
  }

  return (
    <div className="grid gap-4">
      {appointments.map(appt => (
        <Card key={appt.service_appointment_id || appt.id} className="p-5 border-0 shadow-sm">
          <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <StatusBadge status={appt.status} />
              </div>
              <p className="font-semibold text-slate-900">{appt.car_display || 'Автомобиль'}</p>
              <div className="flex flex-wrap gap-3 mt-2 text-sm text-slate-500">
                <span className="flex items-center gap-1.5">
                  <MapPin className="w-3.5 h-3.5" /> {appt.branch_display || '—'}
                </span>
                <span className="flex items-center gap-1.5">
                  <Calendar className="w-3.5 h-3.5" />
                  {appt.appointment_date && format(new Date(appt.appointment_date), 'd MMM yyyy, HH:mm', { locale: ru })}
                </span>
              </div>
              {appt.description && (
                <p className="text-xs text-slate-400 mt-2">{appt.description}</p>
              )}
            </div>
            {appt.status === 'scheduled' && (
              <Button 
                variant="outline" 
                size="sm" 
                onClick={() => handleCancel(appt)} 
                disabled={cancelMutation.isPending}
                className="gap-1.5 rounded-xl text-red-600 hover:text-red-700 hover:bg-red-50"
              >
                <X className="w-3.5 h-3.5" /> Отменить
              </Button>
            )}
          </div>
        </Card>
      ))}
    </div>
  );
}
