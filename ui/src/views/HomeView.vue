<script setup lang="ts">
import {onMounted, reactive, ref} from "vue";
import {queryMqttClientList,saveMqttClient,delMqttClient} from "../api/mqtt";


onMounted(() => {
  getMqttClient()
})


const mqttList = ref()

const getMqttClient = async () => {
  const param = {}
  mqttList.value = await queryMqttClientList(param)
  // mqttList.value =
}


const handleDelete = async (index: number, row: any) => {
  console.log(index, row)
  const param = {
    addr: row.addr,
    port: row.port,
  }
  console.log(param)
  await delMqttClient(param).then(async () => {
    await getMqttClient()
  })
}


// 添加 MQTT 客户端 对话框开关
const dialogAddMqttVisible = ref(false)
const mqttItem = reactive({
  name: '',
  addr: '',
  port: 0,
  portStr: '',
  need_verify: false,
  username: '',
  password: '',
  enable: true,
})
const addMqttClient = async () => {
  mqttItem.port = Number(mqttItem.port)
  console.log(mqttItem)

  const param = {
    name: mqttItem.name,
    addr: mqttItem.addr,
    port: mqttItem.port,
    need_verify: mqttItem.need_verify,
    username: mqttItem.username,
    password: mqttItem.password,
    enable: mqttItem.enable,
  }

  await saveMqttClient(param).then(async () => {
    await getMqttClient()
    dialogAddMqttVisible.value = false
  })
}


</script>

<template>
  <main>
    <!--    <TheWelcome />-->
    <el-button type="primary" plain @click="dialogAddMqttVisible = true">新增 MQTT 客户端</el-button>
    <el-table :data="mqttList">
      <el-table-column prop="name" label="名称" width="180"/>
      <el-table-column prop="addr" label="地址"/>
      <el-table-column prop="port" label="端口" width="180"/>
      <el-table-column
          prop="enable"
          label="状态"
          width="180"
      >
        <template #default="scope">
          <el-tag
              :type="scope.row.enable ? 'success' : 'danger'"
          >{{ scope.row.enable ? '启用' : '关闭' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作按钮">
        <template #default="scope">
          <el-button
              size="small"
              type="danger"
              @click="handleDelete(scope.$index, scope.row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>


    <el-dialog v-model="dialogAddMqttVisible" title="添加 MQTT 客户端" draggable width="30%">
      <el-form :model="mqttItem" label-width="100px">
        <el-form-item label="客户端名称">
          <el-input v-model="mqttItem.name" placeholder="名称"></el-input>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="mqttItem.addr" placeholder="ip地址"></el-input>
        </el-form-item>
        <el-form-item label="端口">
          <el-input v-model="mqttItem.port" placeholder="端口号"></el-input>
        </el-form-item>
        <el-form-item label="需要验证？">
          <el-switch v-model="mqttItem.need_verify"></el-switch>
        </el-form-item>
        <el-form-item v-if="mqttItem.need_verify" label="用户名">
          <el-input v-model="mqttItem.username" placeholder="用户名"></el-input>
        </el-form-item>
        <el-form-item v-if="mqttItem.need_verify" label="密码">
          <el-input v-model="mqttItem.password" placeholder="密码"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <div>
          <el-button @click="dialogAddMqttVisible = false">取消</el-button>
          <el-button type="primary" @click="addMqttClient()">
            确认
          </el-button>
        </div>
      </template>
    </el-dialog>
  </main>
</template>


<style scoped>

.main {
  width: 100%;
}


</style>