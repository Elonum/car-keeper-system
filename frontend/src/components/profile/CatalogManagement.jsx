import React, { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { catalogService } from '@/services/catalogService';
import { serviceService } from '@/services/serviceService';
import { adminCatalogService } from '@/services/adminCatalogService';
import { PERMISSIONS, hasPermission } from '@/lib/authz';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { ErrorNotice, FieldErrorText } from '@/components/common/ErrorNotice';

const SERVICE_CATEGORIES = ['maintenance', 'repair', 'diagnostics', 'detailing', 'tires'];
const BRAND_DEFAULT_FORM = { name: '', country: '' };
const MODEL_DEFAULT_FORM = { brand_id: '', name: '', segment: '', description: '' };
const SERVICE_DEFAULT_FORM = {
  name: '',
  category: 'maintenance',
  description: '',
  price: '',
  duration_minutes: '',
  is_available: true,
};

function normalizeBrandForm(form) {
  return {
    name: String(form.name || '').trim(),
    country: String(form.country || '').trim(),
  };
}

function normalizeServiceForm(form) {
  return {
    name: String(form.name || '').trim(),
    category: String(form.category || 'maintenance'),
    description: String(form.description || '').trim(),
    price: Number(form.price),
    duration_minutes:
      form.duration_minutes === '' || form.duration_minutes == null
        ? null
        : Number(form.duration_minutes),
    is_available: Boolean(form.is_available),
  };
}

function normalizeModelForm(form) {
  return {
    brand_id: String(form.brand_id || '').trim(),
    name: String(form.name || '').trim(),
    segment: String(form.segment || '').trim(),
    description: String(form.description || '').trim(),
  };
}

export default function CatalogManagement({ role }) {
  const qc = useQueryClient();
  const canCatalog = hasPermission(role, PERMISSIONS.CATALOG_MANAGE);
  const canService = hasPermission(role, PERMISSIONS.SERVICE_MANAGE);

  const { data: brands = [] } = useQuery({
    queryKey: ['catalog', 'brands'],
    queryFn: () => catalogService.getBrands(),
    enabled: canCatalog,
  });

  const { data: serviceTypes = [] } = useQuery({
    queryKey: ['service', 'types', 'admin'],
    queryFn: () => serviceService.getServiceTypes(),
    enabled: canService,
  });

  const { data: branches = [] } = useQuery({
    queryKey: ['service', 'branches'],
    queryFn: () => serviceService.getBranches(),
    enabled: canService,
  });
  const { data: models = [] } = useQuery({
    queryKey: ['catalog', 'models', 'admin'],
    queryFn: () => adminCatalogService.listModels(),
    enabled: canCatalog,
  });

  const [brandForm, setBrandForm] = useState(BRAND_DEFAULT_FORM);
  const [modelForm, setModelForm] = useState(MODEL_DEFAULT_FORM);
  const [modelImageFile, setModelImageFile] = useState(null);
  const [stForm, setStForm] = useState(SERVICE_DEFAULT_FORM);
  const [brandDialogOpen, setBrandDialogOpen] = useState(false);
  const [modelDialogOpen, setModelDialogOpen] = useState(false);
  const [serviceDialogOpen, setServiceDialogOpen] = useState(false);
  const [editingBrandId, setEditingBrandId] = useState(null);
  const [editingModelId, setEditingModelId] = useState(null);
  const [editingServiceTypeId, setEditingServiceTypeId] = useState(null);
  const [brandInitialSnapshot, setBrandInitialSnapshot] = useState(
    JSON.stringify(normalizeBrandForm(BRAND_DEFAULT_FORM))
  );
  const [serviceInitialSnapshot, setServiceInitialSnapshot] = useState(
    JSON.stringify(normalizeServiceForm(SERVICE_DEFAULT_FORM))
  );
  const [modelInitialSnapshot, setModelInitialSnapshot] = useState(
    JSON.stringify(normalizeModelForm(MODEL_DEFAULT_FORM))
  );
  const [brandErrors, setBrandErrors] = useState({});
  const [modelErrors, setModelErrors] = useState({});
  const [serviceErrors, setServiceErrors] = useState({});
  const [manageError, setManageError] = useState(null);
  const updateBrand = useMutation({
    mutationFn: ({ id, payload }) => adminCatalogService.updateBrand(id, payload),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['catalog', 'brands'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось обновить бренд')),
  });


  const createBrand = useMutation({
    mutationFn: (payload) => adminCatalogService.createBrand(payload),
    onSuccess: () => {
      setBrandForm(BRAND_DEFAULT_FORM);
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['catalog', 'brands'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось создать бренд')),
  });

  const deleteBrand = useMutation({
    mutationFn: (id) => adminCatalogService.deleteBrand(id),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['catalog', 'brands'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось удалить')),
  });
  const createModel = useMutation({
    mutationFn: (payload) => adminCatalogService.createModel(payload),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['catalog', 'models', 'admin'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось создать автомобиль')),
  });
  const updateModel = useMutation({
    mutationFn: ({ id, payload }) => adminCatalogService.updateModel(id, payload),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['catalog', 'models', 'admin'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось обновить автомобиль')),
  });
  const deleteModel = useMutation({
    mutationFn: (id) => adminCatalogService.deleteModel(id),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['catalog', 'models', 'admin'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось удалить автомобиль')),
  });
  const uploadModelImage = useMutation({
    mutationFn: ({ id, file }) => adminCatalogService.uploadModelImage(id, file),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['catalog', 'models', 'admin'] });
      qc.invalidateQueries({ queryKey: ['trims'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось загрузить изображение')),
  });

  const createSt = useMutation({
    mutationFn: (payload) => adminCatalogService.createServiceType(payload),
    onSuccess: () => {
      setStForm(SERVICE_DEFAULT_FORM);
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['service', 'types', 'admin'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось создать услугу')),
  });

  const patchSt = useMutation({
    mutationFn: ({ id, payload }) => adminCatalogService.updateServiceType(id, payload),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['service', 'types', 'admin'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось обновить')),
  });

  const patchBranch = useMutation({
    mutationFn: ({ id, payload }) => adminCatalogService.updateBranch(id, payload),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['service', 'branches'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось обновить филиал')),
  });

  const deleteSt = useMutation({
    mutationFn: (id) => adminCatalogService.deleteServiceType(id),
    onSuccess: () => {
      setManageError(null);
      qc.invalidateQueries({ queryKey: ['service', 'types', 'admin'] });
    },
    onError: (e) => setManageError(getApiErrorMessage(e, 'Не удалось удалить')),
  });

  if (!canCatalog && !canService) return null;
  const brandPending = createBrand.isPending || updateBrand.isPending;
  const modelPending =
    createModel.isPending || updateModel.isPending || deleteModel.isPending || uploadModelImage.isPending;
  const servicePending = createSt.isPending || patchSt.isPending;
  const brandCurrentSnapshot = useMemo(
    () => JSON.stringify(normalizeBrandForm(brandForm)),
    [brandForm]
  );
  const serviceCurrentSnapshot = useMemo(
    () => JSON.stringify(normalizeServiceForm(stForm)),
    [stForm]
  );
  const modelCurrentSnapshot = useMemo(
    () => JSON.stringify(normalizeModelForm(modelForm)),
    [modelForm]
  );
  const isBrandDirty = brandCurrentSnapshot !== brandInitialSnapshot;
  const isModelDirty = modelCurrentSnapshot !== modelInitialSnapshot || Boolean(modelImageFile);
  const isServiceDirty = serviceCurrentSnapshot !== serviceInitialSnapshot;

  const resetBrandForm = () => {
    setBrandForm(BRAND_DEFAULT_FORM);
    setBrandErrors({});
    setEditingBrandId(null);
    setBrandInitialSnapshot(JSON.stringify(normalizeBrandForm(BRAND_DEFAULT_FORM)));
  };

  const resetServiceForm = () => {
    setStForm(SERVICE_DEFAULT_FORM);
    setServiceErrors({});
    setEditingServiceTypeId(null);
    setServiceInitialSnapshot(JSON.stringify(normalizeServiceForm(SERVICE_DEFAULT_FORM)));
  };
  const resetModelForm = () => {
    setModelForm(MODEL_DEFAULT_FORM);
    setModelImageFile(null);
    setModelErrors({});
    setEditingModelId(null);
    setModelInitialSnapshot(JSON.stringify(normalizeModelForm(MODEL_DEFAULT_FORM)));
  };

  const openCreateBrandDialog = () => {
    resetBrandForm();
    setManageError(null);
    setBrandDialogOpen(true);
  };

  const openEditBrandDialog = (brand) => {
    setBrandForm({
      name: String(brand.name || ''),
      country: String(brand.country || ''),
    });
    setBrandErrors({});
    setEditingBrandId(brand.brand_id);
    setBrandInitialSnapshot(
      JSON.stringify(
        normalizeBrandForm({
          name: String(brand.name || ''),
          country: String(brand.country || ''),
        })
      )
    );
    setManageError(null);
    setBrandDialogOpen(true);
  };

  const openCreateServiceDialog = () => {
    resetServiceForm();
    setManageError(null);
    setServiceDialogOpen(true);
  };
  const openCreateModelDialog = () => {
    resetModelForm();
    setManageError(null);
    setModelDialogOpen(true);
  };
  const openEditModelDialog = (item) => {
    const next = {
      brand_id: String(item.brand_id || ''),
      name: String(item.name || ''),
      segment: String(item.segment || ''),
      description: String(item.description || ''),
    };
    setModelForm(next);
    setModelImageFile(null);
    setModelErrors({});
    setEditingModelId(item.model_id);
    setModelInitialSnapshot(JSON.stringify(normalizeModelForm(next)));
    setManageError(null);
    setModelDialogOpen(true);
  };

  const openEditServiceDialog = (st) => {
    setStForm({
      name: String(st.name || ''),
      category: st.category || 'maintenance',
      description: String(st.description || ''),
      price: String(st.price ?? ''),
      duration_minutes: st.duration_minutes != null ? String(st.duration_minutes) : '',
      is_available: Boolean(st.is_available),
    });
    setServiceErrors({});
    setEditingServiceTypeId(st.service_type_id);
    setServiceInitialSnapshot(
      JSON.stringify(
        normalizeServiceForm({
          name: String(st.name || ''),
          category: st.category || 'maintenance',
          description: String(st.description || ''),
          price: String(st.price ?? ''),
          duration_minutes: st.duration_minutes != null ? String(st.duration_minutes) : '',
          is_available: Boolean(st.is_available),
        })
      )
    );
    setManageError(null);
    setServiceDialogOpen(true);
  };

  const validateBrandForm = () => {
    const next = {};
    if (!brandForm.name.trim()) next.name = 'Укажите название бренда.';
    if (brandForm.name.trim().length > 100) next.name = 'Название: не более 100 символов.';
    if (!brandForm.country.trim()) next.country = 'Укажите страну.';
    if (brandForm.country.trim().length > 100) next.country = 'Страна: не более 100 символов.';
    setBrandErrors(next);
    return Object.keys(next).length === 0;
  };

  const validateServiceForm = () => {
    const next = {};
    const price = Number(stForm.price);
    const duration = stForm.duration_minutes === '' ? null : Number(stForm.duration_minutes);
    if (!stForm.name.trim()) next.name = 'Укажите название услуги.';
    if (stForm.name.trim().length > 120) next.name = 'Название: не более 120 символов.';
    if (!SERVICE_CATEGORIES.includes(stForm.category)) next.category = 'Выберите категорию из списка.';
    if (!Number.isFinite(price) || price < 0) next.price = 'Цена должна быть числом >= 0.';
    if (duration != null && (!Number.isFinite(duration) || duration <= 0 || duration > 1440)) {
      next.duration_minutes = 'Длительность: от 1 до 1440 минут.';
    }
    if (String(stForm.description || '').length > 1000) {
      next.description = 'Описание: не более 1000 символов.';
    }
    setServiceErrors(next);
    return Object.keys(next).length === 0;
  };
  const validateModelForm = () => {
    const next = {};
    const normalized = normalizeModelForm(modelForm);
    if (!normalized.brand_id) next.brand_id = 'Выберите бренд.';
    if (!normalized.name) next.name = 'Укажите название модели.';
    if (normalized.name.length > 150) next.name = 'Название: не более 150 символов.';
    if (normalized.segment.length > 100) next.segment = 'Сегмент: не более 100 символов.';
    if (normalized.description.length > 2000) next.description = 'Описание: не более 2000 символов.';
    if (modelImageFile) {
      const allowed = ['image/jpeg', 'image/png', 'image/webp'];
      if (!allowed.includes(modelImageFile.type)) next.image = 'Допустимы JPG, PNG или WEBP.';
      if (modelImageFile.size > 5 * 1024 * 1024) next.image = 'Размер файла не более 5 МБ.';
    }
    setModelErrors(next);
    return Object.keys(next).length === 0;
  };

  const submitBrandDialog = () => {
    if (!validateBrandForm()) return;
    const payload = normalizeBrandForm(brandForm);
    if (editingBrandId) {
      updateBrand.mutate(
        { id: editingBrandId, payload },
        {
          onSuccess: () => {
            setBrandDialogOpen(false);
            resetBrandForm();
          },
        }
      );
      return;
    }
    createBrand.mutate(payload, {
      onSuccess: () => {
        setBrandDialogOpen(false);
        resetBrandForm();
      },
    });
  };

  const submitServiceDialog = () => {
    if (!validateServiceForm()) return;
    const normalized = normalizeServiceForm(stForm);
    const payload = {
      name: normalized.name,
      category: normalized.category,
      description: normalized.description || null,
      price: normalized.price,
      duration_minutes: normalized.duration_minutes,
      is_available: normalized.is_available,
    };
    if (editingServiceTypeId) {
      patchSt.mutate(
        { id: editingServiceTypeId, payload },
        {
          onSuccess: () => {
            setServiceDialogOpen(false);
            resetServiceForm();
          },
        }
      );
      return;
    }
    createSt.mutate(payload, {
      onSuccess: () => {
        setServiceDialogOpen(false);
        resetServiceForm();
      },
    });
  };
  const submitModelDialog = () => {
    if (!validateModelForm()) return;
    const normalized = normalizeModelForm(modelForm);
    const payload = {
      brand_id: normalized.brand_id,
      name: normalized.name,
      segment: normalized.segment || null,
      description: normalized.description || null,
    };

    const onSuccessWithImage = (id) => {
      if (!modelImageFile) {
        setModelDialogOpen(false);
        resetModelForm();
        return;
      }
      uploadModelImage.mutate(
        { id, file: modelImageFile },
        {
          onSuccess: () => {
            setModelDialogOpen(false);
            resetModelForm();
          },
        }
      );
    };

    if (editingModelId) {
      updateModel.mutate(
        { id: editingModelId, payload },
        {
          onSuccess: () => onSuccessWithImage(editingModelId),
        }
      );
      return;
    }
    createModel.mutate(payload, {
      onSuccess: (created) => onSuccessWithImage(created.model_id),
    });
  };

  const requestCloseBrandDialog = () => {
    if (brandPending) return;
    if (isBrandDirty && !window.confirm('Есть несохранённые изменения. Закрыть без сохранения?')) {
      return;
    }
    setBrandDialogOpen(false);
    resetBrandForm();
  };

  const requestCloseServiceDialog = () => {
    if (servicePending) return;
    if (isServiceDirty && !window.confirm('Есть несохранённые изменения. Закрыть без сохранения?')) {
      return;
    }
    setServiceDialogOpen(false);
    resetServiceForm();
  };
  const requestCloseModelDialog = () => {
    if (modelPending) return;
    if (isModelDirty && !window.confirm('Есть несохранённые изменения. Закрыть без сохранения?')) {
      return;
    }
    setModelDialogOpen(false);
    resetModelForm();
  };

  return (
    <div className="space-y-6">
      <ErrorNotice kind="server" message={manageError} />
      {canCatalog && (
        <Card className="rounded-2xl border-slate-200 p-6 shadow-sm space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
            <h3 className="text-lg font-semibold text-slate-900">Каталог: бренды</h3>
            <p className="text-sm text-slate-600 mt-1">
              Добавление и удаление брендов. Модели и комплектации редактируются отдельно (при необходимости — через БД или расширение API).
            </p>
            </div>
            <Button type="button" onClick={openCreateBrandDialog}>
              Добавить бренд
            </Button>
          </div>
          <div className="flex flex-wrap gap-2 max-h-40 overflow-y-auto">
            {brands.length === 0 && (
              <p className="w-full text-sm text-slate-500 rounded-lg border border-dashed border-slate-300 p-3">
                Бренды пока не добавлены.
              </p>
            )}
            {brands.map((b) => (
              <div
                key={b.brand_id}
                className="flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-2 py-1 text-sm"
              >
                <span>
                  {b.name} ({b.country})
                </span>
                <Button type="button" variant="ghost" size="sm" onClick={() => openEditBrandDialog(b)}>
                  Изм.
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="text-red-600 h-7"
                  onClick={() => {
                    if (window.confirm(`Удалить бренд «${b.name}»?`)) deleteBrand.mutate(b.brand_id);
                  }}
                >
                  ×
                </Button>
              </div>
            ))}
          </div>
        </Card>
      )}
      {canCatalog && (
        <Card className="rounded-2xl border-slate-200 p-6 shadow-sm space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              <h3 className="text-lg font-semibold text-slate-900">Каталог: автомобили</h3>
              <p className="text-sm text-slate-600 mt-1">
                Управление моделями автомобилей: создание, редактирование, удаление и загрузка изображения.
              </p>
            </div>
            <Button type="button" onClick={openCreateModelDialog}>
              Добавить автомобиль
            </Button>
          </div>
          <div className="space-y-2 max-h-72 overflow-y-auto">
            {models.length === 0 && (
              <p className="text-sm text-slate-500 rounded-lg border border-dashed border-slate-300 p-3">
                Автомобили пока не добавлены.
              </p>
            )}
            {models.map((m) => {
              const brandLabel = brands.find((b) => b.brand_id === m.brand_id)?.name || 'Без бренда';
              return (
                <div
                  key={m.model_id}
                  className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-slate-200 bg-white p-2 text-sm"
                >
                  <div className="min-w-0">
                    <span className="font-medium">{brandLabel} {m.name}</span>
                    {m.segment ? <span className="text-slate-500 ml-2">· {m.segment}</span> : null}
                  </div>
                  <div className="flex items-center gap-2">
                    <Button type="button" variant="ghost" size="sm" onClick={() => openEditModelDialog(m)}>
                      Изм.
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="text-red-600"
                      onClick={() => {
                        if (window.confirm(`Удалить автомобиль «${m.name}»?`)) {
                          deleteModel.mutate(m.model_id);
                        }
                      }}
                    >
                      Удалить
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        </Card>
      )}

      {canService && (
        <Card className="rounded-2xl border-slate-200 p-6 shadow-sm space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <h3 className="text-lg font-semibold text-slate-900">Сервис: услуги</h3>
            <Button type="button" onClick={openCreateServiceDialog}>
              Добавить услугу
            </Button>
          </div>
          <div className="space-y-2 max-h-56 overflow-y-auto">
            {serviceTypes.length === 0 && (
              <p className="text-sm text-slate-500 rounded-lg border border-dashed border-slate-300 p-3">
                Услуги пока не добавлены.
              </p>
            )}
            {serviceTypes.map((st) => (
              <div
                key={st.service_type_id}
                className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-slate-200 bg-white p-2 text-sm"
              >
                <div>
                  <span className="font-medium">{st.name}</span>
                  <span className="text-slate-500 ml-2">
                    {st.category} · {st.price} ₽
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => openEditServiceDialog(st)}
                  >
                    Изм.
                  </Button>
                  <Switch
                    checked={st.is_available}
                    onCheckedChange={(v) =>
                      patchSt.mutate({
                        id: st.service_type_id,
                        payload: {
                          name: st.name,
                          category: st.category,
                          description: st.description,
                          price: st.price,
                          duration_minutes: st.duration_minutes,
                          is_available: v,
                        },
                      })
                    }
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="text-red-600"
                    onClick={() => {
                      if (window.confirm(`Удалить услугу «${st.name}»?`)) {
                        deleteSt.mutate(st.service_type_id);
                      }
                    }}
                  >
                    Удалить
                  </Button>
                </div>
              </div>
            ))}
          </div>

          <div className="border-t pt-4">
            <h4 className="font-medium text-slate-900 mb-2">Филиалы</h4>
            <div className="space-y-2">
              {branches.map((br) => (
                <div key={br.branch_id} className="flex flex-wrap items-center justify-between gap-2 rounded border p-2 text-sm">
                  <span>{br.name}</span>
                  <div className="flex items-center gap-2">
                    <span className="text-slate-500">Активен</span>
                    <Switch
                      checked={br.is_active}
                      onCheckedChange={(v) =>
                        patchBranch.mutate({ id: br.branch_id, payload: { is_active: v } })
                      }
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </Card>
      )}

      {canCatalog && (
        <Dialog
          open={brandDialogOpen}
          onOpenChange={(open) => {
            if (!open) {
              requestCloseBrandDialog();
              return;
            }
            setBrandDialogOpen(true);
          }}
        >
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>{editingBrandId ? 'Редактировать бренд' : 'Добавить бренд'}</DialogTitle>
              <DialogDescription>Укажите название бренда и страну происхождения.</DialogDescription>
            </DialogHeader>
            <form
              className="space-y-3"
              onSubmit={(e) => {
                e.preventDefault();
                submitBrandDialog();
              }}
            >
              <div>
                <Label>Название *</Label>
                <Input
                  value={brandForm.name}
                  onChange={(e) => {
                    setBrandForm((f) => ({ ...f, name: e.target.value }));
                    setBrandErrors((prev) => ({ ...prev, name: undefined }));
                    setManageError(null);
                  }}
                />
                <FieldErrorText>{brandErrors.name}</FieldErrorText>
              </div>
              <div>
                <Label>Страна *</Label>
                <Input
                  value={brandForm.country}
                  onChange={(e) => {
                    setBrandForm((f) => ({ ...f, country: e.target.value }));
                    setBrandErrors((prev) => ({ ...prev, country: undefined }));
                    setManageError(null);
                  }}
                />
                <FieldErrorText>{brandErrors.country}</FieldErrorText>
              </div>
              {isBrandDirty && (
                <p className="text-xs text-amber-700 bg-amber-50 border border-amber-100 rounded px-2 py-1">
                  Есть несохранённые изменения.
                </p>
              )}
              <DialogFooter>
                <Button type="button" variant="outline" onClick={requestCloseBrandDialog}>
                  Отмена
                </Button>
                <Button
                  type="submit"
                  disabled={
                    brandPending || (Boolean(editingBrandId) && !isBrandDirty)
                  }
                >
                  {editingBrandId ? 'Сохранить' : 'Создать'}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      )}

      {canService && (
        <Dialog
          open={serviceDialogOpen}
          onOpenChange={(open) => {
            if (!open) {
              requestCloseServiceDialog();
              return;
            }
            setServiceDialogOpen(true);
          }}
        >
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle>
                {editingServiceTypeId ? 'Редактировать услугу' : 'Добавить услугу'}
              </DialogTitle>
              <DialogDescription>
                Заполните параметры услуги. Проверка полей выполняется до отправки на API.
              </DialogDescription>
            </DialogHeader>
            <form
              className="space-y-3"
              onSubmit={(e) => {
                e.preventDefault();
                submitServiceDialog();
              }}
            >
              <div className="grid sm:grid-cols-2 gap-3">
                <div>
                  <Label>Название *</Label>
                  <Input
                    value={stForm.name}
                    onChange={(e) => {
                      setStForm((f) => ({ ...f, name: e.target.value }));
                      setServiceErrors((prev) => ({ ...prev, name: undefined }));
                      setManageError(null);
                    }}
                  />
                  <FieldErrorText>{serviceErrors.name}</FieldErrorText>
                </div>
                <div>
                  <Label>Категория *</Label>
                  <Select
                    value={stForm.category}
                    onValueChange={(v) => {
                      setStForm((f) => ({ ...f, category: v }));
                      setServiceErrors((prev) => ({ ...prev, category: undefined }));
                      setManageError(null);
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {SERVICE_CATEGORIES.map((c) => (
                        <SelectItem key={c} value={c}>
                          {c}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FieldErrorText>{serviceErrors.category}</FieldErrorText>
                </div>
                <div>
                  <Label>Цена (₽) *</Label>
                  <Input
                    type="number"
                    min="0"
                    step="0.01"
                    value={stForm.price}
                    onChange={(e) => {
                      setStForm((f) => ({ ...f, price: e.target.value }));
                      setServiceErrors((prev) => ({ ...prev, price: undefined }));
                      setManageError(null);
                    }}
                  />
                  <FieldErrorText>{serviceErrors.price}</FieldErrorText>
                </div>
                <div>
                  <Label>Длительность (мин)</Label>
                  <Input
                    type="number"
                    min="1"
                    value={stForm.duration_minutes}
                    onChange={(e) => {
                      setStForm((f) => ({ ...f, duration_minutes: e.target.value }));
                      setServiceErrors((prev) => ({ ...prev, duration_minutes: undefined }));
                      setManageError(null);
                    }}
                  />
                  <FieldErrorText>{serviceErrors.duration_minutes}</FieldErrorText>
                </div>
              </div>
              <div>
                <Label>Описание</Label>
                <Textarea
                  value={stForm.description}
                  onChange={(e) => {
                    setStForm((f) => ({ ...f, description: e.target.value }));
                    setServiceErrors((prev) => ({ ...prev, description: undefined }));
                    setManageError(null);
                  }}
                  maxLength={1000}
                />
                <FieldErrorText>{serviceErrors.description}</FieldErrorText>
              </div>
              <div className="flex items-center justify-between rounded border px-3 py-2">
                <Label>Доступна для записи</Label>
                <Switch
                  checked={stForm.is_available}
                  onCheckedChange={(v) => {
                    setStForm((f) => ({ ...f, is_available: v }));
                    setManageError(null);
                  }}
                />
              </div>
              {isServiceDirty && (
                <p className="text-xs text-amber-700 bg-amber-50 border border-amber-100 rounded px-2 py-1">
                  Есть несохранённые изменения.
                </p>
              )}
              <DialogFooter>
                <Button type="button" variant="outline" onClick={requestCloseServiceDialog}>
                  Отмена
                </Button>
                <Button
                  type="submit"
                  disabled={
                    servicePending || (Boolean(editingServiceTypeId) && !isServiceDirty)
                  }
                >
                  {editingServiceTypeId ? 'Сохранить' : 'Создать'}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      )}
      {canCatalog && (
        <Dialog
          open={modelDialogOpen}
          onOpenChange={(open) => {
            if (!open) {
              requestCloseModelDialog();
              return;
            }
            setModelDialogOpen(true);
          }}
        >
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle>{editingModelId ? 'Редактировать автомобиль' : 'Добавить автомобиль'}</DialogTitle>
              <DialogDescription>
                Заполните поля модели и при необходимости загрузите изображение.
              </DialogDescription>
            </DialogHeader>
            <form
              className="space-y-3"
              onSubmit={(e) => {
                e.preventDefault();
                submitModelDialog();
              }}
            >
              <div>
                <Label>Бренд *</Label>
                <Select
                  value={modelForm.brand_id}
                  onValueChange={(v) => {
                    setModelForm((f) => ({ ...f, brand_id: v }));
                    setModelErrors((prev) => ({ ...prev, brand_id: undefined }));
                    setManageError(null);
                  }}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Выберите бренд" />
                  </SelectTrigger>
                  <SelectContent>
                    {brands.map((b) => (
                      <SelectItem key={b.brand_id} value={b.brand_id}>
                        {b.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FieldErrorText>{modelErrors.brand_id}</FieldErrorText>
              </div>
              <div>
                <Label>Название модели *</Label>
                <Input
                  value={modelForm.name}
                  onChange={(e) => {
                    setModelForm((f) => ({ ...f, name: e.target.value }));
                    setModelErrors((prev) => ({ ...prev, name: undefined }));
                    setManageError(null);
                  }}
                />
                <FieldErrorText>{modelErrors.name}</FieldErrorText>
              </div>
              <div>
                <Label>Сегмент</Label>
                <Input
                  value={modelForm.segment}
                  onChange={(e) => {
                    setModelForm((f) => ({ ...f, segment: e.target.value }));
                    setModelErrors((prev) => ({ ...prev, segment: undefined }));
                    setManageError(null);
                  }}
                />
                <FieldErrorText>{modelErrors.segment}</FieldErrorText>
              </div>
              <div>
                <Label>Описание</Label>
                <Textarea
                  value={modelForm.description}
                  maxLength={2000}
                  onChange={(e) => {
                    setModelForm((f) => ({ ...f, description: e.target.value }));
                    setModelErrors((prev) => ({ ...prev, description: undefined }));
                    setManageError(null);
                  }}
                />
                <FieldErrorText>{modelErrors.description}</FieldErrorText>
              </div>
              <div>
                <Label>Изображение (JPG, PNG, WEBP)</Label>
                <Input
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  onChange={(e) => {
                    setModelImageFile(e.target.files?.[0] || null);
                    setModelErrors((prev) => ({ ...prev, image: undefined }));
                    setManageError(null);
                  }}
                />
                <FieldErrorText>{modelErrors.image}</FieldErrorText>
              </div>
              {isModelDirty && (
                <p className="text-xs text-amber-700 bg-amber-50 border border-amber-100 rounded px-2 py-1">
                  Есть несохранённые изменения.
                </p>
              )}
              <DialogFooter>
                <Button type="button" variant="outline" onClick={requestCloseModelDialog}>
                  Отмена
                </Button>
                <Button type="submit" disabled={modelPending || (Boolean(editingModelId) && !isModelDirty)}>
                  {editingModelId ? 'Сохранить' : 'Создать'}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}
