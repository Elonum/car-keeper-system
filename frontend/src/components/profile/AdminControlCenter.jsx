import React, { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { orderService } from '@/services/orderService';
import { roleService } from '@/services/roleService';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { PERMISSIONS, PERMISSION_LABEL_RU, ROLE_TITLE_RU, hasPermission } from '@/lib/authz';
import CatalogManagement from './CatalogManagement';
import { ErrorNotice, FieldErrorText } from '../common/ErrorNotice';

const STATUS_DEFAULT_FORM = {
  code: '',
  customer_label_ru: '',
  admin_label_ru: '',
  description: '',
  sort_order: 100,
  is_active: true,
  is_terminal: false,
};

function normalizeStatusForm(form) {
  return {
    code: String(form.code || '').trim().toLowerCase(),
    customer_label_ru: String(form.customer_label_ru || '').trim(),
    admin_label_ru: String(form.admin_label_ru || '').trim(),
    description: String(form.description || '').trim(),
    sort_order: Number(form.sort_order),
    is_active: Boolean(form.is_active),
    is_terminal: Boolean(form.is_terminal),
  };
}

export default function AdminControlCenter({ role }) {
  const qc = useQueryClient();
  const canManageDict = hasPermission(role, PERMISSIONS.ADMIN_ORDER_STATUSES);
  const canViewRoles = hasPermission(role, PERMISSIONS.ADMIN_ROLES_VIEW);

  const [form, setForm] = useState(STATUS_DEFAULT_FORM);
  const [statusDialogOpen, setStatusDialogOpen] = useState(false);
  const [editingStatusId, setEditingStatusId] = useState(null);
  const [statusInitialSnapshot, setStatusInitialSnapshot] = useState(
    JSON.stringify(normalizeStatusForm(STATUS_DEFAULT_FORM))
  );
  const [formErrors, setFormErrors] = useState({});
  const [manageError, setManageError] = useState(null);

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
    mutationFn: (payload) => orderService.createOrderStatus(payload),
    onSuccess: () => {
      setForm(STATUS_DEFAULT_FORM);
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['order-statuses', 'admin'] });
      qc.invalidateQueries({ queryKey: ['order-statuses', 'public'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось создать статус')),
  });

  const patchMutation = useMutation({
    mutationFn: ({ id, payload }) => orderService.patchOrderStatus(id, payload),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['order-statuses', 'admin'] });
      qc.invalidateQueries({ queryKey: ['order-statuses', 'public'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось обновить статус')),
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => orderService.deleteOrderStatus(id),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['order-statuses', 'admin'] });
      qc.invalidateQueries({ queryKey: ['order-statuses', 'public'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось удалить статус')),
  });
  const canCreateStatus =
    form.code.trim() && form.customer_label_ru.trim() && Number.isFinite(form.sort_order);
  const statusPending = createMutation.isPending || patchMutation.isPending;
  const permissionRows = Object.values(PERMISSIONS).map((code) => ({
    code,
    on: hasPermission(role, code),
    label: PERMISSION_LABEL_RU[code] || code,
  }));
  const enabledPermissionsCount = permissionRows.filter((row) => row.on).length;
  const statusCurrentSnapshot = useMemo(
    () => JSON.stringify(normalizeStatusForm(form)),
    [form]
  );
  const isStatusDirty = statusCurrentSnapshot !== statusInitialSnapshot;

  const resetStatusForm = () => {
    setForm(STATUS_DEFAULT_FORM);
    setFormErrors({});
    setEditingStatusId(null);
    setStatusInitialSnapshot(JSON.stringify(normalizeStatusForm(STATUS_DEFAULT_FORM)));
  };

  const openCreateStatusDialog = () => {
    resetStatusForm();
    setManageError(null);
    setStatusDialogOpen(true);
  };

  const openEditStatusDialog = (status) => {
    setForm({
      code: status.code || '',
      customer_label_ru: status.customer_label_ru || '',
      admin_label_ru: status.admin_label_ru || '',
      description: status.description || '',
      sort_order: Number(status.sort_order ?? 100),
      is_active: Boolean(status.is_active),
      is_terminal: Boolean(status.is_terminal),
    });
    setFormErrors({});
    setManageError(null);
    setEditingStatusId(status.order_status_id);
    setStatusInitialSnapshot(
      JSON.stringify(
        normalizeStatusForm({
          code: status.code || '',
          customer_label_ru: status.customer_label_ru || '',
          admin_label_ru: status.admin_label_ru || '',
          description: status.description || '',
          sort_order: Number(status.sort_order ?? 100),
          is_active: Boolean(status.is_active),
          is_terminal: Boolean(status.is_terminal),
        })
      )
    );
    setStatusDialogOpen(true);
  };

  const validateStatusForm = () => {
    const next = {};
    const code = form.code.trim();
    const customerLabel = form.customer_label_ru.trim();
    const adminLabel = form.admin_label_ru.trim();

    if (!code) next.code = 'Укажите код статуса.';
    if (code && !/^[a-z0-9_]+$/i.test(code)) {
      next.code = 'Код: латиница, цифры и _.';
    }
    if (!customerLabel) next.customer_label_ru = 'Укажите подпись для клиента.';
    if (adminLabel.length > 100) next.admin_label_ru = 'Подпись для админа: не более 100 символов.';
    if (!Number.isFinite(form.sort_order)) next.sort_order = 'Укажите корректный порядок сортировки.';
    if (String(form.description || '').length > 500) {
      next.description = 'Описание: не более 500 символов.';
    }

    setFormErrors(next);
    return Object.keys(next).length === 0;
  };

  const handleStatusDialogSubmit = () => {
    if (!validateStatusForm()) return;
    const normalized = normalizeStatusForm(form);
    const payload = {
      code: normalized.code,
      customer_label_ru: normalized.customer_label_ru,
      admin_label_ru: normalized.admin_label_ru || null,
      description: normalized.description || null,
      sort_order: normalized.sort_order,
      is_active: normalized.is_active,
      is_terminal: normalized.is_terminal,
    };

    if (editingStatusId) {
      patchMutation.mutate(
        { id: editingStatusId, payload },
        {
          onSuccess: () => {
            setStatusDialogOpen(false);
            resetStatusForm();
          },
        }
      );
      return;
    }

    createMutation.mutate(payload, {
      onSuccess: () => {
        setStatusDialogOpen(false);
        resetStatusForm();
      },
    });
  };

  const requestCloseStatusDialog = () => {
    if (statusPending) return;
    if (isStatusDirty && !window.confirm('Есть несохранённые изменения. Закрыть без сохранения?')) {
      return;
    }
    setStatusDialogOpen(false);
    resetStatusForm();
  };

  return (
    <div className="space-y-6">
      <ErrorNotice kind="server" message={manageError} />
      <Card className="rounded-2xl border-slate-200 p-6 shadow-sm">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h3 className="text-lg font-semibold text-slate-900">Ваша роль и возможности</h3>
            <p className="text-sm text-slate-600 mt-1">
              <span className="font-medium text-slate-800">{ROLE_TITLE_RU[role] || role || '—'}</span>
              {' · '}
              доступно {enabledPermissionsCount} из {permissionRows.length} возможностей.
            </p>
          </div>
          <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-700">
            Роль: {ROLE_TITLE_RU[role] || role || '—'}
          </span>
        </div>
        <div className="mt-4 grid gap-2 sm:grid-cols-2">
          {permissionRows.map((row) => (
            <div
              key={row.code}
              className={`rounded-lg border px-3 py-2 text-sm ${
                row.on ? 'border-emerald-200 bg-emerald-50/70 text-emerald-900' : 'border-slate-200 bg-slate-50 text-slate-500'
              }`}
            >
              <span className="font-medium">{row.on ? 'Доступно' : 'Недоступно'}:</span> {row.label}
            </div>
          ))}
        </div>
      </Card>

      <CatalogManagement role={role} />

      {(canManageDict || hasPermission(role, PERMISSIONS.ORDERS_MANAGE_STATUS)) && (
        <Card className="rounded-2xl border-slate-200 p-6 shadow-sm space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <h3 className="text-lg font-semibold text-slate-900">Статусы заказов</h3>
            {canManageDict && (
              <Button type="button" onClick={openCreateStatusDialog}>
                Добавить статус
              </Button>
            )}
          </div>
          <div className="space-y-2">
            {statuses.length === 0 && (
              <p className="text-sm text-slate-500 rounded-lg border border-dashed border-slate-300 p-3">
                Статусы не найдены.
              </p>
            )}
            {statuses.map((s) => (
              <div
                key={s.order_status_id}
                className="rounded-xl border border-slate-200 bg-white p-3 flex flex-wrap gap-2 items-center justify-between"
              >
                <div className="min-w-0">
                  <p className="font-medium text-slate-900">{s.code}</p>
                  <p className="text-sm text-slate-600 truncate">{s.customer_label_ru}</p>
                </div>
                <div className="flex gap-2">
                  {canManageDict && (
                    <>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => openEditStatusDialog(s)}
                      >
                        Редактировать
                      </Button>
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

        </Card>
      )}

      {canViewRoles && (
        <Card className="rounded-2xl border-slate-200 p-6 shadow-sm">
          <h3 className="text-lg font-semibold text-slate-900 mb-3">Справочник ролей</h3>
          <div className="space-y-2">
            {roles.length === 0 && <p className="text-sm text-slate-500">Список ролей пуст.</p>}
            {roles.map((row) => (
              <div key={row.role_id} className="rounded border px-3 py-2">
                <p className="font-medium">{row.code}</p>
                <p className="text-sm text-slate-600">{row.name_ru}</p>
              </div>
            ))}
          </div>
        </Card>
      )}

      {canManageDict && (
        <Dialog
          open={statusDialogOpen}
          onOpenChange={(open) => {
            if (!open) {
              requestCloseStatusDialog();
              return;
            }
            setStatusDialogOpen(true);
          }}
        >
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle>{editingStatusId ? 'Редактирование статуса' : 'Новый статус заказа'}</DialogTitle>
              <DialogDescription>
                Заполните поля статуса. Код используется в API и должен быть стабильным.
              </DialogDescription>
            </DialogHeader>

            <form
              className="space-y-4"
              onSubmit={(e) => {
                e.preventDefault();
                handleStatusDialogSubmit();
              }}
            >
              <div className="grid md:grid-cols-2 gap-3">
                <div>
                  <Label>Код *</Label>
                  <Input
                    value={form.code}
                    onChange={(e) => {
                      setForm((v) => ({ ...v, code: e.target.value }));
                      setFormErrors((prev) => ({ ...prev, code: undefined }));
                      setManageError(null);
                    }}
                    placeholder="pending"
                  />
                  <FieldErrorText>{formErrors.code}</FieldErrorText>
                </div>
                <div>
                  <Label>Порядок сортировки *</Label>
                  <Input
                    type="number"
                    value={form.sort_order}
                    onChange={(e) => {
                      setForm((v) => ({ ...v, sort_order: Number(e.target.value || 0) }));
                      setFormErrors((prev) => ({ ...prev, sort_order: undefined }));
                      setManageError(null);
                    }}
                  />
                  <FieldErrorText>{formErrors.sort_order}</FieldErrorText>
                </div>
                <div>
                  <Label>Подпись для клиента *</Label>
                  <Input
                    value={form.customer_label_ru}
                    onChange={(e) => {
                      setForm((v) => ({ ...v, customer_label_ru: e.target.value }));
                      setFormErrors((prev) => ({ ...prev, customer_label_ru: undefined }));
                      setManageError(null);
                    }}
                  />
                  <FieldErrorText>{formErrors.customer_label_ru}</FieldErrorText>
                </div>
                <div>
                  <Label>Подпись для сотрудников</Label>
                  <Input
                    value={form.admin_label_ru}
                    onChange={(e) => {
                      setForm((v) => ({ ...v, admin_label_ru: e.target.value }));
                      setFormErrors((prev) => ({ ...prev, admin_label_ru: undefined }));
                      setManageError(null);
                    }}
                  />
                  <FieldErrorText>{formErrors.admin_label_ru}</FieldErrorText>
                </div>
              </div>
              <div>
                <Label>Описание</Label>
                <Textarea
                  value={form.description}
                  onChange={(e) => {
                    setForm((v) => ({ ...v, description: e.target.value }));
                    setFormErrors((prev) => ({ ...prev, description: undefined }));
                    setManageError(null);
                  }}
                  maxLength={500}
                />
                <FieldErrorText>{formErrors.description}</FieldErrorText>
              </div>
              <div className="grid md:grid-cols-2 gap-3">
                <div className="flex items-center justify-between rounded border px-3 py-2">
                  <Label>Активный</Label>
                  <Switch
                    checked={form.is_active}
                    onCheckedChange={(v) => {
                      setForm((f) => ({ ...f, is_active: v }));
                      setManageError(null);
                    }}
                  />
                </div>
                <div className="flex items-center justify-between rounded border px-3 py-2">
                  <Label>Терминальный</Label>
                  <Switch
                    checked={form.is_terminal}
                    onCheckedChange={(v) => {
                      setForm((f) => ({ ...f, is_terminal: v }));
                      setManageError(null);
                    }}
                  />
                </div>
              </div>
              {!canCreateStatus && <FieldErrorText>Заполните обязательные поля.</FieldErrorText>}
              {isStatusDirty && (
                <p className="text-xs text-amber-700 bg-amber-50 border border-amber-100 rounded px-2 py-1">
                  Есть несохранённые изменения.
                </p>
              )}

              <DialogFooter>
                <Button type="button" variant="outline" onClick={requestCloseStatusDialog}>
                  Отмена
                </Button>
                <Button
                  type="submit"
                  disabled={
                    statusPending ||
                    !canCreateStatus ||
                    (Boolean(editingStatusId) && !isStatusDirty)
                  }
                >
                  {editingStatusId ? 'Сохранить' : 'Создать'}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}

