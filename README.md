[![Go](https://github.com/sghaida/github-metrics/actions/workflows/go.yml/badge.svg)](https://github.com/sghaida/github-metrics/actions/workflows/go.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/29f8ca49a29e4e8b99bbd61709b5dae6)](https://app.codacy.com/gh/sghaida/github-metrics/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
# GitHub PR metrics
you can use this project to get Repo average PR review time, along bunch of other data
please refer to the generated `metrics.xlsx` for more details.
to see an example of the generated sheets, please refer to [metrics](./data/metrics.xlsx). please note that the sheet is just an example. 

## how to run

refer to [example config](./example.config.yaml) to create `config.yml` under the same place

### running from container

please note that the generated output by default will be under `/tmp/github-metrics/metrics.xlsx`
as the container mounting `/tmp/github-metrics` from the host, if you wish to change that please update [Dockerfile](./Dockerfile) and [Makefile](./Makefile) accordingly


```shell
make docker-build
make docker-run
```

### running locally
please note the following command line arguments 
* **out**: define the output directory of the generated Excel sheet, default to `/tmp`
* **from**: from date that would be used to read the contributions from, default `begining of the current month`
* **to**: to date that would be used to read the contributions from, default `end of the current month`

#### to run
```shell
go build -v -o github-metrics
./github-metrics -out=/tmp -from=2023-01-01 -to=2023-05-01
```

