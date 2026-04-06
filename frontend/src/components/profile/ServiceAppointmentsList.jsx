import React, { useMemo, useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { serviceService } from '@/services/serviceService';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import AppointmentStatusBadge from '../common/AppointmentStatusBadge';
import EmptyState from '../common/EmptyState';
import RescheduleAppointmentDialog from './RescheduleAppointmentDialog';
import CabinetListToolbar from './CabinetListToolbar';
import { filterAppointmentsForCabinet } from '@/lib/cabinetFilters';
import {
  APPOINTMENT_STATUS_LABEL_RU,
  APPOINTMENT_STATUS_CODES,
} from '@/lib/appointmentStatusDisplay';
import { Wrench, MapPin, Calendar, X, CalendarClock } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { toast } from 'sonner';
import { getApiErrorMessage } from '@/lib/apiErrors';

export default function ServiceAppointmentsList({ appointments, isLoading, staffMode = false }) {
  const queryClient = useQueryClient();
  const [rescheduleTarget, setRescheduleTarget] = useState(null);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const statusOptions = useMemo(() => {
    const head = [{ value: 'all', label: 'Все статусы' }];
    const rest = APPOINTMENT_STATUS_CODES.map((code) => ({
      value: code,
      label: APPOINTMENT_STATUS_LABEL_RU[code] || code,
    }));
    return [...head, ...rest];
  }, []);

  const filtered = useMemo(
    () => filterAppointmentsForCabinet(appointments, { search, status: statusFilter }),
    [appointments, search, statusFilter]
  );

  const cancelMutation = useMutation({
    mutationFn: (id) => serviceService.cancelAppointment(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-appointments'] });
      toast.success('Запись отменена');
    },
    onError: (err) => {
      toast.error(getApiErrorMessage(err, 'Ошибка при отмене записи'));
    },
  });

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2].map((i) => (
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
        description={
          staffMode
            ? 'Нет записей или нет доступа к общему списку.'
            : 'У вас пока нет записей на сервисное обслуживание'
        }
      />
    );
  }

  return (
    <div className="space-y-4">
      <CabinetListToolbar
        search={search}
        onSearchChange={setSearch}
        searchPlaceholder="Поиск: VIN, филиал, комментарий, номер записи…"
        statusValue={statusFilter}
        onStatusChange={setStatusFilter}
        statusOptions={statusOptions}
      />

      {filtered.length === 0 ? (
        <p className="text-sm text-slate-500 text-center py-8">
          Нет записей по выбранным условиям. Измените поиск или статус.
        </p>
      ) : (
        filtered.map((appt) => (
          <Card
            key={appt.service_appointment_id || appt.id}
            className="p-6 hover:shadow-md transition-shadow"
          >
            <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
              <div className="flex-1">
                <div className="flex items-start gap-3 mb-3">
                  <div className="w-10 h-10 rounded-lg bg-blue-600 flex items-center justify-center flex-shrink-0">
                    <Wrench className="w-5 h-5 text-white" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="text-lg font-bold text-slate-900 mb-1">
                      {appt.user_car_vin || appt.car_display || 'Сервисная запись'}
                    </h3>
                    {appt.owner_email && (
                      <p className="text-sm text-slate-700 mb-1">
                        <span className="font-medium text-slate-800">Клиент:</span>{' '}
                        {appt.owner_name || '—'} · {appt.owner_email}
                      </p>
                    )}
                    {appt.appointment_date && (
                      <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-slate-500 mb-2">
                        <span className="flex items-center gap-2">
                          <Calendar className="w-3.5 h-3.5 flex-shrink-0" />
                          {format(new Date(appt.appointment_date), 'd MMMM yyyy, HH:mm', {
                            locale: ru,
                          })}
                        </span>
                        {typeof appt.duration_minutes === 'number' && appt.duration_minutes > 0 && (
                          <span className="flex items-center gap-1.5 text-slate-600">
                            <CalendarClock className="w-3.5 h-3.5 flex-shrink-0" />
                            ~{appt.duration_minutes} мин
                          </span>
                        )}
                      </div>
                    )}
                    <AppointmentStatusBadge code={appt.status} />
                  </div>
                </div>

                {(appt.branch_name || appt.branch_display) && (
                  <div className="flex items-start gap-2 text-sm text-slate-600 mb-2">
                    <MapPin className="w-4 h-4 mt-0.5 flex-shrink-0" />
                    <span>
                      {appt.branch_name || appt.branch_display}
                      {appt.branch_address ? `, ${appt.branch_address}` : ''}
                    </span>
                  </div>
                )}

                {Array.isArray(appt.service_types) && appt.service_types.length > 0 && (
                  <div className="mt-3 rounded-lg border border-slate-100 bg-slate-50/80 px-3 py-2">
                    <p className="text-xs font-medium text-slate-500 uppercase tracking-wide mb-1.5">
                      Услуги
                    </p>
                    <ul className="text-sm text-slate-700 space-y-0.5 list-disc list-inside">
                      {appt.service_types.map((st) => (
                        <li key={st.service_type_id || st.id}>{st.name}</li>
                      ))}
                    </ul>
                  </div>
                )}

                {appt.description && (
                  <p className="text-sm text-slate-600 mt-3 p-3 bg-slate-50 rounded-lg">
                    {appt.description}
                  </p>
                )}
              </div>

              {appt.status === 'scheduled' && (
                <div className="flex flex-col sm:flex-row gap-2 shrink-0">
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => setRescheduleTarget(appt)}
                    className="gap-1.5"
                  >
                    <CalendarClock className="w-3.5 h-3.5" />
                    Перенести
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="ghost"
                    onClick={() =>
                      cancelMutation.mutate(appt.service_appointment_id || appt.id)
                    }
                    disabled={cancelMutation.isPending}
                    className="text-red-600 hover:text-red-700 hover:bg-red-50 gap-1.5"
                  >
                    <X className="w-3.5 h-3.5" />
                    Отменить
                  </Button>
                </div>
              )}
            </div>
          </Card>
        ))
      )}

      <RescheduleAppointmentDialog
        open={Boolean(rescheduleTarget)}
        onOpenChange={(open) => {
          if (!open) setRescheduleTarget(null);
        }}
        appointment={rescheduleTarget}
      />
    </div>
  );
}
