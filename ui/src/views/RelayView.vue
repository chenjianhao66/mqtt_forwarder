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

    <!--    <div v-if="connectStatus === true">-->
    <!--      ttttttt-->
    <!--    </div>-->
    <el-text v-if="connectStatus === true" type="success" size="large">DO口</el-text>
    <el-row v-if="connectStatus === true">
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO1</el-text>
        <el-switch v-model="DO1.status"></el-switch>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO2</el-text>
        <el-switch v-model="DO2.status"></el-switch>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO3</el-text>
        <el-switch v-model="DO3.status"></el-switch>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO4</el-text>
        <el-switch v-model="DO4.status"></el-switch>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO5</el-text>
        <el-switch v-model="DO5.status"></el-switch>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO6</el-text>
        <el-switch v-model="DO6.status"></el-switch>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO7</el-text>
        <el-switch v-model="DO7.status"></el-switch>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DO8</el-text>
        <el-switch v-model="DO8.status"></el-switch>
      </el-col>
    </el-row>

    <!--  DI  -->

    <el-text type="success" size="large" v-if="connectStatus === true">DI口</el-text>
    <el-row v-if="connectStatus === true">
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI1</el-text>
        <el-button :type="DI1.status === true?'success':'danger'">
          {{ DI1.status == true ? "在位" : "不在位" }}
        </el-button>

      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI2</el-text>
        <el-button :type="DI2.status === true?'success':'danger'">
          {{ DI2.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI3</el-text>
        <el-button :type="DI3.status === true?'success':'danger'">
          {{ DI3.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI4</el-text>
        <el-button :type="DI4.status === true?'success':'danger'">
          {{ DI4.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI5</el-text>
        <el-button :type="DI5.status === true?'success':'danger'">
          {{ DI5.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI6</el-text>
        <el-button :type="DI6.status === true?'success':'danger'">
          {{ DI6.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI7</el-text>
        <el-button :type="DI7.status === true?'success':'danger'">
          {{ DI7.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
      <el-col :span="4">
        <div class="grid-content ep-bg-purple"/>
        <el-text type="primary" class="text">DI8</el-text>
        <el-button :type="DI8.status === true?'success':'danger'">
          {{ DI8.status == true ? "在位" : "不在位" }}
        </el-button>
      </el-col>
    </el-row>
  </main>

</template>

<script setup>


import {onUnmounted, reactive, ref} from "vue";
import {connectRelay, disconnectRelay,switchStatus, listRelay} from "../api/mqtt";

const dialogRelayVisible = ref(false)
const relay = reactive({
  addr: '192.168.99.113',
  port: '10000',
})

const connectStatus = ref(false)

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
    DO1.value = data.DO1
    DO2.value = data.DO2
    DO3.value = data.DO3
    DO4.value = data.DO4
    DO5.value = data.DO5
    DO6.value = data.DO6
    DO7.value = data.DO7
    DO8.value = data.DO8
    connectSSE()

  }).catch(err => {
    console.log("连接失败")
  })
}

const eventSource = reactive({})

const connectSSE = () => {
  console.log("开始连接sse")
  eventSource.value = new EventSource('http://localhost:8888/mqtt/relay/status');
  eventSource.onmessage = (event) => {
    console.log(event)
  }

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
})

// DO口
const DO1 = reactive({
  name: '',
  status: false,
})
const DO2 = reactive({
  name: '',
  status: false,
})
const DO3 = reactive({
  name: '',
  status: false,
})
const DO4 = reactive({
  name: '',
  status: false,
})
const DO5 = reactive({
  name: '',
  status: false,
})
const DO6 = reactive({
  name: '',
  status: false,
})
const DO7 = reactive({
  name: '',
  status: false,
})
const DO8 = reactive({
  name: '',
  status: false,
})

// DI口
const DI1 = reactive({
  name: '',
  status: false,
})
const DI2 = reactive({
  name: '',
  status: false,
})
const DI3 = reactive({
  name: '',
  status: false,
})
const DI4 = reactive({
  name: '',
  status: false,
})
const DI5 = reactive({
  name: '',
  status: false,
})
const DI6 = reactive({
  name: '',
  status: false,
})
const DI7 = reactive({
  name: '',
  status: false,
})
const DI8 = reactive({
  name: '',
  status: false,
})


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