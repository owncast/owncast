import React, { FC } from 'react';
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
    <span style={style} className={styles.badge} title={title}>
      {badge}
    </span>
  );
};
