import React from 'react';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { orderStatusBadgeClassName } from '@/lib/orderStatusDisplay';

/**
 * @param {object} props
 * @param {string} [props.code] - technical status code (for colors)
 * @param {string} props.label - user-facing label from API / resolveOrderStatusLabel
 */
export default function OrderStatusBadge({ code, label, className: extraClass }) {
  const text = label?.trim() || (code != null && code !== '' ? String(code) : '—');
  return (
    <Badge
      variant="outline"
      className={cn(
        'font-medium text-xs px-2.5 py-0.5 border',
        orderStatusBadgeClassName(code),
        extraClass
      )}
    >
      {text}
    </Badge>
  );
}
