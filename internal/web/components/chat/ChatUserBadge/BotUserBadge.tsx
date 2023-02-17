import dynamic from 'next/dynamic';
import React, { FC } from 'react';
import { ChatUserBadge } from './ChatUserBadge';

// Lazy loaded components

const BulbFilled = dynamic(() => import('@ant-design/icons/BulbFilled'), {
  ssr: false,
});

export type BotBadgeProps = {
  userColor: number;
};

export const BotUserBadge: FC<BotBadgeProps> = ({ userColor }) => (
  <ChatUserBadge badge={<BulbFilled />} userColor={userColor} title="Bot" />
);
