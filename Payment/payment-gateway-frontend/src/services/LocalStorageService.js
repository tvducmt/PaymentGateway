const localstorage = window.localStorage

export const LocalStorageService = {
	set: (key, value) => {
		localstorage.setItem(key, JSON.stringify(value))
	},

	get: (key) => {
		try {
			let data = JSON.parse(localstorage.getItem(key))
			return data
		} catch(e) {
			localstorage.removeItem(key)
			window.location.href = '/'
		} 
	},
	getTK: (key) => {
		try {
			var token = localstorage.getItem(key)
			var base64Url = token.split('.')[1];
			var base64 = base64Url.replace('-', '+').replace('_', '/');
			return JSON.parse(window.atob(base64));
		}
		catch(err) {
			localstorage.removeItem(key)
			window.location.href = '/'
		}
		
	},
	remove: (key) => {
		localstorage.removeItem(key)
	}
}
