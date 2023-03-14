import { RemovalPolicy, Stack } from "aws-cdk-lib";
import { AttributeType, BillingMode, ProjectionType, StreamViewType, Table, TableEncryption } from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";


export class POCDynamoDBStack extends Stack {
    public readonly pocTestTable: Table;

    constructor(scope: Construct, id: string) {

        super(scope, id);

        this.pocTestTable = this.createTestTable();
    };

    private createTestTable(): Table {
        const table = new Table(this, 'POCTestTableDDB-TiDB', {
            tableName: 'POC_TestTableDDB_TiDB',
            removalPolicy: RemovalPolicy.RETAIN,
            billingMode: BillingMode.PAY_PER_REQUEST,
            encryption: TableEncryption.AWS_MANAGED,
            partitionKey: {
                name: 'orgID',
                type: AttributeType.STRING
            },
            sortKey: {
                name: 'docID',
                type: AttributeType.STRING
            },
            stream: StreamViewType.NEW_IMAGE,
        });

        table.addGlobalSecondaryIndex({
            indexName: "TestGSI",
            partitionKey: {name: 'GSIPK', type: AttributeType.STRING},
            sortKey: {name: 'GSISK', type: AttributeType.STRING},
            projectionType: ProjectionType.ALL,
        })

        return table;
    }


}