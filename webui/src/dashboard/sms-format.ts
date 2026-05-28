import type { DecimalMoney } from '@byte-v-forge/common-ui/proto/byte/v/forge/contracts/sms/v1/sms';

export function moneyText(money?: DecimalMoney) {
  if (!money?.amount_decimal) return '-';
  return [money.currency_code, money.amount_decimal].filter(Boolean).join(' ');
}

export function remainingText(expiresAt?: string) {
  if (!expiresAt) return '-';
  const seconds = Math.max(0, Math.floor((new Date(expiresAt).getTime() - Date.now()) / 1000));
  const minutes = Math.floor(seconds / 60);
  return `${minutes}:${String(seconds % 60).padStart(2, '0')}`;
}

export function statusText(status?: string) {
  const labels: Record<string, string> = {
    SMS_ORDER_STATUS_ACQUIRE_REQUESTED: '取号中',
    SMS_ORDER_STATUS_PENDING_CODE: '等待验证码',
    SMS_ORDER_STATUS_MESSAGE_SENT: '已触发短信',
    SMS_ORDER_STATUS_CODE_RECEIVED: '已收到',
    SMS_ORDER_STATUS_ADDITIONAL_CODE_REQUESTED: '重发中',
    SMS_ORDER_STATUS_COMPLETED: '已完成',
    SMS_ORDER_STATUS_CANCELED: '已取消',
    SMS_ORDER_STATUS_EXPIRED: '已过期',
    SMS_ORDER_STATUS_FAILED: '失败'
  };
  return labels[status || ''] || status || '-';
}

export function canCancelStatus(status?: string) {
  return !['SMS_ORDER_STATUS_COMPLETED', 'SMS_ORDER_STATUS_CANCELED', 'SMS_ORDER_STATUS_EXPIRED', 'SMS_ORDER_STATUS_FAILED'].includes(status || '');
}
