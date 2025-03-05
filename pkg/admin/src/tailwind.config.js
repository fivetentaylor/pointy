/** @type {import('tailwindcss').Config} */
const colors = require("tailwindcss/colors");

module.exports = {
  darkMode: ["class"],
  content: ["../templates/*.templ", "./**/*.tsx"],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: {
          DEFAULT: "hsl(var(--background))",
          "opacity-50": "hsl(var(--background-50))",
        },
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
          icon: "hsl(var(--muted-icon))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        highlight: {
          DEFAULT: "hsl(var(--highlight))",
          foreground: "hsl(var(--highlight-foreground))",
        },
        modal: {
          background: "hsla(var(--modal-background))",
        },
        elevated: {
          DEFAULT: "hsl(var(--elevated))",
        },
        reviso: {
          DEFAULT: "hsl(var(--reviso))",
        },
        "reviso-highlight": "hsla(var(--reviso-highlight))",
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      keyframes: {
        "accordion-down": {
          from: { height: 0 },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: 0 },
        },
        fadeOut: {
          "0%": { opacity: 1 },
          "100%": { opacity: 0 },
        },
        fadeIn: {
          "0%": { opacity: 0 },
          "100%": { opacity: 1 },
        },
        slideInY: {
          "0%": { transform: "translateY(100%)", opacity: 0 },
          "100%": { transform: "translateY(0)", opacity: 1 },
        },
        slideOutY: {
          "0%": { transform: "translateY(0)", opacity: 1 },
          "100%": { transform: "translateY(500%)", opacity: 0 },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
        slideInY: "slideInY 0.10s ease-out forwards",
        slideOutY: "slideOutY 0.05s ease-in forwards",
        fadeIn: "fadeIn 0.10s ease-out forwards",
        fadeOut: "fadeOut 0.15s ease-in forwards",
      },
      screens: {
        print: { raw: "print" },
        screen: { raw: "screen" },
      },
    },
    fontFamily: {
      sans: ["var(--font-inter)", "sans-serif"],
      serif: ["var(--font-marat)", "Times New Roman", "serif"],
      mono: [
        "ui-monospace",
        "SFMono-Regular",
        "Menlo",
        "Monaco",
        "Consolas",
        "Liberation Mono",
        "Courier New",
        "monospace",
      ],
    },
    typography: {
      DEFAULT: {
        css: {
          maxWidth: "65ch",
          h1: {
            fontSize: "2.125rem",
            lineHeight: "3rem",
            margin: "1rem 0",
            "&:first-child": {
              marginTop: 0,
            },
          },
          h2: {
            fontSize: "1.875rem",
            lineHeight: "2.625rem",
            margin: "1rem 0 0.5rem",
          },
          h3: {
            fontSize: "1.5rem",
            lineHeight: "2.375rem",
            margin: "1rem 0 0.5rem",
          },
          p: {
            fontSize: "1.125rem",
            lineHeight: "1.875rem",
            "+ p": {
              marginTop: "0.5rem",
            },
          },
          blockquote: {
            fontStyle: "italic",
            borderLeft: "2px solid hsla(0, 0%, 89%, 1)",
            paddingLeft: "1.1875rem",
            margin: "1rem 0",
          },
          a: {
            color: colors.sky[500],
          },
          ul: {
            fontSize: "1.125rem",
            lineHeight: "1.875rem",
            margin: "1rem 0",
            listStyle: "disc",
            paddingLeft: "1rem",
            "> li": {
              paddingLeft: "0.375rem",
            },
          },
          ol: {
            fontSize: "1.125rem",
            lineHeight: "1.875rem",
            margin: "1rem 0",
            listStyle: "decimal",
            paddingLeft: "1rem",
            "> li": {
              paddingLeft: "0.375rem",
            },
          },
        },
      },
    },
  },
  plugins: [
    require("tailwind-scrollbar"),
    require("tailwindcss-animate"),
    require("@tailwindcss/typography"),
  ],
};
