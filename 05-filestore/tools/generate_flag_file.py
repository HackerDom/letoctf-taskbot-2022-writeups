import sys
import uuid

from encrypt import encrypt


# 9c97249a-26d8-49f2-8f1b-92641a8901f8
def generate_flag_file(user_id, flag):
    if len(flag) % 8 != 0:
        flag = flag.encode() + b'\x00' * (8 - (len(flag) % 8))
    b = encrypt(uuid.UUID(user_id), flag)
    with open('flag.enc', 'wb') as file:
        file.write(b)


if __name__ == "__main__":
    generate_flag_file(sys.argv[1], sys.argv[2])
