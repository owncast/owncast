/* eslint-disable react/no-danger */
import { Highlight } from 'react-highlighter-ts';
import { FC } from 'react';
import cn from 'classnames';
import { ChatMessage } from '../../../interfaces/chat-message.model';
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
  <div className={cn([styles.chatSystemMessage, 'chat-message_system'])}>
    <div className={styles.user}>
      <span className={styles.userName}>{displayName}</span>
    </div>
    <Highlight search={highlightString}>
      <div className={styles.message} dangerouslySetInnerHTML={{ __html: body }} />
    </Highlight>
  </div>
);
