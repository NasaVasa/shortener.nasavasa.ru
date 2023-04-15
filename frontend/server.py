from urllib.parse import quote_plus

import requests
from flask import Flask, redirect
from markupsafe import Markup

from auth import auth as auth_blueprint
from index import main as main_blueprint
from api import get_full_url
import pytz
from datetime import datetime

app = Flask(__name__)

app.register_blueprint(auth_blueprint)
app.register_blueprint(main_blueprint)


# -------- format time at template ---------
@app.template_filter("datetimefilter")
def datetimefilter(val):
    dt = datetime.strptime(val, "%Y-%m-%dT%H:%M:%S.%fZ")
    dt = dt.replace(tzinfo=pytz.UTC)
    return dt.astimezone(pytz.timezone("Europe/Moscow")).strftime('%Y-%m-%d %H:%M:%S')


@app.template_filter("url_encode")
def url_encode(s):
    if type(s) == "Markup":
        s = s.unescape()
    return Markup(quote_plus(s.encode()))


@app.route('/<short>')
def fallback(short):
    return redirect(get_full_url(short))


if __name__ == '__main__':
    app.run()
