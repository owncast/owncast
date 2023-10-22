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

  const modalWrapperStyle = {
    zIndex: 799,
    top: 'unset',
  };

  const modalContentStyle = {
    padding: 0,
  };

  const modalHeaderStyle = {
    padding: '16px 24px',
    background: 'var(--theme-color-components-modal-header-background)',
    margin: 0,
    borderBottom: '1px solid #f0f0f0',
  };

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
        className={styles.root}
        styles={{
          header: modalHeaderStyle,
          body: modalBodyStyle,
          mask: { zIndex: 700 },
          content: modalContentStyle,
          // styles.wrapper was added recently (5.10.0) and wrapClassName isn't working
          // So, using ts-ignore for now until the documentation/implementation mismatch is resolved
          // Reported Ant Design issue: https://github.com/ant-design/ant-design/issues/45481
          // @ts-ignore
          wrapper: modalWrapperStyle,
        }}
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
