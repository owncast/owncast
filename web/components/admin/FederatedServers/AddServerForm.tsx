import React, { FC, useState } from 'react';
import { Form, Input, Button, Alert } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import styles from './AddServerForm.module.scss';

export interface AddServerFormProps {
  onAdd: (url: string) => Promise<void>;
  disabled?: boolean;
}

export const AddServerForm: FC<AddServerFormProps> = ({ onAdd, disabled = false }) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const validateUrl = (url: string): boolean => {
    try {
      const parsedUrl = new URL(url);
      // Only allow HTTPS URLs
      if (parsedUrl.protocol !== 'https:') {
        setError('Only HTTPS URLs are supported for federation');
        return false;
      }
      // Check for default port or no port
      if (parsedUrl.port && parsedUrl.port !== '443') {
        setError('Only streams on the default HTTPS port (443) are supported');
        return false;
      }
      return true;
    } catch {
      setError('Please enter a valid URL');
      return false;
    }
  };

  const handleSubmit = async (values: { url: string }) => {
    setError(null);

    // Normalize the URL
    let normalizedUrl = values.url.trim();
    if (!normalizedUrl.startsWith('http://') && !normalizedUrl.startsWith('https://')) {
      normalizedUrl = `https://${normalizedUrl}`;
    }

    if (!validateUrl(normalizedUrl)) {
      return;
    }

    setLoading(true);
    try {
      await onAdd(normalizedUrl);
      form.resetFields();
    } catch (err: any) {
      setError(err.message || 'Failed to feature stream');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Form form={form} layout="inline" onFinish={handleSubmit} className={styles.form}>
        <Form.Item
          name="url"
          rules={[{ required: true, message: 'Please enter a stream URL' }]}
          className={styles.urlInput}
        >
          <Input
            placeholder="Enter Owncast stream URL (e.g., stream.example.com)"
            disabled={disabled || loading}
            size="large"
          />
        </Form.Item>
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            icon={<PlusOutlined />}
            loading={loading}
            disabled={disabled}
            size="large"
          >
            Feature Stream
          </Button>
        </Form.Item>
      </Form>
      {error && (
        <Alert
          message={error}
          type="error"
          showIcon
          closable
          onClose={() => setError(null)}
          className={styles.error}
        />
      )}
      <Alert
        message="Stream Requirements"
        description={
          <ul className={styles.requirements}>
            <li>The stream must be an Owncast instance</li>
            <li>The stream must be accessible via HTTPS</li>
            <li>The stream must be on the default HTTPS port (443)</li>
            <li>The stream must have federation features enabled</li>
          </ul>
        }
        type="info"
        showIcon
        className={styles.info}
      />
    </div>
  );
};
