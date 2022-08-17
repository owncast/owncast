import cn from 'classnames';

import { ServerLogo } from '../../ui';
import SocialLinks from '../../ui/SocialLinks/SocialLinks';
import { SocialLink } from '../../../interfaces/social-link.model';
import s from './ContentHeader.module.scss';

interface Props {
  name: string;
  title: string;
  summary: string;
  tags: string[];
  links: SocialLink[];
  logo: string;
}
export default function ContentHeader({ name, title, summary, logo, tags, links }: Props) {
  return (
    <div className={s.root}>
      <div className={s.logoTitleSection}>
        <div className={s.logo}>
          <ServerLogo src={logo} />
        </div>
        <div className={s.titleSection}>
          <div className={cn(s.title, s.row)}>{name}</div>
          <div className={cn(s.subtitle, s.row)}>{title || summary}</div>
          <div className={cn(s.tagList, s.row)}>
            {tags.length > 0 && tags.map(tag => <span key={tag}>#{tag}&nbsp;</span>)}
          </div>
          <div className={cn(s.socialLinks, s.row)}>
            <SocialLinks links={links} />
          </div>
        </div>
      </div>
    </div>
  );
}
