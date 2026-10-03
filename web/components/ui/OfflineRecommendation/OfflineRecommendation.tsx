import type { FC } from 'react';
import { ArrowRightOutlined } from '@ant-design/icons';
import { useTranslation } from 'next-export-i18n';
import type { FederatedServerResponse } from '../../../hooks/useFederatedServers';
import { Localization } from '../../../types/localization';
import { VideoPoster } from '../../video/VideoPoster/VideoPoster';
import { StreamStatusBadge } from '../StreamStatusBadge/StreamStatusBadge';
import { Translation } from '../Translation/Translation';
import styles from './OfflineRecommendation.module.scss';

export type OfflineRecommendationProps = {
  localStreamName: string;
  server: FederatedServerResponse;
};

const getServerHost = (serverUrl: string) => {
  try {
    return new URL(serverUrl).hostname;
  } catch {
    return serverUrl;
  }
};

export const OfflineRecommendation: FC<OfflineRecommendationProps> = ({
  localStreamName,
  server,
}) => {
  const { t } = useTranslation();
  const serverHost = getServerHost(server.iri);
  const serverName = server.displayName || server.name || serverHost;
  const description = server.streamDescription || server.summary;
  const thumbnail = server.thumbnailUrl;

  return (
    <aside className={styles.panel}>
      <p className={styles.introduction}>
        <Translation
          translationKey={Localization.Frontend.OfflineRecommendation.introduction}
          vars={{ name: localStreamName }}
          defaultText="While this stream is offline, visit this recommended stream."
        />
      </p>
      <a
        className={styles.recommendation}
        href={server.iri}
        target="_blank"
        rel="noopener noreferrer"
        aria-label={t(Localization.Frontend.OfflineRecommendation.linkAria, {
          name: serverName,
          host: serverHost,
        })}
      >
        <div className={styles.image}>
          {thumbnail && (
            <VideoPoster
              online
              initialSrc={thumbnail}
              src={thumbnail}
              refreshInterval={thumbnail.startsWith('data:') ? 0 : 15_000}
              objectFit="cover"
              width="100%"
              height="100%"
              referrerPolicy="no-referrer"
              onError={event => {
                event.currentTarget.hidden = true;
              }}
            />
          )}
        </div>
        <div className={styles.details}>
          <div className={styles.identity}>
            {server.logoUrl && (
              <img
                className={styles.logo}
                src={server.logoUrl}
                alt=""
                referrerPolicy="no-referrer"
                onError={event => {
                  event.currentTarget.hidden = true;
                }}
              />
            )}
            <strong className={styles.name}>{serverName}</strong>
          </div>
          <StreamStatusBadge isOnline className={styles.liveLabel} />
          {description && <span className={styles.description}>{description}</span>}
          <span className={styles.streamTitle}>{server.streamTitle || serverHost}</span>
          <span className={styles.watchAction}>
            <Translation
              translationKey={Localization.Frontend.OfflineRecommendation.watchLive}
              defaultText="Watch live"
            />
            <ArrowRightOutlined aria-hidden="true" />
          </span>
        </div>
      </a>
    </aside>
  );
};
