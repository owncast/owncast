import Image from 'next/image';
import { FC } from 'react';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './SocialLinks.module.scss';

export type SocialLinksProps = {
  links: SocialLink[];
};

export const SocialLinks: FC<SocialLinksProps> = ({ links }) => (
  <div className={styles.links}>
    {links.map(link => (
      <a
        key={link.platform}
        href={link.url}
        className={styles.link}
        target="_blank"
        // eslint-disable-next-line react/no-invalid-html-attribute
        rel="noreferrer me"
      >
        <Image
          src={link.icon || '/img/platformlogos/default.svg'}
          alt={link.platform}
          className={styles.link}
          width="30"
          height="30"
        />
      </a>
    ))}
  </div>
);
