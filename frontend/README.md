# Darts Game Frontend

A React-based frontend for the Darts word guessing game.

## Features

- User authentication (sign up/sign in)
- Word guessing game with semantic distance feedback
- Input validation (2-50 characters, letters only, Unicode support)
- Mock API mode for development without backend
- Real API mode for production use

## Prerequisites

- Node.js (v18 or higher recommended)
- npm or yarn

## Installation

```bash
npm install
```

## Running the Application

### Environment Variables

- `VITE_USE_MOCK_API` - Set to `true` to use mock API, `false` for real backend
- `VITE_API_URL` - Backend API URL (only used when mock is disabled)

To customize, copy `.env.example` to `.env.development` or `.env.production` and modify:

```bash
cp .env.example .env.development
```

### Development Mode (Mock API)

Run with mock data using localStorage - no backend required:

```bash
npm run dev
```

Make sure that `VITE_USE_MOCK_API` is set to true.


The app will be available at `http://localhost:5173`.

This mode:
- Uses mock API with data stored in localStorage
- Simulates network delays for realistic testing
- Generates random distances for guesses

### Development Mode (Backend API)

Connect to the real backend server:

```bash
npm run dev
```

Make sure your backend is running on `http://localhost:8080` or update the `VITE_API_URL` in `.env.development` and `VITE_USE_MOCK_API` is set to false.

## Technologies Used

- React 19 - UI library
- Vite 7 - Build tool and dev server
- Tailwind CSS 3 - Utility-first CSS framework
- Lucide React - Icon library
- ESLint - Code quality and linting
