import { Layout } from 'antd';
import { UserDropdown } from '../../common';
import Logo from '../../logo';
import s from './Header.module.scss';

const { Header } = Layout;

interface Props {
  name: string;
}

export default function HeaderComponent({ name = 'Your stream title' }: Props) {
  return (
    <Header className={`${s.header}`}>
      <div className={`${s.logo}`}>
        <Logo />
        <span>{name}</span>
      </div>
      <UserDropdown />
    </Header>
  );
}
