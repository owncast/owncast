import { ReactElement } from 'react';
import { Main } from '../components/layouts/Main/Main';

export default function Home() {
  return <Main />;
}
Home.getLayout = function getLayout(page: ReactElement) {
  return page;
};
