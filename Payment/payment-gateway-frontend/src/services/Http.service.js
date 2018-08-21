import axios from 'axios';

export const instanceAxios = axios.create({
	// baseURL: 'https://stg-api-heimdall.rockship.co/'
//	baseURL: 'https://dev-api-heimdall.rockship.co/'
	baseURL: 'https://97ced5dc.ngrok.io'
});

export const sendRequest = async (method, url, data, headers = {}) => {
	return new Promise((resolve) => {
		instanceAxios.request({
			url: url,
			method: method,
			data: data,
			headers: headers
		}).then((res) => {
			resolve({
				data: res.data,
				isError: false
			});
		}).catch((error) => {
			resolve({
        		data: null,
				isError: true,
				err: error
    		});
		});
	});
};
