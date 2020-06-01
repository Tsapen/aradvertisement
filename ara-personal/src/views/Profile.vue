<template>
  <div class="container">
    <h1>Create object</h1>
    <table><tbody>
      <div class="large-12 medium-12 small-12 cell">
      <tr>
        <td><input type="text" placeholder='latitude'  v-model="objGLTF.latitude" /></td>
        <td><input type="text" placeholder='longitude' v-model="objGLTF.longitude" /></td>
        <td><input type="text" placeholder='comment'  v-model="objGLTF.comment" /></td>
        <td><label><input type="file" id="file" ref="file" v-on:change="handleFileUpload()"/></label></td>
      </tr>
      <tr><td></td><td></td><td></td><td><button type="button" @click="submitFile()">Create object</button></td></tr>
      </div>
    </tbody></table>  
    <h1>Your objects</h1>
    <table><tbody>
      <tr><th>Location</th><th>Comment</th><th></th></tr>
      <tr v-for="objInfo in objectsInfo" :key=objInfo>
        <td>{{objInfo.location}}</td>
        <td><textarea type="text" cols="40" rows="3" v-model="objInfo.comment" /></td>
        <td><button type="button" @click="updateObject(objInfo)">Update</button></td>
        <td><button type="button" @click="deleteObject(objInfo.id)">Delete</button></td>
      </tr>
    </tbody></table>  
  </div>
</template>
<script>
import ObjGLTF from '../models/object';
import UserServices from '../services/user.service';

export default {
  name: 'Profile',
  data() {
    return {
      objectsInfo: [],
      objGLTF: new ObjGLTF('', '', '', '')
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    }
  },
  mounted() {
    if (!this.currentUser) {
      this.$router.push('/login');
    }
    this.getObjInfo()
  },
  methods: {
    createObject(){
    },

    getObjInfo(){
      UserServices.getObjInfo().then(result => this.objectsInfo = result);
    },

    updateObject(objInfo){
      UserServices.updateObjInfo(objInfo)
    },

    deleteObject(id){
      UserServices.deleteObject(id)
    },
    
    handleFileUpload(){
      this.objGLTF.gltf = this.$refs.file.files[0];
    },

    submitFile(){
      UserServices.uploadFile(this.objGLTF)
      this.objGLTF = new ObjGLTF('', '', '', '')
    }
  }
};
</script>