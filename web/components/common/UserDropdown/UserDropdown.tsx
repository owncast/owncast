import { Menu, Dropdown, Button, Space } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import { useRecoilState } from 'recoil';
import { chatVisibilityAtom } from '../../stores/ClientConfigStore';
import { ChatState, ChatVisibilityState } from '../../../interfaces/application-state';
import s from './UserDropdown.module.scss';

interface Props {
  username?: string;
  chatState?: ChatState;
}

export default function UserDropdown({ username = 'test-user', chatState }: Props) {
  const chatEnabled = chatState !== ChatState.NotAvailable;
  const [chatVisibility, setChatVisibility] =
    useRecoilState<ChatVisibilityState>(chatVisibilityAtom);

  const toggleChatVisibility = () => {
    if (chatVisibility === ChatVisibilityState.Hidden) {
      setChatVisibility(ChatVisibilityState.Visible);
    } else {
      setChatVisibility(ChatVisibilityState.Hidden);
    }
  };

  const menu = (
    <Menu>
      <Menu.Item key="0">Change name</Menu.Item>
      <Menu.Item key="1">Authenticate</Menu.Item>
      {chatEnabled && (
        <Menu.Item key="3" onClick={() => toggleChatVisibility()}>
          Toggle chat
        </Menu.Item>
      )}
    </Menu>
  );

  return (
    <div className={`${s.root}`}>
      <Dropdown overlay={menu} trigger={['click']}>
        <Button>
          <Space>
            {username}
            <DownOutlined />
          </Space>
        </Button>
        {/*
      <button type="button" className="ant-dropdown-link" onClick={e => e.preventDefault()}>
        {username} <DownOutlined />
      </button>
      */}
      </Dropdown>
    </div>
  );
}
