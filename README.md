# **EverDownload: Your Universal Video Downloader**

EverDownload is a sleek and efficient web application built with Go, designed to simplify video downloading from various online platforms. 🚀 Just paste a video URL, and EverDownload fetches essential metadata, offers quality options, and lets you download your favorite content directly. It's fast, user-friendly, and leverages modern web technologies for a smooth experience! ✨

## ⚙️ Installation

Follow these steps to get EverDownload up and running on your local machine.

### Clone the Repository

First, clone the project repository to your local machine:

```bash
git clone https://github.com/samueltuoyo15/EverDownload.git
```

Then, navigate into the project directory:

```bash
cd EverDownload
```

### Prerequisites

Before running the application, ensure you have the following installed:

*   Go (version 1.24.1 or higher recommended)
*   A running Redis instance (local or remote)

### Environment Configuration

The application relies on several environment variables for configuration, especially for connecting to Redis and interacting with the video metadata API.

🌱 Create a `.env` file in the root directory of the project (e.g., `EverDownload/.env`).

📝 Populate the `.env` file with the following variables. You'll need to obtain an API key and host from a video metadata API provider (e.g., RapidAPI) and configure your Redis connection details.

```env
PORT=3000
SECRET_KEY=YOUR_PI_KEY
HOST=YOUR_API_HOST
API_ENDPOINT=YOUR_API_ENDPOINT 
REDIS_URL=localhost:6379 # Your Redis server address
REDIS_PASSWORD= # Your Redis password, if any
```

### Run the Application

Once the `.env` file is set up, you can build and run the Go application:

📦 Download Go module dependencies:

```bash
go mod tidy
```

▶️ Start the application:

```bash
go run main.go
```

You should see a message indicating the server is running, typically on `http://localhost:3000` (or the `PORT` you configured).

## 🚀 Usage

Using EverDownload is straightforward.

1.  **Access the Application**: Open your web browser and navigate to `http://localhost:3000` (or the port you configured in your `.env` file).
2.  **Paste Video URL**: On the homepage, you'll find an input field. Paste the URL of the video you wish to download (e.g., from YouTube, Instagram, etc.).
3.  **Fetch Details**: Click the "Download" button. The application will then dynamically fetch and display details about the video, including its thumbnail, title, source, author, and available quality options. A loading indicator will show while data is being fetched.
4.  **Select Quality & Download**: From the displayed options, select your preferred video quality using the dropdown menu. Once selected, click the "Download Video" button to initiate the download directly to your device.

The user interface leverages htmx and Alpine.js to provide a seamless and responsive experience without full page reloads.

## ✨ Features

*   **Universal Video Fetching**: Designed to download videos from a variety of online platforms by integrating with a powerful metadata API.
*   **Dynamic Metadata Retrieval**: Efficiently fetches crucial video details such as title, thumbnail, author, duration, and a list of available quality formats.
*   **Flexible Quality Selection**: Empowers users to choose their desired video resolution and format from the available options before downloading.
*   **Direct Download Capability**: Provides a direct, browser-initiated download link for the chosen video quality.
*   **Intelligent Redis Caching**: Implements Redis to cache video metadata, significantly reducing API calls and improving response times for subsequent requests.
*   **Modern Frontend**: Boasts a responsive and interactive user interface built with the powerful combination of htmx for dynamic content, Alpine.js for lightweight interactivity, and Tailwind CSS for rapid styling.
*   **Environment-Agnostic Configuration**: Utilizes `.env` files for easy and secure management of environment-specific variables, making deployment flexible.

## 🛠️ Technologies Used

| Technology       | Description                                                                 | Link                                                                   |
| :--------------- | :-------------------------------------------------------------------------- | :--------------------------------------------------------------------- |
| **Go**           | The primary backend language, valued for its performance, concurrency, and robust standard library. | [go.dev](https://go.dev/)                                            |
| **htmx**         | A modern JavaScript-free approach for building dynamic web interfaces directly within HTML. | [htmx.org](https://htmx.org/)                                          |
| **Alpine.js**    | A lightweight, minimalist JavaScript framework for adding reactive behavior to HTML templates. | [alpinejs.dev](https://alpinejs.dev/)                                  |
| **Tailwind CSS** | A utility-first CSS framework that allows for rapid UI development by composing low-level utility classes. | [tailwindcss.com](https://tailwindcss.com/)                          |
| **Redis**        | An open-source, in-memory data structure store used here for caching video metadata to optimize performance. | [redis.io](https://redis.io/)                                          |
| **GoDotEnv**     | A Go package for loading environment variables from `.env` files, simplifying configuration management. | [github.com/joho/godotenv](https://github.com/joho/godotenv)           |
| **RapidAPI**     | (Implied) A third-party API marketplace used to source the video metadata extraction service. | [rapidapi.com](https://rapidapi.com/)                                  |

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## ✍️ Author Info

**Samuel Tuoyo**

A passionate software developer focused on building robust and scalable applications.

*   **Twitter**: [twitter.com/samueltuoyo](https://x.com/TuoyoS26091) (X)

---

[![Go Version](https://img.shields.io/badge/Go-1.24.1-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Made with htmx](https://img.shields.io/badge/Made%20with-htmx-236dd4?style=for-the-badge&logo=htmx&logoColor=white)](https://htmx.org/)

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)