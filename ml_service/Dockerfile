FROM python:3.12-slim

WORKDIR /app

# Системные зависимости для numpy/scipy/torch
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc g++ && \
    rm -rf /var/lib/apt/lists/*

# Копируем зависимости первыми (кэш Docker)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Копируем только код (без модели — она скачается при старте)
COPY server.py loader.py ./

EXPOSE 8000

# Запуск сервера
CMD ["uvicorn", "server:app", "--host", "0.0.0.0", "--port", "8000"]