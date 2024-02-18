/* eslint-disable react/no-danger */
/* eslint-disable jsx-a11y/click-events-have-key-events */
import { FC } from 'react';
import dynamic from 'next/dynamic';
import classNames from 'classnames';
import styles from './OfflineEmbed.module.scss';

// Lazy loaded components

const ClockCircleOutlined = dynamic(() => import('@ant-design/icons/ClockCircleOutlined'), {
  ssr: false,
});

export type OfflineEmbedProps = {
  streamName: string;
  subtitle?: string;
  onFollowClick?: () => void;
};

export const OfflineEmbed: FC<OfflineEmbedProps> = ({ streamName, subtitle, onFollowClick }) => {
  return (
    <div className={classNames(styles.offlineContainer)}>
      <div className={styles.content}>
        <div className={styles.heading}>This stream is not currently live.</div>
        <div className={styles.message}>{subtitle}</div>

        <div className={styles.pageLogo}></div>
        <div className={styles.pageName}>{streamName}</div>
      </div>
    </div>
  );
};
