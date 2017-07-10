import json
import os
import time
import requests
from datetime import datetime

end = 3903900
uri = "https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=your_api_key"
alert = False
timeout = 0.2
while True:
    try:
        r = requests.get(uri)
        latest = int(json.loads(r.text)["result"], 16)
    except Exception, ex:
        print "get error" + str(ex)

    print("BLOCK: %s, DIFF: %s" % (latest, end-latest))
    time.sleep(timeout)
