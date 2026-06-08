import { LoaderCircle, PhoneCall } from 'lucide-react';
import { Badge, Button, Card, EmptyBlock, Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../ui';
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
    <Card className="m-4 mb-3 grid gap-2 p-3">
      <div className="flex flex-wrap items-center gap-3 text-sm">
        <Badge variant="secondary">{loading ? '查询中' : `${total} 条报价`}</Badge>
        <Badge variant="outline">{providerCount} 个平台</Badge>
        <span className="text-muted-foreground">最低价：{best ? `${best.provider_display_name || best.provider_key} · ${moneyText(best.price)}` : '-'}</span>
      </div>
      {error && <div className="rounded-lg border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">{error}</div>}
    </Card>
  );
}

export function OffersTable({ offers, top, loading, queried, error, acquiringOfferId, onAcquire }: OffersTableProps) {
  const bestKey = top ? offerRowKey(top) : '';
  return (
    <div className="min-h-0 flex-1 overflow-auto px-4 pb-4">
      <div className="rounded-xl border border-border bg-card">
        <Table>
          <TableHeader><TableRow><TableHead>平台</TableHead><TableHead>应用</TableHead><TableHead>国家</TableHead><TableHead>价格</TableHead><TableHead>库存</TableHead><TableHead>能力</TableHead><TableHead>观测时间</TableHead><TableHead /></TableRow></TableHeader>
          <TableBody>
            {offers.map((offer) => <OfferRow key={offerRowKey(offer)} offer={offer} bestKey={bestKey} acquiringOfferId={acquiringOfferId} onAcquire={onAcquire} />)}
            {loading && <TableRow><TableCell colSpan={8}><EmptyBlock text={<span className="inline-flex items-center gap-2"><LoaderCircle className="size-4 animate-spin" />查询中</span>} /></TableCell></TableRow>}
            {!loading && queried && offers.length === 0 && <TableRow><TableCell colSpan={8}><EmptyBlock text={error || '暂无可用报价，请调整平台、应用或国家条件'} /></TableCell></TableRow>}
            {!queried && <TableRow><TableCell colSpan={8}><EmptyBlock text="输入应用和国家后查询多个接码平台报价" /></TableCell></TableRow>}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

function OfferRow({ offer, bestKey, acquiringOfferId, onAcquire }: { offer: SmsPriceOffer; bestKey: string; acquiringOfferId?: string; onAcquire: (offer: SmsPriceOffer) => void }) {
  const key = offerRowKey(offer);
  return (
    <TableRow className={key === bestKey ? 'bg-secondary/50' : undefined}>
      <TableCell>{offer.provider_display_name || offer.provider_key}</TableCell>
      <TableCell>{offer.application_name || offer.application_key || '-'}</TableCell>
      <TableCell>{[offer.country_name, offer.country_iso2, offer.country_calling_code && `+${offer.country_calling_code}`].filter(Boolean).join(' · ') || '-'}</TableCell>
      <TableCell>{key === bestKey ? <Badge>{moneyText(offer.price)}</Badge> : moneyText(offer.price)}</TableCell>
      <TableCell>{offer.available_count}</TableCell>
      <TableCell><CapabilityBadges offer={offer} /></TableCell>
      <TableCell>{dateTimeText(offer.observed_at)}</TableCell>
      <TableCell className="text-right"><Button aria-label="按此报价取号" title="按此报价取号" size="icon-sm" disabled={!offer.offer_ref || acquiringOfferId === key} onClick={() => onAcquire(offer)}>{acquiringOfferId === key ? <LoaderCircle className="size-4 animate-spin" /> : <PhoneCall className="size-4" />}</Button></TableCell>
    </TableRow>
  );
}

function CapabilityBadges({ offer }: { offer: SmsPriceOffer }) {
  const tags = [offer.supports_cancel && '可取消', offer.supports_additional_code && '重发', offer.requires_mark_message_sent && '需标记'].filter((value): value is string => !!value);
  if (tags.length === 0) return <span className="text-muted-foreground">-</span>;
  return <div className="flex flex-wrap gap-1">{tags.map((tag) => <Badge key={tag} variant="outline">{tag}</Badge>)}</div>;
}
