import socket
from flask import Flask, request, render_template
import psutil
import shutil
from datetime import datetime
import requests
import json
import sys
import time
import os

start_time=time.time()

client = Flask(__name__, template_folder="/client")

@client.route('/')
def home():
    return render_template("home.html")

@client.route('/disconnect')
def disconnect():
    os.system("docker stop $(docker ps -q)")
    os.system("docker rm $(docker ps -aq)")
    return ""

@client.route('/access', methods=['GET'])
def access():
    global start_time
    while (time.time() - start_time <= 2.0):
        pass

    start_time = time.time()
    hostname = socket.gethostname()
    IPAddr = socket.gethostbyname(hostname)
    processes = [[proc.pid, proc.name()] for proc in psutil.process_iter()]
    (_, _, free_space) = shutil.disk_usage("/")
    last_reboot = psutil.boot_time()
    last_reboot = datetime.fromtimestamp(last_reboot)

    response = requests.get('http://server:8080/?passkey=test')
    res = response.json()
    return render_template("index.html", \
        server_IPAddress = res["IP"],
        server_Processes = res["Processes"],
        server_DiskUsage = res["DiskUsage"] / 1024**3,
        server_LastBootTime = res["LastBootTime"],
        client_IPAddress = IPAddr,
        client_Processes = processes,
        client_DiskUsage = free_space / 1024**3,
        client_LastBootTime = last_reboot,
    )

if __name__ == "__main__":
    client.run(host="0.0.0.0", port=8199)
