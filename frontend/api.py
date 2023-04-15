import requests
import logging
import os

from urllib.parse import urljoin

HOST = os.environ.get("DOMAIN", "http://localhost:5000/")
API_ROOT = os.environ.get("APIHOST", "http://localhost:8080/")
API_BASE = urljoin(API_ROOT, "api/v1/")


def load_secrets(token):
    logging.debug(f"Loading secrets")

    r = requests.get(urljoin(API_BASE, f"get_urls_list/{token}"))

    if r.status_code == 200:
        if r.json()["urls"] is None:
            return []
        return r.json()["urls"]


def get_url_info(secret_key):
    r = requests.get(urljoin(API_BASE, f"admin/{secret_key}"))

    if r.status_code == 200:
        data = r.json()
        data["short_url"] = urljoin(HOST, data["short_url"])
        return data


def del_url(secret_key):
    try:
        r = requests.delete(urljoin(API_BASE, f"admin/{secret_key}"))

        if r.status_code == 204:
            return True

        if r.text:
            if "error" in dict(r.json()):
                return r.json()["error"]

            return r.text

        return r.reason
    except Exception as ex:
        return str(ex)


def make_short_url(url, short_url, ttl, units, token):
    try:
        r = requests.post(
            urljoin(API_BASE, "make_shorter"),
            json={
                "url": url,
                "vip_key": short_url,
                "token": token,
                "time_to_live": int(ttl),
                "time_to_live_unit": units,
            }
        )

        if r.status_code == 200:
            return True

        if r.text:
            if "error" in dict(r.json()):
                return r.json()["error"]

            return r.text

        return r.reason
    except Exception as ex:
        return str(ex)


def make_basic_short_url(url):
    try:
        r = requests.post(
            urljoin(API_BASE, "make_shorter"),
            json={
                "url": url,
            }
        )

        if r.status_code == 200:
            return get_url_info(r.json()["secret_key"])

        if r.text:
            if "error" in dict(r.json()):
                raise ValueError(r.json()["error"])

            raise ValueError(r.text)

        raise ValueError(r.reason)
    except Exception as ex:
        raise ValueError(str(ex))


def signup_user(email, login, password):
    logging.debug(f"Sign up: {email}, {login}, {password}")

    r = requests.post(
        urljoin(API_BASE, "registration"),
        json={
            "email": email,
            "login": login,
            "password": password
        }
    )

    if r.status_code == 200:
        return r.json()["token"]

    raise ValueError(r.json()["error"])


def login_user(login, password):
    logging.debug(f"Login: {login}, {password}")

    r = requests.post(
        urljoin(API_BASE, "auth"),
        json={
            "login": login,
            "password": password
        }
    )

    if r.status_code == 200:
        return r.json()["token"]

    raise ValueError(r.json()["error"])


def get_full_url(short):
    r = requests.get(urljoin(API_ROOT, short))
    if r.status_code == 200:
        return r.json()["long_url"]

    return HOST
