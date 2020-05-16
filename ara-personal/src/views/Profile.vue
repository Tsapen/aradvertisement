<template>
  <div class="container">
    <table><tbody>
      <button type="button" @click="createObject()">New object</button>
      <tr><th>Location</th><th>Comment</th><th></th></tr>
      <tr v-for="object in objects" :key=object>
        <td>{{object.location}}</td>
        <td><textarea type="text" cols="40" rows="3" v-model="object.comment" /></td>
        <td><button type="button" @click="updateObject(object)">Update</button></td>
        <td><button type="button" @click="deleteObject(object.id)">Delete</button></td>
      </tr>
    </tbody></table>  
  </div>
</template>
<script>
import '../models/objectInfo';
import UserServices from '../services/user.service';
import objectInfo from '../models/objectInfo';

export default {
  name: 'Profile',
  data() {
    return {
      obj: new objectInfo,
      objects: []
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
    this.getObjects()
  },
  methods: {
    createObject(){
    },
    getObjects(){
      UserServices.getObjects().then(result => this.objects = result);
    },
    updateObject(obj){
      UserServices.updateObjects(obj)
    },
    deleteObject(id){
      UserServices.deleteObject(id)
    }
  }
};
</script>