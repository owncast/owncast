import React, { FC } from 'react';
import styles from './ChatUserBadge.module.scss';

export type ChatUserBadgeProps = {
  badge: React.ReactNode;
  userColor: number;
};

export const ChatUserBadge: FC<ChatUserBadgeProps> = ({ badge, userColor }) => {
  const color = `var(--theme-user-colors-${userColor})`;
  const style = { color, borderColor: color };

  return (
    <span style={style} className={styles.badge}>
      {badge}
    </span>
  );
};
