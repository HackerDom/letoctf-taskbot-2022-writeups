FROM python:3.10

RUN apt update && apt install -y socat

RUN useradd -s /bin/false -m task

WORKDIR /home/task/

ADD requirements.txt ./

RUN pip install -r requirements.txt

ADD server.py secret.py ./

RUN chmod 555 server.py

ENTRYPOINT ["socat", "TCP-LISTEN:31337,reuseaddr,fork", "EXEC:'timeout 300 ./server.py',su=task"]