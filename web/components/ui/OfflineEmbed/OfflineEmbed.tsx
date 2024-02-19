import { FC } from 'react';
import classNames from 'classnames';
import Head from 'next/head';
import styles from './OfflineEmbed.module.scss';

export type OfflineEmbedProps = {
  streamName: string;
  subtitle?: string;
  image: string;
};

export const OfflineEmbed: FC<OfflineEmbedProps> = ({ streamName, subtitle, image }) => (
  <div>
    <Head>
      <title>{streamName}</title>
    </Head>
    <div className={classNames(styles.offlineContainer)}>
      <div className={styles.content}>
        <div className={styles.heading}>This stream is not currently live.</div>
        <div className={styles.message}>{subtitle}</div>

        <div className={styles.pageLogo} style={{ backgroundImage: `url(${image})` }} />
        <div className={styles.pageName}>{streamName}</div>
      </div>
    </div>
  </div>
);
