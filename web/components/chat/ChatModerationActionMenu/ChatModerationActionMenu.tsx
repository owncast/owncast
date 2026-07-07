import { Dropdown, MenuProps, Space, message, Modal as AntModal } from 'antd';
import { FC, useState } from 'react';
import dynamic from 'next/dynamic';
import { Modal } from '../../ui/Modal/Modal';
import { ChatModerationDetailsModal } from '../ChatModerationDetailsModal/ChatModerationDetailsModal';
import ChatModeration from '../../../services/moderation-service';

const { confirm } = AntModal;

// Lazy loaded components

const CloseCircleOutlined = dynamic(() => import('@ant-design/icons/CloseCircleOutlined'), {
  ssr: false,
});

const ExclamationCircleOutlined = dynamic(
  () => import('@ant-design/icons/ExclamationCircleOutlined'),
  {
    ssr: false,
  },
);

const EyeInvisibleOutlined = dynamic(() => import('@ant-design/icons/EyeInvisibleOutlined'), {
  ssr: false,
});

const SmallDashOutlined = dynamic(() => import('@ant-design/icons/SmallDashOutlined'), {
  ssr: false,
});

export type ChatModerationActionMenuProps = {
  accessToken: string;
  messageID: string;
  userID: string;
  userDisplayName: string;
};

export const ChatModerationActionMenu: FC<ChatModerationActionMenuProps> = ({
  messageID,
  userID,
  userDisplayName,
  accessToken,
}) => {
  const [showUserDetailsModal, setShowUserDetailsModal] = useState(false);

  const handleBanUser = async () => {
    try {
      await ChatModeration.banUser(userID, accessToken);
    } catch (e) {
      console.error(e);
      message.error(e);
    }
  };

  const handleHideMessage = async () => {
    try {
      await ChatModeration.removeMessage(messageID, accessToken);
    } catch (e) {
      console.error(e);
      message.error(e);
    }
  };

  const confirmHideMessage = async () => {
    confirm({
      icon: <ExclamationCircleOutlined />,
      content: `Are you sure you want to remove this message from ${userDisplayName}?`,
      onOk() {
        handleHideMessage();
      },
    });
  };

  const confirmBanUser = async () => {
    confirm({
      icon: <ExclamationCircleOutlined />,
      content: `Are you sure you want to ban ${userDisplayName} from chat?`,
      onOk() {
        handleBanUser();
      },
    });
  };

  const items: MenuProps['items'] = [
    {
      icon: <EyeInvisibleOutlined />,
      label: 'Hide Message',
      key: 'hide-message',
      onClick: confirmHideMessage,
    },
    {
      icon: <CloseCircleOutlined />,
      label: 'Ban User',
      key: 'ban-user',
      onClick: confirmBanUser,
    },
    {
      label: 'More Info...',
      key: 'more-info',
      onClick: () => setShowUserDetailsModal(true),
    },
  ];

  return (
    <>
      <Dropdown menu={{ items }} trigger={['click']}>
        <button
          type="button"
          aria-label="Chat moderation options"
          onClick={e => e.preventDefault()}
        >
          <Space>
            <SmallDashOutlined />
          </Space>
        </button>
      </Dropdown>
      <Modal
        title={userDisplayName}
        open={showUserDetailsModal}
        handleCancel={() => {
          setShowUserDetailsModal(false);
        }}
      >
        <ChatModerationDetailsModal userId={userID} accessToken={accessToken} />
      </Modal>
    </>
  );
};
