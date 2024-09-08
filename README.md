# Go-React Auth Template

This repository provides a template for building web applications with a Go backend and React frontend. It includes cookie-based Google authentication with refresh token functionality and uses Tailwind CSS for styling. The template aims to offer a starting point for developers creating web applications with modern authentication mechanisms.

## 📋 Table of Contents

- [Features](#-features)
- [Using This Template](#-using-this-template)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
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

## 🔧 Using This Template

To use this template for your project:

1. Click the "Use this template" button at the top of this repository.
2. From the dropdown, select "Create a new repository".
3. On the next page:
   - Choose a name for your new repository.
   - Select the owner (your account or an organization).
   - Decide whether you want your new repository to be public or private.
   - Optionally, you can choose to include all branches.
4. Click "Create repository from template".

This process will create a new repository in your account (or chosen organization) with the same files and structure as this template.

## 🛠 Prerequisites

Before you begin, ensure you have the following installed:
- Go (1.16 or later)
- Node.js (14.x or later)
- npm (6.x or later)
- PostgreSQL (12.x or later)

## 📦 Installation

After creating your repository from this template:

1. Clone your new repository:
   ```
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
   ```

2. Install backend dependencies:
   ```
   go mod tidy
   ```

3. Install frontend dependencies:
   ```
   cd dashboard
   npm install
   ```

## ⚙️ Configuration

1. Create a `.env` file in the project root with the following variables:

   ```
   DATABASE_URL=postgresql://username:password@localhost:5432/dbname
   GOOGLE_CLIENT_ID=your_google_client_id
   GOOGLE_CLIENT_SECRET=your_google_client_secret
   HASH_KEY=32_byte_long_hash_key
   BLOCK_KEY=32_byte_long_block_key
   ```

   Ensure `HASH_KEY` and `BLOCK_KEY` are 32 bytes long for secure cookie encryption.

2. Set up your Google OAuth credentials in the Google Developer Console.
   - Create a new project (or select an existing one).
   - Configure the OAuth consent screen.
   - Create OAuth 2.0 Client IDs for Web application.
   - Set the authorized JavaScript origins and redirect URIs.

3. Update the `callback` URL in `auth/handler.go` to match your frontend URL.

## 🚀 Usage

To run the application:

1. Start the backend server:
   ```
   go run main.go
   ```
   The server will run on `http://localhost:8080`.

2. In a separate terminal, start the frontend development server:
   ```
   cd dashboard
   npm run dev
   ```
   The React app will be available at `http://localhost:5173`.

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
├── auth/           # Authentication handlers
├── dashboard/      # React frontend
├── model/          # Database models
├── repository/     # Data access layer
├── main.go         # Main application entry
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

For more information or support, please open an issue on this template repository. If you've created a project using this template and have questions about your specific implementation, please open issues in your project's repository.

Remember to update this README and other documentation to reflect your specific project as you develop it.
