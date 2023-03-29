# Smartraveller API

A simple API for travel advisories published by the Australian DFAT on [Smartraveller](https://www.smartraveller.gov.au/).

**This is not an official API. Visit the Smartraveller website for the latest information.**

## Locally running
```
pip install -r requirements.txt
python wsgi.py
```

## Deployment 

Available at https://smartraveller.api.kevle.xyz/

OR

[![Deploy to Vercel](https://camo.githubusercontent.com/f209ca5cc3af7dd930b6bfc55b3d7b6a5fde1aff/68747470733a2f2f76657263656c2e636f6d2f627574746f6e)](https://vercel.com/import/project?template=https://github.com/kevinle-1/smartraveller-api)

Note: Set Vercel function region to `Sydney, Australia (Southeast) - syd1` as requests may be blocked from certain regions. 

## Endpoint(s)

### [GET /advisory](https://smartraveller.api.kevle.xyz/advisory)

Fetches an advice summary for a country. 

Responses may be cached for up to 1 hour. 

#### Query Parameters 

- country (required)
  - Note: Fuzzy matches for country - E.g. "Spain" or "ES" is accepted. However it is recommend to use [ISO 3166-1 alpha-2 country codes](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2). 

#### Example Response 

https://smartraveller.api.kevle.xyz/advisory?country=es

```
{
  "advisory": "Exercise normal safety precautions",
  "level": 1,
  "country": "Spain",
  "alpha_2": "ES",
  "official_name": "Kingdom of Spain"
}
```

### [GET /advisories](https://smartraveller.api.kevle.xyz/advisories)

Fetches advice summaries for all countries available on Smartraveller. 

Responses may be cached for up to 1 day. 

## Known Problems

Smartraveller uses non-standardised and inconsistent names for countries and does not utilise any standard country codes such as the [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) standard. 

Mappings for certain countries have been temporarily added, but may be inconsistent. 

Need to explore if it is possible to build a mapping programatically using the available ["PublishedPages" API](https://www.smartraveller.gov.au/api/publishedpages). 

## Copyright

Smartraveller material accessed is provided under the latest [Creative Commons Attribution licence](https://creativecommons.org/licenses/by/4.0/). 