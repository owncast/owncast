import { Modal } from 'antd';
import { FC } from 'react';
import { ErrorBoundary } from 'react-error-boundary';

import styles from './ChatModal.module.scss';

import { ComponentError } from '../../ui/ComponentError/ComponentError';
import { ChatContainer } from '../../chat/ChatContainer/ChatContainer';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { CurrentUser } from '../../../interfaces/current-user';
import { UserDropdown } from '../../common/UserDropdown/UserDropdown';

export type ChatModalProps = {
  messages: ChatMessage[];
  currentUser: CurrentUser;
  handleClose: () => void;
};

export const ChatModal: FC<ChatModalProps> = ({ messages, currentUser, handleClose }) => {
  if (!currentUser) {
    return null;
  }
  const { id, displayName, isModerator } = currentUser;

  const modalBodyStyle = {
    padding: '0px',
    height: '55vh',
  };

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="ChatModal"
          message={error.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <Modal
        open
        centered
        maskClosable={false}
        footer={null}
        title={<UserDropdown id="chat-modal-user-menu" showToggleChatOption={false} />}
        maskStyle={{
          zIndex: 700,
        }}
        className={styles.root}
        bodyStyle={modalBodyStyle}
        wrapClassName={styles.modalWrapper}
        onCancel={handleClose}
      >
        <ChatContainer
          messages={messages}
          usernameToHighlight={displayName}
          chatUserId={id}
          isModerator={isModerator}
          chatAvailable
          focusInput={false}
        />
      </Modal>
    </ErrorBoundary>
  );
};
