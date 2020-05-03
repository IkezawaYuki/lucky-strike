
from typing import List, Tuple, Optional

def get_pair(numbers: List[int], target: int) -> Optional[Tuple[int, int]]:
    cache = set()
    for num in numbers:
        cache.add(num)
        print("--num--")
        print(num)
        print("--cache--")
        print(cache)
        val = target - num
        print("--val--")
        print(val)
        if val in cache:
            return val, num


def get_pair_half_sum(numbers: List[int]) -> Optional[Tuple[int, int]]:
    sum_numbers = sum(numbers)
    if sum_numbers % 2 != 0:
        return
    half_sum = int(sum_numbers / 2)

    cache = set()

    for num in numbers:
        cache.add(num)
        val = half_sum - num
        if val in cache:
            return val, num

if __name__ == '__main__':
    l = [11,2,5,9,10,3]
    t = 14

    print(get_pair(l,t))
    print("--answer--")

    print(get_pair_half_sum(l))