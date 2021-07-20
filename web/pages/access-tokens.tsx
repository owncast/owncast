import React, { useState, useEffect } from 'react';
import {
  Table,
  Tag,
  Space,
  Button,
  Modal,
  Checkbox,
  Input,
  Typography,
  Tooltip,
  Row,
  Col,
} from 'antd';
import { DeleteOutlined } from '@ant-design/icons';

import format from 'date-fns/format';

import { fetchData, ACCESS_TOKENS, DELETE_ACCESS_TOKEN, CREATE_ACCESS_TOKEN } from '../utils/apis';

const { Title, Paragraph } = Typography;

const availableScopes = {
  CAN_SEND_SYSTEM_MESSAGES: {
    name: 'System messages',
    description: 'Can send official messages on behalf of the system.',
    color: 'purple',
  },
  CAN_SEND_MESSAGES: {
    name: 'User chat messages',
    description: 'Can send chat messages on behalf of the owner of this token.',
    color: 'green',
  },
  HAS_ADMIN_ACCESS: {
    name: 'Has admin access',
    description: 'Can perform administrative actions such as moderation, get server statuses, etc.',
    color: 'red',
  },
};

function convertScopeStringToTag(scopeString: string) {
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

interface Props {
  onCancel: () => void;
  onOk: any; // todo: make better type
  visible: boolean;
}
function NewTokenModal(props: Props) {
  const { onOk, onCancel, visible } = props;
  const [selectedScopes, setSelectedScopes] = useState([]);
  const [name, setName] = useState('');

  const scopes = Object.keys(availableScopes).map(key => ({
    value: key,
    label: availableScopes[key].description,
  }));

  function onChange(checkedValues) {
    setSelectedScopes(checkedValues);
  }

  function saveToken() {
    onOk(name, selectedScopes);

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
  const checkboxes = scopes.map(singleEvent => (
    <Col span={8} key={singleEvent.value}>
      <Checkbox value={singleEvent.value}>{singleEvent.label}</Checkbox>
    </Col>
  ));

  return (
    <Modal
      title="Create New Access token"
      visible={visible}
      onOk={saveToken}
      onCancel={onCancel}
      okButtonProps={okButtonProps}
    >
      <p>
        <p>
          The name will be displayed as the chat user when sending messages with this access token.
        </p>
        <Input
          value={name}
          placeholder="Name of bot, service, or integration"
          onChange={input => setName(input.currentTarget.value)}
        />
      </p>

      <p>
        Select the permissions this access token will have. It cannot be edited after it&apos;s
        created.
      </p>
      <Checkbox.Group style={{ width: '100%' }} value={selectedScopes} onChange={onChange}>
        <Row>{checkboxes}</Row>
      </Checkbox.Group>

      <p>
        <Button type="primary" onClick={selectAll}>
          Select all
        </Button>
      </p>
    </Modal>
  );
}

export default function AccessTokens() {
  const [tokens, setTokens] = useState([]);
  const [isTokenModalVisible, setIsTokenModalVisible] = useState(false);

  function handleError(error) {
    console.error('error', error);
  }

  async function getAccessTokens() {
    try {
      const result = await fetchData(ACCESS_TOKENS);
      setTokens(result);
    } catch (error) {
      handleError(error);
    }
  }
  useEffect(() => {
    getAccessTokens();
  }, []);

  async function handleDeleteToken(token) {
    try {
      await fetchData(DELETE_ACCESS_TOKEN, {
        method: 'POST',
        data: { token },
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
        data: { name, scopes },
      });
      setTokens(tokens.concat(newToken));
    } catch (error) {
      handleError(error);
    }
  }

  const columns = [
    {
      title: '',
      key: 'delete',
      render: (text, record) => (
        <Space size="middle">
          <Button onClick={() => handleDeleteToken(record.accessToken)} icon={<DeleteOutlined />} />
        </Space>
      ),
    },
    {
      title: 'Name',
      dataIndex: 'displayName',
      key: 'displayName',
    },
    {
      title: 'Token',
      dataIndex: 'accessToken',
      key: 'accessToken',
      render: text => <Input.Password size="small" bordered={false} value={text} />,
    },
    {
      title: 'Scopes',
      dataIndex: 'scopes',
      key: 'scopes',
      // eslint-disable-next-line react/destructuring-assignment
      render: scopes => <>{scopes.map(scope => convertScopeStringToTag(scope))}</>,
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
        <a
          href="https://owncast.online/docs/integrations/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          our documentation
        </a>
        .
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
