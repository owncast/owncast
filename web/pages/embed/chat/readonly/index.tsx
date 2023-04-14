import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  currentUserAtom,
  visibleChatMessagesSelector,
  isChatAvailableSelector,
} from '../../../../components/stores/ClientConfigStore';

export default function ReadOnlyChatEmbed() {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  const isChatAvailable = useRecoilValue(isChatAvailableSelector);

  return (
    <div>
      <ClientConfigStore />
      {currentUser && (
        <ChatContainer
          messages={messages}
          usernameToHighlight={currentUser.displayName}
          chatUserId={currentUser.id}
          isModerator={false}
          showInput={false}
          height="100vh"
          chatAvailable={isChatAvailable}
        />
      )}
    </div>
  );
}
