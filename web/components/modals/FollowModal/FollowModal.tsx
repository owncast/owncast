/* eslint-disable react/no-unescaped-entities */
import { Input, Button, Alert, Spin, Space } from 'antd';
import { FC, useState } from 'react';
import styles from './FollowModal.module.scss';
import { isValidFediverseAccount } from '../../../utils/validators';

const ENDPOINT = '/api/remotefollow';

export type FollowModalProps = {
  handleClose: () => void;
  account: string;
  name: string;
};

export const FollowModal: FC<FollowModalProps> = ({ handleClose, account, name }) => {
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
    <Space direction="vertical" id="follow-modal">
      <div className={styles.header}>
        By following this stream you'll get notified on the Fediverse when it goes live. Now is a
        great time to
        <a href="https://owncast.online/join-fediverse" target="_blank" rel="noreferrer">
          &nbsp;learn about the Fediverse&nbsp;
        </a>
        if it's new to you.
      </div>

      <Spin spinning={loading}>
        {errorMessage && (
          <Alert message="Follow Error" description={errorMessage} type="error" showIcon />
        )}
        <div className={styles.account}>
          <img src="/logo" alt="logo" className={styles.logo} />
          <div className={styles.username}>
            <div className={styles.name}>{name}</div>
            <div>{account}</div>
          </div>
        </div>

        <div>
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
          <Button disabled={!valid} type="primary" onClick={remoteFollowButtonPressed}>
            Follow
          </Button>
          <Button onClick={joinButtonPressed} type="primary">
            Join the Fediverse
          </Button>
        </Space>
      </Spin>
    </Space>
  );
};
