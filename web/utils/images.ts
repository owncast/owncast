export const ACCEPTED_IMAGE_TYPES = ['image/png', 'image/jpeg', 'image/gif'];

export function getBase64(img: File | Blob, callback: (imageUrl: string | ArrayBuffer) => void) {
  const reader = new FileReader();
  reader.addEventListener('load', () => callback(reader.result));
  reader.readAsDataURL(img);
}
