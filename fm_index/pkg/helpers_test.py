import unittest

import pkg.helpers as helpers


class TestHelpers(unittest.TestCase):
    def test_suffix_array(self):
        cases = [
            "banana",
            "ababa",
        ]

        ans = [
            [5, 3, 1, 0, 4, 2],
            [4, 2, 0, 3, 1],
        ]

        for i, case in enumerate(cases):
            self.assertEqual(
                helpers.suffix_array(case), ans[i]
            )

    def test_bwt(self):
        self.assertEqual(
            helpers.bwt("some string"),
            "emnroistg s"
        )

        self.assertEqual(
            helpers.bwt("aabbb"),
            "babba",
        )

    def test_lexicographic_cnt_less(self):
        cases = [
            "randdd",
            "aabbcdde",
        ]

        ans = [
            {"a": 0, "d": 1, "n": 4,"r": 5},
            {"a": 0, "b": 2, "c": 4, "d": 5, "e": 7},
        ]

        for i, case in enumerate(cases):
            cnt = helpers.lexicographic_count_symbols(case)
            for char in set(list(case)):
                true_cnt = ans[i][char]
                self.assertEqual(cnt[char], true_cnt)

    def test_occurance_table(self):
        occ = helpers.SymbolsCountTable("xabacafad")
        self.assertEqual(occ.count("a", 4), 2)
        self.assertEqual(occ.count("a", 3), 1)
        self.assertEqual(occ.count("f", 4), 0)
        self.assertEqual(occ.count("f", 8), 1)
