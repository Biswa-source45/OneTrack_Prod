/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          navy: '#1e3a8a',
          sky: '#38bdf8',
          teal: '#14b8a6',
          cyan: '#06b6d4',
        },
        neutral: {
          white: '#ffffff',
          light: '#f8fafc',
          medium: '#64748b',
          dark: '#1e293b',
        }
      },
    },
  },
  plugins: [],
}

