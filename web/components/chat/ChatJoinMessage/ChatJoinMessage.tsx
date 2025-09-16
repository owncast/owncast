import { FC } from 'react';
import dynamic from 'next/dynamic';
import { ModerationBadge } from '../ChatUserBadge/ModerationBadge';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';

import styles from './ChatJoinMessage.module.scss';

// Lazy loaded components

const UsergroupAddOutlined = dynamic(() => import('@ant-design/icons/UsergroupAddOutlined'), {
  ssr: false,
});

export type ChatJoinMessageProps = {
  isAuthorModerator: boolean;
  userColor: number;
  displayName: string;
};

export const ChatJoinMessage: FC<ChatJoinMessageProps> = ({
  isAuthorModerator,
  userColor,
  displayName,
}) => {
  const color = `var(--theme-color-users-${userColor})`;

  return (
    <div className={styles.root}>
      <span style={{ color }}>
        <span className={styles.icon}>
          <UsergroupAddOutlined />
        </span>
        {isAuthorModerator && (
          <span className={styles.moderatorBadge}>
            <ModerationBadge userColor={userColor} />
          </span>
        )}
        <span className={styles.joinMessage}>
          <Translation
            translationKey={Localization.Frontend.Chat.userJoined}
            vars={{
              name: `<span class="${styles.user}">${displayName}</span>`,
            }}
            defaultText="{{name}} joined the chat."
          />
        </span>
      </span>
    </div>
  );
};
