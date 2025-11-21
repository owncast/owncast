import { Button, Checkbox, Form, Input, Modal, Select, Space, Table, Typography } from 'antd';
import CodeMirror from '@uiw/react-codemirror';
import { bbedit } from '@uiw/codemirror-theme-bbedit';
import { html as codeMirrorHTML } from '@codemirror/lang-html';
import dynamic from 'next/dynamic';
import React, { ReactElement, useContext, useEffect, useState } from 'react';
import { useTranslation } from 'next-export-i18n';
import { FormStatusIndicator } from '../../components/admin/FormStatusIndicator';
import { ExternalAction } from '../../interfaces/external-action';
import {
  API_EXTERNAL_ACTIONS,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
} from '../../utils/config-constants';
import { createInputStatus, STATUS_ERROR, STATUS_SUCCESS } from '../../utils/input-statuses';
import { ServerStatusContext } from '../../utils/server-status-context';
import { isValidUrl, DEFAULT_TEXTFIELD_URL_PATTERN } from '../../utils/validators';
import { Localization } from '../../types/localization';
import { Translation } from '../../components/ui/Translation/Translation';

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

// ActionType is only used here to save either only the URL or only the HTML.
type ActionType = 'url' | 'html';

const ActionModal = (props: Props) => {
  const { onOk, onCancel, open, action } = props;
  const { t } = useTranslation();

  const [actionType, setActionType] = useState<ActionType>('url');

  const [actionUrl, setActionUrl] = useState('');
  const [actionHTML, setActionHTML] = useState('');
  const [actionTitle, setActionTitle] = useState('');
  const [actionDescription, setActionDescription] = useState('');
  const [actionIcon, setActionIcon] = useState('');
  const [actionColor, setActionColor] = useState('');
  const [openExternally, setOpenExternally] = useState(false);

  useEffect(() => {
    setActionType((action?.html?.length || 0) > 0 ? 'html' : 'url');
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
      actionType === 'html' ? '' : actionUrl,
      actionType === 'html' ? actionHTML : '',
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
    if (actionType === 'html') {
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

  const onActionHTMLChanged = (newActionHTML: string) => {
    setActionHTML(newActionHTML);
  };

  return (
    <Modal
      destroyOnClose
      title={
        action == null
          ? t(Localization.Admin.Actions.createNewActionTitle)
          : t(Localization.Admin.Actions.editActionTitle)
      }
      open={open}
      onOk={save}
      onCancel={onCancel}
      okButtonProps={okButtonProps}
    >
      <Form initialValues={action}>
        <Translation translationKey={Localization.Admin.Actions.modalDescription} />{' '}
        <strong>
          <Translation translationKey={Localization.Admin.Actions.onlyHttpsSupported} />
        </strong>
        <p>
          <a
            href="https://owncast.online/thirdparty/actions/"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Translation translationKey={Localization.Admin.Actions.readMoreAboutActions} />
          </a>
        </p>
        <Form.Item>
          <Select
            value={actionType}
            onChange={setActionType}
            placeholder={t(Localization.Admin.Actions.selectActionType)}
            options={[
              { label: t(Localization.Admin.Actions.linkOrEmbedUrl), value: 'url' },
              { label: t(Localization.Admin.Actions.customHtml), value: 'html' },
            ]}
          />
        </Form.Item>
        {actionType === 'html' ? (
          <Form.Item name="html">
            <CodeMirror
              value={actionHTML}
              placeholder={t(Localization.Admin.Actions.htmlEmbedPlaceholder)}
              theme={bbedit}
              height="200px"
              extensions={[codeMirrorHTML()]}
              onChange={onActionHTMLChanged}
            />
          </Form.Item>
        ) : (
          <Form.Item name="url">
            <Input
              required
              placeholder={t(Localization.Admin.Actions.urlPlaceholder)}
              onChange={input => setActionUrl(input.currentTarget.value.trim())}
              type="url"
              pattern={DEFAULT_TEXTFIELD_URL_PATTERN}
            />
          </Form.Item>
        )}
        <Form.Item name="title">
          <Input
            value={actionTitle}
            required
            placeholder={t(Localization.Admin.Actions.titlePlaceholder)}
            onChange={input => setActionTitle(input.currentTarget.value)}
          />
        </Form.Item>
        <Form.Item name="description">
          <Input
            value={actionDescription}
            placeholder={t(Localization.Admin.Actions.descriptionPlaceholder)}
            onChange={input => setActionDescription(input.currentTarget.value)}
          />
        </Form.Item>
        <Form.Item name="icon">
          <Input
            value={actionIcon}
            placeholder={t(Localization.Admin.Actions.iconPlaceholder)}
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
          <Translation translationKey={Localization.Admin.Actions.optionalBackgroundColor} />
        </div>
        {actionType === 'html' ? null : (
          <Form.Item name="openExternally">
            <Checkbox
              checked={openExternally}
              defaultChecked={openExternally}
              onChange={onOpenExternallyChanged}
            >
              <Translation translationKey={Localization.Admin.Actions.openExternally} />
            </Checkbox>
          </Form.Item>
        )}
      </Form>
    </Modal>
  );
};

const Actions = () => {
  const { t } = useTranslation();
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
        setSubmitStatus(
          createInputStatus(STATUS_SUCCESS, t(Localization.Admin.StatusMessages.updated)),
        );
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        console.log(message);
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  }

  async function handleDelete(index) {
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
    _oldAction: ExternalAction | null,
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
      render: (_, record, index) => (
        <Space size="middle">
          <Button onClick={() => handleDelete(index)} icon={<DeleteOutlined />} />
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
      render: (_, record) => (record.html ? 'HTML embed' : record.url),
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
      // Note: embeds will always open in the same tab / in a modal
      render: (openExternally: boolean, record) =>
        !openExternally || record.html ? 'In the same tab' : 'In a new tab',
    },
  ];

  return (
    <div>
      <Title>
        <Translation translationKey={Localization.Admin.Actions.title} />
      </Title>
      <Paragraph>
        <Translation translationKey={Localization.Admin.Actions.description} />
      </Paragraph>
      <Paragraph>
        <Translation translationKey={Localization.Admin.Actions.readMoreLink} />{' '}
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
        <Translation translationKey={Localization.Admin.Actions.createNewAction} />
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
