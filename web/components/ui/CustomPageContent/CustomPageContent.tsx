/* eslint-disable react/no-danger */
import { FC } from 'react';
import styles from './CustomPageContent.module.scss';

export type CustomPageContentProps = {
  content: string;
};

export const CustomPageContent: FC<CustomPageContentProps> = ({ content }) => (
  <div id="custom-page-content">
    <div className={styles.customPageContent} dangerouslySetInnerHTML={{ __html: content }} />
  </div>
);
