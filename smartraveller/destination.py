import json
import logging

import pycountry
import requests
from pydantic import BaseModel


class Country(BaseModel):
    name: str
    alpha_2: str


class SmartravellerDestination(BaseModel):
    country: Country
    page_url: str


# Some countries can't be found with fuzzy search so we have to maintain a manual mapping
MANUAL_MAPPINGS: dict[str, str | list[str]] = {
    "cote-divoire-ivory-coast": "CI",
    "democratic-republic-congo": "CD",
    "federated-states-micronesia": "FM",
    "israel-and-occupied-palestinian-territories": ["PS", "IL"],
    "macau": "MO",
    "north-korea-democratic-peoples-republic-korea": "KP",
    "south-korea-republic-korea": "KR",
    "timor-leste": "TL",
    "united-states-america": "US",
}


def get_destinations(refetch: bool = False) -> dict[str, SmartravellerDestination]:
    if not refetch:
        try:
            with open("data/destinations.json", "r") as destinations:
                destinations = json.load(destinations)["destinations"]
        except FileNotFoundError:
            destinations = _fetch_destinations()
    else:
        destinations = _fetch_destinations()

    return {
        country: SmartravellerDestination(**destination)
        if isinstance(destination, dict)
        else destination
        for country, destination in destinations.items()
    }


def _fetch_destinations() -> dict[str, SmartravellerDestination]:
    response = requests.get("https://www.smartraveller.gov.au/api/publishedpages")
    if response.status_code != 200:
        logging.error("Failed to get destinations")
        return {}

    destinations = {}

    for destination in response.json():
        url_parts = destination["url"].split("/")

        if (
            destination["pageType"] == "location"
            and destination["url"].startswith("/destination")
            and len(url_parts) == 4
        ):
            logging.info(f"Found destination {destination['url']}")

            url = destination["url"]
            country_slug = url_parts[3]
            country_name = country_slug.replace("-", " ")

            try:
                if country_slug in MANUAL_MAPPINGS:
                    country_mapping = MANUAL_MAPPINGS[country_slug]
                    if isinstance(country_mapping, list):
                        for country_code in country_mapping:
                            country = pycountry.countries.get(alpha_2=country_code)
                            destinations[country.alpha_2] = SmartravellerDestination(
                                country=Country(
                                    name=country.name, alpha_2=country.alpha_2
                                ),
                                page_url=url,
                            )

                        continue
                    else:
                        country = pycountry.countries.get(alpha_2=country_mapping)
                else:
                    country = pycountry.countries.search_fuzzy(country_name)[0]

                destinations[country.alpha_2] = SmartravellerDestination(
                    country=Country(name=country.name, alpha_2=country.alpha_2),
                    page_url=url,
                )
            except LookupError:
                logging.warning(f"Could not find country {country_slug}")

    return destinations


if __name__ == "__main__":
    destinations = _fetch_destinations()
    print(destinations)
