package initialize

import (
	"strconv"

	"github.com/goat-project/goat-one/resource"
	"github.com/onego-project/onego/resources"

	"github.com/goat-project/goat-one/constants"
	"github.com/goat-project/goat-one/reader"

	log "github.com/sirupsen/logrus"
)

type benchmark struct {
	bType  string
	bValue string
}

// UserTemplateIdentity returns map of user ID and value in TEMPLATE/IDENTITY.
func UserTemplateIdentity(r reader.Reader) map[int]string {
	objs, err := r.ListAllUsers()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all users")
		return nil
	}

	res := make([]resource.Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return find(res, constants.TemplateIdentity)
}

// ImageTemplateCloudkeeperApplianceMpuri returns map of image ID and
// value in TEMPLATE/CLOUDKEEPER_APPLIANCE_MPURI.
func ImageTemplateCloudkeeperApplianceMpuri(r reader.Reader) map[int]string {
	objs, err := r.ListAllImages()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all images")
		return nil
	}

	res := make([]resource.Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return find(res, constants.TemplateCloudkeeperApplianceMpuri)
}

// HostTemplateBenchmark returns two maps: map of host ID and value in TEMPLATE/BENCHMARK_TYPE and
// map of host ID and value in TEMPLATE/BENCHMARK_VALUE.
func HostTemplateBenchmark(r reader.Reader) (map[int]string, map[int]string) {
	hosts, err := r.ListAllHosts()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error list all hosts")
		return nil, nil
	}

	clustersMap := clustersMap(r)

	hostLength := len(hosts)
	hostTemplateBenchmarkType := make(map[int]string, hostLength)
	hostTemplateBenchmarkValue := make(map[int]string, hostLength)

	for _, host := range hosts {
		id, err := host.ID()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("error get host ID")
			continue
		}

		bType, err := host.Attribute(constants.TemplateBenchmarkType)
		if err != nil {
			bType = typeFromCluster(clustersMap, host)
		}

		hostTemplateBenchmarkType[id] = bType

		bValue, err := host.Attribute(constants.TemplateBenchmarkValue)
		if err != nil {
			bValue = valueFromCluster(clustersMap, host)
		}

		hostTemplateBenchmarkValue[id] = bValue
	}

	return hostTemplateBenchmarkType, hostTemplateBenchmarkValue
}

func find(res []resource.Resource, constant string) map[int]string {
	m := make(map[int]string, len(res))

	for _, r := range res {
		id, err := r.ID()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("error get ID")
			continue
		}

		str, err := r.Attribute(constant)
		if err != nil {
			str = strconv.Itoa(id)
		}

		m[id] = str
	}

	return m
}

func valueFromCluster(clustersMap map[int]benchmark, host *resources.Host) string {
	clusterID, err := host.Cluster()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get cluster ID from host")
		return ""
	}

	return clustersMap[clusterID].bValue
}

func typeFromCluster(clustersMap map[int]benchmark, host *resources.Host) string {
	clusterID, err := host.Cluster()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get cluster ID from host")
		return ""
	}

	return clustersMap[clusterID].bType
}

func clustersMap(r reader.Reader) map[int]benchmark {
	clusters, err := r.ListAllClusters()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error list all clusters")
	}

	idToBenchmark := make(map[int]benchmark, len(clusters))

	for _, cluster := range clusters {
		id, err := cluster.ID()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("error get cluster ID")
			continue
		}

		bType, err := cluster.Attribute(constants.TemplateBenchmarkType)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "cluster": id}).Warn("couldn't get benchmark type from cluster")
		}

		bValue, err := cluster.Attribute(constants.TemplateBenchmarkValue)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "cluster": id}).Warn("couldn't get benchmark value from cluster")
		}

		idToBenchmark[id] = benchmark{bType: bType, bValue: bValue}
	}

	return idToBenchmark
}
