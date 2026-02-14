import React from 'react';
import { Loader2 } from 'lucide-react';

export default function PageLoader() {
  return (
    <div className="flex items-center justify-center min-h-[60vh]">
      <div className="flex flex-col items-center gap-3">
        <Loader2 className="w-8 h-8 animate-spin text-slate-400" />
        <p className="text-sm text-slate-500">Загрузка...</p>
      </div>
    </div>
  );
}