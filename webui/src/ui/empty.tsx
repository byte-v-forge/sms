import { ReactNode } from 'react';
import { cn } from './utils';

export function EmptyBlock({ text, className }: { text: ReactNode; className?: string }) {
  return <div className={cn('rounded-xl border border-dashed border-border p-8 text-center text-sm text-muted-foreground', className)}>{text}</div>;
}

export function DescriptionLine({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div className="flex justify-between gap-3 text-xs">
      <span className="text-muted-foreground">{label}</span>
      <span className="truncate text-right">{value}</span>
    </div>
  );
}
