import json
import logging
import time
from datetime import datetime, timezone

import pycountry
import requests
from destination import _fetch_destinations

BASE_URL = "https://smartraveller.api.kevle.xyz"
OUTPUT_PREFIX = "data"


def _save_destinations():
    response = requests.get(f"{BASE_URL}/destinations?refetch=true")
    if response.status_code == 200:
        with open(f"{OUTPUT_PREFIX}/destinations.json", "w") as outfile:
            json.dump(
                {
                    "last_updated": datetime.now(timezone.utc).isoformat(),
                    "destinations": response.json(),
                },
                outfile,
                indent=4,
            )
        logging.info("Saved destinations")


def _aggregate_advisories():
    advisories = []
    for country in pycountry.countries:
        logging.info(f"Getting advisory for {country.alpha_2}")

        response = requests.get(f"{BASE_URL}/advisory?country={country.alpha_2}")
        if response.status_code == 200:
            advisories.append(response.json())
        time.sleep(0.5)

    with open(f"{OUTPUT_PREFIX}/advisories.json", "w") as outfile:
        json.dump(
            {
                "last_updated": datetime.now(timezone.utc).isoformat(),
                "advisories": advisories,
            },
            outfile,
            indent=4,
        )
        logging.info("Saved advisories")


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)

    _save_destinations()
    _aggregate_advisories()
