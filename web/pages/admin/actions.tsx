import { Button, Checkbox, Form, Input, Modal, Space, Table, Typography } from 'antd';
import _ from 'lodash';
import dynamic from 'next/dynamic';
import React, { ReactElement, useContext, useEffect, useState } from 'react';
import { FormStatusIndicator } from '../../components/admin/FormStatusIndicator';
import { ExternalAction } from '../../interfaces/external-action';
import {
  API_EXTERNAL_ACTIONS,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
} from '../../utils/config-constants';
import { createInputStatus, STATUS_ERROR, STATUS_SUCCESS } from '../../utils/input-statuses';
import { ServerStatusContext } from '../../utils/server-status-context';
import { isValidUrl, DEFAULT_TEXTFIELD_URL_PATTERN } from '../../utils/urls';

import { AdminLayout } from '../../components/layouts/AdminLayout';

const { Title, Paragraph } = Typography;

// Lazy loaded components

const DeleteOutlined = dynamic(() => import('@ant-design/icons/DeleteOutlined'), {
  ssr: false,
});

const EditOutlined = dynamic(() => import('@ant-design/icons/EditOutlined'), {
  ssr: false,
});

let resetTimer = null;

interface Props {
  onCancel: () => void;
  onOk: (
    oldAction: ExternalAction | null,
    actionUrl: string,
    actionTitle: string,
    actionDescription: string,
    actionIcon: string,
    actionColor: string,
    openExternally: boolean,
  ) => void;
  open: boolean;
  action: ExternalAction | null;
}

const ActionModal = (props: Props) => {
  const { onOk, onCancel, open, action } = props;

  const [actionUrl, setActionUrl] = useState('');
  const [actionTitle, setActionTitle] = useState('');
  const [actionDescription, setActionDescription] = useState('');
  const [actionIcon, setActionIcon] = useState('');
  const [actionColor, setActionColor] = useState('');
  const [openExternally, setOpenExternally] = useState(false);

  useEffect(() => {
    setActionUrl(action?.url || '');
    setActionTitle(action?.title || '');
    setActionDescription(action?.description || '');
    setActionIcon(action?.icon || '');
    setActionColor(action?.color || '');
    setOpenExternally(action?.openExternally || false);
  }, [action]);

  function save() {
    onOk(
      action,
      actionUrl,
      actionTitle,
      actionDescription,
      actionIcon,
      actionColor,
      openExternally,
    );
    setActionUrl('');
    setActionTitle('');
    setActionDescription('');
    setActionIcon('');
    setActionColor('');
    setOpenExternally(false);
  }

  function canSave(): Boolean {
    return isValidUrl(actionUrl, ['https:']) && actionTitle !== '';
  }

  const okButtonProps = {
    disabled: !canSave(),
  };

  const onOpenExternallyChanged = checkbox => {
    setOpenExternally(checkbox.target.checked);
  };

  return (
    <Modal
      destroyOnClose
      title={action == null ? 'Create New Action' : 'Edit Action'}
      open={open}
      onOk={save}
      onCancel={onCancel}
      okButtonProps={okButtonProps}
    >
      <Form initialValues={action}>
        Add the URL for the external action you want to present.{' '}
        <strong>Only HTTPS urls are supported.</strong>
        <p>
          <a
            href="https://owncast.online/thirdparty/actions/"
            target="_blank"
            rel="noopener noreferrer"
          >
            Read more about external actions.
          </a>
        </p>
        <Form.Item name="url">
          <Input
            required
            placeholder="https://myserver.com/action (required)"
            onChange={input => setActionUrl(input.currentTarget.value.trim())}
            type="url"
            pattern={DEFAULT_TEXTFIELD_URL_PATTERN}
          />
        </Form.Item>
        <Form.Item name="title">
          <Input
            value={actionTitle}
            required
            placeholder="Your action title (required)"
            onChange={input => setActionTitle(input.currentTarget.value)}
          />
        </Form.Item>
        <Form.Item name="description">
          <Input
            value={actionDescription}
            placeholder="Optional description"
            onChange={input => setActionDescription(input.currentTarget.value)}
          />
        </Form.Item>
        <Form.Item name="icon">
          <Input
            value={actionIcon}
            placeholder="https://myserver.com/action/icon.png (optional)"
            onChange={input => setActionIcon(input.currentTarget.value)}
          />
        </Form.Item>
        <div>
          <Form.Item name="color" style={{ marginBottom: '0px' }}>
            <Input
              type="color"
              value={actionColor}
              onChange={input => setActionColor(input.currentTarget.value)}
            />
          </Form.Item>
          Optional background color of the action button.
        </div>
        <Form.Item name="openExternally">
          <Checkbox
            checked={openExternally}
            defaultChecked={openExternally}
            onChange={onOpenExternallyChanged}
          >
            Open in a new tab instead of within your page.
          </Checkbox>
        </Form.Item>
      </Form>
    </Modal>
  );
};

const Actions = () => {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { externalActions } = serverConfig;
  const [actions, setActions] = useState<ExternalAction[]>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [submitStatus, setSubmitStatus] = useState(null);
  const [editAction, setEditAction] = useState<ExternalAction>(null);

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
    oldAction: ExternalAction | null,
    url: string,
    title: string,
    description: string,
    icon: string,
    color: string,
    openExternally: boolean,
  ) {
    try {
      const actionsData = [...actions];

      const newAction: ExternalAction = {
        url,
        title,
        description,
        icon,
        color,
        openExternally,
      };

      // Replace old action if edited or append the new action
      const index = oldAction ? actions.findIndex(item => _.isEqual(item, oldAction)) : -1;
      if (index >= 0) {
        actionsData[index] = newAction;
      } else {
        actionsData.push(newAction);
      }

      setActions(actionsData);
      await save(actionsData);
    } catch (error) {
      console.error(error);
    }
  }

  async function handleEdit(action: ExternalAction) {
    setEditAction(action);
    setIsModalOpen(true);
  }

  const showCreateModal = () => {
    setEditAction(null);
    setIsModalOpen(true);
  };

  const handleModalSaveButton = (
    oldAction: ExternalAction | null,
    actionUrl: string,
    actionTitle: string,
    actionDescription: string,
    actionIcon: string,
    actionColor: string,
    openExternally: boolean,
  ) => {
    setIsModalOpen(false);
    handleSave(
      oldAction,
      actionUrl,
      actionTitle,
      actionDescription,
      actionIcon,
      actionColor,
      openExternally,
    );
    setEditAction(null);
  };

  const handleModalCancelButton = () => {
    setIsModalOpen(false);
  };

  const columns = [
    {
      title: '',
      key: 'delete-edit',
      render: (text, record) => (
        <Space size="middle">
          <Button onClick={() => handleDelete(record)} icon={<DeleteOutlined />} />
          <Button onClick={() => handleEdit(record)} icon={<EditOutlined />} />
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

      <ActionModal
        action={editAction}
        open={isModalOpen}
        onOk={handleModalSaveButton}
        onCancel={handleModalCancelButton}
      />
    </div>
  );
};

Actions.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
export default Actions;
