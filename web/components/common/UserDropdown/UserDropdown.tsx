import { Menu, Dropdown, Button } from 'antd';
import classnames from 'classnames';

import { useRecoilState, useRecoilValue } from 'recoil';
import { FC, useState } from 'react';
import { useHotkeys } from 'react-hotkeys-hook';
import dynamic from 'next/dynamic';
import { ErrorBoundary } from 'react-error-boundary';
import {
  ChatState,
  chatStateAtom,
  currentUserAtom,
  appStateAtom,
} from '../../stores/ClientConfigStore';
import styles from './UserDropdown.module.scss';
import { AppStateOptions } from '../../stores/application-state';
import { ComponentError } from '../../ui/ComponentError/ComponentError';

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

const ShrinkOutlined = dynamic(() => import('@ant-design/icons/ShrinkOutlined'), {
  ssr: false,
});

const ExpandAltOutlined = dynamic(() => import('@ant-design/icons/ExpandAltOutlined'), {
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
  id: string;
  username?: string;
  hideTitleOnMobile?: boolean;
  showToggleChatOption?: boolean;
};

export const UserDropdown: FC<UserDropdownProps> = ({
  id,
  username: defaultUsername = undefined,
  hideTitleOnMobile = false,
  showToggleChatOption: showHideChatOption = true,
}) => {
  const [showNameChangeModal, setShowNameChangeModal] = useState<boolean>(false);
  const [showAuthModal, setShowAuthModal] = useState<boolean>(false);
  const [chatState, setChatState] = useRecoilState(chatStateAtom);
  const [popupWindow, setPopupWindow] = useState<Window>(null);
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);

  const toggleChatVisibility = () => {
    // If we don't support the hide chat option then don't do anything.
    if (!showHideChatOption) {
      return;
    }

    setChatState(chatState === ChatState.VISIBLE ? ChatState.HIDDEN : ChatState.VISIBLE);
  };

  const handleChangeName = () => {
    setShowNameChangeModal(true);
  };

  const closeChangeNameModal = () => {
    setShowNameChangeModal(false);
  };

  const closeChatPopup = () => {
    if (popupWindow) {
      popupWindow.close();
    }
    setPopupWindow(null);
    setChatState(ChatState.VISIBLE);
  };

  const openChatPopup = () => {
    // close popup (if any) to prevent multiple popup windows.
    closeChatPopup();
    const w = window.open('/embed/chat/readwrite', '_blank', 'popup');
    w.addEventListener('beforeunload', closeChatPopup);
    setPopupWindow(w);
    setChatState(ChatState.POPPED_OUT);
  };

  const canShowHideChat =
    showHideChatOption &&
    appState.chatAvailable &&
    (chatState === ChatState.HIDDEN || chatState === ChatState.VISIBLE);
  const canShowChatPopup =
    showHideChatOption &&
    appState.chatAvailable &&
    (chatState === ChatState.HIDDEN ||
      chatState === ChatState.VISIBLE ||
      chatState === ChatState.POPPED_OUT);

  // Register keyboard shortcut for the space bar to toggle playback
  useHotkeys(
    'c',
    toggleChatVisibility,
    {
      enableOnContentEditable: false,
    },
    [chatState === ChatState.VISIBLE],
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
      {canShowHideChat && (
        <Menu.Item
          key="3"
          icon={<MessageOutlined />}
          onClick={() => toggleChatVisibility()}
          aria-expanded={chatState === ChatState.VISIBLE}
          className={styles.chatToggle}
        >
          {chatState === ChatState.VISIBLE ? 'Hide Chat' : 'Show Chat'}
        </Menu.Item>
      )}
      {canShowChatPopup &&
        (popupWindow ? (
          <Menu.Item key="4" icon={<ShrinkOutlined />} onClick={closeChatPopup}>
            Put chat back
          </Menu.Item>
        ) : (
          <Menu.Item key="4" icon={<ExpandAltOutlined />} onClick={openChatPopup}>
            Pop out chat
          </Menu.Item>
        ))}
    </Menu>
  );

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="UserDropdown"
          message={error.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <div className={styles.root}>
        <Dropdown overlay={menu} trigger={['click']}>
          <Button id={id} type="primary" icon={<UserOutlined className={styles.userIcon} />}>
            <span
              className={classnames([
                styles.username,
                hideTitleOnMobile && styles.hideTitleOnMobile,
              ])}
            >
              {username}
            </span>
            <CaretDownOutlined />
          </Button>
        </Dropdown>
        <Modal
          title="Change Chat Display Name"
          open={showNameChangeModal}
          handleCancel={closeChangeNameModal}
        >
          <NameChangeModal closeModal={closeChangeNameModal} />
        </Modal>
        <Modal
          title="Authenticate"
          open={showAuthModal}
          handleCancel={() => setShowAuthModal(false)}
        >
          <AuthModal />
        </Modal>
      </div>
    </ErrorBoundary>
  );
};
