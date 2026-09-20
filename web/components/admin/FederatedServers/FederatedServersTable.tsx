import { format, parseISO, isValid } from 'date-fns';
import { FC, useState } from 'react';
import { Table, Button, Space, Tag, Popconfirm, message, Tooltip } from 'antd';
import { useTranslation } from 'next-export-i18n';
import {
  DeleteOutlined,
  LinkOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
} from '@ant-design/icons';
import { ColumnsType } from 'antd/es/table';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';
import styles from './FederatedServersTable.module.scss';

export interface FederatedServerData {
  id: number;
  iri: string;
  name?: string;
  displayName?: string;
  isOnline: boolean;
  streamTitle?: string;
  lastStatusUpdate?: string;
  addedAt: string;
  followStatus?: string;
  priority?: number;
}

export interface FederatedServersTableProps {
  servers: FederatedServerData[];
  loading?: boolean;
  onRemove: (id: number) => Promise<void>;
  onReorder: (id: number, direction: 'up' | 'down') => Promise<void>;
}

// Prefer the human-friendly display name, fall back to the federation
// username so the column is never blank for an accepted server.
const serverLabel = (server: FederatedServerData): string =>
  server.displayName || server.name || '';

export const FederatedServersTable: FC<FederatedServersTableProps> = ({
  servers,
  loading = false,
  onRemove,
  onReorder,
}) => {
  const { t } = useTranslation();
  const [removingId, setRemovingId] = useState<number | null>(null);
  const [reorderingId, setReorderingId] = useState<number | null>(null);

  const handleRemove = async (id: number) => {
    setRemovingId(id);
    try {
      await onRemove(id);
      message.success(t(Localization.Admin.FeaturedStreams.streamUnfeaturedSuccess));
    } catch {
      message.error(t(Localization.Admin.FeaturedStreams.failedToUnfeature));
    } finally {
      setRemovingId(null);
    }
  };

  const acceptedCount = servers.filter(s => s.followStatus === 'accepted').length;

  // Sort: accepted servers first (in priority/API order), then pending, then rejected.
  // Defined before columns so the Priority column render can reference it.
  const sortedServers = [
    ...servers.filter(s => s.followStatus === 'accepted'),
    ...servers.filter(s => s.followStatus !== 'accepted'),
  ];
  const columns: ColumnsType<FederatedServerData> = [
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.streamName}
          defaultText="Stream Name"
        />
      ),
      dataIndex: 'name',
      key: 'name',
      render: (_: string, record: FederatedServerData) => (
        <Space>
          <span>{serverLabel(record)}</span>
          <a
            href={record.iri}
            target="_blank"
            rel="noopener noreferrer"
            onClick={e => e.stopPropagation()}
          >
            <LinkOutlined />
          </a>
        </Space>
      ),
    },
    {
      title: (
        <Translation translationKey={Localization.Admin.FeaturedStreams.url} defaultText="URL" />
      ),
      dataIndex: 'iri',
      key: 'iri',
      ellipsis: true,
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.status}
          defaultText="Status"
        />
      ),
      dataIndex: 'isOnline',
      key: 'isOnline',
      render: (isOnline: boolean, record: FederatedServerData) => {
        // A server we've requested to feature but which hasn't accepted our
        // follow yet isn't live anywhere yet -- surface that it's awaiting the
        // remote server's approval instead of showing a misleading Offline.
        if (record.followStatus && record.followStatus !== 'accepted') {
          return (
            <Tag icon={<ClockCircleOutlined />} color="warning">
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.pendingApproval}
                defaultText="Pending approval"
              />
            </Tag>
          );
        }
        return (
          <Tag
            icon={isOnline ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
            color={isOnline ? 'success' : 'default'}
          >
            {isOnline ? (
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.online}
                defaultText="Online"
              />
            ) : (
              <Translation
                translationKey={Localization.Admin.FeaturedStreams.offline}
                defaultText="Offline"
              />
            )}
          </Tag>
        );
      },
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.streamTitle}
          defaultText="Stream Title"
        />
      ),
      dataIndex: 'streamTitle',
      key: 'streamTitle',
      ellipsis: true,
      // Only a live server has a current stream title; hide it for offline or
      // pending servers so a stale title isn't shown.
      render: (streamTitle: string, record: FederatedServerData) =>
        record.isOnline ? streamTitle || '' : '',
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.lastChecked}
          defaultText="Last Checked"
        />
      ),
      dataIndex: 'lastStatusUpdate',
      key: 'lastStatusUpdate',
      render: (text: string) =>
        text && isValid(parseISO(text)) ? (
          format(parseISO(text), 'MMM d, yyyy HH:mm')
        ) : (
          <Translation
            translationKey={Localization.Admin.FeaturedStreams.never}
            defaultText="Never"
          />
        ),
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.added}
          defaultText="Added"
        />
      ),
      dataIndex: 'addedAt',
      key: 'addedAt',
      render: (text: string) =>
        text && isValid(parseISO(text)) ? format(parseISO(text), 'MMM d, yyyy') : (text ?? ''),
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.priority}
          defaultText="Priority"
        />
      ),
      key: 'priority',
      render: (_: unknown, record: FederatedServerData) => {
        if (record.followStatus !== 'accepted') return null;
        const acceptedServers = sortedServers.filter(s => s.followStatus === 'accepted');
        const acceptedIdx = acceptedServers.findIndex(s => s.id === record.id);
        const isFirst = acceptedIdx === 0;
        const isLast = acceptedIdx === acceptedServers.length - 1;
        return (
          <Space size="small">
            <Tooltip
              title={
                <Translation
                  translationKey={Localization.Admin.FeaturedStreams.moveUp}
                  defaultText="Move up"
                />
              }
            >
              <Button
                size="small"
                icon={<ArrowUpOutlined />}
                disabled={isFirst || reorderingId === record.id}
                loading={reorderingId === record.id}
                onClick={async () => {
                  setReorderingId(record.id);
                  try {
                    await onReorder(record.id, 'up');
                  } catch {
                    message.error(t(Localization.Admin.FeaturedStreams.failedToReorder));
                  } finally {
                    setReorderingId(null);
                  }
                }}
                aria-label={t(Localization.Admin.FeaturedStreams.moveUp)}
              />
            </Tooltip>
            <Tooltip
              title={
                <Translation
                  translationKey={Localization.Admin.FeaturedStreams.moveDown}
                  defaultText="Move down"
                />
              }
            >
              <Button
                size="small"
                icon={<ArrowDownOutlined />}
                disabled={isLast || reorderingId === record.id}
                loading={reorderingId === record.id}
                onClick={async () => {
                  setReorderingId(record.id);
                  try {
                    await onReorder(record.id, 'down');
                  } catch {
                    message.error(t(Localization.Admin.FeaturedStreams.failedToReorder));
                  } finally {
                    setReorderingId(null);
                  }
                }}
                aria-label={t(Localization.Admin.FeaturedStreams.moveDown)}
              />
            </Tooltip>
          </Space>
        );
      },
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.actions}
          defaultText="Actions"
        />
      ),
      key: 'actions',
      render: (_: unknown, record: FederatedServerData) => (
        <Popconfirm
          title={
            <Translation
              translationKey={Localization.Admin.FeaturedStreams.unfeatureConfirm}
              defaultText="Unfeature {{name}}?"
              vars={{ name: serverLabel(record) }}
            />
          }
          onConfirm={() => handleRemove(record.id)}
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
          <Button danger size="small" icon={<DeleteOutlined />} loading={removingId === record.id}>
            <Translation
              translationKey={Localization.Admin.FeaturedStreams.unfeatureButton}
              defaultText="Unfeature"
            />
          </Button>
        </Popconfirm>
      ),
    },
  ];

  // Hide the Priority column entirely when there is only one server — move
  // buttons would both be disabled and the column is just visual noise.
  const visibleColumns = acceptedCount > 1 ? columns : columns.filter(c => c.key !== 'priority');

  return (
    <Table
      className={styles.table}
      columns={visibleColumns}
      dataSource={sortedServers}
      rowKey="id"
      loading={loading}
      pagination={{
        pageSize: 10,
        showSizeChanger: true,
        showTotal: (total: number) => `Total ${total} streams`,
      }}
    />
  );
};
