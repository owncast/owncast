import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  currentUserAtom,
  visibleChatMessagesSelector,
} from '../../../../components/stores/ClientConfigStore';

export default function ReadOnlyChatEmbed() {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  if (!currentUser) {
    return null;
  }
  const { id, displayName } = currentUser;
  return (
    <div>
      <ClientConfigStore />
      <ChatContainer
        messages={messages}
        usernameToHighlight={displayName}
        chatUserId={id}
        isModerator={false}
        showInput={false}
        height="100vh"
      />
    </div>
  );
}
