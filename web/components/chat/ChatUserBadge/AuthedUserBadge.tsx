import dynamic from 'next/dynamic';
import React, { FC } from 'react';
import { ChatUserBadge } from './ChatUserBadge';

// Lazy loaded components

const SafetyCertificateFilled = dynamic(() => import('@ant-design/icons/SafetyCertificateFilled'), {
  ssr: false,
});

export type AuthedUserBadgeProps = {
  userColor: number;
};

export const AuthedUserBadge: FC<AuthedUserBadgeProps> = ({ userColor }) => (
  <ChatUserBadge badge={<SafetyCertificateFilled />} userColor={userColor} title="Authenticated" />
);
