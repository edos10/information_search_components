FROM python:3.10

WORKDIR positional_index

COPY ./ ./positional_index

RUN pip install --no-cache-dir -r ./positional_index/requirements.txt

CMD ["pytest", "./"]