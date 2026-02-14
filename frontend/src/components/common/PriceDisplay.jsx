import React from 'react';
import { cn } from "@/lib/utils";

export default function PriceDisplay({ price, size = "md", className, prefix = "" }) {
  const formatted = new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(price || 0);

  const sizes = {
    sm: "text-sm font-medium",
    md: "text-lg font-semibold",
    lg: "text-2xl font-bold",
    xl: "text-3xl font-bold",
  };

  return (
    <span className={cn(sizes[size], "text-slate-900 tabular-nums tracking-tight", className)}>
      {prefix}{formatted}
    </span>
  );
}