import socket
from flask import Flask, request, render_template
import psutil
import shutil
from datetime import datetime
import requests
import json
import sys

server = Flask(__name__, template_folder="/client")

@server.route('/', methods=['GET'])
def home():
    hostname = socket.gethostname()
    IPAddr = socket.gethostbyname(hostname)
    processes = [[proc.pid, proc.name()] for proc in psutil.process_iter()]
    (_, _, free_space) = shutil.disk_usage("/")
    last_reboot = psutil.boot_time()
    last_reboot = datetime.fromtimestamp(last_reboot)

    response = requests.get('http://localhost:8080/?passkey=test')
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
    server.run(host="0.0.0.0", port=8199)
