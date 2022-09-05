import { FC } from 'react';
import styles from './Footer.module.scss';

export type FooterProps = {
  version: string;
};

export const Footer: FC<FooterProps> = ({ version }) => (
  <div className={styles.footer}>
    <div className={styles.text}>
      Powered by <a href="https://owncast.online">{version}</a>
    </div>
    <div className={styles.links}>
      <div className={styles.item}>
        <a href="https://owncast.online/docs" target="_blank" rel="noreferrer">
          Documentation
        </a>
      </div>
      <div className={styles.item}>
        <a href="https://owncast.online/help" target="_blank" rel="noreferrer">
          Contribute
        </a>
      </div>
      <div className={styles.item}>
        <a href="https://github.com/owncast/owncast" target="_blank" rel="noreferrer">
          Source
        </a>
      </div>
    </div>
  </div>
);
export default Footer;
