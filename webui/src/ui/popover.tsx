import * as React from 'react';
import * as PopoverPrimitive from '@radix-ui/react-popover';
import { cn } from './utils';

export const Popover = PopoverPrimitive.Root;
export const PopoverTrigger = PopoverPrimitive.Trigger;

export function PopoverContent({ className, align = 'start', sideOffset = 6, ...props }: React.ComponentPropsWithoutRef<typeof PopoverPrimitive.Content>) {
  return (
    <PopoverPrimitive.Portal>
      <PopoverPrimitive.Content
        align={align}
        sideOffset={sideOffset}
        className={cn('z-50 rounded-xl border border-border bg-popover p-0 text-popover-foreground shadow-lg outline-none', className)}
        {...props}
      />
    </PopoverPrimitive.Portal>
  );
}
