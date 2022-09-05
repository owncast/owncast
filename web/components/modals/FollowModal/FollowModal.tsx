/* eslint-disable react/no-unescaped-entities */
import { Input, Button, Alert, Spin, Space } from 'antd';
import { useState } from 'react';
import s from './FollowModal.module.scss';

const ENDPOINT = '/api/remotefollow';

interface Props {
  handleClose: () => void;
  account: string;
  name: string;
}

function validateAccount(a) {
  const sanitized = a.replace(/^@+/, '');
  const regex =
    /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return regex.test(String(sanitized).toLowerCase());
}

export const FollowModal = (props: Props) => {
  const { handleClose, account, name } = props;
  const [remoteAccount, setRemoteAccount] = useState(null);
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);

  const handleAccountChange = a => {
    setRemoteAccount(a);
    if (validateAccount(a)) {
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
    <Space direction="vertical">
      <div className={s.header}>
        By following this stream you'll get notified on the Fediverse when it goes live. Now is a
        great time to
        <a href="https://owncast.online/join-fediverse" target="_blank" rel="noreferrer">
          learn about the Fediverse
        </a>
        if it's new to you.
      </div>

      <Spin spinning={loading}>
        {errorMessage && (
          <Alert message="Follow Error" description={errorMessage} type="error" showIcon />
        )}
        <div className={s.account}>
          <img src="/logo" alt="logo" className={s.logo} />
          <div className={s.username}>
            <div className={s.name}>{name}</div>
            <div>{account}</div>
          </div>
        </div>

        <div>
          <div className={s.instructions}>Enter your username @server to follow</div>
          <Input
            value={account}
            size="large"
            onChange={e => handleAccountChange(e.target.value)}
            placeholder="Your fediverse account @account@server"
            defaultValue={account}
          />
          <div className={s.footer}>
            You'll be redirected to your Fediverse server and asked to confirm the action.
          </div>
        </div>
        <Space className={s.buttons}>
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
export default FollowModal;
