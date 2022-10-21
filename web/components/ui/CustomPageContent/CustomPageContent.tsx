/* eslint-disable react/no-danger */
import { FC } from 'react';
import { useRecoilValue } from 'recoil';
import Footer from '../Footer/Footer';
import styles from './CustomPageContent.module.scss';
import { isMobileAtom, clientConfigStateAtom } from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';

export type CustomPageContentProps = {
  content: string;
};

export const CustomPageContent: FC<CustomPageContentProps> = ({ content }) => {
  const isMobile = useRecoilValue<boolean | undefined>(isMobileAtom);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { version } = clientConfig;
  return (
    <>
      <div className={styles.pageContentContainer}>
        <div className={styles.customPageContent} dangerouslySetInnerHTML={{ __html: content }} />
      </div>
      {isMobile && <Footer version={version} />}
    </>
  );
};
