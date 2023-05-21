export function pluralize(string, count) {
  if (count === 1) {
    return string;
  }
  return `${string}s`;
}

export function getDiffInDaysFromNow(timestamp) {
  const time = typeof timestamp === 'string' ? new Date(timestamp) : timestamp;
  return (new Date() - time) / (24 * 3600 * 1000);
}

// Take a nested object of state metadata and merge it into
// a single flattened node.
export function mergeMeta(meta) {
  return Object.keys(meta).reduce((acc, key) => {
    const value = meta[key];
    Object.assign(acc, value);

    return acc;
  }, {});
}
