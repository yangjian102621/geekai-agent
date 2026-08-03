import vue from '@vitejs/plugin-vue'
import path from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import { defineConfig, loadEnv } from 'vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  return {
    plugins: [
      vue(),
      AutoImport({
        imports: [
          'vue',
          'vue-router',
          {
            '@vueuse/core': ['useMouse', 'useFetch'],
          },
        ],
        dts: true, // 生成 TypeScript 声明文件
      }),
    ],
    base: env.VITE_BASE_URL,
    build: {
      outDir: 'dist', // 构建输出目录
    },
    css: {
      preprocessorOptions: {
        scss: {
          api: 'modern-compiler', // or 'modern'
        },
      },
    },
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },

    server: {
      port: 8888, // 设置你想要的端口号
      open: false, // 可选：启动服务器时自动打开浏览器
      proxy:
        env.VITE_USE_PROXY === 'true'
          ? {
              '/api': {
                target: 'http://localhost:5678', // 目标服务器的地址
                changeOrigin: true, // 改变请求头中的 Origin 为目标服务器地址
              },
              '/static': {
                target: 'http://localhost:5678', // 目标服务器的地址
                changeOrigin: true, // 改变请求头中的 Origin 为目标服务器地址
              },
            }
          : {},
    },
  }
})
