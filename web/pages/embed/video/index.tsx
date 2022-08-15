import React from 'react';
import { useRecoilValue } from 'recoil';
import {
  clientConfigStateAtom,
  ClientConfigStore,
  isOnlineSelector,
  serverStatusState,
} from '../../../components/stores/ClientConfigStore';
import OfflineBanner from '../../../components/ui/OfflineBanner/OfflineBanner';
import Statusbar from '../../../components/ui/Statusbar/Statusbar';
import OwncastPlayer from '../../../components/video/OwncastPlayer';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { ServerStatus } from '../../../interfaces/server-status.model';

export default function VideoEmbed() {
  const status = useRecoilValue<ServerStatus>(serverStatusState);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);

  const { name } = clientConfig;

  // const { extraPageContent, version, socialHandles, name, title, tags } = clientConfig;
  const { viewerCount, lastConnectTime, lastDisconnectTime } = status;
  const online = useRecoilValue<boolean>(isOnlineSelector);
  return (
    <>
      <ClientConfigStore />
      <div className="video-embed">
        {online && <OwncastPlayer source="/hls/stream.m3u8" online={online} />}
        {!online && <OfflineBanner name={name} text="Stream is offline text goes here." />}{' '}
        <Statusbar
          online={online}
          lastConnectTime={lastConnectTime}
          lastDisconnectTime={lastDisconnectTime}
          viewerCount={viewerCount}
        />
      </div>
    </>
  );
}
