import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  currentUserAtom,
  visibleChatMessagesSelector,
  clientConfigStateAtom,
  appStateAtom,
  serverStatusState,
} from '../../../../components/stores/ClientConfigStore';
import Header from '../../../../components/ui/Header/Header';
import { ClientConfig } from '../../../../interfaces/client-config.model';
import { AppStateOptions } from '../../../../components/stores/application-state';
import { ServerStatus } from '../../../../interfaces/server-status.model';

export default function ReadWriteChatEmbed() {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const clientStatus = useRecoilValue<ServerStatus>(serverStatusState);

  const appState = useRecoilValue<AppStateOptions>(appStateAtom);

  const { name, chatDisabled } = clientConfig;
  const { videoAvailable } = appState;
  const { streamTitle, online } = clientStatus;

  const headerText = online ? streamTitle || name : name;

  return (
    <div>
      <ClientConfigStore />
      <Header name={headerText} chatAvailable chatDisabled={chatDisabled} online={videoAvailable} />
      {currentUser && (
        <div id="chat-container">
          <ChatContainer
            messages={messages}
            usernameToHighlight={currentUser.displayName}
            chatUserId={currentUser.id}
            isModerator={currentUser.isModerator}
            showInput
            height="80vh"
          />
        </div>
      )}
    </div>
  );
}
