import react from '@vitejs/plugin-react'
import process from 'node:process'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  const env = loadEnv(mode, process.cwd(), '')

  return {
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: env.VITE_API_TARGET || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      }
    }
  }
}
})