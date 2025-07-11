/* eslint-disable react/no-danger */
import { FC } from 'react';
import { useTranslation } from 'next-export-i18n';

export interface TranslationProps {
  translationKey: string;
  vars?: Record<string, any>;
  className?: string;
}

export const Translation: FC<TranslationProps> = ({ translationKey, vars, className }) => {
  const { t } = useTranslation();

  const translatedText = t(translationKey, vars);

  return <span className={className} dangerouslySetInnerHTML={{ __html: translatedText }} />;
};

export default Translation;
