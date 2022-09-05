import { Layout, Tag, Tooltip } from 'antd';
import { FC } from 'react';
import { UserDropdown } from '../../common/UserDropdown/UserDropdown';
import { OwncastLogo } from '../../common/OwncastLogo/OwncastLogo';
import s from './Header.module.scss';

const { Header } = Layout;

export type HeaderComponentProps = {
  name: string;
  chatAvailable: boolean;
};

export const HeaderComponent: FC<HeaderComponentProps> = ({
  name = 'Your stream title',
  chatAvailable,
}) => (
  <Header className={`${s.header}`}>
    <div className={`${s.logo}`}>
      <OwncastLogo variant="contrast" />
      <span>{name}</span>
    </div>
    {chatAvailable && <UserDropdown />}
    {!chatAvailable && (
      <Tooltip title="Chat is available when the stream is live." placement="left">
        <Tag color="processing" style={{ cursor: 'pointer' }}>
          Chat offline
        </Tag>
      </Tooltip>
    )}
  </Header>
);
export default HeaderComponent;
