/* eslint-disable react/no-unescaped-entities */
import { Typography, Modal, Button, Row, Col, Alert } from 'antd';
import React, { ReactElement, useContext, useEffect, useState, FC } from 'react';
import {
  TEXTFIELD_TYPE_TEXT,
  TEXTFIELD_TYPE_TEXTAREA,
  TEXTFIELD_TYPE_URL,
} from '../../components/admin/TextField';
import { TextFieldWithSubmit } from '../../components/admin/TextFieldWithSubmit';
import { ToggleSwitch } from '../../components/admin/ToggleSwitch';
import { EditValueArray } from '../../components/admin/EditValueArray';
import { UpdateArgs } from '../../types/config-section';
import {
  FIELD_PROPS_ENABLE_FEDERATION,
  TEXTFIELD_PROPS_FEDERATION_LIVE_MESSAGE,
  TEXTFIELD_PROPS_FEDERATION_DEFAULT_USER,
  FIELD_PROPS_FEDERATION_IS_PRIVATE,
  FIELD_PROPS_SHOW_FEDERATION_ENGAGEMENT,
  TEXTFIELD_PROPS_FEDERATION_INSTANCE_URL,
  FIELD_PROPS_FEDERATION_BLOCKED_DOMAINS,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  API_FEDERATION_BLOCKED_DOMAINS,
  FIELD_PROPS_FEDERATION_NSFW,
} from '../../utils/config-constants';
import { ServerStatusContext } from '../../utils/server-status-context';
import { createInputStatus, STATUS_ERROR, STATUS_SUCCESS } from '../../utils/input-statuses';

import { AdminLayout } from '../../components/layouts/AdminLayout';

export type FederationInfoModalProps = {
  cancelPressed: () => void;
  okPressed: () => void;
};

const FederationInfoModal: FC<FederationInfoModalProps> = ({ cancelPressed, okPressed }) => (
  <Modal
    width="70%"
    title="Enable Social Features"
    visible
    onCancel={cancelPressed}
    footer={
      <div>
        <Button onClick={cancelPressed}>Do not enable</Button>
        <Button type="primary" onClick={okPressed}>
          Enable Social Features
        </Button>
      </div>
    }
  >
    <Typography.Title level={3}>How do Owncast's social features work?</Typography.Title>
    <Typography.Paragraph>
      Owncast's social features are accomplished by having your server join The{' '}
      <a href="https://en.wikipedia.org/wiki/Fediverse" rel="noopener noreferrer" target="_blank">
        Fediverse
      </a>
      , a decentralized, open, collection of independent servers, like yours.
    </Typography.Paragraph>
    Please{' '}
    <a href="https://owncast.online/docs/social" rel="noopener noreferrer" target="_blank">
      read more
    </a>{' '}
    about these features, the details behind them, and how they work.
    <Typography.Paragraph />
    <Typography.Title level={3}>What do you need to know?</Typography.Title>
    <ul>
      <li>
        These features are brand new. Given the variability of interfacing with the rest of the
        world, bugs are possible. Please report anything that you think isn't working quite right.
      </li>
      <li>You must always host your Owncast server with SSL using a https url.</li>
      <li>
        You should not change your server name URL or social username once people begin following
        you, as your server will be seen as a completely different user on the Fediverse, and the
        old user will disappear.
      </li>
      <li>
        Turning on <i>Private mode</i> will allow you to manually approve each follower and limit
        the visibility of your posts to followers only.
      </li>
    </ul>
    <Typography.Title level={3}>Learn more about The Fediverse</Typography.Title>
    <Typography.Paragraph>
      If these concepts are new you should discover more about what this functionality has to offer.
      Visit{' '}
      <a href="https://owncast.online/docs/social" rel="noopener noreferrer" target="_blank">
        our documentation
      </a>{' '}
      to be pointed at some resources that will help get you started on The Fediverse.
    </Typography.Paragraph>
  </Modal>
);

