import re
import uuid

import requests as requests

from encrypt import encrypt

FILENAME = 'file.enc'
SUCCESSFUL_RESP = '{"response":"successful"}'


def test_server():
    def addr(path):
        return f'http://localhost:8080/api{path}'

    text = b"123456789abcdefghlanmpst"

    creds_admin = {'username': 'admin1', 'pass': 'admin'}
    with requests.Session() as s:
        r = s.post(addr('/register'), json=creds_admin)
        assert r.text == SUCCESSFUL_RESP

        r = s.post(addr('/logout'))
        assert r.text == SUCCESSFUL_RESP

        r = s.put(addr('/upload'), files={'file': (FILENAME, b'123')})
        assert r.text == '{"response":"unauthorized"}'

        r = s.post(addr('/login'), json=creds_admin)
        assert r.text == SUCCESSFUL_RESP

        r = s.get(addr('/userid'))
        m = re.match(r'{\"response\":\"([a-f0-9\-]{36})\"}', r.text)
        if not m:
            raise "owner id don't match"
        user_id = uuid.UUID(m.groups()[0])

        large_file = bytes(101)
        r = s.put(addr('/upload'), data={'encrypted': 'false'}, files={'file': (FILENAME, large_file)})
        assert r.text == '{"response":"file must be less than 100 bytes"}'

        r = s.put(addr('/upload'), data={'encrypted': 'true'}, files={'file': (FILENAME, encrypt(user_id, text))})
        assert r.text == SUCCESSFUL_RESP

        r = s.put(addr('/upload'), data={'encrypted': 'true'}, files={'file': (FILENAME, encrypt(user_id, text))})
        assert r.text == '{"response":"file already exist"}'

        r = s.get(addr('/list'))
        assert r.text == '{' + f'"response":["{FILENAME}","flag.enc"]' + '}'

        r = s.get(addr(f'/owner?filename={FILENAME}'))
        m = re.match(r'{\"response\":\"([a-f0-9\-]{36})\"}', r.text)
        if not m:
            raise "owner id don't match"
        assert uuid.UUID(m.groups()[0]) == user_id

        r = s.get(addr(f'/get?filename={FILENAME}'))
        assert r.content == text

        r = s.post(addr('/logout'))
        assert r.text == SUCCESSFUL_RESP

    creds_user = {'username': 'user', 'pass': 'user'}
    with requests.Session() as s:
        r = s.post(addr('/register'), json=creds_user)
        assert r.text == SUCCESSFUL_RESP

        r = s.get(addr(f'/get?filename={FILENAME}'))
        assert r.content == encrypt(user_id, text)


if __name__ == '__main__':
    test_server()
