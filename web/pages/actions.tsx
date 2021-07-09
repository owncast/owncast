import { DeleteOutlined } from '@ant-design/icons';
import { Button, Checkbox, Input, Modal, Space, Table, Typography } from 'antd';
import React, { useContext, useEffect, useState } from 'react';
import FormStatusIndicator from '../components/config/form-status-indicator';
import {
  API_EXTERNAL_ACTIONS,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
} from '../utils/config-constants';
import { createInputStatus, STATUS_ERROR, STATUS_SUCCESS } from '../utils/input-statuses';
import { ServerStatusContext } from '../utils/server-status-context';
import isValidUrl, { DEFAULT_TEXTFIELD_URL_PATTERN } from '../utils/urls';

const { Title, Paragraph } = Typography;
let resetTimer = null;

interface Props {
  onCancel: () => void;
  onOk: any; // todo: make better type
  visible: boolean;
}

function NewActionModal(props: Props) {
  const { onOk, onCancel, visible } = props;

  const [actionUrl, setActionUrl] = useState('');
  const [actionTitle, setActionTitle] = useState('');
  const [actionDescription, setActionDescription] = useState('');
  const [actionIcon, setActionIcon] = useState('');
  const [actionColor, setActionColor] = useState('');
  const [openExternally, setOpenExternally] = useState(false);

  function save() {
    onOk(actionUrl, actionTitle, actionDescription, actionIcon, actionColor, openExternally);
    setActionUrl('');
    setActionTitle('');
    setActionDescription('');
    setActionIcon('');
    setActionColor('');
    setOpenExternally(false);
  }

  function canSave(): Boolean {
    try {
      const validationObject = new URL(actionUrl);
      if (validationObject.protocol !== 'https:') {
        return false;
      }
    } catch {
      return false;
    }

    return isValidUrl(actionUrl) && actionTitle !== '';
  }

  const okButtonProps = {
    disabled: !canSave(),
  };

  const onOpenExternallyChanged = checkbox => {
    setOpenExternally(checkbox.target.checked);
  };

  return (
    <Modal
      title="Create New Action"
      visible={visible}
      onOk={save}
      onCancel={onCancel}
      okButtonProps={okButtonProps}
    >
      <div>
        Add the URL for the external action you want to present.{' '}
        <strong>Only HTTPS urls are supported.</strong>
        <p>
          <a href="https://owncast.online/docs" target="_blank" rel="noopener noreferrer">
            Read more about external actions.
          </a>
        </p>
        <p>
          <Input
            value={actionUrl}
            required
            placeholder="https://myserver.com/action (required)"
            onChange={input => setActionUrl(input.currentTarget.value.trim())}
            type="url"
            pattern={DEFAULT_TEXTFIELD_URL_PATTERN}
          />
        </p>
        <p>
          <Input
            value={actionTitle}
            required
            placeholder="Your action title (required)"
            onChange={input => setActionTitle(input.currentTarget.value)}
          />
        </p>
        <p>
          <Input
            value={actionDescription}
            placeholder="Optional description"
            onChange={input => setActionDescription(input.currentTarget.value)}
          />
        </p>
        <p>
          <Input
            value={actionIcon}
            placeholder="https://myserver.com/action/icon.png (optional)"
            onChange={input => setActionIcon(input.currentTarget.value)}
          />
        </p>
        <p>
          <Input
            type="color"
            value={actionColor}
            onChange={input => setActionColor(input.currentTarget.value)}
          />
          Optional background color of the action button.
        </p>
        <Checkbox
          checked={openExternally}
          defaultChecked={openExternally}
          onChange={onOpenExternallyChanged}
        >
          Open in a new tab instead of within your page.
        </Checkbox>
      </div>
    </Modal>
  );
}

export default function Actions() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { externalActions } = serverConfig;
  const [actions, setActions] = useState([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [submitStatus, setSubmitStatus] = useState(null);

  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  useEffect(() => {
    setActions(externalActions || []);
  }, [externalActions]);

  async function save(actionsData) {
    await postConfigUpdateToAPI({
      apiPath: API_EXTERNAL_ACTIONS,
      data: { value: actionsData },
      onSuccess: () => {
        setFieldInConfigState({ fieldName: 'externalActions', value: actionsData, path: '' });
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Updated.'));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        console.log(message);
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  }

  async function handleDelete(action) {
    const actionsData = [...actions];
    const index = actions.findIndex(item => item.url === action.url);
    actionsData.splice(index, 1);

    try {
      setActions(actionsData);
      save(actionsData);
    } catch (error) {
      console.error(error);
    }
  }

  async function handleSave(
    url: string,
    title: string,
    description: string,
    icon: string,
    color: string,
    openExternally: boolean,
  ) {
    try {
      const actionsData = [...actions];
      const updatedActions = actionsData.concat({
        url,
        title,
        description,
        icon,
        color,
        openExternally,
      });
      setActions(updatedActions);
      await save(updatedActions);
    } catch (error) {
      console.error(error);
    }
  }

  const showCreateModal = () => {
    setIsModalVisible(true);
  };

  const handleModalSaveButton = (
    actionUrl: string,
    actionTitle: string,
    actionDescription: string,
    actionIcon: string,
    actionColor: string,
    openExternally: boolean,
  ) => {
    setIsModalVisible(false);
    handleSave(actionUrl, actionTitle, actionDescription, actionIcon, actionColor, openExternally);
  };

  const handleModalCancelButton = () => {
    setIsModalVisible(false);
  };

  const columns = [
    {
      title: '',
      key: 'delete',
      render: (text, record) => (
        <Space size="middle">
          <Button onClick={() => handleDelete(record)} icon={<DeleteOutlined />} />
        </Space>
      ),
    },
    {
      title: 'Name',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
    },
    {
      title: 'Icon',
      dataIndex: 'icon',
      key: 'icon',
      render: (url: string) => (url ? <img style={{ width: '2vw' }} src={url} alt="" /> : null),
    },
    {
      title: 'Color',
      dataIndex: 'color',
      key: 'color',
      render: (color: string) =>
        color ? <div style={{ backgroundColor: color, height: '30px' }}>{color}</div> : null,
    },
    {
      title: 'Opens',
      dataIndex: 'openExternally',
      key: 'openExternally',
      render: (openExternally: boolean) => (openExternally ? 'In a new tab' : 'In a modal'),
    },
  ];

  return (
    <div>
      <Title>External Actions</Title>
      <Paragraph>
        External action URLs are 3rd party UI you can display, embedded, into your Owncast page when
        a user clicks on a button to launch your action.
      </Paragraph>
      <Paragraph>
        Read more about how to use actions, with examples, at{' '}
        <a
          href="https://owncast.online/thirdparty/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          our documentation
        </a>
        .
      </Paragraph>

      <Table
        rowKey={record => `${record.title}-${record.url}`}
        columns={columns}
        dataSource={actions}
        pagination={false}
      />
      <br />
      <Button type="primary" onClick={showCreateModal}>
        Create New Action
      </Button>
      <FormStatusIndicator status={submitStatus} />

      <NewActionModal
        visible={isModalVisible}
        onOk={handleModalSaveButton}
        onCancel={handleModalCancelButton}
      />
    </div>
  );
}
