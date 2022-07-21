import { Button, Col, Row, Spin } from 'antd';
import { useEffect, useState } from 'react';
import ChatModeration from '../../../services/moderation-service';
import s from './ChatModerationDetailsModal.module.scss';

interface Props {
  userId: string;
  accessToken: string;
}

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

const ChatMessageRow = ({
  id,
  body,
  accessToken,
}: {
  id: string;
  body: string;
  accessToken: string;
}) => (
  <Row justify="space-around" align="middle">
    <Col span={18}>{body}</Col>
    <Col>
      <Button onClick={() => removeMessage(id, accessToken)}>X</Button>
    </Col>
  </Row>
);

const ConnectedClient = ({ client }: { client: Client }) => {
  const { messageCount, userAgent, connectedAt, geo } = client;

  return (
    <div>
      <ValueRow label="Messages Sent" value={`${messageCount}`} />
      <ValueRow label="Geo" value={geo} />
      <ValueRow label="Connected At" value={connectedAt.toString()} />
      <ValueRow label="User Agent" value={userAgent} />
    </div>
  );
};

// eslint-disable-next-line react/prop-types
const UserColorBlock = ({ color }) => {
  const bg = `var(--theme-user-colors-${color})`;
  return (
    <Row justify="space-around" align="middle">
      <Col span={12}>Color</Col>
      <Col span={12}>
        <div className={s.colorBlock} style={{ backgroundColor: bg }}>
          {color}
        </div>
      </Col>
    </Row>
  );
};

export default function ChatModerationDetailsModal(props: Props) {
  const { userId, accessToken } = props;
  const [userDetails, setUserDetails] = useState<UserDetails | null>(null);
  const [loading, setLoading] = useState(true);

  const getDetails = async () => {
    try {
      const response = await (await fetch(`/api/moderation/chat/user/${userId}`)).json();
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
  const { displayName, displayColor, createdAt, previousNames, scopes, isBot, authenticated } =
    user;

  return (
    <div className={s.modalContainer}>
      <Spin spinning={loading}>
        <h1>{displayName}</h1>
        <Row justify="space-around" align="middle">
          {scopes.map(scope => (
            <Col>{scope}</Col>
          ))}
          {authenticated && <Col>Authenticated</Col>}
          {isBot && <Col>Bot</Col>}
        </Row>

        <UserColorBlock color={displayColor} />

        <ValueRow label="User Created" value={createdAt.toString()} />
        <ValueRow label="Previous Names" value={previousNames.join(',')} />

        <hr />

        <h2>Currently Connected</h2>
        {connectedClients.length > 0 && (
          <Row gutter={[15, 15]} wrap>
            {connectedClients.map(client => (
              <Col flex="auto">
                <ConnectedClient client={client} />
              </Col>
            ))}
          </Row>
        )}

        <hr />
        {messages.length > 0 && (
          <div>
            <h1>Recent Chat Messages</h1>

            <div className={s.chatHistory}>
              {messages.map(message => (
                <ChatMessageRow
                  key={message.id}
                  id={message.id}
                  body={message.body}
                  accessToken={accessToken}
                />
              ))}
            </div>
          </div>
        )}
      </Spin>
    </div>
  );
}
