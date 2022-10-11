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
} from '../../stores/ClientConfigStore';

export const AuthModal: FC = () => {
  const currentUser = useRecoilValue(currentUserAtom);
  if (!currentUser) {
    return null;
  }

  const authenticated = useRecoilValue<boolean>(chatAuthenticatedAtom);
  const accessToken = useRecoilValue<string>(accessTokenAtom);
  const federationEnabled = true;
  const { displayName } = currentUser;

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
        renderTabBar={federationEnabled ? null : () => null}
      />
    </div>
  );
};
