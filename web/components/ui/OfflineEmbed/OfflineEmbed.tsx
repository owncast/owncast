/* eslint-disable react/no-danger */
import React, { FC, useCallback, useContext, useState } from 'react';
import Head from 'next/head';
import { Button, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
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

export const OfflineEmbed: FC<OfflineEmbedProps> = ({
  streamName,
  image,
  supportsFollows,
  subtitle,
}: OfflineEmbedProps) => {
  const [showFollowModal, setShowFollowModal] = useState(false);
  const {
    serverConfig: {
      yp: { instanceUrl },
    },
  } = useContext(ServerStatusContext);
  const { t } = useTranslation();

  const Linkable = useCallback(
    ({ href, children, ...props }) =>
      href ? (
        <a href={href} {...props}>
          {children}
        </a>
      ) : (
        <div {...props}>{children}</div>
      ),
    [],
  );

  return (
    <div className={styles.canvas}>
      <Head>
        <title>{streamName}</title>
      </Head>
      <div className={styles.content}>
        <Linkable href={instanceUrl} className={styles.header}>
          <div className={styles.logo} style={{ backgroundImage: `url(${image})` }} />
          <div className={styles.title}>
            <Title level={1} className={styles.text} ellipsis={{ rows: 2 }}>
              {streamName}
            </Title>
          </div>
        </Linkable>
        <div className={styles.body}>
          <Title level={2} className={styles.text}>
            {t('This stream is offline')}
          </Title>
          {subtitle && (
            <div className={styles.message}>
              <div className={styles.text}>{subtitle}</div>
            </div>
          )}
          {supportsFollows && (
            <>
              <Button size="small" onClick={() => setShowFollowModal(true)}>
                {t('Follow')}
              </Button>
              <Modal
                title={`${t('Follow')} ${streamName}`}
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
