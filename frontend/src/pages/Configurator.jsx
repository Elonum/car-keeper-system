import React, { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '@/lib/AuthContext';
import { createPageUrl } from '../utils';
import { catalogService } from '@/services/catalogService';
import { configuratorService } from '@/services/configuratorService';
import { orderService } from '@/services/orderService';
import ColorSelection from '../components/configurator/ColorSelection';
import OptionsSelection from '../components/configurator/OptionsSelection';
import ConfigurationSummary from '../components/configurator/ConfigurationSummary';
import PageLoader from '../components/common/PageLoader';
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ArrowLeft, ArrowRight, Save, ShoppingCart, Check } from 'lucide-react';
import { toast } from 'sonner';

const STEPS = ['Комплектация', 'Цвет', 'Опции', 'Итог'];

export default function Configurator() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user, isAuthenticated, navigateToLogin } = useAuth();
  const [searchParams] = useSearchParams();
  const trimId = searchParams.get('trim_id');
  
  const [currentStep, setCurrentStep] = useState(0);
  const [selectedColorId, setSelectedColorId] = useState(null);
  const [selectedOptionIds, setSelectedOptionIds] = useState([]);

  const { data: trim, isLoading: trimLoading } = useQuery({
    queryKey: ['trim', trimId],
    queryFn: () => catalogService.getTrim(trimId),
    enabled: !!trimId,
  });

  const { data: colors, isLoading: colorsLoading } = useQuery({
    queryKey: ['colors'],
    queryFn: () => configuratorService.getColors(),
  });

  const { data: options, isLoading: optionsLoading } = useQuery({
    queryKey: ['options', trimId],
    queryFn: () => configuratorService.getOptions(trimId),
    enabled: !!trimId,
  });

  const selectedColor = useMemo(() => 
    colors?.find(c => c.color_id === selectedColorId || c.id === selectedColorId),
    [colors, selectedColorId]
  );

  const selectedOptions = useMemo(() => 
    options?.filter(o => selectedOptionIds.includes(o.option_id || o.id)) || [],
    [options, selectedOptionIds]
  );

  const totalPrice = useMemo(() => {
    const basePrice = trim?.base_price || 0;
    const colorDelta = selectedColor?.price_delta || 0;
    const optionsTotal = selectedOptions.reduce((sum, opt) => sum + (opt.price || 0), 0);
    return basePrice + colorDelta + optionsTotal;
  }, [trim, selectedColor, selectedOptions]);

  const saveConfigMutation = useMutation({
    mutationFn: async (status) => {
      if (!isAuthenticated || !user) {
        navigateToLogin();
        return;
      }

      if (!trimId || !selectedColorId) {
        throw new Error('Необходимо выбрать комплектацию и цвет');
      }

      const config = await configuratorService.createConfiguration({
        trim_id: trimId,
        color_id: selectedColorId,
        option_ids: selectedOptionIds,
        status,
        total_price: totalPrice,
      });

      if (status === 'confirmed') {
        await orderService.createOrder({
          configuration_id: config.configuration_id || config.id,
          final_price: totalPrice,
          status: 'pending',
        });
      }

      return { config, status };
    },
    onSuccess: (data) => {
      if (!data) return;
      
      queryClient.invalidateQueries({ queryKey: ['configurations'] });
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      
      if (data.status === 'draft') {
        toast.success('Конфигурация сохранена как черновик');
        navigate(createPageUrl("Profile") + "?tab=configurations");
      } else {
        toast.success('Заказ успешно создан!');
        navigate(createPageUrl("Profile") + "?tab=orders");
      }
    },
    onError: (error) => {
      toast.error(error.message || 'Ошибка при сохранении конфигурации');
    },
  });

  const handleNext = () => {
    if (currentStep === 0 && !trimId) {
      toast.error('Выберите комплектацию');
      return;
    }
    if (currentStep === 1 && !selectedColorId) {
      toast.error('Выберите цвет');
      return;
    }
    setCurrentStep(prev => Math.min(prev + 1, STEPS.length - 1));
  };

  const handleBack = () => {
    setCurrentStep(prev => Math.max(prev - 1, 0));
  };

  const handleToggleOption = (optionId) => {
    setSelectedOptionIds(prev => 
      prev.includes(optionId) 
        ? prev.filter(id => id !== optionId)
        : [...prev, optionId]
    );
  };

  if (trimLoading || colorsLoading || optionsLoading) return <PageLoader />;
  if (!trim) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-slate-50">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-slate-900">Комплектация не найдена</h2>
          <Button className="mt-4" onClick={() => navigate(createPageUrl("Catalog"))}>
            Вернуться в каталог
          </Button>
        </div>
      </div>
    );
  }

  const trimDisplayName = `${trim.brand_name || ''} ${trim.model_name || ''} ${trim.trim_name || trim.name || ''}`.trim();

  return (
    <div className="min-h-screen bg-slate-50">
      {/* Header */}
      <div className="bg-white border-b border-slate-200 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between mb-4">
            <Button variant="ghost" onClick={() => navigate(createPageUrl("Catalog"))} className="gap-2">
              <ArrowLeft className="w-4 h-4" />
              Каталог
            </Button>
            <h1 className="text-xl font-bold text-slate-900 hidden sm:block">
              {trimDisplayName}
            </h1>
          </div>

          {/* Steps */}
          <div className="flex items-center gap-2 overflow-x-auto pb-2">
            {STEPS.map((step, idx) => (
              <React.Fragment key={idx}>
                <button
                  onClick={() => setCurrentStep(idx)}
                  className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all whitespace-nowrap ${
                    idx === currentStep 
                      ? 'bg-slate-900 text-white' 
                      : idx < currentStep
                        ? 'bg-green-100 text-green-700'
                        : 'bg-slate-100 text-slate-500'
                  }`}
                >
                  <span className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${
                    idx === currentStep ? 'bg-white text-slate-900' : idx < currentStep ? 'bg-green-600 text-white' : 'bg-slate-300 text-slate-600'
                  }`}>
                    {idx < currentStep ? <Check className="w-3.5 h-3.5" /> : idx + 1}
                  </span>
                  {step}
                </button>
                {idx < STEPS.length - 1 && (
                  <div className="h-px w-8 bg-slate-200 flex-shrink-0" />
                )}
              </React.Fragment>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid lg:grid-cols-3 gap-8">
          {/* Main content */}
          <div className="lg:col-span-2">
            <Card className="p-6 md:p-8 bg-white shadow-sm rounded-2xl">
              <h2 className="text-2xl font-bold text-slate-900 mb-6">
                {STEPS[currentStep]}
              </h2>

              {currentStep === 0 && (
                <div className="space-y-4">
                  <div className="bg-slate-50 rounded-xl p-6">
                    <h3 className="text-xl font-bold text-slate-900 mb-2">
                      {trimDisplayName}
                    </h3>
                    <p className="text-lg text-slate-700 mb-4">Комплектация: {trim.trim_name || trim.name}</p>
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <p className="text-slate-500">Двигатель</p>
                        <p className="font-semibold">{trim.engine_type_name || trim.engine_type}</p>
                      </div>
                      <div>
                        <p className="text-slate-500">Трансмиссия</p>
                        <p className="font-semibold">{trim.transmission_name || trim.transmission}</p>
                      </div>
                      <div>
                        <p className="text-slate-500">Привод</p>
                        <p className="font-semibold">{trim.drive_type_name || trim.drive_type}</p>
                      </div>
                      <div>
                        <p className="text-slate-500">Базовая цена</p>
                        <p className="font-semibold">{new Intl.NumberFormat('ru-RU').format(trim.base_price || 0)} ₽</p>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {currentStep === 1 && (
                <ColorSelection 
                  colors={colors?.filter(c => c.is_available !== false) || []}
                  selectedColorId={selectedColorId}
                  onSelectColor={(color) => setSelectedColorId(color.color_id || color.id)}
                />
              )}

              {currentStep === 2 && (
                <OptionsSelection 
                  options={options?.filter(o => o.is_available !== false) || []}
                  selectedOptionIds={selectedOptionIds}
                  onToggleOption={handleToggleOption}
                />
              )}

              {currentStep === 3 && (
                <div className="text-center py-8">
                  <div className="w-20 h-20 rounded-full bg-green-100 flex items-center justify-center mx-auto mb-4">
                    <Check className="w-10 h-10 text-green-600" />
                  </div>
                  <h3 className="text-xl font-bold text-slate-900 mb-2">
                    Конфигурация готова!
                  </h3>
                  <p className="text-slate-600 mb-8">
                    Проверьте детали справа и сохраните конфигурацию
                  </p>
                  
                  <div className="flex flex-col sm:flex-row gap-3 justify-center">
                    <Button 
                      variant="outline" 
                      onClick={() => saveConfigMutation.mutate('draft')}
                      disabled={saveConfigMutation.isPending}
                      className="gap-2"
                      size="lg"
                    >
                      <Save className="w-4 h-4" />
                      Сохранить как черновик
                    </Button>
                    <Button 
                      onClick={() => saveConfigMutation.mutate('confirmed')}
                      disabled={saveConfigMutation.isPending}
                      className="gap-2 bg-green-600 hover:bg-green-700"
                      size="lg"
                    >
                      <ShoppingCart className="w-4 h-4" />
                      Оформить заказ
                    </Button>
                  </div>
                </div>
              )}
            </Card>

            {/* Navigation */}
            {currentStep < 3 && (
              <div className="flex justify-between mt-6">
                <Button
                  variant="outline"
                  onClick={handleBack}
                  disabled={currentStep === 0}
                  className="gap-2"
                >
                  <ArrowLeft className="w-4 h-4" />
                  Назад
                </Button>
                <Button
                  onClick={handleNext}
                  className="gap-2 bg-slate-900 hover:bg-slate-800"
                >
                  Далее
                  <ArrowRight className="w-4 h-4" />
                </Button>
              </div>
            )}
          </div>

          {/* Summary sidebar */}
          <div className="lg:col-span-1">
            <div className="sticky top-24">
              <ConfigurationSummary 
                trim={trim}
                color={selectedColor}
                selectedOptions={selectedOptions}
                totalPrice={totalPrice}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
