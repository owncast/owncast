import { FC } from 'react';
import cn from 'classnames';
import { Interweave } from 'interweave';
import { UrlMatcher } from 'interweave-autolink';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import styles from './ChatSystemMessage.module.scss';
import { ChatMessageHighlightMatcher } from '../ChatUserMessage/customMatcher';

export type ChatSystemMessageProps = {
  message: ChatMessage;
  highlightString: string;
};

export const ChatSystemMessage: FC<ChatSystemMessageProps> = ({
  message: {
    body,
    user: { displayName },
  },
  highlightString,
}) => (
  <div className={styles.chatSystemMessagePadding}>
    <div className={cn([styles.chatSystemMessage, 'chat-message_system'])}>
      <div className={styles.user}>
        <span className={styles.userName}>{displayName}</span>
      </div>
      <Interweave
        className={styles.message}
        content={body}
        matchers={[
          new UrlMatcher('url', { customTLDs: ['online'] }),
          new ChatMessageHighlightMatcher('highlight', { highlightString }),
        ]}
      />
    </div>
  </div>
);
