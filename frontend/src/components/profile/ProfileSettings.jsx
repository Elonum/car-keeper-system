import React, { useEffect, useState, useCallback } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ErrorNotice, FieldErrorText } from '@/components/common/ErrorNotice';
import { useAuth } from '@/lib/AuthContext';
import { profileService } from '@/services/profileService';
import {
  validatePersonName,
  validatePhoneOptional,
  validatePassword,
  validatePasswordConfirm,
  FIELD_LIMITS,
} from '@/lib/authValidation';
import { getApiErrorMessage } from '@/lib/apiErrors';
import {
  User,
  Mail,
  Phone,
  Shield,
  Save,
  KeyRound,
  Eye,
  EyeOff,
} from 'lucide-react';

const ROLE_LABELS = {
  customer: 'Клиент',
  manager: 'Менеджер',
  service_advisor: 'Мастер-приёмщик',
  admin: 'Администратор',
};

export default function ProfileSettings() {
  const { user, refreshUser } = useAuth();

  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [phone, setPhone] = useState('');
  const [profileErrors, setProfileErrors] = useState({});

  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [pwdErrors, setPwdErrors] = useState({});
  /** Server-side error (API), shown inline in-card */
  const [profileFormError, setProfileFormError] = useState(null);
  const [passwordFormError, setPasswordFormError] = useState(null);
  const [showCurrent, setShowCurrent] = useState(false);
  const [showNew, setShowNew] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const syncFromUser = useCallback(() => {
    if (!user) return;
    setFirstName(user.first_name ?? '');
    setLastName(user.last_name ?? '');
    setPhone(user.phone ?? '');
  }, [user]);

  useEffect(() => {
    syncFromUser();
  }, [syncFromUser]);

  const profileMutation = useMutation({
    mutationFn: (payload) => profileService.updateProfile(payload),
    onMutate: () => {
      setProfileFormError(null);
    },
    onSuccess: async () => {
      try {
        await refreshUser();
        setProfileErrors({});
        setProfileFormError(null);
      } catch {
        const msg =
          'Данные на сервере обновлены. Обновите страницу, если изменения не отобразились.';
        setProfileFormError(msg);
      }
    },
    onError: (e) => {
      const msg = getApiErrorMessage(e, 'Не удалось сохранить данные профиля');
      setProfileFormError(msg);
    },
  });

  const passwordMutation = useMutation({
    mutationFn: (payload) => profileService.changePassword(payload),
    onMutate: () => {
      setPasswordFormError(null);
    },
    onSuccess: async () => {
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
      setPwdErrors({});
      setPasswordFormError(null);
    },
    onError: (e) => {
      const msg = getApiErrorMessage(e, 'Не удалось сменить пароль');
      setPasswordFormError(msg);
    },
  });

  const validateProfileForm = () => {
    const next = {};
    const fn = validatePersonName(firstName, 'Имя');
    const ln = validatePersonName(lastName, 'Фамилия');
    const ph = validatePhoneOptional(phone);
    if (fn) next.first_name = fn;
    if (ln) next.last_name = ln;
    if (ph) next.phone = ph;
    setProfileErrors(next);
    return Object.keys(next).length === 0;
  };

  const handleProfileSubmit = (e) => {
    e.preventDefault();
    if (!validateProfileForm()) {
      setProfileFormError('Проверьте поля');
      return;
    }
    const payload = {
      first_name: firstName.trim(),
      last_name: lastName.trim(),
    };
    const phoneTrim = phone.trim();
    payload.phone = phoneTrim || null;
    profileMutation.mutate(payload);
  };

  const validatePasswordForm = () => {
    const next = {};
    if (!currentPassword) {
      next.current = 'Введите текущий пароль';
    }
    const np = validatePassword(newPassword);
    if (np) next.new = np;
    const cf = validatePasswordConfirm(newPassword, confirmPassword);
    if (cf) next.confirm = cf;
    setPwdErrors(next);
    return Object.keys(next).length === 0;
  };

  const handlePasswordSubmit = (e) => {
    e.preventDefault();
    if (!validatePasswordForm()) {
      setPasswordFormError('Проверьте поля');
      return;
    }
    passwordMutation.mutate({
      current_password: currentPassword,
      new_password: newPassword,
    });
  };

  const roleLabel = ROLE_LABELS[user?.role] || user?.role || '—';

  return (
    <div className="space-y-8 max-w-2xl">
      <Card className="p-6 sm:p-8 border-0 shadow-sm">
        <div className="flex items-center gap-3 mb-6">
          <div className="w-10 h-10 rounded-xl bg-slate-900 flex items-center justify-center">
            <User className="w-5 h-5 text-white" />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-slate-900">Личные данные</h2>
            <p className="text-sm text-slate-500">
              Имя и телефон используются в заказах и записях на сервис
            </p>
          </div>
        </div>

        <form onSubmit={handleProfileSubmit} className="space-y-5" noValidate>
          <ErrorNotice kind="server" message={profileFormError} />
          <div className="grid sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="prof-fn">Имя</Label>
              <Input
                id="prof-fn"
                value={firstName}
                maxLength={FIELD_LIMITS.NAME}
                onChange={(e) => {
                  setFirstName(e.target.value);
                  setProfileFormError(null);
                  if (profileErrors.first_name) {
                    setProfileErrors((p) => ({ ...p, first_name: undefined }));
                  }
                }}
                autoComplete="given-name"
                disabled={profileMutation.isPending}
                className={profileErrors.first_name ? 'border-red-500' : ''}
              />
              <FieldErrorText>{profileErrors.first_name}</FieldErrorText>
            </div>
            <div className="space-y-2">
              <Label htmlFor="prof-ln">Фамилия</Label>
              <Input
                id="prof-ln"
                value={lastName}
                maxLength={FIELD_LIMITS.NAME}
                onChange={(e) => {
                  setLastName(e.target.value);
                  setProfileFormError(null);
                  if (profileErrors.last_name) {
                    setProfileErrors((p) => ({ ...p, last_name: undefined }));
                  }
                }}
                autoComplete="family-name"
                disabled={profileMutation.isPending}
                className={profileErrors.last_name ? 'border-red-500' : ''}
              />
              <FieldErrorText>{profileErrors.last_name}</FieldErrorText>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="prof-email" className="flex items-center gap-2">
              <Mail className="w-4 h-4 text-slate-400" />
              Email
            </Label>
            <Input
              id="prof-email"
              type="email"
              value={user?.email ?? ''}
              disabled
              className="bg-slate-50 text-slate-600"
              title="Изменение email недоступно в этом интерфейсе"
            />
            <p className="text-xs text-slate-500">
              Email используется для входа. Обратитесь в поддержку за помощью.
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="prof-phone" className="flex items-center gap-2">
              <Phone className="w-4 h-4 text-slate-400" />
              Телефон <span className="text-slate-400 font-normal">(необязательно)</span>
            </Label>
            <Input
              id="prof-phone"
              type="tel"
              inputMode="tel"
              value={phone}
              maxLength={FIELD_LIMITS.PHONE}
              onChange={(e) => {
                setPhone(e.target.value);
                setProfileFormError(null);
                if (profileErrors.phone) {
                  setProfileErrors((p) => ({ ...p, phone: undefined }));
                }
              }}
              autoComplete="tel"
              disabled={profileMutation.isPending}
              placeholder="+7 (999) 123-45-67"
              className={profileErrors.phone ? 'border-red-500' : ''}
            />
            <FieldErrorText>{profileErrors.phone}</FieldErrorText>
          </div>

          <div className="flex flex-wrap items-center gap-3 pt-2">
            <Button
              type="submit"
              disabled={profileMutation.isPending}
              className="gap-2 bg-slate-900 hover:bg-slate-800"
            >
              <Save className="w-4 h-4" />
              {profileMutation.isPending ? 'Сохранение…' : 'Сохранить изменения'}
            </Button>
            <Button
              type="button"
              variant="outline"
              disabled={profileMutation.isPending}
              onClick={() => {
                syncFromUser();
                setProfileErrors({});
                setProfileFormError(null);
              }}
            >
              Сбросить
            </Button>
          </div>
        </form>

        <div className="flex items-start gap-3 pt-6 mt-2 border-t border-slate-100">
          <div className="w-9 h-9 rounded-lg bg-slate-100 flex items-center justify-center shrink-0">
            <Shield className="w-4 h-4 text-slate-700" />
          </div>
          <div>
            <p className="text-sm font-medium text-slate-900">Роль в системе</p>
            <p className="text-sm text-slate-600">{roleLabel}</p>
          </div>
        </div>
      </Card>

      <Card className="p-6 sm:p-8 border-0 shadow-sm">
        <div className="flex items-center gap-3 mb-2">
          <div className="w-10 h-10 rounded-xl bg-amber-50 flex items-center justify-center">
            <KeyRound className="w-5 h-5 text-amber-800" />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-slate-900">Смена пароля</h2>
            <p className="text-sm text-slate-500">
              Укажите текущий пароль — так мы убедимся, что меняете его вы
            </p>
          </div>
        </div>

        <form onSubmit={handlePasswordSubmit} className="space-y-5 mt-6" noValidate>
          <ErrorNotice kind="server" message={passwordFormError} />
          <div className="space-y-2">
            <Label htmlFor="pwd-current">Текущий пароль</Label>
            <div className="relative">
              <Input
                id="pwd-current"
                type={showCurrent ? 'text' : 'password'}
                autoComplete="current-password"
                value={currentPassword}
                maxLength={FIELD_LIMITS.PASSWORD}
                onChange={(e) => {
                  setCurrentPassword(e.target.value);
                  setPasswordFormError(null);
                  if (pwdErrors.current) {
                    setPwdErrors((p) => ({ ...p, current: undefined }));
                  }
                }}
                disabled={passwordMutation.isPending}
                className={`pr-11 ${pwdErrors.current ? 'border-red-500' : ''}`}
              />
              <button
                type="button"
                tabIndex={-1}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-md text-slate-400 hover:text-slate-700"
                onClick={() => setShowCurrent((v) => !v)}
                aria-label={showCurrent ? 'Скрыть пароль' : 'Показать пароль'}
              >
                {showCurrent ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            </div>
            <FieldErrorText>{pwdErrors.current}</FieldErrorText>
          </div>

          <div className="space-y-2">
            <Label htmlFor="pwd-new">Новый пароль</Label>
            <div className="relative">
              <Input
                id="pwd-new"
                type={showNew ? 'text' : 'password'}
                autoComplete="new-password"
                value={newPassword}
                maxLength={FIELD_LIMITS.PASSWORD}
                onChange={(e) => {
                  setNewPassword(e.target.value);
                  setPasswordFormError(null);
                  if (pwdErrors.new) setPwdErrors((p) => ({ ...p, new: undefined }));
                }}
                disabled={passwordMutation.isPending}
                className={`pr-11 ${pwdErrors.new ? 'border-red-500' : ''}`}
              />
              <button
                type="button"
                tabIndex={-1}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-md text-slate-400 hover:text-slate-700"
                onClick={() => setShowNew((v) => !v)}
                aria-label={showNew ? 'Скрыть пароль' : 'Показать пароль'}
              >
                {showNew ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            </div>
            <FieldErrorText>{pwdErrors.new}</FieldErrorText>
            <p className="text-xs text-slate-500">Не менее 6 символов, до 128 символов</p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="pwd-confirm">Подтверждение нового пароля</Label>
            <div className="relative">
              <Input
                id="pwd-confirm"
                type={showConfirm ? 'text' : 'password'}
                autoComplete="new-password"
                value={confirmPassword}
                maxLength={FIELD_LIMITS.PASSWORD}
                onChange={(e) => {
                  setConfirmPassword(e.target.value);
                  setPasswordFormError(null);
                  if (pwdErrors.confirm) {
                    setPwdErrors((p) => ({ ...p, confirm: undefined }));
                  }
                }}
                disabled={passwordMutation.isPending}
                className={`pr-11 ${pwdErrors.confirm ? 'border-red-500' : ''}`}
              />
              <button
                type="button"
                tabIndex={-1}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-md text-slate-400 hover:text-slate-700"
                onClick={() => setShowConfirm((v) => !v)}
                aria-label={showConfirm ? 'Скрыть пароль' : 'Показать пароль'}
              >
                {showConfirm ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            </div>
            <FieldErrorText>{pwdErrors.confirm}</FieldErrorText>
          </div>

          <Button
            type="submit"
            disabled={passwordMutation.isPending}
            variant="outline"
            className="border-slate-300"
          >
            {passwordMutation.isPending ? 'Сохранение…' : 'Обновить пароль'}
          </Button>
        </form>
      </Card>
    </div>
  );
}
