package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
	configs "github.com/erritis/cdk8skit/v3/cdk8skit/configs"
	deployments "github.com/erritis/cdk8skit/v3/cdk8skit/deployments"
	networks "github.com/erritis/cdk8skit/v3/cdk8skit/networks"
	statefulsets "github.com/erritis/cdk8skit/v3/cdk8skit/statefulsets"
	storages "github.com/erritis/cdk8skit/v3/cdk8skit/storages"
	volumes "github.com/erritis/cdk8skit/v3/cdk8skit/volumes"
)

type DbChartProps struct {
	cdk8s.ChartProps
	Environment string
	Network     string
	StorageName string
}

type ServerChartProps struct {
	cdk8s.ChartProps
	Environment   string
	Network       string
	ClusterIssuer string
}

type ClientChartProps struct {
	cdk8s.ChartProps
	Network       string
	ClusterIssuer string
}

type NetworkChartProps struct {
	cdk8s.ChartProps
	Network string
}

func NewDbChart(scope constructs.Construct, id string, props *DbChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	if props.Environment == "Production" {
		storages.NewLocalStorage(chart, props.StorageName, &storages.LocalStorageProps{})
		lpvr := volumes.NewLocalVolume(
			chart,
			"persistent-volume",
			jsii.String("/mnt/vocascandb"),
			&volumes.LocalVolumeProps{
				StorageClassName: jsii.String(props.StorageName),
			},
		)

		statefulsets.NewPostgres(
			chart,
			id,
			&statefulsets.PostgresProps{
				Image:            jsii.String("postgres:12.9"),
				PrefixSecretName: jsii.String("vocascan-db"),
				DBConfig: &statefulsets.DBConfig{
					Name:     jsii.String("{{ .Values.Db.Name }}"),
					Username: jsii.String("{{ .Values.Db.Username }}"),
					Password: jsii.String("{{ .Values.Db.Password }}"),
				},
				VolumeConfig: &statefulsets.VolumeConfig{
					StorageClassName: jsii.String(props.StorageName),
					Volume:           &lpvr.Volume,
				},
				Network: jsii.String(props.Network),
			},
		)
	}

	if props.Environment == "Development" {
		statefulsets.NewPostgres(
			chart,
			id,
			&statefulsets.PostgresProps{
				Image:            jsii.String("postgres:12.9"),
				PrefixSecretName: jsii.String("vocascan-db"),
				DBConfig: &statefulsets.DBConfig{
					Name:     jsii.String("{{ .Values.Db.Name }}"),
					Username: jsii.String("{{ .Values.Db.Username }}"),
					Password: jsii.String("{{ .Values.Db.Password }}"),
				},
				VolumeConfig: &statefulsets.VolumeConfig{
					StorageClassName: jsii.String(props.StorageName),
				},
				Network: jsii.String(props.Network),
			},
		)
	}

	return chart
}

func NewServerChart(scope constructs.Construct, id string, props *ServerChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	volume := volumes.NewSecretVolume(
		chart,
		"vocascan-config-secret",
		jsii.String("vocascan-config"),
		jsii.String("{{ werf_secret_file \"vocascan.config.js\" | b64enc }}"),
		&volumes.SecretVolumeProps{
			Encrypt: jsii.Bool(false),
		},
	)

	deployments.NewFrontend(
		chart,
		id,
		jsii.String("{{ .Values.Server.Domain }}"),
		jsii.String("vocascan/server:latest"),
		&deployments.FrontendProps{
			PortConfig: &configs.ServicePortConfig{
				ContainerPort: jsii.Number(8000),
			},
			Variables: &map[*string]*string{
				jsii.String("VOCASCAN_CONFIG"):          jsii.String("{{ .Values.Server.Config }}"),
				jsii.String("VOCASCAN__DATABASE__HOST"): jsii.String("{{ .Values.Server.Db.Host }}"),
				jsii.String("VOCASCAN__DATABASE__PORT"): jsii.String("{{ .Values.Server.Db.Port }}"),
			},
			Volumes: &map[*string]*cdk8splus26.Volume{
				jsii.String("/etc/vocascan"): &volume,
			},
			ClusterIssuer: &props.ClusterIssuer,
			Network:       jsii.String(props.Network),
		},
	)

	return chart
}

func NewClientChart(scope constructs.Construct, id string, props *ClientChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	deployments.NewFrontend(
		chart,
		id,
		jsii.String("{{ .Values.Client.Domain }}"),
		jsii.String("vocascan/frontend:latest"),
		&deployments.FrontendProps{
			PortConfig: &configs.ServicePortConfig{
				ContainerPort: jsii.Number(80),
			},
			Variables: &map[*string]*string{
				jsii.String("VOCASCAN_BASE_URL"): jsii.String("{{ .Values.Client.BackendUrl }}"),
			},
			ClusterIssuer: &props.ClusterIssuer,
			Network:       jsii.String(props.Network),
		},
	)

	return chart
}

func NewNetworkChart(scope constructs.Construct, id string, props *NetworkChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	networks.NewNetworkPolicy(chart, id, props.Network)

	return chart
}

func main() {

	config, err := LoadConfig()

	if err != nil {
		fmt.Println(err)
	}

	app := cdk8s.NewApp(&cdk8s.AppProps{
		Outdir:              jsii.String(config.Outdir),
		OutputFileExtension: jsii.String(".yaml"),
		YamlOutputType:      cdk8s.YamlOutputType_FILE_PER_CHART,
	})

	cprops := cdk8s.ChartProps{
		DisableResourceNameHashes: jsii.Bool(true),
		Namespace:                 jsii.String("vocascan"),
	}

	NewDbChart(app, "vocascan-db", &DbChartProps{
		ChartProps:  cprops,
		Network:     "io.network/vocascan-network",
		Environment: config.Environment,
		StorageName: config.StorageName,
	})

	NewServerChart(app, "vocascan-server", &ServerChartProps{
		ChartProps:    cprops,
		Network:       "io.network/vocascan-network",
		Environment:   config.Environment,
		ClusterIssuer: config.ClusterIssuer,
	})
	NewClientChart(app, "vocascan-client", &ClientChartProps{
		ChartProps:    cprops,
		Network:       "io.network/vocascan-network",
		ClusterIssuer: config.ClusterIssuer,
	})

	NewNetworkChart(app, "vocascan-network", &NetworkChartProps{
		ChartProps: cprops,
		Network:    "io.network/vocascan-network",
	})

	app.Synth()
}
