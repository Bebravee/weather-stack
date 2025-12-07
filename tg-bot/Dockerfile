FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

# secured_data.json монтируем с хоста, чтобы не класть ключи в образ
CMD ["python", "bot.py"]
