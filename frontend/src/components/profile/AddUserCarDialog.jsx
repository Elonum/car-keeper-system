import React, { useEffect, useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { catalogService } from '@/services/catalogService';
import { configuratorService } from '@/services/configuratorService';
import { profileService } from '@/services/profileService';
import {
  normalizeVIN,
  validateVIN,
  validateVehicleYear,
  validateMileage,
} from '@/lib/userCarValidation';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { toast } from 'sonner';

const GARAGE_QUERY_KEYS = [
  ['my-cars'],
  ['userCars'],
];

function pickId(entity, ...keys) {
  if (!entity) return '';
  for (const k of keys) {
    if (entity[k]) return String(entity[k]);
  }
  return '';
}

export default function AddUserCarDialog({ open, onOpenChange }) {
  const queryClient = useQueryClient();
  const [brandId, setBrandId] = useState('');
  const [modelId, setModelId] = useState('');
  const [generationId, setGenerationId] = useState('');
  const [trimId, setTrimId] = useState('');
  const [colorId, setColorId] = useState('');
  const [vin, setVin] = useState('');
  const [year, setYear] = useState('');
  const [mileage, setMileage] = useState('0');
  const [purchaseDate, setPurchaseDate] = useState('');
  const [fieldErrors, setFieldErrors] = useState({});

  const maxYear = useMemo(() => new Date().getFullYear() + 1, []);

  useEffect(() => {
    if (!open) {
      setBrandId('');
      setModelId('');
      setGenerationId('');
      setTrimId('');
      setColorId('');
      setVin('');
      setYear('');
      setMileage('0');
      setPurchaseDate('');
      setFieldErrors({});
    }
  }, [open]);

  const { data: brands = [] } = useQuery({
    queryKey: ['brands'],
    queryFn: () => catalogService.getBrands(),
    enabled: open,
    staleTime: 5 * 60 * 1000,
  });

  const { data: models = [] } = useQuery({
    queryKey: ['catalog-models', brandId],
    queryFn: () => catalogService.getModels(brandId),
    enabled: open && !!brandId,
  });

  const { data: generations = [] } = useQuery({
    queryKey: ['catalog-generations', modelId],
    queryFn: () => catalogService.getGenerations(modelId),
    enabled: open && !!modelId,
  });

  const { data: trims = [] } = useQuery({
    queryKey: ['catalog-trims-garage', generationId],
    queryFn: () =>
      catalogService.getTrims({
        generation_id: generationId,
        is_available: true,
      }),
    enabled: open && !!generationId,
  });

  const { data: colors = [] } = useQuery({
    queryKey: ['configurator-colors'],
    queryFn: () => configuratorService.getColors(),
    enabled: open,
    staleTime: 5 * 60 * 1000,
  });

  const createMutation = useMutation({
    mutationFn: (payload) => profileService.createUserCar(payload),
    onSuccess: () => {
      GARAGE_QUERY_KEYS.forEach((key) => {
        queryClient.invalidateQueries({ queryKey: key });
      });
      toast.success('Автомобиль добавлен в гараж');
      onOpenChange(false);
    },
    onError: (e) => {
      toast.error(getApiErrorMessage(e, 'Не удалось добавить автомобиль'));
    },
  });

  const validate = () => {
    const next = {};
    if (!trimId) next.trim = 'Выберите комплектацию';
    if (!colorId) next.color = 'Выберите цвет';
    const vErr = validateVIN(vin);
    if (vErr) next.vin = vErr;
    const yErr = validateVehicleYear(year);
    if (yErr) next.year = yErr;
    const mErr = validateMileage(mileage);
    if (mErr) next.mileage = mErr;
    setFieldErrors(next);
    return Object.keys(next).length === 0;
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!validate()) {
      toast.error('Проверьте поля формы');
      return;
    }
    const payload = {
      trim_id: trimId,
      color_id: colorId,
      vin: normalizeVIN(vin),
      year: Number(year),
      current_mileage: mileage === '' ? 0 : Number(mileage),
    };
    if (purchaseDate) {
      payload.purchase_date = purchaseDate;
    }
    createMutation.mutate(payload);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Добавить автомобиль</DialogTitle>
          <DialogDescription>
            Укажите модель из каталога, цвет кузова и данные по VIN. Один VIN может быть только у одного
            пользователя.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4" noValidate>
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label>Марка</Label>
              <Select
                value={brandId || undefined}
                onValueChange={(v) => {
                  setBrandId(v);
                  setModelId('');
                  setGenerationId('');
                  setTrimId('');
                  setFieldErrors({});
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Выберите марку" />
                </SelectTrigger>
                <SelectContent>
                  {brands.map((b) => {
                    const id = pickId(b, 'brand_id', 'id');
                    return (
                      <SelectItem key={id} value={id}>
                        {b.name}
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>Модель</Label>
              <Select
                value={modelId || undefined}
                onValueChange={(v) => {
                  setModelId(v);
                  setGenerationId('');
                  setTrimId('');
                  setFieldErrors({});
                }}
                disabled={!brandId}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Выберите модель" />
                </SelectTrigger>
                <SelectContent>
                  {models.map((m) => {
                    const id = pickId(m, 'model_id', 'id');
                    return (
                      <SelectItem key={id} value={id}>
                        {m.name}
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label>Поколение</Label>
            <Select
              value={generationId || undefined}
              onValueChange={(v) => {
                setGenerationId(v);
                setTrimId('');
                setFieldErrors({});
              }}
              disabled={!modelId}
            >
              <SelectTrigger>
                <SelectValue placeholder="Выберите поколение" />
              </SelectTrigger>
              <SelectContent>
                {generations.map((g) => {
                  const id = pickId(g, 'generation_id', 'id');
                  let label = g.name;
                  if (g.year_from != null) {
                    label =
                      g.year_to != null
                        ? `${g.name} (${g.year_from}–${g.year_to})`
                        : `${g.name} (с ${g.year_from})`;
                  }
                  return (
                    <SelectItem key={id} value={id}>
                      {label}
                    </SelectItem>
                  );
                })}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Комплектация</Label>
            <Select
              value={trimId || undefined}
              onValueChange={(v) => {
                setTrimId(v);
                setFieldErrors((p) => ({ ...p, trim: undefined }));
              }}
              disabled={!generationId}
            >
              <SelectTrigger className={fieldErrors.trim ? 'border-red-500' : ''}>
                <SelectValue placeholder="Выберите комплектацию" />
              </SelectTrigger>
              <SelectContent className="max-h-60">
                {trims.map((t) => {
                  const id = pickId(t, 'trim_id', 'id');
                  const label = [t.brand_name, t.model_name, t.name].filter(Boolean).join(' · ');
                  return (
                    <SelectItem key={id} value={id}>
                      {label}
                    </SelectItem>
                  );
                })}
              </SelectContent>
            </Select>
            {fieldErrors.trim && <p className="text-sm text-red-600">{fieldErrors.trim}</p>}
          </div>

          <div className="space-y-2">
            <Label>Цвет кузова</Label>
            <Select
              value={colorId || undefined}
              onValueChange={(v) => {
                setColorId(v);
                setFieldErrors((p) => ({ ...p, color: undefined }));
              }}
            >
              <SelectTrigger className={fieldErrors.color ? 'border-red-500' : ''}>
                <SelectValue placeholder="Выберите цвет" />
              </SelectTrigger>
              <SelectContent className="max-h-60">
                {colors
                  .filter((c) => c.is_available !== false)
                  .map((c) => {
                    const id = pickId(c, 'color_id', 'id');
                    return (
                      <SelectItem key={id} value={id}>
                        <span className="flex items-center gap-2">
                          {c.hex_code && (
                            <span
                              className="inline-block h-3 w-3 rounded-full border border-slate-200"
                              style={{ backgroundColor: c.hex_code }}
                              aria-hidden
                            />
                          )}
                          {c.name}
                        </span>
                      </SelectItem>
                    );
                  })}
              </SelectContent>
            </Select>
            {fieldErrors.color && <p className="text-sm text-red-600">{fieldErrors.color}</p>}
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="garage-vin">VIN</Label>
              <Input
                id="garage-vin"
                value={vin}
                onChange={(e) => {
                  const v = normalizeVIN(e.target.value).slice(0, 17);
                  setVin(v);
                  setFieldErrors((p) => ({ ...p, vin: undefined }));
                }}
                maxLength={17}
                autoComplete="off"
                spellCheck={false}
                className={`font-mono uppercase ${fieldErrors.vin ? 'border-red-500' : ''}`}
                placeholder="17 символов"
              />
              {fieldErrors.vin && <p className="text-sm text-red-600">{fieldErrors.vin}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="garage-year">Год выпуска</Label>
              <Input
                id="garage-year"
                type="number"
                min={1900}
                max={maxYear}
                value={year}
                onChange={(e) => {
                  setYear(e.target.value);
                  setFieldErrors((p) => ({ ...p, year: undefined }));
                }}
                className={fieldErrors.year ? 'border-red-500' : ''}
              />
              {fieldErrors.year && <p className="text-sm text-red-600">{fieldErrors.year}</p>}
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="garage-mileage">Пробег, км</Label>
              <Input
                id="garage-mileage"
                type="number"
                min={0}
                step={1}
                value={mileage}
                onChange={(e) => {
                  setMileage(e.target.value);
                  setFieldErrors((p) => ({ ...p, mileage: undefined }));
                }}
                className={fieldErrors.mileage ? 'border-red-500' : ''}
              />
              {fieldErrors.mileage && <p className="text-sm text-red-600">{fieldErrors.mileage}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="garage-purchase">Дата покупки (необязательно)</Label>
              <Input
                id="garage-purchase"
                type="date"
                value={purchaseDate}
                max={new Date().toISOString().slice(0, 10)}
                onChange={(e) => setPurchaseDate(e.target.value)}
              />
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Отмена
            </Button>
            <Button type="submit" disabled={createMutation.isPending} className="bg-slate-900 hover:bg-slate-800">
              {createMutation.isPending ? 'Сохранение…' : 'Добавить'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
