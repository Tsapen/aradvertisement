import axios from 'axios';
import authHeader from './auth-header';

const API_URL = 'https://192.168.1.52:8000/api/';

class UserService {
  getObjects() {
    return axios({ 
      method: 'get',
      url: API_URL + 'user_objects',
      headers: authHeader(),
      responseType: 'json',
    })
    .then ( function (response){
      return response.data.response; 
    })
  }

  updateObjects(obj) {
    axios({ 
      method: 'post',
      url: API_URL + 'object/upd',
      headers: authHeader(),
      data: {id: obj.id,comment: obj.comment},
      responseType: 'json',
    })
  }

  deleteObject(id) {
    axios({ 
      method: 'post',
      url: API_URL + 'object/del',
      headers: authHeader(),
      data: {id: id},
      responseType: 'json',
    })
  }
}

export default new UserService();