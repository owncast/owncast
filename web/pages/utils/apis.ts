

const IS_DEV = true;
const ADMIN_USERNAME = 'admin';
const ADMIN_STREAMKEY = 'abc123';

const API_LOCATION = 'http://localhost:8080/api/admin/';

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
    // waits until the request completes...
    // console.log(response);

    if (!response.ok) {
      console.log(response)
      const message = `An error has occured: ${response.status}`;
      throw new Error(message);
    }

    const json = await response.json();
    console.log(json)
    return json;
  } catch (error) {
    console.log(error)
  }
}

// fetch error cases
// json.catch(error => {
//   error.message; // 'An error has occurred: 404'
// });
