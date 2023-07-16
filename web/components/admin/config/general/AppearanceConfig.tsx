import React, { FC, useContext, useCallback, useEffect, useState } from 'react';

import { Button, Col, Collapse, Row, Slider, Space } from 'antd';
import Paragraph from 'antd/lib/typography/Paragraph';
import Title from 'antd/lib/typography/Title';
import { EditCustomStyles } from '../../EditCustomStyles';
import s from './appearance.module.scss';
import { postConfigUpdateToAPI, RESET_TIMEOUT } from '../../../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../../../utils/input-statuses';
import { ServerStatusContext } from '../../../../utils/server-status-context';
import { FormStatusIndicator } from '../../FormStatusIndicator';

const { Panel } = Collapse;

const ENDPOINT = '/appearance';

interface AppearanceVariable {
  value: string;
  description: string;
}

type ColorCollectionProps = {
  variables: { name; description; value }[];
  updateColor: (variable: string, color: string, description: string) => void;
};

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

const componentColorVariables = [
  { name: 'theme-color-background-main', description: 'Background' },
  { name: 'theme-color-action', description: 'Action' },
  { name: 'theme-color-action-hover', description: 'Action Hover' },
  { name: 'theme-color-components-primary-button-border', description: 'Primary Button Border' },
  { name: 'theme-color-components-primary-button-text', description: 'Primary Button Text' },
  { name: 'theme-color-components-chat-background', description: 'Chat Background' },
  { name: 'theme-color-components-chat-text', description: 'Text: Chat' },
  { name: 'theme-color-components-text-on-dark', description: 'Text: Light' },
  { name: 'theme-color-components-text-on-light', description: 'Text: Dark' },
  { name: 'theme-color-background-header', description: 'Header/Footer' },
  { name: 'theme-color-components-content-background', description: 'Page Content' },
  {
    name: 'theme-color-components-video-status-bar-background',
    description: 'Video Status Bar Background',
  },
  {
    name: 'theme-color-components-video-status-bar-foreground',
    description: 'Video Status Bar Foreground',
  },
];

const others = [{ name: 'theme-rounded-corners', description: 'Corner radius' }];

// Create an object so these vars can be indexed by name.
const allAvailableValues = [...componentColorVariables, ...chatColorVariables, ...others].reduce(
  (obj, val) => {
    // eslint-disable-next-line no-param-reassign
    obj[val.name] = { name: val.name, description: val.description };
    return obj;
  },
  {},
);

// eslint-disable-next-line react/function-component-definition
const ColorPicker = React.memo(
  ({
    value,
    name,
    description,
    onChange,
  }: {
    value: string;
    name: string;
    description: string;
    onChange: (name: string, value: string, description: string) => void;
  }) => (
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
  ),
);

const ColorCollection: FC<ColorCollectionProps> = ({ variables, updateColor }) => {
  const cc = variables.map(colorVar => {
    const { name, description, value } = colorVar;

    return (
      <ColorPicker
        key={name}
        value={value}
        name={name}
        description={description}
        onChange={updateColor}
      />
    );
  });
  // eslint-disable-next-line react/jsx-no-useless-fragment
  return <>{cc}</>;
};

// eslint-disable-next-line react/function-component-definition
export default function Appearance() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData;
  const { instanceDetails } = serverConfig;
  const { appearanceVariables } = instanceDetails;

  const [defaultValues, setDefaultValues] = useState<Record<string, AppearanceVariable>>();
  const [customValues, setCustomValues] = useState<Record<string, AppearanceVariable>>();

  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  let resetTimer = null;
  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  const setDefaults = () => {
    const c = {};
    [...componentColorVariables, ...chatColorVariables, ...others].forEach(color => {
      const resolvedColor = getComputedStyle(document.documentElement).getPropertyValue(
        `--${color.name}`,
      );
      c[color.name] = { value: resolvedColor.trim(), description: color.description };
    });
    setDefaultValues(c);
  };

  useEffect(() => {
    setDefaults();
  }, []);

  useEffect(() => {
    if (Object.keys(appearanceVariables).length === 0) return;

    const c = {};
    Object.keys(appearanceVariables).forEach(key => {
      c[key] = {
        value: appearanceVariables[key],
        description: allAvailableValues[key]?.description || '',
      };
    });
    setCustomValues(c);
  }, [appearanceVariables]);

  const updateColor = useCallback((variable: string, color: string, description: string) => {
    setCustomValues(oldCustomValues => ({
      ...oldCustomValues,
      [variable]: { value: color, description },
    }));
  }, []);

  const reset = async () => {
    await postConfigUpdateToAPI({
      apiPath: ENDPOINT,
      data: { value: {} },
      onSuccess: () => {
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Updated.'));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
        setCustomValues({});
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const save = async () => {
    const c = {};
    Object.keys(customValues).forEach(color => {
      c[color] = customValues[color].value;
    });

    await postConfigUpdateToAPI({
      apiPath: ENDPOINT,
      data: { value: c },
      onSuccess: () => {
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Updated.'));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);

        setFieldInConfigState({
          fieldName: 'appearanceVariables',
          value: c,
          path: 'instanceDetails',
        });
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const onBorderRadiusChange = (value: string) => {
    const variableName = 'theme-rounded-corners';

    updateColor(variableName, `${value.toString()}px`, '');
  };

  if (!defaultValues) {
    return <div>Loading...</div>;
  }

  const transformToColorMap = variables =>
    variables.map(colorVar => {
      const source = customValues?.[colorVar.name] ? customValues : defaultValues;
      const { name, description } = colorVar;
      const { value } = source[name];
      return { name, description, value };
    });

  return (
    <>
      <Space direction="vertical">
        <Title>Customize Appearance</Title>
        <Paragraph>The following colors are used across the user interface.</Paragraph>
        <div>
          <Collapse defaultActiveKey={['1']}>
            <Panel header={<strong>Section Colors</strong>} key="1">
              <p>
                Certain sections of the interface can be customized by selecting new colors for
                them.
              </p>
              <Row gutter={[16, 16]}>
                <ColorCollection
                  variables={transformToColorMap(componentColorVariables)}
                  updateColor={updateColor}
                />
              </Row>
            </Panel>
            <Panel header={<strong>Chat User Colors</strong>} key="2">
              <Row gutter={[16, 16]}>
                <ColorCollection
                  variables={transformToColorMap(chatColorVariables)}
                  updateColor={updateColor}
                />
              </Row>
            </Panel>
            <Panel header={<strong>Other Settings</strong>} key="4">
              How rounded should corners be?
              <Row gutter={[16, 16]}>
                <Col span={12}>
                  <Slider
                    min={0}
                    max={20}
                    tooltip={{ formatter: null }}
                    onChange={v => {
                      onBorderRadiusChange(v);
                    }}
                    value={Number(
                      customValues?.['theme-rounded-corners']?.value?.replace('px', '') ??
                        defaultValues?.['theme-rounded-corners']?.value?.replace('px', '') ??
                        0,
                    )}
                  />
                </Col>
                <Col span={4}>
                  <div
                    style={{
                      width: '100px',
                      height: '30px',
                      borderRadius: `${
                        customValues?.['theme-rounded-corners']?.value ??
                        defaultValues?.['theme-rounded-corners']?.value
                      }`,
                      backgroundColor: 'var(--theme-color-palette-7)',
                    }}
                  />
                </Col>
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
      </Space>
      <div className="form-module page-content-module">
        <EditCustomStyles />
      </div>
    </>
  );
}
