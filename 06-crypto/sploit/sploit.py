import json
from contextlib import contextmanager
from pwn import process, remote
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad, unpad

ADDR = "localhost", 4552
DEBUG = False
BLOCK_SIZE = 16

@contextmanager
def get_remote():
    if DEBUG:
        r = remote('0.0.0.0', 1337)
    else:
        r = remote(*ADDR)
    yield r
    r.close()


def check_block(block, iv, r):
    data = json.dumps({"option":"check", "text":block.hex(), "iv":iv.hex()})
    r.sendline(bytes(data, encoding = 'utf-8'))
    M = r.recv(1024)
    if b'False' in M:
        return False
    elif b'True' in M:
        return True

def single_block_attack(block, iv,r):
    zeroing_iv = [0] * BLOCK_SIZE
    for pad_val in range(1, BLOCK_SIZE+1):
        padding_iv = [pad_val ^ b for b in zeroing_iv]
        for candidate in range(256):
            padding_iv[-pad_val] = candidate
            iv = bytes(padding_iv)
            if check_block(block,iv,r):
                if pad_val == 1:
                    padding_iv[-2] ^= 1
                    iv = bytes(padding_iv)
                    if not check_block(block,iv,r):
                        continue
                break
        else:
            raise Exception("no valid padding byte found (is the oracle working correctly?)")
        zeroing_iv[-pad_val] = candidate ^ pad_val
    return zeroing_iv


def full_attack(iv, ct, r):
    assert len(iv) == BLOCK_SIZE and len(ct) % BLOCK_SIZE == 0
    msg = iv + ct
    blocks = [msg[i:i+BLOCK_SIZE] for i in range(0, len(msg), BLOCK_SIZE)]
    result = b''
    iv = blocks[0]
    for ct in blocks[1:]:
        dec = single_block_attack(ct, iv, r)
        pt = bytes(iv_byte ^ dec_byte for iv_byte, dec_byte in zip(iv, dec))
        result += pt
        iv = ct
    return result

def get_enc_password(r):
    M = r.recv(1024)
    option = json.dumps({"option":"get_password"})
    r.sendline(bytes(option, encoding = 'utf-8'))
    m = r.recvuntil(b'\n')
    data = json.loads(m)
    return data

def get_flag(r, result):
    password = unpad(result, BLOCK_SIZE)
    data = json.dumps({"option":"verify_password", "password":str(password)})
    r.sendline(bytes(data, encoding = 'utf-8'))
    m = r.recvuntil(b'\n')
    flag = json.loads(m)
    print(flag["flag"])

def main():
    with get_remote() as r:
        data = get_enc_password(r)
        enc_password = bytes.fromhex(data["password"])
        iv = bytes.fromhex(data["iv"])
        result = full_attack(iv, enc_password, r)
        get_flag(r,result)

        
if __name__ == "__main__":
    main()