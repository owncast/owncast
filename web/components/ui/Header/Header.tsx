import { Tooltip, Avatar } from 'antd';
import { FC, useEffect, useState } from 'react';
import cn from 'classnames';
import dynamic from 'next/dynamic';
import Link from 'next/link';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../../types/localization';
import styles from './Header.module.scss';

// Lazy loaded components

const UserDropdown = dynamic(
  () => import('../../common/UserDropdown/UserDropdown').then(mod => mod.UserDropdown),
  {
    ssr: false,
  },
);

export type HeaderComponentProps = {
  name: string;
  chatAvailable: boolean;
  chatDisabled: boolean;
  online: boolean;
};

export const Header: FC<HeaderComponentProps> = ({ name, chatAvailable, chatDisabled, online }) => {
  const [canHideChat, setCanHideChat] = useState(false);
  const { t } = useTranslation();

  useEffect(() => {
    setCanHideChat(window.innerWidth >= 768);
  }, []);

  return (
    <header className={cn([`${styles.header}`], 'global-header')}>
      {online ? (
        <Link href="#player" className={styles.skipLink}>
          {t(Localization.Frontend.Header.skipToPlayer)}
        </Link>
      ) : (
        <Link href="#offline-message" className={styles.skipLink}>
          {t(Localization.Frontend.Header.skipToOfflineMessage)}
        </Link>
      )}
      <Link href="#skip-to-content" className={styles.skipLink}>
        {t(Localization.Frontend.Header.skipToContent)}
      </Link>
      <Link href="#footer" className={styles.skipLink}>
        {t(Localization.Frontend.Header.skipToFooter)}
      </Link>
      <div className={styles.logo}>
        <div id="header-logo" className={styles.logoImage}>
          <Avatar src="/logo" size="large" shape="circle" className={styles.avatar} />
        </div>
        <h1 className={styles.title} id="global-header-text">
          {name}
        </h1>
      </div>
      {chatAvailable && !chatDisabled && (
        <UserDropdown id="user-menu" hideTitleOnMobile showToggleChatOption={canHideChat} />
      )}
      {!chatAvailable && !chatDisabled && (
        <Tooltip
          overlayClassName={styles.toolTip}
          title={t(Localization.Frontend.Header.chatWillBeAvailable)}
          placement="left"
        >
          <span className={styles.chatOfflineText} id="owncast-chat-offline-text">
            {t(Localization.Frontend.Header.chatOffline)}
          </span>
        </Tooltip>
      )}
    </header>
  );
};
export default Header;
