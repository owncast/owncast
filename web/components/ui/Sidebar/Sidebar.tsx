import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatContainer, ChatTextField } from '../../chat';
import s from './Sidebar.module.scss';

import {
  chatMessagesAtom,
  appStateAtom,
  chatDisplayNameAtom,
  chatUserIdAtom,
} from '../../stores/ClientConfigStore';
import { AppStateOptions } from '../../stores/application-state';

export default function Sidebar() {
  const messages = useRecoilValue<ChatMessage[]>(chatMessagesAtom);
  const appState = useRecoilValue<AppStateOptions>(appStateAtom);
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const chatUserId = useRecoilValue<string>(chatUserIdAtom);

  return (
    <Sider className={s.root} collapsedWidth={0} width={320}>
      <ChatContainer
        messages={messages}
        loading={appState.chatLoading}
        usernameToHighlight={chatDisplayName}
        chatUserId={chatUserId}
        isModerator={false}
      />
      <ChatTextField />
    </Sider>
  );
}
