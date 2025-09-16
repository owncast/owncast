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
  const { displayName } = user;

  return (
    <div className={styles.nameChangeView}>
      <div className={styles.icon}>
        <EditFilled />
      </div>
      <div className={styles.nameChangeText}>
        <Translation
          translationKey={Localization.Frontend.Chat.nameChangeText}
          vars={{ name: oldName, newName: displayName }}
          defaultText="{{name}} is now known as {{newName}}"
        />
      </div>
    </div>
  );
};
