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
import { asArray } from '@/lib/collections';
import { invalidateGarageRelated } from '@/lib/queryKeys';
import { ErrorNotice, FieldErrorText } from '../common/ErrorNotice';

function pickId(entity, ...keys) {
  if (!entity) return '';
  for (const k of keys) {
    if (entity[k]) return String(entity[k]);
  }
  return '';
}

function CatalogSelectField({
  label,
  value,
  onValueChange,
  disabled,
  placeholder,
  items,
  isLoading,
  isError,
  errorMessage,
  emptyHint,
  fieldError,
  getLabel,
  idKeys,
}) {
  const list = asArray(items);
  const showEmpty = !isLoading && !isError && !disabled && list.length === 0;

  return (
    <div className="space-y-2">
      <Label>{label}</Label>
      <Select value={value || undefined} onValueChange={onValueChange} disabled={disabled || isLoading}>
        <SelectTrigger className={fieldError ? 'border-red-500' : ''}>
          <SelectValue placeholder={isLoading ? 'Загрузка…' : placeholder} />
        </SelectTrigger>
        <SelectContent className="max-h-60">
          {list.map((item) => {
            const id = pickId(item, ...idKeys);
            if (!id) return null;
            return (
              <SelectItem key={id} value={id}>
                {getLabel(item)}
              </SelectItem>
            );
          })}
        </SelectContent>
      </Select>
      {isError && (
        <ErrorNotice kind="inline" message={errorMessage || 'Не удалось загрузить список'} />
      )}
      {showEmpty && emptyHint && (
        <p className="text-sm text-amber-800 bg-amber-50 border border-amber-100 rounded-lg px-3 py-2">
          {emptyHint}
        </p>
      )}
      <FieldErrorText>{fieldError}</FieldErrorText>
    </div>
  );
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
  const [formError, setFormError] = useState(null);

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
      setFormError(null);
    }
  }, [open]);

  const brandsQuery = useQuery({
    queryKey: ['brands'],
    queryFn: () => catalogService.getBrands(),
    enabled: open,
    staleTime: 5 * 60 * 1000,
  });

  const modelsQuery = useQuery({
    queryKey: ['catalog-models', brandId],
    queryFn: () => catalogService.getModels(brandId),
    enabled: open && !!brandId,
  });

  const generationsQuery = useQuery({
    queryKey: ['catalog-generations', modelId],
    queryFn: () => catalogService.getGenerations(modelId),
    enabled: open && !!modelId,
  });

  const trimsQuery = useQuery({
    queryKey: ['catalog-trims-garage', generationId],
    queryFn: () =>
      catalogService.getTrims({
        generation_id: generationId,
        is_available: true,
      }),
    enabled: open && !!generationId,
  });

  const colorsQuery = useQuery({
    queryKey: ['configurator-colors'],
    queryFn: () => configuratorService.getColors(),
    enabled: open,
    staleTime: 5 * 60 * 1000,
  });

  const brands = useMemo(() => asArray(brandsQuery.data), [brandsQuery.data]);
  const models = useMemo(() => asArray(modelsQuery.data), [modelsQuery.data]);
  const generations = useMemo(() => asArray(generationsQuery.data), [generationsQuery.data]);
  const trims = useMemo(() => asArray(trimsQuery.data), [trimsQuery.data]);
  const colors = useMemo(() => asArray(colorsQuery.data), [colorsQuery.data]);

  const availableColors = useMemo(
    () => colors.filter((c) => c.is_available !== false),
    [colors]
  );

  const createMutation = useMutation({
    mutationFn: (payload) => profileService.createUserCar(payload),
    onSuccess: () => {
      invalidateGarageRelated(queryClient);
      setFormError(null);
      onOpenChange(false);
    },
    onError: (e) => {
      setFormError(getApiErrorMessage(e, 'Не удалось добавить автомобиль'));
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
    if (Object.keys(next).length > 0) {
      setFormError('Проверьте поля формы');
    }
    return Object.keys(next).length === 0;
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!validate()) {
      return;
    }
    setFormError(null);
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
          <ErrorNotice kind="form" message={formError} />

          <div className="grid gap-4 sm:grid-cols-2">
            <CatalogSelectField
              label="Марка"
              value={brandId}
              onValueChange={(v) => {
                setBrandId(v);
                setModelId('');
                setGenerationId('');
                setTrimId('');
                setFieldErrors({});
              }}
              placeholder="Выберите марку"
              items={brands}
              isLoading={brandsQuery.isLoading}
              isError={brandsQuery.isError}
              errorMessage={getApiErrorMessage(brandsQuery.error, 'Не удалось загрузить марки')}
              emptyHint="Марки каталога пока не добавлены."
              idKeys={['brand_id', 'id']}
              getLabel={(b) => b.name}
            />

            <CatalogSelectField
              label="Модель"
              value={modelId}
              onValueChange={(v) => {
                setModelId(v);
                setGenerationId('');
                setTrimId('');
                setFieldErrors({});
              }}
              disabled={!brandId}
              placeholder="Выберите модель"
              items={models}
              isLoading={modelsQuery.isLoading}
              isError={modelsQuery.isError}
              errorMessage={getApiErrorMessage(modelsQuery.error, 'Не удалось загрузить модели')}
              emptyHint="Для выбранной марки нет моделей."
              idKeys={['model_id', 'id']}
              getLabel={(m) => m.name}
            />
          </div>

          <CatalogSelectField
            label="Поколение"
            value={generationId}
            onValueChange={(v) => {
              setGenerationId(v);
              setTrimId('');
              setFieldErrors({});
            }}
            disabled={!modelId}
            placeholder="Выберите поколение"
            items={generations}
            isLoading={generationsQuery.isLoading}
            isError={generationsQuery.isError}
            errorMessage={getApiErrorMessage(generationsQuery.error, 'Не удалось загрузить поколения')}
            emptyHint="Для выбранной модели нет поколений."
            idKeys={['generation_id', 'id']}
            getLabel={(g) => {
              if (g.year_from == null) return g.name;
              return g.year_to != null
                ? `${g.name} (${g.year_from}–${g.year_to})`
                : `${g.name} (с ${g.year_from})`;
            }}
          />

          <CatalogSelectField
            label="Комплектация"
            value={trimId}
            onValueChange={(v) => {
              setTrimId(v);
              setFieldErrors((p) => ({ ...p, trim: undefined }));
            }}
            disabled={!generationId}
            placeholder="Выберите комплектацию"
            items={trims}
            isLoading={trimsQuery.isLoading}
            isError={trimsQuery.isError}
            errorMessage={getApiErrorMessage(trimsQuery.error, 'Не удалось загрузить комплектации')}
            emptyHint="Для выбранного поколения нет доступных комплектаций."
            fieldError={fieldErrors.trim}
            idKeys={['trim_id', 'id']}
            getLabel={(t) => [t.brand_name, t.model_name, t.name].filter(Boolean).join(' · ')}
          />

          <CatalogSelectField
            label="Цвет кузова"
            value={colorId}
            onValueChange={(v) => {
              setColorId(v);
              setFieldErrors((p) => ({ ...p, color: undefined }));
            }}
            placeholder="Выберите цвет"
            items={availableColors}
            isLoading={colorsQuery.isLoading}
            isError={colorsQuery.isError}
            errorMessage={getApiErrorMessage(colorsQuery.error, 'Не удалось загрузить цвета')}
            emptyHint="Нет доступных цветов в каталоге."
            fieldError={fieldErrors.color}
            idKeys={['color_id', 'id']}
            getLabel={(c) => (
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
            )}
          />

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
              <FieldErrorText>{fieldErrors.vin}</FieldErrorText>
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
              <FieldErrorText>{fieldErrors.year}</FieldErrorText>
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
              <FieldErrorText>{fieldErrors.mileage}</FieldErrorText>
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
