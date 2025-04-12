# Smartraveller API

A simple API for travel advisories published by the Australian DFAT [Smartraveller](https://www.smartraveller.gov.au/).

**This is not an official API. Visit the [Smartraveller website](https://www.smartraveller.gov.au) for the latest information.**

Available at https://smartraveller.kevle.xyz/api/

## Running Locally

```sh
go run .
```

## Deployment

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fkevle1%2Fsmartraveller-api)

## Endpoints

### [GET /api/advisory](https://smartraveller.kevle.xyz/api/advisory)

Get an advice summary for a country. May be cached for up to half an hour.

#### Query Parameters

- country (required)
  - ISO [alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) country codes

#### Example Response

https://smartraveller.kevle.xyz/api/advisory?country=es

```json
{
  "lastFetched": "2025-04-12T09:10:12Z",
  "advisory": {
    "country": {
      "name": "Spain",
      "alpha2": "ES"
    },
    "advice": "Exercise normal safety precautions",
    "latestUpdate": "We continue to advise exercise normal safety precautions in Spain. We advise: Exercise normal safety precautions in Spain.",
    "level": 2,
    "published": "2025-03-27T12:00:00Z",
    "pageUrl": "https://www.smartraveller.gov.au/destinations/europe/spain"
  }
}
```

### [GET /api/advisories](https://smartraveller.kevle.xyz/api/advisories)

Get all available Smartraveller advisories. May be cached for up to 1 hour.

#### Example Response

```json
{
  "lastFetched": "2025-04-12T09:11:06Z",
  "advisories": [
    {
      "country": {
        "name": "France",
        "alpha2": "FR"
      },
      "advice": "Exercise a high degree of caution",
      "latestUpdate": "We've reviewed our travel advice for France and continue to advise exercise a high degree of caution due to the threat of terrorism. France's national terrorist alert warning remains at the highest level. Expect high-level security nationwide (see 'Safety'). If you plan to travel to France to commemorate Anzac Day, understand the risks and plan ahead (see 'Travel'). We advise: Exercise a high degree of caution in France due to the threat of terrorism.",
      "level": 3,
      "published": "2025-03-19T12:00:00Z",
      "pageUrl": "https://www.smartraveller.gov.au/destinations/europe/france"
    },
    ...
  ],
  "length": 178
}
```

## Data

Smartraveller publishes advisories under non-standard names for countries and does not seem to utilise any standard country codes. We attempt to look up the country name and fallback to custom mapping of the page slug.

This method of parsing countries may produce inconsistent or inaccurate results.

If a country can't be found it may be due to either a failed lookup, or an advisory not being published for that country.

## Copyright

Smartraveller material is provided under the latest Creative Commons Attribution licence.

Information: https://www.smartraveller.gov.au/copyright