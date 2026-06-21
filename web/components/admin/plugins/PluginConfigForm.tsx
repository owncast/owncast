import React, { FC, useEffect, useState } from 'react';
import { Alert, Button, Form, Input, InputNumber, Spin, Switch, Typography, message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Plugin } from '../../../interfaces/plugin';
import { Localization } from '../../../types/localization';
import { fetchData, pluginConfigUrl } from '../../../utils/apis';

const { Text } = Typography;

// secretField masks values for keys that look like credentials, so secrets
// aren't shown in plain text in the admin form.
const secretField = (key: string) => /secret|password|token|apikey|api_key|\bkey\b/i.test(key);

// PluginConfigForm auto-renders an editable settings form for a plugin from its
// manifest `config` schema — no plugin-authored HTML. It loads the current
// effective values (admin overrides, else manifest defaults) and saves changes
// back through the host config endpoint, which the plugin then reads via
// owncast.config.get().
export const PluginConfigForm: FC<{ plugin: Plugin }> = ({ plugin }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const schema = plugin.config || {};
  const keys = Object.keys(schema);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);
    fetchData(pluginConfigUrl(plugin.slug))
      .then(values => {
        if (!cancelled) form.setFieldsValue(values);
      })
      .catch(() => {
        if (!cancelled) setError(t(Localization.Admin.Plugins.configLoadError));
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [plugin.slug]);

  const onFinish = async (values: Record<string, any>) => {
    setSaving(true);
    setError(null);
    try {
      await fetchData(pluginConfigUrl(plugin.slug), { method: 'POST', data: values });
      message.success(t(Localization.Admin.Plugins.configSaved));
    } catch (e: any) {
      setError(e?.message || t(Localization.Admin.Plugins.configSaveError));
    } finally {
      setSaving(false);
    }
  };

  if (keys.length === 0) {
    return <Text type="secondary">{t(Localization.Admin.Plugins.configEmpty)}</Text>;
  }
  if (loading) {
    return <Spin />;
  }

  return (
    <Form form={form} layout="vertical" onFinish={onFinish}>
      {error && <Alert type="error" message={error} style={{ marginBottom: 16 }} />}
      {keys.map(key => {
        const field = schema[key];
        const label = field.description || key;
        if (field.type === 'boolean') {
          return (
            <Form.Item key={key} name={key} label={label} valuePropName="checked">
              <Switch />
            </Form.Item>
          );
        }
        if (field.type === 'number') {
          return (
            <Form.Item key={key} name={key} label={label}>
              <InputNumber style={{ width: '100%' }} />
            </Form.Item>
          );
        }
        return (
          <Form.Item key={key} name={key} label={label}>
            {secretField(key) ? <Input.Password autoComplete="new-password" /> : <Input />}
          </Form.Item>
        );
      })}
      <Form.Item>
        <Button type="primary" htmlType="submit" loading={saving}>
          {t(Localization.Admin.Plugins.configSave)}
        </Button>
      </Form.Item>
    </Form>
  );
};
