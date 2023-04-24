// export const ChatSocialMessage: FC<ChatSocialMessageProps> = ({ message }) => {

import dynamic from 'next/dynamic';
import { FC } from 'react';
import { NameChangeEvent } from '../../../interfaces/socket-events';
import styles from './ChatNameChangeMessage.module.scss';

export interface ChatNameChangeMessageProps {
  message: NameChangeEvent;
}

// Lazy loaded components

const EditFilled = dynamic(() => import('@ant-design/icons/EditFilled'), {
  ssr: false,
});

export const ChatNameChangeMessage: FC<ChatNameChangeMessageProps> = ({ message }) => {
  const { oldName, user } = message;
  const { displayName, displayColor } = user;
  const color = `var(--theme-color-users-${displayColor})`;

  return (
    <div className={styles.nameChangeView}>
      <div className={styles.icon}>
        <EditFilled />
      </div>
      <div className={styles.nameChangeText}>
        <span style={{ color }}>{oldName}</span>
        <span className={styles.plain}> is now known as </span>
        <span style={{ color }}>{displayName}</span>
      </div>
    </div>
  );
};
