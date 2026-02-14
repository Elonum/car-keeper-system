import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '@/lib/AuthContext';
import { useSearchParams } from 'react-router-dom';
import { profileService } from '@/services/profileService';
import { orderService } from '@/services/orderService';
import { serviceService } from '@/services/serviceService';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import PageLoader from '../components/common/PageLoader';
import SectionHeader from '../components/common/SectionHeader';
import ConfigurationList from '../components/profile/ConfigurationList';
import OrdersList from '../components/profile/OrdersList';
import UserCarsList from '../components/profile/UserCarsList';
import ServiceAppointmentsList from '../components/profile/ServiceAppointmentsList';
import { Car, Settings, FileText, Wrench, ShoppingCart } from 'lucide-react';

export default function Profile() {
  const [searchParams] = useSearchParams();
  const initialTab = searchParams.get('tab') || 'configurations';
  const [activeTab, setActiveTab] = useState(initialTab);
  const { user, isAuthenticated, navigateToLogin } = useAuth();

  useEffect(() => {
    if (!isAuthenticated) {
      navigateToLogin();
    }
  }, [isAuthenticated, navigateToLogin]);

  const { data: configurations, isLoading: configsLoading } = useQuery({
    queryKey: ['my-configurations'],
    queryFn: () => profileService.getConfigurations(),
    enabled: isAuthenticated,
  });

  const { data: orders, isLoading: ordersLoading } = useQuery({
    queryKey: ['my-orders'],
    queryFn: () => orderService.getOrders(),
    enabled: isAuthenticated,
  });

  const { data: userCars, isLoading: carsLoading } = useQuery({
    queryKey: ['my-cars'],
    queryFn: () => profileService.getUserCars(),
    enabled: isAuthenticated,
  });

  const { data: appointments, isLoading: appointmentsLoading } = useQuery({
    queryKey: ['my-appointments'],
    queryFn: () => serviceService.getAppointments(),
    enabled: isAuthenticated,
  });

  if (!isAuthenticated || !user) return <PageLoader />;

  const getUserDisplayName = () => {
    if (user.full_name) return user.full_name;
    if (user.first_name && user.last_name) {
      return `${user.first_name} ${user.last_name}`;
    }
    return user.email || 'Пользователь';
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <SectionHeader 
          title="Личный кабинет"
          description={`Добро пожаловать, ${getUserDisplayName()}`}
        />

        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="bg-white p-1.5 h-auto rounded-xl shadow-sm">
            <TabsTrigger value="configurations" className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-4 py-2.5">
              <Settings className="w-4 h-4" />
              <span className="hidden sm:inline">Конфигурации</span>
            </TabsTrigger>
            <TabsTrigger value="orders" className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-4 py-2.5">
              <ShoppingCart className="w-4 h-4" />
              <span className="hidden sm:inline">Заказы</span>
            </TabsTrigger>
            <TabsTrigger value="cars" className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-4 py-2.5">
              <Car className="w-4 h-4" />
              <span className="hidden sm:inline">Мои автомобили</span>
            </TabsTrigger>
            <TabsTrigger value="service" className="gap-2 data-[state=active]:bg-slate-900 data-[state=active]:text-white rounded-lg px-4 py-2.5">
              <Wrench className="w-4 h-4" />
              <span className="hidden sm:inline">Записи на ТО</span>
            </TabsTrigger>
          </TabsList>

          <TabsContent value="configurations" className="space-y-4">
            <ConfigurationList 
              configurations={configurations} 
              isLoading={configsLoading} 
            />
          </TabsContent>

          <TabsContent value="orders" className="space-y-4">
            <OrdersList orders={orders} isLoading={ordersLoading} />
          </TabsContent>

          <TabsContent value="cars" className="space-y-4">
            <UserCarsList userCars={userCars} isLoading={carsLoading} />
          </TabsContent>

          <TabsContent value="service" className="space-y-4">
            <ServiceAppointmentsList 
              appointments={appointments} 
              isLoading={appointmentsLoading} 
            />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
