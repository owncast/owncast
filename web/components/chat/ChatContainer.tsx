import { Spin } from 'antd';
import { ChatMessage } from '../../interfaces/chat-message.model';
import { ChatState } from '../../interfaces/application-state';

interface Props {
  messages: ChatMessage[];
  state: ChatState;
}

export default function ChatContainer(props: Props) {
  const { messages, state } = props;
  const loading = state === ChatState.Loading;

  return (
    <div>
      <Spin tip="Loading..." spinning={loading}>
        Chat container with scrolling chat messages go here
      </Spin>
    </div>
  );
}
