FROM python:3.11-slim

WORKDIR /app

RUN pip install --no-cache-dir kubernetes

COPY app.py .

CMD ["python3", "app.py"]
