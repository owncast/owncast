import { useRecoilValue } from 'recoil';
import Head from 'next/head';
import { useTranslation } from 'next-export-i18n';
import { ErrorBoundary } from 'react-error-boundary';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  currentUserAtom,
  visibleChatMessagesSelector,
  clientConfigStateAtom,
  isChatAvailableSelector,
} from '../../../../components/stores/ClientConfigStore';
import { ClientConfig } from '../../../../interfaces/client-config.model';
import { Theme } from '../../../../components/theme/Theme';
import { ComponentError } from '../../../../components/ui/ComponentError/ComponentError';
import { Localization } from '../../../../types/localization';

export default function ReadOnlyChatEmbed() {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const isChatAvailable = useRecoilValue(isChatAvailableSelector);

  const { name } = clientConfig;
  const { t } = useTranslation();

  const pageTitle = name ? t(Localization.Frontend.chatEmbedTitle, { name }) : 'Chat';

  return (
    <div>
      <Head>
        <title>{pageTitle}</title>
      </Head>
      <ErrorBoundary
        // eslint-disable-next-line react/no-unstable-nested-components
        fallbackRender={({ error }) => (
          <ComponentError componentName="ReadOnlyChatEmbed" message={error.message} />
        )}
      >
        <ClientConfigStore />
        <Theme />
        {currentUser && (
          <ChatContainer
            messages={messages}
            usernameToHighlight={currentUser.displayName}
            chatUserId={currentUser.id}
            isModerator={false}
            showInput={false}
            height="100vh"
            chatAvailable={isChatAvailable}
            readonly
          />
        )}
      </ErrorBoundary>
    </div>
  );
}
