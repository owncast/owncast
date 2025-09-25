import React, { FC, useState, useEffect } from 'react';
import { Row, Col, Empty, Spin, Alert } from 'antd';
import { StreamCard } from '../StreamCard/StreamCard';
import { Translation } from '../Translation/Translation';
import { Localization } from '../../../types/localization';
import styles from './StreamsTab.module.scss';

export interface FederatedServer {
  id: string;
  name: string;
  url: string;
  logo?: string;
  isOnline: boolean;
  streamTitle?: string;
  streamDescription?: string;
  tags?: string[];
  thumbnail?: string;
}

export interface StreamsTabProps {
  servers?: FederatedServer[];
  loading?: boolean;
  error?: string;
}

export const StreamsTab: FC<StreamsTabProps> = ({ servers = [], loading = false, error }) => {
  const [federatedServers, setFederatedServers] = useState<FederatedServer[]>(servers);

  useEffect(() => {
    setFederatedServers(servers);
  }, [servers]);

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <Spin
          size="large"
          tip={
            <Translation
              translationKey={Localization.Frontend.StreamsTab.loadingStreams}
              defaultText="Loading featured streams..."
            />
          }
        />
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <Alert
          message={
            <Translation
              translationKey={Localization.Frontend.StreamsTab.errorLoadingStreams}
              defaultText="Error loading streams"
            />
          }
          description={error}
          type="error"
          showIcon
        />
      </div>
    );
  }

  if (federatedServers.length === 0) {
    return (
      <div className={styles.emptyContainer}>
        <Empty
          description={
            <Translation
              translationKey={Localization.Frontend.StreamsTab.noFeaturedStreams}
              defaultText="No featured streams available"
            />
          }
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      </div>
    );
  }

  // Sort servers: online first, then by name
  const sortedServers = [...federatedServers].sort((a, b) => {
    if (a.isOnline !== b.isOnline) {
      return a.isOnline ? -1 : 1;
    }
    return a.name.localeCompare(b.name);
  });

  return (
    <div className={styles.streamsContainer}>
      <Row gutter={[16, 16]}>
        {sortedServers.map(server => (
          <Col key={server.id} xs={24} sm={12} md={8} lg={6} xl={6} xxl={4}>
            <StreamCard
              serverName={server.name}
              serverUrl={server.url}
              serverLogo={server.logo}
              streamTitle={server.streamTitle}
              streamDescription={server.streamDescription}
              tags={server.tags}
              thumbnail={server.thumbnail}
              isOnline={server.isOnline}
            />
          </Col>
        ))}
      </Row>
    </div>
  );
};
