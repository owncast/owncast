import React, { useState, useEffect } from 'react';
import { Table, Tag, Space, Button, Modal, Checkbox, Input, Typography, Tooltip } from 'antd';
import { DeleteOutlined, EyeTwoTone, EyeInvisibleOutlined } from '@ant-design/icons';
const { Title, Paragraph, Text } = Typography;

import format from 'date-fns/format';

import { fetchData, ACCESS_TOKENS, DELETE_ACCESS_TOKEN, CREATE_ACCESS_TOKEN } from '../utils/apis';

const availableScopes = {
  CAN_SEND_SYSTEM_MESSAGES: {
    name: 'System messages',
    description: 'You can send official messages on behalf of the system',
    color: 'purple',
  },
  CAN_SEND_MESSAGES: {
    name: 'User chat messages',
    description: 'You can send messages on behalf of a username',
    color: 'green',
  },
  HAS_ADMIN_ACCESS: {
    name: 'Has admin access',
    description: 'Can perform administrative actions such as moderation, get server statuses, etc',
    color: 'red',
  },
};

function convertScopeStringToTag(scopeString) {
  if (!scopeString || !availableScopes[scopeString]) {
    return null;
  }

  const scope = availableScopes[scopeString];

  return (
    <Tooltip key={scopeString} title={scope.description}>
      <Tag color={scope.color}>{scope.name}</Tag>
    </Tooltip>
  );
}

function NewTokenModal(props) {
  const [selectedScopes, setSelectedScopes] = useState([]);
  const [name, setName] = useState('');

  const scopes = Object.keys(availableScopes).map(function (key) {
    return { value: key, label: availableScopes[key].description };
  });

  function onChange(checkedValues) {
    setSelectedScopes(checkedValues);
  }

  function saveToken() {
    props.onOk(name, selectedScopes);

    // Clear the modal
    setSelectedScopes([]);
    setName('');
  }

  const okButtonProps = {
    disabled: selectedScopes.length === 0 || name === '',
  };

  function selectAll() {
    setSelectedScopes(Object.keys(availableScopes));
  }

  return (
    <Modal
      title="Create New Access token"
      visible={props.visible}
      onOk={saveToken}
      onCancel={props.onCancel}
      okButtonProps={okButtonProps}
    >
      <p>
        <Input
          value={name}
          placeholder="Access token name/description"
          onChange={input => setName(input.currentTarget.value)}
        />
      </p>

      <p>
        Select the permissions this access token will have. It cannot be edited after it's created.
      </p>
      <Checkbox.Group options={scopes} value={selectedScopes} onChange={onChange} />
      <Button type="text" size="small" onClick={selectAll}>
        Select all
      </Button>
    </Modal>
  );
}

export default function AccessTokens() {
  const [tokens, setTokens] = useState([]);
  const [isTokenModalVisible, setIsTokenModalVisible] = useState(false);

  const columns = [
    {
      title: '',
      key: 'delete',
      render: (text, record) => (
        <Space size="middle">
          <Button onClick={() => handleDeleteToken(record.token)} icon={<DeleteOutlined />} />
        </Space>
      ),
    },
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Token',
      dataIndex: 'token',
      key: 'token',
      render: (text, record) => <Input.Password size="small" bordered={false} value={text} />,
    },
    {
      title: 'Scopes',
      dataIndex: 'scopes',
      key: 'scopes',
      render: scopes => (
        <>
          {scopes.map(scope => {
            return convertScopeStringToTag(scope);
          })}
        </>
      ),
    },
    {
      title: 'Last Used',
      dataIndex: 'lastUsed',
      key: 'lastUsed',
      render: lastUsed => {
        if (!lastUsed) {
          return 'Never';
        }
        const dateObject = new Date(lastUsed);
        return format(dateObject, 'P p');
      },
    },
  ];

  const getAccessTokens = async () => {
    try {
      const result = await fetchData(ACCESS_TOKENS);
      setTokens(result);
    } catch (error) {
      handleError(error);
    }
  };

  useEffect(() => {
    getAccessTokens();
  }, []);

  async function handleDeleteToken(token) {
    try {
      const result = await fetchData(DELETE_ACCESS_TOKEN, {
        method: 'POST',
        data: { token: token },
      });
      getAccessTokens();
    } catch (error) {
      handleError(error);
    }
  }

  async function handleSaveToken(name: string, scopes: string[]) {
    try {
      const newToken = await fetchData(CREATE_ACCESS_TOKEN, {
        method: 'POST',
        data: { name: name, scopes: scopes },
      });
      setTokens(tokens.concat(newToken));
    } catch (error) {
      handleError(error);
    }
  }

  function handleError(error) {
    console.error('error', error);
    alert(error);
  }

  const showCreateTokenModal = () => {
    setIsTokenModalVisible(true);
  };

  const handleTokenModalSaveButton = (name, scopes) => {
    setIsTokenModalVisible(false);
    handleSaveToken(name, scopes);
  };

  const handleTokenModalCancel = () => {
    setIsTokenModalVisible(false);
  };

  return (
    <div>
      <Title>Access Tokens</Title>
      <Paragraph>
        Access tokens are used to allow external, 3rd party tools to perform specific actions on
        your Owncast server. They should be kept secure and never included in client code, instead
        they should be kept on a server that you control.
      </Paragraph>
      <Paragraph>
        Read more about how to use these tokens, with examples, at{' '}
        <a href="https://owncast.online/docs/integrations/">our documentation</a>.
      </Paragraph>

      <Table rowKey="token" columns={columns} dataSource={tokens} pagination={false} />
      <br />
      <Button type="primary" onClick={showCreateTokenModal}>
        Create Access Token
      </Button>
      <NewTokenModal
        visible={isTokenModalVisible}
        onOk={handleTokenModalSaveButton}
        onCancel={handleTokenModalCancel}
      />
    </div>
  );
}
