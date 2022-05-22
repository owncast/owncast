import { Menu, Dropdown, Button, Space } from 'antd';
import { CaretDownOutlined, UserOutlined } from '@ant-design/icons';
import { useRecoilState, useRecoilValue } from 'recoil';
import { useState } from 'react';
import Modal from '../../ui/Modal/Modal';
import { chatVisibilityAtom, chatDisplayNameAtom } from '../../stores/ClientConfigStore';
import { ChatState, ChatVisibilityState } from '../../../interfaces/application-state';
import s from './UserDropdown.module.scss';
import NameChangeModal from '../../modals/NameChangeModal';

interface Props {
  username?: string;
  chatState: ChatState;
}

export default function UserDropdown({ username: defaultUsername, chatState }: Props) {
  const [chatVisibility, setChatVisibility] =
    useRecoilState<ChatVisibilityState>(chatVisibilityAtom);
  const username = defaultUsername || useRecoilValue(chatDisplayNameAtom);
  const [showNameChangeModal, setShowNameChangeModal] = useState<boolean>(false);

  const toggleChatVisibility = () => {
    if (chatVisibility === ChatVisibilityState.Hidden) {
      setChatVisibility(ChatVisibilityState.Visible);
    } else {
      setChatVisibility(ChatVisibilityState.Hidden);
    }
  };

  const handleChangeName = () => {
    setShowNameChangeModal(true);
  };

  const menu = (
    <Menu>
      <Menu.Item key="0" onClick={() => handleChangeName()}>
        Change name
      </Menu.Item>
      <Menu.Item key="1">Authenticate</Menu.Item>
      {chatState === ChatState.Available && (
        <Menu.Item key="3" onClick={() => toggleChatVisibility()}>
          Toggle chat
        </Menu.Item>
      )}
    </Menu>
  );

  return (
    <div className={`${s.root}`}>
      <Dropdown overlay={menu} trigger={['click']}>
        <Button icon={<UserOutlined style={{ marginRight: '.5rem' }} />}>
          <Space>
            {username}
            <CaretDownOutlined />
          </Space>
        </Button>
      </Dropdown>
      <Modal
        title="Change Chat Display Name"
        visible={showNameChangeModal}
        handleCancel={() => setShowNameChangeModal(false)}
      >
        <NameChangeModal />
      </Modal>
    </div>
  );
}

UserDropdown.defaultProps = {
  username: undefined,
};
