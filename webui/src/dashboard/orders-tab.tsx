import { useMemo, useState } from 'react';
import { Ban, ChevronLeft, ChevronRight } from 'lucide-react';
import { Badge, Button, EmptyBlock, Select, Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../ui';
import type { SmsOrderCodeView, SmsOrderView } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { canCancelStatus, dateTimeText, moneyText, remainingText, statusText } from './sms-format';

type OrderMode = 'active' | 'history';

type OrdersTabProps = {
  orders: SmsOrderView[];
  codes: SmsOrderCodeView[];
  cancelingId?: string;
  onCancel: (id: string) => void;
};

export function OrdersTab({ orders, codes, cancelingId, onCancel }: OrdersTabProps) {
  const [mode, setMode] = useState<OrderMode>('active');
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(20);
  const codesByOrder = useMemo(() => groupCodes(codes), [codes]);
  const rows = useMemo(() => orders.filter((item) => mode === 'active' ? canCancelStatus(item.order?.status) : !canCancelStatus(item.order?.status)), [mode, orders]);
  const pageCount = Math.max(1, Math.ceil(rows.length / pageSize));
  const visible = rows.slice(page * pageSize, page * pageSize + pageSize);

  function changeMode(next: OrderMode) {
    setMode(next);
    setPage(0);
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col p-4">
      <div className="mb-3 flex flex-wrap items-center justify-between gap-3">
        <div className="flex gap-2">
          <Button size="sm" variant={mode === 'active' ? 'default' : 'outline'} onClick={() => changeMode('active')}>进行中</Button>
          <Button size="sm" variant={mode === 'history' ? 'default' : 'outline'} onClick={() => changeMode('history')}>历史订单</Button>
        </div>
        <OrderPager page={page} pageCount={pageCount} pageSize={pageSize} total={rows.length} onPage={setPage} onPageSize={(size) => { setPageSize(size); setPage(0); }} />
      </div>
      <div className="min-h-0 overflow-auto rounded-xl border border-border bg-card">
        <Table>
          <TableHeader><TableRow><TableHead>号码</TableHead><TableHead>Provider</TableHead><TableHead>状态</TableHead><TableHead>剩余</TableHead><TableHead>最新 OTP</TableHead><TableHead>价格</TableHead><TableHead /></TableRow></TableHeader>
          <TableBody>
            {visible.map((item) => <OrderRow key={item.order?.order_id || item.provider_key} item={item} codes={codesByOrder.get(item.order?.order_id || '') || []} cancelingId={cancelingId} onCancel={onCancel} />)}
            {visible.length === 0 && <TableRow><TableCell colSpan={7}><EmptyBlock text="暂无订单" /></TableCell></TableRow>}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

function OrderPager({ page, pageCount, pageSize, total, onPage, onPageSize }: { page: number; pageCount: number; pageSize: number; total: number; onPage: (page: number) => void; onPageSize: (size: number) => void }) {
  return (
    <div className="flex items-center gap-2 text-xs text-muted-foreground">
      <Badge variant="outline">{total} 条</Badge>
      <Select className="h-8 w-24" value={pageSize} onChange={(event) => onPageSize(Number(event.target.value))}>
        <option value="20">20/页</option><option value="50">50/页</option><option value="100">100/页</option>
      </Select>
      <Button size="icon-sm" variant="outline" disabled={page === 0} onClick={() => onPage(Math.max(0, page - 1))}><ChevronLeft className="size-4" /></Button>
      <span>{page + 1}/{pageCount}</span>
      <Button size="icon-sm" variant="outline" disabled={page + 1 >= pageCount} onClick={() => onPage(Math.min(pageCount - 1, page + 1))}><ChevronRight className="size-4" /></Button>
    </div>
  );
}

function OrderRow({ item, codes, cancelingId, onCancel }: { item: SmsOrderView; codes: SmsOrderCodeView[]; cancelingId?: string; onCancel: (id: string) => void }) {
  const order = item.order;
  const id = order?.order_id || '';
  const cancelable = canCancelStatus(order?.status);
  return (
    <TableRow>
      <TableCell className="font-mono text-xs">{order?.phone_number?.e164_number || order?.phone_number?.national_number || '-'}</TableCell>
      <TableCell>{item.provider_key || '-'}</TableCell>
      <TableCell>{statusText(order?.status)}</TableCell>
      <TableCell>{remainingText(order?.expires_at)}</TableCell>
      <TableCell><CodesCell codes={codes} /></TableCell>
      <TableCell>{moneyText(order?.price)}</TableCell>
      <TableCell className="text-right"><Button aria-label="取消订单" title="取消订单" size="icon-sm" variant="outline" disabled={!cancelable || cancelingId === id} onClick={() => onCancel(id)}><Ban className="size-4" /></Button></TableCell>
    </TableRow>
  );
}

function CodesCell({ codes }: { codes: SmsOrderCodeView[] }) {
  if (codes.length === 0) return <span className="text-muted-foreground">-</span>;
  return <div className="grid gap-1">{codes.slice(0, 3).map((item) => <div key={`${item.order_id}-${item.code?.secret_ref?.secret_id}-${item.code?.received_at}`} className="text-xs">{item.code?.secret_ref?.secret_id ? '已捕获' : '-'} · <span className="text-muted-foreground">{dateTimeText(item.code?.received_at)}</span></div>)}{codes.length > 3 && <div className="text-xs text-muted-foreground">+{codes.length - 3} 条历史</div>}</div>;
}

function groupCodes(codes: SmsOrderCodeView[]) {
  const grouped = new Map<string, SmsOrderCodeView[]>();
  for (const item of codes) {
    const id = item.order_id || '';
    if (!id) continue;
    grouped.set(id, [...(grouped.get(id) || []), item]);
  }
  return grouped;
}
