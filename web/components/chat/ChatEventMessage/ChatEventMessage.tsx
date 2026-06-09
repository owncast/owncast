import { FC } from 'react';
import dynamic from 'next/dynamic';
import { ModerationBadge } from '../ChatUserBadge/ModerationBadge';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';

import styles from './ChatEventMessage.module.scss';

// Lazy loaded components

const UsergroupAddOutlined = dynamic(() => import('@ant-design/icons/UsergroupAddOutlined'), {
  ssr: false,
});

const UsergroupDeleteOutlined = dynamic(() => import('@ant-design/icons/UsergroupDeleteOutlined'), {
  ssr: false,
});

export enum ChatEventType {
  Join = 'join',
  Part = 'part',
}

const variants = {
  [ChatEventType.Join]: {
    Icon: UsergroupAddOutlined,
    translationKey: Localization.Frontend.Chat.userJoined,
    defaultText: 'joined the chat.',
  },
  [ChatEventType.Part]: {
    Icon: UsergroupDeleteOutlined,
    translationKey: Localization.Frontend.Chat.userLeft,
    defaultText: 'left the chat.',
  },
};

export type ChatEventMessageProps = {
  type: ChatEventType;
  isAuthorModerator: boolean;
  userColor: number;
  displayName: string;
};

export const ChatEventMessage: FC<ChatEventMessageProps> = ({
  type,
  isAuthorModerator,
  userColor,
  displayName,
}) => {
  const color = userColor !== undefined ? `var(--theme-color-users-${userColor})` : 'inherit';
  const { Icon, translationKey, defaultText } = variants[type];

  return (
    <div className={styles.root}>
      <span style={{ color }}>
        <span className={styles.icon}>
          <Icon />
        </span>
        {isAuthorModerator && (
          <span className={styles.moderatorBadge}>
            <ModerationBadge userColor={userColor} />
          </span>
        )}
        <span className={styles.message}>
          <span className={styles.user} style={{ color }}>
            {displayName}
          </span>
          <span> </span>
          <Translation translationKey={translationKey} defaultText={defaultText} />
        </span>
      </span>
    </div>
  );
};
