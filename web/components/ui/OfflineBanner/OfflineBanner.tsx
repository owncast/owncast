/* eslint-disable react/no-danger */
/* eslint-disable jsx-a11y/click-events-have-key-events */
import { Divider } from 'antd';
import React, { FC } from 'react';
import { formatDistanceToNow } from 'date-fns';
import dynamic from 'next/dynamic';
import classNames from 'classnames';
import { useTranslation } from 'next-export-i18n';
import { Translation } from '../Translation/Translation';
import { Localization } from '../../../types/localization';
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

  const handleSpanClick = (event: React.MouseEvent<HTMLSpanElement>) => {
    const target = event.target as HTMLSpanElement;
    if (target.classList.contains('notify-link')) {
      onNotifyClick?.();
    } else if (target.classList.contains('follow-link')) {
      onFollowClick?.();
    }
  };

  let text;
  if (customText) {
    text = customText;
  } else if (!customText && notificationsEnabled && fediverseAccount) {
    text = (
      <Translation
        translationKey={Localization.Frontend.offlineNotifyAndFediverse}
        vars={{ streamer: streamName, fediverseAccount }}
        defaultText="This stream is offline. You can <span class='notify-link'>be notified</span> the next time {{streamer}} goes live or <span class='follow-link'>follow</span> {{fediverseAccount}} on the Fediverse."
      />
    );
  } else if (!customText && notificationsEnabled) {
    text = (
      <Translation
        translationKey={Localization.Frontend.offlineNotifyOnly}
        vars={{ streamer: streamName }}
        defaultText="This stream is offline. <span class='notify-link'>Be notified</span> the next time {{streamer}} goes live."
      />
    );
  } else if (!customText && fediverseAccount) {
    text = (
      <Translation
        translationKey={Localization.Frontend.offlineFediverseOnly}
        vars={{ fediverseAccount, streamer: streamName }}
        defaultText="This stream is offline. <span class='follow-link'>Follow</span> {{fediverseAccount}} on the Fediverse to see the next time {{streamer}} goes live."
      />
    );
  } else {
    text = (
      <Translation
        translationKey={Localization.Frontend.offlineBasic}
        defaultText="This stream is offline. Check back soon!"
      />
    );
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
          <div
            className={styles.bodyText}
            onClick={handleSpanClick}
            role="presentation"
            style={{ cursor: 'pointer' }}
          >
            {text}
          </div>
        )}

        {lastLive && (
          <div className={styles.lastLiveDate}>
            <ClockCircleOutlined className={styles.clockIcon} />
            <span id="owncast-offline-last-live-text">
              {`${t(Localization.Frontend.lastLiveAgo, { timeAgo: formatDistanceToNow(new Date(lastLive)) })}`}
            </span>
          </div>
        )}
      </div>
    </div>
  );
};
