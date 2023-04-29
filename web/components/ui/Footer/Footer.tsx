import { FC } from 'react';
import styles from './Footer.module.scss';

export type FooterProps = {
  version: string;
};

export const Footer: FC<FooterProps> = ({ version }) => (
  <footer className={styles.footer} id="footer">
    <span>
      Powered by <a href="https://owncast.online">Owncast v{version}</a>
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
export default Footer;
