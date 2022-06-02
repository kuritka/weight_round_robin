# wrr DEMO

Weight round robin alghoritm produces indexes by predefined weights.

See PDF [30 40 20 10]: The item with the highest probability (index 01) will occur more often at the 01 position that has the highest probability in the PDF
```txt
    [10.0.0.1],[10.1.0.1],[10.2.0.1],[10.3.0.1]
    [30 40 20 10]
    -----------------
 0. [289 401 200 110] 
 1. [298 315 258 129] 
 2. [291 216 307 186] 
 3. [122 68 235 575] 
```
See what DNS query returns :

```shell
# dig @localhost -p 1053 wrr.cloud.example.com +short
10.1.0.1
10.0.0.1
10.3.0.1
10.2.0.1
```


The example matrix was created by 1000x hitting the list of IP adresses with help of WRR. 
The IP at zero index (10.0.0.1) is used 289x on the first position returned by DNS server.
Also 298x used on the second position returned by DNS server.

The address (10.3.0.1) has only 10% probability of to be chosen. It occurs only 110x (cca 10%) on the zero position 
while 575x on the last position. 

The index was calculated 1000 times. When you add individual columns or rows, the result is always 1000, so everything is OK

# WRR From Query To k8gb 

## DNS Query
```shell
# dig @localhost -p 1053 roundrobin.cloud.example.com +tcp +nostats +noedns +nocomments
;roundrobin.cloud.example.com.  IN      A
roundrobin.cloud.example.com. 30 IN     A       172.18.0.8
roundrobin.cloud.example.com. 30 IN     A       172.18.0.5
roundrobin.cloud.example.com. 30 IN     A       172.18.0.6
roundrobin.cloud.example.com. 30 IN     A       172.18.0.9
```

```shell
# dig amazon.com  +nostats +noedns +nocomments
;amazon.com.                    IN      A
amazon.com.             47      IN      A       176.32.103.205
amazon.com.             47      IN      A       54.239.28.85
amazon.com.             47      IN      A       205.251.242.103
```

```shell
# dig app.cloud.example.com +nostats +noedns +nocomments
;app.cloud.example.com.                    IN      A
app.cloud.example.com.             300      IN      A       10.0.0.2
app.cloud.example.com.             300      IN      A       10.1.1.2
app.cloud.example.com.             300      IN      A       10.1.1.1
app.cloud.example.com.             300      IN      A       10.10.10.10
app.cloud.example.com.             300      IN      A       10.0.0.1
```

## CoreDNS side

```shell
# example configuration. For real k8gb config see: 
# https://github.com/k8gb-io/k8gb/blob/master/chart/k8gb/templates/coredns-cm.yaml
app.cloud.example.com:8053 {
    hosts etchosts
    log
    loadbalance weight_round_robin
}
# 1 - 0.8 - 0.2 = 0 everything out of range is divided by a 
# probability of 0%

amazon.com {
    hosts etchosts
    log
    loadbalance weight_round_robin {
    }
}
# 1 - 0.6 = 0.4 everything out of 176.32.103.205 is divided by a
# probability 40% (54.239.28.85 = 20%, 205.251.242.103 = 20%)
```

## k8gb side
```yaml
# Cluster 1
# 10.0.0.1
# 10.0.0.2
apiVersion: ohmyglb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: app-gslb
  namespace: test-gslb
spec:
  ingress:
    rules:
      - host: app.cloud.example.com
        http:
          paths:
            - backend:
                serviceName: app
                servicePort: http
              path: /
  strategy: roundRobin 
    weight: 80%
```

```yaml
# Cluster 2
# 10.1.1.1
# 10.1.1.2
apiVersion: ohmyglb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: app-gslb
  namespace: test-gslb
spec:
  ingress:
    rules:
      - host: app.cloud.example.com
        http:
          paths:
            - backend:
                serviceName: app
                servicePort: http
              path: /
  strategy: roundRobin 
    weight: 20%
```

```yaml
# Cluster 3
# 10.10.10.10
apiVersion: ohmyglb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: app-gslb
  namespace: test-gslb
spec:
  ingress:
    rules:
      - host: app.cloud.example.com
        http:
          paths:
            - backend:
                serviceName: app
                servicePort: http
              path: /
  strategy: roundRobin 
 # not sure what to set here atm. 

# Expected result
The 1000x executed `dig app.cloud.example.com` would return

```
# ± 800x 
10.0.0.1
10.0.0.2
10.1.1.1
10.1.1.2
10.10.10.10


# ±200x 
10.1.1.1
10.1.1.2
10.0.0.1
10.0.0.2
10.10.10.10
```

_There are several corner-cases, getting out of 0% probability is one of them (10.10.10.10). At this point I put the address at the end of the list, but it can be completely discarded (just create some switch to leave the 0% address or discard it)._  