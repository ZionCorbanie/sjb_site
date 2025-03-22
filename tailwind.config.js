const colors = require("tailwindcss/colors");
const { addIconSelectors } = require("@iconify/tailwind");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["internal/templates/*.templ", "internal/templates/*.go"],
  theme: {
    container: {
      center: true,
      padding: {
        DEFAULT: "1rem",
        mobile: "2rem",
        tablet: "4rem",
        desktop: "5rem",
      },
    },
    extend: {
      colors: {
        primary: colors.white,
        secondary: colors.red,
        neutral: colors.stone,
      },
    },
  },
  plugins: [
      require("@tailwindcss/forms"), 
      require("@tailwindcss/typography"),
      addIconSelectors(["mdi-light", "vscode-icons"]),
  ],
};
