import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { Layout, Menu } from 'antd';
import {
  SettingOutlined,
  HomeOutlined,
  LineChartOutlined,
  CloseCircleOutlined,
  PlayCircleFilled,
  StopFilled,
  MinusSquareFilled,
} from '@ant-design/icons';
import classNames from 'classnames';


import OwncastLogo from './logo';
import { BroadcastStatusContext } from '../utils/broadcast-status-context';

import adminStyles from '../../styles/styles.module.css';

export default function MainLayout(props) {
  const { children } = props;

  const context = useContext(BroadcastStatusContext);
  const { broadcastActive } = context || {};

  const router = useRouter();
  const { route } = router || {};

  const { Header, Footer, Content, Sider } = Layout;
  const { SubMenu } = Menu;

  const statusIcon = broadcastActive ?
    <PlayCircleFilled /> : <MinusSquareFilled />;
    const statusMessage = broadcastActive ?
    'Online' : 'Offline';

  const appClass = classNames({
    'owncast-layout': true,
    [adminStyles.online]: broadcastActive,
  })
  return (
    <Layout className={appClass}>
      <Sider
        width={240}
        style={{
          overflow: "auto",
          height: "100vh",
        }}
      >
        <Menu
          theme="dark"
          defaultSelectedKeys={[route.substring(1)]}
          defaultOpenKeys={["current-stream-menu", "utilities-menu"]}
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
            title="Stream Details"
          >
            <Menu.Item key="viewer-info">
              <Link href="/viewer-info">Viewers</Link>
            </Menu.Item>
            <Menu.Item key="hardware-info">
              <Link href="/hardware-info">Hardware</Link>
            </Menu.Item>
            {broadcastActive ? (
              <Menu.Item key="disconnect-stream" icon={<CloseCircleOutlined />}>
                <Link href="/disconnect-stream">Disconnect Stream...</Link>
              </Menu.Item>
            ) : null}
          </SubMenu>

          <SubMenu
            key="utilities-menu"
            icon={<SettingOutlined />}
            title="Utilities"
          >
            <Menu.Item key="update-server-config">
              <Link href="/update-server-config">Server Configuration</Link>
            </Menu.Item>
            <Menu.Item key="video-config">
              <Link href="/video-config">Video Configuration</Link>
            </Menu.Item>
            <Menu.Item key="storage">
              <Link href="/storage">Storage</Link>
            </Menu.Item>
            <Menu.Item key="update-stream-key">
              <Link href="/update-stream-key">Change Stream Key</Link>
            </Menu.Item>
          </SubMenu>
        </Menu>
      </Sider>

      <Layout>
        <Header className={adminStyles.header}>
          <div className={adminStyles.statusIndicatorContainer}>
            <span className={adminStyles.statusLabel}>{statusMessage}</span>
            <span className={adminStyles.statusIcon}>{statusIcon}</span>
          </div>
        </Header>
        <Content className={adminStyles.contentMain}>{children}</Content>

        <Footer style={{ textAlign: "center" }}>
          <a href="https://owncast.online/">About Owncast</a>
        </Footer>
      </Layout>
    </Layout>
  );
}

MainLayout.propTypes = {
  children: PropTypes.element.isRequired,
};