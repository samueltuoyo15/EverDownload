# EverDownload: Universal Video Downloader

## Overview
EverDownload is a high-performance Go-based web application designed for fetching metadata and downloading videos from a multitude of online platforms. It leverages the robust `yt-dlp` utility for media extraction, `Redis` for efficient metadata caching, and provides a modern, interactive frontend experience using `HTMX` and `Alpine.js`.

<img src="/public/demo.png" alt="EverDownload Demo">
    
## Features
- **Go**: High-performance, concurrent backend for handling video processing requests.
- **yt-dlp**: Powerful, versatile command-line program for downloading videos from over 1000 websites.
- **Redis**: In-memory data store used for caching video metadata to reduce API calls and improve response times.
- **Docker & Docker Compose**: Containerization for easy setup, deployment, and consistent environment management.
- **HTMX**: Enhances HTML with AJAX capabilities, enabling dynamic content loading without full page reloads for a smoother user experience.
- **Alpine.js**: A lightweight JavaScript framework providing reactive and declarative templating directly in HTML.
- **TailwindCSS**: Utility-first CSS framework for rapid and responsive UI development.

## Getting Started

### Installation
To get EverDownload up and running on your local machine, follow these steps:

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/samueltuoyo15/EverDownload.git
    cd EverDownload
    ```
2.  **Build and Run with Docker Compose**:
    Ensure you have Docker and Docker Compose installed. This will build the Go application, pull a Redis image, and start both services.
    ```bash
    docker-compose up --build -d
    ```
    The application will be accessible at `http://localhost:8080`.

### Environment Variables
The application requires the following environment variables, typically configured via `docker-compose.yml` or a `.env` file for local development:

*   `REDIS_URL`: The address of the Redis server.
    *   Example: `REDIS_URL=redis:6379` (for Docker Compose setup)
    *   Example: `REDIS_URL=localhost:6379` (for local Redis instance)
*   `REDIS_PASSWORD`: The password for Redis authentication. Leave empty if no password is set.
    *   Example: `REDIS_PASSWORD=YourSecurePassword`
    *   Example: `REDIS_PASSWORD=` (if no password)
*   `PORT`: The port on which the Go application server will listen.
    *   Example: `PORT=8080`

## Usage
Once the EverDownload API is running, you can interact with it through its web interface or directly via its API endpoints.

**Web Interface**:
Navigate to `http://localhost:8080` in your web browser. You will see an input field where you can paste a video URL. After submitting, the application will fetch video details and present you with options to select download quality before initiating the download.

**Supported Platforms**:
EverDownload supports a wide range of platforms thanks to `yt-dlp`, including but not limited to: YouTube, Instagram, X (Twitter), Facebook, TikTok, LinkedIn, Snapchat, Pinterest, Vimeo, Twitch, Threads, Reddit, Discord, Bilibili, Rumble, and Kick.

## API Documentation

### Base URL
`http://localhost:8080` (or the host and port where your service is deployed)

### Endpoints

#### `GET /`
Serves the main HTML page for the EverDownload application.
**Description**: Renders the `templates/index.html` file, providing the user interface for video URL submission.
**Request**: No specific payload required.
**Response**:
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <!-- ... (truncated for brevity) ... -->
</head>
<body class="bg-neutral-900 text-white min-h-screen">
  <div class="container mx-auto px-4 py-8">
    <h2 class="text-2xl font-bold text-center mb-2">EverDownload - Download Videos</h2>
    <!-- ... (form and result container) ... -->
  </div>
  <script src="https://unpkg.com/htmx.org@1.9.6"></script>
  <script src="https://unpkg.com/alpinejs@3.12.0/dist/cdn.min.js"></script>
</body>
</html>
```

#### `POST /submit`
Processes a submitted video URL, fetches its metadata, and renders dynamic HTML content for quality selection.
**Description**: Takes a video URL, uses `yt-dlp` to extract metadata (title, author, thumbnail, available formats), caches this data, and returns an HTML snippet with download options.
**Request**: `application/x-www-form-urlencoded`
```
videoURL=<encoded_video_url>
```
**Example Request Body**:
```
videoURL=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3DDW-pQ2h9_y0
```
**Response**: `text/html`
Returns an HTML fragment (`div`) containing video details, a thumbnail, title, author, a dropdown (`select`) for quality selection, and a download button that uses the selected format.
**Example Success Response (HTML fragment)**:
```html
<div class="mt-6 mb-20 p-4 rounded-lg shadow-2xl" x-data="{ selectedFormat: '303', pageUrl: 'https://www.youtube.com/watch?v=DW-pQ2h9_y0' }">
    <h3 class="text-lg font-bold mb-4">Video Details</h3>
    <img src="https://i.ytimg.com/vi/DW-pQ2h9_y0/hqdefault.jpg" alt="Video Thumbnail" class="w-full rounded-md mb-4" />
    <p class="text-white mb-2"><strong>Title:</strong> Sample Video Title</p>
    <p class="text-white mb-2"><strong>Author:</strong> Sample Author</p>
    <div class="mt-4">
        <label for="qualitySelect" class="block mb-2">Select Quality</label>
        <select id="qualitySelect" x-model="selectedFormat" class="w-full p-2 bg-neutral-800 text-white rounded-md border">
            <option value="303">1080p</option>
            <option value="248">1080p video only</option>
            <option value="247">720p</option>
            <option value="251">Audio only</option>
        </select>
    </div>
    <a 
        x-bind:href="'/download?url=' + encodeURIComponent(pageUrl) + '&filename=Sample Video Title.mp4&format=' + encodeURIComponent(selectedFormat)" 
        class="block mb-32 w-full mt-4 bg-red-900 text-center text-white p-3 rounded-md hover:bg-blue-600"
        download
    >
        Download Video
    </a>
