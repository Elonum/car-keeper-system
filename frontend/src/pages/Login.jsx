import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/lib/AuthContext';
import { createPageUrl } from '../utils';
import { validateEmail, validatePassword, FIELD_LIMITS } from '@/lib/authValidation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { toast } from 'sonner';
import { LogIn, Mail, Lock, Car, AlertCircle, UserPlus, Home } from 'lucide-react';

function mapLoginError(message) {
  if (!message) return 'Ошибка входа. Проверьте email и пароль.';
  const msg = String(message).toLowerCase();
  if (msg.includes('network') || msg.includes('connection')) {
    return 'Ошибка соединения. Проверьте подключение к интернету.';
  }
  if (msg.includes('user not found') || msg.includes('not found')) {
    return 'Пользователь с таким email не найден';
  }
  if (
    msg.includes('invalid') ||
    msg.includes('wrong') ||
    msg.includes('unauthorized') ||
    msg.includes('credentials') ||
    msg.includes('password')
  ) {
    return 'Неверный email или пароль';
  }
  return 'Ошибка входа. Проверьте email и пароль.';
}

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState({});
  const [serverError, setServerError] = useState(null);
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const st = location.state;
    if (!st?.registered) return;
    if (typeof st.email === 'string' && st.email) {
      setEmail(st.email);
    }
    toast.success('Аккаунт создан. Войдите, используя email и пароль.', {
      id: 'login-after-register',
    });
    navigate(location.pathname, { replace: true, state: {} });
  }, [location.state, location.pathname, navigate]);

  const clearServerError = () => setServerError(null);

  const validateForm = () => {
    const next = {};
    const eErr = validateEmail(email);
    const pErr = validatePassword(password);
    if (eErr) next.email = eErr;
    if (pErr) next.password = pErr;
    setErrors(next);
    return Object.keys(next).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    clearServerError();
    if (!validateForm()) return;

    setIsLoading(true);
    try {
      await login(email.trim(), password);
      toast.success('Вход выполнен успешно');
      navigate(createPageUrl('Catalog'));
    } catch (error) {
      const errorMessage = mapLoginError(error?.message);
      setServerError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col items-center justify-center px-4 py-12">
      <Card className="w-full max-w-md border-0 shadow-2xl bg-white/95 backdrop-blur-sm overflow-hidden">
        <div className="h-1.5 bg-gradient-to-r from-slate-800 via-slate-600 to-slate-800" />
        <div className="p-8">
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-slate-900 shadow-lg mb-4">
              <Car className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-slate-900 mb-2">Вход в CarKeeper</h1>
            <p className="text-slate-500 text-sm">Войдите в свой аккаунт</p>
          </div>

          {serverError && (
            <Alert variant="destructive" className="mb-6">
              <div className="flex items-start gap-2">
                <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
                <AlertDescription className="text-sm">{serverError}</AlertDescription>
              </div>
            </Alert>
          )}

          <form onSubmit={handleSubmit} className="space-y-5" noValidate>
            <div className="space-y-2">
              <Label htmlFor="login-email" className="text-sm font-medium text-slate-700">
                Email <span className="text-red-500">*</span>
              </Label>
              <div className="relative">
                <Mail
                  className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.email ? 'text-red-500' : 'text-slate-400'}`}
                />
                <Input
                  id="login-email"
                  type="email"
                  inputMode="email"
                  autoComplete="email"
                  placeholder="email@example.com"
                  value={email}
                  maxLength={FIELD_LIMITS.EMAIL}
                  onChange={(ev) => {
                    setEmail(ev.target.value);
                    if (errors.email) setErrors((prev) => ({ ...prev, email: null }));
                    clearServerError();
                  }}
                  onBlur={() => {
                    const err = validateEmail(email);
                    if (err) setErrors((prev) => ({ ...prev, email: err }));
                  }}
                  className={`pl-10 ${errors.email ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
                  disabled={isLoading}
                  aria-invalid={!!errors.email}
                  aria-describedby={errors.email ? 'login-email-err' : undefined}
                />
              </div>
              {errors.email && (
                <p id="login-email-err" className="text-sm text-red-500 flex items-center gap-1">
                  <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                  {errors.email}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="login-password" className="text-sm font-medium text-slate-700">
                Пароль <span className="text-red-500">*</span>
              </Label>
              <div className="relative">
                <Lock
                  className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.password ? 'text-red-500' : 'text-slate-400'}`}
                />
                <Input
                  id="login-password"
                  type="password"
                  autoComplete="current-password"
                  placeholder="Введите пароль"
                  value={password}
                  maxLength={FIELD_LIMITS.PASSWORD}
                  onChange={(ev) => {
                    setPassword(ev.target.value);
                    if (errors.password) setErrors((prev) => ({ ...prev, password: null }));
                    clearServerError();
                  }}
                  onBlur={() => {
                    const err = validatePassword(password);
                    if (err) setErrors((prev) => ({ ...prev, password: err }));
                  }}
                  className={`pl-10 ${errors.password ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
                  disabled={isLoading}
                  aria-invalid={!!errors.password}
                  aria-describedby={errors.password ? 'login-password-err' : undefined}
                />
              </div>
              {errors.password && (
                <p id="login-password-err" className="text-sm text-red-500 flex items-center gap-1">
                  <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                  {errors.password}
                </p>
              )}
            </div>

            <Button
              type="submit"
              className="w-full h-11 bg-slate-900 hover:bg-slate-800 text-white gap-2"
              disabled={isLoading}
            >
              <LogIn className="w-4 h-4 shrink-0" />
              {isLoading ? 'Вход...' : 'Войти'}
            </Button>

            <div className="flex flex-wrap items-center justify-center gap-x-2 gap-y-0 pt-2 text-sm text-slate-600">
              <span className="leading-none">Нет аккаунта?</span>
              <Link
                to={createPageUrl('Register')}
                className="inline-flex items-center gap-1.5 font-semibold text-slate-900 underline-offset-4 hover:underline leading-none"
              >
                <UserPlus className="w-4 h-4 shrink-0" aria-hidden />
                Зарегистрироваться
              </Link>
            </div>
          </form>
        </div>
      </Card>

      <Link
        to={createPageUrl('Catalog')}
        className="mt-8 inline-flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 rounded-md px-1 py-1"
      >
        <Home className="w-4 h-4 shrink-0" aria-hidden />
        На главную без входа
      </Link>
    </div>
  );
}
