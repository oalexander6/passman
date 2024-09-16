import react from '@vitejs/plugin-react';
import { config as dotenvConfig } from 'dotenv';
import path from 'path';
import { defineConfig } from 'vite';

dotenvConfig({
  path: '../.env'
});

const serverPort = process.env.PORT;

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: `http://localhost:${serverPort}/`,
        changeOrigin: true
      }
    }
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  }
});
