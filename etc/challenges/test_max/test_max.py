import pytest
import answer


@pytest.mark.parametrize('numbers, expected', [
    ((), None),
    (["-0001", "1010101", "0001", "-10101010101", "1111111"], 1111111),
    (["foo", "15"], 15),
])
def test_answer_correct(numbers, expected):
    assert answer.get_max(numbers) == expected
