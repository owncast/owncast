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
        <a href={versionNumber}> .</a>
      </span>
    </footer>
  );
};
