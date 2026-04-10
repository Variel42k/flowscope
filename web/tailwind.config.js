/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        slatebg: '#0f172a',
        panel: '#111827',
        accent: '#16a34a',
        accent2: '#0ea5e9',
      },
      fontFamily: {
        sans: ['"IBM Plex Sans"', '"Segoe UI"', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      },
      boxShadow: {
        panel: '0 10px 30px rgba(2, 6, 23, 0.35)',
      },
    },
  },
  plugins: [],
}
