export as namespace URLToolkit;

declare namespace URLToolkit {
  export type URLParts = {
    scheme: string;
    netLoc: string;
    path: string;
    params: string;
    query: string;
    fragment: string;
  };
}

declare var URLToolkit: {
  buildAbsoluteURL: (
    baseURL: string,
    relativeURL: string,
    opts?: { alwaysNormalize?: boolean }
  ) => string;
  parseURL: (url: string) => URLToolkit.URLParts | null;
  normalizePath: (path: string) => string;
  buildURLFromParts: (parts: URLToolkit.URLParts) => string;
};

export = URLToolkit;
