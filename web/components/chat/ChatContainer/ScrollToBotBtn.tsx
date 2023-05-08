import { Button } from 'antd';
import dynamic from 'next/dynamic';
import { FC } from 'react';
import styles from './ChatContainer.module.scss';

// Lazy loaded components

const VerticalAlignBottomOutlined = dynamic(
  () => import('@ant-design/icons/VerticalAlignBottomOutlined'),
  {
    ssr: false,
  },
);

type Props = {
  onClick: () => void;
};

export const ScrollToBotBtn: FC<Props> = ({ onClick }) => (
  <div className={styles.toBottomWrap} id="scroll-to-chat-bottom">
    <Button
      type="default"
      style={{ color: 'currentColor' }}
      icon={<VerticalAlignBottomOutlined />}
      onClick={onClick}
    >
      Go to last message
    </Button>
  </div>
);
