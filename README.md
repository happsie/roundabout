# Roundabout

## Usage

1. Create a configuration file using `.yml`, see example configuration as an example
2. Start Roundabout `roundabout start --config=config.yml`

## Example config

```
port: 8090
defaultTargetHost: "google.com"
services: [
		{
			name: "My Service",
			targetHost: "localhost:8080",
			paths: ["/api/my-path-v1", "/api/my-other-v1"]
		},
		{
			name: "My Service",
			targetHost: "localhost:8082",
			paths: ["/api2/my-path-v1", "/api2/my-other-v1"]
		}
	]
```