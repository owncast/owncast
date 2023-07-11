import { Avatar, Col, Row } from 'antd';
import dynamic from 'next/dynamic';
import React, { FC } from 'react';
import cn from 'classnames';
import { ChatSocialMessage as ChatMessage } from '../../../interfaces/chat-social-message.model';
import styles from './ChatSocialMessage.module.scss';

const FollowIcon = dynamic(() => import('./follow.svg'), {
  ssr: false,
});
const LikeIcon = dynamic(() => import('./like.svg'), {
  ssr: false,
});
const RepostIcon = dynamic(() => import('./repost.svg'), {
  ssr: false,
});

export interface ChatSocialMessageProps {
  message: ChatMessage;
}

export const ChatSocialMessage: FC<ChatSocialMessageProps> = ({ message }) => {
  const { body, title, image, link, type } = message;

  let Icon;

  switch (type.toString()) {
    case 'FEDIVERSE_ENGAGEMENT_FOLLOW':
      Icon = FollowIcon;
      break;
    case 'FEDIVERSE_ENGAGEMENT_LIKE':
      Icon = LikeIcon;
      break;
    case 'FEDIVERSE_ENGAGEMENT_REPOST':
      Icon = RepostIcon;
      break;
    default:
      break;
  }

  return (
    <div className={styles.followerPadding}>
      <div className={cn([styles.follower, 'chat-message_social'])}>
        <a href={link} target="_blank" rel="noreferrer">
          <Row wrap={false}>
            <Col span={6} className={styles.avatarColumn}>
              <Avatar src={image} alt="Avatar" className={styles.avatar} size="large">
                {title.charAt(0).toUpperCase()}
              </Avatar>

              <Icon className={styles.icon} />
            </Col>
            <Col>
              <Row className={styles.account}>{title}</Row>
              <Row className={styles.body} dangerouslySetInnerHTML={{ __html: body }} />
            </Col>
          </Row>
        </a>
      </div>
    </div>
  );
};
