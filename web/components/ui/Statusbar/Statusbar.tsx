import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import intervalToDuration from 'date-fns/intervalToDuration';
import { FC, useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import classNames from 'classnames';
import styles from './Statusbar.module.scss';
import { pluralize } from '../../../utils/helpers';

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
  const DAY_LABEL = 'day';
  const HOUR_LABEL = 'hour';
  const MINUTE_LABEL = 'minute';
  const SECOND_LABEL = 'second';
  const diff = intervalToDuration({ start: lastConnectTime, end: new Date() });

  if (diff.days >= 1) {
    return `${diff.days} ${pluralize(DAY_LABEL, diff.days)}
			${diff.hours} ${pluralize(HOUR_LABEL, diff.hours)}`;
  }
  if (diff.hours >= 1) {
    return `${diff.hours} ${pluralize(HOUR_LABEL, diff.hours)} ${diff.minutes}
			${pluralize(MINUTE_LABEL, diff.minutes)}`;
  }

  return `${diff.minutes} ${pluralize(MINUTE_LABEL, diff.minutes)}
		${diff.seconds} ${pluralize(SECOND_LABEL, diff.seconds)}`;
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
    onlineMessage = online ? `Live for  ${duration}` : 'Offline';
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
