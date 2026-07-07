// Modal showing the full metadata for a single user (creation date, previous
// names, connection info when available, ban status) plus ban/moderator
// actions. Rendered by UserPopover (opened by clicking a user's name) and by
// the admin users table (opened by clicking a row).

import { FC } from 'react';
import { Divider, Modal, Typography, Row, Col, Space } from 'antd';
import { format, formatDistanceToNow } from 'date-fns';
import { uniq } from 'lodash';

import { BanUserButton } from './BanUserButton';
import { ModeratorUserButton } from './ModeratorUserButton';

import { User, UserConnectionInfo } from '../../types/chat';
import { formatDisplayDate, formatUAstring } from '../../utils/format';

export type UserDetailsModalProps = {
  user: User;
  open: boolean;
  onClose: () => void;
  connectionInfo?: UserConnectionInfo | null;
};

export const UserDetailsModal: FC<UserDetailsModalProps> = ({
  user,
  open,
  onClose,
  connectionInfo,
}) => {
  const { displayName, createdAt, previousNames, nameChangedAt, disabledAt } = user;
  const { connectedAt, messageCount, userAgent } = connectionInfo || {};

  let lastNameChangeDate = null;
  const nameList = previousNames && [...previousNames];

  if (previousNames && previousNames.length > 1 && nameChangedAt) {
    lastNameChangeDate = new Date(nameChangedAt);
    // reverse prev names for display purposes
    nameList.reverse();
  }

  const createdAtDate = format(new Date(createdAt), 'PP pp');

  const lastNameChangeDuration = lastNameChangeDate
    ? formatDistanceToNow(lastNameChangeDate)
    : null;

  return (
    <Modal
      destroyOnClose
      width={600}
      cancelText="Close"
      okButtonProps={{ style: { display: 'none' } }}
      title={`User details: ${displayName}`}
      open={open}
      onOk={onClose}
      onCancel={onClose}
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
                {uniq(nameList).map((name, index) => (
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
        <Space direction="horizontal">
          {disabledAt ? (
            <>
              This user was banned on <code>{formatDisplayDate(disabledAt)}</code>.
              <br />
              <br />
              <BanUserButton
                label="Unban this user"
                user={user}
                isEnabled={false}
                onClick={onClose}
              />
            </>
          ) : (
            <BanUserButton label="Ban this user" user={user} isEnabled onClick={onClose} />
          )}
          <ModeratorUserButton user={user} onClick={onClose} />
        </Space>
      </div>
    </Modal>
  );
};

UserDetailsModal.defaultProps = {
  connectionInfo: null,
};
