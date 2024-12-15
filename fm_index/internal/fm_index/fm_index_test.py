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
