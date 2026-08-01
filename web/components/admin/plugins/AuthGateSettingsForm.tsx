import { FC, useEffect, useState } from 'react';
import { Alert, Button, Form, Radio, Space, Spin, Typography, message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Plugin } from '../../../interfaces/plugin';
import { Localization } from '../../../types/localization';
import { authGateSettingsUrl, fetchData } from '../../../utils/apis';

const { Paragraph, Text } = Typography;

type AuthGateAccessMode = 'website-only' | 'website-and-stream' | 'website-stream-and-status';

type AuthGateSettings = {
  accessMode: AuthGateAccessMode;
};

export const AuthGateSettingsForm: FC<{ plugin: Plugin }> = ({ plugin }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm<AuthGateSettings>();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);
    fetchData(authGateSettingsUrl(plugin.slug))
      .then(values => {
        if (!cancelled) form.setFieldsValue(values);
      })
      .catch(() => {
        if (!cancelled) setError(t(Localization.Admin.Plugins.authSettingsLoadError));
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [plugin.slug]);

  const onFinish = async (values: AuthGateSettings) => {
    setSaving(true);
    setError(null);
    try {
      await fetchData(authGateSettingsUrl(plugin.slug), { method: 'POST', data: values });
      message.success(t(Localization.Admin.Plugins.authSettingsSaved));
    } catch (e: unknown) {
      const errorMessage =
        e instanceof Error ? e.message : t(Localization.Admin.Plugins.authSettingsSaveError);
      setError(errorMessage);
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <Spin />;

  return (
    <Form form={form} layout="vertical" onFinish={onFinish}>
      {error && <Alert type="error" message={error} style={{ marginBottom: 16 }} />}
      <Paragraph>
        <Text>{t(Localization.Admin.Plugins.authSettingsDescription)}</Text>
      </Paragraph>
      <Form.Item name="accessMode" label={t(Localization.Admin.Plugins.authSettingsAccessMode)}>
        <Radio.Group>
          <Space direction="vertical" size="middle">
            <Radio value="website-only">
              <Text strong>{t(Localization.Admin.Plugins.authSettingsWebsiteOnly)}</Text>
              <br />
              <Text type="secondary">
                {t(Localization.Admin.Plugins.authSettingsWebsiteOnlyDescription)}
              </Text>
            </Radio>
            <Radio value="website-and-stream">
              <Text strong>{t(Localization.Admin.Plugins.authSettingsWebsiteAndStream)}</Text>
              <br />
              <Text type="secondary">
                {t(Localization.Admin.Plugins.authSettingsWebsiteAndStreamDescription)}
              </Text>
            </Radio>
            <Radio value="website-stream-and-status">
              <Text strong>{t(Localization.Admin.Plugins.authSettingsWebsiteStreamAndStatus)}</Text>
              <br />
              <Text type="secondary">
                {t(Localization.Admin.Plugins.authSettingsWebsiteStreamAndStatusDescription)}
              </Text>
            </Radio>
          </Space>
        </Radio.Group>
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" loading={saving}>
          {t(Localization.Admin.Plugins.authSettingsSave)}
        </Button>
      </Form.Item>
    </Form>
  );
};
