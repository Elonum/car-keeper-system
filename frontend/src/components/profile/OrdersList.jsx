import React, { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/lib/AuthContext';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import OrderStatusBadge from '../common/OrderStatusBadge';
import PriceDisplay from '../common/PriceDisplay';
import { orderService } from '@/services/orderService';
import { useOrderStatusLabelMap } from '@/hooks/useOrderStatusLabelMap';
import { resolveOrderStatusLabel } from '@/lib/orderStatusDisplay';
import { filterOrdersForCabinet } from '@/lib/cabinetFilters';
import { PERMISSIONS, hasPermission } from '@/lib/authz';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { toast } from 'sonner';
import CabinetListToolbar from './CabinetListToolbar';
import EmptyState from '../common/EmptyState';
import { ShoppingCart } from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

function customerAllowedStatuses(currentStatus) {
  if (currentStatus === 'pending' || currentStatus === 'approved') {
    return ['cancelled'];
  }
  return [];
}

export default function OrdersList({ orders, isLoading }) {
  const qc = useQueryClient();
  const { user } = useAuth();
  const canManageStatuses = hasPermission(user?.role, PERMISSIONS.ORDERS_MANAGE_STATUS);
  const role = user?.role || '';
  const { labelByCode, data: orderStatusRows } = useOrderStatusLabelMap();
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [nextStatusByOrder, setNextStatusByOrder] = useState({});

  const staffStatusesQuery = useQuery({
    queryKey: ['order-statuses', 'admin'],
    queryFn: () => orderService.getAdminOrderStatuses(),
    enabled: canManageStatuses,
  });
  const statusRowsForActions =
    canManageStatuses && Array.isArray(staffStatusesQuery.data)
      ? staffStatusesQuery.data
      : Array.isArray(orderStatusRows)
        ? orderStatusRows
        : [];

  const updateStatusMutation = useMutation({
    mutationFn: ({ orderId, status }) => orderService.updateOrderStatus(orderId, status),
    onSuccess: () => {
      toast.success('Статус заказа обновлён');
      qc.invalidateQueries({ queryKey: ['my-orders'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось обновить статус заказа')),
  });

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
          const customerActions = customerAllowedStatuses(order.status);
          const selectedStatus = nextStatusByOrder[orderIdStr] || order.status;
          const actionStatuses = canManageStatuses
            ? statusRowsForActions
            : statusRowsForActions.filter((s) => customerActions.includes(s.code));
          const canSubmit =
            canManageStatuses ||
            (role === 'customer' &&
              customerActions.includes(selectedStatus) &&
              selectedStatus !== order.status);
          return (
            <Card
              key={orderIdStr}
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
                  {(canManageStatuses || role === 'customer') && (
                    <div className="w-56 flex flex-col gap-2">
                      <Select
                        value={selectedStatus}
                        onValueChange={(v) =>
                          setNextStatusByOrder((prev) => ({ ...prev, [orderIdStr]: v }))
                        }
                        disabled={!canManageStatuses && customerActions.length === 0}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Выберите статус" />
                        </SelectTrigger>
                        <SelectContent>
                          {actionStatuses.map((s) => (
                            <SelectItem key={s.code} value={s.code}>
                              {s.customer_label_ru || s.admin_label_ru || s.code}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Button
                        size="sm"
                        onClick={() =>
                          updateStatusMutation.mutate({
                            orderId: order.order_id || order.id,
                            status: selectedStatus,
                          })
                        }
                        disabled={updateStatusMutation.isPending || !canSubmit}
                      >
                        {role === 'customer' ? 'Отменить заказ' : 'Сменить статус'}
                      </Button>
                      {!canManageStatuses && customerActions.length === 0 && (
                        <p className="text-xs text-slate-500">
                          Для этого заказа смена статуса недоступна.
                        </p>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </Card>
          );
        })
      )}
    </div>
  );
}
