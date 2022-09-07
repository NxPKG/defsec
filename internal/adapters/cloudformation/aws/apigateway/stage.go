package apigateway

import (
	v2 "github.com/aquasecurity/defsec/pkg/providers/aws/apigateway/v2"
	parser2 "github.com/aquasecurity/defsec/pkg/scanners/aws/cloudformation/parser"
	"github.com/aquasecurity/defsec/pkg/types"
)

func getApis(cfFile parser2.FileContext) (apis []v2.API) {

	apiResources := cfFile.GetResourcesByType("AWS::ApiGatewayV2::Api")
	for _, apiRes := range apiResources {
		api := v2.API{
			Metadata:     apiRes.Metadata(),
			Name:         types.StringDefault("", apiRes.Metadata()),
			ProtocolType: types.StringDefault("", apiRes.Metadata()),
			Stages:       getStages(apiRes.ID(), cfFile),
		}
		apis = append(apis, api)
	}

	return apis
}

func getStages(apiId string, cfFile parser2.FileContext) []v2.Stage {
	var apiStages []v2.Stage

	stageResources := cfFile.GetResourcesByType("AWS::ApiGatewayV2::Stage")
	for _, r := range stageResources {
		stageApiId := r.GetStringProperty("ApiId")
		if stageApiId.Value() != apiId {
			continue
		}

		s := v2.Stage{
			Metadata:      r.Metadata(),
			Name:          r.GetStringProperty("StageName"),
			AccessLogging: getAccessLogging(r),
		}
		apiStages = append(apiStages, s)
	}

	return apiStages
}

func getAccessLogging(r *parser2.Resource) v2.AccessLogging {

	loggingProp := r.GetProperty("AccessLogSettings")
	if loggingProp.IsNil() {
		return v2.AccessLogging{
			Metadata:              r.Metadata(),
			CloudwatchLogGroupARN: types.StringDefault("", r.Metadata()),
		}
	}

	destinationProp := r.GetProperty("AccessLogSettings.DestinationArn")

	if destinationProp.IsNil() {
		return v2.AccessLogging{
			Metadata:              loggingProp.Metadata(),
			CloudwatchLogGroupARN: types.StringDefault("", r.Metadata()),
		}
	}
	return v2.AccessLogging{
		Metadata:              destinationProp.Metadata(),
		CloudwatchLogGroupARN: destinationProp.AsStringValue(),
	}
}
