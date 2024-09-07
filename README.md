# goauth

goauth is a robust authentication system built with Go, providing secure user authentication and session management. It integrates Google OAuth for user login and implements refresh token functionality for maintaining user sessions.

## 📋 Table of Contents

- [Installation](#-installation)
- [Usage](#-usage)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Contributing](#-contributing)
- [Testing](#-testing)
- [License](#-license)

## 📦 Installation

To install and set up the goauth project, follow these steps:

1. Clone the repository:

   ```
   git clone https://github.com/yarlson/goauth.git
   cd goauth
   ```

2. Install dependencies:

   ```
   go mod tidy
   ```

3. Set up the PostgreSQL database and update the `.env` file with your database credentials.

4. Set up Google OAuth credentials and update the `.env` file with your client ID and client secret.

## 🚀 Usage

To run the goauth server:

1. Ensure all environment variables are correctly set in the `.env` file.

2. Start the server:
   ```
   go run main.go
   ```

The server will start running on `http://localhost:8080`.

## ⚙️ Configuration

Create a `.env` file in the project root and add the following variables:

```
DATABASE_URL=postgresql://username:password@localhost:5432/dbname
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
HASH_KEY=32_byte_long_hash_key
BLOCK_KEY=32_byte_long_block_key
```

Ensure that `HASH_KEY` and `BLOCK_KEY` are 32 bytes long for secure cookie encryption.

## 📖 API Documentation

### Authentication Endpoints

- `GET /auth/google`: Initiates Google OAuth login
- `GET /auth/google/callback`: Handles Google OAuth callback
- `POST /auth/refresh`: Refreshes the user's session
- `GET /auth/logout`: Logs out the user

### Protected Endpoints

- `GET /api/user/profile`: Retrieves the authenticated user's profile

All protected endpoints require a valid session cookie.

## 🤝 Contributing

Contributions to goauth are welcome. Please follow these steps to contribute:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

Please ensure your code adheres to the project's coding standards and includes appropriate tests.

## 🧪 Testing

To run the tests for goauth:

```
go test ./...
```

Ensure that your database is set up and running before executing the tests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

For more information or support, please open an issue on the [GitHub repository](https://github.com/yarlson/goauth).
