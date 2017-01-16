import { Version } from './version';
import { ContainerSpec } from './container-spec';

export class TaskResources {
    Limits: any;
    Reservations: any;
}

export class TaskRestartPolicy {
    Condition: string;
    MaxAttempts: number;
}

export class TaskPlacement {
    Constraints: string[];
}

export class TaskSpec {
    ContainerSpec: ContainerSpec;
    Resources: TaskResources;
    RestartPolicy: TaskRestartPolicy;
    Placement: TaskPlacement;
}

export class TaskContainerStatus {
    ContainerID: string;
}

export class TaskStatus {
    Timestamp: Date;
    State: string;
    Message: string;
    ContainerStatus: TaskContainerStatus;
}

export class TaskDriverConfiguration {
    Name: string;
}

export class TaskDriver {
    Name: string;
}

export class TaskIPAMOptionsConfigs {
    Subnet: string;
    Gateway: string;
}

export class TaskIPAMOptions {
    Driver: TaskDriver;
    Configs?: TaskIPAMOptionsConfigs[];
}

export class TaskNetworkSpec {
    Name: string;
    DriverConfiguration: TaskDriverConfiguration;
    IPAMOptions: TaskIPAMOptions;
}

export class TaskDriverState {
    Name: string;
    Options: { [key:string]:string; };
}

export class TaskNetwork {
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    Spec: TaskNetworkSpec;
    DriverState: TaskDriverState;
    IPAMOptions: TaskIPAMOptions;
}

export class TaskNetworkAttachment {
    Network: TaskNetwork;
    Addresses: string[];
}

export class Task {
    Name?: string;
    ContainerID?: string;
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    Spec: TaskSpec;
    ServiceID: string;
    Slot: number;
    NodeID: string;
    Status: TaskStatus;
    DesiredState: string;
    NetworksAttachments: TaskNetworkAttachment[];
}
