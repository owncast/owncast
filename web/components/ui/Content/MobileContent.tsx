import React, { ComponentType, FC } from 'react';
import dynamic from 'next/dynamic';
import { TabsProps } from 'antd';
import { ErrorBoundary } from 'react-error-boundary';
import classNames from 'classnames';
import { SocialLink } from '../../../interfaces/social-link.model';
import { PluginTab } from '../../../interfaces/client-config.model';
import styles from './Content.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { PluginTabFrame } from '../PluginTabFrame/PluginTabFrame';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';
import { ComponentError } from '../ComponentError/ComponentError';

export type MobileContentProps = {
  name: string;
  summary: string;
  tags: string[];
  socialHandles: SocialLink[];
  extraPageContent: string;
  pluginTabs: PluginTab[];
  setShowFollowModal: (show: boolean) => void;
  supportFediverseFeatures: boolean;
  online: boolean;
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
  pluginTabs,
  setShowFollowModal,
  supportFediverseFeatures,
  online,
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

  const items: NonNullable<TabsProps['items']> = [];

  items.push({ label: 'About', key: '0', children: aboutTabContent });
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '1', children: followersTabContent });
  }
  // Plugin-contributed tabs render after the built-ins. Key is the
  // slug+title combination; the host's validator rejects duplicate
  // titles within a plugin, so this pair is unique across all
  // plugin tabs and stable across renders (no index-as-key
  // anti-pattern).
  //
  // forceRender mounts each tab's iframe up front instead of on first
  // activation, so the srcdoc loads and the host injects styles while the
  // pane is still hidden — the content is ready (no load flash) the moment
  // the user taps the tab.
  (pluginTabs || []).forEach(tab => {
    items.push({
      label: tab.title,
      key: `plugin-${tab.slug}-${tab.title}`,
      forceRender: true,
      children: (
        <div className={styles.bottomPageContentContainer}>
          <PluginTabFrame content={tab.html} />
        </div>
      ),
    });
  });

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
