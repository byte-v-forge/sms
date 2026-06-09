import path from 'path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  base: '/',
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: [{ find: '@', replacement: path.resolve(__dirname, './src') }]
  },
  build: {
    target: 'esnext',
    rollupOptions: {
      output: {
        manualChunks: {
          'react-vendor': ['react', 'react-dom', 'react-router', '@tanstack/react-query'],
          'ui-vendor': ['cmdk', '@radix-ui/react-checkbox', '@radix-ui/react-popover', '@radix-ui/react-slot', '@radix-ui/react-switch', 'lucide-react']
        }
      }
    }
  }
});
