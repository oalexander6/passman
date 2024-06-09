/** @type {import('tailwindcss').Config} */
const defaultTheme = require('tailwindcss/defaultTheme');

module.exports = {
  content: [
    'pkg/**/*.templ'
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter var', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        blue: {
          50: '#f0f3fe',
          100: '#dee4fb',
          200: '#c4d1f9',
          300: '#9cb4f4',
          400: '#6c8cee',
          500: '#4a67e7',
          600: '#3548db',
          700: '#2c36c9',
          800: '#2a2ea3',
          900: '#272c81',
          950: '#1c1e4f'
        }
      }
    }
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
};
