/* eslint-disable prefer-destructuring */
const ADMIN_USERNAME = process.env.NEXT_PUBLIC_ADMIN_USERNAME;
const ADMIN_STREAMKEY = process.env.NEXT_PUBLIC_ADMIN_STREAMKEY;
const NEXT_PUBLIC_API_HOST = process.env.NEXT_PUBLIC_API_HOST;

const API_LOCATION = `${NEXT_PUBLIC_API_HOST}api/admin/`;

export const FETCH_INTERVAL = 15000;

// Current inbound broadcaster info
export const BROADCASTER = `${API_LOCATION}broadcaster`;

// Disconnect inbound stream
export const DISCONNECT = `${API_LOCATION}disconnect`;

// Change the current streaming key in memory
export const STREAMKEY_CHANGE = `${API_LOCATION}changekey`;

// Current server config
export const SERVER_CONFIG = `${API_LOCATION}serverconfig`;

// Get viewer count over time
export const VIEWERS_OVER_TIME = `${API_LOCATION}viewersOverTime`;

// Get currently connected clients
export const CONNECTED_CLIENTS = `${API_LOCATION}clients`;


// Get hardware stats
export const HARDWARE_STATS = `${API_LOCATION}hardwarestats`;



// Current Stream status (no auth)
// use `admin/broadcaster` instead
// export const STREAM_STATUS = '/api/status';

export async function fetchData(url) {
  const encoded = btoa(`${ADMIN_USERNAME}:${ADMIN_STREAMKEY}`);

  try {
    const response = await fetch(url, {
      headers: {
        'Authorization': `Basic ${encoded}`,
      },
      mode: 'cors',
      credentials: 'include',
    });
    if (!response.ok) {
      const message = `An error has occured: ${response.status}`;
      throw new Error(message);
    }

    const json = await response.json();
    return json;
  } catch (error) {
    console.log(error)
  }
  return {};
}
