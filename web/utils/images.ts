export const MAX_IMAGE_FILESIZE = 2097152;
export const ACCEPTED_IMAGE_TYPES = ['image/png', 'image/jpeg', 'image/gif'];

export function getBase64(img: File | Blob, callback: (imageUrl: string | ArrayBuffer) => void) {
  const reader = new FileReader();
  reader.addEventListener('load', () => callback(reader.result));
  reader.readAsDataURL(img);
}

export function readableBytes(bytes: number): string {
  const index = Math.floor(Math.log(bytes) / Math.log(1024));
  const SIZE_UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const size = Number((bytes / Math.pow(1024, index)).toFixed(2)) * 1;

  return `${size} ${SIZE_UNITS[index]}`;
}
