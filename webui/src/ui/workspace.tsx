import type { ReactNode } from 'react';
import { NavLink, useLocation } from 'react-router';
import { Button } from './button';
import { Card } from './card';
import { cn } from './utils';

type RouteTab = {
  to: string;
  label: string;
  end?: boolean;
};

export function WorkspaceRoutedPanel({ title, meta, tabs, children }: { title: ReactNode; meta?: string; tabs: RouteTab[]; children: ReactNode }) {
  const location = useLocation();
  return (
    <main className="mx-auto flex h-screen max-w-7xl flex-col p-4">
      <Card className="flex min-h-0 flex-1 flex-col overflow-hidden bg-card/95 shadow-lg shadow-slate-200/60">
        <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border bg-background/70 px-4 py-3 backdrop-blur">
          <div>
            <h1 className="text-lg font-semibold tracking-tight">{title}</h1>
            {meta && <p className="text-xs text-muted-foreground">{meta}</p>}
          </div>
          <nav aria-label="SMS 页面导航" className="flex gap-2 rounded-xl bg-muted/40 p-1">
            {tabs.map((tab) => (
              <Button key={tab.to} asChild size="sm" variant={routeMatches(location.pathname, tab) ? 'default' : 'ghost'}>
                <NavLink to={tab.to} className={({ isActive }) => cn(!isActive && 'text-foreground')}>{tab.label}</NavLink>
              </Button>
            ))}
          </nav>
        </header>
        <div className="flex min-h-0 flex-1 flex-col">{children}</div>
      </Card>
    </main>
  );
}

function routeMatches(pathname: string, tab: RouteTab) {
  return tab.end ?? true ? pathname === tab.to : pathname.startsWith(tab.to);
}
