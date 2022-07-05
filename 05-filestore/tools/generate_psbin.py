import struct

import blowfish


def generate_psbin():
    with open('ps.bin', 'wb') as file:
        for p in blowfish.PI_P_ARRAY:
            file.write(struct.pack('I', p))

        for s in blowfish.PI_S_BOXES:
            for b in s:
                file.write(struct.pack('I', b))


if __name__ == '__main__':
    generate_psbin()
