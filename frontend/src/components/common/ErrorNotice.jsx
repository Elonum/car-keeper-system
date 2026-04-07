import React from 'react';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, ShieldAlert, FileWarning } from 'lucide-react';

const VARIANT_META = {
  server: {
    title: 'Ошибка сервера',
    Icon: AlertCircle,
  },
  form: {
    title: 'Проверьте поля формы',
    Icon: FileWarning,
  },
  permission: {
    title: 'Недостаточно прав',
    Icon: ShieldAlert,
  },
};

export function ErrorNotice({ kind = 'server', message, className = '' }) {
  if (!message) return null;
  const meta = VARIANT_META[kind] || VARIANT_META.server;
  const Icon = meta.Icon;

  return (
    <Alert variant="destructive" className={`border-red-200 bg-red-50/90 ${className}`.trim()}>
      <div className="flex gap-2">
        <Icon className="h-4 w-4 shrink-0 mt-0.5" />
        <AlertDescription className="text-sm text-red-900">
          <span className="font-medium">{meta.title}: </span>
          {message}
        </AlertDescription>
      </div>
    </Alert>
  );
}

export function FieldErrorText({ children, id }) {
  if (!children) return null;
  return (
    <p id={id} className="text-sm text-red-600 flex items-center gap-1">
      <AlertCircle className="w-3.5 h-3.5 shrink-0" />
      {children}
    </p>
  );
}

