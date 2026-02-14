import React from 'react';
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

const statusConfig = {
  // Configuration statuses
  draft: { label: "Черновик", className: "bg-slate-100 text-slate-700 border-slate-200" },
  confirmed: { label: "Подтверждено", className: "bg-blue-100 text-blue-700 border-blue-200" },
  ordered: { label: "Заказано", className: "bg-indigo-100 text-indigo-700 border-indigo-200" },
  cancelled: { label: "Отменено", className: "bg-red-100 text-red-700 border-red-200" },
  purchased: { label: "Куплено", className: "bg-green-100 text-green-700 border-green-200" },
  // Order statuses
  pending: { label: "Ожидание", className: "bg-amber-100 text-amber-700 border-amber-200" },
  approved: { label: "Одобрено", className: "bg-blue-100 text-blue-700 border-blue-200" },
  paid: { label: "Оплачено", className: "bg-emerald-100 text-emerald-700 border-emerald-200" },
  completed: { label: "Завершено", className: "bg-green-100 text-green-700 border-green-200" },
  // Service statuses
  scheduled: { label: "Запланировано", className: "bg-sky-100 text-sky-700 border-sky-200" },
};

export default function StatusBadge({ status, className: extraClass }) {
  const config = statusConfig[status] || { label: status, className: "bg-gray-100 text-gray-700" };
  return (
    <Badge variant="outline" className={cn("font-medium text-xs px-2.5 py-0.5 border", config.className, extraClass)}>
      {config.label}
    </Badge>
  );
}