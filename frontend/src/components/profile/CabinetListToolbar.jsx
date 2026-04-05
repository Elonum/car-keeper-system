import React from 'react';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Search } from 'lucide-react';

/**
 * @param {object} props
 * @param {string} props.search
 * @param {(v: string) => void} props.onSearchChange
 * @param {string} props.searchPlaceholder
 * @param {string} [props.statusValue] — 'all' или код
 * @param {(v: string) => void} [props.onStatusChange]
 * @param {Array<{ value: string, label: string }>} [props.statusOptions] — первая обычно { value: 'all', label: 'Все статусы' }
 */
export default function CabinetListToolbar({
  search,
  onSearchChange,
  searchPlaceholder,
  statusValue,
  onStatusChange,
  statusOptions,
}) {
  const showStatus = Array.isArray(statusOptions) && statusOptions.length > 0 && onStatusChange;

  return (
    <div
      className={`flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between ${
        showStatus ? 'sm:gap-4' : ''
      }`}
    >
      <div className="relative flex-1 max-w-md">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
        <Input
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder={searchPlaceholder}
          className="pl-9 bg-white"
          aria-label="Поиск"
        />
      </div>
      {showStatus ? (
        <Select value={statusValue ?? 'all'} onValueChange={onStatusChange}>
          <SelectTrigger className="w-full sm:w-[220px] bg-white" aria-label="Фильтр по статусу">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {statusOptions.map((opt) => (
              <SelectItem key={opt.value} value={opt.value}>
                {opt.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ) : null}
    </div>
  );
}
