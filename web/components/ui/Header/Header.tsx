import { Layout, Tag, Tooltip } from 'antd';
import { FC } from 'react';
import { UserDropdown } from '../../common/UserDropdown/UserDropdown';
import { OwncastLogo } from '../../common/OwncastLogo/OwncastLogo';
import styles from './Header.module.scss';

const { Header: AntHeader } = Layout;

export type HeaderComponentProps = {
  name: string;
  chatAvailable: boolean;
};

export const Header: FC<HeaderComponentProps> = ({ name = 'Your stream title', chatAvailable }) => (
  <AntHeader className={`${styles.header}`}>
    <div className={`${styles.logo}`}>
      <OwncastLogo variant="contrast" />
      <span>{name}</span>
    </div>
    {chatAvailable && <UserDropdown />}
    {!chatAvailable && (
      <Tooltip title="Chat is available when the stream is live." placement="left">
        <Tag style={{ cursor: 'pointer' }}>Chat offline</Tag>
      </Tooltip>
    )}
  </AntHeader>
);
export default Header;
