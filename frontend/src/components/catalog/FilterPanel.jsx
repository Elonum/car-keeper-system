import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Slider } from "@/components/ui/slider";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { SlidersHorizontal, X, RotateCcw } from 'lucide-react';
import PriceDisplay from '../common/PriceDisplay';

const engineTypes = [
  { value: "petrol", label: "Бензин" },
  { value: "diesel", label: "Дизель" },
  { value: "hybrid", label: "Гибрид" },
  { value: "electric", label: "Электро" },
];

const transmissions = [
  { value: "manual", label: "МКПП" },
  { value: "automatic", label: "АКПП" },
  { value: "robot", label: "Робот" },
  { value: "cvt", label: "Вариатор" },
];

const driveTypes = [
  { value: "fwd", label: "Передний" },
  { value: "rwd", label: "Задний" },
  { value: "awd", label: "Полный" },
];

export default function FilterPanel({ filters, setFilters, brands }) {
  const [open, setOpen] = useState(false);

  const toggleFilter = (key, value) => {
    setFilters(prev => {
      const arr = prev[key] || [];
      return {
        ...prev,
        [key]: arr.includes(value) ? arr.filter(v => v !== value) : [...arr, value]
      };
    });
  };

  const resetFilters = () => {
    setFilters({
      brands: [],
      engine_types: [],
      transmissions: [],
      drive_types: [],
      price_range: [0, 15000000],
      available_only: true,
    });
  };

  const activeCount = (filters.brands?.length || 0) + 
    (filters.engine_types?.length || 0) + 
    (filters.transmissions?.length || 0) + 
    (filters.drive_types?.length || 0) +
    (!filters.available_only ? 1 : 0);

  const FilterContent = () => (
    <div className="space-y-6">
      {/* Brands */}
      {brands && brands.length > 0 && (
        <div>
          <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Бренд</Label>
          <div className="space-y-2">
            {brands.map(brand => {
              const brandId = brand.brand_id || brand.id;
              return (
                <label key={brandId} className="flex items-center gap-2.5 cursor-pointer group">
                  <Checkbox
                    checked={(filters.brands || []).includes(brandId)}
                    onCheckedChange={() => toggleFilter('brands', brandId)}
                  />
                  <span className="text-sm text-slate-700 group-hover:text-slate-900">{brand.name}</span>
                </label>
              );
            })}
          </div>
        </div>
      )}

      {/* Engine Type */}
      <div>
        <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Двигатель</Label>
        <div className="space-y-2">
          {engineTypes.map(et => (
            <label key={et.value} className="flex items-center gap-2.5 cursor-pointer group">
              <Checkbox
                checked={(filters.engine_types || []).includes(et.value)}
                onCheckedChange={() => toggleFilter('engine_types', et.value)}
              />
              <span className="text-sm text-slate-700 group-hover:text-slate-900">{et.label}</span>
            </label>
          ))}
        </div>
      </div>

      {/* Transmission */}
      <div>
        <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Трансмиссия</Label>
        <div className="space-y-2">
          {transmissions.map(t => (
            <label key={t.value} className="flex items-center gap-2.5 cursor-pointer group">
              <Checkbox
                checked={(filters.transmissions || []).includes(t.value)}
                onCheckedChange={() => toggleFilter('transmissions', t.value)}
              />
              <span className="text-sm text-slate-700 group-hover:text-slate-900">{t.label}</span>
            </label>
          ))}
        </div>
      </div>

      {/* Drive Type */}
      <div>
        <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Привод</Label>
        <div className="space-y-2">
          {driveTypes.map(d => (
            <label key={d.value} className="flex items-center gap-2.5 cursor-pointer group">
              <Checkbox
                checked={(filters.drive_types || []).includes(d.value)}
                onCheckedChange={() => toggleFilter('drive_types', d.value)}
              />
              <span className="text-sm text-slate-700 group-hover:text-slate-900">{d.label}</span>
            </label>
          ))}
        </div>
      </div>

      {/* Price Range */}
      <div>
        <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Цена</Label>
        <div className="px-1">
          <Slider
            min={0}
            max={15000000}
            step={100000}
            value={filters.price_range || [0, 15000000]}
            onValueChange={(val) => setFilters(prev => ({ ...prev, price_range: val }))}
            className="mb-3"
          />
          <div className="flex justify-between">
            <PriceDisplay price={(filters.price_range || [0])[0]} size="sm" className="text-slate-500" />
            <PriceDisplay price={(filters.price_range || [0, 15000000])[1]} size="sm" className="text-slate-500" />
          </div>
        </div>
      </div>

      {/* Available only */}
      <label className="flex items-center gap-2.5 cursor-pointer">
        <Checkbox
          checked={filters.available_only !== false}
          onCheckedChange={(checked) => setFilters(prev => ({ ...prev, available_only: checked }))}
        />
        <span className="text-sm text-slate-700">Только в наличии</span>
      </label>

      <Button variant="outline" onClick={resetFilters} className="w-full gap-2 text-sm">
        <RotateCcw className="w-3.5 h-3.5" />
        Сбросить фильтры
      </Button>
    </div>
  );

  return (
    <>
      {/* Desktop sidebar */}
      <div className="hidden lg:block w-64 flex-shrink-0">
        <div className="sticky top-20 bg-white rounded-2xl p-5 shadow-sm border border-slate-100">
          <div className="flex items-center justify-between mb-5">
            <h3 className="font-semibold text-slate-900">Фильтры</h3>
            {activeCount > 0 && (
              <span className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full font-medium">
                {activeCount}
              </span>
            )}
          </div>
          <FilterContent />
        </div>
      </div>

      {/* Mobile filter sheet */}
      <div className="lg:hidden">
        <Sheet open={open} onOpenChange={setOpen}>
          <SheetTrigger asChild>
            <Button variant="outline" className="gap-2 rounded-xl">
              <SlidersHorizontal className="w-4 h-4" />
              Фильтры
              {activeCount > 0 && (
                <span className="ml-1 text-xs bg-blue-600 text-white px-1.5 py-0.5 rounded-full">
                  {activeCount}
                </span>
              )}
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-80 overflow-y-auto">
            <SheetHeader>
              <SheetTitle>Фильтры</SheetTitle>
            </SheetHeader>
            <div className="mt-6">
              <FilterContent />
            </div>
          </SheetContent>
        </Sheet>
      </div>
    </>
  );
}