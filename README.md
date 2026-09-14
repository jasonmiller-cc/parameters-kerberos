# parameters-kerberos

REST API for Kerberos KDC management, built on [parameters-core](https://github.com/jasonmiller-cc/parameters-core).

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/principals` | List all principals |
| GET | `/api/v1/principals/{principal}` | Get principal details |
| POST | `/api/v1/principals` | Add a principal |
| DELETE | `/api/v1/principals/{principal}` | Delete a principal |
| POST | `/api/v1/principals/{principal}/password` | Change password |
| POST | `/api/v1/principals/{principal}/keytab` | Generate keytab |
| GET | `/api/v1/policies` | List policies |
| GET | `/api/v1/policies/{name}` | Get policy |
| POST | `/api/v1/policies` | Create policy |
| DELETE | `/api/v1/policies/{name}` | Delete policy |
| GET | `/api/v1/realm` | Get realm info |

## Configuration

```yaml
server:
  port: 8080

kerberos:
  realm: EXAMPLE.COM
  kadmin_server: kdc.example.com
  kadmin_principal: admin/admin@EXAMPLE.COM
  keytab_path: /etc/krb5.keytab
  kdc_host: kdc.example.com
```

Environment variable prefix: `PARAMS_KERBEROS_`.

## Development

```bash
go mod tidy
make build
make test
./dist/parameters-kerberos -config config.yaml
```

## Docker

```bash
make docker
docker run -p 8080:8080 parameters-kerberos:latest
```
