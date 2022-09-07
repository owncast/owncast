import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  chatDisplayNameAtom,
  chatUserIdAtom,
  visibleChatMessagesSelector,
} from '../../../../components/stores/ClientConfigStore';

export default function ReadOnlyChatEmbed() {
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const chatUserId = useRecoilValue<string>(chatUserIdAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);

  return (
    <div>
      <ClientConfigStore />
      <ChatContainer
        messages={messages}
        usernameToHighlight={chatDisplayName}
        chatUserId={chatUserId}
        isModerator={false}
        showInput={false}
        height="100vh"
      />
    </div>
  );
}
