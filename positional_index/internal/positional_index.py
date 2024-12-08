from array import array
import array

import nltk
import pymorphy3
from nltk import PorterStemmer

from nltk.stem import StemmerI
from pyroaring import BitMap
from typing import Dict, List, Set
from enum import Enum
from pymystem3 import Mystem


nltk.download("punkt")
nltk.download("punkt_tab")
nltk.download('stopwords')

EMPTY_BITMAP: BitMap = BitMap()

custom_stop_words = ["и", "а", "но"]

class Lang(Enum):
    RU = "russian"
    EN = "english"

class PositionalIndex:

    def __init__(self, stop_words=None, max_diff_length=1):
        if stop_words is None:
            self._stop_words: Set[str] = set()
        else:
            self._stop_words = stop_words

        self._current_number_id: int = 1
        self._max_diff_length: int = max_diff_length

        self._word_doc: Dict[str, BitMap] = dict()
        self._word_doc_to_positions: Dict[tuple[str, int], BitMap] = dict()

        self.stemmers: Dict[Lang, StemmerI] = {
            Lang.EN: PorterStemmer(),
        }

        self._count_words_doc: Dict[int, int] = dict()
        self.lemmatizer: pymorphy3.MorphAnalyzer = pymorphy3.MorphAnalyzer()

    @staticmethod
    def replacer_symbols(text: str) -> str:
        for sym in ["?", "!", ";", ":", ".", ","]:
            text = text.replace(sym, "", -1)
        return text

    @staticmethod
    def is_word(word: str) -> bool:
        return word.isalpha()

    def processing_text(self, text: str, lang: Lang = Lang.EN) -> list[str]:
        if len(self._stop_words) == 0:
            self._stop_words = custom_stop_words

        text = self.replacer_symbols(text)
        words = text.split(" ")

        words_clean: List[str] = []
        for word in words:
            if word.lower() not in self._stop_words:
                p = self.lemmatizer.parse(word)[0]
                word = p.normal_form
                words_clean.append(word)
        return words_clean

    def add_document(self, text: str, lang: Lang = Lang.EN, num: int = -1) -> None:
        lst_words: list[str] = self.processing_text(text, lang=lang)
        current_doc_id: int = self._current_number_id
        if num != -1:
            current_doc_id = num
        else:
            self._current_number_id += 1

        count_words_doc: int = self._count_words_doc.get(current_doc_id, 0)

        for index, word in enumerate(lst_words):
            self._word_doc.setdefault(word, BitMap()).add(current_doc_id)
            self._word_doc_to_positions.setdefault((word, current_doc_id), BitMap()).add(index + count_words_doc)

        self._count_words_doc[current_doc_id] = count_words_doc + len(lst_words)


    def search(self, body: str, lang: Lang = Lang.EN) -> array:
        words: List[str] = self.processing_text(body, lang)

        if not words:
            return array.array("I")

        bitmap_for_check_positions: BitMap = self._word_doc[words[0]].copy()
        for word in words[1:]:
            if word not in self._word_doc:
                return array.array("I")
            bitmap_for_check_positions &= self._word_doc[word]

        ans_bitmap: BitMap = BitMap()
        for id_document in bitmap_for_check_positions:
            word_positions: List[List[int]] = \
                list(sorted(self._word_doc_to_positions[(word, id_document)]) for word in words)

            if self.check_positions_array(word_positions):
                ans_bitmap.add(id_document)

        return ans_bitmap.to_array()


    def check_positions_array(self, words_positions: List[List[int]]) -> bool:
        pointers = [0] * len(words_positions)
        positions_lengths = [len(pos) for pos in words_positions]

        while pointers[0] < positions_lengths[0]:
            current_pos = words_positions[0][pointers[0]]

            is_exists: bool = True
            for i in range(1, len(words_positions)):
                while pointers[i] < positions_lengths[i] and words_positions[i][pointers[i]] < current_pos:
                    pointers[i] += 1
                if pointers[i] == positions_lengths[i]:
                    return False

                next_pos = words_positions[i][pointers[i]]
                if next_pos - current_pos > self._max_diff_length:
                    is_exists = False
                current_pos = next_pos

            if is_exists:
                return True

            pointers[0] += 1

        return False
