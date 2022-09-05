import { SocialLink } from '../../../interfaces/social-link.model';
import s from './SocialLinks.module.scss';

interface Props {
  // eslint-disable-next-line react/no-unused-prop-types
  links: SocialLink[];
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const SocialLinks = (props: Props) => {
  const { links } = props;
  return (
    <div className={s.links}>
      {links.map(link => (
        <a key={link.platform} href={link.url} className={s.link} target="_blank" rel="noreferrer">
          <img src={link.icon} alt={link.platform} title={link.platform} className={s.link} />
        </a>
      ))}
    </div>
  );
};
export default SocialLinks;
