import React, { useState, useContext, useEffect, FC } from 'react';
import { Button, Modal, Typography } from 'antd';
import { markdown, markdownLanguage } from '@codemirror/lang-markdown';
import CodeMirror, { EditorView } from '@uiw/react-codemirror';
import { bbedit } from '@uiw/codemirror-theme-bbedit';
import { languages } from '@codemirror/language-data';
import {
  TextFieldWithSubmit,
  TEXTFIELD_TYPE_TEXTAREA,
  TEXTFIELD_TYPE_URL,
} from '../../TextFieldWithSubmit';
import { ServerStatusContext } from '../../../../utils/server-status-context';
import {
  postConfigUpdateToAPI,
  TEXTFIELD_PROPS_INSTANCE_URL,
  TEXTFIELD_PROPS_SERVER_NAME,
  TEXTFIELD_PROPS_SERVER_SUMMARY,
  TEXTFIELD_PROPS_SERVER_OFFLINE_MESSAGE,
  API_YP_SWITCH,
  FIELD_PROPS_YP,
  FIELD_PROPS_NSFW,
  FIELD_PROPS_HIDE_VIEWER_COUNT,
  API_SERVER_OFFLINE_MESSAGE,
  FIELD_PROPS_DISABLE_SEARCH_INDEXING,
} from '../../../../utils/config-constants';
import { UpdateArgs } from '../../../../types/config-section';
import { ToggleSwitch } from '../../ToggleSwitch';
import { EditLogo } from '../../EditLogo';
import FormStatusIndicator from '../../FormStatusIndicator';
import { createInputStatus, STATUS_SUCCESS } from '../../../../utils/input-statuses';

const { Title } = Typography;

export type DirectoryInfoModalProps = {
  cancelPressed: () => void;
  okPressed: () => void;
};

const DirectoryInfoModal: FC<DirectoryInfoModalProps> = ({ cancelPressed, okPressed }) => (
  <Modal
    width="70%"
    title="Owncast Directory"
    visible
    onCancel={cancelPressed}
    footer={
      <div>
        <Button onClick={cancelPressed}>Do not share my server.</Button>
        <Button type="primary" onClick={okPressed}>
          Understood. Share my server publicly.
        </Button>
      </div>
    }
  >
    <Typography.Title level={3}>What is the Owncast Directory?</Typography.Title>
    <Typography.Paragraph>
      Owncast operates a public directory at{' '}
      <a href="https://directory.owncast.online">directory.owncast.online</a> to share your video
      streams with more people, while also using these as examples for others. Live streams and
      servers listed on the directory may optionally be shared on other platforms and applications.
    </Typography.Paragraph>

    <Typography.Title level={3}>Disclaimers and Responsibility</Typography.Title>
    <Typography.Paragraph>
      <ul>
        <li>
          By enabling this feature you are granting explicit rights to Owncast to share your stream
          to the public via the directory, as well as other sites, applications and any platform
          where the Owncast project may be promoting Owncast-powered streams including social media.
        </li>
        <li>
          There is no obligation to list any specific server or topic. Servers can and will be
          removed at any time for any reason.
        </li>
        <li>
          Any server that is streaming Not Safe For Work (NSFW) content and does not have the NSFW
          toggle enabled on their server will be removed.
        </li>
        <li>
          Any server streaming harmful, hurtful, misleading or hateful content in any way will not
          be listed.
        </li>
        <li>
          You may reach out to the Owncast team to report any objectionable content or content that
          you believe should not be be publicly listed.
        </li>
        <li>
          You have the right to free software and to build whatever you want with it. But there is
          no obligation for others to share it.
        </li>
      </ul>
    </Typography.Paragraph>
  </Modal>
);

