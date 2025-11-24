import {heroui} from "@heroui/theme"
import textshadow from "tailwindcss-textshadow"

/** @type {import('tailwindcss').Config} */
const config = {
  content: [
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-sans)"],
        mono: ["var(--font-mono)"],
      },
      textShadow: {
        sm: "1px 1px 2px rgba(0,0,0,0.25)",
        DEFAULT: "2px 2px 4px rgba(0,0,0,0.3)",
        lg: "3px 3px 6px rgba(0,0,0,0.4)",
      },
    },
  },
  darkMode: "class",
  plugins: [heroui()],
}

module.exports = config;