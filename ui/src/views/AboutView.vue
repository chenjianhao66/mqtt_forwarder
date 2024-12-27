<template>
  <div class="about">
    <!--    <h1>This is an about page</h1>-->

    <el-button type="primary" plain @click="dialogAddForwarderVisible = true">新增转发器</el-button>
    <el-table :data="forwarderList">
      <el-table-column prop="name" fixed label="名称" width="140"/>
      <el-table-column prop="sourceItemName" label="源"/>
      <el-table-column prop="sourceTopic" label="源主题" width="140"/>
      <el-table-column prop="targetItemName" label="目标"/>
      <el-table-column prop="targetTopic" label="目标主题" width="140"/>

      <el-table-column
          prop="enable"
          label="状态"
          width="80"
      >
        <template #default="scope">
          <el-tag
              :type="scope.row.enable ? 'success' : 'danger'"
          >{{ scope.row.enable ? '启用' : '关闭' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作按钮" fixed="right" width="150">
        <template #default="scope">
          <el-button
              size="small"
              type="primary"
              @click="handleSwitch(scope.$index, scope.row)"
          >切换
          </el-button>

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

    <el-dialog v-model="dialogAddForwarderVisible" title="添加转发器" draggable width="30%" @open="getMqttClientR()">
      <el-form :model="forwarderItem" label-width="100px">
        <el-form-item label="转发器名称">
          <el-input v-model="forwarderItem.name" placeholder="名称"></el-input>
        </el-form-item>

        <el-form-item label="源">
          <el-select v-model="forwarderItem.source" value-key="name" placeholder="选择要转发的MQTT客户端">
            <el-option
                v-for="client in mqttList"
                :key="client.name"
                :label="client.name"
                :value="client"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="源主题">
          <el-input v-model="forwarderItem.sourceTopic" placeholder="源topic"></el-input>
        </el-form-item>
        <el-form-item label="目标">
          <el-select v-model="forwarderItem.target" value-key="name" placeholder="选择转发的目标MQTT客户端">
            <el-option
                v-for="item in mqttList"
                :key="item.name"
                :label="item.name"
                :value="item"
            />

          </el-select>
        </el-form-item>
        <el-form-item label="目标主题">
          <el-input v-model="forwarderItem.targetTopic" placeholder="目标topic"></el-input>
        </el-form-item>


      </el-form>
      <template #footer>
        <div>
          <el-button @click="dialogAddForwarderVisible = false">取消</el-button>
          <el-button type="primary" @click="addForwarder()">
            确认
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>


<script setup>

import {onMounted, reactive, ref} from "vue";
import {listForwarder, enableForwarder, disableForwarder, queryMqttClientList,saveForwarder,delForwarder} from "../api/mqtt";


onMounted(() => {
  getForwardList()
})


const forwarderList = ref()
const getForwardList = async () => {

  const params = {}
  forwarderList.value = await listForwarder(params)
}
const handleSwitch = async (index, row) => {

  if (row.enable) {
    await disableForwarder(row.id)
  } else {
    await enableForwarder(row.id)
  }
  await getForwardList()

}

const handleDelete = async (index, row) => {
  await delForwarder(row)
  await getForwardList()
}


const dialogAddForwarderVisible = ref(false)
const forwarderItem = reactive({
  name: '',
  source: {},
  sourceTopic: '',
  target: {},
  targetTopic: '',
});
const mqttList = ref([])

const getMqttClientR = async () => {
  mqttList.value = []
  mqttList.value = await queryMqttClientList()
}

const addForwarder = async () => {
  const data = {
    name: forwarderItem.name,
    enable: true,
    sourceItemAddr: forwarderItem.source.addr,
    sourceItemPort: forwarderItem.source.port,
    sourceItemName: forwarderItem.source.name,
    sourceTopic: forwarderItem.sourceTopic,
    targetItemAddr: forwarderItem.target.addr,
    targetItemPort: forwarderItem.target.port,
    targetItemName: forwarderItem.target.name,
    targetTopic: forwarderItem.targetTopic,
  }
  console.log(data)
  await saveForwarder(data)
  // 重新获取列表
  await getForwardList()
  dialogAddForwarderVisible.value = false
}

</script>

<style scoped>
.about {
  width: 100%;
}
</style>
