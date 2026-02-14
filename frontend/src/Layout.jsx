import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { createPageUrl } from './utils';
import { useAuth } from '@/lib/AuthContext';
import { 
  Car, Menu, X, User, LogOut, ChevronDown,
  Newspaper, Wrench, Settings, FileText, ShoppingCart
} from 'lucide-react';
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const navItems = [
  { label: "Каталог", page: "Catalog", icon: Car },
  { label: "Новости", page: "News", icon: Newspaper },
  { label: "Сервис", page: "Services", icon: Wrench },
];

export default function Layout({ children, currentPageName }) {
  const { user, logout, isAuthenticated } = useAuth();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    setMobileMenuOpen(false);
  }, [currentPageName]);

  const hideLayout = currentPageName === "Configurator";

  const handleLogout = () => {
    logout();
  };

  const handleLogin = () => {
    window.location.href = '/login';
  };

  const getUserDisplayName = () => {
    if (!user) return 'U';
    if (user.full_name) return user.full_name[0].toUpperCase();
    if (user.first_name && user.last_name) {
      return `${user.first_name[0]}${user.last_name[0]}`.toUpperCase();
    }
    if (user.email) return user.email[0].toUpperCase();
    return 'U';
  };

  const getUserName = () => {
    if (!user) return '';
    if (user.full_name) return user.full_name;
    if (user.first_name && user.last_name) {
      return `${user.first_name} ${user.last_name}`;
    }
    if (user.email) return user.email;
    return 'Пользователь';
  };

  return (
    <div className="min-h-screen bg-slate-50 font-sans">
      <style>{`
        :root {
          --color-primary: #0f172a;
          --color-accent: #2563eb;
          --color-accent-light: #3b82f6;
        }
        body {
          font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
          -webkit-font-smoothing: antialiased;
        }
        * { scrollbar-width: thin; scrollbar-color: #cbd5e1 transparent; }
      `}</style>

      {/* Header */}
      <header className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        scrolled ? 'bg-white/95 backdrop-blur-md shadow-sm border-b border-slate-100' : 'bg-white border-b border-slate-100'
      }`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link to={createPageUrl("Catalog")} className="flex items-center gap-2.5 group">
              <div className="w-9 h-9 rounded-xl bg-slate-900 flex items-center justify-center group-hover:bg-blue-600 transition-colors">
                <Car className="w-5 h-5 text-white" />
              </div>
              <span className="text-lg font-bold text-slate-900 tracking-tight hidden sm:block">
                CarKeeper
              </span>
            </Link>

            {/* Desktop nav */}
            <nav className="hidden md:flex items-center gap-1">
              {navItems.map((item) => (
                <Link
                  key={item.page}
                  to={createPageUrl(item.page)}
                  className={`px-3.5 py-2 rounded-lg text-sm font-medium transition-all ${
                    currentPageName === item.page
                      ? 'bg-slate-900 text-white'
                      : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'
                  }`}
                >
                  {item.label}
                </Link>
              ))}
            </nav>

            {/* Right side */}
            <div className="flex items-center gap-2">
              {isAuthenticated && user ? (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="flex items-center gap-2 px-3 h-9">
                      <div className="w-7 h-7 rounded-full bg-slate-900 flex items-center justify-center">
                        <span className="text-xs font-semibold text-white">
                          {getUserDisplayName()}
                        </span>
                      </div>
                      <span className="text-sm font-medium text-slate-700 hidden sm:block max-w-[120px] truncate">
                        {getUserName()}
                      </span>
                      <ChevronDown className="w-3.5 h-3.5 text-slate-400" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56">
                    <div className="px-3 py-2.5 border-b">
                      <p className="text-sm font-semibold text-slate-900">{getUserName()}</p>
                      <p className="text-xs text-slate-500 truncate">{user.email || ''}</p>
                    </div>
                    <DropdownMenuItem asChild>
                      <Link to={createPageUrl("Profile")} className="flex items-center gap-2 cursor-pointer">
                        <User className="w-4 h-4" /> Личный кабинет
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <Link to={createPageUrl("Profile") + "?tab=orders"} className="flex items-center gap-2 cursor-pointer">
                        <ShoppingCart className="w-4 h-4" /> Мои заказы
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <Link to={createPageUrl("Profile") + "?tab=service"} className="flex items-center gap-2 cursor-pointer">
                        <Wrench className="w-4 h-4" /> Записи на ТО
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onClick={handleLogout} className="text-red-600 cursor-pointer">
                      <LogOut className="w-4 h-4 mr-2" /> Выйти
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : (
                <Button 
                  onClick={handleLogin} 
                  className="bg-slate-900 hover:bg-slate-800 text-white h-9 px-4 text-sm"
                >
                  Войти
                </Button>
              )}

              {/* Mobile menu button */}
              <Button
                variant="ghost"
                size="icon"
                className="md:hidden h-9 w-9"
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              >
                {mobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
              </Button>
            </div>
          </div>
        </div>

        {/* Mobile nav */}
        {mobileMenuOpen && (
          <div className="md:hidden border-t bg-white px-4 py-3 space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.page}
                to={createPageUrl(item.page)}
                className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all ${
                  currentPageName === item.page
                    ? 'bg-slate-900 text-white'
                    : 'text-slate-600 hover:bg-slate-100'
                }`}
              >
                <item.icon className="w-4 h-4" />
                {item.label}
              </Link>
            ))}
          </div>
        )}
      </header>

      {/* Content */}
      <main className={hideLayout ? "" : "pt-16"}>
        {children}
      </main>

      {/* Footer */}
      {!hideLayout && (
        <footer className="bg-slate-900 text-white mt-auto">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div>
                <div className="flex items-center gap-2.5 mb-4">
                  <div className="w-9 h-9 rounded-xl bg-white/10 flex items-center justify-center">
                    <Car className="w-5 h-5 text-white" />
                  </div>
                  <span className="text-lg font-bold tracking-tight">CarKeeper</span>
                </div>
                <p className="text-slate-400 text-sm leading-relaxed">
                  Ваш надёжный партнёр в мире автомобилей. Конфигурация, покупка и сервисное обслуживание в одном месте.
                </p>
              </div>
              <div>
                <h4 className="font-semibold text-sm mb-3 text-slate-300 uppercase tracking-wider">Навигация</h4>
                <div className="space-y-2">
                  {navItems.map(item => (
                    <Link key={item.page} to={createPageUrl(item.page)} 
                      className="block text-sm text-slate-400 hover:text-white transition-colors">
                      {item.label}
                    </Link>
                  ))}
                </div>
              </div>
              <div>
                <h4 className="font-semibold text-sm mb-3 text-slate-300 uppercase tracking-wider">Контакты</h4>
                <div className="space-y-2 text-sm text-slate-400">
                  <p>+7 (800) 123-45-67</p>
                  <p>info@carkeeper.ru</p>
                  <p>Москва, ул. Автомобильная, 1</p>
                </div>
              </div>
            </div>
            <div className="border-t border-slate-800 mt-10 pt-6 text-center">
              <p className="text-xs text-slate-500">© 2026 CarKeeper. Все права защищены.</p>
            </div>
          </div>
        </footer>
      )}
    </div>
  );
}
