import requests
import time
from bs4 import BeautifulSoup

advisory_level ={
    "Exercise normal safety precautions": 1,
    "Exercise a high degree of caution": 2,
    "Reconsider your need to travel": 3,
    "Do not travel": 4
}

def get_overall_advisory() -> dict:
    country = "austria"
    
    html = requests.get(f'https://www.smartraveller.gov.au/destinations/{country}').content
    site = BeautifulSoup(html, 'html.parser')
    
    advisory_block = site.findAll('div', { 'class': 'views-field views-field-field-overall-advice-level'})
    advisory = advisory_block[0].find('strong').getText()
    
    return {
        "advisory": advisory,
        "level": advisory_level[advisory]
    }

if __name__ == "__main__":
    print(get_overall_advisory())