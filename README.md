# mqttstore

Store MQTT messages into MySQL database

## Usage

Run it from console:

```bash
$ mqttstore -dsn testuser:testpass1234@/mqtt_test -host iot.eclipse.org -topic test/#
```

The command above will try to connect to local database, then try to connect to MQTT broker and print stats about stored messages:

```
2018/05/07 20:42:17 db.connected: testuser:testpass1234@/mqtt_test
2018/05/07 20:42:17 mqtt.server: iot.eclipse.org:1883 (dd44fdee-90d3-427c-9438-a9a319589fb7)
2018/05/07 20:42:18 mqtt.topic: test/#
2018/05/07 20:42:18 mqtt.client: ignoring_retained_messages (this can take a long time depending on your broker stored retained messages)
2018/05/07 20:42:28 mqtt.stats.received_messages_total: 3
2018/05/07 20:42:38 mqtt.stats.received_messages_total: 6
2018/05/07 20:42:48 mqtt.stats.received_messages_total: 9
```

### Options

For complete set of options, use the `-help` flag:

```bash
$ mqttstore -help

Usage of mqttstore:

  -clientid string
        MQTT Client ID (defaults to random UUID4)

  -dsn string
        Data source name (e.g.: user:password@hostname/database_name) (default "root:root@/mqtt")

  -host string
        MQTT server hostname or IP address (default "iot.eclipse.org")

  -port string
        MQTT server port (default "1883")

  -topic string
        MQTT topic (default "#")
```

## Installation

You will need to have *golang* installed and the run:

```
go install github.com/atarantini/mqttstore
```

The binary `mqttstore` should be in your `bin` directory.

A MySQL database with a `message` table is also required:

```sql
CREATE TABLE message
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    topic VARCHAR(256),
    payload TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## CHANGELOG

#### 0.0.1 - 2018/05/07

- Initial release

## License

Released under MIT license, see https://opensource.org/licenses/MIT

## Author

Andres I. Tarantini (atarantini@gmail.com)