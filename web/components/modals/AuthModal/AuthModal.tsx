import { Modal, Tabs } from 'antd';
import { useAtomValue } from 'jotai';
import { FC } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { IndieAuthModal } from '../IndieAuthModal/IndieAuthModal';
import { FediAuthModal } from '../FediAuthModal/FediAuthModal';
import { getPendingFediverseAuth } from '../../../utils/fediverseAuthSession';

import styles from './AuthModal.module.scss';
import {
  currentUserAtom,
  chatAuthenticatedAtom,
  accessTokenAtom,
  clientConfigStateAtom,
} from '../../stores/ClientConfigStore';
import { ComponentError } from '../../ui/ComponentError/ComponentError';

export type AuthModalProps = {
  open: boolean;
  handleClose: () => void;
  forceTabs?: boolean;
};

export const AuthModal: FC<AuthModalProps> = ({ open, handleClose, forceTabs }) => {
  const authenticated = useAtomValue(chatAuthenticatedAtom);
  const accessToken = useAtomValue(accessTokenAtom);
  const currentUser = useAtomValue(currentUserAtom);
  const clientConfig = useAtomValue(clientConfigStateAtom);

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

  // If a fediverse verification is still in progress (restored after a reload),
  // open straight to that tab so the recovered code-entry step is visible. Only
  // when federation is enabled, since the FediAuth tab is unreachable otherwise.
  const defaultActiveKey = fediverseEnabled && getPendingFediverseAuth() ? '2' : '1';

  return (
    <Modal
      title="Authenticate"
      open={open}
      onCancel={handleClose}
      maskClosable={false}
      zIndex={999}
      footer={null}
      centered
    >
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
            defaultActiveKey={defaultActiveKey}
            items={items}
            type="card"
            size="small"
            renderTabBar={fediverseEnabled || forceTabs ? undefined : () => null}
          />
        </div>
      </ErrorBoundary>
    </Modal>
  );
};
