FROM golang:1.12-alpine

ARG branch=master
ARG version

ENV name="goat-one" \
    user="goat"
ENV project="/go/src/github.com/goat-project/${name}/" \
    homeDir="/var/lib/${user}/" \
    logDir="/var/${name}/log/"

LABEL application=${name} \
      description="Exporting OpenNebula accounting data" \
      maintainer="svetlovska@cesnet.cz" \
      version=${version} \
      branch=${branch}

# Install tools required for project
RUN apk add --no-cache git shadow
RUN go get github.com/golang/dep/cmd/dep

# List project dependencies with Gopkg.toml and Gopkg.lock
COPY Gopkg.lock Gopkg.toml ${project}
# Install library dependencies
WORKDIR ${project}
RUN dep ensure -vendor-only

# Create user and log directory
RUN useradd --system --shell /bin/false --home ${homeDir} --create-home --uid 1000 ${user} && \
    usermod -L ${user} && \
    mkdir -p ${logDir} && \
    chown -R ${user}:${user} ${logDir}

# Copy the entire project and build it
COPY . ${project}
RUN go build -o /bin/${name}

# Switch user
USER ${user}

# Run main command with subcommands and options:
# No subcommand runs goat-one, configures it from config file (goat-one.yml)
# and extracts virtual machine, network and storage data in the same time.
# To extract only specific data, use subcommand:
#   network     Extract network data
#   storage     Extract storage data
#   vm          Extract virtual machine data
#   help        Help about any command
#
# The configuration from file should be rewrite using the following flags:
#  -d, --debug string                 debug
#  -e, --endpoint string              goat server [GOAT_SERVER_ENDPOINT] (required)
#  -h, --help                         help for goat-one
#  -i, --identifier string            goat identifier [IDENTIFIER] (required)
#      --log-path string              path to log file
#  -o, --opennebula-endpoint string   OpenNebula endpoint [OPENNEBULA_ENDPOINT] (required)
#  -s, --opennebula-secret string     OpenNebula secret [OPENNEBULA_SECRET] (required)
#      --opennebula-timeout string    timeout for OpenNebula calls [TIMEOUT_FOR_OPENNEBULA_CALLS] (required)
#  -p, --records-for-period string    records for period [TIME PERIOD]
#  -f, --records-from string          records from [TIME]
#  -t, --records-to string            records to [TIME]
#      --version                      version for goat-one
#
# Example:
# - extract virtual machine data from the last 5 years and save it with idetifier 'goat-vm'
# CMD /bin/goat-one vm --log-path=${logDir}${name}.log -p=5y -i=goat-vm
CMD /bin/goat-one vm --log-path=${logDir}${name}.log
