import requests
from bs4 import BeautifulSoup
from pydantic import BaseModel

from .destination import Country, get_destinations


class SmartravellerAdvisory(BaseModel):
    country: Country
    advisory: str
    level: int
    page_url: str


ADVISORY_LEVELS = {
    "Exercise normal safety precautions": 1,
    "Exercise a high degree of caution": 2,
    "Reconsider your need to travel": 3,
    "Do not travel": 4,
}


def get_advisory(iso_alpha_2: str) -> SmartravellerAdvisory | None:
    destinations = get_destinations()

    if iso_alpha_2 in destinations:
        destination = destinations[iso_alpha_2]

        url = f"https://www.smartraveller.gov.au{destination.page_url}"
    else:
        return None

    response = requests.get(url, timeout=2)
    html = response.content

    site = BeautifulSoup(html, "html.parser")

    advice = site.findAll(
        "div", {"class": "views-field views-field-field-overall-advice-level"}
    )

    if not advice or len(advice) == 0:
        return None

    advisory = advice[0].find("strong").getText()

    return SmartravellerAdvisory(
        country=Country(
            name=destination.country.name, alpha_2=destination.country.alpha_2
        ),
        advisory=advisory,
        level=ADVISORY_LEVELS[advisory],
        page_url=url,
    )
