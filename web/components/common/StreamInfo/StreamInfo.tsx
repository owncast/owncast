import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../../stores/ClientConfigStore';
import { ServerLogo } from '../../ui';
import CategoryIcon from '../../ui/CategoryIcon/CategoryIcon';
import SocialLinks from '../../ui/SocialLinks/SocialLinks';
import s from './StreamInfo.module.scss';

export default function StreamInfo() {
  const { socialHandles, name, title, tags } = useRecoilValue<ClientConfig>(clientConfigStateAtom);

  return (
    <div className={s.streamInfo}>
      <div className={s.logoTitleSection}>
        <ServerLogo src="/logo" />
        <div className={s.titleSection}>
          <div className={s.title}>{name}</div>
          <div className={s.subtitle}>
            {title}
            <CategoryIcon tags={tags} />
          </div>
          <div>{tags.length > 0 && tags.map(tag => <span key={tag}>#{tag}&nbsp;</span>)}</div>
          <SocialLinks links={socialHandles} />
        </div>
      </div>
    </div>
  );
}
