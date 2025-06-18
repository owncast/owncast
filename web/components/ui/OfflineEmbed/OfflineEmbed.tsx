/* eslint-disable react/no-danger */
import React, { useContext, useState } from 'react';
import Head from 'next/head';
import { Typography } from 'antd';
import styles from './OfflineEmbed.module.scss';
import { Modal } from '../Modal/Modal';
import { FollowForm } from '../../modals/FollowModal/FollowForm';
import { ServerStatusContext } from '../../../utils/server-status-context';

const { Title } = Typography;

export type OfflineEmbedProps = {
  streamName: string;
  subtitle?: string;
  image: string;
  supportsFollows: boolean;
};

const Header = ({ instanceUrl, image, streamName }) =>
  instanceUrl ? (
    <a href={instanceUrl} className={styles.header}>
      <div className={styles.logo} style={{ backgroundImage: `url(${image})` }} />
      <div className={styles.title}>
        <Title level={1} className={styles.text} ellipsis={{ rows: 2 }}>
          {streamName}
        </Title>
      </div>
    </a>
  ) : (
    <div className={styles.header}>
      <div className={styles.logo} style={{ backgroundImage: `url(${image})` }} />
      <div className={styles.title}>
        <Title level={1} className={styles.text} ellipsis={{ rows: 2 }}>
          {streamName}
        </Title>
      </div>
    </div>
  );

export const OfflineEmbed = ({
  streamName,
  image,
  supportsFollows,
  subtitle,
}: OfflineEmbedProps) => {
  const [showFollowModal, setShowFollowModal] = useState(false);

  const { serverConfig } = useContext(ServerStatusContext);

  return (
    <div className={styles.canvas}>
      <Head>
        <title>{streamName}</title>
      </Head>
      <div className={styles.content}>
        <Header instanceUrl={serverConfig.yp.instanceUrl} image={image} streamName={streamName} />
        <div className={styles.body}>
          <Title level={2} className={styles.text}>
            This stream is not currently live.
          </Title>
          {subtitle && (
            <div className={styles.message}>
              <div className={styles.text}>{subtitle}</div>
            </div>
          )}
          {supportsFollows && (
            <>
              <button type="button" onClick={() => setShowFollowModal(true)}>
                Follow
              </button>
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
