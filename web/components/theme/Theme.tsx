/* eslint-disable react/no-danger */
import Head from 'next/head';
import { FC, useEffect, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';

export const Theme: FC = () => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { appearanceVariables, customStyles, pluginStyles } = clientConfig;

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
      {/*
        Appearance cascade, low to high priority (the later style block
        wins): plugin styles (baseline), then admin appearance
        variables, then admin custom CSS. So a plugin theme provides a
        baseline and the admin's explicit colors and CSS layer on top
        and win on overlap.

        pluginStyles is every loaded plugin's manifest.styles followed
        by its on_page_styles output, concatenated server-side with a
        per-plugin delimiter comment for devtools attribution.
      */}
      <style
        dangerouslySetInnerHTML={{
          __html: `
				${pluginStyles}
			`,
        }}
      />
      <style
        dangerouslySetInnerHTML={{
          __html: `
				:root {
					${appearanceVars.join(';\n')}
				}
			`,
        }}
      />
      {/* customStyles is the admin's own custom CSS, rendered last. */}
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
