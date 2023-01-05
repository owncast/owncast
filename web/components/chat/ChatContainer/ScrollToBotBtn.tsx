import { VerticalAlignBottomOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import { FC, MutableRefObject } from 'react';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import styles from './ChatContainer.module.scss';

type Props = {
  chatContainerRef: MutableRefObject<any>;
  messages: ChatMessage[];
};

export const ScrollToBotBtn: FC<Props> = ({ chatContainerRef, messages }) => (
  <div className={styles.toBottomWrap}>
    <Button
      type="default"
      style={{ color: 'currentColor' }}
      icon={<VerticalAlignBottomOutlined />}
      onClick={() =>
        chatContainerRef.current.scrollToIndex({
          index: messages.length - 1,
          behavior: 'auto',
        })
      }
    >
      Go to last message
    </Button>
  </div>
);
