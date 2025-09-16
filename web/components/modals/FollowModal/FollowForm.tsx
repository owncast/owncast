/* eslint-disable react/no-unescaped-entities */
import { FC, useState } from 'react';
import { Input, Button, Alert, Spin, Space } from 'antd';
import { useTranslation } from 'next-export-i18n';
import styles from './FollowModal.module.scss';
import { isValidFediverseAccount } from '../../../utils/validators';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';

const ENDPOINT = '/api/remotefollow';

export type FollowFormProps = {
  handleClose?: () => void;
};

export const FollowForm: FC<FollowFormProps> = ({ handleClose }: FollowFormProps) => {
  const [remoteAccount, setRemoteAccount] = useState(null);
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const { t } = useTranslation();

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
        setErrorMessage(t(Localization.Frontend.FollowModal.unableToFollow) || 'Unable to follow.');
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
          message={
            <Translation
              translationKey={Localization.Frontend.FollowModal.followError}
              defaultText="Follow Error"
            />
          }
          description={errorMessage}
          type="error"
          closable
          className={styles.errorAlert}
        />
      )}

      <div className={styles.inputContainer}>
        <div className={styles.instructions}>
          <Translation
            translationKey={Localization.Frontend.FollowModal.instructions}
            defaultText="Enter your username @server to follow"
          />
        </div>
        <Input
          value={remoteAccount}
          size="large"
          onChange={e => handleAccountChange(e.target.value)}
          placeholder={
            t(Localization.Frontend.FollowModal.placeholder) ||
            'Your fediverse account @account@server'
          }
          defaultValue={remoteAccount}
        />
        <div className={styles.footer}>
          <Translation
            translationKey={Localization.Frontend.FollowModal.redirectMessage}
            defaultText="You'll be redirected to your Fediverse server and asked to confirm the action."
          />
        </div>
      </div>
      <Space className={styles.buttons}>
        <Button onClick={joinButtonPressed} type="text">
          <Translation
            translationKey={Localization.Frontend.FollowModal.joinFediverse}
            defaultText="Join the Fediverse"
          />
        </Button>
        <Button disabled={!valid} type="primary" onClick={remoteFollowButtonPressed}>
          <Translation
            translationKey={Localization.Frontend.FollowModal.followButton}
            defaultText="Follow"
          />
        </Button>
      </Space>
    </Spin>
  );
};
