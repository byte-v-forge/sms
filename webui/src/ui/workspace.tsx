import type { ReactNode } from 'react';
import { NavLink } from 'react-router';
import { buttonClassName } from './primitives';

type RouteTab = {
  to: string;
  label: string;
  end?: boolean;
};

export function WorkspaceRoutedPanel({ title, meta, tabs, children }: { title: ReactNode; meta?: string; tabs: RouteTab[]; children: ReactNode }) {
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
              <NavLink key={tab.to} to={tab.to} end={tab.end ?? true} className={({ isActive }) => buttonClassName(isActive ? 'default' : 'outline', 'sm')}>
                {tab.label}
              </NavLink>
            ))}
          </nav>
        </header>
        <div className="flex min-h-0 flex-1 flex-col">{children}</div>
      </section>
    </main>
  );
}
