/* eslint-disable react/no-danger */
import { Highlight } from 'react-highlighter-ts';
import { FC } from 'react';
import { ChatMessage } from '~/interfaces/chat-message.model';
import styles from './ChatSystemMessage.module.scss';

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
  <div className={styles.chatSystemMessage}>
    <div className={styles.user}>
      <span className={styles.userName}>{displayName}</span>
    </div>
    <Highlight search={highlightString}>
      <div className={styles.message} dangerouslySetInnerHTML={{ __html: body }} />
    </Highlight>
  </div>
);
