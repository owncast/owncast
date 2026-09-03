import Head from 'next/head';
import Link from 'next/link';
import { useTranslation } from 'next-export-i18n';
import { useEffect, useMemo, useState } from 'react';
import type { CSSProperties, FC } from 'react';
import { EventCountdown } from '../EventCountdown/EventCountdown';
import type { ScheduledEvent } from '../../../interfaces/scheduled-event.model';
import type { ClientConfig } from '../../../interfaces/client-config.model';
import { Localization } from '../../../types/localization';
import { fetchData, PUBLIC_SCHEDULE } from '../../../utils/apis';
import styles from './ScheduleEventPage.module.scss';

const getText = (
  translate: (key: string, options?: Record<string, string>) => string,
  key: string,
  fallback: string,
) => {
  const translated = translate(key);
  return translated === key ? fallback : translated;
};

const getEventID = () => {
  if (typeof window === 'undefined') {
    return '';
  }

  const parts = window.location.pathname.split('/').filter(Boolean);
  return parts[0] === 'schedule' && parts[1] ? decodeURIComponent(parts[1]) : '';
};

const formatDateTime = (event: ScheduledEvent) =>
  new Intl.DateTimeFormat(undefined, {
    dateStyle: 'full',
    timeStyle: 'short',
    timeZone: event.timezone,
  }).format(new Date(event.startTime));

export interface ScheduleEventPageProps {
  eventID?: string;
  logoSrc?: string;
}

export const ScheduleEventPage: FC<ScheduleEventPageProps> = ({
  eventID: providedEventID,
  logoSrc = '/logo',
}) => {
  const { t } = useTranslation();
  const [event, setEvent] = useState<ScheduledEvent | null>(null);
  const [serverName, setServerName] = useState('Owncast');
  const [config, setConfig] = useState<ClientConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    const eventID = providedEventID || getEventID();
    if (!eventID) {
      setLoading(false);
      return;
    }

    const from = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString();
    const to = new Date(Date.now() + 90 * 24 * 60 * 60 * 1000).toISOString();
    const url = new URL(PUBLIC_SCHEDULE, window.location.origin);
    url.searchParams.set('from', from);
    url.searchParams.set('to', to);

    Promise.all([
      fetchData<ScheduledEvent[]>(url.toString(), { auth: false }),
      fetchData<ClientConfig>('/api/config', { auth: false }).catch(() => null),
    ])
      .then(([events, loadedConfig]) => {
        setEvent(events.find(candidate => candidate.id === eventID) || null);
        if (loadedConfig) {
          setConfig(loadedConfig);
          setServerName(loadedConfig.name || 'Owncast');
        }
      })
      .catch(() => setError(true))
      .finally(() => setLoading(false));
  }, [providedEventID]);

  useEffect(() => {
    if (!event) {
      return undefined;
    }

    const timer = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(timer);
  }, [event]);

  const themeStyle = useMemo(() => {
    const variables = config?.appearanceVariables || {};
    return Object.fromEntries(
      Object.entries(variables).map(([name, value]) => [`--${name}`, value]),
    ) as CSSProperties;
  }, [config]);

  const loadingText = getText(t, Localization.Frontend.Schedule.loadingEvent, 'Loading event...');
  const notFoundText = getText(
    t,
    Localization.Frontend.Schedule.eventNotFound,
    'This scheduled event could not be found',
  );
  const errorText = getText(t, Localization.Frontend.Schedule.error, 'Unable to load the schedule');
  const backText = getText(t, Localization.Frontend.Schedule.backToStream, 'Back to stream');
  const cancelledText = getText(t, Localization.Frontend.Schedule.cancelled, 'Cancelled');
  const endedText = getText(t, Localization.Frontend.Schedule.eventEnded, 'This event has ended');
  const pageEyebrow = getText(t, Localization.Frontend.Schedule.pageEyebrow, 'Owncast schedule');
  const timezoneText = getText(t, Localization.Frontend.Schedule.timezone, 'Schedule timezone');
  const durationText = getText(t, Localization.Frontend.Schedule.duration, '{{minutes}} minutes');
  const pageTitle = event ? `${event.name} | ${serverName}` : serverName;
  const pageURL =
    typeof window === 'undefined' ? '' : `${window.location.origin}${window.location.pathname}`;
  const imageURL =
    typeof window === 'undefined' ? logoSrc : new URL(logoSrc, window.location.origin).toString();
  const description = event?.description || pageTitle;
  const start = event ? new Date(event.startTime).getTime() : 0;
  const end = event ? start + event.durationMinutes * 60 * 1000 : 0;
  const eventHasEnded = Boolean(event && now >= end);

  return (
    <div className={styles.page} style={themeStyle}>
      <Head>
        <title>{pageTitle}</title>
        {event && (
          <>
            <meta property="og:title" content={pageTitle} />
            <meta property="og:description" content={description} />
            <meta property="og:type" content="event" />
            <meta property="og:url" content={pageURL} />
            <meta property="og:site_name" content={serverName} />
            <meta property="og:image" content={imageURL} />
            <meta property="og:image:alt" content={pageTitle} />
            <meta property="event:start_time" content={new Date(start).toISOString()} />
            <meta property="event:end_time" content={new Date(end).toISOString()} />
            <meta name="twitter:card" content="summary" />
            <meta name="twitter:title" content={pageTitle} />
            <meta name="twitter:description" content={description} />
            <meta name="twitter:image" content={imageURL} />
          </>
        )}
        {config?.pluginStyles && <style>{config.pluginStyles}</style>}
        {config?.customStyles && <style>{config.customStyles}</style>}
      </Head>
      <header className={styles.header}>
        <Link className={styles.brand} href="/" aria-label={backText}>
          <img src={logoSrc} alt="" />
          <span>{serverName}</span>
        </Link>
        <Link className={styles.backLink} href="/">
          {backText}
        </Link>
      </header>
      <main className={styles.main}>
        {loading && (
          <div className={styles.message} role="status" aria-live="polite">
            {loadingText}
          </div>
        )}
        {!loading && error && <div className={styles.message}>{errorText}</div>}
        {!loading && !error && !event && <div className={styles.message}>{notFoundText}</div>}
        {event && (
          <>
            <div className={styles.eyebrow}>{pageEyebrow}</div>
            {event.status === 'cancelled' ? (
              <section className={styles.cancelledHero} aria-labelledby="event-title">
                <span className={styles.cancelledLabel}>{cancelledText}</span>
                <h1 id="event-title">{event.name}</h1>
                {event.description && <p>{event.description}</p>}
              </section>
            ) : eventHasEnded ? (
              <section className={styles.endedHero} aria-labelledby="event-title">
                <h1 id="event-title">{event.name}</h1>
                <p>{endedText}</p>
              </section>
            ) : (
              <EventCountdown event={{ ...event, chatOpen: false }} />
            )}
            <section className={styles.details} aria-labelledby="details-title">
              <h2 id="details-title">{formatDateTime(event)}</h2>
              <p className={styles.timezone}>
                {timezoneText}: {event.timezone}
              </p>
              <p>{durationText.replace('{{minutes}}', `${event.durationMinutes}`)}</p>
            </section>
          </>
        )}
      </main>
    </div>
  );
};
