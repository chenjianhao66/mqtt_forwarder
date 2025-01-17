<template>
  <main>
    <el-form :inline="true" :model="relay">
      <el-form-item label="ip地址">
        <el-input v-model="relay.addr" placeholder="聚英ip地址" clearable/>
      </el-form-item>
      <el-form-item label="端口">
        <el-input v-model="relay.port" placeholder="聚英端口" clearable/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="connect" :disabled="connectStatus">连接</el-button>
      </el-form-item>
    </el-form>

    <el-text v-if="connectStatus === true" type="success" size="large">DO口</el-text>
    <el-row v-if="connectStatus === true">
      <el-col :span="4" v-for="(item,index) in doProps">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO {{ index + 1 }}</el-text>
        <el-switch v-model="item.status" @change="switchChange(index)"></el-switch>
      </el-col>
    </el-row>

    <!--  DI  -->
    <el-text type="success" size="large" v-if="connectStatus === true">DI口</el-text>
    <el-row v-if="connectStatus === true">
      <el-col :span="4" v-for="(item,index) in diProps">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI {{ index + 1 }}</el-text>
        <el-button :type="item.status === true?'success':'danger'">
          {{ item.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
    </el-row>
  </main>

</template>

<script setup>


import {onUnmounted, reactive, ref} from "vue";
import {connectRelay, disconnectRelay, switchStatus, listRelay, relayStatus} from "../api/mqtt";

const dialogRelayVisible = ref(false)
const relay = reactive({
  addr: '',
  port: '',
})

const connectStatus = ref(false)
// 存放do口数据
const doProps = ref([])
const diProps = ref([])

const connect = async () => {
  console.log("addr:", relay.addr)
  console.log("port:", Number(relay.port))

  const data = {
    addr: relay.addr,
    port: Number(relay.port),
  }
  await connectRelay(data).then(async () => {
    console.log("连接聚英成功, address: ", relay.addr, "port: ", relay.port)
    connectStatus.value = true
    const param = {}
    const list = await listRelay(param)
    const data = list[0]
    console.log(data)
    doProps.value.push(data.DO1)
    doProps.value.push(data.DO2)
    doProps.value.push(data.DO3)
    doProps.value.push(data.DO4)
    doProps.value.push(data.DO5)
    doProps.value.push(data.DO6)
    doProps.value.push(data.DO7)
    doProps.value.push(data.DO8)

    diProps.value.push(data.DI1)
    diProps.value.push(data.DI2)
    diProps.value.push(data.DI3)
    diProps.value.push(data.DI4)
    diProps.value.push(data.DI5)
    diProps.value.push(data.DI6)
    diProps.value.push(data.DI7)
    diProps.value.push(data.DI8)


    console.log(doProps.value)
    console.log(diProps.value)
    connectSSE()

  }).catch(err => {
    console.log("连接失败")
  })
}

const eventSource = reactive({})

const connectSSE = () => {
  console.log("开始连接sse")
  const sseUrl = `/mqtt/relay/status?addr=${relay.addr}&port=${Number(relay.port)}`
  eventSource.value = new EventSource(sseUrl);
  eventSource.value.onmessage = (event) => {
    receiveSseData(JSON.parse(event.data))
  }

  // 监听错误事件
  eventSource.onerror = (error) => {
    console.error('SSE 错误:', error);
    eventSource.value.close(); // 关闭连接
  };
}

const receiveSseData = (data) => {
  data.DO1 && (doProps.value[0] = data.DO1)
  data.DO2 && (doProps.value[1] = data.DO2)
  data.DO3 && (doProps.value[2] = data.DO3)
  data.DO4 && (doProps.value[3] = data.DO4)
  data.DO5 && (doProps.value[4] = data.DO5)
  data.DO6 && (doProps.value[5] = data.DO6)
  data.DO7 && (doProps.value[6] = data.DO7)
  data.DO8 && (doProps.value[7] = data.DO8)

  data.DI1 && (diProps.value[0] = data.DI1)
  data.DI2 && (diProps.value[1] = data.DI2)
  data.DI3 && (diProps.value[2] = data.DI3)
  data.DI4 && (diProps.value[3] = data.DI4)
  data.DI5 && (diProps.value[4] = data.DI5)
  data.DI6 && (diProps.value[5] = data.DI6)
  data.DI7 && (diProps.value[6] = data.DI7)
  data.DI8 && (diProps.value[7] = data.DI8)
}

onUnmounted(() => {
  if (connectStatus.value === false) {
    return
  }
  console.log("聚英继电器组件销毁")
  const data = {
    addr: relay.addr,
    port: Number(relay.port),
  }
  disconnectRelay(data)
  if (eventSource.value) {
    console.log("关闭sse")
    eventSource.value.close();
  }
})

const switchChange = async (index) => {
  const doStatus = doProps.value.at(index);
  console.log("switchChange", index, doStatus)
  const data = {
    addr: relay.addr,
    port: Number(relay.port),
    pointNumber: index + 1,
    status: doStatus.status,
  }
  await switchStatus(data).then(() => {
    console.log("切换成功")
  }).catch(err => {
    console.log("切换失败")
  })
}


</script>


<style scoped>
.main {
  width: 100%;
}


/*布局*/
.el-row {
  margin-bottom: 20px;
}

.el-row:last-child {
  margin-bottom: 0;
}

.el-col {
  border-radius: 4px;
}

.grid-content {
  border-radius: 4px;
  min-height: 36px;
}

.text {
  margin-right: 5px;
}
</style>