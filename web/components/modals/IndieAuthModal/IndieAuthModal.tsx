import { Alert, Input, Space, Spin, Collapse, Typography, Button } from 'antd';
import dynamic from 'next/dynamic';
import React, { FC, useState } from 'react';
import { isValidUrl } from '../../../utils/validators';

const { Panel } = Collapse;
const { Link } = Typography;

// Lazy loaded components

const CheckCircleOutlined = dynamic(() => import('@ant-design/icons/CheckCircleOutlined'), {
  ssr: false,
});

export type IndieAuthModalProps = {
  authenticated: boolean;
  displayName: string;
  accessToken: string;
};

export const IndieAuthModal: FC<IndieAuthModalProps> = ({
  authenticated,
  displayName: username,
  accessToken,
}) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [valid, setValid] = useState(false);
  const [host, setHost] = useState('');

  const message = !authenticated ? (
    <span>
      Use your own domain to authenticate <span>{username}</span> or login as a previously{' '}
      authenticated chat user using IndieAuth.
    </span>
  ) : (
    <span>
      <b>You are already authenticated</b>. However, you can add other domains or log in as a
      different user.
    </span>
  );

  let errorMessageText = errorMessage;
  if (errorMessageText) {
    if (errorMessageText.includes('url does not support indieauth')) {
      errorMessageText = 'The provided URL is either invalid or does not support IndieAuth.';
    }
  }

  const validate = (url: string) => {
    if (!isValidUrl(url)) {
      setValid(false);
      return;
    }

    if (!url.includes('.')) {
      setValid(false);
      return;
    }

    setValid(true);
  };

  const onInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Don't allow people to type custom ports or protocols.
    const char = (e.nativeEvent as any).data;
    if (char === ':') {
      return;
    }

    setHost(e.target.value);
    const h = `https://${e.target.value}`;
    validate(h);
  };

  const submitButtonPressed = async () => {
    if (!valid) {
      return;
    }

    setLoading(true);

    try {
      const url = `/api/auth/indieauth?accessToken=${accessToken}`;
      const h = `https://${host}`;
      const data = { authHost: h };
      const rawResponse = await fetch(url, {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });

      const content = await rawResponse.json();
      if (content.message) {
        setErrorMessage(content.message);
        setLoading(false);
        return;
      }
      if (!content.redirect) {
        setErrorMessage('Auth provider did not return a redirect URL.');
        setLoading(false);
        return;
      }

      if (content.redirect) {
        const { redirect } = content;
        window.location = redirect;
      }
    } catch (e) {
      setErrorMessage(e.message);
    }

    setLoading(false);
  };

  return (
    <Spin spinning={loading}>
      <Space direction="vertical">
        {message}
        {errorMessageText && (
          <Alert message="Error" description={errorMessageText} type="error" showIcon />
        )}
        <div>Your domain</div>
        <Input.Search
          addonBefore="https://"
          onInput={onInput}
          type="url"
          value={host}
          placeholder="yoursite.com"
          status={!valid && host.length > 0 ? 'error' : undefined}
          onSearch={submitButtonPressed}
          enterButton={
            <Button type={valid ? 'primary' : 'default'} disabled={!valid || host.length === 0}>
              <CheckCircleOutlined />
            </Button>
          }
        />

        <Collapse ghost>
          <Panel key="header" header="Learn more about using IndieAuth to authenticate with chat.">
            <p>
              IndieAuth allows for a completely independent and decentralized way of identifying
              yourself using your own domain.
            </p>

            <p>
              If you run an Owncast instance, you can use that domain here. Otherwise,{' '}
              <Link href="https://indieauth.net/#providers">
                learn more about how you can support IndieAuth
              </Link>
              .
            </p>
          </Panel>
        </Collapse>
        <div>
          <strong>Note</strong>: This is for authentication purposes only, and no personal
          information will be accessed or stored.
        </div>
      </Space>
    </Spin>
  );
};
