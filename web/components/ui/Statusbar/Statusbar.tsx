import { intervalToDuration, formatDistanceToNow, formatDuration } from 'date-fns';
import { FC, useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import classNames from 'classnames';
import styles from './Statusbar.module.scss';

// Lazy loaded components

const EyeFilled = dynamic(() => import('@ant-design/icons/EyeFilled'), {
  ssr: false,
});

export type StatusbarProps = {
  online: Boolean;
  lastConnectTime?: Date;
  lastDisconnectTime?: Date;
  viewerCount: number;
  className?: string;
};

function makeDurationString(lastConnectTime: Date): string {
  const diff = intervalToDuration({ start: lastConnectTime, end: new Date() });
  if (diff.days >= 1) {
    return formatDuration({
      days: diff.days,
      hours: diff.hours > 0 ? diff.hours : 0,
    });
  }
  if (diff.hours >= 1) {
    return formatDuration({
      hours: diff.hours,
      minutes: diff.minutes > 0 ? diff.minutes : 0,
    });
  }
  return formatDuration({
    minutes: diff.minutes > 0 ? diff.minutes : 0,
    seconds: diff.seconds > 0 ? diff.seconds : 0,
  });
}

export const Statusbar: FC<StatusbarProps> = ({
  online,
  lastConnectTime,
  lastDisconnectTime,
  viewerCount,
  className,
}) => {
  const [, setNow] = useState(new Date());

  // Set a timer to update the status bar.
  useEffect(() => {
    const interval = setInterval(() => setNow(new Date()), 1000);
    return () => {
      clearInterval(interval);
    };
  }, []);

  let onlineMessage = '';
  let rightSideMessage: any;
  if (online && lastConnectTime) {
    const duration = makeDurationString(new Date(lastConnectTime));
    onlineMessage = `Live for  ${duration}`;
    rightSideMessage = viewerCount > 0 && (
      <>
        <span className={styles.viewerIcon}>
          <EyeFilled />
        </span>
        <span>{` ${viewerCount}`}</span>
      </>
    );
  } else if (!online) {
    onlineMessage = 'Offline';
    if (lastDisconnectTime) {
      rightSideMessage = `Last live ${formatDistanceToNow(new Date(lastDisconnectTime))} ago.`;
    }
  }

  return (
    <div className={classNames(styles.statusbar, className)} role="status">
      <span aria-live="off" className={styles.onlineMessage}>
        {onlineMessage}
      </span>
      <span className={styles.viewerCount}>{rightSideMessage}</span>
    </div>
  );
};

Statusbar.defaultProps = {
  lastConnectTime: null,
  lastDisconnectTime: null,
};
