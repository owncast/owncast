/* eslint-disable react/no-unescaped-entities */
import { Input, Button, Alert, Spin, Space } from 'antd';
import { FC, useState } from 'react';
import styles from './FollowModal.module.scss';
import { isValidFediverseAccount } from '../../../utils/validators';

const ENDPOINT = '/api/remotefollow';

export type FollowFormProps = {
  handleClose?: () => void;
};

export const FollowForm: FC<FollowFormProps> = ({ handleClose }: FollowFormProps) => {
  const [remoteAccount, setRemoteAccount] = useState(null);
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);

  const handleAccountChange = a => {
    setRemoteAccount(a);
    if (isValidFediverseAccount(a)) {
      setValid(true);
    } else {
      setValid(false);
    }
  };

  const joinButtonPressed = () => {
    window.open('https://owncast.online/join-fediverse', '_blank');
  };

  const remoteFollowButtonPressed = async () => {
    if (!valid) {
      return;
    }

    setLoading(true);

    try {
      const sanitizedAccount = remoteAccount.replace(/^@+/, '');
      const request = { account: sanitizedAccount };
      const rawResponse = await fetch(ENDPOINT, {
        method: 'POST',
        body: JSON.stringify(request),
      });
      const result = await rawResponse.json();

      if (result.redirectUrl) {
        window.open(result.redirectUrl, '_blank');
        handleClose();
      }
      if (!result.success) {
        setErrorMessage(result.message);
        setLoading(false);
        return;
      }
      if (!result.redirectUrl) {
        setErrorMessage('Unable to follow.');
        setLoading(false);
        return;
      }
    } catch (e) {
      setErrorMessage(e.message);
    }
    setLoading(false);
  };

  return (
    <Spin spinning={loading}>
      {errorMessage && (
        <Alert
          message="Follow Error"
          description={errorMessage}
          type="error"
          closable
          className={styles.errorAlert}
        />
      )}

      <div className={styles.inputContainer}>
        <div className={styles.instructions}>Enter your username @server to follow</div>
        <Input
          value={remoteAccount}
          size="large"
          onChange={e => handleAccountChange(e.target.value)}
          placeholder="Your fediverse account @account@server"
          defaultValue={remoteAccount}
        />
        <div className={styles.footer}>
          You'll be redirected to your Fediverse server and asked to confirm the action.
        </div>
      </div>
      <Space className={styles.buttons}>
        <Button onClick={joinButtonPressed} type="text">
          Join the Fediverse
        </Button>
        <Button disabled={!valid} type="primary" onClick={remoteFollowButtonPressed}>
          Follow
        </Button>
      </Space>
    </Spin>
  );
};
