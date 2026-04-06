import React, { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { orderService } from '@/services/orderService';
import { roleService } from '@/services/roleService';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { PERMISSIONS, hasPermission } from '@/lib/authz';
import { toast } from 'sonner';

export default function AdminControlCenter({ role }) {
  const qc = useQueryClient();
  const canManageDict = hasPermission(role, PERMISSIONS.ADMIN_ORDER_STATUSES);
  const canViewRoles = hasPermission(role, PERMISSIONS.ADMIN_ROLES_VIEW);

  const [form, setForm] = useState({
    code: '',
    customer_label_ru: '',
    admin_label_ru: '',
    description: '',
    sort_order: 100,
    is_active: true,
    is_terminal: false,
  });

  const { data: statuses = [] } = useQuery({
    queryKey: ['order-statuses', 'admin'],
    queryFn: () => orderService.getAdminOrderStatuses(),
    enabled: canManageDict || hasPermission(role, PERMISSIONS.ORDERS_MANAGE_STATUS),
  });

  const { data: roles = [] } = useQuery({
    queryKey: ['roles', 'admin'],
    queryFn: () => roleService.listRoles(),
    enabled: canViewRoles,
  });

  const createMutation = useMutation({
    mutationFn: () => orderService.createOrderStatus(form),
    onSuccess: () => {
      toast.success('Статус создан');
      setForm({
        code: '',
        customer_label_ru: '',
        admin_label_ru: '',
        description: '',
        sort_order: 100,
        is_active: true,
        is_terminal: false,
      });
      qc.invalidateQueries({ queryKey: ['order-statuses', 'admin'] });
      qc.invalidateQueries({ queryKey: ['order-statuses', 'public'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось создать статус')),
  });

  const patchMutation = useMutation({
    mutationFn: ({ id, payload }) => orderService.patchOrderStatus(id, payload),
    onSuccess: () => {
      toast.success('Статус обновлён');
      qc.invalidateQueries({ queryKey: ['order-statuses', 'admin'] });
      qc.invalidateQueries({ queryKey: ['order-statuses', 'public'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось обновить статус')),
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => orderService.deleteOrderStatus(id),
    onSuccess: () => {
      toast.success('Статус удалён');
      qc.invalidateQueries({ queryKey: ['order-statuses', 'admin'] });
      qc.invalidateQueries({ queryKey: ['order-statuses', 'public'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось удалить статус')),
  });
  const canCreateStatus =
    form.code.trim() && form.customer_label_ru.trim() && Number.isFinite(form.sort_order);

  return (
    <div className="space-y-6">
      <Card className="p-6">
        <h3 className="text-lg font-semibold text-slate-900">Права текущей роли</h3>
        <p className="text-sm text-slate-600 mt-1">Роль: {role || 'unknown'}</p>
        <div className="mt-3 flex flex-wrap gap-2 text-xs">
          {Object.values(PERMISSIONS).map((p) => (
            <span
              key={p}
              className={`rounded-full px-2.5 py-1 ${
                hasPermission(role, p) ? 'bg-green-100 text-green-800' : 'bg-slate-100 text-slate-500'
              }`}
            >
              {p}
            </span>
          ))}
        </div>
      </Card>

      {(canManageDict || hasPermission(role, PERMISSIONS.ORDERS_MANAGE_STATUS)) && (
        <Card className="p-6 space-y-4">
          <h3 className="text-lg font-semibold text-slate-900">Статусы заказов</h3>
          <div className="space-y-2">
            {statuses.length === 0 && (
              <p className="text-sm text-slate-500">Статусы не найдены.</p>
            )}
            {statuses.map((s) => (
              <div key={s.order_status_id} className="rounded-lg border p-3 flex flex-wrap gap-2 items-center justify-between">
                <div>
                  <p className="font-medium">{s.code}</p>
                  <p className="text-sm text-slate-500">{s.customer_label_ru}</p>
                </div>
                <div className="flex gap-2">
                  {canManageDict && (
                    <>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() =>
                          patchMutation.mutate({
                            id: s.order_status_id,
                            payload: { is_active: !s.is_active },
                          })
                        }
                      >
                        {s.is_active ? 'Деактивировать' : 'Активировать'}
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="text-red-600"
                        onClick={() => {
                          if (window.confirm(`Удалить статус "${s.code}"?`)) {
                            deleteMutation.mutate(s.order_status_id);
                          }
                        }}
                      >
                        Удалить
                      </Button>
                    </>
                  )}
                </div>
              </div>
            ))}
          </div>

          {canManageDict && (
            <div className="pt-4 border-t space-y-3">
              <p className="font-medium">Новый статус</p>
              <div className="grid md:grid-cols-2 gap-3">
                <div>
                  <Label>Код</Label>
                  <Input value={form.code} onChange={(e) => setForm((v) => ({ ...v, code: e.target.value }))} />
                </div>
                <div>
                  <Label>Лейбл для клиента</Label>
                  <Input
                    value={form.customer_label_ru}
                    onChange={(e) => setForm((v) => ({ ...v, customer_label_ru: e.target.value }))}
                  />
                </div>
                <div>
                  <Label>Лейбл для админа</Label>
                  <Input
                    value={form.admin_label_ru}
                    onChange={(e) => setForm((v) => ({ ...v, admin_label_ru: e.target.value }))}
                  />
                </div>
                <div>
                  <Label>Sort order</Label>
                  <Input
                    type="number"
                    value={form.sort_order}
                    onChange={(e) => setForm((v) => ({ ...v, sort_order: Number(e.target.value || 0) }))}
                  />
                </div>
              </div>
              <div className="grid md:grid-cols-2 gap-3">
                <div className="flex items-center justify-between rounded border px-3 py-2">
                  <Label>Активный</Label>
                  <Switch checked={form.is_active} onCheckedChange={(v) => setForm((f) => ({ ...f, is_active: v }))} />
                </div>
                <div className="flex items-center justify-between rounded border px-3 py-2">
                  <Label>Терминальный</Label>
                  <Switch checked={form.is_terminal} onCheckedChange={(v) => setForm((f) => ({ ...f, is_terminal: v }))} />
                </div>
              </div>
              <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending || !canCreateStatus}>
                Создать статус
              </Button>
            </div>
          )}
        </Card>
      )}

      {canViewRoles && (
        <Card className="p-6">
          <h3 className="text-lg font-semibold text-slate-900 mb-3">Справочник ролей</h3>
          <div className="space-y-2">
            {roles.length === 0 && <p className="text-sm text-slate-500">Список ролей пуст.</p>}
            {roles.map((r) => (
              <div key={r.role_id} className="rounded border px-3 py-2">
                <p className="font-medium">{r.code}</p>
                <p className="text-sm text-slate-600">{r.name_ru}</p>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
}

