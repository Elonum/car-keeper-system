import React, { useState, useCallback } from 'react';
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { SlidersHorizontal, RotateCcw } from 'lucide-react';

export default function FilterPanel({ filters, setFilters, brands, engineTypes, transmissions, driveTypes }) {
  const [open, setOpen] = useState(false);

  const toggleFilter = useCallback((key, value) => {
    setFilters(prev => {
      const arr = prev[key] || [];
      return {
        ...prev,
        [key]: arr.includes(value) ? arr.filter(v => v !== value) : [...arr, value]
      };
    });
  }, [setFilters]);

  const resetFilters = useCallback(() => {
    setFilters({
      brands: [],
      engine_types: [],
      transmissions: [],
      drive_types: [],
      available_only: true,
    });
  }, [setFilters]);

  const activeCount = (filters.brands?.length || 0) + 
    (filters.engine_types?.length || 0) + 
    (filters.transmissions?.length || 0) + 
    (filters.drive_types?.length || 0) +
    (!filters.available_only ? 1 : 0);

  const FilterContent = () => (
    <div className="space-y-6">
      {brands && brands.length > 0 && (
        <div>
          <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Бренд</Label>
          <div className="space-y-2 max-h-48 overflow-y-auto">
            {brands.map(brand => {
              const brandId = brand.brand_id || brand.id;
              return (
                <label key={brandId} className="flex items-center gap-2.5 cursor-pointer group hover:bg-slate-50 p-1.5 rounded transition-colors">
                  <Checkbox
                    checked={(filters.brands || []).includes(brandId)}
                    onCheckedChange={() => toggleFilter('brands', brandId)}
                  />
                  <span className="text-sm text-slate-700 group-hover:text-slate-900 flex-1">{brand.name}</span>
                </label>
              );
            })}
          </div>
        </div>
      )}

      {engineTypes && engineTypes.length > 0 && (
        <div>
          <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Двигатель</Label>
          <div className="space-y-2">
            {engineTypes.map(et => {
              const engineTypeId = et.engine_type_id || et.id;
              return (
                <label key={engineTypeId} className="flex items-center gap-2.5 cursor-pointer group hover:bg-slate-50 p-1.5 rounded transition-colors">
                  <Checkbox
                    checked={(filters.engine_types || []).includes(engineTypeId)}
                    onCheckedChange={() => toggleFilter('engine_types', engineTypeId)}
                  />
                  <span className="text-sm text-slate-700 group-hover:text-slate-900 flex-1">{et.name}</span>
                </label>
              );
            })}
          </div>
        </div>
      )}

      {transmissions && transmissions.length > 0 && (
        <div>
          <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Трансмиссия</Label>
          <div className="space-y-2">
            {transmissions.map(t => {
              const transmissionId = t.transmission_id || t.id;
              return (
                <label key={transmissionId} className="flex items-center gap-2.5 cursor-pointer group hover:bg-slate-50 p-1.5 rounded transition-colors">
                  <Checkbox
                    checked={(filters.transmissions || []).includes(transmissionId)}
                    onCheckedChange={() => toggleFilter('transmissions', transmissionId)}
                  />
                  <span className="text-sm text-slate-700 group-hover:text-slate-900 flex-1">{t.name}</span>
                </label>
              );
            })}
          </div>
        </div>
      )}

      {driveTypes && driveTypes.length > 0 && (
        <div>
          <Label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 block">Привод</Label>
          <div className="space-y-2">
            {driveTypes.map(d => {
              const driveTypeId = d.drive_type_id || d.id;
              return (
                <label key={driveTypeId} className="flex items-center gap-2.5 cursor-pointer group hover:bg-slate-50 p-1.5 rounded transition-colors">
                  <Checkbox
                    checked={(filters.drive_types || []).includes(driveTypeId)}
                    onCheckedChange={() => toggleFilter('drive_types', driveTypeId)}
                  />
                  <span className="text-sm text-slate-700 group-hover:text-slate-900 flex-1">{d.name}</span>
                </label>
              );
            })}
          </div>
        </div>
      )}

      <label className="flex items-center gap-2.5 cursor-pointer hover:bg-slate-50 p-1.5 rounded transition-colors">
        <Checkbox
          checked={filters.available_only !== false}
          onCheckedChange={(checked) => setFilters(prev => ({ ...prev, available_only: checked }))}
        />
        <span className="text-sm text-slate-700">Только в наличии</span>
      </label>

      {activeCount > 0 && (
        <Button variant="outline" onClick={resetFilters} className="w-full gap-2 text-sm" type="button">
          <RotateCcw className="w-3.5 h-3.5" />
          Сбросить фильтры ({activeCount})
        </Button>
      )}
    </div>
  );

  return (
    <>
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

      <div className="lg:hidden">
        <Sheet open={open} onOpenChange={setOpen}>
          <SheetTrigger asChild>
            <Button variant="outline" className="gap-2 rounded-xl" type="button">
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
