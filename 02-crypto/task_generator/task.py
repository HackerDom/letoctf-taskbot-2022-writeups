from Crypto.Util.number import getStrongPrime, bytes_to_long
import sys


def generate_params():
    N = getStrongPrime(1024)
    e = 65537
    return N, e

def encode(flag,e,N):
    return pow(flag,e,N)

def cout_params(c,e,N):
    with open("static.txt", 'w') as f:
        f.write(f"N = {N}\n")
        f.write(f"e = {e}\n")
        f.write(f"c = {c}\n")

def generate_task(flag):
    N, e = generate_params()
    c = encode(flag,e,N)
    task = cout_params(c,e,N)


def main():
    flag = bytes(sys.argv[1], encoding = 'utf-8')
    generate_task(bytes_to_long(flag))
    

if __name__ == '__main__':
    main()