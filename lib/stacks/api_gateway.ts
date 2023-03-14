import {aws_lambda, Stack, StackProps} from "aws-cdk-lib";
import {aws_apigateway as apig } from "aws-cdk-lib";
import { aws_lambda as lambda } from "aws-cdk-lib";
import {Construct} from "constructs";


export interface APIProperties extends StackProps {
    receiveShipmentFunc: lambda.Function;
    searchFunc: lambda.Function;
    stageName: string;
}

export class APIStack extends Stack {

    constructor(scope: Construct, id: string, props: APIProperties) {
        super(scope, id, props);

        const api = new apig.RestApi(this, "POC_DDB_TiDB_API", {
            restApiName: "POC-DDB-TiDB-API",
            deploy: true,
            apiKeySourceType: apig.ApiKeySourceType.HEADER,
            deployOptions: {
                tracingEnabled: true,
                dataTraceEnabled: true,
                stageName: props.stageName,
                loggingLevel: apig.MethodLoggingLevel.ERROR,
            },
            cloudWatchRole: true,
        });

        const receiver_path = api.root.addResource("shipments");
        receiver_path.addMethod("POST", new apig.LambdaIntegration(props.receiveShipmentFunc, {
            requestParameters: {
                "integration.request.header.ORGID": "method.request.header.ORGID"
            }
        }), {
            apiKeyRequired: true,
            requestParameters: {
                "method.request.header.ORGID": true
            }
        });

        const search_path = api.root.addResource("search");
        receiver_path.addMethod("GET", new apig.LambdaIntegration(props.searchFunc, {
            requestParameters: {
                "integration.request.header.ORGID": "method.request.header.ORGID"
            }
        }), {
            apiKeyRequired: true,
            requestParameters: {
                "method.request.header.ORGID": true
            }
        });

    }

}