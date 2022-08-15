import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatContainer } from '../../chat';
import s from './Sidebar.module.scss';

import {
  chatMessagesAtom,
  chatDisplayNameAtom,
  chatUserIdAtom,
} from '../../stores/ClientConfigStore';

export default function Sidebar() {
  const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const chatUserId = useRecoilValue<string>(chatUserIdAtom);

  return (
    <Sider className={s.root} collapsedWidth={0} width={320}>
      <ChatContainer
        messages={messages}
        usernameToHighlight={chatDisplayName}
        chatUserId={chatUserId}
        isModerator={false}
        isMobile={false}
      />
    </Sider>
  );
}
