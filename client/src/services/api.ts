export interface WeatherData {
  startTime: string
  endTime: string
  detailedForecast: string
}

export const getWeather = async (
  street: string,
  zip: string,
): Promise<WeatherData> => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL
  const url = new URL("/weather", baseUrl)

  url.searchParams.append("street", street)
  url.searchParams.append("zip", zip)

  try {
    const response = await fetch(url.toString())

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }

    const data: WeatherData = await response.json()
    return data
  } catch (error) {
    console.error("Error fetching weather data:", error)
    throw new Error("Failed to fetch weather data")
  }
}
