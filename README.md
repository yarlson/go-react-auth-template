# Go-React Auth Template

This repository is a GitHub template for building web applications with a Go backend and React frontend. It includes cookie-based Google authentication with refresh token functionality and uses Tailwind CSS for styling. The template aims to offer a starting point for developers creating web applications with modern authentication mechanisms.

## 📋 Table of Contents

- [Features](#-features)
- [Getting Started](#-getting-started)
- [Prerequisites](#-prerequisites)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- Go backend using Gin web framework
- React frontend with TypeScript
- Tailwind CSS for utility-first styling
- Google OAuth integration
- Secure cookie-based authentication
- Refresh token mechanism for persistent sessions
- PostgreSQL database for user and token storage
- CORS configuration for local development
- Monorepo structure using Turborepo

## 🚀 Getting Started

To use this template:

1. Click the "Use this template" button at the top of this repository.
2. Choose a name for your new repository and select its visibility.
3. Click "Create repository from template".

After creating your repository, clone it locally:

3. To run tests:

   ```
   npm run test
   ```

## 📖 API Documentation

### Authentication Endpoints

- `GET /auth/google`: Initiates Google OAuth login
- `GET /auth/google/callback`: Handles Google OAuth callback
- `POST /auth/refresh`: Refreshes the user's session
- `GET /auth/logout`: Logs out the user

### Protected Endpoints

- `GET /api/user/profile`: Retrieves the authenticated user's profile

All protected endpoints require a valid session cookie.

## 🗂 Project Structure

```
.
├── services/
│   ├── backend/
│   │   ├── auth/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── main.go
│   │   └── ...
│   └── frontend/
│       ├── src/
│       ├── public/
│       └── ...
├── package.json
├── turbo.json
└── README.md
```

## 🤝 Contributing

Contributions to improve this template are welcome. Please follow these steps:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

For more information or support, please open an issue on this repository.
