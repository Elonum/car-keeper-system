import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { configuratorService } from '@/services/configuratorService';
import { orderService } from '@/services/orderService';
import { Link } from 'react-router-dom';
import { createPageUrl } from '../../utils';
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import StatusBadge from '../common/StatusBadge';
import PriceDisplay from '../common/PriceDisplay';
import EmptyState from '../common/EmptyState';
import { Settings, Edit, Trash2, ShoppingCart } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { toast } from 'sonner';

export default function ConfigurationList({ configurations, isLoading }) {
  const queryClient = useQueryClient();

  const deleteMutation = useMutation({
    mutationFn: (id) => configuratorService.deleteConfiguration(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-configurations'] });
      toast.success('Конфигурация удалена');
    },
    onError: () => {
      toast.error('Ошибка при удалении конфигурации');
    },
  });

  const createOrderMutation = useMutation({
    mutationFn: async (config) => {
      await configuratorService.updateConfiguration(config.configuration_id || config.id, { status: 'confirmed' });
      await orderService.createOrder({
        configuration_id: config.configuration_id || config.id,
        final_price: config.total_price,
        status: 'pending',
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-configurations'] });
      queryClient.invalidateQueries({ queryKey: ['my-orders'] });
      toast.success('Заказ создан!');
    },
    onError: () => {
      toast.error('Ошибка при создании заказа');
    },
  });

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map(i => (
          <Card key={i} className="p-6 animate-pulse">
            <div className="h-6 bg-slate-200 rounded w-1/3 mb-4" />
            <div className="h-4 bg-slate-100 rounded w-1/2" />
          </Card>
        ))}
      </div>
    );
  }

  if (!configurations || configurations.length === 0) {
    return (
      <EmptyState 
        icon={Settings}
        title="Нет конфигураций"
        description="Вы ещё не создали ни одной конфигурации автомобиля"
        action={
          <Link to={createPageUrl("Catalog")}>
            <Button className="gap-2">
              <Settings className="w-4 h-4" />
              Перейти в каталог
            </Button>
          </Link>
        }
      />
    );
  }

  return (
    <div className="space-y-4">
      {configurations.map(config => (
        <Card key={config.configuration_id || config.id} className="p-6 hover:shadow-md transition-shadow">
          <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
            <div className="flex-1">
              <div className="flex items-start gap-3 mb-3">
                <div className="w-10 h-10 rounded-lg bg-slate-900 flex items-center justify-center flex-shrink-0">
                  <Settings className="w-5 h-5 text-white" />
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-bold text-slate-900 mb-1">
                    {config.trim_name || 'Конфигурация'}
                  </h3>
                  <p className="text-sm text-slate-500 mb-2">
                    Создано: {format(new Date(config.created_at || config.created_date), 'd MMMM yyyy', { locale: ru })}
                  </p>
                  <StatusBadge status={config.status} />
                </div>
              </div>

              {config.color_name && (
                <p className="text-sm text-slate-600 mb-1">
                  <span className="font-medium">Цвет:</span> {config.color_name}
                </p>
              )}
              {config.option_names && config.option_names.length > 0 && (
                <p className="text-sm text-slate-600">
                  <span className="font-medium">Опции:</span> {config.option_names.join(', ')}
                </p>
              )}
            </div>

            <div className="flex flex-col items-end gap-3">
              <PriceDisplay price={config.total_price} size="lg" />
              
              <div className="flex flex-wrap gap-2">
                {config.status === 'draft' && (
                  <>
                    <Link to={createPageUrl("Configurator") + `?trim_id=${config.trim_id}&config_id=${config.configuration_id || config.id}`}>
                      <Button size="sm" variant="outline" className="gap-1.5">
                        <Edit className="w-3.5 h-3.5" />
                        Редактировать
                      </Button>
                    </Link>
                    <Button 
                      size="sm" 
                      onClick={() => createOrderMutation.mutate(config)}
                      disabled={createOrderMutation.isPending}
                      className="gap-1.5 bg-green-600 hover:bg-green-700"
                    >
                      <ShoppingCart className="w-3.5 h-3.5" />
                      Заказать
                    </Button>
                  </>
                )}
                <Button 
                  size="sm" 
                  variant="ghost" 
                  onClick={() => deleteMutation.mutate(config.configuration_id || config.id)}
                  disabled={deleteMutation.isPending}
                  className="text-red-600 hover:text-red-700 hover:bg-red-50"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                </Button>
              </div>
            </div>
          </div>
        </Card>
      ))}
    </div>
  );
}
