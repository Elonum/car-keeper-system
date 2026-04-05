import React from 'react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import {
  appointmentStatusClassName,
  resolveAppointmentStatusLabel,
} from '@/lib/appointmentStatusDisplay';

/** @param {{ code: string }} props */
export default function AppointmentStatusBadge({ code, className: extraClass }) {
  const label = resolveAppointmentStatusLabel(code);
  return (
    <Badge
      variant="outline"
      className={cn(
        'font-medium text-xs px-2.5 py-0.5 border',
        appointmentStatusClassName(code),
        extraClass
      )}
    >
      {label}
    </Badge>
  );
}
