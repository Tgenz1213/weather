# Weather Forecast Application

This application consists of a Go backend API and a React frontend, containerized using Docker for easy deployment.

## Architecture

### Backend (Go)
- Located in the `server` directory
- Provides a RESTful API endpoint for weather forecasts
- Uses Redis for caching to improve performance

### Frontend (React)
- Located in the `client` directory
- Built with React and TypeScript
- Communicates with the backend API to fetch and display weather data

## How to Run

1. Clone the project to a directory of your choice.

```bash
cd C:/your/directory
git clone https://github.com/Tgenz1213/weather
```

2. Open a terminal or command prompt and navigate to the extracted directory.

```bash
cd C:/your/directory/weather
```

3. Rename or duplicate .env.example to ".env.development"

```bash
cp .env.example .env.development
```

4. Ensure the mode in ./.env matches the .env.[MODE] that you wish you use.

5. Ensure Docker and Docker Compose are installed on your system.

6. Run the following command to start the application:

```bash
docker compose up
```

This command will build and start both the frontend and backend services.

7. Once the containers are running, open a web browser and navigate to:

[http://localhost:3000/](http://localhost:3000/)


8. You should now see the Weather Forecast application. Enter a street address (including the building number) and zip code to get the weather forecast.

## Troubleshooting

- If you encounter any issues, ensure that ports 3000 and 8080 are not in use by other applications.
- Check the Docker logs for any error messages:

```bash
docker compose logs
```

For further assistance, please open an issue on the [GitHub repository](https://github.com/Tgenz1213/weather/issues) or contact Timothy Genz at timothy.genz@yahoo.com.
