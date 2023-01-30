import React from 'react';
import { useRecoilValue } from 'recoil';
import { useRouter } from 'next/router';
import {
  clientConfigStateAtom,
  ClientConfigStore,
  isOnlineSelector,
  serverStatusState,
} from '../../../components/stores/ClientConfigStore';
import { OfflineBanner } from '../../../components/ui/OfflineBanner/OfflineBanner';
import { Statusbar } from '../../../components/ui/Statusbar/Statusbar';
import { OwncastPlayer } from '../../../components/video/OwncastPlayer/OwncastPlayer';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { ServerStatus } from '../../../interfaces/server-status.model';

export default function VideoEmbed() {
  const status = useRecoilValue<ServerStatus>(serverStatusState);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);

  const { name } = clientConfig;

  const { offlineMessage } = clientConfig;
  const { viewerCount, lastConnectTime, lastDisconnectTime, streamTitle } = status;
  const online = useRecoilValue<boolean>(isOnlineSelector);

  const router = useRouter();

  /**
   * router.query isn't initialized until hydration
   * (see https://github.com/vercel/next.js/discussions/11484)
   * but router.asPath is initialized earlier, so we parse the
   * query parameters ourselves
   */
  const path = router.asPath.split('?')[1] ?? '';
  const query = path.split('&').reduce((currQuery, part) => {
    const [key, value] = part.split('=');
    return { ...currQuery, [key]: value };
  }, {} as Record<string, string>);

  const initiallyMuted = query.initiallyMuted === 'true';

  return (
    <>
      <ClientConfigStore />
      <div className="video-embed">
        {online && (
          <OwncastPlayer
            source="/hls/stream.m3u8"
            online={online}
            initiallyMuted={initiallyMuted}
            title={streamTitle || name}
          />
        )}
        {!online && (
          <OfflineBanner
            streamName={name}
            customText={offlineMessage}
            notificationsEnabled={false}
          />
        )}
        {online && (
          <Statusbar
            online={online}
            lastConnectTime={lastConnectTime}
            lastDisconnectTime={lastDisconnectTime}
            viewerCount={viewerCount}
          />
        )}
      </div>
    </>
  );
}
