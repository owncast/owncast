import { FC } from 'react';
import { useRecoilValue } from 'recoil';
import styles from './Footer.module.scss';
import { ServerStatus } from '../../../interfaces/server-status.model';
import { serverStatusState } from '../../stores/ClientConfigStore';

export const Footer: FC = () => {
  const clientStatus = useRecoilValue<ServerStatus>(serverStatusState);
  const { versionNumber } = clientStatus;
  return (
    <footer className={styles.footer} id="footer">
      <span>
        Powered by <a href="https://owncast.online">Owncast v{versionNumber}</a>
      </span>
      <span className={styles.links}>
        <a href="https://owncast.online/docs" target="_blank" rel="noreferrer">
          Documentation
        </a>
        <a href="https://owncast.online/help" target="_blank" rel="noreferrer">
          Contribute
        </a>
        <a href="https://github.com/owncast/owncast" target="_blank" rel="noreferrer">
          Source
        </a>
      </span>
    </footer>
  );
};
