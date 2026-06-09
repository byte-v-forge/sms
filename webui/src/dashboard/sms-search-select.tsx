import { useState } from 'react';
import { Check, ChevronsUpDown, X } from 'lucide-react';
import { Button, Command, CommandEmpty, CommandInput, CommandItem, CommandList, Popover, PopoverContent, PopoverTrigger } from '../ui';
import { cn } from '../ui/utils';

export type SearchSelectOption = {
  value: string;
  label: string;
  description?: string;
  badge?: string;
  keywords?: string[];
};

type SearchSelectProps = {
  label: string;
  placeholder: string;
  emptyText: string;
  value: string;
  searchValue?: string;
  options: SearchSelectOption[];
  contentClassName?: string;
  shouldFilter?: boolean;
  onSearchChange?: (value: string) => void;
  onValueChange: (value: string) => void;
};

export function SearchSelect(props: SearchSelectProps) {
  const [open, setOpen] = useState(false);
  const [draftSearch, setDraftSearch] = useState('');
  const selected = props.options.find((item) => item.value === props.value);
  const searchValue = props.searchValue ?? draftSearch;
  const selectedLabel = selected?.label || props.searchValue || props.placeholder;
  function changeSearch(value: string) {
    if (props.onSearchChange) props.onSearchChange(value);
    else setDraftSearch(value);
  }

  function selectValue(value: string) {
    props.onValueChange(value === props.value ? '' : value);
    if (!props.onSearchChange) setDraftSearch('');
    setOpen(false);
  }

  function clearValue() {
    props.onValueChange('');
    changeSearch('');
  }

  return (
    <div className="grid gap-1 text-xs font-medium text-muted-foreground">
      <span>{props.label}</span>
      <div className="flex gap-1">
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button variant="outline" className="h-9 min-w-0 flex-1 justify-between bg-background px-3 font-normal text-foreground">
              <span className={cn('truncate', !selected && !props.searchValue && 'text-muted-foreground')}>{selectedLabel}</span>
              <ChevronsUpDown className="size-4 shrink-0 text-muted-foreground" />
            </Button>
          </PopoverTrigger>
          <PopoverContent className={cn('w-[min(28rem,calc(100vw-2rem))] bg-white', props.contentClassName)}>
            <Command shouldFilter={props.shouldFilter ?? true} loop className="bg-white">
              <CommandInput value={searchValue} onValueChange={changeSearch} placeholder={props.placeholder} />
              <CommandList>
                <CommandEmpty>{props.emptyText}</CommandEmpty>
                {props.options.map((item) => <SearchSelectItem key={item.value} option={item} selected={item.value === props.value} onSelect={selectValue} />)}
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>
        {(props.value || props.searchValue) && (
          <Button aria-label={`清空${props.label}`} title={`清空${props.label}`} size="icon" variant="outline" onClick={clearValue}>
            <X className="size-4" />
          </Button>
        )}
      </div>
    </div>
  );
}

function SearchSelectItem({ option, selected, onSelect }: { option: SearchSelectOption; selected: boolean; onSelect: (value: string) => void }) {
  return (
    <CommandItem value={option.value} keywords={[option.label, option.description || '', ...(option.keywords || [])]} onSelect={onSelect}>
      <Check className={cn('mr-2 size-4', selected ? 'opacity-100' : 'opacity-0')} />
      <span className="min-w-0 flex-1 truncate">{option.label}</span>
      {option.description && <span className="ml-2 shrink-0 text-xs text-muted-foreground">{option.description}</span>}
      {option.badge && <span className="ml-2 shrink-0 rounded-full bg-muted px-2 py-0.5 text-[11px] text-muted-foreground">{option.badge}</span>}
    </CommandItem>
  );
}
