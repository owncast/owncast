import { FC, useState } from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import dynamic from 'next/dynamic';
import { useHotkeys } from 'react-hotkeys-hook';
import {
  chatVisibleToggleAtom,
  chatDisplayNameAtom,
  appStateAtom,
} from '~/components//stores/ClientConfigStore';
import { AppStateOptions } from '~/components/stores/application-state';
import { Menu, Dropdown, Button, Space } from 'antd';
import {
  CaretDownOutlined,
  EditOutlined,
  LockOutlined,
  MessageOutlined,
  UserOutlined,
} from '@ant-design/icons';
import styles from './UserDropdown.module.scss';

// Lazy loaded components
const Modal = dynamic(() => import('../../ui/Modal/Modal').then(mod => mod.Modal));

const NameChangeModal = dynamic(() =>
  import('../../modals/NameChangeModal/NameChangeModal').then(mod => mod.NameChangeModal),
);

const AuthModal = dynamic(() =>
  import('../../modals/AuthModal/AuthModal').then(mod => mod.AuthModal),
);

export type UserDropdownProps = {
  username?: string;
};

export const UserDropdown: FC<UserDropdownProps> = ({ username: defaultUsername = undefined }) => {
  const username = defaultUsername || useRecoilValue(chatDisplayNameAtom);
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

  const menu = (
    <Menu>
      <Menu.Item key="0" icon={<EditOutlined />} onClick={() => handleChangeName()}>
        Change name
      </Menu.Item>
      <Menu.Item key="1" icon={<LockOutlined />} onClick={() => setShowAuthModal(true)}>
        Authenticate
      </Menu.Item>
      {appState.chatAvailable && (
        <Menu.Item key="3" icon={<MessageOutlined />} onClick={() => toggleChatVisibility()}>
          Toggle chat
        </Menu.Item>
      )}
    </Menu>
  );

  return (
    <div className={`${styles.root}`}>
      <Dropdown overlay={menu} trigger={['click']}>
        <Button type="primary" icon={<UserOutlined style={{ marginRight: '.5rem' }} />}>
          <Space>
            {username}
            <CaretDownOutlined />
          </Space>
        </Button>
      </Dropdown>
      <Modal
        title="Change Chat Display Name"
        visible={showNameChangeModal}
        handleCancel={() => setShowNameChangeModal(false)}
      >
        <NameChangeModal />
      </Modal>
      <Modal
        title="Authenticate"
        visible={showAuthModal}
        handleCancel={() => setShowAuthModal(false)}
      >
        <AuthModal />
      </Modal>
    </div>
  );
};
