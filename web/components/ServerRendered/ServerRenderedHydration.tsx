/* eslint-disable react/no-danger */
import { FC } from 'react';

export const ServerRenderedHydration: FC = () => (
  <script
    id="server-side-hydration"
    nonce="{{.Nonce}}"
    dangerouslySetInnerHTML={{
      __html: `
	window.configHydration = {{.ServerConfigJSON}};
	window.statusHydration = {{.StatusJSON}};
	`,
    }}
  />
);
