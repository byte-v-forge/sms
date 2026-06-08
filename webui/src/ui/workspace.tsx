import { ReactNode, useState } from 'react';
import { Button } from './primitives';

type Tab = {
  value: string;
  label: string;
  content: ReactNode;
};

export function WorkspaceTabbedPanel({ defaultValue, title, meta, tabs }: { defaultValue: string; title: ReactNode; meta?: string; tabs: Tab[] }) {
  const [active, setActive] = useState(defaultValue);
  const current = tabs.find((tab) => tab.value === active) || tabs[0];
  return (
    <main className="mx-auto flex h-screen max-w-7xl flex-col p-4">
      <section className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-2xl border border-border bg-card shadow-sm">
        <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-4 py-3">
          <div>
            <h1 className="text-lg font-semibold">{title}</h1>
            {meta && <p className="text-xs text-muted-foreground">{meta}</p>}
          </div>
          <nav className="flex gap-2">
            {tabs.map((tab) => (
              <Button key={tab.value} size="sm" variant={tab.value === current.value ? 'default' : 'outline'} onClick={() => setActive(tab.value)}>
                {tab.label}
              </Button>
            ))}
          </nav>
        </header>
        <div className="flex min-h-0 flex-1 flex-col">{current.content}</div>
      </section>
    </main>
  );
}
