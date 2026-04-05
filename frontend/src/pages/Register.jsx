import React, { useState, useCallback } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/lib/AuthContext';
import { createPageUrl } from '../utils';
import {
  validateEmail,
  validatePassword,
  validatePersonName,
  validatePhoneOptional,
  validatePasswordConfirm,
  normalizeEmail,
  FIELD_LIMITS,
} from '@/lib/authValidation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { toast } from 'sonner';
import { getApiErrorMessage } from '@/lib/apiErrors';
import {
  Car,
  Mail,
  Lock,
  User,
  Phone,
  UserPlus,
  AlertCircle,
  ArrowLeft,
  Eye,
  EyeOff,
} from 'lucide-react';

const initialErrors = () => ({
  first_name: null,
  last_name: null,
  email: null,
  phone: null,
  password: null,
  password_confirm: null,
});

export default function Register() {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirm, setPasswordConfirm] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showPasswordConfirm, setShowPasswordConfirm] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState(initialErrors);
  const [serverError, setServerError] = useState(null);

  const { register } = useAuth();
  const navigate = useNavigate();

  const clearFieldError = useCallback((key) => {
    setErrors((prev) => ({ ...prev, [key]: null }));
    setServerError(null);
  }, []);

  const validateForm = () => {
    const next = initialErrors();
    next.first_name = validatePersonName(firstName, 'Имя');
    next.last_name = validatePersonName(lastName, 'Фамилия');
    next.email = validateEmail(email);
    next.phone = validatePhoneOptional(phone);
    next.password = validatePassword(password);
    next.password_confirm = validatePasswordConfirm(password, passwordConfirm);
    setErrors(next);
    return Object.values(next).every((e) => e == null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setServerError(null);
    if (!validateForm()) {
      toast.error('Проверьте поля формы');
      return;
    }

    setIsLoading(true);
    try {
      const payload = {
        first_name: firstName.trim(),
        last_name: lastName.trim(),
        email: normalizeEmail(email),
        password,
      };
      const phoneTrim = phone.trim();
      if (phoneTrim) {
        payload.phone = phoneTrim;
      }

      await register(payload);
      const emailNorm = normalizeEmail(email);
      navigate(createPageUrl('Login'), {
        replace: true,
        state: { registered: true, email: emailNorm },
      });
    } catch (error) {
      const text = getApiErrorMessage(error, 'Не удалось завершить регистрацию');
      setServerError(text);
      toast.error(text);
    } finally {
      setIsLoading(false);
    }
  };

  const fieldClass = (key) =>
    errors[key] ? 'border-red-500 focus-visible:ring-red-500' : '';

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center px-4 py-10 sm:py-14">
      <div className="w-full max-w-lg">
        <div className="text-center mb-6">
          <Link
            to={createPageUrl('Login')}
            className="inline-flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors mb-6"
          >
            <ArrowLeft className="w-4 h-4" />
            Назад к входу
          </Link>
        </div>

        <Card className="border-0 shadow-2xl bg-white/95 backdrop-blur-sm overflow-hidden">
          <div className="h-1.5 bg-gradient-to-r from-slate-800 via-slate-600 to-slate-800" />
          <div className="p-8 sm:p-10">
            <div className="text-center mb-8">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-slate-900 shadow-lg mb-4">
                <Car className="w-8 h-8 text-white" />
              </div>
              <h1 className="text-2xl sm:text-3xl font-bold text-slate-900 tracking-tight">
                Создать аккаунт
              </h1>
              <p className="text-slate-500 text-sm mt-2">
                CarKeeper — конфигурации, заказы и сервис в одном месте
              </p>
            </div>

            {serverError && (
              <Alert variant="destructive" className="mb-6">
                <div className="flex items-start gap-2">
                  <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
                  <AlertDescription className="text-sm">{serverError}</AlertDescription>
                </div>
              </Alert>
            )}

            <form onSubmit={handleSubmit} className="space-y-6" noValidate>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="first_name" className="text-slate-700">
                    Имя <span className="text-red-500">*</span>
                  </Label>
                  <div className="relative">
                    <User
                      className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.first_name ? 'text-red-500' : 'text-slate-400'}`}
                    />
                    <Input
                      id="first_name"
                      name="given-name"
                      autoComplete="given-name"
                      value={firstName}
                      onChange={(e) => {
                        setFirstName(e.target.value);
                        clearFieldError('first_name');
                      }}
                      onBlur={() =>
                        setErrors((p) => ({
                          ...p,
                          first_name: validatePersonName(firstName, 'Имя'),
                        }))
                      }
                      className={`pl-10 ${fieldClass('first_name')}`}
                      disabled={isLoading}
                      aria-invalid={!!errors.first_name}
                      aria-describedby={errors.first_name ? 'err-first' : undefined}
                      maxLength={FIELD_LIMITS.NAME}
                    />
                  </div>
                  {errors.first_name && (
                    <p id="err-first" className="text-sm text-red-600 flex items-center gap-1">
                      <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                      {errors.first_name}
                    </p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="last_name" className="text-slate-700">
                    Фамилия <span className="text-red-500">*</span>
                  </Label>
                  <div className="relative">
                    <User
                      className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.last_name ? 'text-red-500' : 'text-slate-400'}`}
                    />
                    <Input
                      id="last_name"
                      name="family-name"
                      autoComplete="family-name"
                      value={lastName}
                      onChange={(e) => {
                        setLastName(e.target.value);
                        clearFieldError('last_name');
                      }}
                      onBlur={() =>
                        setErrors((p) => ({
                          ...p,
                          last_name: validatePersonName(lastName, 'Фамилия'),
                        }))
                      }
                      className={`pl-10 ${fieldClass('last_name')}`}
                      disabled={isLoading}
                      aria-invalid={!!errors.last_name}
                      maxLength={FIELD_LIMITS.NAME}
                    />
                  </div>
                  {errors.last_name && (
                    <p className="text-sm text-red-600 flex items-center gap-1">
                      <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                      {errors.last_name}
                    </p>
                  )}
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="reg-email" className="text-slate-700">
                  Email <span className="text-red-500">*</span>
                </Label>
                <div className="relative">
                  <Mail
                    className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.email ? 'text-red-500' : 'text-slate-400'}`}
                  />
                  <Input
                    id="reg-email"
                    name="email"
                    type="email"
                    inputMode="email"
                    autoComplete="email"
                    placeholder="you@example.com"
                    value={email}
                    onChange={(e) => {
                      setEmail(e.target.value);
                      clearFieldError('email');
                    }}
                    onBlur={() =>
                      setErrors((p) => ({ ...p, email: validateEmail(email) }))
                    }
                    className={`pl-10 ${fieldClass('email')}`}
                    disabled={isLoading}
                    aria-invalid={!!errors.email}
                    maxLength={FIELD_LIMITS.EMAIL}
                  />
                </div>
                {errors.email && (
                  <p className="text-sm text-red-600 flex items-center gap-1">
                    <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                    {errors.email}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone" className="text-slate-700">
                  Телефон <span className="text-slate-400 font-normal">(необязательно)</span>
                </Label>
                <div className="relative">
                  <Phone
                    className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.phone ? 'text-red-500' : 'text-slate-400'}`}
                  />
                  <Input
                    id="phone"
                    name="tel"
                    type="tel"
                    autoComplete="tel"
                    placeholder="+7 (999) 123-45-67"
                    value={phone}
                    onChange={(e) => {
                      setPhone(e.target.value);
                      clearFieldError('phone');
                    }}
                    onBlur={() =>
                      setErrors((p) => ({
                        ...p,
                        phone: validatePhoneOptional(phone),
                      }))
                    }
                    className={`pl-10 ${fieldClass('phone')}`}
                    disabled={isLoading}
                    maxLength={FIELD_LIMITS.PHONE}
                  />
                </div>
                {errors.phone && (
                  <p className="text-sm text-red-600 flex items-center gap-1">
                    <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                    {errors.phone}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="reg-password" className="text-slate-700">
                  Пароль <span className="text-red-500">*</span>
                </Label>
                <div className="relative">
                  <Lock
                    className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.password ? 'text-red-500' : 'text-slate-400'}`}
                  />
                  <Input
                    id="reg-password"
                    name="new-password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete="new-password"
                    value={password}
                    onChange={(e) => {
                      setPassword(e.target.value);
                      clearFieldError('password');
                      if (errors.password_confirm) {
                        setErrors((p) => ({ ...p, password_confirm: null }));
                      }
                    }}
                    onBlur={() =>
                      setErrors((p) => ({
                        ...p,
                        password: validatePassword(password),
                      }))
                    }
                    className={`pl-10 pr-11 ${fieldClass('password')}`}
                    disabled={isLoading}
                    aria-invalid={!!errors.password}
                    maxLength={FIELD_LIMITS.PASSWORD}
                  />
                  <button
                    type="button"
                    tabIndex={-1}
                    className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-md text-slate-400 hover:text-slate-700 hover:bg-slate-100"
                    onClick={() => setShowPassword((v) => !v)}
                    aria-label={showPassword ? 'Скрыть пароль' : 'Показать пароль'}
                  >
                    {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                {errors.password && (
                  <p className="text-sm text-red-600 flex items-center gap-1">
                    <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                    {errors.password}
                  </p>
                )}
                <p className="text-xs text-slate-500">
                  Не менее 6 символов. Рекомендуется комбинация букв и цифр.
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="password_confirm" className="text-slate-700">
                  Подтверждение пароля <span className="text-red-500">*</span>
                </Label>
                <div className="relative">
                  <Lock
                    className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${errors.password_confirm ? 'text-red-500' : 'text-slate-400'}`}
                  />
                  <Input
                    id="password_confirm"
                    name="confirm-password"
                    type={showPasswordConfirm ? 'text' : 'password'}
                    autoComplete="new-password"
                    value={passwordConfirm}
                    onChange={(e) => {
                      setPasswordConfirm(e.target.value);
                      clearFieldError('password_confirm');
                    }}
                    onBlur={() =>
                      setErrors((p) => ({
                        ...p,
                        password_confirm: validatePasswordConfirm(password, passwordConfirm),
                      }))
                    }
                    className={`pl-10 pr-11 ${fieldClass('password_confirm')}`}
                    disabled={isLoading}
                    maxLength={FIELD_LIMITS.PASSWORD}
                  />
                  <button
                    type="button"
                    tabIndex={-1}
                    className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-md text-slate-400 hover:text-slate-700 hover:bg-slate-100"
                    onClick={() => setShowPasswordConfirm((v) => !v)}
                    aria-label={
                      showPasswordConfirm ? 'Скрыть пароль' : 'Показать пароль'
                    }
                  >
                    {showPasswordConfirm ? (
                      <EyeOff className="w-4 h-4" />
                    ) : (
                      <Eye className="w-4 h-4" />
                    )}
                  </button>
                </div>
                {errors.password_confirm && (
                  <p className="text-sm text-red-600 flex items-center gap-1">
                    <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                    {errors.password_confirm}
                  </p>
                )}
              </div>

              <p className="text-xs text-slate-500 leading-relaxed">
                Регистрируясь, вы подтверждаете достоверность указанных данных. Обработка
                персональных данных осуществляется в соответствии с политикой конфиденциальности
                сервиса.
              </p>

              <Button
                type="submit"
                className="w-full h-11 bg-slate-900 hover:bg-slate-800 text-white gap-2 text-base font-medium"
                disabled={isLoading}
              >
                <UserPlus className="w-4 h-4" />
                {isLoading ? 'Создание аккаунта…' : 'Зарегистрироваться'}
              </Button>

              <div className="flex flex-wrap items-center justify-center gap-x-2 gap-y-0 pt-1 text-sm text-slate-600">
                <span className="leading-none">Уже есть аккаунт?</span>
                <Link
                  to={createPageUrl('Login')}
                  className="font-semibold text-slate-900 underline-offset-4 hover:underline leading-none"
                >
                  Войти
                </Link>
              </div>
            </form>
          </div>
        </Card>
      </div>
    </div>
  );
}
