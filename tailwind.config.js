/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/app/**/*.templ"],
  theme: {
    extend: {
      colors: {
        // https://clrs.cc/
        "clr-navy": "#001f3f",
        "clr-blue": "#0074D9",
        "clr-aqua": "#7FDBFF",
        "clr-teal": "#39CCCC",
        "clr-purple": "#B10DC9",
        "clr-fuchsia": "#F012BE",
        "clr-maroon": "#85144b",
        "clr-red": "#FF4136",
        "clr-orange": "#FF851B",
        "clr-yellow": "#FFDC00",
        "clr-olive": "#3D9970",
        "clr-green": "#2ECC40",
        "clr-lime": "#01FF70",
        "clr-black": "#111111",
        "clr-gray": "#AAAAAA",
        "clr-silver": "#DDDDDD",
      },
      fontFamily: {
        montserrat: ["Montserrat", "sans-serif"],
      }
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
    require("daisyui"),
  ],
  daisyui: {
    themes: ["light"]
  }
}

