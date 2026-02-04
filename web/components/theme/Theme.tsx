/* eslint-disable react/no-danger */
import Head from 'next/head';
import { FC, useEffect, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';

export const Theme: FC = () => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { appearanceVariables, customStyles } = clientConfig;

  // Map Owncast appearance variables to CSS variables
  const appearanceVars = Object.keys(appearanceVariables || {})
    .filter(variable => !!appearanceVariables[variable])
    .map(variable => `--${variable}: ${appearanceVariables[variable]}`);

  // Also map to Ant Design CSS variables so Ant components pick up the theme
  const antVars: string[] = [];
  if (appearanceVariables?.['theme-color-background-main']) {
    antVars.push(`--ant-layout-body-bg: ${appearanceVariables['theme-color-background-main']}`);
    antVars.push(`--ant-color-bg-layout: ${appearanceVariables['theme-color-background-main']}`);
  }
  if (appearanceVariables?.['theme-color-action']) {
    antVars.push(`--ant-color-primary: ${appearanceVariables['theme-color-action']}`);
    antVars.push(`--ant-color-link: ${appearanceVariables['theme-color-action']}`);
  }
  if (appearanceVariables?.['theme-color-action-hover']) {
    antVars.push(`--ant-color-primary-hover: ${appearanceVariables['theme-color-action-hover']}`);
    antVars.push(`--ant-color-link-hover: ${appearanceVariables['theme-color-action-hover']}`);
  }
  if (appearanceVariables?.['theme-rounded-corners']) {
    antVars.push(`--ant-border-radius: ${appearanceVariables['theme-rounded-corners']}`);
  }

  const allVars = [...appearanceVars, ...antVars];

  const [themeColor, setThemeColor] = useState('#fff');

  useEffect(() => {
    const color = getComputedStyle(document.documentElement).getPropertyValue(
      '--theme-color-background-header',
    );
    setThemeColor(color);
  }, [appearanceVars]);

  return (
    <>
      <Head>
        <meta name="theme-color" content={themeColor} />
      </Head>
      <style
        dangerouslySetInnerHTML={{
          __html: `
				:root {
					${allVars.join(';\n')}
				}
			`,
        }}
      />
      <style
        dangerouslySetInnerHTML={{
          __html: `
				${customStyles}
			`,
        }}
      />
    </>
  );
};
