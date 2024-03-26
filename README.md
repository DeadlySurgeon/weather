# Weather Service Coding Challenge

Coding challenge for a basic weather service API caller. Currently uses REST and
is exposed over at `/weather/at` with query parameters `lat` and `lon`.

## Environmental Variable Configuration

This program is configured based off of environmental variables as described
below. They can be set ahead of time or in a .env file, which an `example.env`
is included.

| ENV              | Desc                              |
| ---------------- | --------------------------------- |
| Net              | Network based configuration       |
| NET_BIND         | Network address binding           |
| NET_TLSCERT      | TLS Certificate Pem file location |
| NET_TLSKEY       | TLS Key Pem file location         |
|                  |                                   |
| Weather          | Weather service configuration     |
| WEATHER_ENDPOINT | Endpoint to hit                   |
| WEATHER_APIKEY   | API Key for openweathermap.org    |

## API

**Endpoint**: `{SERVICE}/weather/at`  
**Query Params:**:

| Key   | Type    | Desc                                       |
| ----- | ------- | ------------------------------------------ |
| `lat` | float64 | Latitude you want to check the weather of  |
| `lon` | float64 | Longitude you want to check the weather of |

**Response**:

```jsonc
{
  "condition": "[snow|rain|thunder]", // This is up to the openweatherapi condition field
  "temperature": "[freezing|cold|moderate|hot|burning]",
  "temperature_raw": 0 // float32 in metric
}
```
