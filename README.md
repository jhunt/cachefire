Cache Fire
==========

`cachefire` is a Cloud Foundry Firehose Nozzle that caches all
ValueMetric and CounterEvents (a tallying metric), and exposes a
simple, authentication-protected JSON api for retrieving said
cache.

It makes the most sense when integrated with a monitoring system
that will regularly poll it, and wants a full picture of all the
intervening stream.

A Warning To My Friends
-----------------------

**THIS IS ALPHA SOFTWARE**

Currently, `cacehfire` is exploratory.  It definitely needs some
work in the scaling department, and may be missing features
entirely.  Here's what we know we want that hasn't yet been
implemented:

  1. Scalability - `cachefire` is limited to a single instance.
     If you run more than one, they will divide the firehose
     stream amongst them, but never re-integrate the cache.

  2. Metric aging - if a metric stops being emitted through the
     firehose, `cachefire` holds onto it and won't let go.  It
     would be preferable to assign each cache entry a "last seen"
     timestamp and regularly prune the cache accordingly.

Deployment
----------

To deploy `cachefire`, you need the code and a Cloud Foundry.

```
git clone https://github.com/jhunt/cachefire
cd cachefire

# push the code...
cf push --no-start

# now, configure according to your environment.
# see "Configuration", below, for a full  rundown.

# first: a username and a password for accessing the cachefire
# JSON API to retrieve cached metric values.
cf set-env cachefire CACHE_FIRE_USERNAME  cachefire
cf set-env cachefire CACHE_FIRE_PASSWORD  its-a-secret

# next, configure the UAA authentication to use when
# hooking up to the cloud foundry firehose
cf set-env cachefire NOZZLE_UAA_URL       https://uaa.<system-domain>
cf set-env cachefire NOZZLE_UAA_CLIENT    a-client-id
cf set-env cachefire NOZZLE_UAA_SECRET    its-also-a-secret-we-hope

# finally, configure the traffic controller endpoint
# and the firehose subscription parameters.
cf set-env cachefire NOZZLE_TRAFFIC_CONTROLLER_URL wss://doppler.<system-doain>:4443
cf set-env cachefire NOZZLE_SUBSCRIPTION           cachefire-1

# start the app
cf start cachefire

# celebrate!
```

Configuring
-----------

`cachefire` is configured entirely through environment variables.

There are environment variables for controlling the security and
authentication parameters of the JSON API:

- `$CACHE_FIRE_USERNAME` - The HTTP Basic Auth username that must
  be used to access the API.  This is **required**.
- `$CACHE_FIRE_PASSWORD` - The HTTP Basic Auth password that must
  be used to access the API.  This is **required**.

There are environment variables for configuring the firehose / UAA
integration as well:

- `$NOZZLE_UAA_URL` - The full (https) URL of the Cloud Foundry UAA
  endpoint.  This is **required**.
- `$NOZZLE_UAA_CLIENT` - What client ID to authenticate to the UAA as.
  This is **required**.
- `$NOZZLE_UAA_SECRET` - The client secret for the configured UAA client ID.
  This is **required**.
- `$NOZZLE_UAA_SKIP_VERIFY` - Set to "yes" to skip TLS/SSL validation of the
   UAA endpoint.  _Not recommended for production environments!_
- `$NOZZLE_SUBSCRIPTION` - Firehose subscription name.  This is
   **required**, and must be unique to thi cachefire instance.
- `$NOZZLE_TRAFFIC_CONTROLLER_URL` - The web sockets URL of the Cloud
  Foundry Loggregator / Traffic Controller endpoint.  This is **required**.