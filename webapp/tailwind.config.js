const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    screens: {
      xs: "475px",
      ...defaultTheme.screens,
    },
    fontFamily: {
      sans: ["Grandstander", "sans-serif", "Noto Color Emoji"],
    },
    extend: {
      keyframes: {
        "bounce-loop": {
          "0%, 100%": {
            transform: "translateY(0)",
            animationTimingFunction: "cubic-bezier(0.8, 0, 1, 1)",
          },
          "40%, 60%": {
            transform: "translateY(0)",
            animationTimingFunction: "cubic-bezier(0.8, 0, 1, 1)",
          },
          "50%": {
            transform: "translateY(-25%)",
            animationTimingFunction: "cubic-bezier(0, 0, 0.2, 1)",
          },
        },
        wiggle: {
          "0%, 100%": { transform: "rotate(0)" },
          "17%": { transform: "rotate(-7deg)" },
          "33%": { transform: "rotate(7deg)" },
          "47%": { transform: "rotate(-15deg)" },
          "66%": { transform: "rotate(15deg)" },
          "84%": { transform: "rotate(-15deg)" },
        },
      },
      animation: {
        "bounce-loop": "bounce-loop 3s ease-in-out infinite",
        "wiggle-short": "wiggle 1s ease-in-out 1",
      },
    },
  },
  variants: {
    extend: {
      backgroundColor: ["disabled"],
      borderColor: ["disabled"],
      textColor: ["disabled"],
      cursor: ["disabled"],
    },
  },
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};
