import { useEffect } from 'react';
import { QueryKey, useQueryClient } from '@tanstack/react-query';

type HotStreamRule = {
  queryKey: QueryKey;
  eventTypes?: string[];
  resourceTypes?: string[];
};

type HotStreamEvent = {
  metadata?: { type?: string };
  resource_type?: string;
};

export function createHotStreamURL(base: string, filters: { eventTypes?: string[] } = {}) {
  const params = new URLSearchParams();
  for (const eventType of filters.eventTypes || []) params.append('event_type', eventType);
  const suffix = params.size > 0 ? `?${params.toString()}` : '';
  return `${base.replace(/\/$/, '')}/streams/state${suffix}`;
}

export function useHotStreamInvalidation({ url, rules }: { url: string; rules: HotStreamRule[] }) {
  const queryClient = useQueryClient();
  useEffect(() => {
    const source = new EventSource(url);
    source.addEventListener('hotstream', (message) => {
      const event = parseEvent(message);
      for (const rule of rules) {
        if (matches(rule, event)) queryClient.invalidateQueries({ queryKey: rule.queryKey });
      }
    });
    return () => source.close();
  }, [queryClient, rules, url]);
}

function parseEvent(message: Event): HotStreamEvent {
  if (!(message instanceof MessageEvent) || typeof message.data !== 'string') return {};
  try {
    return JSON.parse(message.data) as HotStreamEvent;
  } catch {
    return {};
  }
}

function matches(rule: HotStreamRule, event: HotStreamEvent) {
  return includesOrEmpty(rule.eventTypes, event.metadata?.type) && includesOrEmpty(rule.resourceTypes, event.resource_type);
}

function includesOrEmpty(values: string[] | undefined, value: string | undefined) {
  return !values?.length || (!!value && values.includes(value));
}
