import json

from flask import Flask, Response, request

from smartraveller.advisory import get_advisory

app = Flask(__name__, static_url_path="")


@app.route("/")
def index():
    return '<h1>Unofficial Smartraveller API</h1> \
            <p>Smartraveller material is provided under the latest Creative Commons Attribution licence.<p/> \
            <a href="https://www.smartraveller.gov.au">Smartraveller</a>'


@app.route("/advisory")
def advisory():
    iso_alpha_2 = request.args.get("country")
    if not iso_alpha_2:
        return _return_error("Missing country code", 400)

    if len(iso_alpha_2) != 2 and not iso_alpha_2.isalpha():
        return _return_error("Invalid country code", 400)

    advisory = get_advisory(iso_alpha_2.upper())

    if advisory is None:
        return _return_error(
            "Smartraveller does not have published advisory for selected country", 404
        )

    return _return_success(advisory.model_dump())


@app.route("/advisories")
def advisories():
    with open("data/advisories.json", "r") as advisories_file:
        json_advisories = json.load(advisories_file)
    return _return_success(json_advisories)


@app.route("/destinations")
def destinations():
    with open("data/destinations.json", "r") as destinations_file:
        destinations = json.load(destinations_file)
    return _return_success(destinations)


def _return_success(data: dict, cache_duration_min: int = 60) -> Response:
    response = Response(json.dumps(data))

    # https://vercel.com/docs/edge-network/caching
    response.headers["Cache-Control"] = f"s-maxage={cache_duration_min*60}"
    response.headers["Content-Type"] = "application/json"
    response.headers["Access-Control-Allow-Origin"] = "*"

    return response


def _return_error(message: str, status: int) -> Response:
    return Response(
        json.dumps({"error": message}), status=status, mimetype="application/json"
    )
