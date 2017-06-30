#!/usr/bin/env python2

import json, urllib2, hashlib, struct, sha, time


class ChbtcApi:
    def __init__(self, akey, skey):
        self.access_key = akey
        self.secret_key = skey

    def _fill(self, value, length, fill_byte):
        if len(value) >= length:
            return value
        else:
            fill_size = length - len(value)
        return value + chr(fill_byte) * fill_size

    def _do_xor(self, s, value):
        slist = list(s)
        for index in xrange(len(slist)):
            slist[index] = chr(ord(slist[index]) ^ value)
        return "".join(slist)

    def _hmac_sign(self, aValue, aKey):
        keyb = struct.pack("%ds" % len(aKey), aKey)
        value = struct.pack("%ds" % len(aValue), aValue)
        k_ipad = self._do_xor(keyb, 0x36)
        k_opad = self._do_xor(keyb, 0x5c)
        k_ipad = self._fill(k_ipad, 64, 54)
        k_opad = self._fill(k_opad, 64, 92)
        m = hashlib.md5()
        m.update(k_ipad)
        m.update(value)
        dg = m.digest()

        m = hashlib.md5()
        m.update(k_opad)
        sub_str = dg[0:16]
        m.update(sub_str)
        dg = m.hexdigest()
        return dg

    def _digest(self, a_value):
        value = struct.pack("%ds" % len(a_value), a_value)
        # print(value)
        h = sha.new()
        h.update(value)
        dg = h.hexdigest()
        return dg

    def _api_call(self, path, params=''):
        try:
            sha_secret = self._digest(self.secret_key)
            sign = self._hmac_sign(params, sha_secret)
            reqTime = (int)(time.time() * 1000)
            params += '&sign=%s&reqTime=%d' % (sign, reqTime)
            url = 'https://trade.chbtc.com/api/' + path + '?' + params
            request = urllib2.Request(url)
            response = urllib2.urlopen(request, timeout=2)
            doc = json.loads(response.read())
            return doc
        except Exception, ex:
            print('chbtc request ex: ', ex)
            return None

    def query_account(self):
        try:
            params = "method=getAccountInfo&accesskey=" + self.access_key
            path = 'getAccountInfo'

            obj = self._api_call(path, params)
            return obj
        except Exception, ex:
            print('chbtc query_account exception,', ex)
            return None


if __name__ == '__main__':
    access_key = 'your access key'
    secret_key = 'your secret key'

    api = ChbtcApi(access_key, secret_key)

    print(api.query_account())
