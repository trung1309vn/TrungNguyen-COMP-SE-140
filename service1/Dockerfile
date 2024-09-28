FROM python:3.12
WORKDIR /client
ADD client.py /client/
ADD index.html /client/
ADD requirements.txt /client/
RUN pip install -r requirements.txt
CMD ["python", "client.py"]
