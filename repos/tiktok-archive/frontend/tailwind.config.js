/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'tt-pink': '#fe2c55',
        'tt-cyan': '#25f4ee',
        'tt-dark': '#121212',
        'tt-gray': '#1f1f1f',
      }
    },
  },
  plugins: [],
}
