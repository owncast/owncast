import { Divider } from 'antd';
import { ClockCircleOutlined } from '@ant-design/icons';
import { FC } from 'react';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import styles from './OfflineBanner.module.scss';
import { FollowButton } from '../../action-buttons/FollowButton';
import { NotifyButton } from '../../action-buttons/NotifyButton';

export type OfflineBannerProps = {
  streamName: string;
  customText?: string;
  lastLive?: Date;
  notificationsEnabled: boolean;
  fediverseAccount?: string;
  onNotifyClick?: () => void;
};

export const OfflineBanner: FC<OfflineBannerProps> = ({
  streamName,
  customText,
  lastLive,
  notificationsEnabled,
  fediverseAccount,
  onNotifyClick,
}) => {
  let text;
  if (customText) {
    text = customText;
  } else if (!customText && notificationsEnabled && fediverseAccount) {
    text = `This stream is offline. You can be notified the next time ${streamName} goes live or follow ${fediverseAccount} on the Fediverse.`;
  } else if (!customText && notificationsEnabled) {
    text = `This stream is offline. Be notified the next time ${streamName} goes live.`;
  } else {
    text = `This stream is offline. Check back soon!`;
  }

  return (
    <div className={styles.outerContainer}>
      <div className={styles.innerContainer}>
        <div className={styles.header}>{streamName}</div>
        <Divider className={styles.separator} />
        <div className={styles.bodyText}>{text}</div>
        {lastLive && (
          <div className={styles.lastLiveDate}>
            <ClockCircleOutlined className={styles.clockIcon} />
            {`Last live ${formatDistanceToNow(new Date(lastLive))} ago.`}
          </div>
        )}

        <div className={styles.footer}>
          {fediverseAccount && <FollowButton />}

          {notificationsEnabled && <NotifyButton text="Notify when live" onClick={onNotifyClick} />}
        </div>
      </div>
    </div>
  );
};
