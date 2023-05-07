import { useRecoilValue } from 'recoil';
import { ErrorBoundary } from 'react-error-boundary';
import { ChatMessage } from '../../../../interfaces/chat-message.model';
import { ChatContainer } from '../../../../components/chat/ChatContainer/ChatContainer';
import {
  ClientConfigStore,
  currentUserAtom,
  visibleChatMessagesSelector,
  isChatAvailableSelector,
} from '../../../../components/stores/ClientConfigStore';
import { Theme } from '../../../../components/theme/Theme';
import { ComponentError } from '../../../../components/ui/ComponentError/ComponentError';

export default function ReadOnlyChatEmbed() {
  const currentUser = useRecoilValue(currentUserAtom);
  const messages = useRecoilValue<ChatMessage[]>(visibleChatMessagesSelector);
  const isChatAvailable = useRecoilValue(isChatAvailableSelector);

  return (
    <div>
      <ErrorBoundary
        // eslint-disable-next-line react/no-unstable-nested-components
        fallbackRender={({ error }) => (
          <ComponentError componentName="ReadWriteChatEmbed" message={error.message} />
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
          />
        )}
      </ErrorBoundary>
    </div>
  );
}
