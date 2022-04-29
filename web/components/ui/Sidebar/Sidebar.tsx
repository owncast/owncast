import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import ChatContainer from '../../chat/ChatContainer';
import { chatMessages, chatVisibility as chatVisibilityAtom } from '../../stores/ClientConfigStore';
import { ChatVisibilityState } from '../../../interfaces/application-state';

export default function Sidebar() {
  const messages = useRecoilValue<ChatMessage[]>(chatMessages);
  const chatVisibility = useRecoilValue<ChatVisibilityState>(chatVisibilityAtom);

  return (
    <Sider
      collapsed={chatVisibility === ChatVisibilityState.Hidden}
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
