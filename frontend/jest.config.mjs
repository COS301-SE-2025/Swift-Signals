export default {
  preset: 'ts-jest/presets/js-with-ts',
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.(ts|tsx)$': ['ts-jest', { tsconfig: '<rootDir>/tsconfig.jest.json', useESM: true }],
    '^.+\\.(js|jsx)$': 'babel-jest',
  },
  moduleNameMapper: {
    '\\.css$': 'identity-obj-proxy', // maps CSS imports correctly
    '\\.(jpg|jpeg|png|gif|svg)$': '<rootDir>/jest.styleMock.js', 
  },
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
};
