import {
    aws_events_targets as event_targets,
    aws_lambda as lambda,
    Duration,
    Stack,
    StackProps,
} from "aws-cdk-lib";
import { Table } from "aws-cdk-lib/aws-dynamodb";
import { DynamoEventSource, SqsDlq, SqsEventSource } from "aws-cdk-lib/aws-lambda-event-sources";
import { Queue, QueueEncryption } from "aws-cdk-lib/aws-sqs";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { Construct } from "constructs";
import { StartingPosition, Tracing } from "aws-cdk-lib/aws-lambda";

export interface LambdaProperties extends StackProps {
    testTable: Table
}

export class LambdaStack extends Stack {
    public readonly receiveShipmentFunc: lambda.Function;
    public readonly searchFunc: lambda.Function;
    
    constructor(scope: Construct, id: string, props: LambdaProperties) {

        super(scope, id, props);
        const queue = this.createFIFOQueue("POC-DynamoDB-Stream-Queue");

        this.receiveShipmentFunc = this.receiveShipment(props.testTable);
        this.sendDDBRecord(queue);
        this.streamReceiver(props.testTable, queue);
        this.searchFunc = this.search();
    };

    private receiveShipment(pocTable: Table,) :GoFunction {
        const func = new GoFunction(this, "POC_ReceiveShipmentFunc", {
            functionName: "poc-receive-shipment-func",
            timeout: Duration.seconds(60),
            entry: "./cmd/receiveShipment",
            tracing: Tracing.PASS_THROUGH,
            architecture: lambda.Architecture.ARM_64,
            environment: {
                POC_TABLE: pocTable.tableName,
            },
            bundling: {
                goBuildFlags: ['-ldflags "-s -w"', "-trimpath"],
            },
        });

        pocTable.grantReadWriteData(func);

        return func;
    }

    private streamReceiver(pocTable: Table, queue: Queue) {
        const streamErrDlq = new Queue(this, "POC_Table_DDB_Error_Stream_Handler.dlq",
        {visibilityTimeout: Duration.seconds(60)},);

        const onFailed = new SqsDlq(streamErrDlq);

        const func = new GoFunction(this, "POC_StreamReceiverFunc", {
            functionName: "poc-stream-receiver-func",
            timeout: Duration.seconds(60),
            entry: "./cmd/streamReceiver",
            tracing: Tracing.PASS_THROUGH,
            architecture: lambda.Architecture.ARM_64,
            environment: {
                QUEUE_URL: queue.queueUrl,
            },
            bundling: {
                goBuildFlags: ['-ldflags "-s -w"', "-trimpath"],
            },
            events: [
                new DynamoEventSource(pocTable, {
                    batchSize: 10,
                    enabled: true,
                    startingPosition: StartingPosition.LATEST,
                    onFailure: onFailed,
                    retryAttempts: 3,
                    reportBatchItemFailures: true,
                    bisectBatchOnError: true,
                }),
            ],
            deadLetterQueue: streamErrDlq,
        });

        queue.grantSendMessages(func);
    }

    private sendDDBRecord(queue: Queue) {
        const dlq = new Queue(this, "POC_SendDDB_Error_Handler.dlq", {visibilityTimeout: Duration.seconds(60)});

        const func = new GoFunction(this, "POC_SendDDB_Record_Func", {
            functionName: "poc-send-ddb-record-func",
            timeout: Duration.seconds(60),
            entry: "./cmd/sendDDBRecord",
            tracing: Tracing.PASS_THROUGH,
            architecture: lambda.Architecture.ARM_64,
            environment: {
                QUEUE_URL: queue.queueUrl,
            },
            bundling: {
                goBuildFlags: ['-ldflags "-s -w"', "-trimpath"],
            },
            events: [
                new SqsEventSource(queue, {
                    batchSize: 10,
                    enabled: true,
                    reportBatchItemFailures: true,
                }),
            ],
            deadLetterQueue: dlq,
        });
    }

    private search() :GoFunction {
        return new GoFunction(this, "POC_Search_Func", {
            functionName: "poc-search-func",
            timeout: Duration.seconds(60),
            entry: "./cmd/search",
            tracing: Tracing.PASS_THROUGH,
            architecture: lambda.Architecture.ARM_64,
            memorySize: 1024,
            bundling: {
                goBuildFlags: ['-ldflags "-s -w"', "-trimpath"],
            },
        });
    }

    private createFIFOQueue(queueName: string): Queue {
        const dlq = new Queue(this, queueName + "-DLQ.fifo", {
            queueName: queueName + "-DLQ.fifo",
            fifo: true,
            encryption: QueueEncryption.SQS_MANAGED,
            visibilityTimeout: Duration.seconds(60),
        });

        const queue = new Queue(this, queueName, {
            queueName: queueName + ".fifo",
            fifo: true,
            contentBasedDeduplication: true,
            visibilityTimeout: Duration.seconds(60),
            deadLetterQueue: {
                queue: dlq,
                maxReceiveCount: 5,
            },
            encryption: QueueEncryption.SQS_MANAGED,
        });

        return queue;
    }

}