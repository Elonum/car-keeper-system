import React, { useMemo, useState } from 'react';
import { Card } from '@/components/ui/card';
import OrderStatusBadge from '../common/OrderStatusBadge';
import PriceDisplay from '../common/PriceDisplay';
import { useOrderStatusLabelMap } from '@/hooks/useOrderStatusLabelMap';
import { resolveOrderStatusLabel } from '@/lib/orderStatusDisplay';
import { filterOrdersForCabinet } from '@/lib/cabinetFilters';
import CabinetListToolbar from './CabinetListToolbar';
import EmptyState from '../common/EmptyState';
import { ShoppingCart } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export default function OrdersList({ orders, isLoading }) {
  const { labelByCode, data: orderStatusRows } = useOrderStatusLabelMap();
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');

  const statusOptions = useMemo(() => {
    const head = [{ value: 'all', label: 'Все статусы' }];
    const rows = Array.isArray(orderStatusRows) ? orderStatusRows : [];
    const rest = rows.map((s) => ({
      value: s.code,
      label: s.customer_label_ru || s.code,
    }));
    return [...head, ...rest];
  }, [orderStatusRows]);

  const filteredOrders = useMemo(
    () => filterOrdersForCabinet(orders, { search, status: statusFilter }),
    [orders, search, statusFilter]
  );

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

  if (!orders || orders.length === 0) {
    return (
      <EmptyState
        icon={ShoppingCart}
        title="Нет заказов"
        description="У вас пока нет оформленных заказов"
      />
    );
  }

  return (
    <div className="space-y-4">
      <CabinetListToolbar
        search={search}
        onSearchChange={setSearch}
        searchPlaceholder="Поиск: номер заказа, комплектация, сумма…"
        statusValue={statusFilter}
        onStatusChange={setStatusFilter}
        statusOptions={statusOptions}
      />

      {filteredOrders.length === 0 ? (
        <p className="text-sm text-slate-500 text-center py-8">
          Нет заказов по выбранным условиям. Измените поиск или статус.
        </p>
      ) : (
        filteredOrders.map((order) => {
          const orderId = order.order_id || order.id;
          const orderIdStr =
            typeof orderId === 'string' ? orderId : orderId?.toString() || 'N/A';
          return (
            <Card
              key={orderId || Math.random()}
              className="p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
                <div className="flex-1">
                  <div className="flex items-start gap-3 mb-3">
                    <div className="w-10 h-10 rounded-lg bg-blue-600 flex items-center justify-center flex-shrink-0">
                      <ShoppingCart className="w-5 h-5 text-white" />
                    </div>
                    <div className="flex-1">
                      <h3 className="text-lg font-bold text-slate-900 mb-1">
                        Заказ #{orderIdStr.slice(0, 8)}
                      </h3>
                      <p className="text-sm text-slate-500 mb-2">
                        {order.created_at
                          ? format(new Date(order.created_at), 'd MMMM yyyy', {
                              locale: ru,
                            })
                          : 'N/A'}
                      </p>
                      <OrderStatusBadge
                        code={order.status}
                        label={resolveOrderStatusLabel(order, labelByCode)}
                      />
                    </div>
                  </div>

                  {order.configuration_summary && (
                    <p className="text-sm text-slate-600 mb-1">
                      {order.configuration_summary}
                    </p>
                  )}
                  {order.manager_email && (
                    <p className="text-sm text-slate-500">
                      <span className="font-medium">Менеджер:</span> {order.manager_email}
                    </p>
                  )}
                </div>

                <div className="flex flex-col items-end gap-2">
                  <PriceDisplay price={order.total_price || order.final_price} size="lg" />
                </div>
              </div>
            </Card>
          );
        })
      )}
    </div>
  );
}
