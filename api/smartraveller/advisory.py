import requests
import time
from bs4 import BeautifulSoup

advisory_level ={
    "Exercise normal safety precautions": 1,
    "Exercise a high degree of caution": 2,
    "Reconsider your need to travel": 3,
    "Do not travel": 4
}

# Smartraveller uses non-standardised names for countries - see README.md
# Temporary permanent solution... 
special_mappings = {
    "congo": "democratic-republic-congo",
    "iran,-islamic-republic-of": "iran",
    "israel": "middle-east/israel-and-palestinian-territories",
    "kyrgyzstan": "kyrgyz-republic",
    "korea,-democratic-people's-republic-of": "north-korea-democratic-peoples-republic-korea",
    "korea,-republic-of": "south-korea-republic-korea",
    "slovenia": "europe/slovenia",
    "slovakia": "europe/slovakia",
    "taiwan,-province-of-china": "taiwan",
    "united-states": "americas/united-states-america",
    "viet-nam": "vietnam",
    "russian-federation": "russia",
    "greenland": "denmark",
    "lao-people's-democratic-republic": "laos",
    "kyrgyzstan": "kyrgyz-republic"
}

def get_overall_advisory(country: str) -> dict:
    if country in special_mappings:
        country_query = special_mappings[country]
    else:
        country_query = country
    print(country_query)
    
    url = f'https://www.smartraveller.gov.au/destinations/{country_query}'
    response = requests.get(url, timeout=2)
    html = response.content
    
    site = BeautifulSoup(html, 'html.parser')
    
    advisory_block = site.findAll('div', { 'class': 'views-field views-field-field-overall-advice-level'})
    print(advisory_block)
    
    if not advisory_block: 
        return None
    
    advisory = advisory_block[0].find('strong').getText()
    
    return {
        "advisory": advisory,
        "level": advisory_level[advisory],
        "page_url": url
    }