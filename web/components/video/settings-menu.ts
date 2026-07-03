/* eslint-disable max-classes-per-file */
import { getLocalStorage } from '../../utils/localStorage';

export const LATENCY_COMPENSATION_ENABLED = 'latencyCompensatorEnabled';

const toggleBadgeStyle = (enabled: boolean): string =>
  `display: inline-block; min-width: 34px; text-align: center; padding: 1px 8px; ` +
  `border-radius: 9px; font-size: 0.75em; font-weight: 600; letter-spacing: 0.5px; ${
    enabled
      ? 'background: #12b76a; color: #fff; border: 1px solid #12b76a;'
      : 'background: transparent; color: #999; border: 1px solid #777;'
  }`;

export function createVideoSettingsMenuButton(
  player,
  videojs,
  qualities,
  latencyItemPressed: () => boolean,
): unknown {
  const VjsMenuItem = videojs.getComponent('MenuItem');
  const MenuItem = videojs.getComponent('MenuItem');
  const MenuButtonClass = videojs.getComponent('MenuButton');

  class MenuSeparator extends VjsMenuItem {
    // eslint-disable-next-line no-useless-constructor
    constructor(p: unknown, options: { selectable: boolean }) {
      super(p, options);
    }

    createEl(tag = 'button', props = {}, attributes = {}) {
      const el = super.createEl(tag, props, attributes);
      el.innerHTML = '<hr style="opacity: 0.3; margin-left: 10px; margin-right: 10px;" />';
      return el;
    }
  }

  // The "Minimize latency" toggle renders its own on/off pill so the state
  // is obvious at a glance, instead of relying on the subtle vjs-selected
  // styling.
  class LowLatencyMenuItem extends VjsMenuItem {
    // eslint-disable-next-line no-useless-constructor
    constructor(p: unknown, options: { selectable: boolean; label: string }) {
      super(p, options);
    }

    createEl(tag = 'button', props = {}, attributes = {}) {
      const el = super.createEl(tag, props, attributes);
      el.setAttribute(
        'title',
        'Experimental: slightly speeds up playback to keep you closer to live',
      );
      el.innerHTML =
        '<span style="display: flex; align-items: center; justify-content: space-between; gap: 12px; margin: 0 10px;">' +
        '<span>Minimize latency</span>' +
        `<span class="latency-toggle-state" style="${toggleBadgeStyle(false)}">OFF</span>` +
        '</span>';
      return el;
    }

    setToggleState(enabled: boolean) {
      this.selected(enabled);
      const badge = this.el().querySelector('.latency-toggle-state') as HTMLElement | null;
      if (badge) {
        badge.textContent = enabled ? 'ON' : 'OFF';
        badge.setAttribute('style', toggleBadgeStyle(enabled));
      }
    }
  }

  const lowLatencyItem = new LowLatencyMenuItem(player, {
    selectable: true,
    label: 'Minimize latency',
  });
  // Reflect the saved preference: the player auto-starts the compensator
  // from local storage, and the menu item should agree with it on load.
  lowLatencyItem.setToggleState(getLocalStorage(LATENCY_COMPENSATION_ENABLED) === 'true');
  lowLatencyItem.on('click', () => {
    const enabled: boolean = latencyItemPressed();
    lowLatencyItem.setToggleState(enabled);
  });

  const separator = new MenuSeparator(player, {
    selectable: false,
  });

  class MenuButton extends MenuButtonClass {
    constructor() {
      super(player);
    }

    // eslint-disable-next-line class-methods-use-this
    createItems() {
      const tech = player.tech({ IWillNotUseThisInPlugins: true });

      const defaultAutoItem = new MenuItem(player, {
        selectable: true,
        selected: true,
        label: 'Auto',
      });

      const items = Array(qualities.length);
      qualities.forEach(item => {
        items[item.index] = new MenuItem(player, {
          selectable: true,
          label: item.name,
        });
      });

      let clickEvent;
      if ('ontouchstart' in window) {
        clickEvent = 'touchend'; // Use touchend event for touch devices
      } else {
        clickEvent = 'click'; // Use click for all other devices
      }

      for (let i = 0; i < items.length; i += 1) {
        const item = items[i];
        // Quality selected
        item.on(clickEvent, () => {
          // If for some reason tech doesn't exist, then don't do anything
          if (!tech) {
            console.warn('Invalid attempt to access null player tech');
            return;
          }
          // Only enable and highlight this single, selected representation.
          tech.vhs.representations().forEach((rep, index) => {
            const isCurrent: boolean = index === i;
            rep.enabled(isCurrent);
            items[index].selected(isCurrent);
          });
          defaultAutoItem.selected(false);
        });
      }

      defaultAutoItem.on(clickEvent, () => {
        // Re-enable all representations.
        tech.vhs.representations().forEach(rep => {
          rep.enabled(true);
        });
        // Only highlight "Auto"
        items.forEach(item => item.selected(false));
        defaultAutoItem.selected(true);
      });

      const supportsLatencyCompensator = !!tech && !!tech.vhs;

      // Only show the quality selector if there is more than one option.
      if (qualities.length < 2 && supportsLatencyCompensator) {
        return [lowLatencyItem];
      }

      if (qualities.length > 1 && supportsLatencyCompensator) {
        return [defaultAutoItem, ...items, separator, lowLatencyItem];
      }
      if (!supportsLatencyCompensator && qualities.length === 1) {
        return [];
      }

      return [defaultAutoItem, ...items];
    }
  }

  const menuButton = new MenuButton();
  menuButton.el().setAttribute('aria-label', 'Settings');

  menuButton.addClass('vjs-quality-selector');
  videojs.registerComponent('MenuButton', MenuButton);

  return menuButton;
}
