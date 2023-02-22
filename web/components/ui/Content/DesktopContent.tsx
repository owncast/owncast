import React, { ComponentType, FC } from 'react';
import dynamic from 'next/dynamic';
import { TabsProps } from 'antd';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './Content.module.scss';
import { CustomPageContent } from '../CustomPageContent/CustomPageContent';
import { ContentHeader } from '../../common/ContentHeader/ContentHeader';

export type DesktopContentProps = {
  name: string;
  streamTitle: string;
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
  streamTitle,
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

  const items = [{ label: 'About', key: '2', children: aboutTabContent }];
  if (supportFediverseFeatures) {
    items.push({ label: 'Followers', key: '3', children: followersTabContent });
  }

  return (
    <>
      <div className={styles.lowerHalf} id="skip-to-content">
        <ContentHeader
          name={name}
          title={streamTitle}
          summary={summary}
          tags={tags}
          links={socialHandles}
          logo="/logo"
        />
      </div>

      <div className={styles.lowerSection}>
        {items.length > 1 ? <Tabs defaultActiveKey="0" items={items} /> : aboutTabContent}
      </div>
    </>
  );
};
