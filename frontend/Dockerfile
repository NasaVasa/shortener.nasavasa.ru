FROM python:3.8-slim-buster

WORKDIR /py

COPY frontend/requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY frontend/ .

CMD [ "python3", "-m" , "flask", "--app", "server", "run", "--host=0.0.0.0"]