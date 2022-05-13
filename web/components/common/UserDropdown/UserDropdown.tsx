import { Menu, Dropdown, Button, Space } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import { useRecoilState, useRecoilValue } from 'recoil';
import { chatVisibilityAtom, chatDisplayNameAtom } from '../../stores/ClientConfigStore';
import { ChatState, ChatVisibilityState } from '../../../interfaces/application-state';
import s from './UserDropdown.module.scss';

interface Props {
  username?: string;
  chatState: ChatState;
}

export default function UserDropdown({ username: defaultUsername, chatState }: Props) {
  const [chatVisibility, setChatVisibility] =
    useRecoilState<ChatVisibilityState>(chatVisibilityAtom);
  const username = defaultUsername || useRecoilValue(chatDisplayNameAtom);

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
        <Button>
          <Space>
            {username}
            <DownOutlined />
          </Space>
        </Button>
      </Dropdown>
    </div>
  );
}

UserDropdown.defaultProps = {
  username: undefined,
};
