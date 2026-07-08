import { FC, useEffect, useState } from 'react';
import { Meta, StoryFn } from '@storybook/nextjs';
import { useSetRecoilState } from 'recoil';
import {
  Alert,
  Avatar,
  Badge,
  Button,
  Card,
  Checkbox,
  Divider,
  Input,
  Pagination,
  Popconfirm,
  Select,
  Slider,
  Space,
  Switch,
  Table,
  Tabs,
  Tag,
  Tooltip,
  Typography,
  message,
} from 'antd';
import { Theme } from '../components/theme/Theme';
import { clientConfigStateAtom } from '../components/stores/ClientConfigStore';
import { makeEmptyClientConfig } from '../interfaces/client-config.model';

const { Title, Text } = Typography;

// Each preset is a set of the same appearance variables an admin can set
// from the Appearance config page. They flow through the two real theming
// paths at once: the Theme component turns them into --theme-* CSS
// variables, and AntdProvider maps them onto Ant Design design tokens.
const THEMES: Record<string, { label: string; variables: Record<string, string> }> = {
  default: {
    label: 'Owncast default',
    variables: {},
  },
  midnight: {
    label: 'Midnight',
    variables: {
      'theme-color-action': '#22d3ee',
      'theme-color-action-hover': '#67e8f9',
      'theme-color-background-main': '#1f2937',
      'theme-color-background-header': '#0b1220',
      'theme-color-components-text-on-light': '#e5e7eb',
      'theme-color-components-form-field-background': '#111827',
      'theme-color-components-modal-header-background': '#0b1220',
      'theme-color-components-modal-header-text': '#f9fafb',
      'theme-color-components-modal-content-background': '#1f2937',
      'theme-color-components-menu-background': '#111827',
      'theme-color-components-menu-item-text': '#e5e7eb',
      'theme-rounded-corners': '4px',
    },
  },
  forest: {
    label: 'Forest',
    variables: {
      'theme-color-action': '#15803d',
      'theme-color-action-hover': '#22c55e',
      'theme-color-background-main': '#f0fdf4',
      'theme-color-background-header': '#14532d',
      'theme-color-components-text-on-light': '#14532d',
      'theme-color-components-form-field-background': '#dcfce7',
      'theme-color-components-modal-header-background': '#14532d',
      'theme-color-components-modal-header-text': '#f0fdf4',
      'theme-color-components-modal-content-background': '#f0fdf4',
      'theme-rounded-corners': '12px',
    },
  },
  terminal: {
    label: 'Terminal',
    variables: {
      'theme-color-action': '#00ff41',
      'theme-color-action-hover': '#7dff9b',
      'theme-color-background-main': '#0a0a0a',
      'theme-color-background-header': '#000000',
      'theme-color-components-text-on-light': '#00ff41',
      'theme-color-components-form-field-background': '#001a06',
      'theme-color-components-modal-header-background': '#000000',
      'theme-color-components-modal-header-text': '#00ff41',
      'theme-color-components-modal-content-background': '#0a0a0a',
      'theme-rounded-corners': '0px',
    },
  },
};

const SWATCH_VARS = [
  'theme-color-action',
  'theme-color-action-hover',
  'theme-color-background-main',
  'theme-color-background-header',
  'theme-color-components-form-field-background',
];

const sampleColumns = [
  { title: 'Name', dataIndex: 'name', key: 'name', sorter: true },
  { title: 'Status', dataIndex: 'status', key: 'status' },
];
const sampleRows = [
  { key: '1', name: 'stream-a', status: 'live' },
  { key: '2', name: 'stream-b', status: 'offline' },
  { key: '3', name: 'stream-c', status: 'live' },
];

export type ThemePlaygroundProps = {
  theme: string;
};

