/* eslint-disable react/no-danger */

import { FC, useEffect, useState } from 'react';
import classNames from 'classnames';
import Head from 'next/head';
import { Button, Typography } from 'antd';
import styles from './OfflineEmbed.module.scss';
import { Modal } from '../Modal/Modal';
import { FollowForm } from '../../modals/FollowModal/FollowForm';

const { Title } = Typography;

export type OfflineEmbedProps = {
  streamName: string;
  subtitle?: string;
  image: string;
  supportsFollows: boolean;
};

enum EmbedMode {
  CannotFollow = 1,
  CanFollow,
}

export const OfflineEmbed: FC<OfflineEmbedProps> = ({
  streamName,
  subtitle,
  image,
  supportsFollows,
}) => {
  const [currentMode, setCurrentMode] = useState(EmbedMode.CanFollow);
  const [showFollowModal, setShowFollowModal] = useState(false);

  useEffect(() => {
    if (!supportsFollows) {
      setCurrentMode(EmbedMode.CannotFollow);
    } else if (currentMode === EmbedMode.CannotFollow) {
      setCurrentMode(EmbedMode.CanFollow);
    }
  }, [supportsFollows]);

  const followButtonPressed = async () => {
    setShowFollowModal(true);
  };

  return (
    <div>
      <Head>
        <title>{streamName}</title>
      </Head>
      <div className={classNames(styles.offlineContainer)}>
        <div className={styles.content}>
          <Title level={1} className={styles.headerContainer}>
            <div className={styles.pageLogo} style={{ backgroundImage: `url(${image})` }} />
            <div className={styles.streamName}>{streamName}</div>
          </Title>

          <div className={styles.messageContainer}>
            <Title level={2} className={styles.offlineTitle}>
              This stream is not currently live.
            </Title>
            <div className={styles.message} dangerouslySetInnerHTML={{ __html: subtitle }} />
          </div>

          {currentMode === EmbedMode.CanFollow && (
            <>
              <Button className={styles.followButton} type="primary" onClick={followButtonPressed}>
                Follow Server
              </Button>
              <Modal
                title={`Follow ${streamName}`}
                open={showFollowModal}
                handleCancel={() => setShowFollowModal(false)}
              >
                <FollowForm />
              </Modal>
            </>
          )}
        </div>
      </div>
    </div>
  );
};
