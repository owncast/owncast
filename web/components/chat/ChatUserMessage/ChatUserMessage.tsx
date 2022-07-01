/* eslint-disable react/no-danger */
import { useEffect, useState } from 'react';
import { Highlight } from 'react-highlighter-ts';
import he from 'he';
import cn from 'classnames';
import { Button, Dropdown, Menu } from 'antd';
import { DeleteOutlined, EllipsisOutlined, StopOutlined } from '@ant-design/icons';
import s from './ChatUserMessage.module.scss';
import { formatTimestamp } from './messageFmt';
import { ChatMessage } from '../../../interfaces/chat-message.model';

interface Props {
  message: ChatMessage;
  showModeratorMenu: boolean;
  highlightString: string;
  sentBySelf: boolean;
  sameUserAsLast: boolean;
}

export default function ChatUserMessage({
  message,
  highlightString,
  showModeratorMenu,
  sentBySelf, // Move the border to the right and render a background
  sameUserAsLast,
}: Props) {
  const { body, user, timestamp } = message;
  const { displayName, displayColor } = user;

  const color = `var(--theme-user-colors-${displayColor})`;
  const formattedTimestamp = `Sent at ${formatTimestamp(timestamp)}`;
  const [formattedMessage, setFormattedMessage] = useState<string>(body);

  useEffect(() => {
    setFormattedMessage(he.decode(body));
  }, [message]);

  return (
    <div style={{ padding: 3.5 }}>
      <div
        className={cn(s.root, {
          [s.ownMessage]: sentBySelf,
        })}
        style={{ borderColor: color }}
        title={formattedTimestamp}
      >
        {!sameUserAsLast && (
          <div className={s.user} style={{ color }}>
            {displayName}
          </div>
        )}
        <Highlight search={highlightString}>
          <div className={s.message}>{formattedMessage}</div>
        </Highlight>
        {showModeratorMenu && <ModeratorMenu />}
        <div className={s.customBorder} style={{ color }} />
        <div className={s.background} style={{ color }} />
      </div>
    </div>
  );
}

function ModeratorMenu() {
  const menu = (
    <Menu
      items={[
        {
          label: (
            <div>
              <DeleteOutlined style={{ marginRight: 5 }} />
              Hide message
            </div>
          ),
          key: 0,
        },

        {
          label: (
            <div>
              <StopOutlined style={{ marginRight: 5 }} />
              Ban User
            </div>
          ),
          key: 1,
        },
      ]}
    />
  );

  return (
    <div className={s.modMenuWrapper}>
      <Dropdown overlay={menu} trigger={['click']}>
        <Button ghost icon={<EllipsisOutlined />} />
      </Dropdown>
    </div>
  );
}
