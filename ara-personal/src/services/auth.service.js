import axios from 'axios';

const API_URL = 'http://192.168.1.52:8000/api/auth/';

class AuthService {
  login(user) {
    return axios
      .post(API_URL + 'login', {
        username: user.username,
        password: user.password
      })
      .then(response => {
        if (response.data) {
          localStorage.setItem('user', JSON.stringify(response.data));
        }

        return response.data;
      });
  }

  logout() {
    localStorage.removeItem('user');
  }

  register(user) {
    return axios.post(API_URL + 'registration', {
      username: user.username,
      email: user.email,
      password: user.password
    });
  }
}

export default new AuthService();