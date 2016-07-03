#!/usr/bin/env python

import sys
import tempfile

def main():
    if len(sys.argv) != 3:
        print 'bad args %s' % sys.argv
        return 1

    filename, amount = sys.argv[1:]
    amount = float(amount)
    new_contents = []
    with open(filename) as f:
        lines = f.readlines()
        for line in lines:
            line = line.strip()
            if line.startswith('v '):
                split = line.split(' ')
                parsed = [str(float(x)*amount) for x in split[1:]]
                line = ' '.join(split[:1] + parsed)
            new_contents.append(line)

    with open(filename, 'w') as f:
        f.write('\n'.join(new_contents))


if __name__ == '__main__':
    sys.exit(main())

