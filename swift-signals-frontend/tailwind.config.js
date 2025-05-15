/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class', // Enable dark mode support
  theme: {
    extend: {
      colors: {
        customIndigo: '#6D53FF',
        customGreen: '#2E8244',
        customPurple: '#8749AA',
        statusGreen: '#82DFB0',
        statusYellow: '#EDD9A3',
        statusRed: '#FAC0C1',
        statusTextGreen: '#006234',
        statusTextYellow: '#886800',
        statusTextRed: '#8E0F11',
      },
    },
  },
  plugins: [],
}
