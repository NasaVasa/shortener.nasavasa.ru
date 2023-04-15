from flask import Blueprint

from flask import render_template, url_for
from flask import request, redirect, make_response

from api import load_secrets, signup_user, login_user

auth = Blueprint('auth', __name__)


@auth.route("/login")
def login():
    if request.cookies.get("token") is not None:
        return redirect("/")

    secrets = load_secrets(request.cookies.get("token"))
    if secrets is not None:
        return redirect("/")

    return render_template(
        "login.html",
    )


@auth.route('/login', methods=['POST'])
def login_post():
    login = request.form.get('login')
    password = request.form.get('password')

    try:
        resp = redirect(url_for('main.custom'))
        resp.set_cookie("token", login_user(login, password))

        return resp
    except ValueError as ex:
        return render_template(
            "login.html",
            error=ex
        )


@auth.route('/signup')
def signup():
    return render_template(
        "signup.html",
    )


@auth.route('/signup', methods=['POST'])
def signup_post():
    email = request.form.get('email')
    login = request.form.get('login')
    password = request.form.get('password')

    try:
        token = signup_user(email, login, password)

        resp = redirect(url_for('main.index'))
        resp.set_cookie("token", token)

        return resp

    except ValueError as ex:
        return render_template(
            "signup.html",
            error=ex
        )


@auth.route('/logout')
def logout():
    resp = redirect(url_for("main.index"))
    resp.set_cookie("token", "", expires=0)
    return resp
