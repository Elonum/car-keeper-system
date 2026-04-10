import React, { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/lib/AuthContext';
import { useNavigate } from 'react-router-dom';
import { serviceService } from '@/services/serviceService';
import { createPageUrl } from '../utils';
import { pagesConfig } from '../pages.config';
import { formatAppointmentInTimeZone } from '@/lib/appointmentDisplay';
import StepIndicator from '../components/configurator/StepIndicator';
import CarSelector from '../components/service/CarSelector';
import BranchSelector from '../components/service/BranchSelector';
import PageLoader from '../components/common/PageLoader';
import { ErrorNotice } from '@/components/common/ErrorNotice';
import PriceDisplay from '../components/common/PriceDisplay';
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { UI_LIMITS } from '@/lib/authValidation';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { ArrowLeft, ArrowRight, Check, Wrench, Car, MapPin, Clock, Package, CalendarDays } from 'lucide-react';
import { format, startOfDay, isBefore } from 'date-fns';
import { Calendar } from "@/components/ui/calendar";
import { Label } from "@/components/ui/label";

const STEPS = ["Автомобиль", "Филиал", "Услуги", "Дата", "Подтверждение"];

function truncateDescription(s, maxChars) {
  const arr = Array.from(String(s ?? ''));
  if (arr.length <= maxChars) return s || undefined;
  return arr.slice(0, maxChars).join('');
}

export default function ServiceAppointment() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { isAuthenticated, navigateToLogin } = useAuth();
  
  const [step, setStep] = useState(0);
  const [selectedCarId, setSelectedCarId] = useState(null);
  const [selectedBranchId, setSelectedBranchId] = useState(null);
  const [selectedServiceIds, setSelectedServiceIds] = useState([]);
  const [selectedDate, setSelectedDate] = useState(null);
  const [selectedSlotISO, setSelectedSlotISO] = useState(null);
  const [description, setDescription] = useState('');
  const [formError, setFormError] = useState(null);

  const dateKey = selectedDate ? format(selectedDate, 'yyyy-MM-dd') : null;
  const serviceKey = [...selectedServiceIds].sort().join(',');

  const {
    data: availability,
    isLoading: slotsLoading,
    isError: slotsError,
    error: slotsQueryError,
  } = useQuery({
    queryKey: ['branch-availability', selectedBranchId, dateKey, serviceKey],
    queryFn: () =>
      serviceService.getBranchAvailability(selectedBranchId, {
        date: dateKey,
        service_type_ids: selectedServiceIds,
      }),
    enabled:
      Boolean(isAuthenticated) &&
      step === 3 &&
      !!selectedBranchId &&
      !!dateKey &&
      selectedServiceIds.length > 0,
    staleTime: 30 * 1000,
  });

  const slotStarts = Array.isArray(availability?.slot_starts) ? availability.slot_starts : [];
  const branchTimezone = availability?.timezone || 'Europe/Moscow';
  const plannedDurationMin = availability?.duration_minutes ?? null;

  useEffect(() => {
    setSelectedSlotISO(null);
  }, [selectedBranchId, dateKey, serviceKey]);

  useEffect(() => {
    if (!selectedSlotISO || !slotStarts.length) return;
    const still = slotStarts.some((s) => s === selectedSlotISO);
    if (!still) setSelectedSlotISO(null);
  }, [slotStarts, selectedSlotISO]);

  const { data: userCars, isLoading: carsLoading } = useQuery({
    queryKey: ['userCars'],
    queryFn: () => serviceService.getUserCars(),
    enabled: isAuthenticated,
  });

  const { data: branches = [], isLoading: branchesLoading } = useQuery({
    queryKey: ['branches'],
    queryFn: () => serviceService.getBranches({ is_active: true }),
  });

  const { data: serviceTypes = [], isLoading: servicesLoading } = useQuery({
    queryKey: ['serviceTypes', { is_available: true }],
    queryFn: () => serviceService.getServiceTypes({ is_available: true }),
  });

  const safeUserCars = Array.isArray(userCars) ? userCars : [];
  const selectedCar = safeUserCars.find(c => (c.user_car_id || c.id) === selectedCarId);
  const selectedBranch = branches.find(b => (b.branch_id || b.id) === selectedBranchId);
  const getServiceId = (s) => s.service_type_id || s.id;
  const selectedServices = serviceTypes.filter(s => selectedServiceIds.includes(getServiceId(s)));
  const totalCost = selectedServices.reduce((sum, s) => sum + (s.price || 0), 0);

  const createMutation = useMutation({
    mutationFn: async () => {
      if (!selectedCarId || !selectedBranchId || !selectedSlotISO) {
        throw new Error('Заполните все обязательные поля');
      }
      if (!slotStarts.includes(selectedSlotISO)) {
        throw new Error('Выберите доступное время из списка');
      }

      const desc = truncateDescription(description.trim(), UI_LIMITS.APPOINTMENT_DESCRIPTION);

      return await serviceService.createAppointment({
        user_car_id: selectedCarId,
        branch_id: selectedBranchId,
        appointment_date: selectedSlotISO,
        service_type_ids: selectedServiceIds,
        description: desc || undefined,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-appointments'] });
      queryClient.invalidateQueries({ queryKey: ['branch-availability'] });
      setFormError(null);
      navigate(createPageUrl("Profile") + "?tab=service");
    },
    onError: (error) => {
      setSelectedSlotISO(null);
      queryClient.invalidateQueries({ queryKey: ['branch-availability'] });
      setFormError(getApiErrorMessage(error, 'Ошибка при создании записи'));
    },
  });

  if (!isAuthenticated) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[70vh] px-4">
        <Wrench className="w-16 h-16 text-slate-300 mb-4" />
        <h2 className="text-xl font-bold text-slate-900 mb-2">Войдите для записи на ТО</h2>
        <p className="text-slate-500 mb-6 text-center">Необходимо авторизоваться для записи на обслуживание</p>
        <Button onClick={navigateToLogin} className="bg-slate-900 hover:bg-blue-600">
          Войти
        </Button>
      </div>
    );
  }

  if (carsLoading || branchesLoading || servicesLoading) return <PageLoader />;

  if (!carsLoading && safeUserCars.length === 0) {
    return (
      <div className="min-h-screen bg-slate-50">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8">
          <div className="mb-8">
            <h1 className="text-2xl md:text-3xl font-bold text-slate-900 tracking-tight">Запись на ТО</h1>
            <p className="text-slate-500 mt-1">Запишите ваш автомобиль на обслуживание</p>
          </div>
          <div className="text-center py-12 bg-white rounded-2xl shadow-sm border border-slate-100">
            <Car className="w-16 h-16 text-slate-300 mx-auto mb-4" />
            <h2 className="text-xl font-bold text-slate-900 mb-2">У вас нет автомобилей</h2>
            <p className="text-slate-500 mb-6">Для записи на обслуживание необходимо добавить автомобиль в личном кабинете</p>
            <Button 
              onClick={() => navigate(createPageUrl("Profile") + "?tab=cars")}
              className="bg-slate-900 hover:bg-slate-800"
            >
              Перейти в личный кабинет
            </Button>
          </div>
        </div>
      </div>
    );
  }

  const canProceed = () => {
    switch (step) {
      case 0: return !!selectedCarId;
      case 1: return !!selectedBranchId;
      case 2: return selectedServiceIds.length > 0;
      case 3: return !!selectedDate && !!selectedSlotISO && slotStarts.includes(selectedSlotISO);
      default: return true;
    }
  };

  const isCalendarDisabled = (date) => {
    const today = startOfDay(new Date());
    return isBefore(startOfDay(date), today) || date.getDay() === 0;
  };

  const slotsEmptyMessage = () => {
    if (plannedDurationMin != null && selectedDate) {
      const b = branches.find((x) => (x.branch_id || x.id) === selectedBranchId);
      const startM = b?.workday_start_minutes ?? 540;
      const endM = b?.workday_end_minutes ?? 1080;
      if (plannedDurationMin > endM - startM) {
        return 'Суммарная длительность выбранных услуг не помещается в рабочий день филиала. Уберите часть услуг или выберите другой филиал.';
      }
    }
    return 'На этот день нет свободных окон. Выберите другую дату.';
  };

  const goToServices = () => {
    navigate(createPageUrl('Services'));
  };

  /** Предыдущий шаг мастера; на шаге 0 — шаг назад в истории браузера (не то же самое, что «К услугам»). */
  const wizardBack = () => {
    if (step > 0) {
      setStep((s) => s - 1);
      return;
    }
    if (typeof window !== 'undefined' && window.history.length > 1) {
      navigate(-1);
      return;
    }
    navigate(createPageUrl(pagesConfig.mainPage));
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8 pb-32">
        <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <h1 className="text-2xl md:text-3xl font-bold text-slate-900 tracking-tight">Запись на ТО</h1>
            <p className="text-slate-500 mt-1">Запишите ваш автомобиль на обслуживание</p>
          </div>
          <Button
            type="button"
            variant="outline"
            className="shrink-0 gap-2 rounded-xl self-start sm:self-center"
            onClick={goToServices}
          >
            <Wrench className="h-4 w-4" />
            К услугам
          </Button>
        </div>

        <div className="mb-8">
          <StepIndicator steps={STEPS} currentStep={step} />
        </div>
        <ErrorNotice kind="server" message={formError} className="mb-6" />

        {step === 0 && (
          <CarSelector 
            cars={safeUserCars} 
            selectedCarId={selectedCarId} 
            onSelect={setSelectedCarId} 
          />
        )}
        
        {step === 1 && (
          <BranchSelector 
            branches={branches} 
            selectedBranchId={selectedBranchId} 
            onSelect={setSelectedBranchId} 
          />
        )}
        
        {step === 2 && (
          <div>
            <h2 className="text-xl font-bold text-slate-900 mb-1">Выберите услуги</h2>
            <p className="text-sm text-slate-500 mb-5">Отметьте необходимые сервисные работы</p>

            {serviceTypes.length === 0 ? (
              <Card className="p-6 text-center">
                <p className="text-slate-500">Список услуг временно недоступен</p>
              </Card>
            ) : (
              <div className="space-y-3">
                {serviceTypes.map(service => {
                  const serviceId = getServiceId(service);
                  const isSelected = selectedServiceIds.includes(serviceId);
                  return (
                    <button
                      key={serviceId}
                      type="button"
                      onClick={() => {
                        setSelectedServiceIds(prev => 
                          prev.includes(serviceId)
                            ? prev.filter(id => id !== serviceId)
                            : [...prev, serviceId]
                        );
                      }}
                      className={`w-full text-left p-4 rounded-xl border-2 transition-all ${
                        isSelected 
                          ? 'border-slate-900 bg-slate-50' 
                          : 'border-slate-200 bg-white hover:border-slate-300'
                      }`}
                    >
                      <div className="flex items-start justify-between gap-4">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="font-semibold text-slate-900">{service.name}</h3>
                            {service.duration_minutes != null && (
                              <span className="text-xs text-slate-400 flex items-center gap-1">
                                <Clock className="w-3 h-3" /> {service.duration_minutes} мин
                              </span>
                            )}
                          </div>
                          {service.description && (
                            <p className="text-sm text-slate-600">{service.description}</p>
                          )}
                        </div>
                        <div className="flex flex-col items-end gap-2">
                          <PriceDisplay price={service.price} size="md" className="" />
                          {isSelected && (
                            <div className="w-6 h-6 rounded-full bg-slate-900 flex items-center justify-center">
                              <Check className="w-3.5 h-3.5 text-white" />
                            </div>
                          )}
                        </div>
                      </div>
                    </button>
                  );
                })}
              </div>
            )}

            {selectedServiceIds.length > 0 && (
              <Card className="p-4 mt-6 bg-slate-900 text-white border-0">
                <div className="flex items-center justify-between">
                  <span className="font-medium">Итого ({selectedServiceIds.length} услуг)</span>
                  <PriceDisplay price={totalCost} size="lg" className="text-white" />
                </div>
              </Card>
            )}
          </div>
        )}
        
        {step === 3 && (
          <div>
            <h2 className="text-xl font-bold text-slate-900 mb-1">Выберите дату</h2>
            <p className="text-sm text-slate-500 mb-5">
              Укажите дату и время визита. Воскресенье — выходной. Слоты рассчитываются для выбранного филиала с учётом загрузки.
            </p>

            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <Label className="text-sm font-medium mb-2 block">Дата</Label>
                <Card className="p-3 border-0 shadow-sm inline-block">
                  <Calendar
                    mode="single"
                    selected={selectedDate}
                    onSelect={setSelectedDate}
                    disabled={isCalendarDisabled}
                  />
                </Card>
              </div>
              <div>
                <Label className="text-sm font-medium mb-2 block">Время</Label>
                {slotsLoading && (
                  <p className="text-sm text-slate-500 py-4">Загрузка доступных слотов…</p>
                )}
                {slotsError && (
                  <ErrorNotice
                    kind="server"
                    message={getApiErrorMessage(slotsQueryError, 'Не удалось загрузить слоты')}
                  />
                )}
                {!slotsLoading && !slotsError && selectedDate && selectedServiceIds.length > 0 && slotStarts.length === 0 && (
                  <p className="text-sm text-amber-800 bg-amber-50 border border-amber-100 rounded-xl px-4 py-3">
                    {slotsEmptyMessage()}
                  </p>
                )}
                {!slotsLoading && !slotsError && slotStarts.length > 0 && (
                  <div className="grid grid-cols-3 gap-2">
                    {slotStarts.map((iso) => (
                      <button
                        key={iso}
                        type="button"
                        onClick={() => setSelectedSlotISO(iso)}
                        className={`px-3 py-2.5 rounded-xl text-sm font-medium transition-all ${
                          selectedSlotISO === iso
                            ? 'bg-slate-900 text-white'
                            : 'bg-white border border-slate-200 text-slate-700 hover:border-slate-300'
                        }`}
                      >
                        {formatAppointmentInTimeZone(iso, branchTimezone, { hour: '2-digit', minute: '2-digit', hour12: false })}
                      </button>
                    ))}
                  </div>
                )}
                {plannedDurationMin != null && !slotsLoading && (
                  <p className="text-xs text-slate-500 mt-3">
                    Ориентировочная длительность работ: {plannedDurationMin} мин (часовой пояс филиала: {branchTimezone})
                  </p>
                )}

                <div className="mt-6">
                  <Label className="text-sm font-medium mb-2 block">Описание проблемы (необязательно)</Label>
                  <Textarea
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="Опишите, что вас беспокоит..."
                    className="min-h-[100px] rounded-xl"
                    maxLength={UI_LIMITS.APPOINTMENT_DESCRIPTION}
                  />
                </div>
              </div>
            </div>
          </div>
        )}

        {step === 4 && (
          <div>
            <h2 className="text-xl font-bold text-slate-900 mb-1">Подтверждение</h2>
            <p className="text-sm text-slate-500 mb-5">Проверьте данные записи</p>

            <div className="space-y-3">
              <Card className="p-5 border-0 shadow-sm">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center">
                    <Car className="w-5 h-5 text-blue-600" />
                  </div>
                  <div>
                    <p className="text-xs text-slate-400 uppercase tracking-wider font-medium">Автомобиль</p>
                    <p className="font-semibold text-slate-900">
                      {selectedCar?.vin || 'Автомобиль'}
                    </p>
                  </div>
                </div>
              </Card>

              <Card className="p-5 border-0 shadow-sm">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-green-50 flex items-center justify-center">
                    <MapPin className="w-5 h-5 text-green-600" />
                  </div>
                  <div>
                    <p className="text-xs text-slate-400 uppercase tracking-wider font-medium">Филиал</p>
                    <p className="font-semibold text-slate-900">{selectedBranch?.name}</p>
                    <p className="text-xs text-slate-500">{selectedBranch?.address}</p>
                  </div>
                </div>
              </Card>

              {selectedServices.length > 0 && (
                <Card className="p-5 border-0 shadow-sm">
                  <div className="flex items-start gap-3">
                    <div className="w-10 h-10 rounded-xl bg-orange-50 flex items-center justify-center">
                      <Package className="w-5 h-5 text-orange-600" />
                    </div>
                    <div className="flex-1">
                      <p className="text-xs text-slate-400 uppercase tracking-wider font-medium mb-2">Услуги</p>
                      <div className="space-y-1.5">
                        {selectedServices.map(service => (
                          <div key={getServiceId(service)} className="flex items-center justify-between text-sm">
                            <span className="text-slate-700">{service.name}</span>
                            <PriceDisplay price={service.price} size="sm" className="" />
                          </div>
                        ))}
                      </div>
                      <div className="mt-3 pt-3 border-t border-slate-200 flex items-center justify-between font-semibold">
                        <span>Итого</span>
                        <PriceDisplay price={totalCost} size="md" className="" />
                      </div>
                    </div>
                  </div>
                </Card>
              )}

              <Card className="p-5 border-0 shadow-sm">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-purple-50 flex items-center justify-center">
                      <CalendarDays className="w-5 h-5 text-purple-600" />
                    </div>
                    <div>
                      <p className="text-xs text-slate-400 uppercase tracking-wider font-medium">Дата</p>
                      <p className="font-semibold text-slate-900">
                        {selectedSlotISO &&
                          formatAppointmentInTimeZone(selectedSlotISO, branchTimezone, {
                            day: 'numeric',
                            month: 'long',
                            year: 'numeric',
                          })}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-slate-100 flex items-center justify-center">
                      <Clock className="w-5 h-5 text-slate-600" />
                    </div>
                    <div>
                      <p className="text-xs text-slate-400 uppercase tracking-wider font-medium">Время</p>
                      <p className="font-semibold text-slate-900">
                        {selectedSlotISO &&
                          formatAppointmentInTimeZone(selectedSlotISO, branchTimezone, {
                            hour: '2-digit',
                            minute: '2-digit',
                            hour12: false,
                          })}
                      </p>
                      {plannedDurationMin != null && (
                        <p className="text-xs text-slate-500 mt-1">≈ {plannedDurationMin} мин работ</p>
                      )}
                    </div>
                  </div>
                </div>
              </Card>

              {description && (
                <Card className="p-5 border-0 shadow-sm">
                  <p className="text-xs text-slate-400 uppercase tracking-wider font-medium mb-1">Описание</p>
                  <p className="text-sm text-slate-700">{description}</p>
                </Card>
              )}
            </div>
          </div>
        )}

        <div className="sticky bottom-0 z-20 -mx-4 flex items-center justify-between gap-4 border-t border-slate-200 bg-slate-50/95 px-4 py-4 backdrop-blur supports-[backdrop-filter]:bg-slate-50/90 sm:-mx-6 sm:px-6 mt-8">
          <Button
            type="button"
            variant="outline"
            onClick={wizardBack}
            className="gap-2 rounded-xl min-w-[7.5rem]"
          >
            <ArrowLeft className="w-4 h-4" /> Назад
          </Button>

          {step === 4 ? (
            <Button
              onClick={() => createMutation.mutate()}
              disabled={createMutation.isPending}
              className="bg-blue-600 hover:bg-blue-700 gap-2 rounded-xl"
            >
              <Check className="w-4 h-4" /> Записаться
            </Button>
          ) : (
            <Button
              onClick={() => setStep(s => s + 1)}
              disabled={!canProceed()}
              className="bg-slate-900 hover:bg-slate-800 gap-2 rounded-xl"
            >
              Далее <ArrowRight className="w-4 h-4" />
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
