/* eslint-disable react/no-unused-prop-types */
/* eslint-disable @typescript-eslint/no-unused-vars */
// TODO remove unused props
import { FC } from 'react';
import { ChatMessage } from '../../../interfaces/chat-message.model';

export interface ChatSocialMessageProps {
  message: ChatMessage;
}

export const ChatSocialMessage: FC<ChatSocialMessageProps> = () => <div>Component goes here</div>;
