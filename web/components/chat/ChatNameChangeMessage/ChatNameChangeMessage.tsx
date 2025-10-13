import dynamic from 'next/dynamic';
import { FC } from 'react';
import { NameChangeEvent } from '../../../interfaces/socket-events';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';
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
        <span> </span>
        <Translation
          translationKey={Localization.Frontend.Chat.nameChangeText}
          className={styles.plain}
          id="owncast-name-change-is-now-known-text"
          defaultText="is now known as"
        />
        <span> </span>
        <span style={{ color }}>{displayName}</span>
      </div>
    </div>
  );
};
