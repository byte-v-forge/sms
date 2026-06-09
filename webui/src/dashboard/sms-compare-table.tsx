import { LoaderCircle, PhoneCall } from 'lucide-react';
import { Badge, Button, Card, EmptyBlock, Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../ui';
import type { SmsPriceOffer, SmsProviderLookupError } from '../proto/byte/v/forge/contracts/sms/v1/sms';
import { normalizeChoiceToken } from './sms-compare-text';
import { dateTimeText, moneyText } from './sms-format';
import { availableCount, offerRowKey } from './sms-compare-data';

type CompareSummaryProps = {
  loading: boolean;
  total: number;
  providerCount: number;
  best?: SmsPriceOffer;
  error?: string;
  providerErrors: SmsProviderLookupError[];
};

type OffersTableProps = {
  offers: SmsPriceOffer[];
  top?: SmsPriceOffer;
  loading: boolean;
  queried: boolean;
  error?: string;
  acquiringOfferId?: string;
  serviceName: string;
  onAcquire: (offer: SmsPriceOffer) => void;
};

export function CompareSummary({ loading, total, providerCount, best, error, providerErrors }: CompareSummaryProps) {
  return (
    <Card className="m-4 mb-3 grid gap-3 p-3">
      <div className="grid gap-2 md:grid-cols-3">
        <SummaryMetric label="报价" value={loading ? '查询中' : `${total} 条`} loading={loading} />
        <SummaryMetric label="覆盖平台" value={`${providerCount} 个`} />
        <SummaryMetric label="最低价" value={best ? moneyText(best.price) : '-'} hint={best ? best.provider_display_name || best.provider_key : undefined} />
      </div>
      {error && <div className="rounded-lg border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">{error}</div>}
      {providerErrors.length > 0 && <ProviderErrors errors={providerErrors} />}
    </Card>
  );
}

function ProviderErrors({ errors }: { errors: SmsProviderLookupError[] }) {
  return (
    <div className="flex flex-wrap gap-2">
      {errors.map((item, index) => (
        <Badge key={`${item.provider_key || item.provider_display_name || 'provider'}-${index}`} variant="outline" className="border-destructive/30 text-destructive">
          {(item.provider_display_name || item.provider_key || '接码源')}：{item.error?.message || item.error?.code || '查询失败'}
        </Badge>
      ))}
    </div>
  );
}

function SummaryMetric({ label, value, hint, loading }: { label: string; value: string; hint?: string; loading?: boolean }) {
  return (
    <div className="rounded-lg border border-border bg-background p-3">
      <div className="text-xs text-muted-foreground">{label}</div>
      <div className="mt-1 flex items-center gap-2 text-lg font-semibold">{loading && <LoaderCircle className="size-4 animate-spin text-muted-foreground" />}{value}</div>
      {hint && <div className="mt-1 truncate text-xs text-muted-foreground">{hint}</div>}
    </div>
  );
}

export function OffersTable({ offers, top, loading, queried, error, acquiringOfferId, serviceName, onAcquire }: OffersTableProps) {
  const bestKey = top ? offerRowKey(top) : '';
  return (
    <div className="min-h-0 flex-1 overflow-auto px-4 pb-4">
      <div className="rounded-xl border border-border bg-card">
        <div className="flex flex-wrap items-center justify-between gap-2 border-b border-border px-3 py-2">
          <div className="text-sm font-medium">报价明细</div>
          {top && <Badge variant="secondary">最低价已高亮</Badge>}
        </div>
        <Table>
          <TableHeader><TableRow><TableHead>平台</TableHead><TableHead>应用</TableHead><TableHead>国家</TableHead><TableHead>价格</TableHead><TableHead>库存</TableHead><TableHead>能力</TableHead><TableHead>观测时间</TableHead><TableHead /></TableRow></TableHeader>
          <TableBody>
            {offers.map((offer) => <OfferRow key={offerRowKey(offer)} offer={offer} bestKey={bestKey} acquiringOfferId={acquiringOfferId} serviceName={serviceName} onAcquire={onAcquire} />)}
            {loading && <TableRow><TableCell colSpan={8}><EmptyBlock text={<span className="inline-flex items-center gap-2"><LoaderCircle className="size-4 animate-spin" />查询中</span>} /></TableCell></TableRow>}
            {!loading && queried && offers.length === 0 && <TableRow><TableCell colSpan={8}><EmptyBlock text={error || '没有匹配报价，请换个关键词或放宽库存/平台筛选'} /></TableCell></TableRow>}
            {!queried && <TableRow><TableCell colSpan={8}><EmptyBlock text="启用接码平台后搜索服务名称即可查询报价" /></TableCell></TableRow>}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

function OfferRow({ offer, bestKey, acquiringOfferId, serviceName, onAcquire }: { offer: SmsPriceOffer; bestKey: string; acquiringOfferId?: string; serviceName: string; onAcquire: (offer: SmsPriceOffer) => void }) {
  const key = offerRowKey(offer);
  return (
    <TableRow className={key === bestKey ? 'bg-secondary/50' : undefined}>
      <TableCell>{offer.provider_display_name || offer.provider_key}</TableCell>
      <TableCell>{offerApplicationName(offer, serviceName)}</TableCell>
      <TableCell>{[offer.country_name, offer.country_iso2, offer.country_calling_code && `+${offer.country_calling_code}`].filter(Boolean).join(' · ') || '-'}</TableCell>
      <TableCell>{key === bestKey ? <Badge>{moneyText(offer.price)}</Badge> : moneyText(offer.price)}</TableCell>
      <TableCell>{availableCount(offer)}</TableCell>
      <TableCell><CapabilityBadges offer={offer} /></TableCell>
      <TableCell>{dateTimeText(offer.observed_at)}</TableCell>
      <TableCell className="text-right"><Button aria-label="按此报价取号" title="按此报价取号" size="icon-sm" disabled={!offer.offer_ref || acquiringOfferId === key} onClick={() => onAcquire(offer)}>{acquiringOfferId === key ? <LoaderCircle className="size-4 animate-spin" /> : <PhoneCall className="size-4" />}</Button></TableCell>
    </TableRow>
  );
}

function offerApplicationName(offer: SmsPriceOffer, fallback: string) {
  const name = visibleApplicationName(offer.application_name, offer.application_key);
  if (name) return name;
  return visibleApplicationName(fallback, offer.application_key) || '-';
}

function visibleApplicationName(value: string, key: string) {
  const name = value.trim();
  if (!name || looksLikeShortCode(name, key)) return '';
  return name;
}

function looksLikeShortCode(value: string, key: string) {
  return value.length <= 3 && value === value.toLowerCase() && normalizeChoiceToken(value) === normalizeChoiceToken(key);
}

function CapabilityBadges({ offer }: { offer: SmsPriceOffer }) {
  const tags = [offer.supports_cancel && '可取消', offer.supports_additional_code && '重发', offer.requires_mark_message_sent && '需标记'].filter((value): value is string => !!value);
  if (tags.length === 0) return <span className="text-muted-foreground">-</span>;
  return <div className="flex flex-wrap gap-1">{tags.map((tag) => <Badge key={tag} variant="outline">{tag}</Badge>)}</div>;
}
