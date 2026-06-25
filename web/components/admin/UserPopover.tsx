// This displays a clickable user name (or whatever children element you provide), and displays a simple tooltip of created time. OnClick a modal with more information about the user is displayed.

import { useState, ReactNode, FC } from 'react';
import { Tooltip } from 'antd';
import { format } from 'date-fns';

import { UserDetailsModal } from './UserDetailsModal';

import { User, UserConnectionInfo } from '../../types/chat';

export type UserPopoverProps = {
  user: User;
  connectionInfo?: UserConnectionInfo | null;
  children: ReactNode;
};

export const UserPopover: FC<UserPopoverProps> = ({ user, connectionInfo, children }) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const handleShowModal = () => {
    setIsModalOpen(true);
  };
  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  const createdAtDate = format(new Date(user.createdAt), 'PP pp');

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

      <UserDetailsModal
        user={user}
        connectionInfo={connectionInfo}
        open={isModalOpen}
        onClose={handleCloseModal}
      />
    </>
  );
};

UserPopover.defaultProps = {
  connectionInfo: null,
};
