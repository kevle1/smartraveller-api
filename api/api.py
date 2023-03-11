import json
import pycountry
from flask import Flask, request, Response
from api.smartraveller.advisory import get_overall_advisory
import logging

app = Flask(__name__, static_url_path="")
@app.route("/")
def index(): 
    return '<h1>Smartraveller API</h1> \
            <p>Smartraveller material accessed is provided under the latest Creative Commons Attribution licence.<p/> \
            <a href="https://www.smartraveller.gov.au/">Smartraveller</a>'

@app.route("/advisory")
def advisory():
    country_query = request.args.get("country")
    if country_query:
        try:
            country = pycountry.countries.lookup(country_query)
        except Exception as e:
            logging.error(e)
            country = None
        
        if country is not None:
            country_query = country.name.lower().replace(" ", "-")
            logging.debug(f"Querying for country: {country_query}")
            
            advisory = get_overall_advisory(country_query)
            
            if advisory is None:
                return "Smartraveller does not have published advisory for selected country", 404
            
            advisory["country"] = country.name
            advisory["alpha_2"] = country.alpha_2
            
            if hasattr(country, "official_name"):
                advisory["official_name"] = country.official_name
            
            response = Response(json.dumps(advisory))
            response.headers["Cache-Control"] = "s-maxage=3600" # Vercel cache for 1 hour
            response.headers["Content-Type"] = "application/json"
            
            return response
        else:
            return "Invalid country", 404
    else:
        return "Missing country query parameter", 400

@app.route("/advisories")
def advisories():
    advisories = None
    with open("smartraveller-compounded.json", "r") as advisories:
        json_advisories = json.load(advisories)
    
    response = Response(json.dumps(json_advisories))
    response.headers["Cache-Control"] = "s-maxage=86400" # Vercel cache for 1 day
    response.headers["Content-Type"] = "application/json"
    
    return response