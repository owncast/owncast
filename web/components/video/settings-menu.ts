/* eslint-disable max-classes-per-file */
export function createVideoSettingsMenuButton(player, videojs, qualities, latencyItemPressed): any {
  const VjsMenuItem = videojs.getComponent('MenuItem');
  const MenuItem = videojs.getComponent('MenuItem');
  const MenuButtonClass = videojs.getComponent('MenuButton');

  class MenuSeparator extends VjsMenuItem {
    // eslint-disable-next-line no-useless-constructor
    constructor(p: any, options: { selectable: boolean }) {
      super(p, options);
    }

    createEl(tag = 'button', props = {}, attributes = {}) {
      const el = super.createEl(tag, props, attributes);
      el.innerHTML = '<hr style="opacity: 0.3; margin-left: 10px; margin-right: 10px;" />';
      return el;
    }
  }

  const lowLatencyItem = new MenuItem(player, {
    selectable: true,
  });
  lowLatencyItem.setAttribute('class', 'latency-toggle-item');
  lowLatencyItem.on('click', () => {
    latencyItemPressed();
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
        label: 'Auto',
      });

      const items = qualities.map(item => {
        const newMenuItem = new MenuItem(player, {
          selectable: true,
          label: item.name,
        });

        // Quality selected
        newMenuItem.on('click', () => {
          // If for some reason tech doesn't exist, then don't do anything
          if (!tech) {
            console.warn('Invalid attempt to access null player tech');
            return;
          }
          // Only enable this single, selected representation.
          tech.vhs.representations().forEach((rep, index) => {
            rep.enabled(index === item.index);
          });
          newMenuItem.selected(false);
        });

        return newMenuItem;
      });

      defaultAutoItem.on('click', () => {
        // Re-enable all representations.
        tech.vhs.representations().forEach(rep => {
          rep.enabled(true);
        });
        defaultAutoItem.selected(false);
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

  // If none of the settings in this menu are applicable then don't show it.
  const tech = player.tech({ IWillNotUseThisInPlugins: true });
  menuButton.addClass('vjs-quality-selector');
  videojs.registerComponent('MenuButton', MenuButton);

  if (!tech.vhs) {
    return menuButton;
  }

  if (qualities.length < 2 && (!tech || !tech.vhs)) {
    return menuButton;
  }

  // eslint-disable-next-line consistent-return
  return menuButton;
}
