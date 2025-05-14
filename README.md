# Swift-Signals

# Project Structure
```
swift-signals-frontend/
├── public/                # Static assets
│   ├── favicon.svg
│   └── index.html
├── src/
│   ├── assets/            # Images, fonts, etc.
│   ├── components/        # Reusable UI components (buttons, charts, etc.)
│   ├── constants/         # Static config or enums (e.g., traffic light states)
│   ├── context/           # React context providers (e.g., auth, UI themes)
│   ├── hooks/             # Custom React hooks
│   ├── layouts/           # Shared page layouts (e.g., dashboard shell)
│   ├── pages/             # Main routes/views (e.g., /simulate, /config)
│   ├── routes/            # Route configuration
│   ├── services/          # API calls, axios instances, etc.
│   ├── store/             # Zustand, Redux, or context for global state
│   ├── styles/            # Tailwind config, global styles
│   ├── utils/             # Utility/helper functions
│   ├── main.tsx           # App entry point
│   └── App.tsx            # Main app component
├── .env                   # Environment variables
├── tailwind.config.js     # TailwindCSS config
├── postcss.config.js      # PostCSS config
├── tsconfig.json          # TypeScript config
├── vite.config.ts         # Vite config
├── package.json
└── README.md
```
