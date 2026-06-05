import React from 'react';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { AlertTriangle, Info, Trash2 } from 'lucide-react';

const VARIANT_META = {
  destructive: {
    Icon: Trash2,
    iconClassName: 'bg-red-50 text-red-600',
    actionClassName: 'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-600',
  },
  warning: {
    Icon: AlertTriangle,
    iconClassName: 'bg-amber-50 text-amber-600',
    actionClassName: 'bg-slate-900 text-white hover:bg-slate-800',
  },
  default: {
    Icon: Info,
    iconClassName: 'bg-slate-100 text-slate-700',
    actionClassName: 'bg-slate-900 text-white hover:bg-slate-800',
  },
};

export default function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel = 'Подтвердить',
  cancelLabel = 'Отмена',
  variant = 'default',
  onConfirm,
  disabled = false,
}) {
  const meta = VARIANT_META[variant] || VARIANT_META.default;
  const Icon = meta.Icon;

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="max-w-md rounded-2xl border-slate-200 p-0 shadow-2xl">
        <div className="p-6">
          <AlertDialogHeader className="space-y-3">
            <div className={`flex h-11 w-11 items-center justify-center rounded-xl ${meta.iconClassName}`}>
              <Icon className="h-5 w-5" />
            </div>
            <AlertDialogTitle className="text-xl text-slate-950">{title}</AlertDialogTitle>
            <AlertDialogDescription className="text-sm leading-6 text-slate-600">
              {description}
            </AlertDialogDescription>
          </AlertDialogHeader>
        </div>
        <AlertDialogFooter className="border-t border-slate-100 bg-slate-50/80 px-6 py-4">
          <AlertDialogCancel disabled={disabled} className="mt-0">
            {cancelLabel}
          </AlertDialogCancel>
          <AlertDialogAction
            disabled={disabled}
            className={meta.actionClassName}
            onClick={(event) => {
              event.preventDefault();
              onConfirm?.();
            }}
          >
            {confirmLabel}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
