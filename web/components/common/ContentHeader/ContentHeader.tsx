import cn from 'classnames';
import { FC } from 'react';
import Linkify from 'react-linkify';
import { Logo } from '../../ui/Logo/Logo';
import { SocialLinks } from '../../ui/SocialLinks/SocialLinks';
import { SocialLink } from '../../../interfaces/social-link.model';
import styles from './ContentHeader.module.scss';

export type ContentHeaderProps = {
  name: string;
  title: string;
  summary: string;
  tags: string[];
  links: SocialLink[];
  logo: string;
};

export const ContentHeader: FC<ContentHeaderProps> = ({
  name,
  title,
  summary,
  logo,
  tags,
  links,
}) => (
  <div className={styles.root}>
    <div className={styles.logoTitleSection}>
      <div className={styles.logo}>
        <Logo src={logo} />
      </div>
      <div className={styles.titleSection}>
        <div className={cn(styles.title, styles.row, 'header-title')}>{name}</div>
        <div className={cn(styles.subtitle, styles.row, 'header-subtitle')}>
          <Linkify>{title || summary}</Linkify>
        </div>
        <div className={cn(styles.tagList, styles.row)}>
          {tags.length > 0 && tags.map(tag => <span key={tag}>#{tag}&nbsp;</span>)}
        </div>
        <div className={cn(styles.socialLinks, styles.row)}>
          <SocialLinks links={links} />
        </div>
      </div>
    </div>
  </div>
);
