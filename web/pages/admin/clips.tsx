import { Alert, Button, Popconfirm, Select, Space, Table, Tabs, Typography } from 'antd';
import type { ColumnsType } from 'antd/lib/table';
import { ReactElement, useCallback, useContext, useEffect, useState } from 'react';
import { format } from 'date-fns';
import { AdminLayout } from '../../components/layouts/AdminLayout';
import { ToggleSwitch } from '../../components/admin/ToggleSwitch';
import { TextFieldWithSubmit } from '../../components/admin/TextFieldWithSubmit';
import { FormStatusIndicator } from '../../components/admin/FormStatusIndicator';
import { ServerStatusContext } from '../../utils/server-status-context';
import {
  CLIP_PERMISSION_OPTIONS,
  FIELD_PROPS_CLIPS_ENABLED,
  FIELD_PROPS_CLIP_PERMISSIONS,
  FIELD_PROPS_MAX_CLIP_DURATION,
  FIELD_PROPS_REPLAY_ENABLED,
  RESET_TIMEOUT,
  postConfigUpdateToAPI,
} from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { ADMIN_CLIPS, CLIP_DELETE, REPLAYS, REPLAY_DELETE, fetchData } from '../../utils/apis';
import type { Clip, Replay } from '../../interfaces/clip.model';

const { Title, Paragraph } = Typography;

// formatBytes renders a byte count in the largest sensible unit.
function formatBytes(bytes: number): string {
  if (!bytes) {
    return '0 MB';
  }
  const megabytes = bytes / 1024 / 1024;
  if (megabytes < 1024) {
    return `${megabytes.toFixed(1)} MB`;
  }
  return `${(megabytes / 1024).toFixed(2)} GB`;
}

// formatDuration renders a number of seconds as h:mm:ss or m:ss.
function formatDuration(seconds: number): string {
  const total = Math.max(Math.round(seconds || 0), 0);
  const hours = Math.floor(total / 3600);
  const minutes = Math.floor((total % 3600) / 60);
  const remainder = total % 60;

  if (hours > 0) {
    return `${hours}:${String(minutes).padStart(2, '0')}:${String(remainder).padStart(2, '0')}`;
  }
  return `${minutes}:${String(remainder).padStart(2, '0')}`;
}

function formatTimestamp(value?: string): string {
  if (!value) {
    return '';
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return '';
  }
  return format(date, 'PP p');
}

