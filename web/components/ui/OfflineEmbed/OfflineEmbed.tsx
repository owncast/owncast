/* eslint-disable react/no-danger */

import { FC, useEffect, useState } from 'react';
import classNames from 'classnames';
import Head from 'next/head';
import { Button, Input, Space, Spin, Alert, Typography } from 'antd';
import styles from './OfflineEmbed.module.scss';
import { isValidFediverseAccount } from '../../../utils/validators';

const { Title } = Typography;
const ENDPOINT = '/api/remotefollow';

export type OfflineEmbedProps = {
  streamName: string;
  subtitle?: string;
  image: string;
  supportsFollows: boolean;
};

enum EmbedMode {
  CannotFollow = 1,
  CanFollow,
  FollowPrompt,
  InProgress,
}

export const OfflineEmbed: FC<OfflineEmbedProps> = ({
  streamName,
  subtitle,
  image,
  supportsFollows,
}) => {
  const [currentMode, setCurrentMode] = useState(EmbedMode.CanFollow);
  const [remoteAccount, setRemoteAccount] = useState(null);
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);

  useEffect(() => {
    if (!supportsFollows) {
      setCurrentMode(EmbedMode.CannotFollow);
    }
  }, [supportsFollows]);

  const followButtonPressed = async () => {
    setCurrentMode(EmbedMode.FollowPrompt);
  };

  const remoteFollowButtonPressed = async () => {
    setLoading(true);
    setCurrentMode(EmbedMode.CannotFollow);

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

  const handleAccountChange = a => {
    setRemoteAccount(a);
    if (isValidFediverseAccount(a)) {
      setValid(true);
    } else {
      setValid(false);
    }
  };

  return (
    <div>
      <Head>
        <title>{streamName}</title>
      </Head>
      <div className={classNames(styles.offlineContainer)}>
        <Spin spinning={loading}>
          <div className={styles.content}>
            <div className={styles.heading}>This stream is not currently live.</div>
            <div className={styles.message} dangerouslySetInnerHTML={{ __html: subtitle }} />

            <div className={styles.pageLogo} style={{ backgroundImage: `url(${image})` }} />
            <div className={styles.pageName}>{streamName}</div>

            {errorMessage && (
              <Alert message="Follow Error" description={errorMessage} type="error" showIcon />
            )}

            {currentMode === EmbedMode.CanFollow && (
              <Button className={styles.submitButton} type="primary" onClick={followButtonPressed}>
                Follow Server
              </Button>
            )}

            {currentMode === EmbedMode.InProgress && (
              <Title level={4} className={styles.heading}>
                Follow the instructions on your Fediverse server to complete the follow.
              </Title>
            )}

            {currentMode === EmbedMode.FollowPrompt && (
              <div>
                <Input
                  value={remoteAccount}
                  size="large"
                  onChange={e => handleAccountChange(e.target.value)}
                  placeholder="Your fediverse account @account@server"
                  defaultValue={remoteAccount}
                />
                <div className={styles.footer}>
                  You&apos;ll be redirected to your Fediverse server and asked to confirm the
                  action.
                </div>
                <Space className={styles.buttons}>
                  <Button
                    className={styles.submitButton}
                    disabled={!valid}
                    type="primary"
                    onClick={remoteFollowButtonPressed}
                  >
                    Submit and Follow
                  </Button>
                </Space>
              </div>
            )}
          </div>
        </Spin>
      </div>
    </div>
  );
};
