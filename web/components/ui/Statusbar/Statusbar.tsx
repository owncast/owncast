import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import formatDistanceToNowStrict from 'date-fns/formatDistanceToNowStrict';
import { useEffect, useState } from 'react';

import s from './Statusbar.module.scss';

interface Props {
  online: Boolean;
  lastConnectTime?: Date;
  lastDisconnectTime?: Date;
  viewerCount: number;
}
export default function Statusbar(props: Props) {
  const [now, setNow] = useState(new Date());

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
    const diff = formatDistanceToNowStrict(new Date(lastConnectTime));
    onlineMessage = online ? `Streaming for  ${diff}` : 'Offline';
    rightSideMessage = `${viewerCount} viewers`;
  } else {
    onlineMessage = 'Offline';
    if (lastDisconnectTime) {
      rightSideMessage = `Last live ${formatDistanceToNow(lastDisconnectTime)} ago.`;
    }
  }

  return (
    <div className={s.statusbar}>
      {/* <div>{streamStartedAt}</div> */}
      <div>{onlineMessage}</div>
      <div>{rightSideMessage}</div>
    </div>
  );
}

Statusbar.defaultProps = {
  lastConnectTime: null,
  lastDisconectTime: null,
};
