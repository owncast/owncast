import { Menu, Dropdown, Button } from 'antd';

import { useRecoilState, useRecoilValue } from 'recoil';
import { FC, useState } from 'react';
import { useHotkeys } from 'react-hotkeys-hook';
import dynamic from 'next/dynamic';
import {
  chatVisibleToggleAtom,
  currentUserAtom,
  appStateAtom,
} from '../../stores/ClientConfigStore';
import styles from './UserDropdown.module.scss';
import { AppStateOptions } from '../../stores/application-state';

// Lazy loaded components

const CaretDownOutlined = dynamic(() => import('@ant-design/icons/CaretDownOutlined'), {
  ssr: false,
});

const EditOutlined = dynamic(() => import('@ant-design/icons/EditOutlined'), {
  ssr: false,
});

const LockOutlined = dynamic(() => import('@ant-design/icons/LockOutlined'), {
  ssr: false,
});

const MessageOutlined = dynamic(() => import('@ant-design/icons/MessageOutlined'), {
  ssr: false,
});

const UserOutlined = dynamic(() => import('@ant-design/icons/UserOutlined'), {
  ssr: false,
});

const Modal = dynamic(() => import('../../ui/Modal/Modal').then(mod => mod.Modal), {
  ssr: false,
});

const NameChangeModal = dynamic(
  () => import('../../modals/NameChangeModal/NameChangeModal').then(mod => mod.NameChangeModal),
  {
    ssr: false,
  },
);

const AuthModal = dynamic(
  () => import('../../modals/AuthModal/AuthModal').then(mod => mod.AuthModal),
  {
    ssr: false,
  },
);

export type UserDropdownProps = {
  username?: string;
};

export const UserDropdown: FC<UserDropdownProps> = ({ username: defaultUsername = undefined }) => {
  const [showNameChangeModal, setShowNameChangeModal] = useState<boolean>(false);
  const [showAuthModal, setShowAuthModal] = useState<boolean>(false);
  const [chatToggleVisible, setChatToggleVisible] = useRecoilState(chatVisibleToggleAtom);
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);

  const toggleChatVisibility = () => {
    setChatToggleVisible(!chatToggleVisible);
  };

  const handleChangeName = () => {
    setShowNameChangeModal(true);
  };

  // Register keyboard shortcut for the space bar to toggle playback
  useHotkeys(
    'c',
    toggleChatVisibility,
    {
      enableOnContentEditable: false,
    },
    [chatToggleVisible],
  );

  const currentUser = useRecoilValue(currentUserAtom);
  if (!currentUser) {
    return null;
  }

  const { displayName } = currentUser;
  const username = defaultUsername || displayName;
  const menu = (
    <Menu>
      <Menu.Item key="0" icon={<EditOutlined />} onClick={() => handleChangeName()}>
        Change name
      </Menu.Item>
      <Menu.Item key="1" icon={<LockOutlined />} onClick={() => setShowAuthModal(true)}>
        Authenticate
      </Menu.Item>
      {appState.chatAvailable && (
        <Menu.Item
          key="3"
          icon={<MessageOutlined />}
          onClick={() => toggleChatVisibility()}
          aria-expanded={chatToggleVisible}
        >
          {chatToggleVisible ? 'Hide Chat' : 'Show Chat'}
        </Menu.Item>
      )}
    </Menu>
  );

  return (
    <div id="user-menu" className={`${styles.root}`}>
      <Dropdown overlay={menu} trigger={['click']}>
        <Button type="primary" icon={<UserOutlined className={styles.userIcon} />}>
          <span className={styles.username}>{username}</span>
          <CaretDownOutlined />
        </Button>
      </Dropdown>
      <Modal
        title="Change Chat Display Name"
        open={showNameChangeModal}
        handleCancel={() => setShowNameChangeModal(false)}
      >
        <NameChangeModal />
      </Modal>
      <Modal title="Authenticate" open={showAuthModal} handleCancel={() => setShowAuthModal(false)}>
        <AuthModal />
      </Modal>
    </div>
  );
};
