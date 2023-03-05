# Smartraveller API

A simple API for travel advisories published by the Australian Government on [Smartraveller](https://www.smartraveller.gov.au/).

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

## Endpoint(s)

### [GET /advisory](https://smartraveller.api.kevle.xyz?country=es)

Fetches an advice summary for a country. 

Responses may be cached for up to 1 hour. 

#### Query Parameters 

- country (required)
  - Note: Fuzzy matches for country - E.g. "Spain" or "ES" is accepted. However it is recommend to use Alpha 3 country codes. 

#### Example Response 

```
{
  "advisory": "Exercise normal safety precautions",
  "level": 1,
  "country": "Spain",
  "alpha_2": "ES",
  "official_name": "Kingdom of Spain"
}
```

## Copyright

Smartraveller material accessed is provided under the latest [Creative Commons Attribution licence](https://creativecommons.org/licenses/by/4.0/). 