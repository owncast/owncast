import Sider from 'antd/lib/layout/Sider';
import { useRecoilValue } from 'recoil';
import { chatCurrentlyVisible } from '../../stores/ClientConfigStore';

export default function Sidebar() {
  let chatOpen = useRecoilValue(chatCurrentlyVisible);
  return (
    <Sider
      collapsed={!chatOpen}
      width={300}
      style={{
        position: 'fixed',
        right: 0,
        top: 0,
        bottom: 0,
      }}
    />
  );
}
