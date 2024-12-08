from array import array

import pytest

from positional_index import PositionalIndex, Lang


def test_base_search():
    index = PositionalIndex()
    index.add_document("text a b c d text d text rand")
    index.add_document("text a b c d text d text rand d")
    index.add_document("rand text d", Lang.EN, 1)
    res = index.search("rand text d")

    assert res == array("I", [1])


def test_base_case_not_found_then_found():
    index = PositionalIndex()
    index.add_document("text a b c d text d text rand")
    index.add_document("text a b c d text d text rand d")

    res = index.search("rand text d")
    assert res == array("I", [])

    index.add_document("rand text d", Lang.EN, 1)
    res = index.search("rand text d")

    assert res == array("I", [1])

def test_add_many_documents_and_big_query_on_several_results():
    index = PositionalIndex(stop_words=["и"])
    index.add_document("Один, два. Трое. Четырех. Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Два, Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Один, Пять и шесть, четыре и два, пять и шесть.",lang=Lang.RU)


    res = index.search("четыре два пять шесть", Lang.RU)
    assert res == array("I", [2, 3])

def test_add_many_documents_and_without_stop_word_fail():
    index = PositionalIndex(stop_words=["а"])
    index.add_document("Один, два. Трое. Четырех. Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Два, Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Один, Пять и шесть, четыре и два, пять и шесть.",lang=Lang.RU)


    res = index.search("четыре два пять шесть", Lang.RU)
    assert res == array("I", [])

def test_add_many_documents_and_max_diff_two():
    index = PositionalIndex(max_diff_length=2, stop_words=["а"])
    index.add_document("Один, один. Трое. шесть. Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Два, Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Один, Пять, шесть, четыре и два, пять и шесть.",lang=Lang.RU)


    res = index.search("один один шесть", Lang.RU)
    assert res == array("I", [1, 3])


def test_add_many_documents_and_max_diff_two_with_stop_word_in_text():
    index = PositionalIndex(max_diff_length=2, stop_words=["один"])
    index.add_document("Один, один. Трое. шесть. Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Два, Пять и шесть, семь и восемь, девять и десять.", lang=Lang.RU)
    index.add_document("Три. Один. Четырех. Один, Пять, шесть, четыре и два, пять и шесть.",lang=Lang.RU)


    res = index.search("один один шесть", Lang.RU)
    assert res == array("I", [1, 2, 3])
