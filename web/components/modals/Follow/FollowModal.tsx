import { Input, Button, Alert, Spin } from 'antd';
import { useState } from 'react';

const ENDPOINT = '/api/remotefollow';

interface Props {
  handleClose: () => void;
}

function validateAccount(a) {
  const sanitized = a.replace(/^@+/, '');
  const regex =
    /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return regex.test(String(sanitized).toLowerCase());
}

export default function FollowModal(props: Props) {
  const { handleClose } = props;
  const [account, setAccount] = useState(null);
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);

  const handleAccountChange = a => {
    setAccount(a);
    if (validateAccount(a)) {
      setValid(true);
    } else {
      setValid(false);
    }
  };

  const remoteFollowButtonPressed = async () => {
    if (!valid) {
      return;
    }

    setLoading(true);

    try {
      const sanitizedAccount = account.replace(/^@+/, '');
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
        <Alert message="Follow Error" description={errorMessage} type="error" showIcon />
      )}
      <Input
        value={account}
        size="large"
        onChange={e => handleAccountChange(e.target.value)}
        placeholder="Your fediverse account @account@server"
        defaultValue={account}
      />
      <Button disabled={!valid} onClick={remoteFollowButtonPressed}>
        Follow
      </Button>
      <div>
        Information about following a Fediverse account and next steps how to create a Fediverse
        account goes here.
      </div>
    </Spin>
  );
}
