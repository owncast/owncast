import { FC } from 'react';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './SocialLinks.module.scss';

export type SocialLinksProps = {
  links: SocialLink[];
};

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const SocialLinks: FC<SocialLinksProps> = ({ links }) => (
  <div className={styles.links}>
    {links.map(link => (
      <a
        key={link.platform}
        href={link.url}
        className={styles.link}
        target="_blank"
        rel="noreferrer"
      >
        <img src={link.icon} alt={link.platform} title={link.platform} className={styles.link} />
      </a>
    ))}
  </div>
);
export default SocialLinks;
