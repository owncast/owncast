import React, { ComponentType, FC } from 'react';
import dynamic from 'next/dynamic';
import { TabsProps } from 'antd';
import { ErrorBoundary } from 'react-error-boundary';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './Content.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ComponentError } from '../ComponentError/ComponentError';

export type DesktopContentProps = {
  name: string;
  summary: string;
  tags: string[];
  socialHandles: SocialLink[];
  extraPageContent: string;
  setShowFollowModal: (show: boolean) => void;
  supportFediverseFeatures: boolean;
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

export const DesktopContent: FC<DesktopContentProps> = ({
  name,
  summary,
  tags,
  socialHandles,
  extraPageContent,
  setShowFollowModal,
  supportFediverseFeatures,
  federatedServers = [],
}) => {
  const aboutTabContent = (
    <div className={styles.bottomPageContentContainer}>
      <CustomPageContent content={extraPageContent} />
    </div>
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

  const items = [!!extraPageContent && { label: 'About', key: '2', children: aboutTabContent }];
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '3', children: followersTabContent });
  }
  // Add Featured tab if there are featured streams
  if (federatedServers && federatedServers.length > 0) {
    items.push({ label: 'Featured Streams', key: '4', children: streamsTabContent });
  }

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="DesktopContent"
          message={error.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <div id="skip-to-content">
        <ContentHeader
          name={name}
          summary={summary}
          tags={tags}
          links={socialHandles}
          logo="/logo"
        />
      </div>

      <div>
        {items.length > 1 ? (
          <Tabs defaultActiveKey="0" items={items} />
        ) : (
          !!extraPageContent && aboutTabContent
        )}
      </div>
    </ErrorBoundary>
  );
};
