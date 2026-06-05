import React from 'react';

export default function SectionHeader({ title, description, action }) {
  return (
    <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between mb-6 sm:mb-8">
      <div className="min-w-0 flex-1">
        <h1 className="text-xl sm:text-2xl md:text-3xl font-bold text-slate-900 tracking-tight break-words">
          {title}
        </h1>
        {description && (
          <p className="text-slate-500 mt-1.5 text-sm md:text-base break-words">{description}</p>
        )}
      </div>
      {action && <div className="w-full sm:w-auto sm:ml-4 shrink-0">{action}</div>}
    </div>
  );
}