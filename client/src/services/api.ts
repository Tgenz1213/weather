// src/services/api.ts
import axios from 'axios';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
});

export interface WeatherData {
    startTime: string;
    endTime: string;
    detailedForecast: string;
}

export const getWeather = async (street: string, zip: string): Promise<WeatherData> => {
    try {
        const response = await api.get<WeatherData>('/weather', {
            params: { street, zip }
        });
        return response.data;
    } catch (error) {
        console.error('Error fetching weather data:', error);
        throw new Error('Failed to fetch weather data');
    }
};
