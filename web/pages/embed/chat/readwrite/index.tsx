import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  chatDisplayNameAtom,
  chatUserIdAtom,
  visibleChatMessagesSelector,
  clientConfigStateAtom,
  isChatModeratorAtom,
} from '../../../../components/stores/ClientConfigStore';
import Header from '../../../../components/ui/Header/Header';
import { ClientConfig } from '../../../../interfaces/client-config.model';

export default function ReadWriteChatEmbed() {
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const chatUserId = useRecoilValue<string>(chatUserIdAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const isModerator = useRecoilValue<boolean>(isChatModeratorAtom);

  const { name, chatDisabled } = clientConfig;

  return (
    <div>
      <ClientConfigStore />
      <Header name={name} chatAvailable chatDisabled={chatDisabled} />
      <ChatContainer
        messages={messages}
        usernameToHighlight={chatDisplayName}
        chatUserId={chatUserId}
        isModerator={isModerator}
        showInput
        height="80vh"
      />
    </div>
  );
}
