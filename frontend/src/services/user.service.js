import api from './api';

class UserService {
    getPublicContent() {
        return api.get('/api/me');
    }

    getUserBoard() {
        return api.get('/test/user');
    }

    getModeratorBoard() {
        return api.get('/test/mod');
    }

    getAdminBoard() {
        return api.get('/test/admin');
    }
}

export default new UserService();