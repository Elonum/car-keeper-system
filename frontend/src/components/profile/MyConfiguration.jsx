import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { configuratorService } from '@/services/configuratorService';
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import StatusBadge from '../common/StatusBadge';
import PriceDisplay from '../common/PriceDisplay';
import EmptyState from '../common/EmptyState';
import { Settings, Trash2 } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { toast } from "sonner";

export default function MyConfigurations({ configurations, refetch }) {
  const queryClient = useQueryClient();

  const cancelMutation = useMutation({
    mutationFn: (id) => configuratorService.updateConfiguration(id, { status: 'cancelled' }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-configurations'] });
      toast.success("Конфигурация отменена");
      if (refetch) refetch();
    },
    onError: (error) => {
      const errorMessage = error?.message || error?.data?.error || 'Ошибка при отмене конфигурации';
      toast.error(errorMessage);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => configuratorService.deleteConfiguration(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-configurations'] });
      toast.success("Конфигурация удалена");
      if (refetch) refetch();
    },
    onError: (error) => {
      const errorMessage = error?.message || error?.data?.error || 'Ошибка при удалении конфигурации';
      toast.error(errorMessage);
    },
  });

  const handleCancel = (config) => {
    cancelMutation.mutate(config.configuration_id || config.id);
  };

  const handleDelete = (config) => {
    if (window.confirm('Вы уверены, что хотите удалить эту конфигурацию?')) {
      deleteMutation.mutate(config.configuration_id || config.id);
    }
  };

  if (!configurations || configurations.length === 0) {
    return <EmptyState icon={Settings} title="Нет конфигураций" description="Создайте конфигурацию в каталоге" />;
  }

  return (
    <div className="grid gap-4">
      {configurations.map(config => (
        <Card key={config.configuration_id || config.id} className="p-5 border-0 shadow-sm">
          <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <h3 className="font-bold text-slate-900">{config.trim_name || 'Конфигурация'}</h3>
                <StatusBadge status={config.status} />
              </div>
              <p className="text-sm text-slate-500">Цвет: {config.color_name || '—'}</p>
              {config.option_names && Array.isArray(config.option_names) && config.option_names.length > 0 && (
                <p className="text-xs text-slate-400 mt-1">
                  Опции: {config.option_names.join(', ')}
                </p>
              )}
              <div className="flex items-center gap-4 mt-3">
                <PriceDisplay price={config.total_price} size="md" />
                <span className="text-xs text-slate-400">
                  {config.created_at && format(new Date(config.created_at), 'd MMM yyyy', { locale: ru })}
                </span>
              </div>
            </div>
            <div className="flex gap-2">
              {config.status === 'draft' && (
                <Button 
                  variant="outline" 
                  size="sm" 
                  onClick={() => handleDelete(config)} 
                  disabled={deleteMutation.isPending}
                  className="gap-1.5 rounded-xl text-red-600 hover:text-red-700 hover:bg-red-50"
                >
                  <Trash2 className="w-3.5 h-3.5" /> Удалить
                </Button>
              )}
              {config.status !== 'purchased' && config.status !== 'cancelled' && config.status !== 'draft' && (
                <Button 
                  variant="outline" 
                  size="sm" 
                  onClick={() => handleCancel(config)} 
                  disabled={cancelMutation.isPending}
                  className="gap-1.5 rounded-xl"
                >
                  Отменить
                </Button>
              )}
            </div>
          </div>
        </Card>
      ))}
    </div>
  );
}
