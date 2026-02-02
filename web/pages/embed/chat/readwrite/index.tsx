/* eslint-disable react/no-unknown-property */
import { useRecoilValue } from 'recoil';
import { useEffect } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { useTranslation } from 'next-export-i18n';
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
  chatAuthenticatedAtom,
} from '../../../../components/stores/ClientConfigStore';
import Header from '../../../../components/ui/Header/Header';
import { ClientConfig } from '../../../../interfaces/client-config.model';
import { AppStateOptions } from '../../../../components/stores/application-state';
import { ServerStatus } from '../../../../interfaces/server-status.model';
import { Theme } from '../../../../components/theme/Theme';
import { ComponentError } from '../../../../components/ui/ComponentError/ComponentError';
import { Localization } from '../../../../types/localization';

export default function ReadWriteChatEmbed() {
  const { t } = useTranslation();
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const clientStatus = useRecoilValue<ServerStatus>(serverStatusState);

  const appState = useRecoilValue<AppStateOptions>(appStateAtom);
  const isChatAvailable = useRecoilValue(isChatAvailableSelector);
  const isUserAuthenticated = useRecoilValue<boolean>(chatAuthenticatedAtom);

  const { name, chatDisabled, chatRequireAuthentication } = clientConfig;

  // Determine if chat input should be enabled based on authentication requirements.
  // Moderators bypass the authentication requirement.
  const chatInputEnabled = !!(
    isChatAvailable &&
    (!chatRequireAuthentication || isUserAuthenticated || currentUser?.isModerator)
  );
  const chatInputDisabledMessage = chatRequireAuthentication
    ? t(Localization.Frontend.Chat.authenticateToChat)
    : t(Localization.Frontend.chatDisabled);
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
      <ErrorBoundary
        // eslint-disable-next-line react/no-unstable-nested-components
        fallbackRender={({ error }) => (
          <ComponentError componentName="ReadWriteChatEmbed" message={error.message} />
        )}
      >
        <ClientConfigStore />
        <Theme />
        <Header
          name={headerText}
          chatAvailable
          chatDisabled={chatDisabled}
          online={videoAvailable}
        />
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
              inputEnabled={chatInputEnabled}
              inputDisabledPlaceholder={chatInputDisabledMessage}
            />
          </div>
        )}
      </ErrorBoundary>
    </div>
  );
}
