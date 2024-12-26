<script setup lang="ts">
import {onMounted, ref} from "vue";
import axios from "axios"


onMounted(() => {
  getMqttClient()
})


const mqttList = ref([{
  name:"",
  addr:"",
  port:0
}])

const getMqttClient = () => {
  axios({
    method: "get",
    url: "/mqtt/client/list"
  })
  .then(resp => {
    for (let idx in resp.data) {
      mqttList.value.push(resp.data[idx])
    }
  })
  .catch(err => {
    console.error(err)
  })
}


</script>

<template>
<!--  <main>-->
<!--    <TheWelcome />-->
    <el-table :data="mqttList" stripe style="width: 100%">
      <el-table-column prop="name" label="Name" width="180" />
      <el-table-column prop="addr" label="Address" />
      <el-table-column prop="port" label="port" width="180" />
      <el-table-column prop="enable" label="Enable" width="180" />
    </el-table>
<!--  </main>-->
</template>
