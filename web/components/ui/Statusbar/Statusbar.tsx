import { intervalToDuration, formatDistanceToNow, formatDuration, Locale } from 'date-fns';
import { FC, useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import classNames from 'classnames';
import styles from './Statusbar.module.scss';
import {Translation} from '../Translation/Translation';
import {Localization} from '../../../types/localization';
import {useDateFnsLocale} from '../../../hooks/useDateFnsLocale';




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

function makeDurationString(lastConnectTime: Date, locale?: Locale): string {
  const diff = intervalToDuration({ start: lastConnectTime, end: new Date() });
  if (diff.days >= 1) {
    return formatDuration(
      {
        days: diff.days,
        hours: diff.hours > 0 ? diff.hours : 0,
      },
      { locale },
    );
  }
  if (diff.hours >= 1) {
    return formatDuration(
      {
        hours: diff.hours,
        minutes: diff.minutes > 0 ? diff.minutes : 0,
      },
      { locale },
    );
  }
  return formatDuration(
    {
      minutes: diff.minutes > 0 ? diff.minutes : 0,
      seconds: diff.seconds > 0 ? diff.seconds : 0,
    },
    { locale },
  );
}

export const Statusbar: FC<StatusbarProps> = ({
  online,
  lastConnectTime = null,
  lastDisconnectTime = null,
  viewerCount,
  className,
}) => {
  const dateFnsLocale = useDateFnsLocale();
  const [, setNow] = useState(new Date());
  // Set a timer to update the status bar.
  useEffect(() => {
    const interval = setInterval(() => setNow(new Date()), 1000);
    return () => {
      clearInterval(interval);
    };
  }, []);

  let onlineMessage: any;
  let rightSideMessage: any;
  if (online && lastConnectTime) {
    const duration = makeDurationString(new Date(lastConnectTime), dateFnsLocale);
    onlineMessage = (
      <Translation
        translationKey={Localization.Frontend.liveNow}
        vars={{timeSince: duration}}
        defaultText="Live for {{timeSince}}"
      />
    );
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
      rightSideMessage = (
        <Translation
          translationKey={Localization.Frontend.lastLiveAgo}
          vars={{ timeAgo: formatDistanceToNow(new Date(lastDisconnectTime), { locale: dateFnsLocale }) }}
          defaultText="Last live {{timeAgo}} ago."
        />
      );
    }
  }

  return (
    <div className={classNames(styles.statusbar, className)}>
      <span className={styles.onlineMessage}>{onlineMessage}</span>
      <span className={styles.viewerCount}>{rightSideMessage}</span>
    </div>
  );
};
