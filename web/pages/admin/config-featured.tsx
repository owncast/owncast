/* eslint-disable react/no-unescaped-entities */
import { Typography, Alert, Button, Space, Tabs, Badge } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import Link from 'next/link';
import React, { ReactElement, useContext, useState } from 'react';
import { Translation } from '../../components/ui/Translation/Translation';
import { Localization } from '../../types/localization';
import { AdminLayout } from '../../components/layouts/AdminLayout';
import { FederatedServersTable } from '../../components/admin/FederatedServers/FederatedServersTable';
import { FeatureStreamModal } from '../../components/admin/FederatedServers/FeatureStreamModal';
import { FeatureRequests } from '../../components/admin/FederatedServers/FeatureRequests';
import { DirectoryListings } from '../../components/admin/FederatedServers/DirectoryListings';
import { useFederatedServers } from '../../hooks/useFederatedServers';
import { useFeatureRequests } from '../../hooks/useFeatureRequests';
import { useDirectoryFollowers } from '../../hooks/useDirectoryFollowers';
import { ServerStatusContext } from '../../utils/server-status-context';

const ConfigFeatured = () => {
  const { Title, Paragraph } = Typography;
  const [modalOpen, setModalOpen] = useState(false);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};
  const { federation } = serverConfig || {};
  const { enabled: federationEnabled } = federation || {};

  const {
    servers: federatedServers,
    loading: serversLoading,
    addServer,
    removeServer,
  } = useFederatedServers(true);

  const {
    requests: featureRequests,
    loading: requestsLoading,
    approve: approveFeatureRequest,
    reject: rejectFeatureRequest,
  } = useFeatureRequests();

  const {
    directories: directoryFollowers,
    loading: directoriesLoading,
    remove: removeDirectoryFollower,
  } = useDirectoryFollowers();

  const handleFeatureStream = async (url: string) => {
    await addServer(url);
    setModalOpen(false);
  };

  // Counts shown in the tab labels so the operator can see the size of each
  // list at a glance without opening it. "Featuring you" combines the pending
  // requests and the directories already listing this server, since both live
  // under that tab.
  const featuredCount = federatedServers.length;
  const pendingRequestCount = featureRequests.length;
  const featuringYouCount = pendingRequestCount + directoryFollowers.length;

  // Neutral count chip so it doesn't compete with the red attention dot used to
  // flag pending requests below.
  const countBadgeStyle = { backgroundColor: '#8c8c8c' };

  return (
    <div>
      <Title>
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.pageTitle}
          defaultText="Featured Streams"
        />
      </Title>
      <Paragraph>
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.pageDescription}
          defaultText='Feature other Owncast streams to display their streaming status on your server. When enabled, visitors can discover and navigate to other featured streams through a dedicated "Featured" tab on your main page.'
        />
      </Paragraph>
      <Paragraph>
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.pageDescriptionSecondary}
          defaultText="This creates a network effect, allowing users to discover and navigate to other Owncast servers easily. It also gives visitors somewhere to go if your stream is offline."
        />
      </Paragraph>

      {federationEnabled ? (
        <>
          <Tabs
            defaultActiveKey="featuring"
            items={[
              {
                key: 'featuring',
                label: (
                  <Space size={8}>
                    <Translation
                      translationKey={Localization.Admin.FeaturedStreams.streamsYouFeatureTab}
                      defaultText="Streams you feature"
                    />
                    {featuredCount > 0 && (
                      <Badge count={featuredCount} overflowCount={999} style={countBadgeStyle} />
                    )}
                  </Space>
                ),
                children: (
                  <Space direction="vertical" size="large" style={{ width: '100%' }}>
                    <Button
                      type="primary"
                      icon={<PlusOutlined />}
                      onClick={() => setModalOpen(true)}
                      size="large"
                    >
                      <Translation
                        translationKey={Localization.Admin.FeaturedStreams.featureStreamButton}
                        defaultText="Feature Live Stream"
                      />
                    </Button>

                    <FederatedServersTable
                      servers={federatedServers}
                      loading={serversLoading}
                      onRemove={removeServer}
                    />
                  </Space>
                ),
              },
              {
                key: 'featuring-you',
                label: (
                  <Space size={8}>
                    {/* Red dot on the title draws the eye when requests are
                        waiting for approval; it clears once none are pending. */}
                    <Badge dot={pendingRequestCount > 0} offset={[2, 0]}>
                      <Translation
                        translationKey={Localization.Admin.FeaturedStreams.featuringYouTab}
                        defaultText="Featuring you"
                      />
                    </Badge>
                    {featuringYouCount > 0 && (
                      <Badge
                        count={featuringYouCount}
                        overflowCount={999}
                        style={countBadgeStyle}
                      />
                    )}
                  </Space>
                ),
                children: (
                  <Space direction="vertical" size="large" style={{ width: '100%' }}>
                    <FeatureRequests
                      requests={featureRequests}
                      loading={requestsLoading}
                      onApprove={approveFeatureRequest}
                      onReject={rejectFeatureRequest}
                    />

                    <DirectoryListings
                      directories={directoryFollowers}
                      loading={directoriesLoading}
                      onRemove={removeDirectoryFollower}
                    />
                  </Space>
                ),
              },
            ]}
          />

          <FeatureStreamModal
            open={modalOpen}
            onCancel={() => setModalOpen(false)}
            onOk={handleFeatureStream}
          />
        </>
      ) : (
        <Alert
          message={
            <Translation
              translationKey={Localization.Admin.FeaturedStreams.socialFeaturesRequired}
              defaultText="Social features must be enabled"
            />
          }
          description={
            <>
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.socialFeaturesRequiredDesc}
                defaultText="You must enable social features in the"
              />{' '}
              <Link href="/admin/config-federation">
                <Translation
                  translationKey={Localization.Admin.FeaturedStreams.federationSettings}
                  defaultText="Federation settings"
                />
              </Link>{' '}
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.socialFeaturesRequiredDesc}
                defaultText="before you can feature other streams."
              />
            </>
          }
          type="warning"
          showIcon
        />
      )}
    </div>
  );
};

ConfigFeatured.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};

export default ConfigFeatured;
