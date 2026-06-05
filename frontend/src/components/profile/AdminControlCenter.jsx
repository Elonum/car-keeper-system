import React, { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
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
import { LayoutGrid, ListChecks, Shield, Tags } from 'lucide-react';
import ConfirmDialog from '@/components/common/ConfirmDialog';
import {
  ADMIN_LIMITS,
  validateOrderStatusForm,
} from '@/lib/adminValidation';

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
  const canCatalog = hasPermission(role, PERMISSIONS.CATALOG_MANAGE);
  const canService = hasPermission(role, PERMISSIONS.SERVICE_MANAGE);
  const canStatuses =
    canManageDict || hasPermission(role, PERMISSIONS.ORDERS_MANAGE_STATUS);

  const defaultTab = canCatalog || canService ? 'catalog' : canStatuses ? 'statuses' : 'access';

  const [form, setForm] = useState(STATUS_DEFAULT_FORM);
  const [statusDialogOpen, setStatusDialogOpen] = useState(false);
  const [editingStatusId, setEditingStatusId] = useState(null);
  const [statusInitialSnapshot, setStatusInitialSnapshot] = useState(
    JSON.stringify(normalizeStatusForm(STATUS_DEFAULT_FORM))
  );
  const [formErrors, setFormErrors] = useState({});
  const [manageError, setManageError] = useState(null);
  const [confirmIntent, setConfirmIntent] = useState(null);

  const { data: statuses = [] } = useQuery({
    queryKey: ['order-statuses', 'admin'],
    queryFn: () => orderService.getAdminOrderStatuses(),
    enabled: canStatuses,
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
    const next = validateOrderStatusForm(form);
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
    if (isStatusDirty) {
      setConfirmIntent({
        variant: 'warning',
        title: 'Закрыть без сохранения?',
        description: 'Внесённые изменения в статус заказа будут потеряны.',
        confirmLabel: 'Закрыть',
        onConfirm: () => {
          setStatusDialogOpen(false);
          resetStatusForm();
        },
      });
      return;
    }
    setStatusDialogOpen(false);
    resetStatusForm();
  };

  const requestDeleteStatus = (status) => {
    setConfirmIntent({
      variant: 'destructive',
      title: 'Удалить статус заказа?',
      description: `Статус «${status.code}» будет удалён из справочника. Если он уже используется в заказах, сервер безопасно отклонит удаление.`,
      confirmLabel: 'Удалить',
      onConfirm: () => deleteMutation.mutate(status.order_status_id),
    });
  };

  return (
    <div className="space-y-4">
      <ErrorNotice kind="server" message={manageError} />

      <Tabs defaultValue={defaultTab} className="space-y-4">
        <TabsList className="bg-white p-1.5 h-auto rounded-xl shadow-sm border border-slate-200 flex-wrap w-full justify-start">
          <TabsTrigger
            value="access"
            className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-3 py-2 text-sm"
          >
            <Shield className="w-4 h-4" />
            Доступ
          </TabsTrigger>
          {(canCatalog || canService) && (
            <TabsTrigger
              value="catalog"
              className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-3 py-2 text-sm"
            >
              <LayoutGrid className="w-4 h-4" />
              Каталог и сервис
            </TabsTrigger>
          )}
          {canStatuses && (
            <TabsTrigger
              value="statuses"
              className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-3 py-2 text-sm"
            >
              <ListChecks className="w-4 h-4" />
              Статусы заказов
            </TabsTrigger>
          )}
          {canViewRoles && (
            <TabsTrigger
              value="roles"
              className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-3 py-2 text-sm"
            >
              <Tags className="w-4 h-4" />
              Роли
            </TabsTrigger>
          )}
        </TabsList>

        <TabsContent value="access" className="mt-0">
          <Card className="rounded-2xl border-slate-200 p-5 shadow-sm">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <div>
                <h3 className="text-base font-semibold text-slate-900">Ваша роль</h3>
                <p className="text-sm text-slate-600 mt-0.5">
                  {ROLE_TITLE_RU[role] || role || '—'} · {enabledPermissionsCount} из{' '}
                  {permissionRows.length} прав
                </p>
              </div>
              <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-700">
                {ROLE_TITLE_RU[role] || role}
              </span>
            </div>
            <ul className="grid gap-1.5 sm:grid-cols-2">
              {permissionRows.map((row) => (
                <li
                  key={row.code}
                  className={`rounded-lg border px-3 py-2 text-sm ${
                    row.on
                      ? 'border-emerald-200/80 bg-emerald-50/60 text-emerald-950'
                      : 'border-slate-200 bg-slate-50 text-slate-500'
                  }`}
                >
                  {row.label}
                </li>
              ))}
            </ul>
          </Card>
        </TabsContent>

        {(canCatalog || canService) && (
          <TabsContent value="catalog" className="mt-0">
            <CatalogManagement role={role} />
          </TabsContent>
        )}

        {canStatuses && (
          <TabsContent value="statuses" className="mt-0">
            <Card className="rounded-2xl border-slate-200 p-5 shadow-sm space-y-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <h3 className="text-base font-semibold text-slate-900">Статусы заказов</h3>
                {canManageDict && (
                  <Button type="button" size="sm" onClick={openCreateStatusDialog}>
                    Добавить статус
                  </Button>
                )}
              </div>
              <div className="space-y-2">
                {statuses.length === 0 && (
                  <p className="text-sm text-slate-500 rounded-lg border border-dashed border-slate-300 p-4 text-center">
                    Статусы не найдены
                  </p>
                )}
                {statuses.map((s) => (
                  <div
                    key={s.order_status_id}
                    className="rounded-xl border border-slate-200 bg-white p-3 flex flex-wrap gap-2 items-center justify-between"
                  >
                    <div className="min-w-0">
                      <p className="font-medium text-slate-900 font-mono text-sm">{s.code}</p>
                      <p className="text-sm text-slate-600 truncate">{s.customer_label_ru}</p>
                    </div>
                    {canManageDict && (
                      <div className="flex flex-wrap gap-2">
                        <Button size="sm" variant="outline" onClick={() => openEditStatusDialog(s)}>
                          Изменить
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
                          {s.is_active ? 'Выкл.' : 'Вкл.'}
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          className="text-red-600 hover:text-red-700"
                          onClick={() => requestDeleteStatus(s)}
                        >
                          Удалить
                        </Button>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </Card>
          </TabsContent>
        )}

        {canViewRoles && (
          <TabsContent value="roles" className="mt-0">
            <Card className="rounded-2xl border-slate-200 p-5 shadow-sm">
              <h3 className="text-base font-semibold text-slate-900 mb-3">Справочник ролей</h3>
              <div className="grid gap-2 sm:grid-cols-2">
                {roles.length === 0 && (
                  <p className="text-sm text-slate-500 col-span-full">Список ролей пуст</p>
                )}
                {roles.map((row) => (
                  <div
                    key={row.role_id}
                    className="rounded-xl border border-slate-200 bg-white px-3 py-2.5"
                  >
                    <p className="font-medium text-slate-900 font-mono text-sm">{row.code}</p>
                    <p className="text-sm text-slate-600">{row.name_ru}</p>
                  </div>
                ))}
              </div>
            </Card>
          </TabsContent>
        )}
      </Tabs>

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
                Код используется в API и должен оставаться стабильным после публикации.
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
                    maxLength={ADMIN_LIMITS.ORDER_STATUS_CODE}
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
                  <Label>Порядок *</Label>
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
                    maxLength={ADMIN_LIMITS.ORDER_STATUS_LABEL}
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
                    maxLength={ADMIN_LIMITS.ORDER_STATUS_LABEL}
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
                  maxLength={ADMIN_LIMITS.ORDER_STATUS_DESCRIPTION}
                />
                <FieldErrorText>{formErrors.description}</FieldErrorText>
              </div>
              <div className="grid md:grid-cols-2 gap-3">
                <div className="flex items-center justify-between rounded-lg border px-3 py-2">
                  <Label>Активный</Label>
                  <Switch
                    checked={form.is_active}
                    onCheckedChange={(v) => {
                      setForm((f) => ({ ...f, is_active: v }));
                      setManageError(null);
                    }}
                  />
                </div>
                <div className="flex items-center justify-between rounded-lg border px-3 py-2">
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
              {!canCreateStatus && <FieldErrorText>Заполните обязательные поля</FieldErrorText>}
              {isStatusDirty && (
                <p className="text-xs text-amber-700 bg-amber-50 border border-amber-100 rounded px-2 py-1">
                  Есть несохранённые изменения
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
      <ConfirmDialog
        open={Boolean(confirmIntent)}
        onOpenChange={(open) => {
          if (!open && !deleteMutation.isPending) setConfirmIntent(null);
        }}
        variant={confirmIntent?.variant}
        title={confirmIntent?.title}
        description={confirmIntent?.description}
        confirmLabel={deleteMutation.isPending ? 'Удаление…' : confirmIntent?.confirmLabel}
        disabled={deleteMutation.isPending}
        onConfirm={() => {
          const action = confirmIntent?.onConfirm;
          action?.();
          if (!deleteMutation.isPending) setConfirmIntent(null);
        }}
      />
    </div>
  );
}