const ConfigFederation = () => {
  const { Title } = Typography;
  const [formDataValues, setFormDataValues] = useState(null);
  const [isInfoModalOpen, setIsInfoModalOpen] = useState(false);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const [blockedDomainSaveState, setBlockedDomainSaveState] = useState(null);

  const { federation, yp, instanceDetails } = serverConfig;
  const { enabled, isPrivate, username, goLiveMessage, showEngagement, blockedDomains } =
    federation;
  const { instanceUrl } = yp;
  const { nsfw } = instanceDetails;

  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  };

  const handleUsernameChange = ({ fieldName, value }: UpdateArgs) => {
    handleFieldChange({
      fieldName,
      value,
    });
    setFormDataValues({
      ...formDataValues,
      username: value.replace(/\W/g, ''),
    });
  };

  const handleEnabledSwitchChange = (value: boolean) => {
    if (!value) {
      setFormDataValues({
        ...formDataValues,
        enabled: false,
      });
    } else {
      setIsInfoModalOpen(true);
    }
  };

  // if instanceUrl is empty, we should also turn OFF the `enabled` field of directory.
  const handleSubmitInstanceUrl = () => {
    const hasInstanceUrl = formDataValues.instanceUrl !== '';
    const isInstanceUrlSecure = formDataValues.instanceUrl.startsWith('https://');

    if (!hasInstanceUrl || !isInstanceUrlSecure) {
      postConfigUpdateToAPI({
        apiPath: FIELD_PROPS_ENABLE_FEDERATION.apiPath,
        data: { value: false },
      });
      setFormDataValues({
        ...formDataValues,
        enabled: false,
      });
    }
  };

  function federationInfoModalCancelPressed() {
    setIsInfoModalOpen(false);
    setFormDataValues({
      ...formDataValues,
      enabled: false,
    });
  }

  function federationInfoModalOkPressed() {
    setIsInfoModalOpen(false);
    setFormDataValues({
      ...formDataValues,
      enabled: true,
    });
  }

  function resetBlockedDomainsSaveState() {
    setBlockedDomainSaveState(null);
  }

  function saveBlockedDomains() {
    try {
      postConfigUpdateToAPI({
        apiPath: API_FEDERATION_BLOCKED_DOMAINS,
        data: { value: formDataValues.blockedDomains },
        onSuccess: () => {
          setFieldInConfigState({
            fieldName: 'forbiddenUsernames',
            value: formDataValues.forbiddenUsernames,
          });
          setBlockedDomainSaveState(STATUS_SUCCESS);
          setTimeout(resetBlockedDomainsSaveState, RESET_TIMEOUT);
        },
        onError: (message: string) => {
          setBlockedDomainSaveState(createInputStatus(STATUS_ERROR, message));
          setTimeout(resetBlockedDomainsSaveState, RESET_TIMEOUT);
        },
      });
    } catch (e) {
      console.error(e);
      setBlockedDomainSaveState(STATUS_ERROR);
    }
  }

  function handleDeleteBlockedDomain(index: number) {
    formDataValues.blockedDomains.splice(index, 1);
    saveBlockedDomains();
  }

  function handleCreateBlockedDomain(domain: string) {
    let newDomain;
    try {
      const u = new URL(domain);
      newDomain = u.host;
    } catch {
      newDomain = domain;
    }

    formDataValues.blockedDomains.push(newDomain);
    handleFieldChange({
      fieldName: 'blockedDomains',
      value: formDataValues.blockedDomains,
    });
    saveBlockedDomains();
  }

  useEffect(() => {
    setFormDataValues({
      enabled,
      isPrivate,
      username,
      goLiveMessage,
      showEngagement,
      blockedDomains,
      nsfw,
      instanceUrl: yp.instanceUrl,
    });
  }, [serverConfig, yp]);

  if (!formDataValues) {
    return null;
  }

  const hasInstanceUrl = instanceUrl !== '';
  const isInstanceUrlSecure = instanceUrl.startsWith('https://');
  const configurationWarning = !isInstanceUrlSecure && (
    <>
      <Alert
        message="You must set your server URL before you can enable this feature."
        type="warning"
        showIcon
      />
      <br />
      <TextFieldWithSubmit
        fieldName="instanceUrl"
        {...TEXTFIELD_PROPS_FEDERATION_INSTANCE_URL}
        value={formDataValues.instanceUrl}
        initialValue={yp.instanceUrl}
        type={TEXTFIELD_TYPE_URL}
        onChange={handleFieldChange}
        onSubmit={handleSubmitInstanceUrl}
        required
      />
    </>
  );

  const invalidPortWarning = (
    <Alert
      message="Only Owncast instances available on the default SSL port 443 support this feature."
      type="warning"
      showIcon
    />
  );

  const hasInvalidPort =
    instanceUrl && new URL(instanceUrl).port !== '' && new URL(instanceUrl).port !== '443';

  return (
    <div>
      <Title>Configure Social Features</Title>
      <p>
        Owncast provides the ability for people to follow and engage with your instance. It's a
        great way to promote alerting, sharing and engagement of your stream.
      </p>
      <p>
        Once enabled you'll alert your followers when you go live as well as gain the ability to
        compose custom posts to share any information you like.
      </p>
      <p>
        <a href="https://owncast.online/docs/social" rel="noopener noreferrer" target="_blank">
          Read more about the specifics of these social features.
        </a>
      </p>
      <Row>
        <Col span={15} className="form-module" style={{ marginRight: '15px' }}>
          {configurationWarning}
          {hasInvalidPort && invalidPortWarning}
          <ToggleSwitch
            fieldName="enabled"
            onChange={handleEnabledSwitchChange}
            {...FIELD_PROPS_ENABLE_FEDERATION}
            checked={formDataValues.enabled}
            disabled={hasInvalidPort || !hasInstanceUrl || !isInstanceUrlSecure}
          />
          <ToggleSwitch
            fieldName="isPrivate"
            {...FIELD_PROPS_FEDERATION_IS_PRIVATE}
            checked={formDataValues.isPrivate}
            disabled={!enabled}
          />
          <ToggleSwitch
            fieldName="nsfw"
            useSubmit
            {...FIELD_PROPS_FEDERATION_NSFW}
            checked={formDataValues.nsfw}
            disabled={hasInvalidPort || !hasInstanceUrl}
          />
          <TextFieldWithSubmit
            required
            fieldName="username"
            type={TEXTFIELD_TYPE_TEXT}
            {...TEXTFIELD_PROPS_FEDERATION_DEFAULT_USER}
            value={formDataValues.username}
            initialValue={username}
            onChange={handleUsernameChange}
            disabled={!enabled}
          />
          <TextFieldWithSubmit
            fieldName="goLiveMessage"
            {...TEXTFIELD_PROPS_FEDERATION_LIVE_MESSAGE}
            type={TEXTFIELD_TYPE_TEXTAREA}
            value={formDataValues.goLiveMessage}
            initialValue={goLiveMessage}
            onChange={handleFieldChange}
            disabled={!enabled}
          />
          <ToggleSwitch
            fieldName="showEngagement"
            {...FIELD_PROPS_SHOW_FEDERATION_ENGAGEMENT}
            checked={formDataValues.showEngagement}
            disabled={!enabled}
          />
        </Col>
        <Col span={8} className="form-module">
          <EditValueArray
            title={FIELD_PROPS_FEDERATION_BLOCKED_DOMAINS.label}
            placeholder={FIELD_PROPS_FEDERATION_BLOCKED_DOMAINS.placeholder}
            description={FIELD_PROPS_FEDERATION_BLOCKED_DOMAINS.tip}
            values={formDataValues.blockedDomains}
            handleDeleteIndex={handleDeleteBlockedDomain}
            handleCreateString={handleCreateBlockedDomain}
            submitStatus={createInputStatus(blockedDomainSaveState)}
          />
        </Col>
      </Row>
      {isInfoModalOpen && (
        <FederationInfoModal
          cancelPressed={federationInfoModalCancelPressed}
          okPressed={federationInfoModalOkPressed}
        />
      )}
    </div>
  );
};

ConfigFederation.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};

export default ConfigFederation;
