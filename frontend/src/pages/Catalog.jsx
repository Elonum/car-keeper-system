import React, { useState, useMemo, useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { catalogService } from '@/services/catalogService';
import CatalogCard from '../components/catalog/CatalogCard';
import FilterPanel from '../components/catalog/FilterPanel';
import SortSelect from '../components/catalog/SortSelect';
import PageLoader from '../components/common/PageLoader';
import EmptyState from '../components/common/EmptyState';
import SectionHeader from '../components/common/SectionHeader';
import { Car, Search } from 'lucide-react';
import { Input } from "@/components/ui/input";

export default function Catalog() {
  const [filters, setFilters] = useState({
    brands: [],
    engine_types: [],
    transmissions: [],
    drive_types: [],
    available_only: true,
  });
  const [sortBy, setSortBy] = useState('newest');
  const [searchQuery, setSearchQuery] = useState('');

  const queryParams = useMemo(() => {
    const params = {};
    
    if (filters.brands && filters.brands.length > 0) {
      params.brand_id = filters.brands.join(',');
    }
    
    if (filters.engine_types && filters.engine_types.length > 0) {
      params.engine_type_id = filters.engine_types.join(',');
    }
    
    if (filters.transmissions && filters.transmissions.length > 0) {
      params.transmission_id = filters.transmissions.join(',');
    }
    
    if (filters.drive_types && filters.drive_types.length > 0) {
      params.drive_type_id = filters.drive_types.join(',');
    }
    
    if (filters.available_only !== undefined) {
      params.is_available = filters.available_only;
    }
    
    return params;
  }, [filters]);

  const queryKey = useMemo(() => ['trims', queryParams], [queryParams]);

  const { data: trims, isLoading: trimsLoading } = useQuery({
    queryKey,
    queryFn: () => catalogService.getTrims(queryParams),
    keepPreviousData: true,
    refetchOnWindowFocus: false,
    refetchOnMount: false,
    refetchOnReconnect: false,
    staleTime: 5 * 60 * 1000,
  });

  const { data: brands } = useQuery({
    queryKey: ['brands'],
    queryFn: () => catalogService.getBrands(),
  });

  const { data: engineTypes } = useQuery({
    queryKey: ['engineTypes'],
    queryFn: () => catalogService.getEngineTypes(),
  });

  const { data: transmissions } = useQuery({
    queryKey: ['transmissions'],
    queryFn: () => catalogService.getTransmissions(),
  });

  const { data: driveTypes } = useQuery({
    queryKey: ['driveTypes'],
    queryFn: () => catalogService.getDriveTypes(),
  });

  const filteredAndSortedTrims = useMemo(() => {
    if (!trims) return [];

    let result = [...trims];

    if (searchQuery) {
      const query = searchQuery.toLowerCase().trim();
      result = result.filter(t => 
        (t.brand_name?.toLowerCase().includes(query)) ||
        (t.model_name?.toLowerCase().includes(query)) ||
        (t.name?.toLowerCase().includes(query)) ||
        (t.generation_name?.toLowerCase().includes(query))
      );
    }

    if (sortBy === 'price_asc') {
      result.sort((a, b) => (a.base_price || 0) - (b.base_price || 0));
    } else if (sortBy === 'price_desc') {
      result.sort((a, b) => (b.base_price || 0) - (a.base_price || 0));
    } else if (sortBy === 'name_asc') {
      result.sort((a, b) => {
        const nameA = `${a.brand_name || ''} ${a.model_name || ''} ${a.name || ''}`.toLowerCase();
        const nameB = `${b.brand_name || ''} ${b.model_name || ''} ${b.name || ''}`.toLowerCase();
        return nameA.localeCompare(nameB, 'ru');
      });
    } else if (sortBy === 'newest') {
      result.sort((a, b) => {
        const dateA = a.created_at ? new Date(a.created_at).getTime() : 0;
        const dateB = b.created_at ? new Date(b.created_at).getTime() : 0;
        return dateB - dateA;
      });
    }

    return result;
  }, [trims, sortBy, searchQuery]);

  if (trimsLoading) return <PageLoader />;

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <SectionHeader 
          title="Каталог автомобилей"
          description="Выберите автомобиль и создайте свою идеальную конфигурацию"
        />

        <div className="flex gap-6 lg:gap-8">
          <FilterPanel 
            filters={filters} 
            setFilters={setFilters} 
            brands={brands || []}
            engineTypes={engineTypes || []}
            transmissions={transmissions || []}
            driveTypes={driveTypes || []}
          />

          <div className="flex-1 min-w-0">
            <div className="flex flex-col sm:flex-row gap-3 mb-6">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                <Input 
                  placeholder="Поиск по бренду, модели..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 rounded-xl bg-white"
                />
              </div>
              <SortSelect value={sortBy} onChange={setSortBy} />
            </div>

            <p className="text-sm text-slate-500 mb-4">
              Найдено автомобилей: <span className="font-semibold text-slate-900">{filteredAndSortedTrims.length}</span>
            </p>

            {filteredAndSortedTrims.length === 0 ? (
              <EmptyState 
                icon={Car}
                title="Ничего не найдено"
                description="Попробуйте изменить фильтры или поисковый запрос"
              />
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {filteredAndSortedTrims.map(trim => (
                  <CatalogCard key={trim.trim_id || trim.id} trim={trim} />
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
