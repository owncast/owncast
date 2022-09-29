import { FC } from 'react';
import styles from './IconLink.module.scss';

export type IconLinkProps = {
  href: string;
  icon: string;
  alt: string;
  title: string;
};

/**
 * Component that renders an icon within a link.
 *
 * @component
 */
export const IconLink: FC<IconLinkProps> = ({ href, icon, title, alt }) => (
  <a href={href} className={styles.link} target="_blank" rel="noreferrer">
    <img src={icon} alt={alt} title={title} className={styles.link} />
  </a>
);
