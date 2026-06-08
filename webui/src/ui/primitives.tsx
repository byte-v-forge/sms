import { ButtonHTMLAttributes, HTMLAttributes, InputHTMLAttributes, ReactNode } from 'react';

type ButtonVariant = 'default' | 'outline' | 'secondary';
type ButtonSize = 'sm' | 'icon-sm' | 'icon';

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant;
  size?: ButtonSize;
};

export function Button({ className, variant = 'default', size = 'sm', ...props }: ButtonProps) {
  const base = 'inline-flex items-center justify-center rounded-lg border font-medium transition disabled:cursor-not-allowed disabled:opacity-50';
  const variants = {
    default: 'border-primary bg-primary text-primary-foreground hover:bg-primary/90',
    outline: 'border-border bg-background hover:bg-muted/60',
    secondary: 'border-secondary bg-secondary text-secondary-foreground hover:bg-secondary/80'
  } satisfies Record<ButtonVariant, string>;
  const sizes = {
    sm: 'h-8 px-3 text-xs',
    'icon-sm': 'h-8 w-8 p-0',
    icon: 'h-9 w-9 p-0'
  } satisfies Record<ButtonSize, string>;
  return <button className={cn(base, variants[variant], sizes[size], className)} type="button" {...props} />;
}

export function Badge({ className, variant = 'default', ...props }: HTMLAttributes<HTMLSpanElement> & { variant?: ButtonVariant }) {
  const variants = {
    default: 'border-primary bg-primary text-primary-foreground',
    outline: 'border-border bg-background text-muted-foreground',
    secondary: 'border-secondary bg-secondary text-secondary-foreground'
  } satisfies Record<ButtonVariant, string>;
  return <span className={cn('inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium', variants[variant], className)} {...props} />;
}

export function Input(props: InputHTMLAttributes<HTMLInputElement>) {
  return <input {...props} className={cn('h-9 rounded-lg border border-border bg-background px-3 text-sm outline-none focus:border-primary', props.className)} />;
}

export function Switch({ checked, onCheckedChange, ...props }: { checked: boolean; onCheckedChange: (checked: boolean) => void } & Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange'>) {
  return (
    <label className="relative inline-flex h-5 w-9 cursor-pointer items-center">
      <input {...props} checked={checked} className="peer sr-only" type="checkbox" onChange={(event) => onCheckedChange(event.target.checked)} />
      <span className="h-5 w-9 rounded-full bg-muted transition peer-checked:bg-primary" />
      <span className="absolute left-0.5 h-4 w-4 rounded-full bg-white shadow transition peer-checked:translate-x-4" />
    </label>
  );
}

export function Table(props: HTMLAttributes<HTMLTableElement>) {
  return <table {...props} className={cn('w-full border-collapse text-sm', props.className)} />;
}

export function TableHeader(props: HTMLAttributes<HTMLTableSectionElement>) {
  return <thead {...props} className={cn('sticky top-0 bg-card', props.className)} />;
}

export function TableBody(props: HTMLAttributes<HTMLTableSectionElement>) {
  return <tbody {...props} className={props.className} />;
}

export function TableRow(props: HTMLAttributes<HTMLTableRowElement>) {
  return <tr {...props} className={cn('border-b border-border/70', props.className)} />;
}

export function TableHead(props: HTMLAttributes<HTMLTableCellElement>) {
  return <th {...props} className={cn('px-3 py-2 text-left text-xs font-semibold text-muted-foreground', props.className)} />;
}

export function TableCell(props: HTMLAttributes<HTMLTableCellElement> & { colSpan?: number }) {
  return <td {...props} className={cn('px-3 py-2 align-top', props.className)} />;
}

export function Item({ className, variant: _variant, ...props }: HTMLAttributes<HTMLDivElement> & { variant?: 'outline' }) {
  return <div {...props} className={cn('rounded-xl border border-border p-3', className)} />;
}

export function ItemContent(props: HTMLAttributes<HTMLDivElement>) {
  return <div {...props} className={cn('flex flex-col', props.className)} />;
}

export function ItemTitle(props: HTMLAttributes<HTMLDivElement>) {
  return <div {...props} className={cn('flex items-center gap-2 text-sm font-semibold', props.className)} />;
}

export function ItemDescription(props: HTMLAttributes<HTMLDivElement>) {
  return <div {...props} className={cn('text-xs text-muted-foreground', props.className)} />;
}

export function DescriptionLine({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div className="flex justify-between gap-3 text-xs">
      <span className="text-muted-foreground">{label}</span>
      <span className="truncate text-right">{value}</span>
    </div>
  );
}

export function EmptyBlock({ text }: { text: string }) {
  return <div className="rounded-xl border border-dashed border-border p-8 text-sm text-muted-foreground">{text}</div>;
}

export function cn(...values: Array<string | false | undefined>) {
  return values.filter(Boolean).join(' ');
}
