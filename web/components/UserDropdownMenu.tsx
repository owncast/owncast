import { Menu, Dropdown } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import { useRecoilState } from 'recoil';
import { ChatVisibilityState, ChatState } from '../interfaces/application-state';
import { chatVisibilityAtom as chatVisibilityAtom } from './stores/ClientConfigStore';

interface Props {
  username: string;
  chatState: ChatState;
}

export default function UserDropdown(props: Props) {
  const { username, chatState } = props;

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
    <Dropdown overlay={menu} trigger={['click']}>
      <button type="button" className="ant-dropdown-link" onClick={e => e.preventDefault()}>
        {username} <DownOutlined />
      </button>
    </Dropdown>
  );
}
