import { FC } from 'react';
import dynamic from 'next/dynamic';
import { ModerationBadge } from '../ChatUserBadge/ModerationBadge';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';

import styles from './ChatPartMessage.module.scss';

// Lazy loaded components

const UsergroupDeleteOutlined = dynamic(() => import('@ant-design/icons/UsergroupDeleteOutlined'), {
  ssr: false,
});

export type ChatPartMessageProps = {
  isAuthorModerator: boolean;
  userColor: number;
  displayName: string;
};

export const ChatPartMessage: FC<ChatPartMessageProps> = ({
  isAuthorModerator,
  userColor,
  displayName,
}) => {
  const color = `var(--theme-color-users-${userColor})`;

  return (
    <div className={styles.root}>
      <span style={{ color }}>
        <span className={styles.icon}>
          <UsergroupDeleteOutlined />
        </span>
        {isAuthorModerator && (
          <span className={styles.moderatorBadge}>
            <ModerationBadge userColor={userColor} />
          </span>
        )}
        <span className={styles.partMessage}>
          <Translation
            translationKey={Localization.Frontend.Chat.userLeft}
            vars={{
              name: displayName,
              userClass: styles.user,
            }}
            defaultText='<span class="{{userClass}}">{{name}}</span> left the chat.'
          />
        </span>
      </span>
    </div>
  );
};
