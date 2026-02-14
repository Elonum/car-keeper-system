import React from 'react';
import { Link } from 'react-router-dom';
import { createPageUrl } from '../../utils';
import PriceDisplay from '../common/PriceDisplay';
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Fuel, Cog, CircleDot, ArrowRight } from 'lucide-react';

const engineLabels = { petrol: "Бензин", diesel: "Дизель", hybrid: "Гибрид", electric: "Электро" };
const transLabels = { manual: "МКПП", automatic: "АКПП", robot: "Робот", cvt: "Вариатор" };
const driveLabels = { fwd: "Передний", rwd: "Задний", awd: "Полный" };

export default function CatalogCard({ trim }) {
  const imageUrl = trim.image_url || `https://images.unsplash.com/photo-1494976388531-d1058494cdd8?w=600&h=400&fit=crop`;

  return (
    <Card className="group overflow-hidden border-0 shadow-sm hover:shadow-xl transition-all duration-500 bg-white rounded-2xl">
      {/* Image */}
      <div className="relative aspect-[16/10] overflow-hidden bg-slate-100">
        <img
          src={imageUrl}
          alt={`${trim.brand_name || ''} ${trim.model_name || ''}`}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
        />
        {!trim.is_available && (
          <div className="absolute inset-0 bg-black/40 flex items-center justify-center">
            <Badge className="bg-white text-slate-900 text-xs">Нет в наличии</Badge>
          </div>
        )}
        {trim.segment && (
          <Badge className="absolute top-3 left-3 bg-white/90 backdrop-blur text-slate-700 border-0 text-xs font-medium">
            {trim.segment === 'suv' ? 'SUV' : trim.segment === 'sedan' ? 'Седан' : trim.segment === 'hatchback' ? 'Хэтчбек' : trim.segment === 'coupe' ? 'Купе' : trim.segment === 'crossover' ? 'Кроссовер' : trim.segment}
          </Badge>
        )}
      </div>

      {/* Content */}
      <div className="p-5">
        <div className="mb-1">
          <p className="text-xs font-medium text-blue-600 uppercase tracking-wider">
            {trim.brand_name || "Бренд"}
          </p>
          <h3 className="text-lg font-bold text-slate-900 tracking-tight mt-0.5">
            {trim.model_name || "Модель"} {trim.name}
          </h3>
          {trim.generation_name && (
            <p className="text-xs text-slate-400 mt-0.5">{trim.generation_name}</p>
          )}
        </div>

        {/* Specs */}
        <div className="flex flex-wrap gap-2 mt-3 mb-4">
          {trim.engine_type && (
            <div className="flex items-center gap-1 text-xs text-slate-500 bg-slate-50 px-2 py-1 rounded-md">
              <Fuel className="w-3 h-3" />
              <span>{engineLabels[trim.engine_type] || trim.engine_type}</span>
              {trim.engine_volume && <span>• {trim.engine_volume}</span>}
            </div>
          )}
          {trim.transmission && (
            <div className="flex items-center gap-1 text-xs text-slate-500 bg-slate-50 px-2 py-1 rounded-md">
              <Cog className="w-3 h-3" />
              <span>{transLabels[trim.transmission] || trim.transmission}</span>
            </div>
          )}
          {trim.drive_type && (
            <div className="flex items-center gap-1 text-xs text-slate-500 bg-slate-50 px-2 py-1 rounded-md">
              <CircleDot className="w-3 h-3" />
              <span>{driveLabels[trim.drive_type] || trim.drive_type}</span>
            </div>
          )}
          {trim.horsepower && (
            <div className="text-xs text-slate-500 bg-slate-50 px-2 py-1 rounded-md">
              {trim.horsepower} л.с.
            </div>
          )}
        </div>

        {/* Price + CTA */}
        <div className="flex items-end justify-between pt-3 border-t border-slate-100">
          <div>
            <p className="text-xs text-slate-400 mb-0.5">от</p>
            <PriceDisplay price={trim.base_price} size="md" />
          </div>
          <Link to={createPageUrl("Configurator") + `?trim_id=${trim.id}`}>
            <Button size="sm" className="bg-slate-900 hover:bg-blue-600 transition-colors text-white gap-1.5 rounded-xl h-9 px-4">
              Собрать
              <ArrowRight className="w-3.5 h-3.5" />
            </Button>
          </Link>
        </div>
      </div>
    </Card>
  );
}