import { useMemo, useState } from 'react';
import { StopOutlined } from '@ant-design/icons';
import { Button, Empty, Flex, Popconfirm, Segmented, Table, Typography } from 'antd';
import type { TableColumnsType } from 'antd';
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
  const codesByOrder = useMemo(() => groupCodes(codes), [codes]);
  const rows = useMemo(() => orders.filter((item) => mode === 'active' ? canCancelStatus(item.order?.status) : !canCancelStatus(item.order?.status)), [mode, orders]);
  const columns = useMemo<TableColumnsType<SmsOrderView>>(() => [
    { title: '号码', dataIndex: ['order', 'phone_number', 'e164_number'], render: (_, item) => orderNumber(item) },
    { title: 'Provider', dataIndex: 'provider_key', render: (value: string) => value || '-' },
    { title: '状态', dataIndex: ['order', 'status'], render: (value: string) => statusText(value) },
    { title: '剩余', dataIndex: ['order', 'expires_at'], render: (value: string) => remainingText(value) },
    { title: '最新 OTP', render: (_, item) => <CodesCell codes={codesByOrder.get(item.order?.order_id || '') || []} /> },
    { title: '价格', dataIndex: ['order', 'price'], render: (_, item) => moneyText(item.order?.price) },
    { title: '', width: 64, align: 'right', render: (_, item) => <CancelAction item={item} cancelingId={cancelingId} onCancel={onCancel} /> }
  ], [cancelingId, codesByOrder, onCancel]);
  return (
    <div className="sms-fill" style={{ padding: 16 }}>
      <Flex align="center" justify="space-between" style={{ marginBottom: 12 }}>
        <Segmented<OrderMode> value={mode} onChange={setMode} options={[{ label: '进行中', value: 'active' }, { label: '历史订单', value: 'history' }]} />
        <Typography.Text type="secondary">{rows.length} 条订单</Typography.Text>
      </Flex>
      <Table
        rowKey={(item) => item.order?.order_id || item.provider_key}
        columns={columns}
        dataSource={rows}
        size="middle"
        sticky
        locale={{ emptyText: <Empty description="暂无订单" /> }}
        pagination={{ defaultPageSize: 20, showSizeChanger: true, pageSizeOptions: [20, 50, 100], showTotal: (total) => `${total} 条` }}
        scroll={{ y: 'calc(100vh - 260px)', x: 980 }}
      />
    </div>
  );
}

function CancelAction({ item, cancelingId, onCancel }: { item: SmsOrderView; cancelingId?: string; onCancel: (id: string) => void }) {
  const id = item.order?.order_id || '';
  const disabled = !canCancelStatus(item.order?.status) || cancelingId === id;
  return (
    <Popconfirm title="取消号码订单？" okText="取消订单" cancelText="保留" disabled={disabled} onConfirm={() => onCancel(id)}>
      <Button aria-label="取消订单" title="取消订单" icon={<StopOutlined />} size="small" disabled={disabled} />
    </Popconfirm>
  );
}

function CodesCell({ codes }: { codes: SmsOrderCodeView[] }) {
  if (codes.length === 0) return <Typography.Text type="secondary">-</Typography.Text>;
  return (
    <Flex vertical gap={2}>
      {codes.slice(0, 3).map((item) => (
        <Typography.Text key={`${item.order_id}-${item.code?.secret_ref?.secret_id}-${item.code?.received_at}`} style={{ fontSize: 12 }}>
          {item.code?.secret_ref?.secret_id ? '已捕获' : '-'} · {dateTimeText(item.code?.received_at)}
        </Typography.Text>
      ))}
      {codes.length > 3 && <Typography.Text type="secondary" style={{ fontSize: 12 }}>+{codes.length - 3} 条历史</Typography.Text>}
    </Flex>
  );
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

function orderNumber(item: SmsOrderView) {
  return <Typography.Text code>{item.order?.phone_number?.e164_number || item.order?.phone_number?.national_number || '-'}</Typography.Text>;
}
