import { Alert, Button, Input, Space, Spin, Collapse } from 'antd';
import React, { FC, useState } from 'react';
import dynamic from 'next/dynamic';
import styles from './FediAuthModal.module.scss';
import { isValidFediverseAccount } from '../../../utils/validators';

const { Panel } = Collapse;

// Lazy loaded components

const CheckCircleOutlined = dynamic(() => import('@ant-design/icons/CheckCircleOutlined'), {
  ssr: false,
});

export type FediAuthModalProps = {
  authenticated: boolean;
  displayName: string;
  accessToken: string;
};

export const FediAuthModal: FC<FediAuthModalProps> = ({
  authenticated,
  displayName,
  accessToken,
}) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [valid, setValid] = useState(false);
  const [account, setAccount] = useState('');
  const [code, setCode] = useState('');
  const [verifyingCode, setVerifyingCode] = useState(false);

  const message = !authenticated ? (
    <span>
      Receive a direct message on the Fediverse to link your account to{' '}
      <strong>{displayName}</strong>, or login as a previously linked chat user.
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

  const validate = (acct: string) => {
    setValid(isValidFediverseAccount(acct));
  };

  const onInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setAccount(e.target.value);
    validate(e.target.value);
  };

  const makeRequest = async (url, data) => {
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
    }
  };

  const submitCodePressed = async () => {
    setLoading(true);
    const url = `/api/auth/fediverse/verify?accessToken=${accessToken}`;
    const data = { code };

    try {
      await makeRequest(url, data);

      // Success. Reload the page.
      window.location.href = '/';
    } catch (e) {
      console.error(e);
      setErrorMessage(e);
    }
    setLoading(false);
  };

  const submitAccountPressed = async () => {
    if (!valid) {
      return;
    }

    setLoading(true);
    setErrorMessage(null);
    const url = `/api/auth/fediverse?accessToken=${accessToken}`;
    const normalizedAccount = account.replace(/^@+/, '');
    const data = { account: normalizedAccount };

    try {
      await makeRequest(url, data);
      setVerifyingCode(true);
    } catch (e) {
      console.error(e);
      setErrorMessage(e);
    }
    setLoading(false);
  };

  const inputCodeStep = (
    <div>
      Paste in the code that was sent to your Fediverse account. If you did not receive a code, make
      sure you can accept direct messages.
      <div className={styles.codeInputContainer}>
        <Input
          value={code}
          onChange={e => setCode(e.target.value)}
          className={styles.codeInput}
          placeholder="123456"
          maxLength={6}
        />
        <Button
          type="primary"
          onClick={submitCodePressed}
          disabled={code.length < 6}
          className={styles.submitButton}
        >
          Verify Code
        </Button>
      </div>
    </div>
  );

  const inputAccountStep = (
    <>
      <div>Your Fediverse Account</div>
      <Input.Search
        addonBefore="@"
        onInput={onInput}
        value={account}
        placeholder="youraccount@yourserver.com"
        status={!valid && account.length > 0 ? 'error' : undefined}
        onSearch={submitAccountPressed}
        enterButton={
          <Button type={valid ? 'primary' : 'default'} disabled={!valid || account.length === 0}>
            <CheckCircleOutlined />
          </Button>
        }
      />
    </>
  );

  return (
    <Spin spinning={loading}>
      <Space direction="vertical">
        {message}
        {errorMessageText && (
          <Alert message="Error" description={errorMessageText} type="error" showIcon />
        )}
        {verifyingCode ? inputCodeStep : inputAccountStep}
        <Collapse ghost>
          <Panel
            key="header"
            header="Learn more about using the Fediverse to authenticate with chat."
          >
            <p>
              You can link your chat identity with your Fediverse identity. Next time you want to
              use this chat identity you can again go through the Fediverse authentication.
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