export default function Clips(): ReactElement {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { replay } = serverConfig || {};

  const replayEnabled = !!replay?.enabled;

  const [replays, setReplays] = useState<Replay[]>([]);
  const [clips, setClips] = useState<Clip[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [maxClipDuration, setMaxClipDuration] = useState('');
  const [clipPermissionStatus, setClipPermissionStatus] = useState<StatusState>(null);

  const saveClipPermissions = async (value: string) => {
    setClipPermissionStatus(createInputStatus(STATUS_PROCESSING));
    await postConfigUpdateToAPI({
      apiPath: FIELD_PROPS_CLIP_PERMISSIONS.apiPath,
      data: { value },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'clipPermissions',
          value,
          path: FIELD_PROPS_CLIP_PERMISSIONS.configPath,
        });
        setClipPermissionStatus(createInputStatus(STATUS_SUCCESS));
      },
      onError: (message: string) => {
        setClipPermissionStatus(createInputStatus(STATUS_ERROR, message));
      },
    });
    setTimeout(() => setClipPermissionStatus(null), RESET_TIMEOUT);
  };

  // Seed the editable field from the server's value once it arrives.
  useEffect(() => {
    if (replay?.maxClipDurationSeconds) {
      setMaxClipDuration(String(replay.maxClipDurationSeconds));
    }
  }, [replay?.maxClipDurationSeconds]);

  const loadData = useCallback(async () => {
    if (!replayEnabled) {
      setReplays([]);
      setClips([]);
      return;
    }

    setLoading(true);
    setError('');
    try {
      const [replayResult, clipResult] = await Promise.all([
        fetchData(REPLAYS),
        fetchData(ADMIN_CLIPS),
      ]);
      setReplays(replayResult || []);
      setClips(clipResult || []);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Unable to load replays and clips.');
    } finally {
      setLoading(false);
    }
  }, [replayEnabled]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const deleteItem = async (url: string, id: string) => {
    setError('');
    try {
      const result = await fetchData(url, {
        method: 'POST',
        data: { id },
        auth: true,
      });

      if (result?.success === false) {
        setError(result.message || 'Unable to delete.');
        return;
      }

      await loadData();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Unable to delete.');
    }
  };

  const replayColumns: ColumnsType<Replay> = [
    {
      title: 'Stream',
      key: 'title',
      render: (_, record) => (
        <Space orientation="vertical" size={0}>
          <span>{record.title || 'Untitled stream'}</span>
          <Typography.Text type="secondary">{formatTimestamp(record.startTime)}</Typography.Text>
        </Space>
      ),
    },
    {
      title: 'Duration',
      key: 'duration',
      render: (_, record) =>
        record.inProgress ? 'In progress' : formatDuration(record.durationSeconds),
    },
    {
      title: 'Size',
      key: 'size',
      render: (_, record) => formatBytes(record.totalBytes),
      sorter: (a, b) => a.totalBytes - b.totalBytes,
    },
    {
      title: 'Clips',
      dataIndex: 'clipCount',
      key: 'clipCount',
      sorter: (a, b) => a.clipCount - b.clipCount,
    },
    {
      title: '',
      key: 'actions',
      render: (_, record) => (
        <Popconfirm
          title="Delete this replay?"
          description={
            record.clipCount > 0
              ? `This permanently deletes the recorded video and the ${record.clipCount} clip(s) taken from it.`
              : 'This permanently deletes the recorded video.'
          }
          okText="Delete"
          okButtonProps={{ danger: true }}
          onConfirm={() => deleteItem(REPLAY_DELETE, record.id)}
        >
          <Button danger size="small" disabled={record.inProgress}>
            Delete
          </Button>
        </Popconfirm>
      ),
    },
  ];

  const clipColumns: ColumnsType<Clip> = [
    {
      title: 'Clip',
      key: 'title',
      render: (_, record) => (
        <Space size="middle">
          {record.thumbnail && (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={record.thumbnail}
              alt={record.title || 'Clip'}
              style={{ width: 96, aspectRatio: '16 / 9', objectFit: 'cover', borderRadius: 4 }}
            />
          )}
          <Space orientation="vertical" size={0}>
            <a href={`/clips/${record.id}`} target="_blank" rel="noreferrer">
              {record.title || 'Untitled clip'}
            </a>
            <Typography.Text type="secondary">
              {record.streamTitle || record.streamId}
            </Typography.Text>
          </Space>
        </Space>
      ),
    },
    {
      title: 'Duration',
      key: 'duration',
      render: (_, record) => formatDuration(record.durationSeconds),
      sorter: (a, b) => a.durationSeconds - b.durationSeconds,
    },
    {
      title: 'Clipped by',
      key: 'clippedBy',
      render: (_, record) => record.clippedBy || 'Anonymous',
    },
    {
      title: 'Created',
      key: 'timestamp',
      render: (_, record) => formatTimestamp(record.timestamp),
      sorter: (a, b) => (a.timestamp > b.timestamp ? 1 : -1),
    },
    {
      title: '',
      key: 'actions',
      render: (_, record) => (
        <Popconfirm
          title="Delete this clip?"
          description="The recorded video it came from is not affected."
          okText="Delete"
          okButtonProps={{ danger: true }}
          onConfirm={() => deleteItem(CLIP_DELETE, record.id)}
        >
          <Button danger size="small">
            Delete
          </Button>
        </Popconfirm>
      ),
    },
  ];

  const settings = (
    <>
      <Title level={3}>Replays and Clips</Title>
      <Paragraph>
        Replays are recordings of your past streams. Clips are short, shareable moments taken from
        them.
      </Paragraph>

      <Alert
        type="warning"
        showIcon
        message="Saving replays keeps every video segment on disk"
        description="Video from saved replays is kept on disk instead of being cleaned up as a stream progresses, so disk use grows with every stream you record. Delete replays you no longer need."
        style={{ marginBottom: '1rem' }}
      />

      {replay?.forcedByCommandLine && (
        <Alert
          type="info"
          showIcon
          message="Replays are enabled by a command line flag"
          description="This server was started with -enableReplayFeatures, so replays stay on regardless of this setting."
          style={{ marginBottom: '1rem' }}
        />
      )}

      <ToggleSwitch
        fieldName="enabled"
        useSubmit
        {...FIELD_PROPS_REPLAY_ENABLED}
        checked={replayEnabled}
        disabled={replay?.forcedByCommandLine}
      />
      <ToggleSwitch
        fieldName="clipsEnabled"
        useSubmit
        {...FIELD_PROPS_CLIPS_ENABLED}
        checked={!!replay?.clipsEnabled}
        disabled={!replayEnabled}
      />
      <div className="formfield-container">
        <div className="label-side">
          <span className="formfield-label">{FIELD_PROPS_CLIP_PERMISSIONS.label}</span>
        </div>
        <div className="input-side">
          <div className="input-group">
            <Select
              style={{ width: 260 }}
              value={replay?.clipPermissions || 'established'}
              onChange={saveClipPermissions}
              disabled={!replayEnabled || !replay?.clipsEnabled}
              options={CLIP_PERMISSION_OPTIONS.map(({ value, label }) => ({ value, label }))}
            />
            <FormStatusIndicator status={clipPermissionStatus} />
          </div>
          <p className="field-tip">
            {CLIP_PERMISSION_OPTIONS.find(
              o => o.value === (replay?.clipPermissions || 'established'),
            )?.description || FIELD_PROPS_CLIP_PERMISSIONS.tip}{' '}
            Moderators can always create clips. The button is hidden from viewers who are not
            allowed.
          </p>
        </div>
      </div>
      <TextFieldWithSubmit
        fieldName="maxClipDurationSeconds"
        {...FIELD_PROPS_MAX_CLIP_DURATION}
        value={maxClipDuration}
        initialValue={String(replay?.maxClipDurationSeconds ?? '')}
        type="number"
        useTrim
        onChange={({ value }) => setMaxClipDuration(value)}
      />
    </>
  );

  return (
    <div>
      {settings}

      {error && <Alert type="error" message={error} showIcon style={{ margin: '1rem 0' }} />}

      {!replayEnabled ? (
        <Alert
          type="info"
          showIcon
          message="Replays are not being saved"
          description="Turn on the setting above to start keeping recorded video, which is what clips are made from."
          style={{ marginTop: '1rem' }}
        />
      ) : (
        <Tabs
          defaultActiveKey="clips"
          style={{ marginTop: '1rem' }}
          items={[
            {
              key: 'clips',
              label: `Clips (${clips.length})`,
              children: (
                <Table
                  rowKey="id"
                  loading={loading}
                  dataSource={clips}
                  columns={clipColumns}
                  pagination={{ pageSize: 20, hideOnSinglePage: true }}
                />
              ),
            },
            {
              key: 'replays',
              label: `Replays (${replays.length})`,
              children: (
                <Table
                  rowKey="id"
                  loading={loading}
                  dataSource={replays}
                  columns={replayColumns}
                  pagination={{ pageSize: 20, hideOnSinglePage: true }}
                />
              ),
            },
          ]}
        />
      )}
    </div>
  );
}

Clips.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
