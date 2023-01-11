import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { FC } from 'react';
import dynamic from 'next/dynamic';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import styles from './Sidebar.module.scss';

import { currentUserAtom, visibleChatMessagesSelector } from '../../stores/ClientConfigStore';

// Lazy loaded components
const ChatContainer = dynamic(
  () => import('../../chat/ChatContainer/ChatContainer').then(mod => mod.ChatContainer),
  {
    ssr: false,
  },
);

export const Sidebar: FC = () => {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  if (!currentUser) {
    return <Sider className={styles.root} collapsedWidth={0} width={320} />;
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
