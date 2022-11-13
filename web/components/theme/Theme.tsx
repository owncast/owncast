/* eslint-disable react/no-danger */
import { FC } from 'react';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';

export const Theme: FC = () => {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { appearanceVariables, customStyles } = clientConfig;

  const appearanceVars = Object.keys(appearanceVariables)
    .filter(variable => !!appearanceVariables[variable])
    .map(variable => `--${variable}: ${appearanceVariables[variable]}`);

  return (
    <style
      dangerouslySetInnerHTML={{
        __html: `
				:root {
					${appearanceVars.join(';\n')}
				}
				${customStyles}
			`,
      }}
    />
  );
};
