import pytest

from internal.fm_index.fm_index import FMIndex


def test_simple_search():
    index = FMIndex("simple text for search")
    lst = list(index.search_query("text"))
    lst.sort(key=lambda i: i[0])
    assert lst == [(7, 11)]


def test_simple_search_with_two_occurencies():
    index = FMIndex("text a b c d text d text rand")
    lst = list(index.search_query("d"))
    lst.sort(key=lambda i: i[0])
    assert lst == [(11, 12), (18, 19), (28, 29)]


def test_on_zero_res():
    index = FMIndex("simple text for search text")
    lst = list(index.search_query("aaa"))
    lst.sort(key=lambda i: i[0])
    assert lst == []

def test_on_hard_search():
    index = FMIndex("""
    С точки зрения банальной эрудиции каждый индивидуум, критически
    мотивирующий абстракцию, не может игнорировать текст из 100 слов,
    концептуально интерпретируя общепринятые дефанизирующие поляризаторы,
    поэтому консенсус, достигнутый диалектической материальной классификацией
    всеобщих мотиваций в парадогматических связях предикатов, решает проблему
    усовершенствования формирующих геотрансплантационных квазипузлистатов
    всех кинетически кореллирующих аспектов. Эта проблема существует уже давно.
    """)
    lst = list(index.search_query("ааа"))
    lst.sort(key=lambda i: i[0])
    assert lst == []

    lst = list(index.search_query("100"))
    lst.sort(key=lambda i: i[0])
    assert lst == [(129, 132)]

    lst = list(index.search_query("ирующих"))
    lst.sort(key=lambda i: i[0])
    assert lst == [(396, 403), (470, 477)]

    lst = list(index.search_query(" ирующих"))
    lst.sort(key=lambda i: i[0])
    assert lst == []

    lst = list(index.search_query("проблем"))
    lst.sort(key=lambda i: i[0])
    assert lst == [(360, 367), (492, 499)]
