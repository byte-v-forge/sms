import { useState } from 'react';
import { Ban, ChevronLeft, ChevronRight } from 'lucide-react';
import { Button, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@byte-v-forge/common-ui';
import type { SmsOrderView } from '../proto/byte/v/forge/sms/internal/v1/sms_internal';
import { canCancelStatus, moneyText, remainingText, statusText } from './sms-format';

type OrdersTabProps = {
  orders: SmsOrderView[];
  cancelingId?: string;
  onCancel: (id: string) => void;
};

export function OrdersTab({ orders, cancelingId, onCancel }: OrdersTabProps) {
  const [mode, setMode] = useState<'active' | 'history'>('active');
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(20);
  const rows = orders.filter((item) => mode === 'active' ? canCancelStatus(item.order?.status) : !canCancelStatus(item.order?.status));
  const pageCount = Math.max(1, Math.ceil(rows.length / pageSize));
  const visible = rows.slice(page * pageSize, page * pageSize + pageSize);

  function changeMode(next: 'active' | 'history') {
    setMode(next);
    setPage(0);
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <div className="flex items-center justify-between border-b border-border/70 px-3 py-2">
        <div className="flex gap-2">
          <Button size="sm" variant={mode === 'active' ? 'default' : 'outline'} onClick={() => changeMode('active')}>进行中</Button>
          <Button size="sm" variant={mode === 'history' ? 'default' : 'outline'} onClick={() => changeMode('history')}>历史订单</Button>
        </div>
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <Select value={String(pageSize)} onValueChange={(value) => { setPageSize(Number(value)); setPage(0); }}>
            <SelectTrigger className="h-8 w-24"><SelectValue /></SelectTrigger>
            <SelectContent>
              <SelectItem value="20">20/页</SelectItem>
              <SelectItem value="50">50/页</SelectItem>
              <SelectItem value="100">100/页</SelectItem>
            </SelectContent>
          </Select>
          <Button size="icon-sm" variant="outline" disabled={page === 0} onClick={() => setPage((value) => Math.max(0, value - 1))}><ChevronLeft className="size-4" /></Button>
          <span>{page + 1}/{pageCount}</span>
          <Button size="icon-sm" variant="outline" disabled={page + 1 >= pageCount} onClick={() => setPage((value) => Math.min(pageCount - 1, value + 1))}><ChevronRight className="size-4" /></Button>
        </div>
      </div>
      <div className="min-h-0 overflow-auto">
        <Table>
          <TableHeader><TableRow><TableHead>号码</TableHead><TableHead>Provider</TableHead><TableHead>状态</TableHead><TableHead>剩余</TableHead><TableHead>最新OTP</TableHead><TableHead>价格</TableHead><TableHead /></TableRow></TableHeader>
          <TableBody>
            {visible.map((item) => <OrderRow key={item.order?.order_id || item.provider_key} item={item} cancelingId={cancelingId} onCancel={onCancel} />)}
            {visible.length === 0 && <TableRow><TableCell colSpan={7} className="h-24 text-center text-muted-foreground">暂无订单</TableCell></TableRow>}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

function OrderRow({ item, cancelingId, onCancel }: { item: SmsOrderView; cancelingId?: string; onCancel: (id: string) => void }) {
  const order = item.order;
  const id = order?.order_id || '';
  const cancelable = canCancelStatus(order?.status);
  return (
    <TableRow>
      <TableCell className="font-mono text-xs">{order?.phone_number?.e164_number || order?.phone_number?.national_number || '-'}</TableCell>
      <TableCell>{item.provider_key || '-'}</TableCell>
      <TableCell>{statusText(order?.status)}</TableCell>
      <TableCell>{remainingText(order?.expires_at)}</TableCell>
      <TableCell className="font-mono text-xs">-</TableCell>
      <TableCell>{moneyText(order?.price)}</TableCell>
      <TableCell className="text-right">
        <Button size="icon-sm" variant="outline" disabled={!cancelable || cancelingId === id} onClick={() => onCancel(id)}><Ban className="size-4" /></Button>
      </TableCell>
    </TableRow>
  );
}
