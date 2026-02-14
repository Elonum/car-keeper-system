import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/lib/AuthContext';
import { useNavigate } from 'react-router-dom';
import { serviceService } from '@/services/serviceService';
import { createPageUrl } from '../utils';
import StepIndicator from '../components/configurator/StepIndicator';
import CarSelector from '../components/service/CarSelector';
import BranchSelector from '../components/service/BranchSelector';
import PageLoader from '../components/common/PageLoader';
import PriceDisplay from '../components/common/PriceDisplay';
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";
import { ArrowLeft, ArrowRight, Check, Wrench, Car, MapPin, Clock, Package } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Calendar } from "@/components/ui/calendar";
import { Label } from "@/components/ui/label";

const STEPS = ["Автомобиль", "Филиал", "Услуги", "Дата и время", "Подтверждение"];

const timeSlots = [
  "09:00", "09:30", "10:00", "10:30", "11:00", "11:30",
  "12:00", "12:30", "13:00", "13:30", "14:00", "14:30",
  "15:00", "15:30", "16:00", "16:30", "17:00", "17:30",
];

export default function ServiceAppointment() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user, isAuthenticated, navigateToLogin } = useAuth();
  
  const [step, setStep] = useState(0);
  const [selectedCarId, setSelectedCarId] = useState(null);
  const [selectedBranchId, setSelectedBranchId] = useState(null);
  const [selectedServiceIds, setSelectedServiceIds] = useState([]);
  const [selectedDate, setSelectedDate] = useState(null);
  const [selectedTime, setSelectedTime] = useState('');
  const [description, setDescription] = useState('');

  const { data: userCars = [], isLoading: carsLoading } = useQuery({
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

  const selectedCar = userCars.find(c => (c.user_car_id || c.id) === selectedCarId);
  const selectedBranch = branches.find(b => (b.branch_id || b.id) === selectedBranchId);
  const getServiceId = (s) => s.service_type_id || s.id;
  const selectedServices = serviceTypes.filter(s => selectedServiceIds.includes(getServiceId(s)));
  const totalCost = selectedServices.reduce((sum, s) => sum + (s.price || 0), 0);

  const createMutation = useMutation({
    mutationFn: async () => {
      if (!selectedCarId || !selectedBranchId || !selectedDate || !selectedTime) {
        throw new Error('Заполните все обязательные поля');
      }

      const appointmentDate = new Date(selectedDate);
      const [hours, minutes] = selectedTime.split(':');
      appointmentDate.setHours(parseInt(hours), parseInt(minutes));

      return await serviceService.createAppointment({
        user_car_id: selectedCarId,
        branch_id: selectedBranchId,
        appointment_date: appointmentDate.toISOString(),
        service_type_ids: selectedServiceIds,
        status: 'scheduled',
        description: description || undefined,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-appointments'] });
      toast.success("Вы успешно записаны на обслуживание!");
      navigate(createPageUrl("Profile") + "?tab=service");
    },
    onError: (error) => {
      toast.error(error.message || 'Ошибка при создании записи');
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

  const canProceed = () => {
    switch (step) {
      case 0: return !!selectedCarId;
      case 1: return !!selectedBranchId;
      case 2: return selectedServiceIds.length > 0;
      case 3: return !!selectedDate && !!selectedTime;
      default: return true;
    }
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8">
        <div className="mb-8">
          <h1 className="text-2xl md:text-3xl font-bold text-slate-900 tracking-tight">Запись на ТО</h1>
          <p className="text-slate-500 mt-1">Запишите ваш автомобиль на обслуживание</p>
        </div>

        <div className="mb-8">
          <StepIndicator steps={STEPS} currentStep={step} />
        </div>

        {step === 0 && (
          <CarSelector 
            cars={userCars} 
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
                <p className="text-slate-500">Список услуг будет доступен после настройки backend</p>
              </Card>
            ) : (
              <div className="space-y-3">
                {serviceTypes.map(service => {
                  const serviceId = getServiceId(service);
                  const isSelected = selectedServiceIds.includes(serviceId);
                  return (
                    <button
                      key={serviceId}
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
                            {service.duration_minutes && (
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
                          <PriceDisplay price={service.price} size="md" />
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
            <h2 className="text-xl font-bold text-slate-900 mb-1">Выберите дату и время</h2>
            <p className="text-sm text-slate-500 mb-5">Выберите удобные дату и время визита</p>

            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <Label className="text-sm font-medium mb-2 block">Дата</Label>
                <Card className="p-3 border-0 shadow-sm inline-block">
                  <Calendar
                    mode="single"
                    selected={selectedDate}
                    onSelect={setSelectedDate}
                    disabled={(date) => date < new Date() || date.getDay() === 0}
                  />
                </Card>
              </div>
              <div>
                <Label className="text-sm font-medium mb-2 block">Время</Label>
                <div className="grid grid-cols-3 gap-2">
                  {timeSlots.map(time => (
                    <button
                      key={time}
                      onClick={() => setSelectedTime(time)}
                      className={`px-3 py-2.5 rounded-xl text-sm font-medium transition-all ${
                        selectedTime === time
                          ? 'bg-slate-900 text-white'
                          : 'bg-white border border-slate-200 text-slate-700 hover:border-slate-300'
                      }`}
                    >
                      {time}
                    </button>
                  ))}
                </div>

                <div className="mt-6">
                  <Label className="text-sm font-medium mb-2 block">Описание проблемы (необязательно)</Label>
                  <Textarea
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="Опишите, что вас беспокоит..."
                    className="min-h-[100px] rounded-xl"
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
                          <div key={service.id} className="flex items-center justify-between text-sm">
                            <span className="text-slate-700">{service.name}</span>
                            <PriceDisplay price={service.price} size="sm" />
                          </div>
                        ))}
                      </div>
                      <div className="mt-3 pt-3 border-t border-slate-200 flex items-center justify-between font-semibold">
                        <span>Итого</span>
                        <PriceDisplay price={totalCost} size="md" />
                      </div>
                    </div>
                  </div>
                </Card>
              )}

              <Card className="p-5 border-0 shadow-sm">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-purple-50 flex items-center justify-center">
                    <Clock className="w-5 h-5 text-purple-600" />
                  </div>
                  <div>
                    <p className="text-xs text-slate-400 uppercase tracking-wider font-medium">Дата и время</p>
                    <p className="font-semibold text-slate-900">
                      {selectedDate && format(selectedDate, 'd MMMM yyyy', { locale: ru })} в {selectedTime}
                    </p>
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

        <div className="flex items-center justify-between mt-8 pt-6 border-t border-slate-200">
          <Button
            variant="outline"
            onClick={() => setStep(s => s - 1)}
            disabled={step === 0}
            className="gap-2 rounded-xl"
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
