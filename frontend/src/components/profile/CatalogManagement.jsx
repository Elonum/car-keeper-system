import React, { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { catalogService } from '@/services/catalogService';
import { serviceService } from '@/services/serviceService';
import { adminCatalogService } from '@/services/adminCatalogService';
import { PERMISSIONS, hasPermission } from '@/lib/authz';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { toast } from 'sonner';

const SERVICE_CATEGORIES = ['maintenance', 'repair', 'diagnostics', 'detailing', 'tires'];

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

  const [brandForm, setBrandForm] = useState({ name: '', country: '' });
  const [stForm, setStForm] = useState({
    name: '',
    category: 'maintenance',
    price: '',
    duration_minutes: '',
    is_available: true,
  });

  const createBrand = useMutation({
    mutationFn: () => adminCatalogService.createBrand({ name: brandForm.name.trim(), country: brandForm.country.trim() }),
    onSuccess: () => {
      toast.success('Бренд добавлен');
      setBrandForm({ name: '', country: '' });
      qc.invalidateQueries({ queryKey: ['catalog', 'brands'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось создать бренд')),
  });

  const deleteBrand = useMutation({
    mutationFn: (id) => adminCatalogService.deleteBrand(id),
    onSuccess: () => {
      toast.success('Бренд удалён');
      qc.invalidateQueries({ queryKey: ['catalog', 'brands'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось удалить')),
  });

  const createSt = useMutation({
    mutationFn: () =>
      adminCatalogService.createServiceType({
        name: stForm.name.trim(),
        category: stForm.category,
        price: Number(stForm.price) || 0,
        duration_minutes: stForm.duration_minutes ? Number(stForm.duration_minutes) : null,
        is_available: stForm.is_available,
      }),
    onSuccess: () => {
      toast.success('Услуга добавлена');
      setStForm({
        name: '',
        category: 'maintenance',
        price: '',
        duration_minutes: '',
        is_available: true,
      });
      qc.invalidateQueries({ queryKey: ['service', 'types', 'admin'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось создать услугу')),
  });

  const patchSt = useMutation({
    mutationFn: ({ id, payload }) => adminCatalogService.updateServiceType(id, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['service', 'types', 'admin'] }),
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось обновить')),
  });

  const patchBranch = useMutation({
    mutationFn: ({ id, payload }) => adminCatalogService.updateBranch(id, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['service', 'branches'] }),
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось обновить филиал')),
  });

  const deleteSt = useMutation({
    mutationFn: (id) => adminCatalogService.deleteServiceType(id),
    onSuccess: () => {
      toast.success('Услуга удалена');
      qc.invalidateQueries({ queryKey: ['service', 'types', 'admin'] });
    },
    onError: (e) => toast.error(getApiErrorMessage(e, 'Не удалось удалить')),
  });

  if (!canCatalog && !canService) return null;

  return (
    <div className="space-y-6">
      {canCatalog && (
        <Card className="p-6 space-y-4">
          <div>
            <h3 className="text-lg font-semibold text-slate-900">Каталог: бренды</h3>
            <p className="text-sm text-slate-600 mt-1">
              Добавление и удаление брендов. Модели и комплектации редактируются отдельно (при необходимости — через БД или расширение API).
            </p>
          </div>
          <div className="flex flex-wrap gap-2 max-h-40 overflow-y-auto">
            {brands.map((b) => (
              <div
                key={b.brand_id}
                className="flex items-center gap-2 rounded-lg border px-2 py-1 text-sm"
              >
                <span>
                  {b.name} ({b.country})
                </span>
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
          <div className="grid sm:grid-cols-2 gap-3 max-w-lg">
            <div>
              <Label>Название</Label>
              <Input value={brandForm.name} onChange={(e) => setBrandForm((f) => ({ ...f, name: e.target.value }))} />
            </div>
            <div>
              <Label>Страна</Label>
              <Input value={brandForm.country} onChange={(e) => setBrandForm((f) => ({ ...f, country: e.target.value }))} />
            </div>
          </div>
          <Button
            type="button"
            onClick={() => createBrand.mutate()}
            disabled={!brandForm.name.trim() || !brandForm.country.trim() || createBrand.isPending}
          >
            Добавить бренд
          </Button>
        </Card>
      )}

      {canService && (
        <Card className="p-6 space-y-4">
          <h3 className="text-lg font-semibold text-slate-900">Сервис: услуги</h3>
          <div className="space-y-2 max-h-56 overflow-y-auto">
            {serviceTypes.map((st) => (
              <div
                key={st.service_type_id}
                className="flex flex-wrap items-center justify-between gap-2 rounded-lg border p-2 text-sm"
              >
                <div>
                  <span className="font-medium">{st.name}</span>
                  <span className="text-slate-500 ml-2">
                    {st.category} · {st.price} ₽
                  </span>
                </div>
                <div className="flex items-center gap-2">
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
          <div className="grid sm:grid-cols-2 gap-3 max-w-2xl border-t pt-4">
            <div>
              <Label>Название</Label>
              <Input value={stForm.name} onChange={(e) => setStForm((f) => ({ ...f, name: e.target.value }))} />
            </div>
            <div>
              <Label>Категория</Label>
              <Select value={stForm.category} onValueChange={(v) => setStForm((f) => ({ ...f, category: v }))}>
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
            </div>
            <div>
              <Label>Цена (₽)</Label>
              <Input
                type="number"
                min="0"
                step="0.01"
                value={stForm.price}
                onChange={(e) => setStForm((f) => ({ ...f, price: e.target.value }))}
              />
            </div>
            <div>
              <Label>Длительность (мин), опционально</Label>
              <Input
                type="number"
                min="1"
                value={stForm.duration_minutes}
                onChange={(e) => setStForm((f) => ({ ...f, duration_minutes: e.target.value }))}
              />
            </div>
            <div className="flex items-center justify-between sm:col-span-2 rounded border px-3 py-2">
              <Label>Доступна для записи</Label>
              <Switch
                checked={stForm.is_available}
                onCheckedChange={(v) => setStForm((f) => ({ ...f, is_available: v }))}
              />
            </div>
          </div>
          <Button
            type="button"
            onClick={() => createSt.mutate()}
            disabled={!stForm.name.trim() || createSt.isPending}
          >
            Добавить услугу
          </Button>

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
    </div>
  );
}
