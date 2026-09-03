import { ComponentType, FC, useState } from 'react';
import dynamic from 'next/dynamic';
import { Button, Dropdown } from 'antd';
import type { TabsProps } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { ErrorBoundary } from 'react-error-boundary';
import classNames from 'classnames';
import { SocialLink } from '../../../interfaces/social-link.model';
import { PluginTab } from '../../../interfaces/client-config.model';
import type { FederatedServer } from '../StreamsTab/StreamsTab';
import { Localization } from '../../../types/localization';
import styles from './Content.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { PluginTabFrame } from '../PluginTabFrame/PluginTabFrame';
import { Translation } from '../Translation/Translation';
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
  showFollowersTab: boolean;
  scheduleEnabled?: boolean;
  online: boolean;
  federatedServers?: FederatedServer[];
};

// lazy loaded components

const Tabs: ComponentType<TabsProps> = dynamic(() => import('antd').then(mod => mod.Tabs), {
  ssr: false,
});

const UnorderedListOutlined = dynamic(() => import('@ant-design/icons/UnorderedListOutlined'), {
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

const ScheduleTab = dynamic(
  () => import('../ScheduleTab/ScheduleTab').then(mod => mod.ScheduleTab),
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
  pluginTabs,
  setShowFollowModal,
  showFollowersTab,
  online,
  scheduleEnabled = false,
  federatedServers = [],
}) => {
  const { t } = useTranslation();
  const [activeKey, setActiveKey] = useState('0');

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

  const scheduleTabContent = (
    <div className={styles.bottomPageContentContainer}>
      <ScheduleTab />
    </div>
  );

  const streamsTabContent = (
    <div className={styles.bottomPageContentContainer}>
      <StreamsTab servers={federatedServers} />
    </div>
  );

  const items: NonNullable<TabsProps['items']> = [];

  items.push({ label: 'About', key: '0', children: aboutTabContent });
  if (showFollowersTab) {
    items.push({ label: 'Followers', key: '1', children: followersTabContent });
  }
  if (scheduleEnabled) {
    items.push({
      label: (
        <Translation translationKey={Localization.Frontend.Schedule.tab} defaultText="Schedule" />
      ),
      key: 'schedule',
      children: scheduleTabContent,
    });
  }
  // Add Featured tab if there are featured streams
  if (federatedServers && federatedServers.length > 0) {
    items.push({ label: 'Featured', key: '2', children: streamsTabContent });
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

  // Menu listing every tab so any of them stays reachable when the tab
  // bar overflows the screen. rc-tabs turns its own "more" dropdown off
  // on mobile devices, leaving swipe-scrolling as the only way to reach
  // far tabs. This button is shown via CSS only while the tab bar
  // actually overflows (the nav-wrap ping classes).
  const translatedShowAllTabs = t(Localization.Frontend.showAllTabs);
  const showAllTabsLabel =
    translatedShowAllTabs === Localization.Frontend.showAllTabs
      ? 'Show all tabs'
      : translatedShowAllTabs;

  const allTabsMenu = (
    <Dropdown
      menu={{
        items: items.map(({ key, label }) => ({ key, label })),
        onClick: ({ key }) => setActiveKey(key),
        selectedKeys: [activeKey],
      }}
      placement="bottomRight"
      trigger={['click']}
    >
      <Button type="text" icon={<UnorderedListOutlined />} aria-label={showAllTabsLabel} />
    </Dropdown>
  );

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentErrorFallback error={error} resetErrorBoundary={resetErrorBoundary} />
      )}
    >
      {items.length > 1 ? (
        <div className={classNames([styles.lowerSectionMobileTabbed, online && styles.online])}>
          <Tabs
            activeKey={activeKey}
            onChange={setActiveKey}
            items={items}
            tabBarExtraContent={{ right: allTabsMenu }}
          />
        </div>
      ) : (
        <div style={{ width: '100%' }}>{aboutTabContent}</div>
      )}
    </ErrorBoundary>
  );
};
