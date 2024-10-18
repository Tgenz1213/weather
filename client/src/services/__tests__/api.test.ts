import { getWeather, WeatherData } from "../api";

describe("getWeather", () => {
  const mockWeatherData: WeatherData = {
    startTime: "2024-08-15T00:00:00Z",
    endTime: "2024-08-15T12:00:00Z",
    detailedForecast: "Sunny with a chance of rain.",
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("returns weather data on successful API call", async () => {
    global.fetch = jest.fn().mockResolvedValueOnce({
      ok: true,
      json: async () => mockWeatherData,
    });

    const data = await getWeather("123 Main St", "12345");
    expect(data).toEqual(mockWeatherData);
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/weather")
    );
  });

  it("throws an error on HTTP error response", async () => {
    global.fetch = jest.fn().mockResolvedValueOnce({
      ok: false,
      status: 404,
    });

    await expect(getWeather("123 Main St", "12345")).rejects.toThrow(
      "Failed to fetch weather data"
    );
  });

  it("throws an error on network error", async () => {
    global.fetch = jest.fn().mockRejectedValueOnce(new Error("Network error"));

    await expect(getWeather("123 Main St", "12345")).rejects.toThrow(
      "Failed to fetch weather data"
    );
  });

  it("handles invalid data structure gracefully", async () => {
    global.fetch = jest.fn().mockResolvedValueOnce({
      ok: true,
      json: async () => ({}),
    });

    await expect(getWeather("123 Main St", "12345")).resolves.toEqual({});
  });

  it("handles partial weather data gracefully", async () => {
    const partialData = { startTime: "2024-08-15T00:00:00Z" };
    global.fetch = jest.fn().mockResolvedValueOnce({
      ok: true,
      json: async () => partialData,
    });

    const data = await getWeather("123 Main St", "12345");
    expect(data).toEqual(partialData);
  });
});
