/* eslint-disable react/no-unknown-property */
import { useRecoilValue } from 'recoil';
import { useEffect } from 'react';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  currentUserAtom,
  visibleChatMessagesSelector,
  clientConfigStateAtom,
  appStateAtom,
  serverStatusState,
  isChatAvailableSelector,
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
  const isChatAvailable = useRecoilValue(isChatAvailableSelector);

  const { name, chatDisabled } = clientConfig;
  const { videoAvailable } = appState;
  const { streamTitle, online } = clientStatus;

  const headerText = online ? streamTitle || name : name;

  // This is a hack to force a specific body background color for just this page.
  useEffect(() => {
    document.body.classList.add('body-background');
  }, []);

  return (
    <div>
      <style jsx global>
        {`
          .body-background {
            background: var(--theme-color-components-chat-background);
          }
        `}
      </style>
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
            height="92vh"
            chatAvailable={isChatAvailable}
          />
        </div>
      )}
    </div>
  );
}
