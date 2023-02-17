async function fetchData(url, options) {
	const ADMIN_USERNAME = 'admin';
	const ADMIN_STREAMKEY = 'abc123';

	const { data, method = 'GET', auth = true } = options || {};

	// eslint-disable-next-line no-undef
	const requestOptions = {
		method,
	};

	if (data) {
		requestOptions.body = JSON.stringify(data);
	}

	if (auth && ADMIN_USERNAME && ADMIN_STREAMKEY) {
		const encoded = btoa(`${ADMIN_USERNAME}:${ADMIN_STREAMKEY}`);
		requestOptions.headers = {
			Authorization: `Basic ${encoded}`,
		};
		requestOptions.mode = 'cors';
		requestOptions.credentials = 'include';
	}

	try {
		const response = await fetch(url, requestOptions);
		const json = await response.json();

		if (!response.ok) {
			const message =
				json.message || `An error has occurred: ${response.status}`;
			throw new Error(message);
		}
		return json;
	} catch (error) {
		console.error(error);
		return error;
		// console.log(error)
		// throw new Error(error)
	}
}

export default fetchData;
