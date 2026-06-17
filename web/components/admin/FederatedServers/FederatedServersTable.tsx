import React, { FC, useState } from 'react';
import { Table, Button, Space, Tag, Popconfirm, message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import {
  DeleteOutlined,
  LinkOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
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
  lastStatusUpdate?: string;
  addedAt: string;
}

export interface FederatedServersTableProps {
  servers: FederatedServerData[];
  loading?: boolean;
  onRemove: (id: number) => Promise<void>;
}

// Prefer the human-friendly display name, fall back to the federation
// username so the column is never blank for an accepted server.
const serverLabel = (server: FederatedServerData): string =>
  server.displayName || server.name || '';

export const FederatedServersTable: FC<FederatedServersTableProps> = ({
  servers,
  loading = false,
  onRemove,
}) => {
  const { t } = useTranslation();
  const [removingId, setRemovingId] = useState<number | null>(null);

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
      render: (isOnline: boolean) => (
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
      ),
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
        text || (
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
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.FeaturedStreams.actions}
          defaultText="Actions"
        />
      ),
      key: 'actions',
      render: (_: any, record: FederatedServerData) => (
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

  return (
    <Table
      className={styles.table}
      columns={columns}
      dataSource={servers}
      rowKey="id"
      loading={loading}
      pagination={{
        pageSize: 10,
        showSizeChanger: true,
        showTotal: (total: number) => `Total ${total} streams`, // Note: This would need custom translation handling
      }}
    />
  );
};
