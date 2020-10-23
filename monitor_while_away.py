# usage: monitor_while_away.py

import requests, sys, syslog, time
from is_home import *

URL = MOTION_URL_HERE
SLEEP = 5 # seconds
AWAY_THRESHOLD = 300 # seconds

def main():
        while True:
                try:
                        time.sleep(SLEEP) # give it a few seconds to come up
                        monitor()
                except Exception as e:
                        syslog.syslog(syslog.LOG_ERR, f'problem with monitor_while_away.py err={e}')

def monitor():
        # motion comes up active and gets set to pause as soon as one person is home.
        # this guarantees it stays active as long as no one is home.
        last_home_time = time.time() - (AWAY_THRESHOLD + 1)

        while True:
                if is_home(["anyone"]):
                        last_home_time = time.time()
                syslog.syslog(syslog.LOG_DEBUG, f"the last time someone was home was {time.strftime('%D %H:%M:%S', time.localtime(last_home_time))}")
                active = motion_active()
                syslog.syslog(syslog.LOG_DEBUG, f"motion active={active}")

                if recently_home(last_home_time) and active:
                        syslog.syslog(syslog.LOG_INFO, "someone is home, pausing motion")
                        requests.get(URL+"pause")
                elif not recently_home(last_home_time) and not active:
                        syslog.syslog(syslog.LOG_INFO, "no one is home, starting motion")
                        requests.get(URL+"start")

                time.sleep(SLEEP)

def recently_home(last_home_time):
        return (time.time() - last_home_time) < AWAY_THRESHOLD

def motion_active():
        resp = requests.get(url=URL+"status")
        syslog.syslog(syslog.LOG_DEBUG, f"motion responded to status request with {resp.text}")
        if "Detection status ACTIVE" in resp.text:
                return True
        return False

if __name__ == '__main__':
        main()
