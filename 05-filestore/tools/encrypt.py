import blowfish


def encrypt(user_id, text):
    cipher = blowfish.Cipher(user_id.bytes, byte_order='little')
    return b"".join(cipher.encrypt_ecb(text))
