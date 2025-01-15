/* eslint-disable react/no-danger */
/* eslint-disable jsx-a11y/click-events-have-key-events */
import { Divider } from 'antd';
import { FC } from 'react';
import { formatDistanceToNow } from 'date-fns';
import dynamic from 'next/dynamic';
import classNames from 'classnames';
import { useTranslation } from 'next-export-i18n';
import styles from './OfflineBanner.module.scss';

// Lazy loaded components

const ClockCircleOutlined = dynamic(() => import('@ant-design/icons/ClockCircleOutlined'), {
  ssr: false,
});

export type OfflineBannerProps = {
  streamName: string;
  customText?: string;
  lastLive?: Date;
  notificationsEnabled: boolean;
  fediverseAccount?: string;
  showsHeader?: boolean;
  onNotifyClick?: () => void;
  onFollowClick?: () => void;
  className?: string;
};

export const OfflineBanner: FC<OfflineBannerProps> = ({
  streamName,
  customText,
  lastLive,
  notificationsEnabled,
  fediverseAccount,
  showsHeader = true,
  onNotifyClick,
  onFollowClick,
  className,
}) => {
  const { t } = useTranslation();

  let text;
  if (customText) {
    text = customText;
  } else if (!customText && notificationsEnabled && fediverseAccount) {
    text = (
      <span>
        {t('This stream is offline. You can')}{' '}
        <span role="link" tabIndex={0} className={styles.actionLink} onClick={onNotifyClick}>
          be notified
        </span>{' '}
        the next time {streamName} goes live or{' '}
        <span role="link" tabIndex={0} className={styles.actionLink} onClick={onFollowClick}>
          follow
        </span>{' '}
        {fediverseAccount} on the Fediverse.
      </span>
    );
  } else if (!customText && notificationsEnabled) {
    text = (
      <span>
        {t('This stream is offline')}.{' '}
        <span role="link" tabIndex={0} className={styles.actionLink} onClick={onNotifyClick}>
          Be notified
        </span>{' '}
        {t('the next time goes live', { streamer: streamName })}.
      </span>
    );
  } else if (!customText && fediverseAccount) {
    text = (
      <span>
        {t('This stream is offline.')}{' '}
        <span role="link" tabIndex={0} className={styles.actionLink} onClick={onFollowClick}>
          {t('Follow')}
        </span>{' '}
        {t('on the Fediverse to see the next time goes live', {
          fediverseAccount,
          streamer: streamName,
        })}
        .
      </span>
    );
  } else {
    text = `This stream is offline. Check back soon!`;
  }

  return (
    <div id="offline-banner" className={classNames(styles.outerContainer, className)}>
      <div className={styles.innerContainer}>
        {showsHeader && (
          <>
            <div className={styles.header}>{streamName}</div>
            <Divider className={styles.separator} />
          </>
        )}
        {customText ? (
          <div className={styles.bodyText} dangerouslySetInnerHTML={{ __html: text }} />
        ) : (
          <div className={styles.bodyText}>{text}</div>
        )}

        {lastLive && (
          <div className={styles.lastLiveDate}>
            <ClockCircleOutlined className={styles.clockIcon} />
            {`${t('Last live ago', { timeAgo: formatDistanceToNow(new Date(lastLive)) })}`}
          </div>
        )}
      </div>
    </div>
  );
};
