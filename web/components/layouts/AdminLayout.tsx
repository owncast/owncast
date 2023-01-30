/* eslint-disable @next/next/no-css-tags */
import { FC, ReactElement } from 'react';
import ServerStatusProvider from '../../utils/server-status-context';
import AlertMessageProvider from '../../utils/alert-message-context';
import { MainLayout } from '../admin/MainLayout';

/*
NOTE: A bunch of compiled css is loaded here for the Admin UI.
These are old stylesheets that were converted from sass and should not be
edited or maintained. Instead we are using css modules everywhere. So if you
need to change a style rewrite the css file as a css module and import it
into the component that needs it, removing it from this global list.
*/

type AdminLayoutProps = {
  page: ReactElement;
};

export const AdminLayout: FC<AdminLayoutProps> = ({ page }) => (
  <>
    <link rel="stylesheet" href="/styles/admin/main-layout.css" />
    <link rel="stylesheet" href="/styles/admin/form-textfields.css" />
    <link rel="stylesheet" href="/styles/admin/config-socialhandles.css" />
    <link rel="stylesheet" href="/styles/admin/config-storage.css" />
    <link rel="stylesheet" href="/styles/admin/config-edit-string-tags.css" />
    <link rel="stylesheet" href="/styles/admin/config-video-variants.css" />
    <link rel="stylesheet" href="/styles/admin/config-public-details.css" />
    <link rel="stylesheet" href="/styles/admin/home.css" />
    <link rel="stylesheet" href="/styles/admin/chat.css" />

    <ServerStatusProvider>
      <AlertMessageProvider>
        <MainLayout>{page}</MainLayout>
      </AlertMessageProvider>
    </ServerStatusProvider>
  </>
);
