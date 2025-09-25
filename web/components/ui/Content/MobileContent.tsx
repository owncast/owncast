import React, { ComponentType, FC } from 'react';
import dynamic from 'next/dynamic';
import { TabsProps } from 'antd';
import { ErrorBoundary } from 'react-error-boundary';
import classNames from 'classnames';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './Content.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ComponentError } from '../ComponentError/ComponentError';

export type MobileContentProps = {
  name: string;
  summary: string;
  tags: string[];
  socialHandles: SocialLink[];
  extraPageContent: string;
  setShowFollowModal: (show: boolean) => void;
  supportFediverseFeatures: boolean;
  online: boolean;
  federatedServers?: any[]; // Will be properly typed when API is implemented
};

// lazy loaded components

const Tabs: ComponentType<TabsProps> = dynamic(() => import('antd').then(mod => mod.Tabs), {
  ssr: false,
});

const FollowerCollection = dynamic(
  () =>
    import('../followers/FollowerCollection/FollowerCollection').then(
      mod => mod.FollowerCollection,
    ),
  {
    ssr: false,
  },
);

const StreamsTab = dynamic(() => import('../StreamsTab/StreamsTab').then(mod => mod.StreamsTab), {
  ssr: false,
});

const ComponentErrorFallback = ({ error, resetErrorBoundary }) => (
  <ComponentError
    message={error}
    componentName="MobileContent"
    retryFunction={resetErrorBoundary}
  />
);

export const MobileContent: FC<MobileContentProps> = ({
  name,
  summary,
  tags,
  socialHandles,
  extraPageContent,
  setShowFollowModal,
  supportFediverseFeatures,
  online,
  federatedServers = [],
}) => {
  const aboutTabContent = (
    <>
      <ContentHeader name={name} summary={summary} tags={tags} links={socialHandles} logo="/logo" />
      {!!extraPageContent && (
        <div className={styles.bottomPageContentContainer}>
          <CustomPageContent content={extraPageContent} />
        </div>
      )}
    </>
  );
  const followersTabContent = (
    <div className={styles.bottomPageContentContainer}>
      <FollowerCollection name={name} onFollowButtonClick={() => setShowFollowModal(true)} />
    </div>
  );

  const streamsTabContent = (
    <div className={styles.bottomPageContentContainer}>
      <StreamsTab servers={federatedServers} />
    </div>
  );

  const items = [];

  items.push({ label: 'About', key: '0', children: aboutTabContent });
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '1', children: followersTabContent });
  }
  // Add Featured tab if there are featured streams
  if (federatedServers && federatedServers.length > 0) {
    items.push({ label: 'Featured', key: '2', children: streamsTabContent });
  }

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentErrorFallback error={error} resetErrorBoundary={resetErrorBoundary} />
      )}
    >
      {items.length > 1 ? (
        <div className={classNames([styles.lowerSectionMobileTabbed, online && styles.online])}>
          <Tabs defaultActiveKey="0" items={items} />
        </div>
      ) : (
        <div>{aboutTabContent}</div>
      )}
    </ErrorBoundary>
  );
};
