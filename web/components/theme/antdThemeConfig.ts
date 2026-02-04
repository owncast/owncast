import type { ThemeConfig } from 'antd';

/**
 * Ant Design v5 theme configuration for Owncast.
 *
 * IMPORTANT: Ant Design v5 tokens do NOT support CSS variable references directly.
 * These are the DEFAULT values (matching Style Dictionary defaults).
 *
 * Runtime customization happens via the CSS variable bridge (antd-css-var-bridge.css)
 * which maps Owncast CSS variables (set by Theme.tsx) to Ant Design's CSS variables.
 *
 * With cssVar: true, Ant Design outputs its tokens as CSS variables (e.g., --ant-color-primary)
 * which we can then override using our own variables through CSS cascade.
 */
export const antdTheme: ThemeConfig = {
  token: {
    colorPrimary: '#6544e9',
    colorLink: '#6544e9',
    colorLinkHover: '#7a5cf3',
    colorLinkActive: '#7a5cf3',
    borderRadius: 9,
    colorBgContainer: '#ffffff',
    colorBgLayout: '#e2e8f0',
    colorBgElevated: '#ffffff',
    colorText: '#12161d',
    colorTextSecondary: '#5d5f72',
    colorError: '#ff4b39',
    colorWarning: '#ffc655',
    fontFamily:
      "Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif",
  },
  components: {
    Button: {
      borderRadius: 9,
      controlHeight: 36,
      fontSize: 14,
      fontSizeLG: 16,
      paddingContentHorizontal: 16,
    },
    Modal: {
      headerBg: '#2d3748',
      titleColor: '#e2e8f0',
      contentBg: '#e2e8f0',
      wireframe: true,
    },
    Input: {
      colorBgContainer: '#ffffff',
    },
    Dropdown: {
      borderRadiusLG: 9,
    },
    Menu: {
      itemBg: 'transparent',
      itemHoverBg: 'rgba(0, 0, 0, 0.05)',
      itemSelectedBg: 'rgba(0, 0, 0, 0.1)',
    },
    Tabs: {
      itemActiveColor: '#6544e9',
      itemHoverColor: '#7a5cf3',
      itemSelectedColor: '#6544e9',
      inkBarColor: '#6544e9',
    },
    Collapse: {
      headerBg: '#f0f3f8',
    },
    Table: {
      headerBg: '#2d3748',
      headerColor: '#e2e8f0',
    },
  },
};
