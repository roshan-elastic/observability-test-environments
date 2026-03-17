FROM python:3-alpine

COPY policies.py /user/bin/policies.py
COPY requirements.txt /user/share/requirements.txt
COPY policies.yml /etc/elastic-agent/policies.yml

RUN pip install -r /user/share/requirements.txt

ENTRYPOINT ["python", "/user/bin/policies.py"]
CMD [ "--filename", "/etc/elastic-agent/policies.yml" ]
