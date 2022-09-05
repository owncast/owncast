import { FC } from 'react';
import s from './Footer.module.scss';

export type FooterComponentProps = {
  version: string;
};

export const FooterComponent: FC<FooterComponentProps> = ({ version }) => (
  <div className={s.footer}>
    <div className={s.text}>
      Powered by <a href="https://owncast.online">{version}</a>
    </div>
    <div className={s.links}>
      <div className={s.item}>
        <a href="https://owncast.online/docs" target="_blank" rel="noreferrer">
          Documentation
        </a>
      </div>
      <div className={s.item}>
        <a href="https://owncast.online/help" target="_blank" rel="noreferrer">
          Contribute
        </a>
      </div>
      <div className={s.item}>
        <a href="https://github.com/owncast/owncast" target="_blank" rel="noreferrer">
          Source
        </a>
      </div>
    </div>
  </div>
);
export default FooterComponent;
