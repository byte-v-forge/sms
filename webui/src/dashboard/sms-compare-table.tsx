import { PhoneOutlined } from '@ant-design/icons';
import { Alert, Button, Card, Empty, Flex, Space, Statistic, Table, Tag, Typography } from 'antd';
import type { TableColumnsType } from 'antd';
import type { SmsPriceOffer } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { dateTimeText, moneyText } from './sms-format';
import { offerRowKey } from './sms-compare-data';

type CompareSummaryProps = {
  loading: boolean;
  total: number;
  providerCount: number;
  best?: SmsPriceOffer;
  error?: string;
};

type OffersTableProps = {
  offers: SmsPriceOffer[];
  top?: SmsPriceOffer;
  loading: boolean;
  queried: boolean;
  error?: string;
  acquiringOfferId?: string;
  onAcquire: (offer: SmsPriceOffer) => void;
};

export function CompareSummary({ loading, total, providerCount, best, error }: CompareSummaryProps) {
  return (
    <Card size="small" style={{ margin: '12px 16px' }}>
      <Flex gap={24} wrap="wrap" align="center">
        <Statistic title="报价" value={loading ? '查询中' : total} suffix={loading ? '' : '条'} />
        <Statistic title="平台" value={providerCount} suffix="个" />
        <Statistic title="最低价" value={best ? moneyText(best.price) : '-'} />
        {best && <Typography.Text type="secondary">{best.provider_display_name || best.provider_key}</Typography.Text>}
      </Flex>
      {error && <Alert style={{ marginTop: 12 }} type="error" message={error} showIcon />}
    </Card>
  );
}

export function OffersTable({ offers, top, loading, queried, error, acquiringOfferId, onAcquire }: OffersTableProps) {
  const bestKey = top ? offerRowKey(top) : '';
  const columns: TableColumnsType<SmsPriceOffer> = [
    { title: '平台', dataIndex: 'provider_display_name', render: (_, offer) => offer.provider_display_name || offer.provider_key },
    { title: '应用', dataIndex: 'application_name', render: (_, offer) => offer.application_name || offer.application_key || '-' },
    { title: '国家', render: (_, offer) => countryText(offer) },
    { title: '价格', dataIndex: 'price', render: (_, offer) => <PriceCell offer={offer} bestKey={bestKey} /> },
    { title: '库存', dataIndex: 'available_count' },
    { title: '能力', render: (_, offer) => <CapabilityTags offer={offer} /> },
    { title: '观测时间', dataIndex: 'observed_at', render: (value: string) => dateTimeText(value) },
    { title: '', width: 64, align: 'right', render: (_, offer) => <AcquireButton offer={offer} acquiringOfferId={acquiringOfferId} onAcquire={onAcquire} /> }
  ];
  return (
    <div style={{ minHeight: 0, flex: 1, padding: '0 16px 16px' }}>
      <Table
        rowKey={offerRowKey}
        columns={columns}
        dataSource={offers}
        loading={loading}
        size="middle"
        sticky
        locale={{ emptyText: <Empty description={queried ? error || '暂无可用报价，请调整平台、应用或国家条件' : '输入应用和国家后查询多个接码平台报价'} /> }}
        pagination={{ defaultPageSize: 20, showSizeChanger: true, showTotal: (total) => `${total} 条` }}
        scroll={{ y: 'calc(100vh - 390px)', x: 1100 }}
      />
    </div>
  );
}

function PriceCell({ offer, bestKey }: { offer: SmsPriceOffer; bestKey: string }) {
  const content = moneyText(offer.price);
  return offerRowKey(offer) === bestKey ? <Tag color="green">最低 · {content}</Tag> : content;
}

function AcquireButton({ offer, acquiringOfferId, onAcquire }: { offer: SmsPriceOffer; acquiringOfferId?: string; onAcquire: (offer: SmsPriceOffer) => void }) {
  const key = offerRowKey(offer);
  return (
    <Button aria-label="按此报价取号" title="按此报价取号" icon={<PhoneOutlined />} size="small" disabled={!offer.offer_ref} loading={acquiringOfferId === key} onClick={() => onAcquire(offer)} />
  );
}

function CapabilityTags({ offer }: { offer: SmsPriceOffer }) {
  const tags = [
    offer.supports_cancel && '可取消',
    offer.supports_additional_code && '重发',
    offer.requires_mark_message_sent && '需标记'
  ].filter((value): value is string => !!value);
  if (tags.length === 0) return <Typography.Text type="secondary">-</Typography.Text>;
  return <Space size={[4, 4]} wrap>{tags.map((tag) => <Tag key={tag}>{tag}</Tag>)}</Space>;
}

function countryText(offer: SmsPriceOffer) {
  return [offer.country_name, offer.country_iso2, offer.country_calling_code && `+${offer.country_calling_code}`].filter(Boolean).join(' · ') || '-';
}
