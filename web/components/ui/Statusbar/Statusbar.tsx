import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import intervalToDuration from 'date-fns/intervalToDuration';
import { FC, useEffect, useState } from 'react';
import { EyeFilled } from '@ant-design/icons';
import styles from './Statusbar.module.scss';

export type StatusbarProps = {
  online: Boolean;
  lastConnectTime?: Date;
  lastDisconnectTime?: Date;
  viewerCount: number;
};

function makeDurationString(lastConnectTime: Date): string {
  const diff = intervalToDuration({ start: lastConnectTime, end: new Date() });
  if (diff.days > 1) {
    return `${diff.days} days ${diff.hours} hours`;
  }
  if (diff.hours >= 1) {
    return `${diff.hours} hours ${diff.minutes} minutes`;
  }

  return `${diff.minutes} minutes ${diff.seconds} seconds`;
}

export const Statusbar: FC<StatusbarProps> = ({
  online,
  lastConnectTime,
  lastDisconnectTime,
  viewerCount,
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
      <div className={styles.right}>
        <span>
          <EyeFilled />
        </span>
        <span>{` ${viewerCount}`}</span>
      </div>
    );
  } else if (!online) {
    onlineMessage = 'Offline';
    if (lastDisconnectTime) {
      rightSideMessage = `Last live ${formatDistanceToNow(new Date(lastDisconnectTime))} ago.`;
    }
  }

  return (
    <div className={styles.statusbar}>
      <div>{onlineMessage}</div>
      <div>{rightSideMessage}</div>
    </div>
  );
};
export default Statusbar;

Statusbar.defaultProps = {
  lastConnectTime: null,
  lastDisconnectTime: null,
};
