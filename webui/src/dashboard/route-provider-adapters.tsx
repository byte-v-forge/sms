import { Input, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/dashboard/module-kit';
import {
  SmsRouteFieldScope,
  SmsRouteOptionSource,
  type SmsProviderRouteField,
  type SmsProviderRouteOptions,
  type SmsRouteCandidate,
  type SmsRouteOption
} from '@/proto/byte/v/forge/sms/internal/v1/sms_internal';

type Props = {
  route: SmsRouteCandidate;
  fields: SmsProviderRouteField[];
  options?: SmsProviderRouteOptions;
  onChange: (route: SmsRouteCandidate) => void;
};

export function ProviderRouteFields({ route, fields, options, onChange }: Props) {
  return (
    <div className="grid grid-cols-2 gap-2">
      {fields.map((field) => (
        <FieldInput key={`${field.scope}-${field.field_key}-${field.label}`} field={field} route={route} fields={fields} options={options} onChange={onChange} />
      ))}
    </div>
  );
}

function FieldInput({ field, route, options, onChange }: Props & { field: SmsProviderRouteField }) {
  const value = readField(route, field);
  const choices = routeOptions(options, field.option_source).filter((item) => item.value);
  return (
    <div className="grid gap-1">
      <Label>{field.label}</Label>
      {choices.length > 0 ? (
        <Select value={value || undefined} onValueChange={(next) => onChange(writeChoice(route, field, choices.find((item) => item.value === next)))}>
          <SelectTrigger><SelectValue placeholder={choices.length ? '选择' : '无可选项'} /></SelectTrigger>
          <SelectContent>{choices.map((item) => <SelectItem key={item.value} value={item.value}>{item.label || item.value}</SelectItem>)}</SelectContent>
        </Select>
      ) : <Input value={value} onChange={(event) => onChange(writeField(route, field, event.target.value))} />}
    </div>
  );
}

function readField(route: SmsRouteCandidate, field: SmsProviderRouteField) {
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_OPTION) return route.provider_options?.[field.field_key] || '';
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_MIN_PRICE) return route.min_price?.amount_decimal || '';
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_MAX_PRICE) return route.max_price?.amount_decimal || '';
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_ROUTE) return field.field_key === 'upstream_service_key' ? route.upstream_service_key : route.provider_country_id;
  return '';
}

function writeChoice(route: SmsRouteCandidate, field: SmsProviderRouteField, choice?: SmsRouteOption): SmsRouteCandidate {
  if (!choice) return route;
  const next = writeField(route, field, choice.value);
  if (field.option_source === SmsRouteOptionSource.SMS_ROUTE_OPTION_SOURCE_COUNTRIES) {
    next.target = {
      application_key: next.target?.application_key || '',
      country_iso2: choice.metadata?.country_iso2 || next.target?.country_iso2 || '',
      country_calling_code: choice.metadata?.country_calling_code || next.target?.country_calling_code || '',
      min_price: next.target?.min_price,
      max_price: next.target?.max_price
    };
  }
  if (field.option_source === SmsRouteOptionSource.SMS_ROUTE_OPTION_SOURCE_SERVICES && choice.metadata?.application_key) {
    next.target = { application_key: choice.metadata.application_key, country_iso2: next.target?.country_iso2 || '', country_calling_code: next.target?.country_calling_code || '', min_price: next.target?.min_price, max_price: next.target?.max_price };
  }
  return next;
}

function writeField(route: SmsRouteCandidate, field: SmsProviderRouteField, value: string): SmsRouteCandidate {
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_OPTION) {
    return { ...route, provider_options: { ...(route.provider_options || {}), [field.field_key]: value } };
  }
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_MIN_PRICE) {
    return { ...route, min_price: { currency_code: route.min_price?.currency_code || 'USD', amount_decimal: value } };
  }
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_MAX_PRICE) {
    return { ...route, max_price: { currency_code: route.max_price?.currency_code || 'USD', amount_decimal: value } };
  }
  if (field.scope === SmsRouteFieldScope.SMS_ROUTE_FIELD_SCOPE_ROUTE && field.field_key === 'upstream_service_key') return { ...route, upstream_service_key: value };
  return { ...route, provider_country_id: value };
}

function routeOptions(options: SmsProviderRouteOptions | undefined, source: SmsRouteOptionSource) {
  if (source === SmsRouteOptionSource.SMS_ROUTE_OPTION_SOURCE_SERVICES) return options?.services || [];
  if (source === SmsRouteOptionSource.SMS_ROUTE_OPTION_SOURCE_COUNTRIES) return options?.countries || [];
  if (source === SmsRouteOptionSource.SMS_ROUTE_OPTION_SOURCE_OPERATORS) return options?.operators || [];
  if (source === SmsRouteOptionSource.SMS_ROUTE_OPTION_SOURCE_UPSTREAM_PROVIDERS) return options?.upstream_providers || [];
  return [];
}
