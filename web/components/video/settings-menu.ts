/* eslint-disable max-classes-per-file */
export function createVideoSettingsMenuButton(
  player,
  videojs,
  qualities,
  latencyItemPressed: () => boolean,
): any {
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
    label: 'minimize latency (experimental)',
  });
  lowLatencyItem.on('click', () => {
    const enabled: boolean = latencyItemPressed();
    lowLatencyItem.selected(enabled);
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
