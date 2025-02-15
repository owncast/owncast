import { FC } from 'react';
import { useRecoilValue } from 'recoil';
import { useTranslation } from 'next-export-i18n';
import styles from './Footer.module.scss';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { serverStatusState } from '../../stores/ClientConfigStore';

export const Footer: FC = () => {
  const clientStatus = useRecoilValue<ServerStatus>(serverStatusState);
  const { versionNumber } = clientStatus;
  const { t } = useTranslation();
  return (
    <footer className={styles.footer} id="footer">
      <span>
        {t('Powered by Owncast')}
        <a href="https://owncast.online">v{versionNumber}</a>
      </span>
      <span className={styles.links}>
        <a href="https://owncast.online/docs" target="_blank" rel="noreferrer">
          {t('Documentation')}
        </a>
        <a href="https://owncast.online/help" target="_blank" rel="noreferrer">
          {t('Contribute')}
        </a>
        <a href="https://github.com/owncast/owncast" target="_blank" rel="noreferrer">
          {t('Source')}
        </a>
      </span>
    </footer>
  );
};
