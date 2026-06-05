import React from 'react';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, ShieldAlert, FileWarning } from 'lucide-react';

const VARIANT_META = {
  server: {
    title: 'Не удалось выполнить действие',
    Icon: AlertCircle,
    className: 'border-red-200 bg-red-50 text-red-950',
    iconClassName: 'text-red-600',
  },
  form: {
    title: 'Проверьте поля формы',
    Icon: FileWarning,
    className: 'border-amber-200 bg-amber-50 text-amber-950',
    iconClassName: 'text-amber-600',
  },
  permission: {
    title: 'Недостаточно прав',
    Icon: ShieldAlert,
    className: 'border-slate-200 bg-slate-50 text-slate-900',
    iconClassName: 'text-slate-600',
  },
};

export function ErrorNotice({ kind = 'server', message, className = '' }) {
  if (!message) return null;
  const meta = VARIANT_META[kind] || VARIANT_META.server;
  const Icon = meta.Icon;

  return (
    <Alert className={`${meta.className} rounded-xl shadow-sm ${className}`.trim()}>
      <div className="flex gap-3">
        <Icon className={`h-5 w-5 shrink-0 mt-0.5 ${meta.iconClassName}`} />
        <AlertDescription className="text-sm">
          <span className="block font-semibold">{meta.title}</span>
          <span className="mt-0.5 block leading-5">{message}</span>
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

