/*
This is a finite state machine model that is used by xstate. https://xstate.js.org/
You send events to it and it changes state based on the pre-determined
modeling.
This allows for a clean and reliable way to model the current state of the
web application, and a single place to determine the flow of states.

You can paste this code into https://stately.ai/viz to see a visual state
map or install the VS Code plugin:
https://marketplace.visualstudio.com/items?itemName=statelyai.stately-vscode
*/

import { createMachine } from 'xstate';

export interface AppStateOptions {
  chatAvailable: boolean;
  chatLoading?: boolean;
  videoAvailable: boolean;
  appLoading?: boolean;
}

export function makeEmptyAppState(): AppStateOptions {
  return {
    chatAvailable: false,
    chatLoading: true,
    videoAvailable: false,
    appLoading: true,
  };
}

const OFFLINE_STATE: AppStateOptions = {
  chatAvailable: false,
  chatLoading: false,
  videoAvailable: false,
  appLoading: false,
};

const ONLINE_STATE: AppStateOptions = {
  chatAvailable: true,
  chatLoading: false,
  videoAvailable: true,
  appLoading: false,
};

const LOADING_STATE: AppStateOptions = {
  chatAvailable: false,
  chatLoading: false,
  videoAvailable: false,
  appLoading: true,
};

const GOODBYE_STATE: AppStateOptions = {
  chatAvailable: true,
  chatLoading: false,
  videoAvailable: false,
  appLoading: false,
};

export enum AppStateEvent {
  Loading = 'LOADING', // Have not pulled configuration data from the server.
  Loaded = 'LOADED', // Configuration data has been pulled from the server.
  Online = 'ONLINE', // Stream is live
  Offline = 'OFFLINE', // Stream is not live
  NeedsRegister = 'NEEDS_REGISTER', // Needs to register a chat user
  Fail = 'FAIL', // Error
  ChatUserDisabled = 'CHAT_USER_DISABLED', // Chat user is disabled
}

const appStateModel =
  /** @xstate-layout N4IgpgJg5mDOIC5QEMAOqDKAXZWwDoAbAe2QgEsA7KAYgCUBRAcQEkMAVBxgEUVFWKxyWcsUp8QAD0QBGAGwz8ABgCscpUoDsAZgAcKgEwrtATgA0IAJ6zNS-CZMLtcuQBY9Jg5t0BfHxbRMHDwiUgpqGgA5BgZuDAB9RlYOLgkBIRExCWkEGRVFVXUtPUNjcytEAxNXfB0DbSNNORMG119-EEDsXAISMipaABkAeQBBbli0wWFRcSQpWQVlNQ0dfSNTC2tc+vwZVwMZWxNbA5kDAz8A9G6QvvDaADFRlkGpjNns2VcCleL1spbRDaA74FS6ORVfYHTSObRXTo3YIEABOYCg5FgeBRA3ozDYnB47xmWXmOQOKj22hUnl02iajjk2iBCCqdgO2n2VRcbQhCK6yPwaIxWLAOIiz1exMyc1A5KMVJpBjpDJczIqCG0enwXk0MiUENMjiUBjk-KRPSFYDIlnwYkIVDANGGj0egxY0WlnzJiF0TXwulU9LqWhMMl0LIM7ipmguIIObU85qClrRNrtADMMw7KE7hpF3Z75ukSbKFgg5Pp7PSQdTXBp9uVtlHtDG464E7okx0BanrRBbVBiMQIAAjSxOyRYy3IDPYgAU2g0y4AlDReyE0wP8EOR+OwF7SXLff7A8ZNCHYeGWedW3qlDIWkz61rLgjKCO4BIN70wgND2W8o3pyyjKrCdZ6q4sYqMmtyouimLYv+xbTDKXyakuyiuHI+QmCoJquCo2FyCyS52PURzqI+mgqIYpqwYKW62vajoAehsJyFS+paPhfoRhqUa6PgTIuLoRznComiEWaPYWpu-bMVmOYHihHxHuWRQ6pCJz0sqGgNMBBjCfohiEUoIJ0kyDF9umu5jhObE+rkqh2BCLS8XhYa6K4wGUuGtEviYZ5qNZ8k2o5x4IJJnGmlUOixoG5kGCy+G1JyXgHGJJhKByrihQQsBigAbmKjzIOQhAAK5ohF5YXJx+q2MYbieFB4KkW0yj6NBfr6vU3n5fglWFSiGblVVNWqaW6H5HY4bNFJbhKI4HV3k0rgnNlS6xm+1wpngtU5NeGoALQXDqbVBct+g5Qofh+EAA */
  createMachine({
    id: 'appState',
    initial: 'loading',
    predictableActionArguments: true,
    states: {
      loading: {
        meta: LOADING_STATE,
        on: {
          NEEDS_REGISTER: {
            target: 'loading',
          },
          LOADED: {
            target: 'ready',
          },
          FAIL: {
            target: 'serverFailure',
          },
        },
      },
      ready: {
        initial: 'offline',
        states: {
          online: {
            meta: {
              ...ONLINE_STATE,
            },
            on: {
              OFFLINE: {
                target: 'goodbye',
              },
              CHAT_USER_DISABLED: {
                target: 'chatUserDisabled',
              },
            },
          },
          offline: {
            meta: {
              ...OFFLINE_STATE,
            },
            on: {
              ONLINE: {
                target: 'online',
              },
            },
          },
          goodbye: {
            on: {
              ONLINE: {
                target: 'online',
              },
            },
            meta: {
              ...GOODBYE_STATE,
            },
            after: {
              '300000': {
                target: 'offline',
              },
            },
          },
          chatUserDisabled: {
            meta: {
              ...ONLINE_STATE,
              chatAvailable: false,
            },
          },
        },
      },
      serverFailure: {
        type: 'final',
      },
      userfailure: {
        type: 'final',
      },
    },
  });

export default appStateModel;
