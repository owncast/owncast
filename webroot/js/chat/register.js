// import {setLocalStorage} from "../utils/helpers";
import {URL_CHAT_REGISTRATION} from "../utils/constants.js";

// TODO: Be able to pass in the username they want to register as
// in case they have one stored in localstorage already.
export async function registerChat(username) {
    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({displayName: username})
    }

    try {
        const response = await fetch(URL_CHAT_REGISTRATION, options);
        const result = await response.json();
        console.log(result)
        return result;
    } catch(e) {
        console.error(e);
    }
}