from flask import Blueprint

from flask import render_template
from flask import request, redirect, url_for
from urllib.parse import urljoin

from api import load_secrets, get_url_info, make_short_url, make_basic_short_url, del_url

main = Blueprint('main', __name__)


@main.route("/")
def index(url_data=None):
    return render_template(
        "index.html",
        error=request.args.get("error"),
        url_data=url_data
    )


@main.route("/custom")
def custom():
    is_auth = request.cookies.get("token") is not None

    secrets = load_secrets(request.cookies.get("token"))
    if secrets is None:
        is_auth = False

    if not is_auth:
        return redirect(url_for("auth.login"))

    return render_template(
        "custom.html",
        links=list(zip(secrets, [get_url_info(secret) for secret in secrets])),
        error=request.args.get("error"),
    )


@main.route("/make_url", methods=["POST"])
def make_url():
    url = request.form.get('url')
    short_url = request.form.get('short-url')
    ttl = request.form.get('ttl')
    units = request.form.get('units')

    r = make_short_url(url, short_url, ttl, units, request.cookies.get("token"))

    if r is True:
        return redirect(url_for("main.custom"))

    return redirect(f'{url_for("main.custom")}?error={r}')


@main.route("/del_url", methods=["POST"])
def del_url_post():
    secret = request.form.get('del')

    r = del_url(secret)

    if r is True:
        return redirect(url_for("main.custom"))

    return redirect(f'{url_for("main.custom")}?error={r}')


@main.route("/", methods=["POST"])
def make_basic_url():
    url = request.form.get('url')

    try:
        return index(make_basic_short_url(url))
    except ValueError as ex:
        return redirect(f"/?error={ex}")
