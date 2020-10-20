import window from 'global/window';

export default function decodeB64ToUint8Array(b64Text) {
  const decodedString = window.atob(b64Text || '');
  const array = new Uint8Array(decodedString.length);

  for (let i = 0; i < decodedString.length; i++) {
    array[i] = decodedString.charCodeAt(i);
  }
  return array;
}
