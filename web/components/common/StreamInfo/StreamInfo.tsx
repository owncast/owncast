import { useRecoilValue } from 'recoil';
import cn from 'classnames';
import { EyeFilled } from '@ant-design/icons';
import { useEffect } from 'react';
import { ClientConfig } from '../../../interfaces/client-config.model';
import {
  clientConfigStateAtom,
  isOnlineSelector,
  serverStatusState,
} from '../../stores/ClientConfigStore';
import { ServerLogo } from '../../ui';
import SocialLinks from '../../ui/SocialLinks/SocialLinks';
import s from './StreamInfo.module.scss';
import { ServerStatus } from '../../../interfaces/server-status.model';

interface Props {
  isMobile: boolean;
}
export default function StreamInfo({ isMobile }: Props) {
  const { socialHandles, name, title, tags, summary } =
    useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { viewerCount } = useRecoilValue<ServerStatus>(serverStatusState);
  const online = useRecoilValue<boolean>(isOnlineSelector);

  useEffect(() => {
    console.log({ online });
  }, [online]);

  return isMobile ? (
    <div className={cn(s.root, s.mobile)}>
      <div className={s.mobileInfo}>
        <ServerLogo src="/logo" />
        <div className={s.title}>{name}</div>
      </div>
      <div className={s.mobileStatus}>
        <div className={s.viewerCount}>
          {online && (
            <>
              <span>{viewerCount}</span>
              <EyeFilled />
            </>
          )}
        </div>
        <div className={s.liveStatus}>
          {online && <div className={s.liveCircle} />}
          <span>{online ? 'LIVE' : 'OFFLINE'}</span>
        </div>
      </div>
    </div>
  ) : (
    <div className={s.root}>
      <div className={s.logoTitleSection}>
        <ServerLogo src="/logo" />
        <div className={s.titleSection}>
          <div className={s.title}>{name}</div>
          <div className={s.subtitle}>{title || summary}</div>
          <div className={s.tagList}>
            {tags.length > 0 && tags.map(tag => <span key={tag}>#{tag}&nbsp;</span>)}
          </div>
          <SocialLinks links={socialHandles} />
        </div>
      </div>
    </div>
  );
}
