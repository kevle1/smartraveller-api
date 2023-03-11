import requests
import json
import pycountry
from datetime import datetime
import time

all_codes = [country.alpha_2 for country in pycountry.countries]

if __name__ == "__main__":
    advisories = {}
    for code in all_codes:
        print(f"Getting advisory for {code}")
        response = requests.get(f"https://smartraveller.api.kevle.xyz/advisory?country={code}")
        print(response)
        
        if response.status_code == 200:
            advisories[code] = response.json()
            
        time.sleep(0.5) 

    with open(f"smartraveller-compounded.json", "w") as outfile:
        json.dump({
            "last_updated": datetime.now().isoformat(),
            "advisories": advisories
        }, outfile, indent=4)
