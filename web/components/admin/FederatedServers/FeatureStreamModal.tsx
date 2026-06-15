import React, { FC, useState } from 'react';
import { Modal, Form, Input, Alert, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { isValidUrl } from '../../../utils/validators';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';

const { Title } = Typography;

export interface FeatureStreamModalProps {
  open: boolean;
  onCancel: () => void;
  onOk: (url: string) => Promise<void>;
}

export const FeatureStreamModal: FC<FeatureStreamModalProps> = ({ open, onCancel, onOk }) => {
  const [form] = Form.useForm();
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      setError(null);

      // Normalize the URL
      let normalizedUrl = values.url.trim();
      if (!normalizedUrl.startsWith('http://') && !normalizedUrl.startsWith('https://')) {
        normalizedUrl = `https://${normalizedUrl}`;
      }

      // Additional validation for Owncast requirements
      try {
        const parsedUrl = new URL(normalizedUrl);

        // Only allow HTTPS URLs
        if (parsedUrl.protocol !== 'https:') {
          setError(t(Localization.Admin.FeaturedStreams.onlyHttpsSupported));
          return;
        }

        // Check for default port or no port
        if (parsedUrl.port && parsedUrl.port !== '443') {
          setError(t(Localization.Admin.FeaturedStreams.onlyDefaultPortSupported));
          return;
        }
      } catch {
        setError(t(Localization.Admin.FeaturedStreams.enterValidUrl));
        return;
      }

      setLoading(true);
      await onOk(normalizedUrl);
      form.resetFields();
      onCancel();
    } catch (err: any) {
      if (err.errorFields) {
        // Form validation error - handled by the form
        return;
      }
      setError(err.message || t(Localization.Admin.FeaturedStreams.failedToFeature));
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    setError(null);
    onCancel();
  };

  const validateUrl = (_: any, value: string) => {
    if (!value) {
      return Promise.reject(new Error(t(Localization.Admin.FeaturedStreams.enterStreamUrl)));
    }

    // Allow URLs without protocol (will be normalized later)
    let urlToValidate = value.trim();
    if (!urlToValidate.startsWith('http://') && !urlToValidate.startsWith('https://')) {
      urlToValidate = `https://${urlToValidate}`;
    }

    if (!isValidUrl(urlToValidate)) {
      return Promise.reject(new Error(t(Localization.Admin.FeaturedStreams.enterValidUrl)));
    }

    return Promise.resolve();
  };

  return (
    <Modal
      title={
        <Title level={4}>
          <Translation
            translationKey={Localization.Admin.FeaturedStreams.modalTitle}
            defaultText="Feature Live Stream"
          />
        </Title>
      }
      open={open}
      onOk={handleOk}
      onCancel={handleCancel}
      okText={
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.featureStreamAction}
          defaultText="Feature Stream"
        />
      }
      confirmLoading={loading}
      width={600}
    >
      <Form form={form} layout="vertical" autoComplete="off">
        <Form.Item
          name="url"
          label={
            <Translation
              translationKey={Localization.Admin.FeaturedStreams.streamUrlLabel}
              defaultText="Stream URL"
            />
          }
          rules={[{ validator: validateUrl }]}
          help={
            <Translation
              translationKey={Localization.Admin.FeaturedStreams.streamUrlHelp}
              defaultText="Enter the URL of an Owncast stream (e.g., stream.example.com)"
            />
          }
        >
          <Input
            placeholder={t(Localization.Admin.FeaturedStreams.streamUrlPlaceholder)}
            size="large"
          />
        </Form.Item>
      </Form>

      {error && (
        <Alert
          message={error}
          type="error"
          showIcon
          closable
          onClose={() => setError(null)}
          style={{ marginTop: 16 }}
        />
      )}

      <Alert
        message={
          <Translation
            translationKey={Localization.Admin.FeaturedStreams.streamRequirements}
            defaultText="Stream Requirements"
          />
        }
        description={
          <ul style={{ marginBottom: 0, paddingLeft: 20 }}>
            <li>
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.requirementOwncast}
                defaultText="The stream must be an Owncast instance"
              />
            </li>
            <li>
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.requirementHttps}
                defaultText="The stream must be accessible via HTTPS"
              />
            </li>
            <li>
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.requirementDefaultPort}
                defaultText="The stream must be on the default HTTPS port (443)"
              />
            </li>
            <li>
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.requirementFeaturedStreams}
                defaultText="The stream must support featured streams"
              />
            </li>
          </ul>
        }
        type="info"
        showIcon
        style={{ marginTop: 16 }}
      />
    </Modal>
  );
};
