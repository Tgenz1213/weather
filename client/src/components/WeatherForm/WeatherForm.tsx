import React, { useState, useEffect } from 'react';
import { getWeather } from '../../services/api';

interface WeatherData {
    startTime: string;
    endTime: string;
    detailedForecast: string;
}

const WeatherForm: React.FC = () => {
    const [street, setStreet] = useState<string>('');
    const [zip, setZip] = useState<string>('');
    const [weatherData, setWeatherData] = useState<WeatherData | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);

    useEffect(() => {
        console.log('Current weatherData:', weatherData);
    }, [weatherData]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setWeatherData(null);

        try {
            console.log('Fetching weather data for:', { street, zip });
            const data = await getWeather(street, zip);
            console.log('API response:', data);

            if (data && data.startTime && data.endTime && data.detailedForecast) {
                setWeatherData(data);
                setError(null);
            } else {
                setError('Invalid data structure received from API');
                console.log(data);
                setWeatherData(null);
            }
        } catch (err) {
            console.error('API error:', err);
            setError('An error occurred while fetching weather data');
        } finally {
            setLoading(false);
        }
    };

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
                <button type="submit" disabled={loading}>Get Forecast</button>
            </form>

            {loading && <div>Loading...</div>}
            {error && <div className="error">{error}</div>}
            {weatherData && (
                <div className="weather-data">
                    <h2>Weather Details</h2>
                    <p><strong>Start Time:</strong> {weatherData.startTime}</p>
                    <p><strong>End Time:</strong> {weatherData.endTime}</p>
                    <p><strong>Detailed Forecast:</strong> {weatherData.detailedForecast}</p>
                </div>
            )}
        </div>
    );
};

export default WeatherForm;
