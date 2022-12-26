import { Tabs } from 'antd';
import { useRecoilValue } from 'recoil';
import { FC } from 'react';
import { IndieAuthModal } from '../IndieAuthModal/IndieAuthModal';
import { FediAuthModal } from '../FediAuthModal/FediAuthModal';

import FediverseIcon from '../../../assets/images/fediverse-black.png';
import IndieAuthIcon from '../../../assets/images/indieauth.png';

import styles from './AuthModal.module.scss';
import {
  currentUserAtom,
  chatAuthenticatedAtom,
  accessTokenAtom,
  clientConfigStateAtom,
} from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';

export const AuthModal: FC = () => {
  const authenticated = useRecoilValue<boolean>(chatAuthenticatedAtom);
  const accessToken = useRecoilValue<string>(accessTokenAtom);
  const currentUser = useRecoilValue(currentUserAtom);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);

  if (!currentUser) {
    return null;
  }
  const { displayName } = currentUser;
  const { federation } = clientConfig;
  const { enabled: fediverseEnabled } = federation;

  const indieAuthTabTitle = (
    <span className={styles.tabContent}>
      <img className={styles.icon} src={IndieAuthIcon.src} alt="IndieAuth" />
      IndieAuth
    </span>
  );

  const indieAuthTab = (
    <IndieAuthModal
      authenticated={authenticated}
      displayName={displayName}
      accessToken={accessToken}
    />
  );

  const fediAuthTabTitle = (
    <span className={styles.tabContent}>
      <img className={styles.icon} src={FediverseIcon.src} alt="Fediverse auth" />
      FediAuth
    </span>
  );

  const fediAuthTab = (
    <FediAuthModal
      authenticated={authenticated}
      displayName={displayName}
      accessToken={accessToken}
    />
  );

  const items = [
    { label: indieAuthTabTitle, key: '1', children: indieAuthTab },
    { label: fediAuthTabTitle, key: '2', children: fediAuthTab },
  ];

  return (
    <div>
      <Tabs
        defaultActiveKey="1"
        items={items}
        type="card"
        size="small"
        renderTabBar={fediverseEnabled ? null : () => null}
      />
    </div>
  );
};
