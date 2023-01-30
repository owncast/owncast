import dynamic from 'next/dynamic';
import React, { FC } from 'react';
import { ChatUserBadge } from './ChatUserBadge';

// Lazy loaded components

const StarFilled = dynamic(() => import('@ant-design/icons/StarFilled'), {
  ssr: false,
});

export type ModerationBadgeProps = {
  userColor: number;
};

export const ModerationBadge: FC<ModerationBadgeProps> = ({ userColor }) => (
  <ChatUserBadge badge={<StarFilled />} userColor={userColor} title="Moderator" />
);
