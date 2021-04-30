module github.com/pingcap/tidb

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Jeffail/gabs/v2 v2.5.1
	github.com/blacktear23/go-proxyprotocol v0.0.0-20180807104634-af7a81e8dd0d
	github.com/coocood/freecache v1.1.1
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548
	github.com/cznic/sortutil v0.0.0-20181122101858-f5f958428db8
	github.com/danjacques/gofslock v0.0.0-20200623023034-5d0bd0fa6ef0
	github.com/dgraph-io/ristretto v0.0.1
	github.com/dgryski/go-farm v0.0.0-20190423205320-6a90982ecee2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.2-0.20190904063534-ff6b7dc882cf
	github.com/google/btree v1.0.0
	github.com/google/pprof v0.0.0-20200407044318-7d83b28da2e9
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/klauspost/cpuid v1.2.1
	github.com/ngaut/pools v0.0.0-20180318154953-b7bc8c42aac7
	github.com/ngaut/sync2 v0.0.0-20141008032647-7a24ed77b2ef
	github.com/opentracing/basictracer-go v1.0.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/oraluben/go-fuzz v0.0.0-20210430101957-6ebb538fc058
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pingcap/badger v1.5.1-0.20200908111422-2e78ee155d19
	github.com/pingcap/br v5.0.0-nightly.0.20210419090151-03762465b589+incompatible
	github.com/pingcap/check v0.0.0-20200212061837-5e12011dc712
	github.com/pingcap/errors v0.11.5-0.20201126102027-b0a155152ca3
	github.com/pingcap/failpoint v0.0.0-20210316064728-7acb0f0a3dfd
	github.com/pingcap/fn v0.0.0-20200306044125-d5540d389059
	github.com/pingcap/goleveldb v0.0.0-20191226122134-f82aafb29989
	github.com/pingcap/kvproto v0.0.0-20210416062510-5a0d6e96603c
	github.com/pingcap/log v0.0.0-20210317133921-96f4fcab92a4
	github.com/pingcap/parser v0.0.0-20210427084954-8e8ed7927bde
	github.com/pingcap/sysutil v0.0.0-20210315073920-cc0985d983a3
	github.com/pingcap/tidb-tools v4.0.9-0.20201127090955-2707c97b3853+incompatible
	github.com/pingcap/tipb v0.0.0-20210422074242-57dd881b81b1
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.20.0
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/soheilhy/cmux v0.1.4
	github.com/tiancaiamao/appdash v0.0.0-20181126055449-889f96f722a2
	github.com/tikv/pd v1.1.0-beta.0.20210323121136-78679e5e209d
	github.com/twmb/murmur3 v1.1.3
	github.com/uber-go/atomic v1.4.0
	github.com/uber/jaeger-client-go v2.26.0+incompatible
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/atomic v1.7.0
	go.uber.org/automaxprocs v1.2.0
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c
	golang.org/x/text v0.3.6
	golang.org/x/tools v0.1.0
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	sourcegraph.com/sourcegraph/appdash v0.0.0-20190731080439-ebfcffb1b5c0
	sourcegraph.com/sourcegraph/appdash-data v0.0.0-20151005221446-73f23eafcf67
)

go 1.13

replace (
	github.com/oraluben/go-fuzz => ../go-fuzz
	github.com/pragmatwice/go-squirrel => ../go-squirrel
)
