import type { Meta, StoryObj } from '@storybook/nextjs';
import type { FederatedServerResponse } from '../../../hooks/useFederatedServers';
import { OfflineBanner } from '../OfflineBanner/OfflineBanner';
import { OfflineRecommendation } from './OfflineRecommendation';
import styles from './OfflineRecommendation.stories.module.scss';

const featuredServer: FederatedServerResponse = {
  id: 1,
  iri: 'https://radio.gamethattune.com',
  displayName: 'Game That Tune Radio',
  streamDescription: 'Video game deep cuts and listener requests',
  streamTitle: 'Always-on video game music radio',
  logoUrl: 'https://watch.owncast.online/logo',
  thumbnailUrl: 'https://watch.owncast.online/thumbnail.jpg',
  isOnline: true,
  addedAt: new Date().toISOString(),
};

const brightThumbnailServer: FederatedServerResponse = {
  ...featuredServer,
  displayName: 'Daylight Community TV',
  streamDescription: 'A very bright outdoor broadcast',
  iri: 'https://daylight.example',
  thumbnailUrl: `data:image/svg+xml,${encodeURIComponent(`
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 960 540">
      <rect width="960" height="540" fill="#fff9db"/>
      <circle cx="760" cy="120" r="95" fill="#ffffff"/>
      <path d="M0 430 240 180 420 350 620 150 960 440V540H0Z" fill="#dff5ff"/>
      <path d="M0 470 260 300 480 440 700 260 960 470V540H0Z" fill="#f2f8df"/>
    </svg>
  `)}`,
};

const longMetadataServer: FederatedServerResponse = {
  ...featuredServer,
  displayName: 'The International Community Game Preservation and Independent Radio Collective',
  streamDescription:
    'Restoring forgotten games live while taking listener requests and answering community questions',
  streamTitle:
    'The complete international history of independent games, preserved live with community guests',
  iri: 'https://live.broadcast.archive.community-media.gamethattune.example',
};

const failedThumbnailServer: FederatedServerResponse = {
  ...featuredServer,
  thumbnailUrl: '/missing-featured-thumbnail.jpg',
};

const lastLive = new Date(Date.now() - 3 * 24 * 60 * 60 * 1000);

const OfflineRecommendationStory = ({
  server = featuredServer,
}: {
  server?: FederatedServerResponse | null;
}) => (
  <div className={styles.viewerFrame}>
    <header className={styles.viewerHeader}>
      <strong>Gabek&apos;s Owncast</strong>
      <span>gabek</span>
    </header>
    <main className={styles.offlineStage}>
      <div className={server ? styles.splitPanel : styles.offlineOnly}>
        <OfflineBanner
          showsHeader={false}
          streamName="Gabek's Owncast"
          customText="No stream today. I'll be back Thursday for a live Q&amp;A and project update."
          notificationsEnabled={false}
          lastLive={lastLive}
        />
        {server && <OfflineRecommendation localStreamName="Gabek's Owncast" server={server} />}
      </div>
    </main>
    <section className={styles.viewerContent}>
      <nav className={styles.viewerTabs} aria-label="Viewer sections">
        <span className={styles.activeTab}>About</span>
        <span>Schedule</span>
        <span>Featured</span>
      </nav>
      <div className={styles.viewerSummary}>
        <h2>Gabek&apos;s Owncast</h2>
        <p>Independent live video, community updates, and project streams.</p>
      </div>
    </section>
  </div>
);

const meta = {
  title: 'owncast/Layout/Offline Recommendation Concepts',
  component: OfflineRecommendation,
  args: {
    localStreamName: "Gabek's Owncast",
    server: featuredServer,
  },
  parameters: {
    layout: 'fullscreen',
    chromatic: { viewports: [375, 1280] },
    docs: {
      description: {
        component:
          'The responsive offline viewer layout with one live featured-stream recommendation.',
      },
    },
  },
} satisfies Meta<typeof OfflineRecommendation>;

export default meta;
type Story = StoryObj<typeof meta>;

export const SplitPanel: Story = {
  render: () => <OfflineRecommendationStory />,
};

export const SplitPanelWithoutLiveFeaturedStream: Story = {
  render: () => <OfflineRecommendationStory server={null} />,
};

export const SplitPanelBrightThumbnail: Story = {
  render: () => <OfflineRecommendationStory server={brightThumbnailServer} />,
};

export const SplitPanelLongMetadata: Story = {
  render: () => <OfflineRecommendationStory server={longMetadataServer} />,
};

export const SplitPanelFailedThumbnail: Story = {
  render: () => <OfflineRecommendationStory server={failedThumbnailServer} />,
};
