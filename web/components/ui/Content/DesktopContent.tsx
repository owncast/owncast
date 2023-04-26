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

export const DesktopContent: FC<DesktopContentProps> = ({
  name,
  summary,
  tags,
  socialHandles,
  extraPageContent,
  setShowFollowModal,
  supportFediverseFeatures,
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
  const items = [!!extraPageContent && { label: 'About', key: '2', children: aboutTabContent }];
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '3', children: followersTabContent });
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
