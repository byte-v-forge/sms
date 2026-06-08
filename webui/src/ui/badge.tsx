import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from './utils';

const badgeVariants = cva('inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium', {
  variants: {
    variant: {
      default: 'border-primary bg-primary text-primary-foreground',
      secondary: 'border-secondary bg-secondary text-secondary-foreground',
      outline: 'border-border bg-background text-muted-foreground',
      destructive: 'border-destructive bg-destructive text-destructive-foreground'
    }
  },
  defaultVariants: { variant: 'default' }
});

export function Badge({ className, variant, ...props }: React.HTMLAttributes<HTMLSpanElement> & VariantProps<typeof badgeVariants>) {
  return <span className={cn(badgeVariants({ variant }), className)} {...props} />;
}
