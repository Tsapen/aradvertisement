import axios from 'axios';
import authHeader from './auth-header';

const API_URL = 'http://192.168.1.52:8000/api/';

class UserService {
  getObjInfo() {
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

  updateObjInfo(obj) {
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

  uploadFile(obj) {
    let formData = new FormData();
    
    formData.append('gltf', obj.gltf);
    let info = {
      latitude: parseFloat(obj.latitude),
      longitude: parseFloat(obj.longitude),
      comment: obj.comment
    };
    formData.append('info', JSON.stringify(info));
    
    axios.post(  API_URL + 'object',
      formData,
      {
      headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: authHeader().Authorization
        }   
      }
    ).then(function(){
      console.log('SUCCESS!!');
    })
    .catch(function(){
      console.log('FAILURE!!');
    });
  }
}

export default new UserService();