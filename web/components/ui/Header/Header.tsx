import { Tag } from 'antd';
import { FC } from 'react';
import cn from 'classnames';
import dynamic from 'next/dynamic';
import { OwncastLogo } from '../../common/OwncastLogo/OwncastLogo';
import styles from './Header.module.scss';

// Lazy loaded components

const UserDropdown = dynamic(() =>
  import('../../common/UserDropdown/UserDropdown').then(mod => mod.UserDropdown),
);

const Tooltip = dynamic(() => import('antd').then(mod => mod.Tooltip), {
  ssr: false,
});

export type HeaderComponentProps = {
  name: string;
  chatAvailable: boolean;
  chatDisabled: boolean;
};

export const Header: FC<HeaderComponentProps> = ({
  name = 'Your stream title',
  chatAvailable,
  chatDisabled,
}) => (
  <div className={cn([`${styles.header}`], 'global-header')}>
    <div className={styles.logo}>
      <div id="header-logo" className={styles.logoImage}>
        <OwncastLogo variant="contrast" />
      </div>
      <h1 className={styles.title} id="global-header-text" title={name}>
        {name}
      </h1>
    </div>
    {chatAvailable && !chatDisabled && <UserDropdown />}
    {!chatAvailable && !chatDisabled && (
      <Tooltip title="Chat is available when the stream is live." placement="left">
        <Tag style={{ cursor: 'pointer' }}>Chat offline</Tag>
      </Tooltip>
    )}
  </div>
);
export default Header;
