import { Typography } from 'antd';
import { FC, useEffect, useState } from 'react';
import { formatDistanceToNow } from 'date-fns';
import dynamic from 'next/dynamic';
import classNames from 'classnames';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../../types/localization';
import type { ScheduledEventStatus } from '../../../interfaces/server-status.model';
import { PUBLIC_SCHEDULE_ICS } from '../../../utils/apis';
import styles from './EventCountdown.module.scss';

const ClockCircleOutlined = dynamic(() => import('@ant-design/icons/ClockCircleOutlined'), {
  ssr: false,
});

const getRemainingSeconds = (startTime: string, now: number) =>
  Math.max(0, Math.ceil((new Date(startTime).getTime() - now) / 1000));

export const EventCountdown: FC<{
  event: ScheduledEventStatus;
  lastLive?: Date;
  className?: string;
}> = ({ event, lastLive, className }) => {
  const { t } = useTranslation();
  const [remainingSeconds, setRemainingSeconds] = useState<number | null>(null);

  useEffect(() => {
    const update = () => setRemainingSeconds(getRemainingSeconds(event.startTime, Date.now()));
    update();
    const interval = window.setInterval(update, 1000);
    return () => window.clearInterval(interval);
  }, [event.startTime]);
  const secondsRemaining = remainingSeconds ?? 0;
  const days = Math.floor(secondsRemaining / 86400);
  const hours = Math.floor((secondsRemaining % 86400) / 3600);
  const minutes = Math.floor((secondsRemaining % 3600) / 60);
  const seconds = secondsRemaining % 60;
  const translate = (key: string, fallback: string, count: number) => {
    const translated = t(key, { count: `${count}` });
    return translated === key ? fallback : translated;
  };
  const time = [
    translate(Localization.Frontend.Schedule.countdownDays, `${days} days`, days),
    translate(Localization.Frontend.Schedule.countdownHours, `${hours} hours`, hours),
    translate(Localization.Frontend.Schedule.countdownMinutes, `${minutes} minutes`, minutes),
    translate(Localization.Frontend.Schedule.countdownSeconds, `${seconds} seconds`, seconds),
  ].join(', ');

  const countdownText =
    remainingSeconds === null
      ? null
      : remainingSeconds === 0
        ? t(Localization.Frontend.Schedule.countdownLiveNow)
        : t(Localization.Frontend.Schedule.countdownLiveIn, { time });

  return (
    <div className={classNames(styles.container, className)} aria-live="polite">
      <div className={styles.content}>
        <Typography.Title level={2} className={styles.title}>
          {event.name}
        </Typography.Title>
        {countdownText && (
          <Typography.Text className={styles.countdown}>{countdownText}</Typography.Text>
        )}
        {event.description && <Typography.Paragraph>{event.description}</Typography.Paragraph>}
        <a className={styles.calendarLink} href={PUBLIC_SCHEDULE_ICS}>
          {t(Localization.Frontend.Schedule.addToCalendar)}
        </a>
        {lastLive && (
          <div className={styles.lastLiveDate}>
            <ClockCircleOutlined className={styles.clockIcon} />
            <span id="owncast-offline-last-live-text">
              {t(Localization.Frontend.lastLiveAgo, {
                timeAgo: formatDistanceToNow(new Date(lastLive)),
              })}
            </span>
          </div>
        )}
      </div>
    </div>
  );
};
