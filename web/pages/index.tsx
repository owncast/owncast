<<<<<<< HEAD
import { ReactElement } from 'react';
import { Main } from '../components/layouts/Main';
=======
import { Main } from '../components/layouts/Main/Main';
>>>>>>> 4efa8716a (fix scrolling issues on mobile)

export default function Home() {
  return <Main />;
}
Home.getLayout = function getLayout(page: ReactElement) {
  return page;
};
