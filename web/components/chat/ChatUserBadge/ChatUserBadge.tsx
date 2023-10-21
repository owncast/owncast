import React, { FC } from 'react';
import cn from 'classnames';
import styles from './ChatUserBadge.module.scss';

export type ChatUserBadgeProps = {
  badge: React.ReactNode;
  userColor: number;
  title: string;
};

export const ChatUserBadge: FC<ChatUserBadgeProps> = ({ badge, userColor, title }) => {
  const color = `var(--theme-color-users-${userColor})`;
  const style = { color };

  return (
    <span style={style} className={cn([styles.badge, 'chat-user-badge'])} title={title}>
      {badge}
    </span>
  );
};
