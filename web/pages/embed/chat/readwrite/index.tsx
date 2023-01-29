import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  currentUserAtom,
  visibleChatMessagesSelector,
  clientConfigStateAtom,
  appStateAtom,
} from '../../../../components/stores/ClientConfigStore';
import Header from '../../../../components/ui/Header/Header';
import { ClientConfig } from '../../../../interfaces/client-config.model';
import { AppStateOptions } from '../../../../components/stores/application-state';

export default function ReadWriteChatEmbed() {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);

  const { name, chatDisabled } = clientConfig;
  const { videoAvailable } = appState;

  if (!currentUser) {
    return null;
  }

  const { id, displayName, isModerator } = currentUser;

  return (
    <div>
      <ClientConfigStore />
      <Header name={name} chatAvailable chatDisabled={chatDisabled} online={videoAvailable} />
      <ChatContainer
        messages={messages}
        usernameToHighlight={displayName}
        chatUserId={id}
        isModerator={isModerator}
        showInput
        height="80vh"
      />
    </div>
  );
}
