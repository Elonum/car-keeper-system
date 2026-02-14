import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { serviceService } from '@/services/serviceService';
import { Link } from 'react-router-dom';
import { createPageUrl } from '../utils';
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import PageLoader from '../components/common/PageLoader';
import PriceDisplay from '../components/common/PriceDisplay';
import SectionHeader from '../components/common/SectionHeader';
import { Wrench, Clock } from 'lucide-react';

const categoryLabels = {
  maintenance: "ТО и обслуживание",
  repair: "Ремонт",
  diagnostics: "Диагностика",
  detailing: "Детейлинг",
  tires: "Шиномонтаж"
};

const categoryColors = {
  maintenance: "bg-blue-50 text-blue-700 border-blue-200",
  repair: "bg-red-50 text-red-700 border-red-200",
  diagnostics: "bg-purple-50 text-purple-700 border-purple-200",
  detailing: "bg-green-50 text-green-700 border-green-200",
  tires: "bg-orange-50 text-orange-700 border-orange-200"
};

export default function Services() {
  const [selectedCategory, setSelectedCategory] = useState('all');

  const { data: services = [], isLoading } = useQuery({
    queryKey: ['serviceTypes', { is_available: true }],
    queryFn: () => serviceService.getServiceTypes({ is_available: true }),
  });

  const categories = ['all', ...Object.keys(categoryLabels)];

  const filteredServices = selectedCategory === 'all' 
    ? services 
    : services.filter(s => s.category === selectedCategory);

  const serviceId = (s) => s.service_type_id || s.id;

  const groupedServices = filteredServices.reduce((acc, service) => {
    const cat = service.category || 'other';
    if (!acc[cat]) acc[cat] = [];
    acc[cat].push(service);
    return acc;
  }, {});

  if (isLoading) return <PageLoader />;

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <SectionHeader 
          title="Сервисные услуги"
          description="Выберите услуги для вашего автомобиля"
          action={
            <Link to={createPageUrl("ServiceAppointment")}>
              <Button className="bg-blue-600 hover:bg-blue-700 gap-2">
                <Wrench className="w-4 h-4" />
                Записаться на ТО
              </Button>
            </Link>
          }
        />

        {/* Category filters */}
        <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
          {categories.map(cat => (
            <button
              key={cat}
              onClick={() => setSelectedCategory(cat)}
              className={`px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-all ${
                selectedCategory === cat
                  ? 'bg-slate-900 text-white'
                  : 'bg-white border border-slate-200 text-slate-600 hover:border-slate-300'
              }`}
            >
              {cat === 'all' ? 'Все услуги' : categoryLabels[cat]}
            </button>
          ))}
        </div>

        {/* Services grid */}
        {Object.keys(groupedServices).length > 0 ? (
          Object.entries(groupedServices).map(([category, items]) => (
            <div key={category} className="mb-8">
              <div className="flex items-center gap-2 mb-4">
                <Badge variant="outline" className={`${categoryColors[category]} border text-sm px-3 py-1.5`}>
                  {categoryLabels[category] || category}
                </Badge>
                <div className="h-px flex-1 bg-slate-200" />
              </div>

              <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                {items.map(service => (
                  <Card key={serviceId(service)} className="p-5 border-0 shadow-sm hover:shadow-md transition-shadow">
                    <div className="flex items-start justify-between mb-3">
                      <h3 className="font-bold text-slate-900 text-lg">{service.name}</h3>
                      <PriceDisplay price={service.price} size="md" />
                    </div>

                    {service.description && (
                      <p className="text-sm text-slate-600 mb-3 line-clamp-2">{service.description}</p>
                    )}

                    {service.duration_minutes && (
                      <div className="flex items-center gap-1.5 text-xs text-slate-400">
                        <Clock className="w-3.5 h-3.5" />
                        <span>{service.duration_minutes} мин</span>
                      </div>
                    )}
                  </Card>
                ))}
              </div>
            </div>
          ))
        ) : (
          <div className="text-center py-12">
            <Wrench className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <p className="text-slate-500">Список услуг будет доступен после настройки backend</p>
          </div>
        )}

        {filteredServices.length === 0 && services.length > 0 && (
          <div className="text-center py-12">
            <Wrench className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <p className="text-slate-500">Нет доступных услуг в этой категории</p>
          </div>
        )}
      </div>
    </div>
  );
}
