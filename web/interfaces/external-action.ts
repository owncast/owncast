export interface ExternalAction {
  title: string;
  description?: string;
  color?: string;
  url?: string;
  html?: string;
  icon?: string;
  openExternally?: boolean;
}

/**
 * Utility functions for working with ExternalAction objects
 */
export namespace ExternalActionUtils {
  /**
   * Generates a unique key for an external action based on its properties.
   * This ensures each action has a unique identifier for React rendering.
   */
  export function generateKey(action: ExternalAction): string {
    return `${action.title}-${action.url || action.html || 'no-url'}-${action.color || 'no-color'}`;
  }

  /**
   * Finds an action from an array based on a generated key.
   * Useful for finding actions in event handlers when only the key is available.
   */
  export function findByKey(actions: ExternalAction[], key: string): ExternalAction | undefined {
    return actions.find(action => generateKey(action) === key);
  }
}
