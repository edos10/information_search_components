import pymorphy3
morph = pymorphy3.MorphAnalyzer()

text = "Четверо"


def lemmatize(text):
    words = text.split() # разбиваем текст на слова
    res = list()
    for word in words:
        p = morph.parse(word)[0]
        res.append(p)

    return res

def main():
    print(lemmatize(text))


if __name__ == "__main__":
    main()
