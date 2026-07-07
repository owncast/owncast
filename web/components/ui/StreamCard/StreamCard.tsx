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

  // Show the live stream thumbnail only while the server is live; when offline
  // (or if no thumbnail is available) fall back to the server logo. Gating on
  // isOnline rather than thumbnail presence avoids showing a stale preview for
  // a server that has since gone offline.
  const cardCover = (
    <div className={styles.coverContainer}>
      {isOnline && thumbnail ? (
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
      {streamDescription && (
        <Paragraph ellipsis={{ rows: 3 }} className={styles.description}>
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

  // Render the card as a real anchor so the destination is visible on hover
  // (status bar), supports middle-click / open-in-new-tab, and is keyboard
  // focusable. The link target is the vetted, immutable server URL.
  return (
    <a
      className={styles.cardLink}
      href={serverUrl}
      target="_blank"
      rel="noopener noreferrer"
      onClick={onClick}
    >
      <Card
        hoverable
        role="article"
        className={classNames(styles.streamCard, {
          [styles.online]: isOnline,
          [styles.offline]: !isOnline,
        })}
        cover={cardCover}
        bodyStyle={{ padding: 0 }}
      >
        {cardDescription}
      </Card>
    </a>
  );
};
