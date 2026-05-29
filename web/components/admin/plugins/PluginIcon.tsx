import React from 'react';
import { Plugin } from '../../../interfaces/plugin';
import s from './PluginIcon.module.scss';

export type PluginIconProps = {
  plugin: Pick<Plugin, 'slug' | 'hasIcon'>;
  // size controls which class set the wrapper uses; the default sidebar
  // entries want a small inline glyph, the plugins-list table cells want
  // a larger square.
  size?: 'list' | 'sidebar';
};

// PuzzlePiece is the fallback when a plugin doesn't bundle its own
// icon.png. The single path describes a classic puzzle piece (knob on
// the top, socket on the right) and uses currentColor so it inherits
// the surrounding text color, including the AntD menu's hover and
// selected states. Exported so the registry browse modal can render
// the same placeholder for plugins whose icon URL isn't set.
export const PuzzlePiece = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    className={className}
    role="presentation"
    fill="none"
    stroke="currentColor"
    strokeWidth={1.6}
    strokeLinejoin="round"
    strokeLinecap="round"
  >
    <path d="M10 4V3a2 2 0 1 1 4 0v1h4a2 2 0 0 1 2 2v4h-1a2 2 0 1 0 0 4h1v4a2 2 0 0 1-2 2h-4v-1a2 2 0 1 0-4 0v1H6a2 2 0 0 1-2-2v-4h1a2 2 0 1 0 0-4H4V6a2 2 0 0 1 2-2h4z" />
  </svg>
);

// PluginIcon renders the plugin's icon.png when one was bundled, or a
// puzzle-piece placeholder otherwise. Used by both the admin plugin
// list (large size) and the sidebar (small size). `size` selects the
// wrapper class; layout lives in PluginIcon.module.scss so the same
// img/svg works in both contexts.
export const PluginIcon = ({ plugin, size = 'list' }: PluginIconProps) => {
  const className = size === 'sidebar' ? s.sidebar : s.list;
  if (plugin.hasIcon) {
    return (
      <img
        src={`/api/plugins/${encodeURIComponent(plugin.slug)}/icon`}
        alt=""
        className={className}
      />
    );
  }
  return <PuzzlePiece className={className} />;
};
