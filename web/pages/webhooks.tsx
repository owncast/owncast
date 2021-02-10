import React, { useState, useEffect } from 'react';
import {
  Table,
  Tag,
  Space,
  Button,
  Modal,
  Checkbox,
  Input,
  Typography,
  Tooltip,
  Select,
} from 'antd';
import { DeleteOutlined, EyeTwoTone, EyeInvisibleOutlined } from '@ant-design/icons';
import { isValidUrl } from '../utils/urls';

const { Title, Paragraph, Text } = Typography;
const { Option } = Select;

import { fetchData, DELETE_WEBHOOK, CREATE_WEBHOOK, WEBHOOKS } from '../utils/apis';

const availableEvents = {
  CHAT: { name: 'Chat messages', description: 'When a user sends a chat message', color: 'purple' },
  USER_JOINED: { name: 'User joined', description: 'When a user joins the chat', color: 'green' },
  NAME_CHANGE: {
    name: 'User name changed',
    description: 'When a user changes their name',
    color: 'blue',
  },
  'VISIBILITY-UPDATE': {
    name: 'Message visibility changed',
    description: 'When a message visibility changes, likely due to moderation',
    color: 'red',
  },
  STREAM_STARTED: { name: 'Stream started', description: 'When a stream starts', color: 'orange' },
  STREAM_STOPPED: { name: 'Stream stopped', description: 'When a stream stops', color: 'cyan' },
};

function convertEventStringToTag(eventString) {
  if (!eventString || !availableEvents[eventString]) {
    return null;
  }

  const event = availableEvents[eventString];

  return (
    <Tooltip key={eventString} title={event.description}>
      <Tag color={event.color}>{event.name}</Tag>
    </Tooltip>
  );
}

function NewWebhookModal(props) {
  const [selectedEvents, setSelectedEvents] = useState([]);
  const [webhookUrl, setWebhookUrl] = useState('');

  const events = Object.keys(availableEvents).map(function (key) {
    return { value: key, label: availableEvents[key].description };
  });

  function onChange(checkedValues) {
    setSelectedEvents(checkedValues);
  }

  function selectAll() {
    setSelectedEvents(Object.keys(availableEvents));
  }

  function save() {
    props.onOk(webhookUrl, selectedEvents);

    // Reset the modal
    setWebhookUrl('');
    setSelectedEvents(null);
  }

  const okButtonProps = {
    disabled: selectedEvents?.length === 0 || !isValidUrl(webhookUrl),
  };

  return (
    <Modal
      title="Create New Webhook"
      visible={props.visible}
      onOk={save}
      onCancel={props.onCancel}
      okButtonProps={okButtonProps}
    >
      <div>
        <Input
          value={webhookUrl}
          placeholder="https://myserver.com/webhook"
          onChange={input => setWebhookUrl(input.currentTarget.value)}
        />
      </div>

      <p>Select the events that will be sent to this webhook.</p>
      <Checkbox.Group options={events} value={selectedEvents} onChange={onChange} />
      <Button type="text" size="small" onClick={selectAll}>
        Select all
      </Button>
    </Modal>
  );
}

export default function Webhooks() {
  const [webhooks, setWebhooks] = useState([]);
  const [isModalVisible, setIsModalVisible] = useState(false);

  const columns = [
    {
      title: '',
      key: 'delete',
      render: (text, record) => (
        <Space size="middle">
          <Button onClick={() => handleDelete(record.id)} icon={<DeleteOutlined />} />
        </Space>
      ),
    },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
    },
    {
      title: 'Events',
      dataIndex: 'events',
      key: 'events',
      render: events => (
        <>
          {events.map(event => {
            return convertEventStringToTag(event);
          })}
        </>
      ),
    },
  ];

  const getWebhooks = async () => {
    try {
      const result = await fetchData(WEBHOOKS);
      setWebhooks(result);
    } catch (error) {
      handleError(error);
    }
  };

  useEffect(() => {
    getWebhooks();
  }, []);

  async function handleDelete(id) {
    try {
      const result = await fetchData(DELETE_WEBHOOK, { method: 'POST', data: { id: id } });
      getWebhooks();
    } catch (error) {
      handleError(error);
    }
  }

  async function handleSave(url: string, events: string[]) {
    try {
      const newHook = await fetchData(CREATE_WEBHOOK, {
        method: 'POST',
        data: { url: url, events: events },
      });
      setWebhooks(webhooks.concat(newHook));
    } catch (error) {
      handleError(error);
    }
  }

  function handleError(error) {
    console.error('error', error);
    alert(error);
  }

  const showCreateModal = () => {
    setIsModalVisible(true);
  };

  const handleModalSaveButton = (url, events) => {
    setIsModalVisible(false);
    handleSave(url, events);
  };

  const handleModalCancelButton = () => {
    setIsModalVisible(false);
  };

  return (
    <div>
      <Title>Webhooks</Title>
      <Paragraph>
        A webhook is a callback made to an external API in response to an event that takes place
        within Owncast. This can be used to build chat bots or sending automatic notifications that
        you've started streaming.
      </Paragraph>
      <Paragraph>
        Read more about how to use webhooks, with examples, at{' '}
        <a href="https://owncast.online/docs/integrations/">our documentation</a>.
      </Paragraph>

      <Table rowKey="id" columns={columns} dataSource={webhooks} pagination={false} />
      <br />
      <Button type="primary" onClick={showCreateModal}>
        Create Webhook
      </Button>
      <NewWebhookModal
        visible={isModalVisible}
        onOk={handleModalSaveButton}
        onCancel={handleModalCancelButton}
      />
    </div>
  );
}
