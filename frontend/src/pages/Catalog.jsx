import React, { useState, useMemo } from 'react';
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
    price_range: [0, 15000000],
    available_only: true,
  });
  const [sortBy, setSortBy] = useState('newest');
  const [searchQuery, setSearchQuery] = useState('');

  const { data: trims, isLoading: trimsLoading } = useQuery({
    queryKey: ['trims', filters],
    queryFn: () => catalogService.getTrims({
      brand_id: filters.brands.length > 0 ? filters.brands : undefined,
      engine_type_id: filters.engine_types.length > 0 ? filters.engine_types : undefined,
      transmission_id: filters.transmissions.length > 0 ? filters.transmissions : undefined,
      drive_type_id: filters.drive_types.length > 0 ? filters.drive_types : undefined,
      min_price: filters.price_range[0],
      max_price: filters.price_range[1],
      is_available: filters.available_only,
    }),
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

    // Search filter
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      result = result.filter(t => 
        (t.brand_name?.toLowerCase().includes(query)) ||
        (t.model_name?.toLowerCase().includes(query)) ||
        (t.name?.toLowerCase().includes(query)) ||
        (t.trim_name?.toLowerCase().includes(query))
      );
    }

    // Sorting
    if (sortBy === 'price_asc') {
      result.sort((a, b) => (a.base_price || 0) - (b.base_price || 0));
    } else if (sortBy === 'price_desc') {
      result.sort((a, b) => (b.base_price || 0) - (a.base_price || 0));
    } else if (sortBy === 'name_asc') {
      result.sort((a, b) => {
        const nameA = `${a.brand_name || ''} ${a.model_name || ''}`.toLowerCase();
        const nameB = `${b.brand_name || ''} ${b.model_name || ''}`.toLowerCase();
        return nameA.localeCompare(nameB);
      });
    } else if (sortBy === 'newest') {
      result.sort((a, b) => new Date(b.created_at || 0) - new Date(a.created_at || 0));
    }

    return result;
  }, [trims, filters, sortBy, searchQuery]);

  if (trimsLoading) return <PageLoader />;

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <SectionHeader 
          title="Каталог автомобилей"
          description="Выберите автомобиль и создайте свою идеальную конфигурацию"
        />

        <div className="flex gap-6 lg:gap-8">
          {/* Filters */}
          <FilterPanel 
            filters={filters} 
            setFilters={setFilters} 
            brands={brands || []}
            engineTypes={engineTypes || []}
            transmissions={transmissions || []}
            driveTypes={driveTypes || []}
          />

          {/* Main content */}
          <div className="flex-1 min-w-0">
            {/* Search + Sort */}
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

            {/* Results count */}
            <p className="text-sm text-slate-500 mb-4">
              Найдено автомобилей: <span className="font-semibold text-slate-900">{filteredAndSortedTrims.length}</span>
            </p>

            {/* Grid */}
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
