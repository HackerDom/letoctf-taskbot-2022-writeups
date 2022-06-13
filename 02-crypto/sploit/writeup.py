from Crypto.Util.number import inverse, long_to_bytes


def read_params():
    with open("static.txt", 'r') as f:
        N = int(f.readline()[4:])
        e = int(f.readline()[4:])
        c = int(f.readline()[4:])
    return N, e, c

def get_private_key(N,e):
    phi = N - 1
    d = inverse(e,phi)
    return d 

def get_flag(c,d,N):
    return pow(c,d,N)

def solve_task():
    N,e,c = read_params()
    d = get_private_key(N,e)
    flag = long_to_bytes(get_flag(c,d,N))
    print(flag)

def main():
    solve_task()



if __name__ == '__main__':
    main()