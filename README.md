# Important

- all relatives paths in `config.yaml` are converted into absolute paths relative to the directory of `config.yaml`

# Requirements

- **config.yaml**: contains all information about the startup of the server
  -- certificates
  -- log level
  -- server endpoint

# UI

The ui endpoint is hosted on http://localhost:8080

- To change the configuration, of nodes, send an array of matching node and cfg ids:

curl -X POST http://localhost:8888/api/trigger \
 -H "Content-Type: application/json" \
 -d '{
[
{"cfg_id": "config1", "node_id": "node1"},
{"cfg_id": "config2", "node_id": "node2"}
]
}'

# CMD

## Reconfigure kritis3m_scale
with the `import` command, the policies are parsed and stored into `db.sqlite`. To avoid conflicts, delete and recreate db.sqlite
- The server uses the .json file containing the new database state defined in:
  - `config.yaml` -> `acl_policy_path: ./startup.json`
- **import**: `./kritis3m_scale --config <path/to/config.yaml> import`

## list

- **list active nodes**: `./kritis3m_scale --config <path/to/config.yaml> nodes lsa`
- **list nodes**: `./kritis3m_scale --config <path/to/config.yaml> nodes list`
- **list configs**: `./kritis3m_scale --config <path/to/config.yaml> nodes configs list`
- **list configs and appls**: `./kritis3m_scale --config <path/to/config.yaml> nodes configs list --with_appls`

# activate a certain config for a node
- **activate**: `./kritis3m_scale --config <path/to/config.yaml> nodes activate --cfg_id <id> --node_id<node_id>`


