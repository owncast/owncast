import { Button, Col, Collapse, Row, Spin, Table, Tag } from 'antd';
import { FC, useEffect, useState } from 'react';
import format from 'date-fns/format';
import { ColumnsType } from 'antd/lib/table';
import dynamic from 'next/dynamic';
import { ErrorBoundary } from 'react-error-boundary';
import ChatModeration from '../../../services/moderation-service';
import styles from './ChatModerationDetailsModal.module.scss';
import { formatUAstring } from '../../../utils/format';
import { ComponentError } from '../../ui/ComponentError/ComponentError';

const { Panel } = Collapse;

// Lazy loaded components

const DeleteOutlined = dynamic(() => import('@ant-design/icons/DeleteOutlined'), {
  ssr: false,
});

export type ChatModerationDetailsModalProps = {
  userId: string;
  accessToken: string;
};

export interface UserDetails {
  user: User;
  connectedClients: Client[];
  messages: Message[];
}

export interface Client {
  messageCount: number;
  userAgent: string;
  connectedAt: Date;
  geo: string;
  id: number;
}

export interface Message {
  id: string;
  timestamp: Date;
  user: null;
  body: string;
}

export interface User {
  id: string;
  displayName: string;
  displayColor: number;
  createdAt: Date;
  previousNames: string[];
  scopes: string[];
  isBot: boolean;
  authenticated: boolean;
}

const removeMessage = async (messageId: string, accessToken: string) => {
  try {
    ChatModeration.removeMessage(messageId, accessToken);
  } catch (e) {
    console.error(e);
  }
};

const ValueRow = ({ label, value }: { label: string; value: string }) => (
  <Row justify="space-around" align="middle">
    <Col span={12}>{label}</Col>
    <Col span={12}>{value}</Col>
  </Row>
);

const ConnectedClient = ({ client }: { client: Client }) => {
  const { messageCount, connectedAt, geo } = client;
  const connectedAtDate = format(new Date(connectedAt), 'PP pp');

  return (
    <div>
      <ValueRow label="Messages Sent" value={messageCount.toString()} />
      {geo !== 'N/A' && <ValueRow label="Geo" value={geo} />}
      <ValueRow label="Connected At" value={connectedAtDate} />
    </div>
  );
};

const UserColorBlock = ({ color }) => {
  const bg = `var(--theme-color-users-${color})`;
  return (
    <div className={styles.colorBlock} style={{ backgroundColor: bg }}>
      Color {color}
    </div>
  );
};

export const ChatModerationDetailsModal: FC<ChatModerationDetailsModalProps> = ({
  userId,
  accessToken,
}) => {
  const [userDetails, setUserDetails] = useState<UserDetails | null>(null);
  const [loading, setLoading] = useState(true);

  const getDetails = async () => {
    try {
      const response = await (
        await fetch(`/api/moderation/chat/user/${userId}?accessToken=${accessToken}`)
      ).json();
      setUserDetails(response);
      setLoading(false);
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    getDetails();
  }, []);

  if (!userDetails) {
    return null;
  }

  const { user, connectedClients, messages } = userDetails;
  const { displayColor, createdAt, previousNames, scopes, isBot, authenticated } = user;

  const totalMessagesSent = connectedClients.reduce((acc, client) => acc + client.messageCount, 0);
  const createdAtDate = format(new Date(createdAt), 'PP pp');

  const chatMessageColumns: ColumnsType<Message> = [
    {
      title: 'Message',
      dataIndex: 'body',
      key: 'body',
    },
    {
      title: 'Sent At',
      dataIndex: 'timestamp',
      key: 'timestamp',
      render: timestamp => format(new Date(timestamp), 'PP pp'),
    },
    {
      title: 'Delete',
      key: 'delete',
      render: (_text, record) => (
        <Button
          type="primary"
          ghost
          icon={<DeleteOutlined />}
          onClick={() => removeMessage(record.id, accessToken)}
        />
      ),
    },
  ];
  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="ChatModerationDetailsModal"
          message={error.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <Spin spinning={loading}>
        <UserColorBlock color={displayColor} />
        {scopes?.map(scope => <Tag key={scope}>{scope}</Tag>)}
        {authenticated && <Tag>Authenticated</Tag>}
        {isBot && <Tag>Bot</Tag>}
        <ValueRow label="Messages Sent Across Clients" value={totalMessagesSent.toString()} />
        <ValueRow label="User Created" value={createdAtDate} />
        <ValueRow label="Known As" value={previousNames.join(',')} />
        <Collapse accordion>
          <Panel header="Currently Connected Clients" key="connected-clients">
            <Collapse accordion>
              {connectedClients.map(client => (
                <Panel header={formatUAstring(client.userAgent)} key={client.id}>
                  <ConnectedClient client={client} />
                </Panel>
              ))}
            </Collapse>
          </Panel>
          <Collapse accordion>
            <Panel header="Recent Chat Messages" key="chat-messages">
              <Table
                size="small"
                pagination={null}
                columns={chatMessageColumns}
                dataSource={messages}
                rowKey="id"
              />
            </Panel>
          </Collapse>
        </Collapse>
      </Spin>
    </ErrorBoundary>
  );
};
