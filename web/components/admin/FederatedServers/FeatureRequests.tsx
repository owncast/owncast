import React, { FC, useState } from 'react';
import { List, Button, Avatar, Typography, Popconfirm, message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { CheckOutlined, CloseOutlined } from '@ant-design/icons';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';
import { FeatureRequest } from '../../../hooks/useFeatureRequests';
import { isValidUrl } from '../../../utils/validators';

const { Text, Title, Paragraph } = Typography;

export interface FeatureRequestsProps {
  requests: FeatureRequest[];
  loading?: boolean;
  onApprove: (actorIRI: string) => Promise<void>;
  onReject: (actorIRI: string) => Promise<void>;
}

// FeatureRequests lists incoming requests from other Owncast servers asking to
// feature this server's stream in their directory, with approve/reject
// actions. It renders nothing when there are no pending requests.
export const FeatureRequests: FC<FeatureRequestsProps> = ({
  requests,
  loading = false,
  onApprove,
  onReject,
}) => {
  const { t } = useTranslation();
  const [pendingIRI, setPendingIRI] = useState<string | null>(null);

  if (!loading && requests.length === 0) {
    return null;
  }

  const handle = async (
    actorIRI: string,
    action: (iri: string) => Promise<void>,
    failKey: string,
  ) => {
    setPendingIRI(actorIRI);
    try {
      await action(actorIRI);
    } catch {
      message.error(t(failKey));
    } finally {
      setPendingIRI(null);
    }
  };

  return (
    <div>
      <Title level={3}>
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.featureRequestsTitle}
          defaultText="Requests to feature your stream"
        />
      </Title>
      <Paragraph type="secondary">
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.featureRequestsDescription}
          defaultText="These Owncast servers have asked to show your stream in their featured streams directory. Approve a server to let it display your live status."
        />
      </Paragraph>
      <List
        loading={loading}
        dataSource={requests}
        rowKey="link"
        renderItem={request => (
          <List.Item
            actions={[
              <Popconfirm
                key="approve"
                title={
                  <Translation
                    translationKey={Localization.Admin.FeaturedStreams.confirmYes}
                    defaultText="Yes"
                  />
                }
                onConfirm={() =>
                  handle(
                    request.link,
                    onApprove,
                    Localization.Admin.FeaturedStreams.failedToApprove,
                  )
                }
                okText={
                  <Translation
                    translationKey={Localization.Admin.FeaturedStreams.confirmYes}
                    defaultText="Yes"
                  />
                }
                cancelText={
                  <Translation
                    translationKey={Localization.Admin.FeaturedStreams.confirmNo}
                    defaultText="No"
                  />
                }
              >
                <Button
                  type="primary"
                  size="small"
                  icon={<CheckOutlined />}
                  loading={pendingIRI === request.link}
                >
                  <Translation
                    translationKey={Localization.Admin.FeaturedStreams.approveButton}
                    defaultText="Approve"
                  />
                </Button>
              </Popconfirm>,
              <Button
                key="reject"
                danger
                size="small"
                icon={<CloseOutlined />}
                loading={pendingIRI === request.link}
                onClick={() =>
                  handle(request.link, onReject, Localization.Admin.FeaturedStreams.failedToReject)
                }
              >
                <Translation
                  translationKey={Localization.Admin.FeaturedStreams.rejectButton}
                  defaultText="Reject"
                />
              </Button>,
            ]}
          >
            <List.Item.Meta
              avatar={
                request.image ? (
                  <Avatar src={request.image} />
                ) : (
                  <Avatar>{(request.name || request.username || '?').charAt(0)}</Avatar>
                )
              }
              title={
                // Only render the remote-supplied link as a clickable anchor
                // when it is a valid http(s) URL; otherwise show plain text so
                // a hostile value (e.g. a javascript: URL) can't reach href.
                isValidUrl(request.link) ? (
                  <a href={request.link} target="_blank" rel="noopener noreferrer">
                    {request.name || request.username || request.link}
                  </a>
                ) : (
                  <span>{request.name || request.username || request.link}</span>
                )
              }
              description={<Text type="secondary">{request.username || request.link}</Text>}
            />
          </List.Item>
        )}
      />
    </div>
  );
};
