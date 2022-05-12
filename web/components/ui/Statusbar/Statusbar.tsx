import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import intervalToDuration from 'date-fns/intervalToDuration';
import { useEffect, useState } from 'react';
import s from './Statusbar.module.scss';

interface Props {
  online: Boolean;
  lastConnectTime?: Date;
  lastDisconnectTime?: Date;
  viewerCount: number;
}

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
export default function Statusbar(props: Props) {
  const [, setNow] = useState(new Date());

  // Set a timer to update the status bar.
  useEffect(() => {
    const interval = setInterval(() => setNow(new Date()), 1000);
    return () => {
      clearInterval(interval);
    };
  }, []);

  const { online, lastConnectTime, lastDisconnectTime, viewerCount } = props;

  let onlineMessage = '';
  let rightSideMessage = '';
  if (online && lastConnectTime) {
    const duration = makeDurationString(new Date(lastConnectTime));
    onlineMessage = online ? `Live for  ${duration}` : 'Offline';
    rightSideMessage = `${viewerCount > 0 ? `${viewerCount}` : 'No'} ${
      viewerCount === 1 ? 'viewer' : 'viewers'
    }`;
  } else {
    onlineMessage = 'Offline';
    if (lastDisconnectTime) {
      rightSideMessage = `Last live ${formatDistanceToNow(new Date(lastDisconnectTime))} ago.`;
    }
  }

  return (
    <div className={s.statusbar}>
      <div>{onlineMessage}</div>
      <div>{rightSideMessage}</div>
    </div>
  );
}

Statusbar.defaultProps = {
  lastConnectTime: null,
  lastDisconnectTime: null,
};
