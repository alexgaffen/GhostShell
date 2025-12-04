# brain/Dockerfile
FROM python:3.9-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the code
COPY . .

# Run the app
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5000"]