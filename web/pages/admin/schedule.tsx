/* eslint-disable react/destructuring-assignment */
import {
  Button,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from 'antd';
import dynamic from 'next/dynamic';
import { useContext, useEffect, useMemo, useState } from 'react';
import type { ReactElement } from 'react';
import { useTranslation } from 'next-export-i18n';
import {
  SCHEDULE_ADMIN,
  SCHEDULE_DELETE_EVENT,
  SCHEDULE_PREVIEW,
  SCHEDULE_UPSERT_EVENT,
  fetchData,
} from '../../utils/apis';
import { API_SCHEDULE_ENABLED, API_SCHEDULE_REMINDER_MESSAGE } from '../../utils/config-constants';
import {
  WEEKDAYS,
  composeWeeklyRecurrence,
  parseWeeklyRecurrence,
  timezoneChoices,
  wallTimeInZone,
  wallTimeInZoneToUTC,
} from '../../utils/schedule';
import type { WeeklyRecurrence } from '../../utils/schedule';
import { Localization } from '../../types/localization';
import { AdminLayout } from '../../components/layouts/AdminLayout';
import { ServerStatusContext } from '../../utils/server-status-context';
import { ToggleSwitch } from '../../components/admin/ToggleSwitch';
import { TextFieldWithSubmit } from '../../components/admin/TextFieldWithSubmit';
import { TEXTFIELD_TYPE_TEXTAREA } from '../../components/admin/TextField';
import type { UpdateArgs } from '../../types/config-section';

const { Title, Paragraph } = Typography;

// Lazy loaded components

const DeleteOutlined = dynamic(() => import('@ant-design/icons/DeleteOutlined'), {
  ssr: false,
});

const EditOutlined = dynamic(() => import('@ant-design/icons/EditOutlined'), {
  ssr: false,
});

interface ScheduledEvent {
  id: string;
  seriesId?: string;
  name: string;
  description: string;
  startTime: string;
  durationMinutes: number;
  timezone: string;
  status: string;
}

interface ScheduledEventSeries {
  id: string;
  name: string;
  description: string;
  recurrence: string;
  durationMinutes: number;
  active: boolean;
}

interface AdminSchedule {
  series: ScheduledEventSeries[];
  events: ScheduledEvent[];
}

// The modal edits one of three things: a brand new entry, an existing
// concrete event, or an existing recurring series.
type EventModalTarget =
  | { kind: 'create' }
  | { kind: 'edit-event'; event: ScheduledEvent }
  | { kind: 'edit-series'; series: ScheduledEventSeries };

interface EventModalProps {
  target: EventModalTarget;
  open: boolean;
  onCancel: () => void;
  onSaved: (schedule: AdminSchedule) => void;
}

const EventModal = (props: EventModalProps) => {
  const { target, open, onCancel, onSaved } = props;
  const { t } = useTranslation();

  const browserZone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [duration, setDuration] = useState(60);
  const [repeats, setRepeats] = useState<'none' | 'weekly'>('none');
  const [timezone, setTimezone] = useState(browserZone);
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [days, setDays] = useState<string[]>([]);
  const [endsOn, setEndsOn] = useState('');
  const [preview, setPreview] = useState<string[]>([]);
  const [previewError, setPreviewError] = useState('');
  const [saving, setSaving] = useState(false);

  // A series whose recurrence this UI cannot parse (never produced by v1)
  // degrades to details-only editing.
  const seriesRule: WeeklyRecurrence | null =
    target.kind === 'edit-series' ? parseWeeklyRecurrence(target.series.recurrence) : null;
  const recurrenceEditable = target.kind !== 'edit-series' || seriesRule !== null;

  // Reset the form whenever the modal opens for a new target.
  useEffect(() => {
    if (!open) {
      return;
    }
    if (target.kind === 'create') {
      setName('');
      setDescription('');
      setDuration(60);
      setRepeats('none');
      setTimezone(browserZone);
      setDate('');
      setTime('');
      setDays([]);
      setEndsOn('');
    } else if (target.kind === 'edit-event') {
      const wall = wallTimeInZone(target.event.startTime, target.event.timezone);
      setName(target.event.name);
      setDescription(target.event.description);
      setDuration(target.event.durationMinutes);
      setRepeats('none');
      setTimezone(target.event.timezone);
      setDate(wall.date);
      setTime(wall.time);
    } else {
      setName(target.series.name);
      setDescription(target.series.description);
      setDuration(target.series.durationMinutes);
      setRepeats('weekly');
      if (seriesRule) {
        setTimezone(seriesRule.timezone);
        setDate(seriesRule.startsOn);
        setTime(seriesRule.time);
        setDays(seriesRule.days);
        setEndsOn(seriesRule.endsOn || '');
      }
    }
    setPreview([]);
    setPreviewError('');
  }, [open, target]);

  const weeklyFormComplete = days.length > 0 && date !== '' && time !== '';
  const composedRecurrence = useMemo(() => {
    if (repeats !== 'weekly' || !weeklyFormComplete) {
      return '';
    }
    const rule: WeeklyRecurrence = { days, time, startsOn: date, timezone };
    if (endsOn) {
      rule.endsOn = endsOn;
    }
    return composeWeeklyRecurrence(rule);
  }, [repeats, days, time, date, timezone, endsOn, weeklyFormComplete]);

  // Live preview: the server expands the rule with the same parser the
  // materializer uses, so what shows here is exactly what will be scheduled.
  useEffect(() => {
    if (!open || !composedRecurrence || !recurrenceEditable) {
      setPreview([]);
      setPreviewError('');
      return undefined;
    }
    let stale = false;
    const timer = setTimeout(async () => {
      try {
        const result = await fetchData(SCHEDULE_PREVIEW, {
          method: 'POST',
          data: { recurrence: composedRecurrence },
        });
        if (stale) {
          return;
        }
        setPreview(result.occurrences || []);
        setPreviewError('');
      } catch (error) {
        if (!stale) {
          setPreview([]);
          setPreviewError(`${error}`);
        }
      }
    }, 400);
    return () => {
      stale = true;
      clearTimeout(timer);
    };
  }, [open, composedRecurrence, recurrenceEditable]);

  async function save() {
    const payload: Record<string, unknown> = {
      name: name.trim(),
      description,
      durationMinutes: duration,
    };

    if (target.kind === 'edit-event') {
      payload.id = target.event.id;
      const wall = wallTimeInZone(target.event.startTime, target.event.timezone);
      if (date !== wall.date || time !== wall.time) {
        payload.start = wallTimeInZoneToUTC(date, time, target.event.timezone).toISOString();
      }
    } else if (target.kind === 'edit-series') {
      payload.id = target.series.id;
      if (recurrenceEditable) {
        payload.recurrence = composedRecurrence;
      }
    } else if (repeats === 'weekly') {
      payload.recurrence = composedRecurrence;
    } else {
      payload.start = wallTimeInZoneToUTC(date, time, timezone).toISOString();
      payload.timezone = timezone;
    }

    setSaving(true);
    try {
      const schedule = await fetchData(SCHEDULE_UPSERT_EVENT, { method: 'POST', data: payload });
      onSaved(schedule);
    } catch (error) {
      message.error(`${error}`);
    }
    setSaving(false);
  }

  const timeFormComplete =
    repeats === 'weekly' ? weeklyFormComplete && recurrenceEditable : date !== '' && time !== '';
  const okButtonProps = {
    disabled: name.trim() === '' || (recurrenceEditable ? !timeFormComplete : false),
    loading: saving,
  };

  const timezoneOptions = useMemo(
    () => timezoneChoices().map(zone => ({ label: zone, value: zone })),
    [],
  );

  const dayChips = WEEKDAYS.map(day => (
    <Tag.CheckableTag
      key={day.code}
      checked={days.includes(day.code)}
      onChange={checked =>
        setDays(checked ? [...days, day.code] : days.filter(code => code !== day.code))
      }
    >
      {day.label}
    </Tag.CheckableTag>
  ));

  const title =
    target.kind === 'create'
      ? t(Localization.Admin.Schedule.addEvent)
      : t(Localization.Admin.Schedule.editEvent);

  return (
    <Modal title={title} open={open} onOk={save} onCancel={onCancel} okButtonProps={okButtonProps}>
      <p>{t(Localization.Admin.Schedule.nameLabel)}</p>
      <Input
        value={name}
        placeholder={t(Localization.Admin.Schedule.namePlaceholder)}
        maxLength={100}
        onChange={input => setName(input.currentTarget.value)}
      />

      <p>{t(Localization.Admin.Schedule.descriptionLabel)}</p>
      <Input.TextArea
        value={description}
        rows={2}
        maxLength={500}
        onChange={input => setDescription(input.currentTarget.value)}
      />

      {target.kind === 'create' && (
        <>
          <p>{t(Localization.Admin.Schedule.repeatsLabel)}</p>
          <Select
            style={{ width: '100%' }}
            value={repeats}
            onChange={setRepeats}
            options={[
              { label: t(Localization.Admin.Schedule.doesNotRepeat), value: 'none' },
              { label: t(Localization.Admin.Schedule.weekly), value: 'weekly' },
            ]}
          />
        </>
      )}

      {recurrenceEditable && repeats === 'weekly' && (
        <>
          <p>{t(Localization.Admin.Schedule.onDays)}</p>
          <div>{dayChips}</div>
        </>
      )}

      {recurrenceEditable && (
        <>
          <p>
            {repeats === 'weekly'
              ? t(Localization.Admin.Schedule.startingFrom)
              : t(Localization.Admin.Schedule.dateLabel)}
          </p>
          <Input type="date" value={date} onChange={input => setDate(input.currentTarget.value)} />

          <p>
            {repeats === 'weekly'
              ? t(Localization.Admin.Schedule.atTime)
              : t(Localization.Admin.Schedule.timeLabel)}
          </p>
          <Input type="time" value={time} onChange={input => setTime(input.currentTarget.value)} />
        </>
      )}

      {recurrenceEditable && repeats === 'weekly' && (
        <>
          <p>{t(Localization.Admin.Schedule.endsOnOptional)}</p>
          <Input
            type="date"
            value={endsOn}
            onChange={input => setEndsOn(input.currentTarget.value)}
          />
        </>
      )}

      <p>{t(Localization.Admin.Schedule.durationLabel)}</p>
      <InputNumber
        min={1}
        max={10080}
        value={duration}
        onChange={value => setDuration(value ?? 60)}
      />

      {target.kind === 'create' ? (
        <>
          <p>{t(Localization.Admin.Schedule.timezoneLabel)}</p>
          <Select
            style={{ width: '100%' }}
            showSearch
            value={timezone}
            onChange={setTimezone}
            options={timezoneOptions}
          />
        </>
      ) : (
        // A series whose rule this form can't parse has an unknown TZID;
        // showing leftover state would be misinformation.
        (target.kind === 'edit-event' || recurrenceEditable) && (
          <p>
            {t(Localization.Admin.Schedule.timezoneLabel)}:{' '}
            {target.kind === 'edit-event' ? target.event.timezone : timezone}
          </p>
        )
      )}

      {preview.length > 0 && (
        <>
          <p>{t(Localization.Admin.Schedule.nextOccurrences)}</p>
          <ul>
            {preview.map(occurrence => (
              // Render in the rule's own timezone: these are that zone's
              // wall times, and the browser's zone would mislead.
              <li key={occurrence}>
                {new Date(occurrence).toLocaleString(undefined, { timeZone: timezone })}
              </li>
            ))}
          </ul>
        </>
      )}
      {previewError && <Paragraph type="danger">{previewError}</Paragraph>}
      {!recurrenceEditable && (
        <Paragraph type="secondary">{t(Localization.Admin.Schedule.uneditableRule)}</Paragraph>
      )}
    </Modal>
  );
};

const Schedule = () => {
  const { t } = useTranslation();
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || ({} as never);
  const scheduleConfig = serverConfig?.schedule || { enabled: false, reminderMessage: '' };

  const [schedule, setSchedule] = useState<AdminSchedule>({ series: [], events: [] });
  const [modalTarget, setModalTarget] = useState<EventModalTarget>({ kind: 'create' });
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [reminderMessage, setReminderMessage] = useState('');

  useEffect(() => {
    setReminderMessage(scheduleConfig.reminderMessage);
  }, [scheduleConfig.reminderMessage]);

  async function getSchedule() {
    try {
      const result = await fetchData(SCHEDULE_ADMIN);
      setSchedule({ series: result.series || [], events: result.events || [] });
    } catch (error) {
      console.error('error fetching schedule', error);
    }
  }

  useEffect(() => {
    getSchedule();
  }, []);

  function applyScheduleResponse(result: AdminSchedule) {
    setSchedule({ series: result.series || [], events: result.events || [] });
  }

  async function removeOrCancel(id: string, cancel: boolean) {
    try {
      const result = await fetchData(SCHEDULE_DELETE_EVENT, {
        method: 'POST',
        data: cancel ? { id, cancel: true } : { id },
      });
      applyScheduleResponse(result);
      message.success(
        cancel
          ? t(Localization.Admin.Schedule.cancelledToast)
          : t(Localization.Admin.Schedule.deletedToast),
      );
    } catch (error) {
      message.error(`${error}`);
    }
  }

  const openCreateModal = () => {
    setModalTarget({ kind: 'create' });
    setIsModalOpen(true);
  };

  const eventColumns = [
    {
      title: t(Localization.Admin.Schedule.columnName),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t(Localization.Admin.Schedule.columnWhen),
      dataIndex: 'startTime',
      key: 'startTime',
      // Render the wall time in the event's own timezone so the value
      // matches the zone label beside it, wherever the admin's browser is.
      render: (startTime: string, record: ScheduledEvent) =>
        `${new Date(startTime).toLocaleString(undefined, { timeZone: record.timezone })} (${
          record.timezone
        })`,
    },
    {
      title: t(Localization.Admin.Schedule.columnDuration),
      dataIndex: 'durationMinutes',
      key: 'durationMinutes',
      render: (minutes: number) =>
        t(Localization.Admin.Schedule.durationValue, { minutes: `${minutes}` }),
    },
    {
      title: t(Localization.Admin.Schedule.columnStatus),
      key: 'status',
      render: (_, record: ScheduledEvent) => (
        <>
          {record.status === 'cancelled' ? (
            <Tag color="red">{t(Localization.Admin.Schedule.statusCancelled)}</Tag>
          ) : (
            <Tag color="green">{t(Localization.Admin.Schedule.statusScheduled)}</Tag>
          )}
          {record.seriesId && <Tag>{t(Localization.Admin.Schedule.statusRecurring)}</Tag>}
        </>
      ),
    },
    {
      title: '',
      key: 'actions',
      render: (_, record: ScheduledEvent) => (
        <Space size="middle">
          <Button
            size="small"
            icon={<EditOutlined />}
            onClick={() => {
              setModalTarget({ kind: 'edit-event', event: record });
              setIsModalOpen(true);
            }}
          />
          {record.status !== 'cancelled' && (
            <Button size="small" onClick={() => removeOrCancel(record.id, true)}>
              {t(Localization.Admin.Schedule.cancelAction)}
            </Button>
          )}
          <Popconfirm
            title={t(Localization.Admin.Schedule.deleteEventConfirm)}
            onConfirm={() => removeOrCancel(record.id, false)}
          >
            <Button size="small" icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const seriesColumns = [
    {
      title: t(Localization.Admin.Schedule.columnName),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t(Localization.Admin.Schedule.columnRepeats),
      dataIndex: 'recurrence',
      key: 'recurrence',
      render: (recurrence: string) => {
        const rule = parseWeeklyRecurrence(recurrence);
        if (!rule) {
          return recurrence;
        }
        const days = rule.days
          .map(code => WEEKDAYS.find(day => day.code === code)?.label || code)
          .join(', ');
        return (
          <>
            {t(Localization.Admin.Schedule.weekly)} {t(Localization.Admin.Schedule.recurrenceOn)}{' '}
            {days} {t(Localization.Admin.Schedule.recurrenceAt)} {rule.time} ({rule.timezone})
            {rule.endsOn && (
              <>
                {' '}
                {t(Localization.Admin.Schedule.recurrenceUntil)} {rule.endsOn}
              </>
            )}
          </>
        );
      },
    },
    {
      title: t(Localization.Admin.Schedule.columnDuration),
      dataIndex: 'durationMinutes',
      key: 'durationMinutes',
      render: (minutes: number) =>
        t(Localization.Admin.Schedule.durationValue, { minutes: `${minutes}` }),
    },
    {
      title: '',
      key: 'actions',
      render: (_, record: ScheduledEventSeries) => (
        <Space size="middle">
          <Button
            size="small"
            icon={<EditOutlined />}
            onClick={() => {
              setModalTarget({ kind: 'edit-series', series: record });
              setIsModalOpen(true);
            }}
          />
          <Popconfirm
            title={t(Localization.Admin.Schedule.deleteSeriesConfirm)}
            onConfirm={() => removeOrCancel(record.id, false)}
          >
            <Button size="small" icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const sortedEvents = [...schedule.events].sort(
    (a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime(),
  );

  const handleReminderChange = ({ value }: UpdateArgs) => {
    setReminderMessage(value);
  };

  return (
    <div>
      <Title>{t(Localization.Admin.Schedule.title)}</Title>
      <Paragraph>{t(Localization.Admin.Schedule.pageDescription)}</Paragraph>

      <ToggleSwitch
        fieldName="enabled"
        configPath="schedule"
        apiPath={API_SCHEDULE_ENABLED}
        checked={scheduleConfig.enabled}
        useSubmit
        label={t(Localization.Admin.Schedule.enableLabel)}
        tip={t(Localization.Admin.Schedule.enableTip)}
      />

      <TextFieldWithSubmit
        fieldName="reminderMessage"
        configPath="schedule"
        apiPath={API_SCHEDULE_REMINDER_MESSAGE}
        type={TEXTFIELD_TYPE_TEXTAREA}
        label={t(Localization.Admin.Schedule.reminderLabel)}
        tip={t(Localization.Admin.Schedule.reminderTip)}
        value={reminderMessage}
        initialValue={scheduleConfig.reminderMessage}
        maxLength={500}
        onChange={handleReminderChange}
      />

      <Title level={3}>{t(Localization.Admin.Schedule.upcomingEvents)}</Title>
      <Table
        rowKey={(record: ScheduledEvent) => record.id}
        columns={eventColumns}
        dataSource={sortedEvents}
        pagination={false}
        size="small"
      />
      <br />
      <Button type="primary" onClick={openCreateModal}>
        {t(Localization.Admin.Schedule.addEvent)}
      </Button>

      <Title level={3}>{t(Localization.Admin.Schedule.recurringSchedules)}</Title>
      <Table
        rowKey={(record: ScheduledEventSeries) => record.id}
        columns={seriesColumns}
        dataSource={schedule.series}
        pagination={false}
        size="small"
      />

      <EventModal
        target={modalTarget}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        onSaved={result => {
          setIsModalOpen(false);
          applyScheduleResponse(result);
          message.success(t(Localization.Admin.Schedule.savedToast));
        }}
      />
    </div>
  );
};

Schedule.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};

export default Schedule;
