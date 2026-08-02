import { Space, Table, Tag, Tooltip, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import { SortOrder } from 'antd/lib/table/interface';
import { formatDistanceToNow } from 'date-fns';
import { FC } from 'react';
import { useTranslation } from 'next-export-i18n';
import { PlaybackClient } from '../../types/playback';
import { Localization } from '../../types/localization';
import { Translation } from '../ui/Translation/Translation';
import { formatUAstring } from '../../utils/format';

export type PlaybackClientTableProps = {
  data: PlaybackClient[];
};

// Shown instead of a number when a value was never measured for a client,
// which is different from it having been measured as zero.
const UNKNOWN = '—';

const formatNumber = (value: number | null | undefined, digits: number): string =>
  value === null || value === undefined ? UNKNOWN : value.toFixed(digits);

// Keeps unknown values at the bottom in both directions. Ant Design
// negates the comparator for a descending sort, so the sentinel has to
// flip with it or a column of mostly-unmeasured clients fills its first
// screen with blanks.
const compareMeasurements = (
  a: number | null | undefined,
  b: number | null | undefined,
  sortOrder?: SortOrder,
): number => {
  const last = sortOrder === 'descend' ? -1 : 1;
  if (a === null || a === undefined) return b === null || b === undefined ? 0 : last;
  if (b === null || b === undefined) return -last;
  return a - b;
};
const playerStateLabel = (state: string | undefined) => {
  const translationKey = {
    p: Localization.Admin.PlaybackClients.playing,
    a: Localization.Admin.PlaybackClients.paused,
    w: Localization.Admin.PlaybackClients.buffering,
    k: Localization.Admin.PlaybackClients.seeking,
    e: Localization.Admin.PlaybackClients.ended,
  }[state || ''];

  return (
    <Translation
      translationKey={translationKey || Localization.Admin.PlaybackClients.unknownState}
      defaultText={
        {
          p: 'Playing',
          a: 'Paused',
          w: 'Buffering',
          k: 'Seeking',
          e: 'Ended',
        }[state || ''] || 'Unknown'
      }
    />
  );
};

export const PlaybackClientTable: FC<PlaybackClientTableProps> = ({ data }) => {
  const { t } = useTranslation();

  const columns: ColumnsType<PlaybackClient> = [
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.player}
          defaultText="Player"
        />
      ),
      dataIndex: 'userAgent',
      key: 'userAgent',
      render: (userAgent: string, client: PlaybackClient) => {
        const { playback } = client;
        let tag = (
          <Tag color="default">
            <Translation
              translationKey={Localization.Admin.PlaybackClients.sourceUnknown}
              defaultText="No data"
            />
          </Tag>
        );

        if (playback?.measurementStatus === 'unmeasurable') {
          tag = (
            <Tooltip title={t(Localization.Admin.PlaybackClients.serverUnmeasurableDescription)}>
              <Tag color="warning">
                <Translation
                  translationKey={Localization.Admin.PlaybackClients.serverUnmeasurable}
                  defaultText="Server can't measure"
                />
              </Tag>
            </Tooltip>
          );
        } else if (playback) {
          tag = (
            <Tooltip
              title={t(Localization.Admin.PlaybackClients.lastMeasured, {
                age: formatDistanceToNow(new Date(playback.lastUpdate)),
              })}
            >
              <Tag color={playback.source === 'client' ? 'blue' : 'default'}>
                {playback.source === 'client' ? (
                  <Translation
                    translationKey={Localization.Admin.PlaybackClients.sourceClient}
                    defaultText="Player reporting"
                  />
                ) : (
                  <Translation
                    translationKey={Localization.Admin.PlaybackClients.sourceServer}
                    defaultText="Server measuring"
                  />
                )}
              </Tag>
            </Tooltip>
          );
        }

        return (
          <Space direction="vertical" size={0}>
            <span>{formatUAstring(userAgent)}</span>
            {tag}
          </Space>
        );
      },
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.location}
          defaultText="Location"
        />
      ),
      dataIndex: 'geo',
      key: 'geo',
      render: geo => (geo ? `${geo.regionName}, ${geo.countryCode}` : UNKNOWN),
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.watchTime}
          defaultText="Watch time"
        />
      ),
      dataIndex: 'firstSeen',
      key: 'firstSeen',
      defaultSortOrder: 'ascend' as SortOrder,
      render: (time: string) => formatDistanceToNow(new Date(time)),
      sorter: (a: PlaybackClient, b: PlaybackClient) =>
        new Date(a.firstSeen).getTime() - new Date(b.firstSeen).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.state}
          defaultText="State"
        />
      ),
      key: 'playerState',
      render: (_, client: PlaybackClient) => playerStateLabel(client.playback?.playerState),
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.speed}
          defaultText="Speed (kbps)"
        />
      ),
      key: 'bandwidth',
      className: 'number-col',
      render: (_, client: PlaybackClient) => formatNumber(client.playback?.bandwidthKbps, 0),
      sorter: (a: PlaybackClient, b: PlaybackClient, sortOrder?: SortOrder) =>
        compareMeasurements(a.playback?.bandwidthKbps, b.playback?.bandwidthKbps, sortOrder),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.quality}
          defaultText="Quality (kbps)"
        />
      ),
      key: 'bitrate',
      className: 'number-col',
      render: (_, client: PlaybackClient) => formatNumber(client.playback?.bitrateKbps, 0),
      sorter: (a: PlaybackClient, b: PlaybackClient, sortOrder?: SortOrder) =>
        compareMeasurements(a.playback?.bitrateKbps, b.playback?.bitrateKbps, sortOrder),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.latency}
          defaultText="Latency (s)"
        />
      ),
      key: 'latency',
      className: 'number-col',
      render: (_, client: PlaybackClient) => formatNumber(client.playback?.latencySeconds, 1),
      sorter: (a: PlaybackClient, b: PlaybackClient, sortOrder?: SortOrder) =>
        compareMeasurements(a.playback?.latencySeconds, b.playback?.latencySeconds, sortOrder),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.segmentDownload}
          defaultText="Segment download (s)"
        />
      ),
      key: 'downloadDuration',
      className: 'number-col',
      render: (_, client: PlaybackClient) => formatNumber(client.playback?.downloadSeconds, 2),
      sorter: (a: PlaybackClient, b: PlaybackClient, sortOrder?: SortOrder) =>
        compareMeasurements(a.playback?.downloadSeconds, b.playback?.downloadSeconds, sortOrder),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.errors}
          defaultText="Errors"
        />
      ),
      key: 'errors',
      className: 'number-col',
      render: (_, client: PlaybackClient) => {
        const errors = client.playback?.errorCount;
        if (errors === null || errors === undefined) {
          return UNKNOWN;
        }
        return <Typography.Text type={errors > 0 ? 'danger' : undefined}>{errors}</Typography.Text>;
      },
      sorter: (a: PlaybackClient, b: PlaybackClient, sortOrder?: SortOrder) =>
        compareMeasurements(a.playback?.errorCount, b.playback?.errorCount, sortOrder),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: (
        <Translation
          translationKey={Localization.Admin.PlaybackClients.qualityChanges}
          defaultText="Quality changes"
        />
      ),
      key: 'qualityVariantChanges',
      className: 'number-col',
      render: (_, client: PlaybackClient) =>
        formatNumber(client.playback?.qualityVariantChanges, 0),
      sorter: (a: PlaybackClient, b: PlaybackClient, sortOrder?: SortOrder) =>
        compareMeasurements(
          a.playback?.qualityVariantChanges,
          b.playback?.qualityVariantChanges,
          sortOrder,
        ),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
  ];

  return (
    <Table
      pagination={{ hideOnSinglePage: true }}
      className="table-container"
      columns={columns}
      dataSource={data}
      size="small"
      rowKey="clientID"
    />
  );
};
