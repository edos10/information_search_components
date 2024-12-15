import bisect

from collections import Counter
from typing import SupportsIndex, Iterable, Sized


def suffix_array(s: SupportsIndex | Sized) -> list[int]:
    """
    :param s: строка
    :return: суффиксный массив
    такой навиный алгоритм работает за n^2 + nlogn
    """

    return sorted(range(len(s)), key=lambda i: s[i:])


def bwt(s: SupportsIndex, suff: list[int] = None):
    """
    :param s: входная строка s
    :param suff: суффиксный массив если есть
    :return: строку в результате преобразования Барроуза-Уиллера
    """
    if suff is None:
        suff = suffix_array(s)

    result = [s[i - 1] for i in suff]
    if isinstance(s, str):
        result = ''.join(result)

    return result


def lexicographic_count_symbols(s: SupportsIndex | Iterable) -> dict:
    """
    :param s: принимает строку (пока)
    :return: словарь, где ключ - символ, а его значение - количество символов в тексте, лексикографически меньших, чем он
    """
    counter: Counter = Counter(s)
    res: dict = dict()

    total: int = 0
    for sym in sorted(counter):
        res[sym] = total
        total += counter[sym]

    return res


# could be designed in a more time-efficient way
class SymbolsCountTable:
    def __init__(self, s: SupportsIndex | Iterable):
        self.poses = {}
        for i, char in enumerate(s):
            self.poses.setdefault(char, []).append(i)

    def count(self, char, right_bound: int) -> int:
        if char not in self.poses:
            return 0
        return bisect.bisect_right(self.poses[char], right_bound - 1)
