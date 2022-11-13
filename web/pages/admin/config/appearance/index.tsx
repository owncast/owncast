import React, { useContext, useEffect, useState } from 'react';

import { Button, Col, Collapse, Row, Space } from 'antd';
import Paragraph from 'antd/lib/typography/Paragraph';
import Title from 'antd/lib/typography/Title';
import { EditCustomStyles } from '../../../../components/config/EditCustomStyles';
import s from './appearance.module.scss';
import { postConfigUpdateToAPI, RESET_TIMEOUT } from '../../../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../../../utils/input-statuses';
import { ServerStatusContext } from '../../../../utils/server-status-context';
import { FormStatusIndicator } from '../../../../components/config/FormStatusIndicator';

const { Panel } = Collapse;

const ENDPOINT = '/appearance';

interface AppearanceVariable {
  value: string;
  description: string;
}

function ColorPicker({
  value,
  name,
  description,
  onChange,
}: {
  value: string;
  name: string;
  description: string;
  onChange: (name: string, value: string, description: string) => void;
}) {
  return (
    <Col span={3} key={name}>
      <input
        type="color"
        id={name}
        name={description}
        title={description}
        value={value}
        className={s.colorPicker}
        onChange={e => onChange(name, e.target.value, description)}
      />
      <div style={{ padding: '2px' }}>{description}</div>
    </Col>
  );
}
export default function Appearance() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};
  const { instanceDetails } = serverConfig;
  const { appearanceVariables } = instanceDetails;

  const chatColorVariables = [
    { name: 'theme-color-users-0', description: '' },
    { name: 'theme-color-users-1', description: '' },
    { name: 'theme-color-users-2', description: '' },
    { name: 'theme-color-users-3', description: '' },
    { name: 'theme-color-users-4', description: '' },
    { name: 'theme-color-users-5', description: '' },
    { name: 'theme-color-users-6', description: '' },
    { name: 'theme-color-users-7', description: '' },
  ];

  const paletteVariables = [
    { name: 'theme-color-palette-0', description: '' },
    { name: 'theme-color-palette-1', description: '' },
    { name: 'theme-color-palette-2', description: '' },
    { name: 'theme-color-palette-3', description: '' },
    { name: 'theme-color-palette-4', description: '' },
    { name: 'theme-color-palette-5', description: '' },
    { name: 'theme-color-palette-6', description: '' },
    { name: 'theme-color-palette-7', description: '' },
    { name: 'theme-color-palette-8', description: '' },
    { name: 'theme-color-palette-9', description: '' },
    { name: 'theme-color-palette-10', description: '' },
    { name: 'theme-color-palette-11', description: '' },
    { name: 'theme-color-palette-12', description: '' },
  ];

  const componentColorVariables = [
    { name: 'theme-color-background-main', description: 'Background' },
    { name: 'theme-color-action', description: 'Action' },
    { name: 'theme-color-action-hover', description: 'Action Hover' },
    { name: 'theme-color-components-chat-background', description: 'Chat Background' },
    { name: 'theme-color-components-chat-text', description: 'Text: Chat' },
    { name: 'theme-color-components-text-on-dark', description: 'Text: Light' },
    { name: 'theme-color-components-text-on-light', description: 'Text: Dark' },
    { name: 'theme-color-background-header', description: 'Header/Footer' },
  ];

  const [colors, setColors] = useState<Record<string, AppearanceVariable>>();
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  let resetTimer = null;
  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  useEffect(() => {
    const c = {};
    [...paletteVariables, ...componentColorVariables, ...chatColorVariables].forEach(color => {
      const resolvedColor = getComputedStyle(document.documentElement).getPropertyValue(
        `--${color.name}`,
      );
      c[color.name] = { value: resolvedColor.trim(), description: color.description };
    });
    setColors(c);
  }, []);

  useEffect(() => {
    if (!appearanceVariables || !colors) return;

    const c = colors;
    Object.keys(appearanceVariables).forEach(key => {
      c[key] = { value: appearanceVariables[key], description: colors[key]?.description || '' };
    });
    setColors(c);
  }, [appearanceVariables]);

  const updateColor = (variable: string, color: string, description: string) => {
    setColors({
      ...colors,
      [variable]: { value: color, description },
    });
  };

  const reset = async () => {
    setColors({});
    await postConfigUpdateToAPI({
      apiPath: ENDPOINT,
      data: { value: {} },
      onSuccess: () => {
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Updated.'));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const save = async () => {
    const c = {};
    Object.keys(colors).forEach(color => {
      c[color] = colors[color].value;
    });

    await postConfigUpdateToAPI({
      apiPath: ENDPOINT,
      data: { value: c },
      onSuccess: () => {
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Updated.'));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  if (!colors) {
    return <div>Loading...</div>;
  }

  return (
    <Space direction="vertical">
      <Title>Customize Appearance</Title>
      <Paragraph>
        The following colors are used across the user interface. You can change them.
      </Paragraph>
      <div>
        <Collapse defaultActiveKey={['1']}>
          <Panel header={<Title level={3}>Section Colors</Title>} key="1">
            <p>
              Certain specific sections of the interface changed by selecting new colors for them
              here.
            </p>
            <Row gutter={[16, 16]}>
              {componentColorVariables.map(colorVar => {
                const { name } = colorVar;
                const c = colors[name];
                return (
                  <ColorPicker
                    key={name}
                    value={c.value}
                    name={name}
                    description={c.description}
                    onChange={updateColor}
                  />
                );
              })}
            </Row>
          </Panel>
          <Panel header={<Title level={3}>Chat User Colors</Title>} key="2">
            <Row gutter={[16, 16]}>
              {chatColorVariables.map(colorVar => {
                const { name } = colorVar;
                const c = colors[name];
                return (
                  <ColorPicker
                    key={name}
                    value={c.value}
                    name={name}
                    description={c.description}
                    onChange={updateColor}
                  />
                );
              })}
            </Row>
          </Panel>
          <Panel header={<Title level={3}>Theme Colors</Title>} key="3">
            <Row gutter={[16, 16]}>
              {paletteVariables.map(colorVar => {
                const { name } = colorVar;
                const c = colors[name];
                return (
                  <ColorPicker
                    key={name}
                    value={c.value}
                    name={name}
                    description={c.description}
                    onChange={updateColor}
                  />
                );
              })}
            </Row>
          </Panel>
        </Collapse>
      </div>

      <Space direction="horizontal">
        <Button type="primary" onClick={save}>
          Save Colors
        </Button>
        <Button type="ghost" onClick={reset}>
          Reset to Defaults
        </Button>
      </Space>
      <FormStatusIndicator status={submitStatus} />
      <div className="form-module page-content-module">
        <EditCustomStyles />
      </div>
    </Space>
  );
}
