import {defineConfig} from 'vite'
import reactRefresh from '@vitejs/plugin-react-refresh'
import svgr from './plugins/svgr';

export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true,
      },
    },
  },
  plugins: [reactRefresh(), svgr()]
})
