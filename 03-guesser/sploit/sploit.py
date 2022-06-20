import hashlib
import string

def first_part():
    love = bytes("Hackerdom", "utf-8")
    check1 = bytes([4, 4, 23, 4, 38, 38, 34, 20, 30])
    res = []
    for i in range(len(love)):
        res.append(chr(love[i] ^ check1[i]))
    
    return "".join(res)


def second_part():
    check = "0n5b"
    res = []
    for i in range(len(check)):
        res.append(chr(ord(check[i]) - i))

    return "".join(res)


def third_part():
    check = bytes([178, 45, 190, 74, 136, 86, 244, 152, 236, 125, 14, 102, 191, 250, 105, 203])
    for i in string.ascii_lowercase + "_":
        for j in string.ascii_lowercase + "_":
            for k in string.ascii_lowercase + "_":
                if hashlib.md5((i + j + k).encode("utf-8")).digest() == check:
                    return i + j + k


def fourth_part():
    check = -1655835832096201751
    rnd = -1638144296173776243
    return (check ^ rnd).to_bytes(7, "big").decode()


def fifth_part():
    symbols = [chr(x) for x in  [51, 118, 125, 95, 114]]
    res = []
    for i in range(len(symbols)):
        res.append(symbols[(i - 2) % 5])

    return "".join(res)


if __name__ == "__main__":
    print(first_part(), second_part(), third_part(), fourth_part(), fifth_part(), sep="")