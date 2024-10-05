# Smartraveller API

A simple API for travel advisories published by the Australian DFAT [Smartraveller](https://www.smartraveller.gov.au/).

**This is not an official API. Visit the Smartraveller website for the latest information.**

Available at https://smartraveller.api.kevle.xyz/

## Locally running
```
pip install -r requirements.txt
python wsgi.py
```

## Deployment

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fkevle1%2Fsmartraveller-api)

Or configure a WSGI production server

## Endpoints

### [GET /advisory](https://smartraveller.api.kevle.xyz/advisory)

Get an advice summary for a country. May be cached for up to 1 hour.

#### Query Parameters

- country (required)
  - ISO [alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) country codes

#### Example Response

https://smartraveller.api.kevle.xyz/advisory?country=es

```json
{
  "country": {
    "name": "Spain",
    "alpha_2": "ES"
  },
  "advisory": "Exercise normal safety precautions",
  "level": 1,
  "page_url": "https://www.smartraveller.gov.au/destinations/europe/spain"
}
```

### [GET /advisories](https://smartraveller.api.kevle.xyz/advisories)

Get all available Smartraveller advisories. This is updated every 3 hours.

#### Example Response

```json
{
    "last_updated": "2024-10-05T04:55:06.260067+00:00",
    "advisories": [
        {
          "country": {
            "name": "Spain",
            "alpha_2": "ES"
          },
          "advisory": "Exercise normal safety precautions",
          "level": 1,
          "page_url": "https://www.smartraveller.gov.au/destinations/europe/spain"
        },
    ]
}
```

### [GET /destinations](https://smartraveller.api.kevle.xyz/destinations)

Get all available Smartraveller destinations. This is updated every 3 hours.

#### Example Response

```json
{
    "last_updated": "2024-10-05T04:55:06.260067+00:00",
    "destinations": {
      "ES": {
        "country": {
          "name": "Spain",
          "alpha_2": "ES"
        },
        "page_url": "/destinations/europe/spain"
      },
    }
}
```

## Data

Smartraveller publishes advisories under non-standard names for countries and does not seem to utilise any standard country codes. Countries have been mapped using a fuzzy search for the country based on the URL slug fetched from: https://www.smartraveller.gov.au/api/publishedpages.

This method of parsing countries may produce inconsistent or inaccurate results.

If a country can't be found it may be due to a failed lookup, or an advisory not being published for that country.

## Copyright

Smartraveller material is provided under the latest Creative Commons Attribution licence.

Information: https://www.smartraveller.gov.au/copyright