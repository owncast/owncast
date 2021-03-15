// comment

import React, { useState, useEffect, useContext } from 'react';
import { Table, Space, Button, Modal, Checkbox, Input, Typography } from 'antd';
import { ServerStatusContext } from '../utils/server-status-context';
import { DeleteOutlined } from '@ant-design/icons';
import isValidUrl from '../utils/urls';
import FormStatusIndicator from '../components/config/form-status-indicator';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../utils/input-statuses';

import {
  postConfigUpdateToAPI,
  API_EXTERNAL_ACTIONS,
  RESET_TIMEOUT,
} from '../utils/config-constants';

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
  }

  const okButtonProps = {
    disabled: !isValidUrl(actionUrl) || actionTitle === '',
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
        <p>
          <Input
            value={actionUrl}
            required
            placeholder="https://myserver.com/action (required)"
            onChange={input => setActionUrl(input.currentTarget.value)}
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
      render: (url: string) => {
        return url ? <img src={url} /> : null;
      },
    },
    {
      title: 'Color',
      dataIndex: 'color',
      key: 'color',
      render: (color: string) => {
        return color ? (<div style={{backgroundColor: color, height: '30px'}}>{color}</div>) : null;
      },
    },
    {
      title: 'Opens',
      dataIndex: 'openExternally',
      key: 'openExternally',
      render: (openExternally: boolean) => {
        return openExternally ? 'In a new tab' : 'In a modal';
      },
    },
  ];

  async function handleDelete(action) {
    let actionsData = [...actions];
    const index = actions.findIndex(item => item.url === action.url);
    actionsData.splice(index, 1);

    setActions(actionsData);
    save(actionsData);
    try {
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
      let actionsData = [...actions];
      const updatedActions = actionsData.concat({ url, title, description, icon, color, openExternally });
      setActions(updatedActions);
      await save(updatedActions);
    } catch (error) {
      console.error(error);
    }
  }

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

  return (
    <div>
      <Title>External Actions</Title>
      <Paragraph>Description goes here.</Paragraph>
      <Paragraph>
        Read more about how to use actions, with examples, at{' '}
        <a href="https://owncast.online/thirdparty/?source=admin">our documentation</a>.
      </Paragraph>

      <Table rowKey="id" columns={columns} dataSource={actions} pagination={false} />
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
