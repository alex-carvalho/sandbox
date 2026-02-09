FROM python:3.11-slim

WORKDIR /app

RUN pip install --no-cache-dir kubernetes watchdog

COPY sidecar.py .

USER 1000

CMD ["python3", "sidecar.py"]
