import React from 'react';
import { Card } from "@/components/ui/card";
import StatusBadge from '../common/StatusBadge';
import PriceDisplay from '../common/PriceDisplay';
import EmptyState from '../common/EmptyState';
import { ShoppingCart } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export default function MyOrders({ orders }) {
  if (!orders || !Array.isArray(orders) || orders.length === 0) {
    return <EmptyState icon={ShoppingCart} title="Нет заказов" description="Оформите заказ через конфигуратор" />;
  }

  return (
    <div className="grid gap-4">
      {orders.map(order => (
        <Card key={order.id} className="p-5 border-0 shadow-sm">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <h3 className="font-bold text-slate-900">Заказ #{order.id?.slice(-6)}</h3>
                <StatusBadge status={order.status} />
              </div>
              <p className="text-sm text-slate-500">{order.configuration_summary || '—'}</p>
              {order.manager_email && (
                <p className="text-xs text-slate-400 mt-1">Менеджер: {order.manager_email}</p>
              )}
              <div className="flex items-center gap-3 mt-2">
                <span className="text-xs text-slate-400">
                  {order.created_date && format(new Date(order.created_date), 'd MMM yyyy', { locale: ru })}
                </span>
              </div>
            </div>
            <PriceDisplay price={order.final_price} size="lg" />
          </div>
        </Card>
      ))}
    </div>
  );
}