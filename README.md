> This project was archived and moved into [o11yutil](https://github.com/wperron/o11yutil).

# 🧟 Zombie

A _simple_ natural load generator meant to simulate organic traffic on a system.
It's a rudimentary tool meant to use locally to test and application or simply
generate irregular traffic patterns more similar to a human than a simple 
infinite loop.

## Example

example config:

```yaml
api:
  enabled: true
  addr: ":8082"

targets:
  - name: localhost
    url: "https://google.com" # Required
    delay: 10000              # 10,000ms, or 10s
    jitter: 0.2
    headers:
      "Accept":
        - "*/*"
```

run with:

```bash
zombie -config zombie.yaml
```