</div>
```
**Errors**:
- `400 Bad Request`: "Error Parsing Form" or "Invalid or unsupported video URL" if the provided URL is malformed or from an unsupported host.
- `500 Internal Server Error`: "Error fetching video meta data: [details]" if `yt-dlp` fails to retrieve video information.

#### `GET /download`
Initiates the download of a video based on the provided URL and format ID.
**Description**: Streams the video content directly to the client using `yt-dlp` to fetch the specific format.
**Request**: Query Parameters
```
/download?url=<encoded_video_url>&filename=<encoded_filename>&format=<format_id>
```
**Query Parameters**:
- `url` (string, **required**): The original URL of the video.
- `filename` (string, optional): The desired filename for the downloaded video. Defaults to `video.mp4` if not provided.
- `format` (string, **required**): The specific format ID to download, obtained from the `/submit` endpoint's response.
**Example Request**:
```
GET /download?url=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3DDW-pQ2h9_y0&filename=Sample%20Video%20Title.mp4&format=303
```
**Response**: `video/mp4`
The raw video stream data will be sent as the response. The `Content-Disposition` header will be set to prompt a download with the specified filename.
**Errors**:
- `400 Bad Request`: "Missing video page URL" or "Invalid format" if `url` or `format` query parameters are missing or invalid.
- `500 Internal Server Error`: "Failed to download video" if `yt-dlp` encounters an error during the download process.

#### `GET /static/`
Serves static assets like images or CSS files.
**Description**: A file server for assets located in the `static/` directory.
**Request**: Accesses files directly by path (e.g., `/static/animate.png`).
**Response**: Corresponds to the requested file type (e.g., `image/png`, `image/svg+xml`).

## Technologies Used

| Technology         | Description                                        | Link                                                                      |
| :----------------- | :------------------------------------------------- | :------------------------------------------------------------------------ |
| **Go**             | Core backend language                              | [Go Lang](https://golang.org/)                                            |
| **Redis**          | In-memory data store for caching                   | [Redis](https://redis.io/)                                                |
| **yt-dlp**         | Command-line video downloader                      | [yt-dlp GitHub](https://github.com/yt-dlp/yt-dlp)                         |
| **Docker**         | Containerization platform                          | [Docker](https://www.docker.com/)                                         |
| **HTMX**           | HTML-centric AJAX framework                        | [HTMX](https://htmx.org/)                                                 |
| **Alpine.js**      | Lightweight JavaScript framework for interactivity | [Alpine.js](https://alpinejs.dev/)                                        |
| **TailwindCSS**    | Utility-first CSS framework                        | [TailwindCSS](https://tailwindcss.com/)                                   |
| **Go Dotenv**      | Loads environment variables from `.env` file       | [Go Dotenv GitHub](https://github.com/joho/godotenv)                      |
| **Go Redis Client**| Redis client for Go                                | [Go Redis GitHub](https://github.com/redis/go-redis/v9)                   |

## Contributing
Contributions are highly welcome! If you have suggestions for improvements, new features, or bug fixes, please feel free to contribute.

1.  🍴 **Fork the repository**: Start by forking the EverDownload repository to your GitHub account.
2.  🌿 **Create a new branch**: For each new feature or bug fix, create a dedicated branch.
    ```bash
    git checkout -b feature/your-feature-name
    ```
3.  💻 **Implement your changes**: Write your code, ensuring it adheres to the project's coding standards.
4.  🧪 **Test your changes**: Thoroughly test your code to prevent regressions.
5.  ➕ **Commit your changes**: Write clear and concise commit messages.
    ```bash
    git commit -m "feat: Add new feature for X"
    ```
6.  ⬆️ **Push your branch**: Push your branch to your forked repository.
    ```bash
    git push origin feature/your-feature-name
    ```
7.  📝 **Open a Pull Request**: Submit a pull request to the `main` branch of the original repository. Provide a detailed description of your changes.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Author Info
Developed and maintained by:
**Samuel Tuoyo**
*   LinkedIn: [https://www.linkedin.com/in/samuel-tuoyo-8568b62b6/](https://www.linkedin.com/in/samuel-tuoyo-8568b62b6/)
*   X(formely Twitter): [https://x.com/TuoyoS26091](https://x.com/TuoyoS26091)

---
[![Built with Go](https://img.shields.io/badge/Go-1.24.1-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Powered by Redis](https://img.shields.io/badge/Redis-FF4438?style=flat-square&logo=redis&logoColor=white)](https://redis.io/)
[![Docker Compose](https://img.shields.io/badge/Docker%20Compose-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square)](https://github.com/samueltuoyo15/EverDownload/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)
