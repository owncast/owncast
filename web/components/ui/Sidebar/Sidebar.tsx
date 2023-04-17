import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { FC } from 'react';
import dynamic from 'next/dynamic';
import { Spin } from 'antd';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import styles from './Sidebar.module.scss';

import {
  currentUserAtom,
  knownChatUserDisplayNamesAtom,
  visibleChatMessagesSelector,
  isChatAvailableSelector,
} from '../../stores/ClientConfigStore';

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
  const isChatAvailable = useRecoilValue(isChatAvailableSelector);
  const knownChatUserDisplayNames = useRecoilValue(knownChatUserDisplayNamesAtom);

  if (!currentUser) {
    return (
      <Sider className={styles.root} collapsedWidth={0} width={100}>
        <Spin spinning size="large" />
      </Sider>
    );
  }

  const { id, isModerator, displayName } = currentUser;
  return (
    <Sider className={styles.root} collapsedWidth={0} width={320}>
      <ChatContainer
        messages={messages}
        usernameToHighlight={displayName}
        chatUserId={id}
        isModerator={isModerator}
        chatAvailable={isChatAvailable}
        showInput={!!currentUser}
        knownChatUserDisplayNames={knownChatUserDisplayNames}
      />
    </Sider>
  );
};
