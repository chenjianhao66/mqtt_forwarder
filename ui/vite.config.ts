import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  // server: {
  //   proxy: {
  //     '/mqtt': {
  //       target: 'http://127.0.0.1:8888/mqtt', // 代理接口
  //       // target: 'http://192.168.0.69:18080',
  //       changeOrigin: true,
  //       rewrite: path => path.replace(/^\/mqtt/, '')
  //     }
  //   }
  // }
})
