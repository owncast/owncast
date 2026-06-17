import React, { FC } from 'react';
import { Card, Tag, Typography, Badge } from 'antd';
import classNames from 'classnames';
import styles from './StreamCard.module.scss';

const { Text, Paragraph } = Typography;

export interface StreamCardProps {
  serverName: string;
  serverUrl: string;
  serverLogo?: string;
  streamTitle?: string;
  streamDescription?: string;
  tags?: string[];
  thumbnail?: string;
  isOnline: boolean;
  onClick?: () => void;
}

export const StreamCard: FC<StreamCardProps> = ({
  serverName,
  serverUrl,
  serverLogo,
  streamTitle,
  streamDescription,
  tags = [],
  thumbnail,
  isOnline,
  onClick,
}) => {
  const handleClick = () => {
    if (onClick) {
      onClick();
    } else {
      window.open(serverUrl, '_blank', 'noopener,noreferrer');
    }
  };

  // The server-supplied name/logo are not trustworthy (a remote server can
  // call itself anything), but the link target is the immutable URL the admin
  // vetted. Surface its hostname so visitors can see where a card actually
  // leads rather than relying on the spoofable display name.
  const serverHost = (() => {
    try {
      return new URL(serverUrl).hostname;
    } catch {
      return serverUrl;
    }
  })();

  const cardCover = (
    <div className={styles.coverContainer}>
      {thumbnail ? (
        <img alt={serverName} src={thumbnail} className={styles.thumbnail} />
      ) : (
        <div className={styles.placeholderThumbnail}>
          {serverLogo && <img src={serverLogo} alt={serverName} className={styles.logoOverlay} />}
        </div>
      )}
      <Badge
        status={isOnline ? 'success' : 'default'}
        text={isOnline ? 'LIVE' : 'OFFLINE'}
        className={styles.statusBadge}
      />
    </div>
  );

  const cardDescription = (
    <div className={styles.cardContent}>
      <div className={styles.serverInfo}>
        {serverLogo && <img src={serverLogo} alt={serverName} className={styles.serverLogo} />}
        <div className={styles.textInfo}>
          <Text strong className={styles.serverName}>
            {serverName}
          </Text>
          <Text type="secondary" className={styles.serverHost} ellipsis>
            {serverHost}
          </Text>
          {isOnline && streamTitle && (
            <Text className={styles.streamTitle} ellipsis>
              {streamTitle}
            </Text>
          )}
        </div>
      </div>
      {isOnline && streamDescription && (
        <Paragraph ellipsis={{ rows: 2 }} className={styles.description}>
          {streamDescription}
        </Paragraph>
      )}
      {tags.length > 0 && (
        <div className={styles.tags}>
          {tags.slice(0, 3).map(tag => (
            <Tag key={tag} className={styles.tag}>
              {tag}
            </Tag>
          ))}
        </div>
      )}
    </div>
  );

  return (
    <Card
      hoverable
      role="article"
      className={classNames(styles.streamCard, {
        [styles.online]: isOnline,
        [styles.offline]: !isOnline,
      })}
      cover={cardCover}
      onClick={handleClick}
      bodyStyle={{ padding: 0 }}
    >
      {cardDescription}
    </Card>
  );
};
