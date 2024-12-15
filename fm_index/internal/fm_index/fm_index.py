from typing import SupportsIndex, Sized, Reversible

import pkg.helpers as helper_string


class FMIndex:
    def __init__(self, text: str):
        text = text.lower()
        self._text: str = text.lower()

        self._suffix_array = helper_string.suffix_array(text)
        self._bwt: str = helper_string.bwt(text, self._suffix_array)
        self._lexicographic_count:  dict[SupportsIndex, int] = (
            helper_string.lexicographic_count_symbols(text)
        )
        self._count_symbols_table: helper_string.SymbolsCountTable = helper_string.SymbolsCountTable(self._bwt)

    def search(self, pattern: SupportsIndex | Reversible) -> tuple[int, int]:
        start = 0
        end = len(self._bwt) - 1

        for char in reversed(pattern):
            if char not in self._lexicographic_count:
                return -1, -2


            start: int = self._lexicographic_count[char] + self._count_symbols_table.count(char, start)
            end: int = self._lexicographic_count[char] + self._count_symbols_table.count(char, end + 1) - 1
            if start > end:
                return -1, -2

        return start, end

    def generate_results(self, start: int, end: int, n: int, limit: int = -1):
        for i in range(max(start, 0), end + 1):
            if limit == 0:
                break
            limit = max(-1, limit - 1)
            yield (self._suffix_array[i], self._suffix_array[i] + n)

    def search_query(self, pattern: SupportsIndex | Sized):
        start, end = self.search(pattern)
        return self.generate_results(start, end, n=len(pattern))
