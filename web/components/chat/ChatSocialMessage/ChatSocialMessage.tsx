import { Avatar, Col, Row } from 'antd';
import dynamic from 'next/dynamic';
import React, { FC } from 'react';
import { ChatSocialMessage as ChatMessage } from '../../../interfaces/chat-social-message.model';
import styles from './ChatSocialMessage.module.scss';

const FollowIcon = dynamic(() => import('./follow.svg'));
const LikeIcon = dynamic(() => import('./like.svg'));
const RepostIcon = dynamic(() => import('./repost.svg'));

export interface ChatSocialMessageProps {
  message: ChatMessage;
}

export const ChatSocialMessage: FC<ChatSocialMessageProps> = ({ message }) => {
  const { body, title, image, link, type } = message;

  let Icon;

  switch (type.toString()) {
    case 'follow':
      Icon = FollowIcon;
      break;
    case 'like':
      Icon = LikeIcon;
      break;
    case 'repost':
      Icon = RepostIcon;
      break;
    default:
      break;
  }

  return (
    <div className={styles.follower}>
      <a href={link} target="_blank" rel="noreferrer">
        <Row wrap={false}>
          <Col span={6}>
            <Avatar src={image} alt="Avatar" className={styles.avatar}>
              <img src="/logo" alt="Logo" className={styles.placeholder} />
            </Avatar>
            <Icon className={styles.icon} />
          </Col>
          <Col>
            <Row className={styles.account}>{title}</Row>
            <Row className={styles.body}>{body}</Row>
          </Col>
        </Row>
      </a>
    </div>
  );
};
