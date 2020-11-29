import React, { useContext, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import Link from 'next/link';
import Head from 'next/head'
import { differenceInSeconds } from "date-fns";
import { useRouter } from 'next/router';
import { Layout, Menu, Popover } from 'antd';

import {
  SettingOutlined,
  HomeOutlined,
  LineChartOutlined,
  ToolOutlined,
  PlayCircleFilled,
  MinusSquareFilled,
} from '@ant-design/icons';
import classNames from 'classnames';
import { upgradeVersionAvailable } from "../../utils/apis";
import { parseSecondsToDurationString } from '../../utils/format'

import OwncastLogo from './logo';
import { ServerStatusContext } from '../../utils/server-status-context';

import adminStyles from '../../styles/styles.module.css';

let performedUpgradeCheck = false;

export default function MainLayout(props) {
  const { children } = props;

  const context = useContext(ServerStatusContext);
  const { online, broadcaster, versionNumber } = context || {};

  const router = useRouter();
  const { route } = router || {};

  const { Header, Footer, Content, Sider } = Layout;
  const { SubMenu } = Menu;

  const streamDurationString = online ? parseSecondsToDurationString(differenceInSeconds(new Date(), new Date(broadcaster.time))) : "";

  const content = (
    <div>
     <img src="/thumbnail.jpg" width="200px" />
    </div>
  );
  const statusIcon = online ? <PlayCircleFilled /> : <MinusSquareFilled />;
  const statusMessage = online ? `Online ${streamDurationString}` : "Offline";

  const [upgradeVersion, setUpgradeVersion] = useState(null);
  const checkForUpgrade = async () => {
    try {
      const result = await upgradeVersionAvailable(versionNumber);
      setUpgradeVersion(result);
    } catch (error) {
      console.log("==== error", error);
    }
  };

  useEffect(() => {
    if (!performedUpgradeCheck && !context.disableUpgradeChecks) {
      checkForUpgrade();
      performedUpgradeCheck = true
    }
  });

  const appClass = classNames({
    "owncast-layout": true,
    [adminStyles.online]: online,
  });

  const upgradeMenuItemStyle = upgradeVersion ? 'block' : 'none';
  const upgradeVersionString = upgradeVersion || '';

  return (
    <Layout className={appClass}>
      <Head>
        <title>Owncast Admin</title>
        <link rel="icon" type="image/png" sizes="32x32" href="/img/favicon/favicon-32x32.png"/>
      </Head>
      
      <Sider
        width={240}
        className={adminStyles.sideNav}
      >
        <Menu
          theme="dark"
          defaultSelectedKeys={[route.substring(1) || "home"]}
          defaultOpenKeys={["current-stream-menu", "utilities-menu", "configuration"]}
          mode="inline"
        >
          <h1 className={adminStyles.owncastTitleContainer}>
            <span className={adminStyles.logoContainer}>
              <OwncastLogo />
            </span>
            <span className={adminStyles.owncastTitle}>Owncast Admin</span>
          </h1>
          <Menu.Item key="home" icon={<HomeOutlined />}>
            <Link href="/">Home</Link>
          </Menu.Item>

          <SubMenu
            key="current-stream-menu"
            icon={<LineChartOutlined />}
            title="Current stream"
          >
            <Menu.Item key="viewer-info">
              <Link href="/viewer-info">Viewers</Link>
            </Menu.Item>
            {/* {online ? (
              <Menu.Item key="disconnect-stream" icon={<CloseCircleOutlined />}>
                <Link href="/disconnect-stream">Disconnect Stream...</Link>
              </Menu.Item>
            ) : null} */}
          </SubMenu>

          <SubMenu
            key="configuration"
            title="Configuration"
            icon={<SettingOutlined />}
          >
            <Menu.Item key="update-server-config">
              <Link href="/update-server-config">Server</Link>
            </Menu.Item>
            <Menu.Item key="video-config">
              <Link href="/video-config">Video</Link>
            </Menu.Item>
            <Menu.Item key="storage">
              <Link href="/storage">Storage</Link>
            </Menu.Item>
          </SubMenu>

          <SubMenu
            key="utilities-menu"
            icon={<ToolOutlined />}
            title="Utilities"
          >
            <Menu.Item key="hardware-info">
              <Link href="/hardware-info">Hardware</Link>
            </Menu.Item>
            <Menu.Item key="logs">
              <Link href="/logs">Logs</Link>
            </Menu.Item>
            <Menu.Item key="upgrade" style={{ display: upgradeMenuItemStyle }}>
              <Link href="/upgrade">
                <a>Upgrade to v{upgradeVersionString}</a>
              </Link>
            </Menu.Item>
          </SubMenu>
        </Menu>
      </Sider>

      <Layout className={adminStyles.layoutMain}>
        <Header className={adminStyles.header}>
        <Popover content={content} title="Thumbnail" trigger="hover">
          <div className={adminStyles.statusIndicatorContainer}>
            <span className={adminStyles.statusLabel}>{statusMessage}</span>
            <span className={adminStyles.statusIcon}>{statusIcon}</span>
          </div>
        </Popover>
        </Header>
        <Content className={adminStyles.contentMain}>{children}</Content>

        <Footer style={{ textAlign: "center" }}>
          <a href="https://owncast.online/">About Owncast v{versionNumber}</a>
        </Footer>
      </Layout>
    </Layout>
  );
}

MainLayout.propTypes = {
  children: PropTypes.element.isRequired,
};