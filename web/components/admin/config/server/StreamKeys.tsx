import React, { useContext, useState } from 'react';
import { Table, Space, Button, Typography, Alert, Input, Form } from 'antd';
import dynamic from 'next/dynamic';
import { ServerStatusContext } from '../../../../utils/server-status-context';

import { fetchData, UPDATE_STREAM_KEYS } from '../../../../utils/apis';

const { Paragraph } = Typography;
const { Item } = Form;

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
    setError(error);
  }
};

const AddKeyForm = ({ setShowAddKeyForm, setFieldInConfigState, streamKeys, setError }) => {
  const handleAddKey = (newkey: any) => {
    const updatedKeys = [...streamKeys, newkey];

    setFieldInConfigState({
      fieldName: 'streamKeys',
      value: updatedKeys,
    });

    saveKeys(updatedKeys, setError);

    setShowAddKeyForm(false);
  };

  // Default auto-generated key
  let defaultKey = '';
  for (let i = 0; i < 3; i += 1) {
    defaultKey += Math.random().toString(36).substring(2);
  }

  return (
    <Form layout="inline" autoComplete="off" onFinish={handleAddKey}>
      <Item label="Key" name="key" tooltip="The key you provide your broadcasting software">
        <Input placeholder="def456" defaultValue={defaultKey} />
      </Item>
      <Item label="Comment" name="comment" tooltip="For remembering why you added this key">
        <Input placeholder="My OBS Key" />
      </Item>

      <Button type="primary" htmlType="submit">
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
          <Paragraph copyable>{showKeyMap[text] ? text : '**********'}</Paragraph>

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
      render: text => <Button onClick={() => handleDeleteKey(text)} icon={<DeleteOutlined />} />,
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