// eslint-disable-next-line react/function-component-definition
export default function EditInstanceDetails() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { instanceDetails, yp, hideViewerCount, disableSearchIndexing } = serverConfig;
  const { instanceUrl } = yp;

  const [offlineMessageSaveStatus, setOfflineMessageSaveStatus] = useState(null);
  const [isDirectoryInfoModalOpen, setIsDirectoryInfoModalOpen] = useState(false);

  useEffect(() => {
    setFormDataValues({
      ...instanceDetails,
      ...yp,
      hideViewerCount,
      disableSearchIndexing,
    });
  }, [instanceDetails, yp]);

  if (!formDataValues) {
    return null;
  }

  const handleDirectorySwitchChange = (value: boolean) => {
    if (!value) {
      // Disabled. No-op.
    } else {
      setIsDirectoryInfoModalOpen(true);
    }
    setFormDataValues({
      ...formDataValues,
      yp: {
        enabled: value,
      },
    });
  };

  // if instanceUrl is empty, we should also turn OFF the `enabled` field of directory.
  const handleSubmitInstanceUrl = () => {
    if (formDataValues.instanceUrl === '') {
      if (yp.enabled === true) {
        postConfigUpdateToAPI({
          apiPath: API_YP_SWITCH,
          data: { value: false },
        });
      }
    }
  };

  const handleSaveOfflineMessage = () => {
    postConfigUpdateToAPI({
      apiPath: API_SERVER_OFFLINE_MESSAGE,
      data: { value: formDataValues.offlineMessage },
    });
    setOfflineMessageSaveStatus(createInputStatus(STATUS_SUCCESS));
    setTimeout(() => {
      setOfflineMessageSaveStatus(null);
    }, 2000);
  };

  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  };

  function handleHideViewerCountChange(enabled: boolean) {
    handleFieldChange({ fieldName: 'hideViewerCount', value: enabled });
  }

  function handleDisableSearchEngineIndexingChange(enabled: boolean) {
    handleFieldChange({ fieldName: 'disableSearchIndexing', value: enabled });
  }

  function directoryInfoModalCancelPressed() {
    setIsDirectoryInfoModalOpen(false);
    handleDirectorySwitchChange(false);
    handleFieldChange({ fieldName: 'enabled', value: false });
  }

  function directoryInfoModalOkPressed() {
    setIsDirectoryInfoModalOpen(false);
    handleFieldChange({ fieldName: 'enabled', value: true });
    setFormDataValues({
      ...formDataValues,
      yp: {
        enabled: true,
      },
    });
  }

  const hasInstanceUrl = instanceUrl !== '';

  return (
    <div className="edit-general-settings">
      <Title level={3} className="section-title">
        Configure Instance Details
      </Title>
      <br />

      <TextFieldWithSubmit
        fieldName="name"
        {...TEXTFIELD_PROPS_SERVER_NAME}
        value={formDataValues.name}
        initialValue={instanceDetails.name}
        onChange={handleFieldChange}
      />

      <TextFieldWithSubmit
        fieldName="instanceUrl"
        {...TEXTFIELD_PROPS_INSTANCE_URL}
        value={formDataValues.instanceUrl}
        initialValue={yp.instanceUrl}
        type={TEXTFIELD_TYPE_URL}
        onChange={handleFieldChange}
        onSubmit={handleSubmitInstanceUrl}
      />

      <TextFieldWithSubmit
        fieldName="summary"
        {...TEXTFIELD_PROPS_SERVER_SUMMARY}
        type={TEXTFIELD_TYPE_TEXTAREA}
        value={formDataValues.summary}
        initialValue={instanceDetails.summary}
        onChange={handleFieldChange}
      />

      <div style={{ marginBottom: '50px', marginRight: '150px' }}>
        <div
          style={{
            display: 'flex',
            width: '80vh',
            justifyContent: 'space-between',
            alignItems: 'end',
          }}
        >
          <p style={{ margin: '20px', marginRight: '10px', fontWeight: '400' }}>Offline Message:</p>
          <CodeMirror
            value={formDataValues.offlineMessage}
            {...TEXTFIELD_PROPS_SERVER_OFFLINE_MESSAGE}
            theme={bbedit}
            height="150px"
            width="450px"
            onChange={value => {
              handleFieldChange({ fieldName: 'offlineMessage', value });
            }}
            extensions={[
              markdown({ base: markdownLanguage, codeLanguages: languages }),
              EditorView.lineWrapping,
            ]}
          />
        </div>
        <div className="field-tip">
          The offline message is displayed to your page visitors when you&apos;re not streaming.
          Markdown is supported.
        </div>

        <Button
          type="primary"
          onClick={handleSaveOfflineMessage}
          style={{ margin: '10px', float: 'right' }}
        >
          Save Message
        </Button>
        <FormStatusIndicator status={offlineMessageSaveStatus} />
      </div>

      {/* Logo section */}
      <EditLogo />

      <ToggleSwitch
        fieldName="hideViewerCount"
        useSubmit
        {...FIELD_PROPS_HIDE_VIEWER_COUNT}
        checked={formDataValues.hideViewerCount}
        onChange={handleHideViewerCountChange}
      />

      <ToggleSwitch
        fieldName="disableSearchIndexing"
        useSubmit
        {...FIELD_PROPS_DISABLE_SEARCH_INDEXING}
        checked={formDataValues.disableSearchIndexing}
        onChange={handleDisableSearchEngineIndexingChange}
      />

      <br />
      <p className="description">
        Increase your audience by appearing in the{' '}
        <a href="https://directory.owncast.online" target="_blank" rel="noreferrer">
          <strong>Owncast Directory</strong>
        </a>
        . This is an external service run by the Owncast project.{' '}
        <a
          href="https://owncast.online/docs/directory/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn more
        </a>
        .
      </p>
      {!yp.instanceUrl && (
        <p className="description">
          You must set your <strong>Server URL</strong> above to enable the directory.
        </p>
      )}

      <div className="config-yp-container">
        <ToggleSwitch
          fieldName="enabled"
          useSubmit
          {...FIELD_PROPS_YP}
          checked={formDataValues.enabled}
          disabled={!hasInstanceUrl}
          onChange={handleDirectorySwitchChange}
        />
        <ToggleSwitch
          fieldName="nsfw"
          useSubmit
          {...FIELD_PROPS_NSFW}
          checked={formDataValues.nsfw}
          disabled={!hasInstanceUrl}
        />
      </div>
      {isDirectoryInfoModalOpen && (
        <DirectoryInfoModal
          cancelPressed={directoryInfoModalCancelPressed}
          okPressed={directoryInfoModalOkPressed}
        />
      )}
    </div>
  );
}
