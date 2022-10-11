import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { FC } from 'react';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatContainer } from '../../chat/ChatContainer/ChatContainer';
import styles from './Sidebar.module.scss';

import { currentUserAtom, visibleChatMessagesSelector } from '../../stores/ClientConfigStore';

export const Sidebar: FC = () => {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);

  if (!currentUser) {
    return null;
  }

  const { id, isModerator, displayName } = currentUser;
  return (
    <Sider className={styles.root} collapsedWidth={0} width={320}>
      <ChatContainer
        messages={messages}
        usernameToHighlight={displayName}
        chatUserId={id}
        isModerator={isModerator}
      />
    </Sider>
  );
};
