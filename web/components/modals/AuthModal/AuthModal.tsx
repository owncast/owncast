import { Tabs } from 'antd';
import { useRecoilValue } from 'recoil';
import { FC } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { IndieAuthModal } from '../IndieAuthModal/IndieAuthModal';
import { FediAuthModal } from '../FediAuthModal/FediAuthModal';

import styles from './AuthModal.module.scss';
import {
  currentUserAtom,
  chatAuthenticatedAtom,
  accessTokenAtom,
  clientConfigStateAtom,
} from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { ComponentError } from '../../ui/ComponentError/ComponentError';

export type AuthModalProps = {
  forceTabs?: boolean;
};

export const AuthModal: FC<AuthModalProps> = ({ forceTabs }) => {
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
      <img className={styles.icon} src="/img/indieauth.png" alt="IndieAuth" />
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
      <img className={styles.icon} src="/img/fediverse-black.png" alt="Fediverse auth" />
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
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="AuthModal"
          message={error.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <div>
        <Tabs
          defaultActiveKey="1"
          items={items}
          type="card"
          size="small"
          renderTabBar={fediverseEnabled || forceTabs ? null : () => null}
        />
      </div>
    </ErrorBoundary>
  );
};
