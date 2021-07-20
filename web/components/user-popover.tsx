// This displays a clickable user name (or whatever children element you provide), and displays a simple tooltip of created time. OnClick a modal with more information about the user is displayed.

import { useState } from 'react';
import { Divider, Modal, Tooltip, Typography, Row, Col } from 'antd';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import format from 'date-fns/format';
import { ReactNode } from 'react-markdown';
import BlockUserbutton from './ban-user-button';

import { User, UserConnectionInfo } from '../types/chat';
import { formatDisplayDate } from './user-table';
import { formatUAstring } from '../utils/format';

interface UserPopoverProps {
  user: User;
  connectionInfo?: UserConnectionInfo | null;
  children: ReactNode;
}

export default function UserPopover({ user, connectionInfo, children }: UserPopoverProps) {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const handleShowModal = () => {
    setIsModalVisible(true);
  };
  const handleCloseModal = () => {
    setIsModalVisible(false);
  };

  const { displayName, createdAt, previousNames, nameChangedAt, disabledAt } = user;
  const { connectedAt, messageCount, userAgent } = connectionInfo || {};

  let lastNameChangeDate = null;
  const nameList = previousNames && [...previousNames];

  if (previousNames && previousNames.length > 1 && nameChangedAt) {
    lastNameChangeDate = new Date(nameChangedAt);
    // reverse prev names for display purposes
    nameList.reverse();
  }

  const dateObject = new Date(createdAt);
  const createdAtDate = format(dateObject, 'PP pp');

  const lastNameChangeDuration = lastNameChangeDate
    ? formatDistanceToNow(lastNameChangeDate)
    : null;

  return (
    <>
      <Tooltip
        title={
          <>
            Created at: {createdAtDate}.
            <br /> Click for more info.
          </>
        }
        placement="bottomLeft"
      >
        <button
          type="button"
          aria-label="Display more details about this user"
          className="user-item-container"
          onClick={handleShowModal}
        >
          {children}
        </button>
      </Tooltip>

      <Modal
        destroyOnClose
        width={600}
        cancelText="Close"
        okButtonProps={{ style: { display: 'none' } }}
        title={`User details: ${displayName}`}
        visible={isModalVisible}
        onOk={handleCloseModal}
        onCancel={handleCloseModal}
      >
        <div className="user-details">
          <Typography.Title level={4}>{displayName}</Typography.Title>
          <p className="created-at">User created at {createdAtDate}.</p>
          <Row gutter={16}>
            {connectionInfo && (
              <Col md={lastNameChangeDate ? 12 : 24}>
                <Typography.Title level={5}>
                  This user is currently connected to Chat.
                </Typography.Title>
                <ul className="connection-info">
                  <li>
                    <strong>Active for:</strong> {formatDistanceToNow(new Date(connectedAt))}
                  </li>
                  <li>
                    <strong>Messages sent:</strong> {messageCount}
                  </li>
                  <li>
                    <strong>User Agent:</strong>
                    <br />
                    {formatUAstring(userAgent)}
                  </li>
                </ul>
              </Col>
            )}
            {lastNameChangeDate && (
              <Col md={connectionInfo ? 12 : 24}>
                <Typography.Title level={5}>This user is also seen as:</Typography.Title>
                <ul className="previous-names-list">
                  {nameList.map((name, index) => (
                    <li className={index === 0 ? 'latest' : ''}>
                      <span className="user-name-item">{name}</span>
                      {index === 0 ? ` (Changed ${lastNameChangeDuration} ago)` : ''}
                    </li>
                  ))}
                </ul>
              </Col>
            )}
          </Row>
          <Divider />
          {disabledAt ? (
            <>
              This user was banned on <code>{formatDisplayDate(disabledAt)}</code>.
              <br />
              <br />
              <BlockUserbutton
                label="Unban this user"
                user={user}
                isEnabled={false}
                onClick={handleCloseModal}
              />
            </>
          ) : (
            <BlockUserbutton
              label="Ban this user"
              user={user}
              isEnabled
              onClick={handleCloseModal}
            />
          )}
        </div>
      </Modal>
    </>
  );
}

UserPopover.defaultProps = {
  connectionInfo: null,
};
