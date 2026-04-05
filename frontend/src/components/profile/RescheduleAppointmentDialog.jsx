import React, { useState, useEffect } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { serviceService } from '@/services/serviceService';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { formatAppointmentInTimeZone, appointmentServiceTypeIds } from '@/lib/appointmentDisplay';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Calendar } from '@/components/ui/calendar';
import { toast } from 'sonner';
import { format, startOfDay, isBefore } from 'date-fns';
import { ru } from 'date-fns/locale';

function calendarDisabled(date) {
  const today = startOfDay(new Date());
  return isBefore(startOfDay(date), today) || date.getDay() === 0;
}

export default function RescheduleAppointmentDialog({ open, onOpenChange, appointment }) {
  const queryClient = useQueryClient();
  const [selectedDate, setSelectedDate] = useState(null);
  const [selectedSlotISO, setSelectedSlotISO] = useState(null);

  const branchId = appointment?.branch_id;
  const ids = appointment ? appointmentServiceTypeIds(appointment) : [];
  const dateKey = selectedDate ? format(selectedDate, 'yyyy-MM-dd') : null;
  const idKey = [...ids].sort().join(',');

  useEffect(() => {
    if (!open) return;
    setSelectedDate(null);
    setSelectedSlotISO(null);
  }, [open, appointment?.service_appointment_id]);

  const {
    data: availability,
    isLoading: slotsLoading,
    isError: slotsError,
    error: slotsQueryError,
  } = useQuery({
    queryKey: ['branch-availability', 'reschedule', branchId, dateKey, idKey],
    queryFn: () =>
      serviceService.getBranchAvailability(branchId, {
        date: dateKey,
        service_type_ids: ids,
      }),
    enabled: Boolean(open && branchId && dateKey && ids.length > 0),
    staleTime: 15 * 1000,
  });

  const slotStarts = Array.isArray(availability?.slot_starts) ? availability.slot_starts : [];
  const branchTimezone = availability?.timezone || 'Europe/Moscow';

  const rescheduleMutation = useMutation({
    mutationFn: (iso) => {
      const id = appointment?.service_appointment_id || appointment?.id;
      if (!id) {
        return Promise.reject(new Error('Не удалось определить запись'));
      }
      return serviceService.rescheduleAppointment(id, { appointment_date: iso });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-appointments'] });
      queryClient.invalidateQueries({ queryKey: ['branch-availability'] });
      toast.success('Запись перенесена');
      onOpenChange(false);
    },
    onError: (err) => {
      toast.error(getApiErrorMessage(err, 'Не удалось перенести запись'));
    },
  });

  const canSubmit = selectedSlotISO && slotStarts.includes(selectedSlotISO);
  const noServices = ids.length === 0;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Перенести запись</DialogTitle>
          <DialogDescription>
            Выберите новую дату и время. Услуги и филиал остаются прежними.
            {appointment?.appointment_date && (
              <span className="mt-2 block text-slate-600">
                Сейчас:{' '}
                {format(new Date(appointment.appointment_date), 'd MMMM yyyy, HH:mm', { locale: ru })}
              </span>
            )}
          </DialogDescription>
        </DialogHeader>

        {noServices ? (
          <p className="text-sm text-amber-800 bg-amber-50 border border-amber-100 rounded-lg px-3 py-2">
            Не удалось определить услуги для этой записи. Обратитесь в сервис.
          </p>
        ) : (
          <div className="grid gap-6 sm:grid-cols-2">
            <div>
              <Label className="mb-2 block text-sm font-medium">Новая дата</Label>
              <Calendar
                mode="single"
                selected={selectedDate}
                onSelect={setSelectedDate}
                disabled={calendarDisabled}
              />
            </div>
            <div>
              <Label className="mb-2 block text-sm font-medium">Время</Label>
              {slotsLoading && <p className="text-sm text-slate-500 py-2">Загрузка слотов…</p>}
              {slotsError && (
                <p className="text-sm text-red-700 bg-red-50 border border-red-100 rounded-lg px-3 py-2">
                  {getApiErrorMessage(slotsQueryError, 'Не удалось загрузить слоты')}
                </p>
              )}
              {!slotsLoading && !slotsError && selectedDate && slotStarts.length === 0 && (
                <p className="text-sm text-amber-800 bg-amber-50 border rounded-lg px-3 py-2">
                  На этот день нет свободных окон.
                </p>
              )}
              {!slotsLoading && !slotsError && slotStarts.length > 0 && (
                <div className="grid grid-cols-3 gap-2">
                  {slotStarts.map((iso) => (
                    <button
                      key={iso}
                      type="button"
                      onClick={() => setSelectedSlotISO(iso)}
                      className={`rounded-xl px-2 py-2 text-sm font-medium transition-all ${
                        selectedSlotISO === iso
                          ? 'bg-slate-900 text-white'
                          : 'border border-slate-200 bg-white text-slate-700 hover:border-slate-300'
                      }`}
                    >
                      {formatAppointmentInTimeZone(iso, branchTimezone, {
                        hour: '2-digit',
                        minute: '2-digit',
                        hour12: false,
                      })}
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        <DialogFooter className="gap-2 sm:gap-0">
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Отмена
          </Button>
          <Button
            type="button"
            disabled={!canSubmit || rescheduleMutation.isPending || noServices}
            onClick={() => canSubmit && rescheduleMutation.mutate(selectedSlotISO)}
            className="bg-blue-600 hover:bg-blue-700"
          >
            {rescheduleMutation.isPending ? 'Сохранение…' : 'Сохранить'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
