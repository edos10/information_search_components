import sys
import time
import logging
from logging import Logger

from internal.fm_index.fm_index import FMIndex


def read_file(path: str) -> str:
    try:
        with open(path, "r", encoding="utf-8") as f:
            return f.read()
    except FileNotFoundError:
        print(f"Невозможно прочитать данный файл, попробуйте еще раз...")
        sys.exit(1)


def prepare_logger() -> logging.Logger:
    logger: logging.Logger = logging.getLogger(__name__)

    handler: logging.FileHandler = logging.FileHandler("index.log", mode='w')
    handler_out: logging.StreamHandler = logging.StreamHandler(sys.stdout)
    formatter: logging.Formatter = logging.Formatter("[%(asctime)s] LEVEL: %(levelname)s | %(message)s")

    handler.setFormatter(formatter)
    handler_out.setFormatter(formatter)

    logger.addHandler(handler)
    logger.addHandler(handler_out)

    return logger


def main():
    logger: logging.Logger = prepare_logger()

    if len(sys.argv) < 2:
        logger.log(logging.ERROR, "Usage: python3 main.py <filepath>")
        sys.exit(1)

    file_path = sys.argv[1]

    logger.log(logging.INFO, f"Построение индекса из файла: {file_path}...")
    start_time: time = time.time()
    index: FMIndex = FMIndex(read_file(file_path))
    logger.log(logging.INFO,f"Индекс готов. Это заняло {time.time() - start_time:.2f} секунд")

    while True:
        try:
            pattern = input("Ищем паттерн: ").strip()

            start_time: time = time.time()
            start, end = index.search(pattern)
            total_time: time = time.time() - start_time

            logger.log(logging.INFO,f"Найдено {max(0, end - start + 1)} соответствий")
            if end - start + 1 <= 0:
                logger.log(logging.WARN, f"Ничего не найдено по этому запросу... {pattern}")
            else:
                for start, end in index.generate_results(start, end, n=len(pattern)):
                    print(f"Позиции в файле: [{start}; {end - 1}]")
            logger.log(logging.INFO,f"Для обработки данного запроса ушло {total_time:.2f} секунд")

        except KeyboardInterrupt:
            logger.log(logging.INFO, "Работа индекса окочена...")
            break

if __name__ == "__main__":
    main()
