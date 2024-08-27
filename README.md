# PoC that generates a number of RGW requests to keystone

```sh
./ceph-auth-poc -endpoint http://ceph.rgw.endpoint.local -region region1 -access ec2access -secret ec2secret -requests 1000
```
