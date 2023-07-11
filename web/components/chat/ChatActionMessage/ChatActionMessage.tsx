import { FC } from 'react';
import styles from './ChatActionMessage.module.scss';

/* eslint-disable react/no-danger */
export type ChatActionMessageProps = {
  body: string;
};

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const ChatActionMessage: FC<ChatActionMessageProps> = ({ body }) => (
  <div className={styles.chatActionPadding}>
    <div dangerouslySetInnerHTML={{ __html: body }} className={styles.chatAction} />
  </div>
);
