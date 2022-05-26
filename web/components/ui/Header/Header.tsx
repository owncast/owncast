import { Layout, Tag, Tooltip } from 'antd';
import { OwncastLogo, UserDropdown } from '../../common';
import s from './Header.module.scss';

const { Header } = Layout;

interface Props {
  name: string;
  chatAvailable: boolean;
}

export default function HeaderComponent({ name = 'Your stream title', chatAvailable }: Props) {
  return (
    <Header className={`${s.header}`}>
      <div className={`${s.logo}`}>
        <OwncastLogo variant="contrast" />
        <span>{name}</span>
      </div>
      {chatAvailable && <UserDropdown />}
      {!chatAvailable && (
        <Tooltip title="Chat is available when the stream is live." placement="left">
          <Tag color="processing">Chat offline</Tag>
        </Tooltip>
      )}
    </Header>
  );
}
