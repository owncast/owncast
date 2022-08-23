import {
  CloseCircleOutlined,
  ExclamationCircleOutlined,
  EyeInvisibleOutlined,
  SmallDashOutlined,
} from '@ant-design/icons';
import { Dropdown, Menu, MenuProps, Space, Modal, message } from 'antd';
import { useState } from 'react';
import ChatModerationDetailsModal from './ChatModerationDetailsModal';
import s from './ChatModerationActionMenu.module.scss';
import ChatModeration from '../../../services/moderation-service';

const { confirm } = Modal;

interface Props {
  accessToken: string;
  messageID: string;
  userID: string;
  userDisplayName: string;
}

export default function ChatModerationActionMenu(props: Props) {
  const { messageID, userID, userDisplayName, accessToken } = props;
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
              <span className={s.icon}>
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
              <span className={s.icon}>
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
        visible={showUserDetailsModal}
        okText="Ban User"
        okButtonProps={{ danger: true }}
        onOk={handleBanUser}
        onCancel={() => {
          setShowUserDetailsModal(false);
        }}
      >
        <ChatModerationDetailsModal userId={userID} accessToken={accessToken} />
      </Modal>
    </>
  );
}
