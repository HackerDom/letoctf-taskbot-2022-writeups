#!/usr/bin/env python
import os
from Crypto.Cipher import AES
from secret import key, flag, password
from Crypto.Util.Padding import pad, unpad
import json

BLOCK_LENGTH = 16

def check_pad(s, iv):
    cipher = AES.new(key, AES.MODE_CBC, iv)
    t = cipher.decrypt(s)
    try:
        unpad(t, BLOCK_LENGTH)
        return True
    except:
        return False

def get_password():
    iv = os.urandom(BLOCK_LENGTH)
    cipher = AES.new(key, AES.MODE_CBC, iv)
    enc_passwd = cipher.encrypt(pad(password,BLOCK_LENGTH))
    return json.dumps({"password": enc_passwd.hex(), "iv":iv.hex()})

def option_check(text, iv):
    try:
        return check_pad(bytes.fromhex(text), bytes.fromhex(iv))
    except:
        return "not hex format"

def verify_password(your_pass):
    if(your_pass == str(password)):
        data = json.dumps({"flag":flag})
        return data
    else:
        return "incorrect password"

def header():
    print("Hello guys! I stole your flag and will give it only if you send me correct password\n"
        "But I'm too kind today and you can get it in encrypted form\n"
        "Also, you can try to input your cipher and check it padding. Your options:\n")

def print_options():
    print("get_passwd - Get encrypted password in hex\n"
        "check - Check padding of ciphertext in JSON format in hex encoding\n"
        "verify_passwd - Enter password and get flag\n")

def main():
    header()
    while(True):
        print_options()
        try:
            choice = json.loads(input())
            if choice["option"] == "get_password":
                server_answer = get_password()
                print(server_answer)
            elif choice["option"] == "check":
                server_answer = option_check(choice["text"], choice["iv"])
                print(server_answer)
            elif choice["option"] == "verify_password":
                server_answer = verify_password(choice["password"])
                print(server_answer)
            else:
                print("Bad input")
        except:
            print("Incorrect JSON")


if __name__ == "__main__":
    main()