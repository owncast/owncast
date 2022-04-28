import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import ChatContainer from '../../chat/ChatContainer';
import { chatMessages } from '../../stores/ClientConfigStore';

export default function Sidebar() {
  const messages = useRecoilValue<ChatMessage[]>(chatMessages);

  return (
    <Sider
      collapsed={false}
      width={300}
      style={{
        position: 'fixed',
        right: 0,
        top: 0,
        bottom: 0,
      }}
    >
      <ChatContainer messages={messages} />
    </Sider>
  );
}
