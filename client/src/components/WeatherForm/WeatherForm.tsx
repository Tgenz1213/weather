import React, { useState } from "react"
import { useWeather } from "@hooks/useWeather"

const WeatherForm: React.FC = () => {
  const [street, setStreet] = useState<string>("")
  const [zip, setZip] = useState<string>("")
  const { weatherData, error, loading, fetchWeather } = useWeather(street, zip)

  const handleSubmit = (e: React.FormEvent) => {
    console.log("calling handleSubmit")
    e.preventDefault()
    fetchWeather()
  }

  return (
    <div className="weather-form-container">
      <h1>Weather Forecast</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Street Address: </label>
          <input
            type="text"
            value={street}
            onChange={(e) => setStreet(e.target.value)}
            maxLength={100}
            required
          />
        </div>
        <div>
          <label>Zip Code: </label>
          <input
            type="text"
            value={zip}
            onChange={(e) => setZip(e.target.value)}
            size={5}
            required
          />
        </div>
        <button type="submit" disabled={loading}>
          Get Forecast
        </button>
      </form>

      {loading && <div>Loading...</div>}
      {error && <div className="error">{error}</div>}
      {weatherData && (
        <div className="weather-data">
          <h2>Weather Details</h2>
          <p>
            <strong>Start Time:</strong> {weatherData.startTime}
          </p>
          <p>
            <strong>End Time:</strong> {weatherData.endTime}
          </p>
          <p>
            <strong>Detailed Forecast:</strong> {weatherData.detailedForecast}
          </p>
        </div>
      )}
    </div>
  )
}

export default WeatherForm
