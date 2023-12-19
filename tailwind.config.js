/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'templates/**/*.templ',
  ],
  darkMode: 'class',
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
  corePlugins: {
    preflight: true,
  }
}

