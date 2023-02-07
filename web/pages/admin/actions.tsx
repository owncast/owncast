import { Button, Checkbox, Form, Input, Modal, Space, Table, Typography } from 'antd';
import dynamic from 'next/dynamic';
import React, { ReactElement, useContext, useEffect, useState } from 'react';
import TextArea from 'antd/lib/input/TextArea';
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
    oldActionIndex: number | null,
    actionUrl: string,
    actionHTML: string,
    actionTitle: string,
    actionDescription: string,
    actionIcon: string,
    actionColor: string,
    openExternally: boolean,
  ) => void;
  open: boolean;
  action: ExternalAction | null;
  index: number | null;
}

const ActionModal = (props: Props) => {
  const { onOk, onCancel, open, action } = props;

  // useHTML is not part of the action object -- HTML is used if its length is > 0
  const [useHTML, setUseHTML] = useState(false);

  const [actionUrl, setActionUrl] = useState('');
  const [actionHTML, setActionHTML] = useState('');
  const [actionTitle, setActionTitle] = useState('');
  const [actionDescription, setActionDescription] = useState('');
  const [actionIcon, setActionIcon] = useState('');
  const [actionColor, setActionColor] = useState('');
  const [openExternally, setOpenExternally] = useState(false);

  useEffect(() => {
    setUseHTML((action?.html?.length || 0) > 0);
    setActionUrl(action?.url || '');
    setActionHTML(action?.html || '');
    setActionTitle(action?.title || '');
    setActionDescription(action?.description || '');
    setActionIcon(action?.icon || '');
    setActionColor(action?.color || '');
    setOpenExternally(action?.openExternally || false);
  }, [action]);

  function save() {
    onOk(
      action,
      props.index,
      // Save only one of the properties
      useHTML ? '' : actionUrl,
      useHTML ? actionHTML : '',
      actionTitle,
      actionDescription,
      actionIcon,
      actionColor,
      openExternally,
    );
    setActionUrl('');
    setActionHTML('');
    setActionTitle('');
    setActionDescription('');
    setActionIcon('');
    setActionColor('');
    setOpenExternally(false);
  }

  function canSave(): Boolean {
    if (useHTML) {
      return actionHTML !== '' && actionTitle !== '';
    }
    return isValidUrl(actionUrl, ['https:']) && actionTitle !== '';
  }

  const okButtonProps = {
    disabled: !canSave(),
  };

  const onOpenExternallyChanged = checkbox => {
    setOpenExternally(checkbox.target.checked);
  };

  const onUseHTMLChanged = checkbox => {
    setUseHTML(checkbox.target.checked);
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
        <strong>Only HTTPS URLs and embeds are supported.</strong>
        <p>
          <a
            href="https://owncast.online/thirdparty/actions/"
            target="_blank"
            rel="noopener noreferrer"
          >
            Read more about external actions.
          </a>
        </p>
        {useHTML ? (
          <Form.Item name="html">
            <TextArea
              value={actionHTML}
              placeholder="HTML embed code (required)"
              onChange={input => setActionHTML(input.currentTarget.value)}
            />
          </Form.Item>
        ) : (
          <Form.Item name="url">
            <Input
              required
              placeholder="https://myserver.com/action (required)"
              onChange={input => setActionUrl(input.currentTarget.value.trim())}
              type="url"
              pattern={DEFAULT_TEXTFIELD_URL_PATTERN}
            />
          </Form.Item>
        )}
        <Form.Item name="useHTML">
          <Checkbox checked={useHTML} defaultChecked={openExternally} onChange={onUseHTMLChanged}>
            Embed HTML instead of an URL
          </Checkbox>
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
        {useHTML ? null : (
          <Form.Item name="openExternally">
            <Checkbox
              checked={openExternally}
              defaultChecked={openExternally}
              onChange={onOpenExternallyChanged}
            >
              Open in a new tab instead of within your page.
            </Checkbox>
          </Form.Item>
        )}
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
  const [editActionIndex, setEditActionIndex] = useState(-1);

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

  async function handleDelete(action, index) {
    const actionsData = [...actions];
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
    oldActionIndex: number,
    url: string,
    html: string,
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
        html,
        title,
        description,
        icon,
        color,
        openExternally,
      };

      // Replace old action if edited or append the new action
      if (oldActionIndex >= 0) {
        actionsData[oldActionIndex] = newAction;
      } else {
        actionsData.push(newAction);
      }

      setActions(actionsData);
      await save(actionsData);
    } catch (error) {
      console.error(error);
    }
  }

  async function handleEdit(action: ExternalAction, index) {
    setEditActionIndex(index);
    setEditAction(action);
    setIsModalOpen(true);
  }

  const showCreateModal = () => {
    setEditAction(null);
    setEditActionIndex(-1);
    setIsModalOpen(true);
  };

  const handleModalSaveButton = (
    oldAction: ExternalAction | null,
    oldActionIndex: number,
    actionUrl: string,
    actionHTML: string,
    actionTitle: string,
    actionDescription: string,
    actionIcon: string,
    actionColor: string,
    openExternally: boolean,
  ) => {
    setIsModalOpen(false);
    handleSave(
      oldAction,
      oldActionIndex,
      actionUrl,
      actionHTML,
      actionTitle,
      actionDescription,
      actionIcon,
      actionColor,
      openExternally,
    );
    setEditAction(null);
    setEditActionIndex(-1);
  };

  const handleModalCancelButton = () => {
    setIsModalOpen(false);
  };

  const columns = [
    {
      title: '',
      key: 'delete-edit',
      render: (text, record, index) => (
        <Space size="middle">
          <Button onClick={() => handleDelete(record, index)} icon={<DeleteOutlined />} />
          <Button onClick={() => handleEdit(record, index)} icon={<EditOutlined />} />
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
      title: 'URL / Embed',
      key: 'url',
      dataIndex: 'url',
      render: (text, record) => (record.html ? 'HTML embed' : record.url),
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
      key: 'openExternally',
      dataIndex: 'openExternally',
      // Note embeds will alway open in the same tab.
      render: (openExternally: boolean, record) =>
        !openExternally || record.html ? 'In the same tab' : 'In a new tab',
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
        index={editActionIndex}
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
