/* eslint-disable react/no-danger */
import Head from 'next/head';
import { FC, useEffect, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';

export const Theme: FC = () => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { appearanceVariables, customStyles } = clientConfig;

  const appearanceVars = Object.keys(appearanceVariables || {})
    .filter(variable => !!appearanceVariables[variable])
    .map(variable => `--${variable}: ${appearanceVariables[variable]}`);

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
					${appearanceVars.join(';\n')}
				}
			`,
        }}
      />
      {/*
        customStyles is the admin's CSS followed by every loaded
        plugin's manifest.styles content (concatenated server-side).
        One inline <style> block carries both.
      */}
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
