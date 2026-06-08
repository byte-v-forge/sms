import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from './utils';

const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg border text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        default: 'border-primary bg-primary text-primary-foreground hover:bg-primary/90',
        destructive: 'border-destructive bg-destructive text-destructive-foreground hover:bg-destructive/90',
        outline: 'border-border bg-background hover:bg-muted/60',
        secondary: 'border-secondary bg-secondary text-secondary-foreground hover:bg-secondary/80',
        ghost: 'border-transparent hover:bg-muted/60'
      },
      size: {
        sm: 'h-8 px-3 text-xs',
        default: 'h-9 px-4',
        icon: 'h-9 w-9 p-0',
        'icon-sm': 'h-8 w-8 p-0'
      }
    },
    defaultVariants: { variant: 'default', size: 'default' }
  }
);

export type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & VariantProps<typeof buttonVariants> & { asChild?: boolean };

export function Button({ className, variant, size, asChild = false, ...props }: ButtonProps) {
  const Comp = asChild ? Slot : 'button';
  return <Comp className={cn(buttonVariants({ variant, size }), className)} {...props} />;
}
