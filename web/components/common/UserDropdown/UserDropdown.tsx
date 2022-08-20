import { Menu, Dropdown, Button, Space } from 'antd';
import {
  CaretDownOutlined,
  EditOutlined,
  LockOutlined,
  MessageOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useRecoilState, useRecoilValue } from 'recoil';
import { useState } from 'react';
import { useHotkeys } from 'react-hotkeys-hook';
import Modal from '../../ui/Modal/Modal';
import {
  chatVisibleToggleAtom,
  chatDisplayNameAtom,
  appStateAtom,
} from '../../stores/ClientConfigStore';
import s from './UserDropdown.module.scss';
import NameChangeModal from '../../modals/NameChangeModal';
import { AppStateOptions } from '../../stores/application-state';
import AuthModal from '../../modals/AuthModal/AuthModal';

interface Props {
  username?: string;
}

export default function UserDropdown({ username: defaultUsername }: Props) {
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
    <div className={`${s.root}`}>
      <Dropdown overlay={menu} trigger={['click']}>
        <Button icon={<UserOutlined style={{ marginRight: '.5rem' }} />}>
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
}

UserDropdown.defaultProps = {
  username: undefined,
};
