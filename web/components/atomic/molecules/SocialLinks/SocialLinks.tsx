import { FC } from 'react';
import { SocialLink } from '../../../../interfaces/social-link.model';
import styles from './SocialLinks.module.scss';
import { IconLink } from '../../atoms/IconLink/IconLink';

export type SocialLinksProps = {
  links: SocialLink[];
};

export const SocialLinks: FC<SocialLinksProps> = ({ links }) => (
  <div className={styles.links}>
    {links.map(link => (
      <IconLink
        key={link.platform}
        href={link.url}
        icon={link.icon}
        alt={link.platform}
        title={link.platform}
      />
    ))}
  </div>
);
