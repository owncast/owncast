import { MenuProps, Dropdown, Button } from 'antd';
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

  const items: MenuProps['items'] = [
    {
      key: 0,
      icon: <EditOutlined />,
      label: 'Change name',
      onClick: handleChangeName,
    },
    {
      key: 1,
      icon: <LockOutlined />,
      label: 'Authenticate',
      onClick: () => setShowAuthModal(true),
    },
  ];
  if (canShowHideChat)
    items.push({
      key: 3,
      'aria-expanded': chatState === ChatState.VISIBLE,
      className: styles.chatToggle, // TODO why do we hide this button on tablets?
      icon: <MessageOutlined />,
      label: chatState === ChatState.VISIBLE ? 'Hide Chat' : 'Show Chat',
      onClick: toggleChatVisibility,
    } as MenuProps['items'][0]);
  if (canShowChatPopup)
    items.push({
      key: 4,
      icon: popupWindow ? <ShrinkOutlined /> : <ExpandAltOutlined />,
      label: popupWindow ? 'Put chat back' : 'Pop out chat',
      onClick: popupWindow ? closeChatPopup : openChatPopup,
    });

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
        <Dropdown menu={{ items }} trigger={['click']}>
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
