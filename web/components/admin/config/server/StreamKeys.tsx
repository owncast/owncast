import React, { useContext, useEffect, useState } from 'react';
import { Table, Space, Button, Typography, Alert, Input, Form, message } from 'antd';
import dynamic from 'next/dynamic';
import { ServerStatusContext } from '../../../../utils/server-status-context';

import { fetchData, UPDATE_STREAM_KEYS } from '../../../../utils/apis';
import { PASSWORD_COMPLEXITY_RULES, REGEX_PASSWORD } from '../../../../utils/config-constants';

const { Paragraph } = Typography;

// Lazy loaded components

const DeleteOutlined = dynamic(() => import('@ant-design/icons/DeleteOutlined'), {
  ssr: false,
});

const EyeOutlined = dynamic(() => import('@ant-design/icons/EyeOutlined'), {
  ssr: false,
});

const PlusOutlined = dynamic(() => import('@ant-design/icons/PlusOutlined'), {
  ssr: false,
});

const saveKeys = async (keys, setError) => {
  try {
    await fetchData(UPDATE_STREAM_KEYS, {
      method: 'POST',
      auth: true,
      data: { value: keys },
    });
  } catch (error) {
    console.error(error);
    setError(error.message);
  }
};
export const generateRndKey = () => {
  let defaultKey = '';
  let isValidStreamKey = false;
  const streamKeyRegex = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[!@#$^&*]).{8,192}$/;
  const s = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$^&*';

  while (!isValidStreamKey) {
    const temp = Array.apply(20, Array(30))
      .map(() => s.charAt(Math.floor(Math.random() * s.length)))
      .join('');
    if (streamKeyRegex.test(temp)) {
      isValidStreamKey = true;
      defaultKey = temp;
    }
  }
  return defaultKey;
};

const AddKeyForm = ({ setShowAddKeyForm, setFieldInConfigState, streamKeys, setError }) => {
  const [hasChanged, setHasChanged] = useState(true);
  const [form] = Form.useForm();
  const { Item } = Form;

  // Password Complexity rules
  const passwordComplexityRules = [];

  useEffect(() => {
    PASSWORD_COMPLEXITY_RULES.forEach(element => {
      passwordComplexityRules.push(element);
    });
  }, []);

  const handleAddKey = (newkey: any) => {
    const updatedKeys = [...streamKeys, newkey];

    setFieldInConfigState({
      fieldName: 'streamKeys',
      value: updatedKeys,
    });

    saveKeys(updatedKeys, setError);

    setShowAddKeyForm(false);
  };

  const handleInputChange = (event: any) => {
    const val = event.target.value;
    if (REGEX_PASSWORD.test(val)) {
      setHasChanged(true);
    } else {
      setHasChanged(false);
    }
  };

  // Default auto-generated key
  const defaultKey = generateRndKey();

  return (
    <Form
      layout="horizontal"
      autoComplete="off"
      onFinish={handleAddKey}
      form={form}
      style={{ display: 'flex', flexDirection: 'row' }}
      initialValues={{ key: defaultKey, comment: 'My new key' }}
    >
      <Item
        style={{ width: '60%', marginRight: '5px' }}
        label="Key"
        name="key"
        tooltip={
          <p>
            The key you provide your broadcasting software. Please note that the key must be a
            minimum of eight characters and must include at least one uppercase letter, at least one
            lowercase letter, at least one special character, and at least one number.
          </p>
        }
        rules={PASSWORD_COMPLEXITY_RULES}
      >
        <Input placeholder="your key" onChange={handleInputChange} />
      </Item>
      <Item
        style={{ width: '60%', marginRight: '5px' }}
        label="Comment"
        name="comment"
        tooltip="For remembering why you added this key"
      >
        <Input placeholder="My OBS Key" />
      </Item>
      <Button type="primary" htmlType="submit" disabled={!hasChanged}>
        Add
      </Button>
    </Form>
  );
};

const AddKeyButton = ({ setShowAddKeyForm }) => (
  <Button type="default" onClick={() => setShowAddKeyForm(true)}>
    <PlusOutlined />
  </Button>
);
const copyText = (text: string) => {
  navigator.clipboard
    .writeText(text)
    .then(() => message.success('Copied to clipboard'))
    .catch(() => message.error('Failed to copy to clipboard'));
};

const StreamKeys = () => {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { streamKeys } = serverConfig;
  const [showAddKeyForm, setShowAddKeyForm] = useState(false);
  const [showKeyMap, setShowKeyMap] = useState({});
  const [error, setError] = useState(null);

  const handleDeleteKey = keyToRemove => {
    const newKeys = streamKeys.filter(k => k !== keyToRemove);
    setFieldInConfigState({
      fieldName: 'streamKeys',
      value: newKeys,
    });
    saveKeys(newKeys, setError);
  };

  const handleToggleShowKey = key => {
    setShowKeyMap({
      ...showKeyMap,
      [key]: !showKeyMap[key],
    });
  };

  const columns = [
    {
      title: 'Key',
      dataIndex: 'key',
      key: 'key',
      render: text => (
        <Space direction="horizontal">
          <Paragraph
            copyable={{
              text: showKeyMap[text] ? text : '**********',
              onCopy: () => copyText(text),
            }}
          >
            {showKeyMap[text] ? text : '**********'}
          </Paragraph>

          <Button
            type="link"
            style={{ top: '-7px' }}
            icon={<EyeOutlined />}
            onClick={() => handleToggleShowKey(text)}
          />
        </Space>
      ),
    },
    {
      title: 'Comment',
      dataIndex: 'comment',
      key: 'comment',
    },
    {
      title: '',
      key: 'delete',
      render: text => (
        <Button
          disabled={streamKeys.length === 1}
          onClick={() => handleDeleteKey(text)}
          icon={<DeleteOutlined />}
        />
      ),
    },
  ];

  return (
    <div>
      <Paragraph>
        A streaming key is used with your broadcasting software to authenticate itself to Owncast.
        Most people will only need one. However, if you share a server with others or you want
        different keys for different broadcasting sources you can add more here.
      </Paragraph>
      <Paragraph>
        These keys are unrelated to the admin password and will not grant you access to make changes
        to Owncast&apos;s configuration.
      </Paragraph>
      <Paragraph>
        Read more about broadcasting at{' '}
        <a
          href="https://owncast.online/docs/broadcasting/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          the documentation
        </a>
        .
      </Paragraph>

      <Space direction="vertical" style={{ width: '70%' }}>
        {error && <Alert type="error" message="Saving Keys Error" description={error} />}

        {streamKeys.length === 0 && (
          <Alert
            message="No stream keys!"
            description="You will not be able to stream until you create at least one stream key and add it to your broadcasting software."
            type="error"
          />
        )}

        <Table
          rowKey="key"
          columns={columns}
          dataSource={streamKeys}
          pagination={false}
          // eslint-disable-next-line react/no-unstable-nested-components
          footer={() =>
            showAddKeyForm ? (
              <AddKeyForm
                setShowAddKeyForm={setShowAddKeyForm}
                streamKeys={streamKeys}
                setFieldInConfigState={setFieldInConfigState}
                setError={setError}
              />
            ) : (
              <AddKeyButton setShowAddKeyForm={setShowAddKeyForm} />
            )
          }
        />
        <br />
      </Space>
    </div>
  );
};
export default StreamKeys;
