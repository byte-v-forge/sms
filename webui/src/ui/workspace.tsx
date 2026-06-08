import type { ReactNode } from 'react';
import { Card, Flex, Menu, Typography } from 'antd';
import type { MenuProps } from 'antd';
import { Link, useLocation } from 'react-router';

type RouteTab = {
  to: string;
  label: string;
  end?: boolean;
};

export function WorkspaceRoutedPanel({ title, meta, tabs, children }: { title: ReactNode; meta?: string; tabs: RouteTab[]; children: ReactNode }) {
  const location = useLocation();
  const selected = tabs.find((tab) => routeMatches(location.pathname, tab))?.to || tabs[0]?.to;
  const items: MenuProps['items'] = tabs.map((tab) => ({ key: tab.to, label: <Link to={tab.to}>{tab.label}</Link> }));
  return (
    <main className="sms-shell">
      <Card className="sms-shell-card" variant="borderless">
        <Flex align="center" justify="space-between" gap={16} wrap="wrap" style={{ padding: '16px 20px', borderBottom: '1px solid #edf0f5' }}>
          <div>
            <Typography.Title level={4} style={{ margin: 0 }}>{title}</Typography.Title>
            {meta && <Typography.Text type="secondary">{meta}</Typography.Text>}
          </div>
          <Menu className="sms-route-tabs" mode="horizontal" selectedKeys={selected ? [selected] : []} items={items} style={{ minWidth: 280, borderBottom: 0 }} />
        </Flex>
        <div className="sms-route-body">{children}</div>
      </Card>
    </main>
  );
}

function routeMatches(pathname: string, tab: RouteTab) {
  return tab.end ?? true ? pathname === tab.to : pathname.startsWith(tab.to);
}
