import { Dropdown, Menu, MenuProps, Space, message, Modal as AntModal } from 'antd';
import { FC, useState } from 'react';
import dynamic from 'next/dynamic';
import { Modal } from '../../ui/Modal/Modal';
import { ChatModerationDetailsModal } from '../ChatModerationDetailsModal/ChatModerationDetailsModal';
import styles from './ChatModerationActionMenu.module.scss';
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

  const onMenuClick: MenuProps['onClick'] = ({ key }) => {
    if (key === 'hide-message') {
      confirmHideMessage();
    } else if (key === 'ban-user') {
      confirmBanUser();
    } else if (key === 'more-info') {
      setShowUserDetailsModal(true);
    }
  };

  const menu = (
    <Menu
      onClick={onMenuClick}
      items={[
        {
          label: (
            <div>
              <span className={styles.icon}>
                <EyeInvisibleOutlined />
              </span>
              Hide Message
            </div>
          ),
          key: 'hide-message',
        },
        {
          label: (
            <div>
              <span className={styles.icon}>
                <CloseCircleOutlined />
              </span>
              Ban User
            </div>
          ),
          key: 'ban-user',
        },
        {
          label: <div>More Info...</div>,
          key: 'more-info',
        },
      ]}
    />
  );

  return (
    <>
      <Dropdown overlay={menu} trigger={['click']}>
        <button type="button" onClick={e => e.preventDefault()}>
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
