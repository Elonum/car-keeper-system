import React from 'react';
import { cn } from "@/lib/utils";
import { Check } from 'lucide-react';

export default function StepIndicator({ steps, currentStep }) {
  return (
    <div className="flex items-center gap-2">
      {steps.map((step, index) => {
        const isActive = index === currentStep;
        const isCompleted = index < currentStep;
        return (
          <React.Fragment key={index}>
            <div className="flex items-center gap-2">
              <div className={cn(
                "w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold transition-all",
                isActive && "bg-slate-900 text-white ring-4 ring-slate-900/10",
                isCompleted && "bg-green-500 text-white",
                !isActive && !isCompleted && "bg-slate-100 text-slate-400"
              )}>
                {isCompleted ? <Check className="w-4 h-4" /> : index + 1}
              </div>
              <span className={cn(
                "text-sm font-medium hidden sm:block",
                isActive ? "text-slate-900" : "text-slate-400"
              )}>
                {step}
              </span>
            </div>
            {index < steps.length - 1 && (
              <div className={cn(
                "w-8 sm:w-12 h-0.5 rounded-full",
                isCompleted ? "bg-green-500" : "bg-slate-200"
              )} />
            )}
          </React.Fragment>
        );
      })}
    </div>
  );
}