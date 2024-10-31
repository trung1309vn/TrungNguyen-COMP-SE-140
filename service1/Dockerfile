FROM python:3.12

WORKDIR /client

ADD client.py /client/

ADD *.html /client/

ADD requirements.txt /client/

EXPOSE 8199

RUN pip install -r requirements.txt

RUN apt-get update && apt-get install -y docker.io

CMD ["python", "client.py"]
