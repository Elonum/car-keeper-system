import React, { useState } from 'react';
import { useAuth } from '@/lib/AuthContext';
import { useNavigate } from 'react-router-dom';
import { createPageUrl } from '../utils';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { toast } from 'sonner';
import { LogIn, Mail, Lock, Car, AlertCircle } from 'lucide-react';

const validateEmail = (email) => {
  if (!email || email.trim() === '') {
    return 'Email обязателен для заполнения';
  }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    return 'Введите корректный email адрес';
  }
  return null;
};

const validatePassword = (password) => {
  if (!password || password.trim() === '') {
    return 'Пароль обязателен для заполнения';
  }
  if (password.length < 6) {
    return 'Пароль должен содержать минимум 6 символов';
  }
  return null;
};

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState({});
  const [serverError, setServerError] = useState(null);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleEmailChange = (e) => {
    const value = e.target.value;
    setEmail(value);
    if (errors.email) {
      setErrors(prev => ({ ...prev, email: null }));
    }
    setServerError(null);
  };

  const handlePasswordChange = (e) => {
    const value = e.target.value;
    setPassword(value);
    if (errors.password) {
      setErrors(prev => ({ ...prev, password: null }));
    }
    setServerError(null);
  };

  const validateForm = () => {
    const newErrors = {};
    
    const emailError = validateEmail(email);
    if (emailError) {
      newErrors.email = emailError;
    }
    
    const passwordError = validatePassword(password);
    if (passwordError) {
      newErrors.password = passwordError;
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setServerError(null);
    
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);

    try {
      await login(email.trim(), password);
      toast.success('Вход выполнен успешно');
      navigate(createPageUrl('Catalog'));
    } catch (error) {
      let errorMessage = 'Ошибка входа. Проверьте email и пароль.';
      
      if (error.message) {
        const msg = error.message.toLowerCase();
        if (msg.includes('invalid credentials') || msg.includes('invalid') || msg.includes('incorrect') || msg.includes('wrong') || msg.includes('password') || msg.includes('credentials')) {
          errorMessage = 'Неверный email или пароль';
        } else if (msg.includes('user not found') || msg.includes('not found')) {
          errorMessage = 'Пользователь с таким email не найден';
        } else if (msg.includes('unauthorized')) {
          errorMessage = 'Неверный email или пароль';
        } else if (msg.includes('network') || msg.includes('connection')) {
          errorMessage = 'Ошибка соединения. Проверьте подключение к интернету.';
        } else {
          errorMessage = 'Ошибка входа. Проверьте email и пароль.';
        }
      }
      
      setServerError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 flex items-center justify-center px-4 py-12">
      <Card className="w-full max-w-md p-8 shadow-xl border-0">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-900 mb-4">
            <Car className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900 mb-2">Вход в CarKeeper</h1>
          <p className="text-slate-500 text-sm">Войдите в свой аккаунт</p>
        </div>

        {serverError && (
          <Alert variant="destructive" className="mb-6">
            <div className="flex items-start gap-2">
              <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
              <AlertDescription className="text-sm">{serverError}</AlertDescription>
            </div>
          </Alert>
        )}

        <form onSubmit={handleSubmit} className="space-y-5" noValidate>
          <div className="space-y-2">
            <Label htmlFor="email" className="text-sm font-medium text-slate-700">
              Email <span className="text-red-500">*</span>
            </Label>
            <div className="relative">
              <Mail className={`absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 ${errors.email ? 'text-red-500' : 'text-slate-400'}`} />
              <Input
                id="email"
                type="email"
                placeholder="email@example.com"
                value={email}
                onChange={handleEmailChange}
                onBlur={() => {
                  const error = validateEmail(email);
                  if (error) {
                    setErrors(prev => ({ ...prev, email: error }));
                  }
                }}
                required
                className={`pl-10 ${errors.email ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
                disabled={isLoading}
                aria-invalid={!!errors.email}
                aria-describedby={errors.email ? 'email-error' : undefined}
              />
            </div>
            {errors.email && (
              <p id="email-error" className="text-sm text-red-500 mt-1 flex items-center gap-1">
                <AlertCircle className="w-3.5 h-3.5" />
                {errors.email}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="password" className="text-sm font-medium text-slate-700">
              Пароль <span className="text-red-500">*</span>
            </Label>
            <div className="relative">
              <Lock className={`absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 ${errors.password ? 'text-red-500' : 'text-slate-400'}`} />
              <Input
                id="password"
                type="password"
                placeholder="Введите пароль"
                value={password}
                onChange={handlePasswordChange}
                onBlur={() => {
                  const error = validatePassword(password);
                  if (error) {
                    setErrors(prev => ({ ...prev, password: error }));
                  }
                }}
                required
                className={`pl-10 ${errors.password ? 'border-red-500 focus-visible:ring-red-500' : ''}`}
                disabled={isLoading}
                aria-invalid={!!errors.password}
                aria-describedby={errors.password ? 'password-error' : undefined}
              />
            </div>
            {errors.password && (
              <p id="password-error" className="text-sm text-red-500 mt-1 flex items-center gap-1">
                <AlertCircle className="w-3.5 h-3.5" />
                {errors.password}
              </p>
            )}
          </div>

          <Button
            type="submit"
            className="w-full bg-slate-900 hover:bg-slate-800 text-white gap-2"
            disabled={isLoading}
          >
            <LogIn className="w-4 h-4" />
            {isLoading ? 'Вход...' : 'Войти'}
          </Button>
        </form>
      </Card>
    </div>
  );
}
