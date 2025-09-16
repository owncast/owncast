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
        <Translation
          translationKey={Localization.Frontend.Chat.nameChangeText}
          vars={{
            name: `<span style="color: ${color}">${oldName}</span>`,
            newName: `<span style="color: ${color}">${displayName}</span>`,
          }}
          defaultText="{{name}} is now known as {{newName}}"
        />
      </div>
    </div>
  );
};
