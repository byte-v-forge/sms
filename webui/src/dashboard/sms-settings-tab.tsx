import { useEffect, useState } from 'react';
import { DeleteOutlined, SaveOutlined } from '@ant-design/icons';
import { Button, Card, Descriptions, Empty, Flex, Input, Popconfirm, Space, Spin, Switch, Tag, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { getSmsProviderBalance, smsKeys, type SaveSmsProviderSettingRequest, type SmsProviderOption, type SmsProviderSetting } from './sms-api';
import { moneyText } from './sms-format';

type ProviderCardDescriptor = {
  option: SmsProviderOption;
  config?: SmsProviderSetting;
};

type SmsSettingsTabProps = {
  providerOptions: SmsProviderOption[];
  configs: SmsProviderSetting[];
  busy?: boolean;
  savingProviderKey?: string;
  deletingProviderKey?: string;
  onSave: (input: SaveSmsProviderSettingRequest) => void;
  onDelete: (id: string) => void;
};

export function SmsSettingsTab(props: SmsSettingsTabProps) {
  const providers = mergeProviderOptions(props.providerOptions, props.configs);
  const enabledCount = props.configs.filter((config) => config.enabled).length;
  return (
    <div className="sms-fill" style={{ overflow: 'auto', padding: 16 }}>
      <Flex align="center" justify="space-between" gap={12} wrap="wrap" style={{ marginBottom: 16 }}>
        <div>
          <Typography.Title level={5} style={{ margin: 0 }}>SMS Provider</Typography.Title>
          <Typography.Text type="secondary">每个平台在卡片内直接配置，配置仅保存在 SMS 服务。</Typography.Text>
        </div>
        <Space wrap>
          <Tag color="blue">{props.providerOptions.length} 插件</Tag>
          <Tag>{props.configs.length} 已配置</Tag>
          <Tag color={enabledCount > 0 ? 'green' : 'default'}>{enabledCount} 启用</Tag>
        </Space>
      </Flex>
      <Spin spinning={!!props.busy && providers.length === 0}>
        <Flex align="flex-start" gap={16} wrap="wrap">
          {providers.map(({ option, config }) => (
            <ProviderCard
              key={option.provider_key}
              config={config}
              provider={option}
              saving={props.savingProviderKey === option.provider_key}
              deleting={props.deletingProviderKey === option.provider_key}
              onSave={props.onSave}
              onDelete={props.onDelete}
            />
          ))}
          {!props.busy && providers.length === 0 && <Empty description="暂无接码源插件" />}
        </Flex>
      </Spin>
    </div>
  );
}

function ProviderCard({ config, provider, saving, deleting, onSave, onDelete }: {
  config?: SmsProviderSetting;
  provider: SmsProviderOption;
  saving?: boolean;
  deleting?: boolean;
  onSave: (input: SaveSmsProviderSettingRequest) => void;
  onDelete: (id: string) => void;
}) {
  const [apiKey, setAPIKey] = useState('');
  const [enabled, setEnabled] = useState(true);
  useEffect(() => {
    setAPIKey('');
    setEnabled(config?.enabled ?? true);
  }, [config?.enabled, config?.provider_key]);
  const balance = useQuery({
    queryKey: smsKeys.balance(provider.provider_key),
    queryFn: () => getSmsProviderBalance(provider.provider_key),
    enabled: !!config?.provider_key && config.enabled && !!config.api_key_set,
    refetchInterval: 60000
  });
  const hasCredential = !!apiKey.trim() || !!config?.api_key_set;
  const canSave = enabled === false || hasCredential;
  return (
    <Card style={{ width: 340 }} title={provider.display_name || provider.provider_key} extra={providerStatusTag(config, enabled)}>
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        <Typography.Text type="secondary">{provider.provider_key}</Typography.Text>
        <Input.Password placeholder={config?.api_key_set ? '留空保留 API Key' : 'Provider API Key'} value={apiKey} onChange={(event) => setAPIKey(event.target.value)} />
        <Flex align="center" justify="space-between">
          <Typography.Text>启用接码源</Typography.Text>
          <Switch checked={enabled} onChange={setEnabled} />
        </Flex>
        <Descriptions column={1} size="small" bordered>
          <Descriptions.Item label="余额">{balanceText(config, balance.isLoading, moneyText(balance.data?.balance))}</Descriptions.Item>
          <Descriptions.Item label="API Key">{config?.api_key_set ? '已配置' : '未配置'}</Descriptions.Item>
        </Descriptions>
        <Flex justify="flex-end" gap={8}>
          <Button aria-label="保存" title="保存" type="primary" icon={<SaveOutlined />} loading={saving} disabled={!canSave} onClick={() => onSave({ provider_key: provider.provider_key, enabled, api_key: apiKey.trim() || undefined })} />
          <Popconfirm title="删除接码源配置？" okText="删除" cancelText="保留" disabled={!config?.provider_key || deleting} onConfirm={() => onDelete(provider.provider_key)}>
            <Button aria-label="删除" title="删除" danger icon={<DeleteOutlined />} loading={deleting} disabled={!config?.provider_key} />
          </Popconfirm>
        </Flex>
      </Space>
    </Card>
  );
}

function mergeProviderOptions(providerOptions: SmsProviderOption[], configs: SmsProviderSetting[]): ProviderCardDescriptor[] {
  const configsByKey = new Map(configs.map((config) => [config.provider_key, config]));
  const providers = providerOptions.map((option) => ({ option, config: configsByKey.get(option.provider_key) }));
  const optionKeys = new Set(providerOptions.map((option) => option.provider_key));
  for (const config of configs) {
    if (!optionKeys.has(config.provider_key)) providers.push({ option: { provider_key: config.provider_key, display_name: config.provider_key }, config });
  }
  return providers;
}

function balanceText(config: SmsProviderSetting | undefined, loading: boolean, balance: string) {
  if (!config) return '-';
  return loading ? '读取中' : balance;
}

function providerStatusTag(config: SmsProviderSetting | undefined, enabled: boolean) {
  if (!config) return <Tag>未配置</Tag>;
  return enabled ? <Tag color="green">启用</Tag> : <Tag color="default">停用</Tag>;
}
