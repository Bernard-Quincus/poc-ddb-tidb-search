import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { POCDynamoDBStack } from '../stacks/dynamodb';
import { LambdaStack } from '../stacks/lambda';
import { APIStack } from '../stacks/api_gateway';

export class PocDdbTidbSearchStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const ddbStatck = new POCDynamoDBStack(scope, "POC-DDB-TiDB-Stack");

    const lambdas = new LambdaStack(scope, "POC-DDB-TiDB-Lambdas", {
      testTable: ddbStatck.pocTestTable,
    });
    lambdas.node.addDependency(ddbStatck);

    new APIStack(scope, "POC-DDB-TiDB-API",{
      stageName: "POC-DDB-TiDB-DevStage",
      searchFunc: lambdas.searchFunc,
      receiveShipmentFunc: lambdas.receiveShipmentFunc,
    });
    
  }
}
