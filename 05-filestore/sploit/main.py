import sys
import uuid

import blowfish
import requests


def main(url):
    def addr(path):
        return f'{url}/api{path}'

    r = requests.get(addr('/get?filename=flag.enc'))
    encrypted = r.content

    r = requests.get(addr('/owner?filename=flag.enc'))
    admin_id = uuid.UUID(r.json()['response'])

    cipher = blowfish.Cipher(admin_id.bytes, byte_order='little')
    decrypted = b''.join(cipher.decrypt_ecb(encrypted))
    print(decrypted.rstrip(b'\x00').decode())


if __name__ == '__main__':
    main(sys.argv[1])
