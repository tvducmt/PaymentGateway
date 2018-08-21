import { LocalStorageService } from './LocalStorageService'

export const AuthenticateService = {
    setAuthenticateUser: (token, user) => {
        LocalStorageService.set('currentUser', JSON.stringify(user))
        LocalStorageService.set('accesstoken', token)
        window.location.href = '/'
    },
    isAuthenticate: () => {
        const user = LocalStorageService.get('accesstoken')
        if (!user) return false
        return true
    },
    getAuthenticateUser: () => {
        try {
            let data = JSON.parse(LocalStorageService.get('currentUser'))
            return data
        } catch (e) {
            LocalStorageService.remove('currentUser')
            LocalStorageService.remove('accesstoken')
            window.location.href = '/'
            return null
        }
    },
    getAuthenticateGmail: () => {
        try {
            let data = JSON.parse(LocalStorageService.get('currentUser'));
            let gmail = data.username;
            return gmail
        } catch (e) {
            // console.log("Error: ", e);
            LocalStorageService.remove('currentUser');
            LocalStorageService.remove('accesstoken');
            // window.location.href = '/';
            return null;
        }
    },
    removeAuthenticate: () => {
        LocalStorageService.remove('currentUser')
        LocalStorageService.remove('accesstoken')
        window.location.href = '/'
    },
}