import { FC, useState } from 'react';
import { Spin, Button, Avatar, Typography, Popconfirm, message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { DeleteOutlined, SendOutlined } from '@ant-design/icons';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';
import { DirectoryFollower } from '../../../hooks/useDirectoryFollowers';
import { isValidUrl } from '../../../utils/validators';
import styles from './FederatedServerList.module.scss';

const { Text, Title, Paragraph } = Typography;

export interface DirectoryListingsProps {
  directories: DirectoryFollower[];
  loading?: boolean;
  onRemove: (actorIRI: string) => Promise<void>;
  onResendApproval: (actorIRI: string) => Promise<void>;
}

// DirectoryListings shows the directories that are currently listing this
// server (accepted directory followers), with the option to remove this server
// from any of them. It renders nothing when there are none.
export const DirectoryListings: FC<DirectoryListingsProps> = ({
  directories,
  loading = false,
  onRemove,
  onResendApproval,
}) => {
  const { t } = useTranslation();
  const [pendingAction, setPendingAction] = useState<string | null>(null);

  const handleRemove = async (actorIRI: string) => {
    setPendingAction(`remove:${actorIRI}`);
    try {
      await onRemove(actorIRI);
    } catch {
      message.error(t(Localization.Admin.FeaturedStreams.failedToRemoveDirectory));
    } finally {
      setPendingAction(null);
    }
  };

  const handleResendApproval = async (actorIRI: string) => {
    setPendingAction(`resend:${actorIRI}`);
    try {
      await onResendApproval(actorIRI);
      message.success(t(Localization.Admin.FeaturedStreams.approvalResent));
    } catch {
      message.error(t(Localization.Admin.FeaturedStreams.failedToResendApproval));
    } finally {
      setPendingAction(null);
    }
  };

  return (
    <div>
      <Title level={3}>
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.directoryListingsTitle}
          defaultText="Directories listing you"
        />
      </Title>
      <Paragraph type="secondary">
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.directoryListingsDescription}
          defaultText="These directories follow your server so they can list your stream. Each one receives your live status. Remove one to stop sending it your status and drop your server from its listing."
        />
      </Paragraph>
      {!loading && directories.length === 0 ? (
        <Paragraph type="secondary">
          <Translation
            translationKey={Localization.Admin.FeaturedStreams.directoryListingsEmpty}
            defaultText="No directories are featuring your stream yet."
          />
        </Paragraph>
      ) : (
        <Spin spinning={loading}>
          <ul className={styles.list}>
            {directories.map(directory => (
              <li key={directory.link} className={styles.item}>
                <div className={styles.meta}>
                  {directory.image ? (
                    <Avatar src={directory.image} />
                  ) : (
                    <Avatar>{(directory.name || directory.username || '?').charAt(0)}</Avatar>
                  )}
                  <div className={styles.text}>
                    {/* Only render the remote-supplied link as a clickable anchor
                        when it is a valid http(s) URL; otherwise show plain text so a
                        hostile value (e.g. a javascript: URL) can't reach href. */}
                    {isValidUrl(directory.link) ? (
                      <a href={directory.link} target="_blank" rel="noopener noreferrer">
                        {directory.name || directory.username || directory.link}
                      </a>
                    ) : (
                      <span>{directory.name || directory.username || directory.link}</span>
                    )}
                    <Text type="secondary">{directory.username || directory.link}</Text>
                  </div>
                </div>
                <div className={styles.actions}>
                  <Button
                    size="small"
                    icon={<SendOutlined />}
                    loading={pendingAction === `resend:${directory.link}`}
                    disabled={pendingAction === `remove:${directory.link}`}
                    onClick={() => handleResendApproval(directory.link)}
                  >
                    <Translation
                      translationKey={Localization.Admin.FeaturedStreams.resendApprovalButton}
                      defaultText="Resend approval"
                    />
                  </Button>
                  <Popconfirm
                    title={
                      <Translation
                        translationKey={
                          Localization.Admin.FeaturedStreams.removeFromDirectoryConfirm
                        }
                        defaultText="Remove your server from this directory?"
                      />
                    }
                    onConfirm={() => handleRemove(directory.link)}
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
                      danger
                      size="small"
                      icon={<DeleteOutlined />}
                      loading={pendingAction === `remove:${directory.link}`}
                      disabled={pendingAction === `resend:${directory.link}`}
                    >
                      <Translation
                        translationKey={
                          Localization.Admin.FeaturedStreams.removeFromDirectoryButton
                        }
                        defaultText="Remove"
                      />
                    </Button>
                  </Popconfirm>
                </div>
              </li>
            ))}
          </ul>
        </Spin>
      )}
    </div>
  );
};
