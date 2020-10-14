#!/usr/bin/python3
# usage: is_home.py [name ..] [group]
# Determine if a person is home based on certain network devices connected to the router.
# ex:
# python3 is_home.py cody
# python3 is_home.py julia
# python3 is_home.py cody julia
# python3 is_home.py anyone
# python3 is_home.py everyone

import subprocess, sys

people = {
    "cody": [{
        "ip": "XXX.XXX.XXX.XXX",
        "mac": "XX:XX:XX:XX:XX:XX"
    }],
    "julia": [{
        "ip": "XXX.XXX.XXX.XXX",
        "mac": "XX:XX:XX:XX:XX:XX"
    }]
}

def usage():
    print('<name>')
    print('   Manually specify which people should be home.')
    print('   Can use multiple times.')
    print('   Any person specified will satisfy the condition (e.g. if Cody is home but Julia is not, will return true).')
    print('   Valid names: %s', valid_names)
    print('<group>')
    print('   Use instead of specifying people.')
    print('   Can only specify one group.')
    print('   Valid groups: %s', valid_groups)

def main():
    if len(sys.argv) == 1:
        usage()
        sys.exit(2) # error
    if is_home(sys.argv[1:]):
        sys.exit(0) # true
    sys.exit(1) # false

def is_home(args):
    for arg in args:
        if arg == "anyone":
            for name in people:
                if are_any_devices_connected(people[name]):
                    return True # someone is home

        elif arg == "everyone":
            for name in people:
                if not are_any_devices_connected(people[name]):
                    return False # not everyone is home
            return True # everyone is home

        else:
            if are_any_devices_connected(people[arg]):
                return True # individual names only require one to be home

    return False # no one is home

def are_any_devices_connected(devices):
    for device in devices:
        call = subprocess.run(['sudo', 'nmap', '-sn', device["ip"]], stdout=subprocess.PIPE, text=True)
        if device["mac"] in str(call):
            return True
    return False

if __name__ == '__main__':
    main()
