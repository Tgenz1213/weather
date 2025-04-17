import { useState, useCallback } from "react";
import { getWeather } from "@services/api";

interface WeatherData {
  startTime: string;
  endTime: string;
  detailedForecast: string;
}

type FetchWeatherFunction = (
  street: string,
  zip: string
) => Promise<WeatherData>;

export const useWeather = (
  street: string,
  zip: string,
  fetchWeatherFn: FetchWeatherFunction = getWeather
) => {
  const [weatherData, setWeatherData] = useState<WeatherData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const fetchWeather = useCallback(async () => {
    setLoading(true);
    setError(null);
    setWeatherData(null);

    try {
      console.log("attempting to fetch Weather..."); // DELETEME
      const data = await fetchWeatherFn(street, zip);
      if (data && data.startTime && data.endTime && data.detailedForecast) {
        setWeatherData(data);
      } else {
        setError("Invalid data structure received from API");
        console.error("Invalid data structure:", data);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Unknown error";
      setError("An error occurred while fetching weather data");
      console.error("API error:", errorMessage, err);
    } finally {
      setLoading(false);
    }
  }, [street, zip, fetchWeatherFn]);

  return { weatherData, error, loading, fetchWeather };
};
