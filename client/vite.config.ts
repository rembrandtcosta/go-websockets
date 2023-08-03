import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'node:path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      room: resolve(__dirname, './room/index.html'),
      main: resolve(__dirname, './src/index.html'),
      //'@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