const ThemePlayground: FC<ThemePlaygroundProps> = ({ theme }) => {
  const [selected, setSelected] = useState(theme in THEMES ? theme : 'default');
  const setClientConfig = useSetRecoilState(clientConfigStateAtom);

  // Keep the in-canvas dropdown in sync with the Storybook control.
  useEffect(() => {
    if (theme in THEMES) {
      setSelected(theme);
    }
  }, [theme]);

  // Push the preset through the real pipeline: the appearance variables
  // land in the client config atom, exactly as if the admin had saved
  // them, and both the Theme component and AntdProvider react to it.
  useEffect(() => {
    setClientConfig({
      ...makeEmptyClientConfig(),
      appearanceVariables: THEMES[selected].variables,
    });
  }, [selected, setClientConfig]);

  return (
    <>
      <Theme />
      <div
        id="theme-playground"
        style={{
          padding: '1.5rem',
          minHeight: '100vh',
          background: 'var(--theme-color-background-main)',
          color: 'var(--theme-color-components-text-on-light)',
        }}
      >
        <Space orientation="vertical" size="large" style={{ width: '100%' }}>
          <Space>
            <Text style={{ color: 'inherit' }}>Theme:</Text>
            <Select
              id="theme-picker"
              value={selected}
              onChange={setSelected}
              style={{ width: 220 }}
              options={Object.entries(THEMES).map(([value, t]) => ({
                value,
                label: t.label,
              }))}
            />
          </Space>

          <div>
            <Divider titlePlacement="left">CSS variables (non-antd surfaces)</Divider>
            <Space wrap>
              {SWATCH_VARS.map(v => (
                <div key={v} style={{ textAlign: 'center' }}>
                  <div
                    style={{
                      width: 88,
                      height: 44,
                      borderRadius: 'var(--theme-rounded-corners)',
                      border: '1px solid rgb(128 128 128 / 40%)',
                      background: `var(--${v})`,
                    }}
                  />
                  <Text style={{ fontSize: 10, color: 'inherit' }}>
                    {v.replace('theme-color-', '').replace('components-', '')}
                  </Text>
                </div>
              ))}
            </Space>
          </div>

          <div>
            <Divider titlePlacement="left">Buttons and controls</Divider>
            <Space wrap>
              <Button type="primary" id="primary-button">
                Primary
              </Button>
              <Button>Default</Button>
              <Button danger>Danger</Button>
              <Button type="primary" disabled>
                Disabled
              </Button>
              <Popconfirm title="Are you sure?">
                <Button>Popconfirm</Button>
              </Popconfirm>
              <Tooltip title="A themed tooltip">
                <Button>Tooltip</Button>
              </Tooltip>
              <Button
                onClick={() => message.info('Static message APIs render with default tokens.')}
              >
                message.info
              </Button>
              <Switch defaultChecked />
              <Checkbox defaultChecked style={{ color: 'inherit' }}>
                Checkbox
              </Checkbox>
            </Space>
          </div>

          <div>
            <Divider titlePlacement="left">Forms</Divider>
            <Space wrap>
              <Input placeholder="Text input" style={{ width: 180 }} />
              <Input.Search placeholder="Search" style={{ width: 200 }} />
              <Select
                defaultValue="one"
                style={{ width: 140 }}
                options={[
                  { value: 'one', label: 'Option one' },
                  { value: 'two', label: 'Option two' },
                ]}
              />
              <Slider defaultValue={40} style={{ width: 160 }} />
            </Space>
          </div>

          <div>
            <Divider titlePlacement="left">Content</Divider>
            <Space orientation="vertical" style={{ width: '100%' }}>
              <Space wrap>
                <Tag color="processing">Tag</Tag>
                <Badge count={5}>
                  <Avatar shape="square">OC</Avatar>
                </Badge>
                <Pagination total={50} pageSize={10} showSizeChanger={false} />
              </Space>
              <Alert
                type="error"
                title="An error alert, palette-colored via ant-overrides."
                showIcon
              />
              <Alert type="info" title="An info alert, colored by design tokens." showIcon />
              <Tabs
                defaultActiveKey="a"
                items={[
                  {
                    key: 'a',
                    label: 'First tab',
                    children: (
                      <Text style={{ color: 'inherit' }}>
                        Ink bar and active color follow theme-color-action.
                      </Text>
                    ),
                  },
                  {
                    key: 'b',
                    label: 'Second tab',
                    children: <Text style={{ color: 'inherit' }}>Second pane.</Text>,
                  },
                ]}
              />
              <Card title="A card" size="small" style={{ maxWidth: 420 }}>
                Card content on the themed background.
              </Card>
              <Table
                columns={sampleColumns}
                dataSource={sampleRows}
                size="small"
                pagination={false}
              />
            </Space>
          </div>

          <Title level={5} style={{ color: 'inherit' }}>
            Everything above restyles live: CSS variables through the Theme component, antd
            components through the design tokens in AntdProvider.
          </Title>
        </Space>
      </div>
    </>
  );
};

export default {
  title: 'owncast/Style Guide/Theme playground',
  component: ThemePlayground,
  parameters: {
    layout: 'fullscreen',
    docs: {
      description: {
        component:
          'Live testbed for the custom theming pipeline. Picking a theme (from the in-canvas dropdown or the Storybook control) writes appearance variables into the client config atom, exactly as the admin appearance page does. The Theme component turns them into --theme-* CSS variables and AntdProvider maps them onto Ant Design design tokens, so both layers restyle in real time. Note: the static message/notification APIs deliberately render with default tokens (they live in a detached React root).',
      },
    },
  },
  argTypes: {
    theme: {
      control: 'select',
      options: Object.keys(THEMES),
      description: 'Appearance variable preset',
    },
  },
} as Meta<typeof ThemePlayground>;

const Template: StoryFn<ThemePlaygroundProps> = args => <ThemePlayground {...args} />;

export const Playground = Template.bind({});
Playground.args = { theme: 'default' };

export const Midnight = Template.bind({});
Midnight.args = { theme: 'midnight' };

export const Forest = Template.bind({});
Forest.args = { theme: 'forest' };

export const Terminal = Template.bind({});
Terminal.args = { theme: 'terminal' };
